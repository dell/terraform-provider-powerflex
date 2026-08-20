/*
Copyright (c) 2024-2026 Dell Inc., or its subsidiaries. All Rights Reserved.

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

	"terraform-provider-powerflex/clientgen"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &volumeActionV2Resource{}
	_ resource.ResourceWithConfigure = &volumeActionV2Resource{}
)

// NewVolumeActionV2Resource returns a new volume action v2 resource.
func NewVolumeActionV2Resource() resource.Resource {
	return &volumeActionV2Resource{}
}

type volumeActionV2Resource struct {
	client *clientgen.APIClient
}

type volumeActionModel struct {
	ID               types.String `tfsdk:"id"`
	VolumeID         types.String `tfsdk:"volume_id"`
	Action           types.String `tfsdk:"action"`
	SourceSnapshotID types.String `tfsdk:"source_snapshot_id"`
	HostID           types.String `tfsdk:"host_id"`
	AccessMode       types.String `tfsdk:"access_mode"`
}

func (r *volumeActionV2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume_action_v2"
}

func (r *volumeActionV2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource performs actions on volumes for PowerFlex Gen2 systems (5.0+). Supported actions: refresh, restore, map_to_host, unmap_from_host.",
		MarkdownDescription: "This resource performs actions on volumes for PowerFlex Gen2 systems (5.0+). Supported actions: `refresh`, `restore`, `map_to_host`, `unmap_from_host`.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The identifier for this resource.",
				MarkdownDescription: "The identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"volume_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the volume to perform the action on.",
				MarkdownDescription: "The ID of the volume to perform the action on.",
			},
			"action": schema.StringAttribute{
				Required:            true,
				Description:         "The action to perform. Valid values: refresh, restore, map_to_host, unmap_from_host.",
				MarkdownDescription: "The action to perform. Valid values: `refresh`, `restore`, `map_to_host`, `unmap_from_host`.",
				Validators: []validator.String{
					stringvalidator.OneOf("refresh", "restore", "map_to_host", "unmap_from_host"),
				},
			},
			"source_snapshot_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Source snapshot ID for restore action.",
				MarkdownDescription: "Source snapshot ID for `restore` action.",
			},
			"host_id": schema.StringAttribute{
				Optional:            true,
				Description:         "Host ID for map_to_host and unmap_from_host actions.",
				MarkdownDescription: "Host ID for `map_to_host` and `unmap_from_host` actions.",
			},
			"access_mode": schema.StringAttribute{
				Optional:            true,
				Description:         "Access mode for map_to_host action. Valid values: READ_ONLY, READ_WRITE.",
				MarkdownDescription: "Access mode for `map_to_host` action. Valid values: `READ_ONLY`, `READ_WRITE`.",
				Validators: []validator.String{
					stringvalidator.OneOf("READ_ONLY", "READ_WRITE"),
				},
			},
		},
	}
}

func (r *volumeActionV2Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	p, ok := req.ProviderData.(*powerflexProvider)
	if !ok {
		resp.Diagnostics.AddError(
			"Unexpected Resource Configure Type",
			fmt.Sprintf("Expected *powerflexProvider, got: %T.", req.ProviderData),
		)
		return
	}

	if p.genClient == nil {
		resp.Diagnostics.AddError(
			"Gen2 API Client Not Available",
			"The Gen2 API client is required for volume actions. Ensure the PowerFlex system supports Gen2 APIs (5.0+).",
		)
		return
	}

	r.client = p.genClient
}

func (r *volumeActionV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan volumeActionModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	volumeID := plan.VolumeID.ValueString()
	action := plan.Action.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Executing volume action '%s' on volume '%s'", action, volumeID))

	switch action {
	case "refresh":
		emptyBody := map[string]interface{}{}
		_, err := r.client.VolumeAPI.RefreshVolume(ctx, volumeID).Body(emptyBody).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Refreshing Volume",
				fmt.Sprintf("Could not refresh volume %s: %s", volumeID, err.Error()),
			)
			return
		}

	case "restore":
		sourceSnapshotID := plan.SourceSnapshotID.ValueString()
		if sourceSnapshotID == "" {
			resp.Diagnostics.AddError(
				"Missing Required Parameter",
				"source_snapshot_id is required when action is 'restore'.",
			)
			return
		}
		payload := clientgen.RestoreVolumePayload{
			SourceSnapshotId: sourceSnapshotID,
		}
		_, err := r.client.VolumeAPI.RestoreVolume(ctx, volumeID).RestoreVolumePayload(payload).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Restoring Volume",
				fmt.Sprintf("Could not restore volume %s: %s", volumeID, err.Error()),
			)
			return
		}

	case "map_to_host":
		hostID := plan.HostID.ValueString()
		if hostID == "" {
			resp.Diagnostics.AddError(
				"Missing Required Parameter",
				"host_id is required when action is 'map_to_host'.",
			)
			return
		}
		payload := clientgen.MapVolumeToHostPayload{
			HostId:     hostID,
			AccessMode: plan.AccessMode.ValueString(),
		}
		_, err := r.client.VolumeAPI.MapVolumeToHost(ctx, volumeID).MapVolumeToHostPayload(payload).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Mapping Volume to Host",
				fmt.Sprintf("Could not map volume %s to host %s: %s", volumeID, hostID, err.Error()),
			)
			return
		}

	case "unmap_from_host":
		hostID := plan.HostID.ValueString()
		if hostID == "" {
			resp.Diagnostics.AddError(
				"Missing Required Parameter",
				"host_id is required when action is 'unmap_from_host'.",
			)
			return
		}
		payload := clientgen.UnmapVolumeFromHostPayload{
			HostId: hostID,
		}
		_, err := r.client.VolumeAPI.UnmapVolumeFromHost(ctx, volumeID).UnmapVolumeFromHostPayload(payload).Execute()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Unmapping Volume from Host",
				fmt.Sprintf("Could not unmap volume %s from host %s: %s", volumeID, hostID, err.Error()),
			)
			return
		}
	}

	plan.ID = types.StringValue(fmt.Sprintf("%s-%s", volumeID, action))

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (r *volumeActionV2Resource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// Volume actions are fire-and-forget; state is maintained by the volume resource itself.
}

func (r *volumeActionV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Re-execute the action on update
	var plan volumeActionModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delegate to Create logic
	createReq := resource.CreateRequest{Plan: req.Plan}
	createResp := resource.CreateResponse{State: resp.State, Diagnostics: resp.Diagnostics}
	r.Create(ctx, createReq, &createResp)
	resp.State = createResp.State
	resp.Diagnostics = createResp.Diagnostics
}

func (r *volumeActionV2Resource) Delete(_ context.Context, _ resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Actions are fire-and-forget; nothing to clean up
	resp.State.RemoveResource(context.Background())
}
