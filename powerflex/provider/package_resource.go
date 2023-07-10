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
	"strconv"
	"strings"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// NewPackageResource is a helper function to simplify the provider implementation.
func NewPackageResource() resource.Resource {
	return &packageResource{}
}

// packageResource is the resource implementation.
type packageResource struct {
	gatewayClient *goscaleio.GatewayClient
}

func (r *packageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_package"
}

func (r *packageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource can be used to upload packages on a PowerFlex Gateway.",
		MarkdownDescription: "This resource can be used to upload packages on a PowerFlex Gateway.",
		Attributes: map[string]schema.Attribute{
			"file_path": schema.ListAttribute{
				Description:         "The list of path of packages",
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The list of path of packages",
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
				},
			},
			"id": schema.StringAttribute{
				Description:         "The ID of the package.",
				Computed:            true,
				MarkdownDescription: "The ID of the package.",
			},
			"package_details": schema.SetNestedAttribute{
				Description:         "Uploaded Packages details.",
				Computed:            true,
				MarkdownDescription: "Uploaded Packages details.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"file_name": schema.StringAttribute{
							Description:         "The Name of package.",
							Computed:            true,
							MarkdownDescription: "The Name of package.",
						},
						"operating_system": schema.StringAttribute{
							Description:         "Supported OS.",
							Computed:            true,
							MarkdownDescription: "Supported OS.",
						},
						"linux_flavour": schema.StringAttribute{
							Description:         "Type of Linux OS",
							Computed:            true,
							MarkdownDescription: "Type of Linux OS",
						},
						"version": schema.StringAttribute{
							Description:         "Uploaded Package Version.",
							Computed:            true,
							MarkdownDescription: "Uploaded Package Version.",
						},
						"label": schema.StringAttribute{
							Description:         "Uploaded Package Minor Version with OS Combination.",
							Computed:            true,
							MarkdownDescription: "Uploaded Package Minor Version with OS Combination.",
						},
						"type": schema.StringAttribute{
							Description:         "Type of Package.",
							Computed:            true,
							MarkdownDescription: "Type of Package. Like. MDM, LIA, SDS, SDC, etc.",
						},
						"sio_patch_number": schema.Int64Attribute{
							Description:         "Package Patch Number.",
							Computed:            true,
							MarkdownDescription: "Package Patch Number.",
						},
						"size": schema.Int64Attribute{
							Description:         "Size of Package.",
							Computed:            true,
							MarkdownDescription: "Size of Package.",
						},
						"latest": schema.BoolAttribute{
							Description:         "Package Version is latest or not.",
							Computed:            true,
							MarkdownDescription: "Package Version is latest or not",
						},
					},
				},
			},
		},
	}
}

func (r *packageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).gatewayClient != nil {
		r.gatewayClient = req.ProviderData.(*powerflexProvider).gatewayClient
	} else {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)

		return
	}

}

// Create creates the resource and sets the initial Terraform state.
func (r *packageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.PackageModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	filePaths := []string{}
	plan.FilePath.ElementsAs(ctx, &filePaths, true)

	uploadPackageResponse, uploadPackageError := r.gatewayClient.UploadPackages(filePaths)
	if uploadPackageError != nil {
		resp.Diagnostics.AddError(
			"Error getting with file path",
			"unexpected error: "+uploadPackageError.Error(),
		)
		return
	}

	if uploadPackageResponse.StatusCode == 200 {
		packageDetailResponse, packageDetailError := r.gatewayClient.GetPackageDetails()
		if packageDetailError != nil {
			resp.Diagnostics.AddError(
				"Error for getting package details.",
				"unexpected error: "+packageDetailError.Error(),
			)
			return
		}

		// Set refreshed state
		data, dgs := helper.UpdateUploadPackageState(packageDetailResponse, plan)
		resp.Diagnostics.Append(dgs...)

		diags = resp.State.Set(ctx, data)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.Diagnostics.AddError(
			"Error while uploading package :"+uploadPackageResponse.Message+" & Error Code :"+strconv.Itoa(uploadPackageResponse.ErrorCode),
			"Status Code:"+strconv.Itoa(uploadPackageResponse.StatusCode),
		)
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *packageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state models.PackageModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	packageDetailResponse, packageDetailError := r.gatewayClient.GetPackageDetails()
	if packageDetailError != nil {
		resp.Diagnostics.AddError(
			"Error for getting package details.",
			"unexpected error: "+packageDetailError.Error(),
		)
		return
	}

	// Set refreshed state
	data, dgs := helper.UpdateUploadPackageState(packageDetailResponse, state)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *packageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan models.PackageModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.PackageModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	planFilePaths := []string{}
	plan.FilePath.ElementsAs(ctx, &planFilePaths, true)

	stateFilePaths := []string{}
	state.FilePath.ElementsAs(ctx, &stateFilePaths, true)

	planFileMap := make(map[string]string)
	stateFileMap := make(map[string]string)

	// Populate planFileMap with the file paths defined in plan
	for _, filePath := range planFilePaths {
		planFileMap[filePath] = filePath
	}

	// Populate stateFileMap with the file paths stored in state
	for _, filePath := range stateFilePaths {
		stateFileMap[filePath] = filePath
	}

	removePackages := helper.DifferenceMap(stateFileMap, planFileMap)

	if len(removePackages) > 0 {
		for _, packageData := range removePackages {
			packageName := packageData[strings.LastIndex(packageData, "/")+1:]
			_, err := r.gatewayClient.DeletePackage(packageName)
			if err != nil {
				resp.Diagnostics.AddError(
					"Error removing package with Name: "+packageName,
					"unexpected error: "+err.Error(),
				)
				return
			}

		}
	}

	uploadPackageResponse, uploadPackageError := r.gatewayClient.UploadPackages(planFilePaths)
	if uploadPackageError != nil {
		resp.Diagnostics.AddError(
			"Error getting with upload package.",
			"unexpected error: "+uploadPackageError.Error(),
		)
		return
	}

	if uploadPackageResponse.StatusCode == 200 {
		packgeDetailResponse, packgeDetailError := r.gatewayClient.GetPackageDetails()
		if packgeDetailError != nil {
			resp.Diagnostics.AddError(
				"Error for getting package details.",
				"unexpected error: "+packgeDetailError.Error(),
			)
			return
		}

		// Set refreshed state
		data, dgs := helper.UpdateUploadPackageState(packgeDetailResponse, plan)
		resp.Diagnostics.Append(dgs...)

		diags = resp.State.Set(ctx, data)
		resp.Diagnostics.Append(diags...)
	} else {
		resp.Diagnostics.AddError(
			"Error while uploading package :"+uploadPackageResponse.Message+" & Error Code :"+strconv.Itoa(uploadPackageResponse.ErrorCode),
			"Status Code:"+strconv.Itoa(uploadPackageResponse.StatusCode),
		)
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *packageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.PackageModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateFilePaths := []string{}
	state.FilePath.ElementsAs(ctx, &stateFilePaths, true)

	for _, packageData := range stateFilePaths {
		packageName := packageData[strings.LastIndex(packageData, "/")+1:]

		_, err := r.gatewayClient.DeletePackage(packageName)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error removing package with Name: "+packageName,
				"unexpected error: "+err.Error(),
			)
			return
		}
	}
}

// ImportState imports the resource
func (r *packageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}
