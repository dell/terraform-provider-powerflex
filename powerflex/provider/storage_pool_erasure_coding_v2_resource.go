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
	_ resource.Resource              = &storagePoolErasureCodingV2Resource{}
	_ resource.ResourceWithConfigure = &storagePoolErasureCodingV2Resource{}
)

// NewStoragePoolErasureCodingV2Resource returns a new storage pool erasure coding v2 resource.
func NewStoragePoolErasureCodingV2Resource() resource.Resource {
	return &storagePoolErasureCodingV2Resource{}
}

// storagePoolErasureCodingV2Resource is the resource implementation.
type storagePoolErasureCodingV2Resource struct {
	client *clientgen.APIClient
}

// storagePoolErasureCodingModel maps the resource schema data.
type storagePoolErasureCodingModel struct {
	ID                  types.String `tfsdk:"id"`
	StoragePoolID       types.String `tfsdk:"storage_pool_id"`
	ErasureCodingPolicy types.String `tfsdk:"erasure_coding_policy"`
}

func (r *storagePoolErasureCodingV2Resource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storage_pool_erasure_coding"
}

func (r *storagePoolErasureCodingV2Resource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource is used to configure Erasure Coding policy on Storage Pools for PowerFlex Gen2 systems (5.0+).",
		MarkdownDescription: "This resource is used to configure Erasure Coding policy on Storage Pools for PowerFlex Gen2 systems (5.0+).",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed:            true,
				Description:         "The identifier for this resource.",
				MarkdownDescription: "The identifier for this resource.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"storage_pool_id": schema.StringAttribute{
				Required:            true,
				Description:         "The ID of the storage pool to configure erasure coding on.",
				MarkdownDescription: "The ID of the storage pool to configure erasure coding on.",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
			},
			"erasure_coding_policy": schema.StringAttribute{
				Required:            true,
				Description:         "The erasure coding policy to set. Valid values: none, rs_2_1, rs_4_1.",
				MarkdownDescription: "The erasure coding policy to set. Valid values: `none`, `rs_2_1`, `rs_4_1`.",
				Validators: []validator.String{
					stringvalidator.OneOf("none", "rs_2_1", "rs_4_1"),
				},
			},
		},
	}
}

func (r *storagePoolErasureCodingV2Resource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
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
			"The Gen2 API client is required for erasure coding management. Ensure the PowerFlex system supports Gen2 APIs (5.0+).",
		)
		return
	}

	r.client = p.genClient
}

// Create sets the erasure coding policy on a storage pool.
func (r *storagePoolErasureCodingV2Resource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	var plan storagePoolErasureCodingModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolID := plan.StoragePoolID.ValueString()
	policy := plan.ErasureCodingPolicy.ValueString()

	tflog.Debug(ctx, fmt.Sprintf("Setting Erasure Coding policy '%s' on storage pool '%s'", policy, poolID))

	payload := clientgen.SetErasureCodingPolicy{
		ErasureCodingPolicy: policy,
	}
	_, err := r.client.StoragePoolAPI.SetErasureCodingPolicy(ctx, poolID).SetErasureCodingPolicy(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Setting Erasure Coding Policy",
			fmt.Sprintf("Could not set erasure coding policy on pool %s: %s", poolID, err.Error()),
		)
		return
	}

	plan.ID = types.StringValue(poolID)

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	tflog.Debug(ctx, fmt.Sprintf("Set Erasure Coding policy '%s' on storage pool '%s'", policy, poolID))
}

// Read refreshes the Terraform state.
func (r *storagePoolErasureCodingV2Resource) Read(_ context.Context, _ resource.ReadRequest, _ *resource.ReadResponse) {
	// The erasure coding policy is a setting on the storage pool.
	// A full read would require fetching the pool and extracting the EC policy.
	// For now, this is a pass-through since the EC policy is set as an action.
}

// Update updates the erasure coding policy.
func (r *storagePoolErasureCodingV2Resource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan storagePoolErasureCodingModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolID := plan.StoragePoolID.ValueString()
	policy := plan.ErasureCodingPolicy.ValueString()

	payload := clientgen.SetErasureCodingPolicy{
		ErasureCodingPolicy: policy,
	}
	_, err := r.client.StoragePoolAPI.SetErasureCodingPolicy(ctx, poolID).SetErasureCodingPolicy(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Updating Erasure Coding Policy",
			fmt.Sprintf("Could not update erasure coding policy on pool %s: %s", poolID, err.Error()),
		)
		return
	}

	diags = resp.State.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Delete resets the erasure coding policy to "none".
func (r *storagePoolErasureCodingV2Resource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	var state storagePoolErasureCodingModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	poolID := state.StoragePoolID.ValueString()

	// Reset erasure coding policy to "none" on delete
	payload := clientgen.SetErasureCodingPolicy{
		ErasureCodingPolicy: "none",
	}
	_, err := r.client.StoragePoolAPI.SetErasureCodingPolicy(ctx, poolID).SetErasureCodingPolicy(payload).Execute()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Resetting Erasure Coding Policy",
			fmt.Sprintf("Could not reset erasure coding policy on pool %s: %s", poolID, err.Error()),
		)
		return
	}

	resp.State.RemoveResource(ctx)
	tflog.Debug(ctx, fmt.Sprintf("Reset Erasure Coding policy to 'none' on storage pool '%s'", poolID))
}
