/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"fmt"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
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

	diags := req.Config.Get(ctx, &state)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sps, err := helper.GetAllStoragePools(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex Storage Groups",
			err.Error(),
		)
		return
	}

	// Set state for filters
	if state.SPFilter != nil {
		filtered, err := helper.GetDataSourceByValue(*state.SPFilter, sps)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering storage pools: %v please validate the filter", state.SPFilter), err.Error(),
			)
			return
		}
		filteredSps := []scaleiotypes.StoragePool{}
		for _, val := range filtered {
			filteredSps = append(filteredSps, val.(scaleiotypes.StoragePool))
		}
		sps = filteredSps
	}

	for _, val := range sps {
		sp := goscaleio.NewStoragePool(d.client)
		sp.StoragePool = &val
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
		storagePool := helper.GetStoragePoolState(volList, sdsList, val)
		state.StoragePools = append(state.StoragePools, storagePool)
	}

	// this is required for acceptance testing
	state.ID = types.StringValue("storage_pool_datasource")

	// Set state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
