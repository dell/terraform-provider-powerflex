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
	"fmt"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"
)

var (
	_ datasource.DataSource              = &deviceDataSource{}
	_ datasource.DataSourceWithConfigure = &deviceDataSource{}
)

// DeviceDataSource returns the volume data source
func DeviceDataSource() datasource.DataSource {
	return &deviceDataSource{}
}

type deviceDataSource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (d *deviceDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (d *deviceDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = DeviceDataSourceSchema
}

func (d *deviceDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	d.client = req.ProviderData.(*powerflexProvider).client

	system, err := helper.GetFirstSystem(d.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}
	d.system = system
}

func (d *deviceDataSource) ValidateConfig(ctx context.Context, req datasource.ValidateConfigRequest, resp *datasource.ValidateConfigResponse) {
	var config models.DeviceDataSourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !config.StoragePoolName.IsNull() {
		if config.ProtectionDomainID.IsNull() && config.ProtectionDomainName.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("storage_pool_name"),
				"Please provide protection_domain_name or protection_domain_id with storage_pool_name.",
				"Please provide protection_domain_name or protection_domain_id with storage_pool_name.",
			)
		}
	}

	if !config.ProtectionDomainID.IsNull() && config.StoragePoolName.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("protection_domain_id"),
			"Please provide protection_domain_id with storage_pool_name.",
			"Please provide protection_domain_id with storage_pool_name.",
		)
	}

	if !config.ProtectionDomainName.IsNull() && config.StoragePoolName.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("protection_domain_name"),
			"Please provide protection_domain_name with storage_pool_name.",
			"Please provide protection_domain_name with storage_pool_name.",
		)
	}
}

func (d *deviceDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	var (
		state models.DeviceDataSourceModel
		// pd      *goscaleio.ProtectionDomain
		err     error
		devices []scaleiotypes.Device
	)

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	if !state.StoragePoolID.IsNull() || !state.StoragePoolName.IsNull() {
		// Get ProtectionDomain with ID and Name if StoragePoolName is provided
		var sp *goscaleio.StoragePool
		var err error
		if !state.StoragePoolName.IsNull() {
			pd, err := helper.GetNewProtectionDomainEx(d.client, state.ProtectionDomainID.ValueString(), state.ProtectionDomainName.ValueString(), "")
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in getting protection domain details with ID: "+state.ProtectionDomainID.ValueString()+" name: "+state.ProtectionDomainName.ValueString(),
					err.Error(),
				)
				return
			}

			// Find StoragePool by Name
			sp, err := pd.FindStoragePool("", state.StoragePoolName.ValueString(), "")
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in getting storage pool details with name: "+state.StoragePoolName.ValueString(),
					err.Error(),
				)
				return
			}

			state.StoragePoolID = types.StringValue(sp.ID)
		}
		// Get StoragePool by ID
		if !state.StoragePoolID.IsUnknown() {
			sp, err = getStoragePool(d, state.StoragePoolID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error getting storage pool instance with ID: "+state.StoragePoolID.ValueString(),
					"unexpected error: "+err.Error(),
				)
				return
			}
		}
		state.StoragePoolName = types.StringValue(sp.StoragePool.Name)
		devices, err = sp.GetDevice()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting all devices within storage pool instance with ID: "+state.StoragePoolID.ValueString(),
				"unexpected error: "+err.Error(),
			)
		}

	} else if !state.SdsID.IsNull() || !state.SdsName.IsNull() {
		var rsp scaleiotypes.Sds
		var err error
		if !state.SdsName.IsNull() {
			sds, err := d.system.FindSds("Name", state.SdsName.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error in getting sds details with name: "+state.SdsName.ValueString(),
					err.Error(),
				)
				return
			}
			state.SdsID = types.StringValue(sds.ID)
		}
		if !state.SdsID.IsUnknown() {
			rsp, err = d.system.GetSdsByID(state.SdsID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Could not get SDS by ID %s", state.SdsID.ValueString()),
					err.Error(),
				)
				return
			}
		}
		sds := goscaleio.NewSdsEx(d.client, &rsp)
		devices, err = sds.GetDevice()
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not get devices within SDS by ID %s", state.SdsID.ValueString()),
				err.Error(),
			)
		}
	} else if !state.CurrentPath.IsNull() {
		devices, err = d.system.GetDeviceByField("DeviceCurrentPathName", state.CurrentPath.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting device with Current Path: "+state.CurrentPath.ValueString(),
				"unexpected error: "+err.Error(),
			)
			return
		}
	} else if !state.Name.IsNull() {
		devices, err = d.system.GetDeviceByField("Name", state.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting device with Name: "+state.Name.ValueString(),
				"unexpected error: "+err.Error(),
			)
			return
		}
	} else if !state.ID.IsNull() {
		devices = make([]scaleiotypes.Device, 0)
		deviceResponse, err3 := d.system.GetDevice(state.ID.ValueString())
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error getting device with ID: "+state.ID.ValueString(),
				"unexpected error: "+err3.Error(),
			)
			return
		}
		devices = append(devices, *deviceResponse)
	} else {
		devices, err = d.system.GetAllDevice()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting all devices in the system ",
				"unexpected error: "+err.Error(),
			)
			return
		}
	}

	// Set refreshed state
	if state.ID.IsNull() {
		state.ID = types.StringValue("placeholder")
	}
	state.DeviceModel = helper.GetAllDeviceState(devices)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func getStoragePool(d *deviceDataSource, storagePoolID string) (*goscaleio.StoragePool, error) {

	system, err := helper.GetFirstSystem(d.client)
	if err != nil {
		return nil, err
	}

	sp, err := system.GetStoragePoolByID(storagePoolID)
	if err != nil {
		return nil, err
	}

	sp1 := goscaleio.NewStoragePoolEx(d.client, sp)
	return sp1, nil
}
