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
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
)

var (
	_ resource.Resource              = &protectionDomainResource{}
	_ resource.ResourceWithConfigure = &protectionDomainResource{}
)

// NewProtectionDomainResource returns the resource for protection domain
func NewProtectionDomainResource() resource.Resource {
	return &protectionDomainResource{}
}

type protectionDomainResource struct {
	client   *goscaleio.Client
	system   *goscaleio.System
	pdClient *goscaleio.ProtectionDomain
}

func (d *protectionDomainResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_protection_domain"
}

func (d *protectionDomainResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ProtectionDomainResourceSchema
}

func (d *protectionDomainResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data models.ProtectionDomainResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// IOPS Limits validation
	if k := data.OverallIoNetworkThrottlingInKbps.ValueInt64(); k != 0 {
		if v := data.ProtectedMaintenanceModeNetworkThrottlingInKbps.ValueInt64(); v == 0 || v > k {
			resp.Diagnostics.AddAttributeError(
				path.Root("protected_maintenance_mode_network_throttling_in_kbps"),
				"protected_maintenance_mode_network_throttling_in_kbps must be set to a value less than overall_io_network_throttling_in_kbps",
				"",
			)
		}
		if v := data.RebalanceNetworkThrottlingInKbps.ValueInt64(); v == 0 || v > k {
			resp.Diagnostics.AddAttributeError(
				path.Root("rebalance_network_throttling_in_kbps"),
				"rebalance_network_throttling_in_kbps must be set to a value less than overall_io_network_throttling_in_kbps",
				"",
			)
		}
		if v := data.RebuildNetworkThrottlingInKbps.ValueInt64(); v == 0 || v > k {
			resp.Diagnostics.AddAttributeError(
				path.Root("rebuild_network_throttling_in_kbps"),
				"rebuild_network_throttling_in_kbps must be set to a value less than overall_io_network_throttling_in_kbps",
				"",
			)
		}
		if v := data.VTreeMigrationNetworkThrottlingInKbps.ValueInt64(); v == 0 || v > k {
			resp.Diagnostics.AddAttributeError(
				path.Root("vtree_migration_network_throttling_in_kbps"),
				"vtree_migration_network_throttling_in_kbps must be set to a value less than overall_io_network_throttling_in_kbps",
				"",
			)
		}
	}

	// RF cache validation
	if !data.RfCacheEnabled.ValueBool() {
		if !data.RfCacheOperationalMode.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("rf_cache_operational_mode"),
				"rf_cache_operational_mode can be set only when rf_cache_enabled is set to true",
				"",
			)
		}
		if !data.RfCachePageSizeKb.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("rf_cache_page_size_kb"),
				"rf_cache_page_size_kb can be set only when rf_cache_enabled is set to true",
				"",
			)
		}
		if !data.RfCacheMaxIoSizeKb.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("rf_cache_max_io_size_kb"),
				"rf_cache_max_io_size_kb can be set only when rf_cache_enabled is set to true",
				"",
			)
		}
	}

	// FGL Metadata caching validation
	if !data.FglMetadataCacheEnabled.ValueBool() {
		if !data.FglDefaultMetadataCacheSize.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("fgl_default_metadata_cache_size"),
				"fgl_default_metadata_cache_size can be set only when fgl_metadata_cache_enabled is set to true",
				"",
			)
		}
	}
}

func (d *protectionDomainResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			"Unable to Read Powerflex System",
			err.Error(),
		)
		return
	}
	d.system = system
}

// ConfigurePdState receives the previous state and builds the protection domain client internally
func (d *protectionDomainResource) ConfigurePdState(ctx context.Context, state models.ProtectionDomainResourceModel) diag.Diagnostics {
	// for now it only reinstates the links and id fields
	d.pdClient = goscaleio.NewProtectionDomain(d.client)
	d.pdClient.ProtectionDomain.ID = state.ID.ValueString()
	var diags diag.Diagnostics
	d.pdClient.ProtectionDomain.Links, diags = helper.GetLinksFromTfList(ctx, state.Links)
	return diags
}

func (d *protectionDomainResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.ProtectionDomainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch protection domain of given id
	resp.Diagnostics.Append(d.ConfigurePdState(ctx, state)...)
	newState, err := d.ReadByID()
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Unable to Read Powerflex ProtectionDomain of ID %s", state.ID.ValueString()),
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

func (d *protectionDomainResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.ProtectionDomainResourceModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	id, err := d.system.CreateProtectionDomain(plan.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating protection domain",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(d.ConfigurePdState(ctx, models.ProtectionDomainResourceModel{
		ID: types.StringValue(id),
	})...)
	if resp.Diagnostics.HasError() {
		return
	}
	tflog.Info(ctx, "Protection domain created with name "+plan.Name.ValueString()+" and id "+id)

	state, err := d.ReadByID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching protection domain by id after initial create.",
			err.Error(),
		)
		return
	}

	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting update of parameters of protection domain "+plan.Name.ValueString()+" after initial successful create.")
	resp.Diagnostics.Append(d.UpdateResource(ctx, plan, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err = d.ReadByID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching protection domain by id after full create.",
			err.Error(),
		)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

func (d *protectionDomainResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan models.ProtectionDomainResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Retrieve values from state
	var oldState models.ProtectionDomainResourceModel
	diags = req.State.Get(ctx, &oldState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.ConfigurePdState(ctx, oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting update of protection domain "+plan.Name.ValueString())
	resp.Diagnostics.Append(d.UpdateResource(ctx, plan, oldState)...)
	if resp.Diagnostics.HasError() {
		return
	}

	state, err := d.ReadByID()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error fetching protection domain by id after full create.",
			err.Error(),
		)
	}
	resp.Diagnostics.Append(resp.State.Set(ctx, state)...)
}

// Delete deletes the resource and removes the Terraform state on success.
func (d *protectionDomainResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.ProtectionDomainResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resp.Diagnostics.Append(d.ConfigurePdState(ctx, state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	tflog.Info(ctx, "Starting delete of protection domain "+state.Name.ValueString())
	err := d.pdClient.Delete()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error deleting protection domain.",
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "Finished delete of protection domain "+state.Name.ValueString())
	resp.State.RemoveResource(ctx)
}

func (d *protectionDomainResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// ReadByID is a helper function that refreshes the protection domain client and marshalls it into protectionDomainResourceModel
func (d *protectionDomainResource) ReadByID() (models.ProtectionDomainResourceModel, error) {
	// Fetch protection domain of given id
	err := d.pdClient.Refresh()
	if err != nil {
		return models.ProtectionDomainResourceModel{}, err
	}
	response := helper.GetPDResState(d.pdClient.ProtectionDomain)
	return response, nil
}

// UpdateResource is a common function called on create and update of pd to update all of its params
func (d *protectionDomainResource) UpdateResource(ctx context.Context, plan, state models.ProtectionDomainResourceModel) (dia diag.Diagnostics) {
	pd := d.pdClient
	if name := plan.Name.ValueString(); name != state.Name.ValueString() {
		err := pd.SetName(name)
		if err != nil {
			dia.AddError("Could not change name of protection domain.", err.Error())
		}
	}

	// Activate pd if required
	if !plan.Active.IsUnknown() && plan.Active.ValueBool() && !state.Active.ValueBool() {
		err := pd.Activate(true)
		if err != nil {
			dia.AddError("Could not activate protection domain.", err.Error())
		}
	}

	dia.Append(d.UpdateRfCache(ctx, plan, state)...)

	// set default FGL Metadata cache size - must be done before enabling caching for the first time
	if !plan.FglDefaultMetadataCacheSize.IsUnknown() && plan.FglDefaultMetadataCacheSize != state.FglDefaultMetadataCacheSize {
		err := pd.SetDefaultFGLMcacheSize(int(plan.FglDefaultMetadataCacheSize.ValueInt64()))
		if err != nil {
			dia.AddError("Could not set FGL Metadata cache size for protection domain.", err.Error())
		}
	}

	// enable/disable FGL Metadata cache
	if !plan.FglMetadataCacheEnabled.IsUnknown() && plan.FglMetadataCacheEnabled != state.FglMetadataCacheEnabled {
		var (
			err    error
			errMsg string
		)
		if plan.FglMetadataCacheEnabled.ValueBool() {
			err = pd.EnableFGLMcache()
			errMsg = "Could not enable FGL Metadata Cache"
		} else {
			err = pd.DisableFGLMcache()
			errMsg = "Could not disable FGL Metadata Cache"
		}
		if err != nil {
			dia.AddError(errMsg, err.Error())
		}
	}

	dia.Append(d.UpdateIopsLimits(ctx, plan, state)...)

	// InActivate pd if required
	if !plan.Active.IsUnknown() && !plan.Active.ValueBool() && state.Active.ValueBool() {
		err := pd.InActivate(true)
		if err != nil {
			dia.AddError("Could not inactivate protection domain.", err.Error())
		}
	}
	return dia
}

// UpdateRfCache is a common function called on create and update of pd to update its Rf cache params
func (d *protectionDomainResource) UpdateRfCache(ctx context.Context, plan, state models.ProtectionDomainResourceModel) (dia diag.Diagnostics) {
	pd := d.pdClient
	// Rfcache enable / disable
	if b := plan.RfCacheEnabled.ValueBool(); b != state.RfCacheEnabled.ValueBool() {
		var (
			err    error
			errMsg string
		)
		if b {
			err = pd.EnableRfcache()
			errMsg = "Could not enable RF Cache"
		} else {
			err = pd.DisableRfcache()
			errMsg = "Could not disable RF Cache"
		}
		if err != nil {
			dia.AddError(errMsg, err.Error())
		}
	}

	// Rfcache params
	rfcachePayload, ok := scaleiotypes.PDRfCacheParams{}, true
	if !plan.RfCacheOperationalMode.IsUnknown() && plan.RfCacheOperationalMode != state.RfCacheOperationalMode {
		ok = false
		rfcachePayload.RfCacheOperationalMode = scaleiotypes.PDRfCacheOpMode(plan.RfCacheOperationalMode.ValueString())
	}
	if !plan.RfCachePageSizeKb.IsUnknown() && plan.RfCachePageSizeKb != state.RfCachePageSizeKb {
		ok = false
		rfcachePayload.RfCachePageSizeKb = int(plan.RfCachePageSizeKb.ValueInt64())
	}
	if !plan.RfCacheMaxIoSizeKb.IsUnknown() && plan.RfCacheMaxIoSizeKb != state.RfCacheMaxIoSizeKb {
		ok = false
		rfcachePayload.RfCacheMaxIoSizeKb = int(plan.RfCacheMaxIoSizeKb.ValueInt64())
	}

	if !ok {
		err := pd.SetRfcacheParams(rfcachePayload)
		if err != nil {
			dia.AddError("Could not set RF cache params for protection domain.", err.Error())
		}
	}
	return dia
}

// UpdateIopsLimits is a common function called on create and update of pd to update its IOPS params
func (d *protectionDomainResource) UpdateIopsLimits(ctx context.Context, plan, state models.ProtectionDomainResourceModel) (dia diag.Diagnostics) {
	pd := d.pdClient
	// SDS IOPS params
	iopsPayload, ok := scaleiotypes.SdsNetworkLimitParams{}, true

	if !plan.RebuildNetworkThrottlingInKbps.IsUnknown() && plan.RebuildNetworkThrottlingInKbps != state.RebuildNetworkThrottlingInKbps {
		ok = false
		rebuildNw := int(plan.RebuildNetworkThrottlingInKbps.ValueInt64())
		iopsPayload.RebuildNetworkThrottlingInKbps = &rebuildNw
	}
	if !plan.RebalanceNetworkThrottlingInKbps.IsUnknown() && plan.RebalanceNetworkThrottlingInKbps != state.RebalanceNetworkThrottlingInKbps {
		ok = false
		rebalanceNw := int(plan.RebalanceNetworkThrottlingInKbps.ValueInt64())
		iopsPayload.RebalanceNetworkThrottlingInKbps = &rebalanceNw
	}
	if !plan.VTreeMigrationNetworkThrottlingInKbps.IsUnknown() && plan.VTreeMigrationNetworkThrottlingInKbps != state.VTreeMigrationNetworkThrottlingInKbps {
		ok = false
		vTreeNw := int(plan.VTreeMigrationNetworkThrottlingInKbps.ValueInt64())
		iopsPayload.VTreeMigrationNetworkThrottlingInKbps = &vTreeNw
	}
	if !plan.ProtectedMaintenanceModeNetworkThrottlingInKbps.IsUnknown() && plan.ProtectedMaintenanceModeNetworkThrottlingInKbps != state.ProtectedMaintenanceModeNetworkThrottlingInKbps {
		ok = false
		protectedMaintNw := int(plan.ProtectedMaintenanceModeNetworkThrottlingInKbps.ValueInt64())
		iopsPayload.ProtectedMaintenanceModeNetworkThrottlingInKbps = &protectedMaintNw
	}
	if !plan.OverallIoNetworkThrottlingInKbps.IsUnknown() && plan.OverallIoNetworkThrottlingInKbps != state.OverallIoNetworkThrottlingInKbps {
		ok = false
		overallNw := int(plan.OverallIoNetworkThrottlingInKbps.ValueInt64())
		iopsPayload.OverallIoNetworkThrottlingInKbps = &overallNw
	}

	if !ok {
		err := pd.SetSdsNetworkLimits(iopsPayload)
		if err != nil {
			dia.AddError("Could not set IOPS limits for protection domain.", err.Error())
		}
	}
	return dia
}
