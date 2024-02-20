/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

	"reflect"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
)

var (
	_ resource.Resource                = &snapshotPolicyResource{}
	_ resource.ResourceWithConfigure   = &snapshotPolicyResource{}
	_ resource.ResourceWithImportState = &snapshotPolicyResource{}
)

// NewSnapshotPolicyResource - function to return resource interface
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
		Name:                             plan.Name.ValueString(),
		AutoSnapshotCreationCadenceInMin: plan.AutoSnapshotCreationCadenceInMin.String(),
		NumOfRetainedSnapshotsPerLevel:   stringList,
		Paused:                           plan.Paused.String(),
		SnapshotAccessMode:               plan.SnapshotAccessMode.ValueString(),
		SecureSnapshots:                  plan.SecureSnapshots.String(),
	}

	// create the snapshot policy
	snapID, err := r.system.CreateSnapshotPolicy(payload)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating snapshot policy",
			"Could not create snapshot policy, unexpected error: "+err.Error(),
		)
		return
	}

	// Assigning volumes to snapshot policy
	// If there is an issue while assigning the volumes then delete the entire resource
	stringListVol := helper.ListToSliceVol(plan)
	var mappedVols []string
	for _, v := range stringListVol {
		payload2 := &scaleiotypes.AssignVolumeToSnapshotPolicyParam{
			SourceVolumeId: v,
		}
		err2 := r.system.AssignVolumeToSnapshotPolicy(payload2, snapID)
		if err2 != nil {
			if len(mappedVols) == 0 {
				err := r.system.RemoveSnapshotPolicy(snapID)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error Deleting Snapshot Policy",
						"Couldn't Delete Snapshot Policy"+err.Error(),
					)
					return
				}
			} else if len(mappedVols) > 0 {
				for _, v := range mappedVols {
					payload2 := &scaleiotypes.AssignVolumeToSnapshotPolicyParam{
						SourceVolumeId:            v,
						AutoSnapshotRemovalAction: plan.RemoveMode.ValueString(),
					}
					err2 := r.system.UnassignVolumeFromSnapshotPolicy(payload2, snapID)
					if err2 != nil {
						resp.Diagnostics.AddError(
							"Error unassigning volume from snapshot policy",
							"Error unassigning volume from snapshot policy: "+err2.Error(),
						)
						return
					}
				}
				err := r.system.RemoveSnapshotPolicy(snapID)
				if err != nil {
					resp.Diagnostics.AddError(
						"Error Deleting Snapshot Policy",
						"Couldn't Delete Snapshot Policy"+err.Error(),
					)
					return
				}
			}
			resp.Diagnostics.AddError(
				"Error assigning volume to snapshot policy",
				"Error assigning volume to snapshot policy: "+err2.Error(),
			)
			return
		}
		mappedVols = append(mappedVols, v)
	}

	// Fetching the details of the snapshot policy
	response, err := r.client.GetSnapshotPolicy("", snapID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot policy after creation",
			"Could not get snapshot policy, unexpected error: "+err.Error(),
		)
		return
	}
	// fetching the list of volumes assigned to a snaphot policy
	volumes, err := r.system.GetSourceVolume(snapID)
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get volumes assigned to snapshot policy",
			err.Error(),
		)
		return
	}
	// updating the details to the state
	var state models.SnapshotPolicyResourceModel
	state = helper.UpdateSnapshotPolicyResourceState(response, volumes, &state)
	state.RemoveMode = plan.RemoveMode
	//setting the state
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
	// Get the details of the snapshot policy
	sp, err := r.client.GetSnapshotPolicy("", state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not get snapshot policy by ID %s", state.ID.ValueString()),
			err.Error(),
		)
		return
	}
	// get the volumes assigned to snapshot policy
	volumes, err := r.system.GetSourceVolume(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get volumes assigned to snapshot policy",
			err.Error(),
		)
		return
	}
	state = helper.UpdateSnapshotPolicyResourceState(sp, volumes, &state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Update snapshot policy resource
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

	//  Don't allow update operation for secure snapshots and snapshot access mode.
	if plan.SecureSnapshots.ValueBool() != state.SecureSnapshots.ValueBool() {
		resp.Diagnostics.AddError(
			"Cannot Update Secure Snapshots after creation",
			"Secure snapshot attribute cannot be updated once it is created.",
		)
		return
	}

	if plan.SnapshotAccessMode.ValueString() != state.SnapshotAccessMode.ValueString() {
		resp.Diagnostics.AddError(
			"Cannot Update snapshot access mode after creation",
			"Snapshot access mode attribute cannot be updated once it is created.",
		)
		return
	}

	// If there is a change in the name of the snapshot policy then update the name
	if plan.Name.ValueString() != state.Name.ValueString() {
		err := r.system.RenameSnapshotPolicy(state.ID.ValueString(), plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating name of snapshot policy", err.Error(),
			)
		}
	}
	// update auto snapshot creation cadence in min or number of retained snapshots per level or both
	var snapUpdate *scaleiotypes.SnapshotPolicyModifyParam
	if plan.AutoSnapshotCreationCadenceInMin.ValueInt64() != state.AutoSnapshotCreationCadenceInMin.ValueInt64() && reflect.DeepEqual(plan.NumOfRetainedSnapshotsPerLevel, state.NumOfRetainedSnapshotsPerLevel) {
		stringList := helper.ListToSlice(plan)
		snapUpdate = &scaleiotypes.SnapshotPolicyModifyParam{
			AutoSnapshotCreationCadenceInMin: plan.AutoSnapshotCreationCadenceInMin.String(),
			NumOfRetainedSnapshotsPerLevel:   stringList,
		}
		err := r.system.ModifySnapshotPolicy(snapUpdate, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating auto snapshot creation cadence ", err.Error(),
			)
		}
	} else if !reflect.DeepEqual(plan.NumOfRetainedSnapshotsPerLevel, state.NumOfRetainedSnapshotsPerLevel) && plan.AutoSnapshotCreationCadenceInMin.ValueInt64() == state.AutoSnapshotCreationCadenceInMin.ValueInt64() {
		stringList := helper.ListToSlice(plan)
		snapUpdate = &scaleiotypes.SnapshotPolicyModifyParam{
			AutoSnapshotCreationCadenceInMin: state.AutoSnapshotCreationCadenceInMin.String(),
			NumOfRetainedSnapshotsPerLevel:   stringList,
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
			NumOfRetainedSnapshotsPerLevel:   stringList,
		}
		err := r.system.ModifySnapshotPolicy(snapUpdate, state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error while updating auto snapshot creation cadence or num of retained snapshots", err.Error(),
			)
		}
	}

	// pause or resume snapshot policy
	if plan.Paused.ValueBool() != state.Paused.ValueBool() {
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

	// map or unmap volumes to snapshot policy
	var planVolList, stateVolList []string
	planVolList = helper.ListToSliceVol(plan)
	stateVolList = helper.ListToSliceVol(state)

	mapVolIds, unmapVolIds := helper.DifferenceArray(stateVolList, planVolList)

	if len(unmapVolIds) > 0 {
		if !plan.RemoveMode.IsNull() {
			for _, v := range unmapVolIds {
				payload2 := &scaleiotypes.AssignVolumeToSnapshotPolicyParam{
					SourceVolumeId:            v,
					AutoSnapshotRemovalAction: plan.RemoveMode.ValueString(),
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
			for _, v := range unmapVolIds {
				payload2 := &scaleiotypes.AssignVolumeToSnapshotPolicyParam{
					SourceVolumeId:            v,
					AutoSnapshotRemovalAction: state.RemoveMode.ValueString(),
				}
				err2 := r.system.UnassignVolumeFromSnapshotPolicy(payload2, state.ID.ValueString())
				if err2 != nil {
					resp.Diagnostics.AddError(
						"Error unassigning volume from snapshot policy",
						"Error unassigning volume from snapshot policy: "+err2.Error(),
					)
				}
			}
		}
	}
	if len(mapVolIds) > 0 {
		for _, v := range mapVolIds {
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

	response, err := r.client.GetSnapshotPolicy("", state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot policy after updation",
			"Could not get snapshot policy, unexpected error: "+err.Error(),
		)
		return
	}
	volumes, err := r.system.GetSourceVolume(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Could not get volumes assigned to snapshot policy",
			err.Error(),
		)
		return
	}
	state = helper.UpdateSnapshotPolicyResourceState(response, volumes, &state)

	if !plan.RemoveMode.IsNull() {
		state.RemoveMode = plan.RemoveMode
	}
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Delete Snapshot resource
func (r *snapshotPolicyResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.SnapshotPolicyResourceModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateVolList := helper.ListToSliceVol(state)

	if len(stateVolList) > 0 {
		for _, v := range stateVolList {
			payload2 := &scaleiotypes.AssignVolumeToSnapshotPolicyParam{
				SourceVolumeId:            v,
				AutoSnapshotRemovalAction: state.RemoveMode.ValueString(),
			}
			err2 := r.system.UnassignVolumeFromSnapshotPolicy(payload2, state.ID.ValueString())
			if err2 != nil {
				resp.Diagnostics.AddError(
					"Error unassigning volume from snapshot policy",
					"Error unassigning volume from snapshot policy: "+err2.Error(),
				)
			}
		}
	}

	err := r.system.RemoveSnapshotPolicy(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Snapshot Policy",
			"Couldn't Delete Snapshot Policy"+err.Error(),
		)
		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

// Function used to ImportState for snapshot policy Resource
func (r *snapshotPolicyResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
