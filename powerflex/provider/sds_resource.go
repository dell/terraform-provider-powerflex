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

	scaleiotypes "github.com/dell/goscaleio/types/v1"

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &sdsResource{}
	_ resource.ResourceWithConfigure   = &sdsResource{}
	_ resource.ResourceWithImportState = &sdsResource{}
)

// NewSDSResource is a helper function to simplify the provider implementation.
func NewSDSResource() resource.Resource {
	return &sdsResource{}
}

// sdsResource is the resource implementation.
type sdsResource struct {
	client *goscaleio.Client
}

func (r *sdsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sds"
}

func (r *sdsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SDSResourceSchema
}

func (r *sdsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*goscaleio.Client)
}

func (r *sdsResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var data models.SdsResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &data)...)

	if resp.Diagnostics.HasError() {
		return
	}

	// validate that same IP is not added in the set multiple times with different roles
	iplist := data.GetIPList(ctx)
	ipmap := make(map[string]int)
	// count how many times an IP is used in the set
	for _, ipObj := range iplist {
		ipmap[ipObj.IP]++
	}
	// raise errors for duplicate IP entries
	for ip, count := range ipmap {
		if count == 1 {
			continue
		}
		resp.Diagnostics.AddAttributeError(
			path.Root("ip_list"),
			"IP Duplication Error",
			fmt.Sprintf("The IP %s is configured with %d roles, but only 1 role expected.", ip, count),
		)
	}
}

// Create creates the resource and sets the initial Terraform state.
func (r *sdsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.SdsResourceModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if rmcache size is provided but rmcache is not enabled
	if !(plan.RmcacheSizeInMB.IsNull() || plan.RmcacheSizeInMB.IsUnknown()) && !plan.RmcacheEnabled.ValueBool() {
		resp.Diagnostics.AddAttributeError(
			path.Root("rmcache_size_in_mb"),
			"rmcache_size_in_mb cannot be specified while rmcache_enabled is not set to true",
			"Read Ram cache must be enabled in order to configure its size",
		)
	}
	if resp.Diagnostics.HasError() {
		return
	}

	pdm, err := helper.GetNewProtectionDomainEx(r.client, plan.ProtectionDomainID.ValueString(), plan.ProtectionDomainName.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			err.Error(),
		)
		return
	}

	// set the protection domain name in the plan so that it gets propagated to the state
	plan.ProtectionDomainName = types.StringValue(pdm.ProtectionDomain.Name)

	sdsName := plan.Name.ValueString()
	iplist := plan.GetIPList(ctx)

	params := scaleiotypes.Sds{
		Name:   sdsName,
		IPList: iplist,
	}
	if !plan.RmcacheEnabled.IsUnknown() {
		params.RmcacheEnabled = plan.RmcacheEnabled.ValueBool()
	}
	if !plan.RmcacheSizeInMB.IsUnknown() {
		params.RmcacheSizeInKb = int(plan.RmcacheSizeInMB.ValueInt64()) * 1024
	}
	if !plan.DrlMode.IsUnknown() {
		params.DrlMode = plan.DrlMode.ValueString()
	}
	if !plan.Port.IsUnknown() {
		params.Port = int(plan.Port.ValueInt64())
	}
	sdsID, err2 := pdm.CreateSdsWithParams(&params)
	if err2 != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create SDS with name %s and IP list %v", sdsName, iplist),
			err2.Error(),
		)
		return
	}

	// Get created SDS
	rsp, err3 := pdm.FindSds("ID", sdsID)
	if err3 != nil {
		resp.Diagnostics.AddError(
			"Error getting SDS after creation",
			err3.Error(),
		)
		return
	}
	// Set refreshed state
	state, dgs := helper.UpdateSdsState(rsp, plan)
	resp.Diagnostics.Append(dgs...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

	if !plan.RfcacheEnabled.IsUnknown() {
		rfCacheEnabled := plan.RfcacheEnabled.ValueBool()
		err := pdm.SetSdsRfCache(sdsID, rfCacheEnabled)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set SDS Rf Cache settings to %t", rfCacheEnabled),
				err.Error(),
			)
		}
	}

	if !plan.PerformanceProfile.IsUnknown() {
		perfprof := plan.PerformanceProfile.ValueString()
		err := pdm.SetSdsPerformanceProfile(sdsID, perfprof)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set SDS Performance Profile settings to %s", perfprof),
				err.Error(),
			)
		}
	}

	// Get updated SDS
	rsp, err4 := pdm.FindSds("ID", sdsID)
	if err4 != nil {
		resp.Diagnostics.AddError(
			"Error getting SDS after setting Rf cache and Performance Profile",
			err4.Error(),
		)
		return
	}
	// Set refreshed state
	state, dgs = helper.UpdateSdsState(rsp, plan)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *sdsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state models.SdsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get the system on the PowerFlex cluster
	system, err := helper.GetFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}

	// Get SDS
	var rsp scaleiotypes.Sds
	if rsp, err = system.GetSdsByID(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not get SDS by ID %s", state.ID.ValueString()),
			err.Error(),
		)
		return
	}

	// when SDS is imported, protection domain name is not known and this causes a non empty plan
	if state.ProtectionDomainName.IsNull() {
		protectionDomain, err := system.FindProtectionDomain(rsp.ProtectionDomainID, "", "")
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to read name of protection domain of ID %s for SDS %s", rsp.ProtectionDomainID, rsp.Name),
				err.Error(),
			)
		} else {
			state.ProtectionDomainName = types.StringValue(protectionDomain.Name)
		}
	}

	// Set refreshed state
	state, dgs := helper.UpdateSdsState(&rsp, state)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sdsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan models.SdsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state models.SdsResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if rm cache size is provided
	if !(plan.RmcacheSizeInMB.IsNull() || plan.RmcacheSizeInMB.IsUnknown()) {
		if plan.RmcacheEnabled.ValueBool() ||
			((plan.RmcacheEnabled.IsNull() || plan.RmcacheEnabled.IsUnknown()) && state.RmcacheEnabled.ValueBool()) {
			// if plan has explicitly rmcache enabled, no issues
			// if plan does not have explicitly rmcache enabled, but its enabled in state, again no problem
		} else {
			// else throw an error
			resp.Diagnostics.AddAttributeError(
				path.Root("rmcache_size_in_mb"),
				"rmcache_size_in_mb cannot be specified while rmcache_enabled is not set to true",
				"Read Ram cache must be enabled in order to configure its size",
			)
		}
	}

	if resp.Diagnostics.HasError() {
		return
	}

	pdm, err := helper.GetNewProtectionDomainEx(r.client, state.ProtectionDomainID.ValueString(), state.ProtectionDomainName.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			err.Error(),
		)
		return
	}

	// Check if there difference between plan and state
	if plan.Name.ValueString() != state.Name.ValueString() {
		err := pdm.SetSdsName(state.ID.ValueString(), plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not rename SDS",
				err.Error(),
			)
		}
	}

	// Check if there are updates in ip lists
	// Stop updating IPs if one IP updation fails
	func() {
		toAdd, toRmv, changed, _ := helper.SdsIPListDiff(ctx, &plan, &state)
		for _, ip := range toAdd {
			err := pdm.AddSdSIP(state.ID.ValueString(), ip.IP, ip.Role)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error adding IP %s to SDS with role %s", ip.IP, ip.Role),
					err.Error(),
				)
				return
			}
		}
		for _, ip := range changed {
			err := pdm.SetSDSIPRole(state.ID.ValueString(), ip.IP, ip.Role)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error updating IP %s role to %s in SDS", ip.IP, ip.Role),
					err.Error(),
				)
				return
			}
		}
		for _, ip := range toRmv {
			err := pdm.RemoveSDSIP(state.ID.ValueString(), ip.IP)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error removing IP %s with role %s from SDS", ip.IP, ip.Role),
					err.Error(),
				)
				return
			}
		}
	}()

	// check if there is change in sds port
	if !plan.Port.IsUnknown() && !state.Port.Equal(plan.Port) {
		port := plan.Port.ValueInt64()
		err := pdm.SetSdsPort(state.ID.ValueString(), int(port))
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not change SDS port to %d", port),
				err.Error(),
			)
		}
	}

	// check if there is change in sds drl mode
	if !plan.DrlMode.IsUnknown() && !state.DrlMode.Equal(plan.DrlMode) {
		drlMode := plan.DrlMode.ValueString()
		err := pdm.SetSdsDrlMode(state.ID.ValueString(), drlMode)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not change SDS DRL Mode to %s", drlMode),
				err.Error(),
			)
		}
	}

	// check if there is change in sds rmcache
	if !plan.RmcacheEnabled.IsUnknown() && !state.RmcacheEnabled.Equal(plan.RmcacheEnabled) {
		rmCacheEnabled := plan.RmcacheEnabled.ValueBool()
		err := pdm.SetSdsRmCache(state.ID.ValueString(), rmCacheEnabled)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not change SDS Read Ram Cache settings to %t", rmCacheEnabled),
				err.Error(),
			)
		}
	}
	if !plan.RmcacheSizeInMB.IsUnknown() && !state.RmcacheSizeInMB.Equal(plan.RmcacheSizeInMB) {
		rmCacheSize := plan.RmcacheSizeInMB.ValueInt64()
		err := pdm.SetSdsRmCacheSize(state.ID.ValueString(), int(rmCacheSize))
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not change SDS Read Ram Cache size to %d", rmCacheSize),
				err.Error(),
			)
		}
	}

	// check if there is change in sds rfcache
	if !plan.RfcacheEnabled.IsUnknown() && !state.RfcacheEnabled.Equal(plan.RfcacheEnabled) {
		rfCacheEnabled := plan.RfcacheEnabled.ValueBool()
		err := pdm.SetSdsRfCache(state.ID.ValueString(), rfCacheEnabled)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not change SDS Rf Cache settings to %t", rfCacheEnabled),
				err.Error(),
			)
		}
	}

	// Check if performance profile has been changed
	if !plan.PerformanceProfile.IsUnknown() && !state.PerformanceProfile.Equal(plan.PerformanceProfile) {
		perfprof := plan.PerformanceProfile.ValueString()
		err := pdm.SetSdsPerformanceProfile(state.ID.ValueString(), perfprof)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set SDS Performance Profile settings to %s", perfprof),
				err.Error(),
			)
		}
	}

	// Find updated SDS
	rsp, err := pdm.FindSds("ID", state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting SDS after updation",
			err.Error(),
		)
		return
	}

	// Set refreshed state
	state, dgs := helper.UpdateSdsState(rsp, state)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sdsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.SdsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pdm, err := helper.GetNewProtectionDomainEx(r.client, state.ProtectionDomainID.ValueString(), state.ProtectionDomainName.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			err.Error(),
		)
		return
	}

	// Delete SDS
	err = pdm.DeleteSds(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Powerflex SDS",
			err.Error(),
		)

		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)

}

func (r *sdsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
