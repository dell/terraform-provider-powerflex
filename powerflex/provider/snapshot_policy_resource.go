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
	"terraform-provider-powerflex/powerflex/helper"

	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"reflect"
)

var (
	_ resource.Resource                = &snapshotPolicyResource{}
	_ resource.ResourceWithConfigure   = &snapshotPolicyResource{}
	_ resource.ResourceWithImportState = &snapshotPolicyResource{}
)

// NewFaultSetResource - function to return resource interface
func NewSnapshotPolicyResource() resource.Resource {
	return &snapshotPolicyResource{}
}

type snapshotPolicyResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (r *snapshotPolicyResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot_policy"
}

func (r *snapshotPolicyResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SnapshotPolicyResourceSchema
}

func (r *snapshotPolicyResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	r.client = req.ProviderData.(*powerflexProvider).client
	system, err := helper.GetFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex System",
			err.Error(),
		)
		return
	}
	r.system = system
}

// Function used to create snapshot policy resource
func (r *snapshotPolicyResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.SnapshotPolicyResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// converting a list to slice 
	stringList := helper.ListToSlice(plan)
	payload := &scaleiotypes.SnapshotPolicyCreateParam{
		Name:               plan.Name.ValueString(),
		AutoSnapshotCreationCadenceInMin: plan.AutoSnapshotCreationCadenceInMin.String(),
		NumOfRetainedSnapshotsPerLevel: stringList,
		Paused: plan.Paused.String(),
		SnapshotAccessMode: plan.SnapshotAccessMode.ValueString(),
		SecureSnapshots: plan.SecureSnapshots.String(),
	}

	// create the fault set
	snapID, err := r.system.CreateSnapshotPolicy(payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating snapshot policy",
			"Could not create snapshot policy, unexpected error: "+err.Error(),
		)
		return
	}

	if !plan.VolumeId.IsNull() {
		for _,v := range plan.VolumeId.Elements() {
			payload2 := &scaleiotypes.AssignVolumeToSnapshotPolicyParam{
				SourceVolumeId: v.String(),
			}
			err2 := r.system.AssignVolumeToSnapshotPolicy(payload2, snapID)
			if err2 != nil {
				resp.Diagnostics.AddError(
					"Error assigning volume to snapshot policy",
					"Error assigning volume to snapshot policy: "+err2.Error(),
				)
			}
		}
	}

	response, err := r.client.GetSnapshotPolicy("",snapID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot policy after creation",
			"Could not get snapshot policy, unexpected error: "+err.Error(),
		)
		return
	}

	state := helper.UpdateSnapshotPolicyResourceState(response)
	if !plan.Paused.IsNull(){
		state.Paused = plan.Paused
	}

	if !plan.VolumeId.IsNull(){
		state.VolumeId = plan.VolumeId
	}
	
	
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Read Snapshot Policy Resource
func (r *snapshotPolicyResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state models.SnapshotPolicyResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sp, err := r.client.GetSnapshotPolicy("", state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not get snapshot policy by ID %s", state.ID.ValueString()),
			err.Error(),
		)
		return
	}
	faultsetState := helper.UpdateSnapshotPolicyResourceState(sp)
	diags = resp.State.Set(ctx, faultsetState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Update fault set Resource
func (r *snapshotPolicyResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan models.SnapshotPolicyResourceModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//Get Current State
	var state models.SnapshotPolicyResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	if plan.Name.ValueString() != state.Name.ValueString() {
		err := r.system.RenameSnapshotPolicy(state.ID.ValueString(), plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating name of snapshot policy", err.Error(),
			)
		}
	}

	var snapUpdate *scaleiotypes.SnapshotPolicyModifyParam
	if plan.AutoSnapshotCreationCadenceInMin.ValueInt64() != state.AutoSnapshotCreationCadenceInMin.ValueInt64() && reflect.DeepEqual(plan.NumOfRetainedSnapshotsPerLevel, state.NumOfRetainedSnapshotsPerLevel) {
		stringList := helper.ListToSlice(plan)
		snapUpdate = &scaleiotypes.SnapshotPolicyModifyParam{
			AutoSnapshotCreationCadenceInMin: plan.AutoSnapshotCreationCadenceInMin.String(),
			NumOfRetainedSnapshotsPerLevel: stringList,
		}
		err := r.system.ModifySnapshotPolicy(snapUpdate, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating auto snapshot creation cadence ", err.Error(),
			)
		}
	} else if !reflect.DeepEqual(plan.NumOfRetainedSnapshotsPerLevel, state.NumOfRetainedSnapshotsPerLevel) && plan.AutoSnapshotCreationCadenceInMin.ValueInt64() == state.AutoSnapshotCreationCadenceInMin.ValueInt64()  {
		stringList := helper.ListToSlice(plan)
		snapUpdate = &scaleiotypes.SnapshotPolicyModifyParam{
			AutoSnapshotCreationCadenceInMin: state.AutoSnapshotCreationCadenceInMin.String(),
			NumOfRetainedSnapshotsPerLevel: stringList,
		}
		err := r.system.ModifySnapshotPolicy(snapUpdate, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating num of retained snapshots per level ", err.Error(),
			)
		}
	} else if plan.AutoSnapshotCreationCadenceInMin.ValueInt64() != state.AutoSnapshotCreationCadenceInMin.ValueInt64() && !reflect.DeepEqual(plan.NumOfRetainedSnapshotsPerLevel, state.NumOfRetainedSnapshotsPerLevel) {
		stringList := helper.ListToSlice(plan)
		snapUpdate = &scaleiotypes.SnapshotPolicyModifyParam{
			AutoSnapshotCreationCadenceInMin: plan.AutoSnapshotCreationCadenceInMin.String(),
			NumOfRetainedSnapshotsPerLevel: stringList,
		}
		err := r.system.ModifySnapshotPolicy(snapUpdate, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating auto snapshot creation cadence and num of retained snapshots", err.Error(),
			)
	}
}	

	if plan.Paused.ValueBool() != state.Paused.ValueBool(){
		if plan.Paused.ValueBool() {
			err := r.system.PauseSnapshotPolicy(state.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error while pausing the snapshot policy", err.Error(),
				)
		}
		} else {
			err := r.system.ResumeSnapshotPolicy(state.ID.ValueString())
			if err != nil {
				resp.Diagnostics.AddError(
					"Error while resuming the snapshot policy", err.Error(),
				)
		}
		}
	}

	var planVolList,stateVolList []string

	diags = plan.VolumeId.ElementsAs(ctx, &planVolList, true)
	resp.Diagnostics.Append(diags...)

	diags = state.VolumeId.ElementsAs(ctx, &stateVolList, true)
	resp.Diagnostics.Append(diags...)

	mapVolIds,unmapVolIds := helper.DifferenceArray(planVolList, stateVolList)

	if len(mapVolIds) > 0 {
		for _,v := range mapVolIds {
			payload2 := &scaleiotypes.AssignVolumeToSnapshotPolicyParam{
				SourceVolumeId: v,
			}
			err2 := r.system.AssignVolumeToSnapshotPolicy(payload2, state.ID.ValueString())
			if err2 != nil {
				resp.Diagnostics.AddError(
					"Error assigning volume to snapshot policy",
					"Error assigning volume to snapshot policy: "+err2.Error(),
				)
			}
		}
	}

	if len(unmapVolIds) > 0 {
		if !plan.RemoveMode.IsNull(){
			for _,v := range unmapVolIds {
				payload2 := &scaleiotypes.AssignVolumeToSnapshotPolicyParam{
					SourceVolumeId: v,
					AutoSnapshotRemovalAction: plan.RemoveMode.ValueString() ,
				}
				err2 := r.system.UnassignVolumeFromSnapshotPolicy(payload2, state.ID.ValueString())
				if err2 != nil {
					resp.Diagnostics.AddError(
						"Error unassigning volume from snapshot policy",
						"Error unassigning volume from snapshot policy: "+err2.Error(),
					)
				}
			}
		} else {
			resp.Diagnostics.AddError(
				"Missing Remove mode",
				"Missing Remove mode",
			)
		}
		
	}

	response, err := r.client.GetSnapshotPolicy("",state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot policy after updation",
			"Could not get snapshot policy, unexpected error: "+err.Error(),
		)
		return
	}

	state2 := helper.UpdateSnapshotPolicyResourceState(response)
	if plan.Paused.ValueBool() != state.Paused.ValueBool(){
		state.Paused = plan.Paused
	}

	if len(mapVolIds) > 0 || len(unmapVolIds) > 0 {
		state.VolumeId = plan.VolumeId
	}
	
	
	diags = resp.State.Set(ctx, state2)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}


}

// Function used to Delete Fault Set Resource
func (r *snapshotPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.FaultSetResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pd, err := helper.GetNewProtectionDomainEx(r.client, state.ProtectionDomainID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			"Could not get Protection Domain, unexpected err: "+err.Error(),
		)
		return
	}

	err = pd.DeleteFaultSet(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting fault set",
			"Couldn't Delete fault set "+err.Error(),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

// Function used to ImportState for fault set Resource
func (r *snapshotPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
