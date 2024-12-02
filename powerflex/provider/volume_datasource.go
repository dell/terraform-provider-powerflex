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

func (d *volumeDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	d.client = req.ProviderData.(*powerflexProvider).client
}

func (d *volumeDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var state models.VolumeDataSourceModel
	var volumes []scaleiotypes.Volume
	var err error

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get all the volumes
	volumes, err = helper.GetAllVolumes(d.client)

	//check if there is any error while getting the volume
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex Volumes",
			err.Error(),
		)
		return
	}
	// Set state for filters
	if state.VolumeFilter != nil {
		filtered, err := helper.GetDataSourceByValue(*state.VolumeFilter, volumes)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Error in filtering volumes: %v please validate the filter", state.VolumeFilter), err.Error(),
			)
			return
		}
		filteredVol := []scaleiotypes.Volume{}
		for _, val := range filtered {
			filteredVol = append(filteredVol, val.(scaleiotypes.Volume))
		}
		volumes = filteredVol
	}
	state.Volumes = helper.UpdateVolumeState(volumes)
	state.ID = types.StringValue("volume_datasource_id")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
