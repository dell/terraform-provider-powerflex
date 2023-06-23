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
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ datasource.DataSource              = &volumeDataSource{}
	_ datasource.DataSourceWithConfigure = &volumeDataSource{}
)

// VolumeDataSource returns the volume data source
func VolumeDataSource() datasource.DataSource {
	return &volumeDataSource{}
}

type volumeDataSource struct {
	client *goscaleio.Client
}

func (d *volumeDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

func (d *volumeDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = VolumeDataSourceSchema
}

func (d *volumeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, _ *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	d.client = req.ProviderData.(*goscaleio.Client)
}

func (d *volumeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.VolumeDataSourceModel
	var volumes []*scaleiotypes.Volume
	var err error

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//Read the volumes based on volume id/name or storage pool id/name and if nothing
	//is mentioned , then return all volumes
	if state.Name.ValueString() != "" {
		volumes, err = d.client.GetVolume("", "", "", state.Name.ValueString(), false)
	} else if state.ID.ValueString() != "" {
		volumes, err = d.client.GetVolume("", state.ID.ValueString(), "", "", false)
	} else if state.StoragePoolID.ValueString() != "" {
		sps, err1 := d.client.FindStoragePool(state.StoragePoolID.ValueString(), "", "", "")
		if err1 != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Powerflex Volumes",
				err1.Error(),
			)
			return
		}
		sp := goscaleio.NewStoragePool(d.client)
		sp.StoragePool = sps
		volumes, err = sp.GetVolume("", "", "", "", false)
	} else if state.StoragePoolName.ValueString() != "" {
		sps, err1 := d.client.FindStoragePool("", state.StoragePoolName.ValueString(), "", "")
		if err1 != nil {
			resp.Diagnostics.AddError(
				"Unable to Read Powerflex Volumes",
				err1.Error(),
			)
			return
		}
		sp := goscaleio.NewStoragePool(d.client)
		sp.StoragePool = sps
		volumes, err = sp.GetVolume("", "", "", "", false)
	} else {
		volumes, err = d.client.GetVolume("", "", "", "", false)
	}
	//check if there is any error while getting the volume
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex Volumes",
			err.Error(),
		)
		return
	}
	state.Volumes = helper.UpdateVolumeState(volumes)
	state.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
