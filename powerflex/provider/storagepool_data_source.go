/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://mozilla.org/MPL/2.0/


Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"context"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ datasource.DataSource              = &storagepoolDataSource{}
	_ datasource.DataSourceWithConfigure = &storagepoolDataSource{}
)

// StoragePoolDataSource is a helper function to simplify the provider implementation.
func StoragePoolDataSource() datasource.DataSource {
	return &storagepoolDataSource{}
}

// storagepoolDataSource is the data source implementation.
type storagepoolDataSource struct {
	client *goscaleio.Client
}

// Metadata returns the data source type name.
func (d *storagepoolDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_pool"
}

// Schema defines the schema for the data source.
func (d *storagepoolDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DataSourceSchema
}

// Configure adds the provider configured client to the data source.
func (d *storagepoolDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	d.client = req.ProviderData.(*powerflexProvider).client
}

// Read refreshes the Terraform state with the latest data.
func (d *storagepoolDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started storage pool data source read method")
	var state models.StoragepoolDataSourceModel
	var pd *scaleio_types.ProtectionDomain
	var err3 error

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the systems on the PowerFlex cluster
	c2, err := helper.GetFirstSystem(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance",
			err.Error(),
		)
		return
	}

	// Check if protection domain ID or name is provided
	if state.ProtectionDomainID.ValueString() != "" {
		pd, err3 = c2.FindProtectionDomain(state.ProtectionDomainID.ValueString(), "", "")
	} else {
		pd, err3 = c2.FindProtectionDomain("", state.ProtectionDomainName.ValueString(), "")
	}

	if err3 != nil {
		resp.Diagnostics.AddError(
			"Unable to find protection domain",
			err3.Error(),
		)
		return
	}

	p1 := goscaleio.NewProtectionDomainEx(d.client, pd)

	sp := goscaleio.NewStoragePool(d.client)

	spID := []string{}
	// Check if storage pool ID or name is provided
	if !state.StoragePoolIDs.IsNull() {
		diags = state.StoragePoolIDs.ElementsAs(ctx, &spID, true)
	} else if !state.StoragePoolNames.IsNull() {
		diags = state.StoragePoolNames.ElementsAs(ctx, &spID, true)
	} else {
		// Get all the storage pools associated with protection domain
		storagePools, _ := p1.GetStoragePool("")
		for sp := range storagePools {
			spID = append(spID, storagePools[sp].Name)
		}
	}

	if numSP := len(spID); numSP == 0 {
		resp.Diagnostics.AddError("No storage pools found for the specified protection domain", "")
		return
	}

	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
		return
	}

	for _, spIdentifier := range spID {
		var s1 *scaleio_types.StoragePool

		if !state.StoragePoolIDs.IsNull() {
			s1, err3 = p1.FindStoragePool(spIdentifier, "", "")
		} else {
			s1, err3 = p1.FindStoragePool("", spIdentifier, "")
		}

		if err3 != nil {
			resp.Diagnostics.AddError(
				"Unable to read storage pool",
				err3.Error(),
			)
			return
		}
		sp.StoragePool = s1

		volList, err4 := sp.GetVolume("", "", "", "", false)
		if err4 != nil {
			resp.Diagnostics.AddError(
				"Unable to get volumes associated with storage pool",
				err4.Error(),
			)
			return
		}

		sdsList, err5 := sp.GetSDSStoragePool()
		if err5 != nil {
			resp.Diagnostics.AddError(
				"Unable to get SDS associated with storage pool",
				err5.Error(),
			)
			return
		}

		storagePool := helper.GetStoragePoolState(volList, sdsList, s1)
		state.StoragePools = append(state.StoragePools, storagePool)
	}

	// this is required for acceptance testing
	state.ID = types.StringValue("dummyID")

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
