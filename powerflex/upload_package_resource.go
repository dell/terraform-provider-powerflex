package powerflex

import (
	"context"
	"strconv"
	"strings"

	"github.com/dell/goscaleio"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// UploadPackageModel defines the struct for device resource
type UploadPackageModel struct {
	ID             types.String `tfsdk:"id"`
	FilePath       types.List   `tfsdk:"file_path"`
	PackageDetails types.Set    `tfsdk:"package_details"`
}

// NewUploadPackageResource is a helper function to simplify the provider implementation.
func NewUploadPackageResource() resource.Resource {
	return &uploadPackageResource{}
}

// uploadPackageResource is the resource implementation.
type uploadPackageResource struct {
	gatewayClient *goscaleio.GatewayClient
}

func (r *uploadPackageResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_uploadPackage"
}

func (r *uploadPackageResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource can be used to upload packages on a PowerFlex Gateway.",
		MarkdownDescription: "This resource can be used to upload packages on a PowerFlex Gateway.",
		Attributes: map[string]schema.Attribute{
			"file_path": schema.ListAttribute{
				Description:         "The Path of the directory of packages or package file ",
				Required:            true,
				ElementType:         types.StringType,
				MarkdownDescription: "The Path of the directory of packages or package file",
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
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"operating_system": schema.StringAttribute{
							Description:         "Supported OS.",
							Computed:            true,
							MarkdownDescription: "Supported OS.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"linux_flavour": schema.StringAttribute{
							Description:         "Type of Linux OS",
							Computed:            true,
							MarkdownDescription: "Type of Linux OS",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"version": schema.StringAttribute{
							Description:         "Uploaded Package Version.",
							Computed:            true,
							MarkdownDescription: "Uploaded Package Version.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"label": schema.StringAttribute{
							Description:         "Uploaded Package Minor Version with OS Combination.",
							Computed:            true,
							MarkdownDescription: "Uploaded Package Minor Version with OS Combination.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"type": schema.StringAttribute{
							Description:         "Type of Package.",
							Computed:            true,
							MarkdownDescription: "Type of Package. Like. MDM, LIA, SDS, SDC, etc.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
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

func (r *uploadPackageResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.gatewayClient = req.ProviderData.(*goscaleio.GatewayClient)
}

func (r *uploadPackageResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config UploadPackageModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}
}

// ModifyPlan modify resource plan attribute value
func (r *uploadPackageResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	var (
		plan UploadPackageModel
	)

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Create creates the resource and sets the initial Terraform state.
func (r *uploadPackageResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan UploadPackageModel
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
		packageDetailResponse, packageDetailError := r.gatewayClient.GetPackgeDetails()
		if packageDetailError != nil {
			resp.Diagnostics.AddError(
				"Error for getting package details.",
				"unexpected error: "+packageDetailError.Error(),
			)
			return
		}

		// Set refreshed state
		data, dgs := updateUploadPackageState(packageDetailResponse, plan)
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

func updateUploadPackageState(packageDetails []*goscaleio_types.PackageDetails, plan UploadPackageModel) (UploadPackageModel, diag.Diagnostics) {
	state := plan
	var diags diag.Diagnostics

	PackageAttrTypes := getPackageType()
	PackageElemType := types.ObjectType{
		AttrTypes: PackageAttrTypes,
	}

	packages := []attr.Value{}
	for _, vol := range packageDetails {
		objVal, dgs := getPackageValue(vol)
		diags = append(diags, dgs...)
		packages = append(packages, objVal)
	}
	setVal, dgs := types.SetValue(PackageElemType, packages)
	diags = append(diags, dgs...)
	state.PackageDetails = setVal
	state.ID = types.StringValue("placeholder")

	return state, diags
}

// getPackageType returns the Package type required for mapping
func getPackageType() map[string]attr.Type {
	return map[string]attr.Type{
		"file_name":        types.StringType,
		"operating_system": types.StringType,
		"linux_flavour":    types.StringType,
		"version":          types.StringType,
		"label":            types.StringType,
		"type":             types.StringType,
		"sio_patch_number": types.Int64Type,
		"size":             types.Int64Type,
		"latest":           types.BoolType,
	}
}

// getPackageValue returns the Package object required for mapping
func getPackageValue(packageDetails *goscaleio_types.PackageDetails) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(getPackageType(), map[string]attr.Value{
		"file_name":        types.StringValue(packageDetails.Filename),
		"operating_system": types.StringValue(packageDetails.OperatingSystem),
		"linux_flavour":    types.StringValue(packageDetails.LinuxFlavour),
		"version":          types.StringValue(packageDetails.Version),
		"label":            types.StringValue(packageDetails.Label),
		"type":             types.StringValue(packageDetails.Type),
		"sio_patch_number": types.Int64Value(int64(packageDetails.SioPatchNumber)),
		"size":             types.Int64Value(int64(packageDetails.Size)),
		"latest":           types.BoolValue(bool(packageDetails.Latest)),
	})
}

// Read refreshes the Terraform state with the latest data.
func (r *uploadPackageResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state UploadPackageModel

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	packageDetailResponse, packageDetailError := r.gatewayClient.GetPackgeDetails()
	if packageDetailError != nil {
		resp.Diagnostics.AddError(
			"Error for getting package details.",
			"unexpected error: "+packageDetailError.Error(),
		)
		return
	}

	// Set refreshed state
	data, dgs := updateUploadPackageState(packageDetailResponse, state)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *uploadPackageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan UploadPackageModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state UploadPackageModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	planFilePaths := []string{}
	plan.FilePath.ElementsAs(ctx, &planFilePaths, true)

	stateFilePaths := []string{}
	state.FilePath.ElementsAs(ctx, &stateFilePaths, true)

	removePackages := make([]string, 0)

	for _, path := range stateFilePaths {
		if !inslice(path, planFilePaths) {
			removePackages = append(removePackages, path)
		}
	}

	if len(removePackages) > 0 {
		for _, packageData := range removePackages {
			packageName := packageData[strings.LastIndex(packageData, "/")+1:]
			_, err := r.gatewayClient.DeletePackge(packageName)
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
		packgeDetailResponse, packgeDetailError := r.gatewayClient.GetPackgeDetails()
		if packgeDetailError != nil {
			resp.Diagnostics.AddError(
				"Error for getting package details.",
				"unexpected error: "+packgeDetailError.Error(),
			)
			return
		}

		// Set refreshed state
		data, dgs := updateUploadPackageState(packgeDetailResponse, plan)
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
func (r *uploadPackageResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state UploadPackageModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stateFilePaths := []string{}
	state.FilePath.ElementsAs(ctx, &stateFilePaths, true)

	for _, packageData := range stateFilePaths {
		packageName := packageData[strings.LastIndex(packageData, "/")+1:]

		_, err := r.gatewayClient.DeletePackge(packageName)
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
func (r *uploadPackageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}

func inslice(n string, h []string) bool {
	for _, v := range h {
		if v == n {
			return true
		}
	}
	return false
}