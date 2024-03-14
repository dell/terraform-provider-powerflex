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
	"strconv"
	"strings"

	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	pftypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &snapshotResource{}
	_ resource.ResourceWithConfigure   = &snapshotResource{}
	_ resource.ResourceWithImportState = &snapshotResource{}
)

// NewSnapshotResource is a helper function to simplify the provider implementation.
func NewSnapshotResource() resource.Resource {
	return &snapshotResource{}
}

// snapshotResource is the resource implementation.
type snapshotResource struct {
	client *goscaleio.Client
}

// Metadata returns the resource type name.
func (r *snapshotResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_snapshot"
}

// Schema defines the schema for the resource.
func (r *snapshotResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SnapshotResourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *snapshotResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	r.client = req.ProviderData.(*powerflexProvider).client
}

// ModifyPlan modify resource plan attribute value
func (r *snapshotResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	var plan models.SnapshotResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if !plan.Size.IsNull() && !plan.Size.IsUnknown() {
		VSIKB := helper.ConverterKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
		plan.SizeInKb = types.Int64Value(int64(VSIKB))
	}
	if plan.SizeInKb.ValueInt64() == 0 {
		plan.Size = basetypes.NewInt64Unknown()
		plan.SizeInKb = basetypes.NewInt64Unknown()
	}
	if !plan.DesiredRetention.IsNull() && !plan.DesiredRetention.IsUnknown() {
		retentionInMin := helper.ConvertToMin(plan.DesiredRetention.ValueInt64(), plan.RetentionUnit.ValueString())
		plan.RetentionInMin = types.StringValue(retentionInMin)
	} else {
		if plan.DesiredRetention.IsNull() {
			plan.RetentionInMin = basetypes.NewStringNull()
		}
		if plan.DesiredRetention.IsUnknown() {
			plan.RetentionInMin = basetypes.NewStringUnknown()
		}
	}
	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Create creates the resource and sets the initial Terraform state.
func (r *snapshotResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan models.SnapshotResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	errMsg := make(map[string]string, 0)
	sr, err := helper.GetFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting first system",
			"unexpected error: "+err.Error(),
		)
		return
	}

	diags = r.getVolumeID(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	snapshotReqs := make([]*pftypes.SnapshotDef, 0)

	snapReq := &pftypes.SnapshotDef{
		VolumeID:     plan.VolumeID.ValueString(),
		SnapshotName: plan.Name.ValueString(),
	}
	snapshotReqs = append(snapshotReqs, snapReq)
	snapParam := &pftypes.SnapshotVolumesParam{
		SnapshotDefs: snapshotReqs,
		AccessMode:   plan.AccessMode.ValueString(),
	}
	// create a snapshot of one volume, this requires a volume id and snapshot as parameter
	snapResps, err := sr.CreateSnapshotConsistencyGroup(snapParam)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating snapshot",
			"unexpected err: "+err.Error(),
		)
		return
	}
	snapID := snapResps.VolumeIDList[0]
	snapResponse, err2 := r.client.GetVolume("", snapID, "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"unexpected error: "+err2.Error(),
		)
		// add state update here before so we don't lose the creation of snapshot
		return
	}
	snap := snapResponse[0]
	snapResource := goscaleio.NewVolume(r.client)
	snapResource.Volume = snap

	// setting snapshot size. default value will be equal to volume size
	if !plan.Size.IsUnknown() {
		vikb := helper.ConverterKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
		tflog.Info(ctx, "vikb"+strconv.FormatInt(vikb, 10))
		if int64(snapResource.Volume.SizeInKb) != vikb {
			switch plan.CapacityUnit.ValueString() {
			case "TB":
				err3 := snapResource.SetVolumeSize(strconv.FormatInt(plan.Size.ValueInt64()*1000, 10))
				if err3 != nil {
					errMsg["size/capacity_unit"] = err3.Error()
					// In case of failure, the size will be stored in GB capacity unit in state
					plan.SizeInKb = types.Int64Value(int64(snap.SizeInKb))
					plan.CapacityUnit = types.StringValue("GB")
				}
			case "GB":
				err3 := snapResource.SetVolumeSize(strconv.FormatInt(plan.Size.ValueInt64(), 10))
				if err3 != nil {
					errMsg["size/capacity_unit"] = err3.Error()
					// In case of failure, the size will be stored in GB capacity unit in state
					plan.SizeInKb = types.Int64Value(int64(snap.SizeInKb))
					plan.CapacityUnit = types.StringValue("GB")
				}
			}

		}
	}
	// locking the auto snapshot on finding LockedAutoSnapshot parameter as true
	if plan.LockAutoSnapshot.ValueBool() {
		err := snapResource.LockAutoSnapshot()
		if err != nil {
			errMsg["lock_auto_snapshot"] = err.Error()
		}
	}

	// disabling retention in case of error with update
	if (len(errMsg) > 0) && !plan.DesiredRetention.IsNull() {
		errMsg["desired_retention/retention_unit"] = "The specified snapshot can't be retained due to failure in creation."
	}
	if (len(errMsg) == 0) && !plan.DesiredRetention.IsNull() {
		errRetention := snapResource.SetSnapshotSecurity(plan.RetentionInMin.ValueString())
		if errRetention != nil {
			errMsg["desired_retention/retention_unit"] = errRetention.Error()
		}
	}
	snapResponse, err2 = r.client.GetVolume("", snapID, "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		return
	}
	snap = snapResponse[0]
	dgs := helper.RefreshState(snap, &plan)
	resp.Diagnostics.Append(dgs...)
	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if len(errMsg) > 0 {
		failureMessage := ""
		for key, value := range errMsg {
			failureMessage += key + " : " + value + ", "
		}
		failureMessage = strings.TrimSuffix(failureMessage, ", ")
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failure Message: [%v]", failureMessage),
			failureMessage)
	}
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *snapshotResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	var state models.SnapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	snapResponse, err2 := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		return
	}
	snap := snapResponse[0]
	dgs := helper.RefreshState(snap, &state)
	resp.Diagnostics.Append(dgs...)
	resp.Diagnostics.Append(diags...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *snapshotResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan models.SnapshotResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	var state models.SnapshotResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	errMsg := make(map[string]string, 0)

	if plan.VolumeName.ValueString() != state.VolumeName.ValueString() {
		if !plan.VolumeName.IsNull() {
			volResponse, err3 := r.client.GetVolume("", "", "", plan.VolumeName.ValueString(), false)
			if err3 != nil {
				resp.Diagnostics.AddError(
					"Error getting volume details",
					"Could not get volume, unexpected error: "+err3.Error(),
				)
				return
			}
			vol := volResponse[0]
			state.VolumeName = types.StringValue(vol.Name)
		} else if plan.VolumeName.IsNull() {
			state.VolumeName = types.StringNull()
		}
	}

	snapResponse, err2 := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		return
	}
	snap := snapResponse[0]
	snapResource := goscaleio.NewVolume(r.client)
	snapResource.Volume = snap

	// updating the name of volume if there is change in plan
	if plan.Name.ValueString() != state.Name.ValueString() {
		err := snapResource.SetVolumeName(plan.Name.ValueString())
		if err != nil {
			errMsg["name"] = err.Error()
		}
	}
	// updating the size of the volume if there is change in plan
	if !plan.SizeInKb.IsUnknown() && (plan.SizeInKb.ValueInt64() != state.SizeInKb.ValueInt64()) {
		sizeInGb, err1 := strconv.Atoi(strconv.FormatInt(plan.SizeInKb.ValueInt64(), 10))
		if err1 != nil {
			errMsg["size: int-to-string-conversion-error"] = err1.Error()
		}
		sizeInGb = sizeInGb / 1048576
		sizeInGB := strconv.FormatInt(int64(sizeInGb), 10)
		err2 := snapResource.SetVolumeSize(sizeInGB)
		if err2 != nil {
			errMsg["size/capacity_unit"] = err2.Error()
		}
	}
	// locking the snapshot in case of change in LockedAutoSnapshot state to true
	if plan.LockAutoSnapshot.ValueBool() && !state.LockAutoSnapshot.ValueBool() {
		err := snapResource.LockAutoSnapshot()
		if err != nil {
			errMsg["lock_auto_snapshot"] = err.Error()
		}
	}
	// unlocking the snapshot in case of change in LockedAutoSnapshot state to false
	if !plan.LockAutoSnapshot.ValueBool() && state.LockAutoSnapshot.ValueBool() {
		err := snapResource.UnlockAutoSnapshot()
		if err != nil {
			errMsg["lock_auto_snapshot"] = err.Error()
		}
	}

	// changing the access mode in case of change in access mode state
	if !plan.AccessMode.IsUnknown() && plan.AccessMode.ValueString() != state.AccessMode.ValueString() {
		err := snapResource.SetVolumeAccessModeLimit(plan.AccessMode.ValueString())
		if err != nil {
			errMsg["access"] = err.Error()
		}
	}

	// disabling retention in case of error with update
	if (len(errMsg) > 0) && !plan.DesiredRetention.IsNull() {
		errMsg["desired_retention/retention_unit"] = "The specified snapshot can't be retained due to failure in update."
	}

	// updating the retention in min if there is change in plan and state.
	if plan.RetentionInMin.ValueString() != state.RetentionInMin.ValueString() {
		err := snapResource.SetSnapshotSecurity(plan.RetentionInMin.ValueString())

		if err != nil {
			errMsg["desired_retention/retention_unit"] = err.Error()
		} else {
			state.DesiredRetention = plan.DesiredRetention
			state.RetentionUnit = plan.RetentionUnit
			state.RetentionInMin = plan.RetentionInMin
		}
	}

	// getting the updated snapshot instance
	snapResponse, err2 = r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		return
	}
	snap = snapResponse[0]

	// refreshing the state
	dgs := helper.RefreshState(snap, &state)
	resp.Diagnostics.Append(dgs...)
	// setting the state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

	// Adding error if the len of errMsg is greater than zero.
	if len(errMsg) > 0 {
		failureMessage := ""
		for key, value := range errMsg {
			failureMessage += key + " : " + value + ", "
		}
		failureMessage = strings.TrimSuffix(failureMessage, ", ")
		resp.Diagnostics.AddError(
			fmt.Sprintf("Failure Message: [%v]", failureMessage),
			failureMessage)
	}
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *snapshotResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state models.SnapshotResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	snapResponse, err2 := r.client.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting snapshot",
			"Could not get snapshot, unexpected error: "+err2.Error(),
		)
		return
	}
	snapshot := goscaleio.NewVolume(r.client)
	snapshot.Volume = snapResponse[0]

	err := snapshot.RemoveVolume(state.RemoveMode.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Removing Volume",
			"Couldn't remove volume "+err.Error(),
		)
		return
	}
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

func (r *snapshotResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getVolumeID updates the volume ID in the plan
func (r *snapshotResource) getVolumeID(ctx context.Context, plan *models.SnapshotResourceModel) (diags diag.Diagnostics) {
	if plan.VolumeName.ValueString() != "" {
		tflog.Info(ctx, fmt.Sprintf("Volume name is provided: %s", plan.VolumeName.ValueString()))
		snapResponse, err2 := r.client.GetVolume("", "", "", plan.VolumeName.ValueString(), false)
		if err2 != nil {
			diags.AddError(
				"Error getting volume by name",
				"unexpected error: "+err2.Error(),
			)
			return
		}
		plan.VolumeID = types.StringValue(snapResponse[0].ID)
	}
	return
}
