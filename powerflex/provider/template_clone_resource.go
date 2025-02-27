/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"terraform-provider-powerflex/powerflex/constants"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &templateCloneResource{}
	_ resource.ResourceWithConfigure = &templateCloneResource{}
)

// SystemResource - function to return resource interface
func TemplateCloneResource() resource.Resource {
	return &templateCloneResource{}
}

// templateCloneResource - struct to define template clone resource
type templateCloneResource struct {
	gc     *goscaleio.GatewayClient
	sys    *goscaleio.System
	client *goscaleio.Client
}

// Metadata - function to return metadata for system resource.
func (r *templateCloneResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_template_clone"
}

// Schema - function to return Schema for system resource.
func (r *templateCloneResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = TemplateCloneResourceSchema
}

// Configure - function to return Configuration for template Clone resource.
func (r *templateCloneResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}
	if req.ProviderData.(*powerflexProvider).client != nil {
		r.client = req.ProviderData.(*powerflexProvider).client
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
	r.sys = system

	if req.ProviderData.(*powerflexProvider).gatewayClient != nil {
		r.gc = req.ProviderData.(*powerflexProvider).gatewayClient
	} else {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}
}

func (r *templateCloneResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Create")
	var plan models.TemplateCloneModel

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	err := r.gc.CloneTemplate(r.sys, plan.OriginalTemplateID.ValueString(), plan.TemplateName.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Error doing clone on template: %s", plan.TemplateName.ValueString()),
			err.Error(),
		)
		return
	}

	diagsState := resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diagsState...)
}

func (r *templateCloneResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Read")
	var state models.TemplateCloneModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diagsState := resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diagsState...)
}

func (r *templateCloneResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	resp.Diagnostics.AddWarning(
		constants.UpdatesNotSupportedErrorMsg+", Template clone resource does not support update operation",
		constants.UpdatesNotSupportedErrorMsg+", Template clone resource does not support update operation",
	)
}

func (r *templateCloneResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	resp.Diagnostics.AddWarning(
		constants.DeleteIsNotSupportedErrorMsg+", Template clone resource does not support delete operation",
		constants.DeleteIsNotSupportedErrorMsg+", Template clone resource does not support delete operation",
	)
}
