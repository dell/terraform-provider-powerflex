package powerflex

import (
	"context"
	"strings"

	"github.com/dell/goscaleio"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
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
	FilePath      types.String `tfsdk:"file_path"`
	PackageName   types.String `tfsdk:"package_name"`
	PackageParams types.Set    `tfsdk:"package_params"`
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
			"file_path": schema.StringAttribute{
				Description:         "The File Path of the package.",
				Required:            true,
				MarkdownDescription: "The File Path of the package.",
			},
			"package_name": schema.StringAttribute{
				Description:         "The name of package.",
				Optional:            true,
				MarkdownDescription: "The name of package.",
			},
			"package_params": schema.SetNestedAttribute{
				Description:         "The name of package.",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of package.",
				NestedObject: schema.NestedAttributeObject{
					Attributes: map[string]schema.Attribute{
						"file_name": schema.StringAttribute{
							Description:         "The ID of the volume.",
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "The ID of the volume.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"operating_system": schema.StringAttribute{
							Description:         "The name of the volume.",
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The name of the volume.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"linux_flavour": schema.StringAttribute{
							Description:         "The name of the volume.",
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The name of the volume.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"version": schema.StringAttribute{
							Description:         "The name of the volume.",
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The name of the volume.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"label": schema.StringAttribute{
							Description:         "The name of the volume.",
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The name of the volume.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"type": schema.StringAttribute{
							Description:         "The name of the volume.",
							Computed:            true,
							Optional:            true,
							MarkdownDescription: "The name of the volume.",
							Validators: []validator.String{
								stringvalidator.LengthAtLeast(1),
							},
						},
						"sio_patch_number": schema.Int64Attribute{
							Description:         "IOPS limit. Valid values are 0 or integers greater than 10. 0 represents unlimited IOPS. Default value is 0.",
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "IOPS limit. Valid values are 0 or integers greater than 10. 0 represents unlimited IOPS. Default value is 0.",
						},
						"size": schema.Int64Attribute{
							Description:         "IOPS limit. Valid values are 0 or integers greater than 10. 0 represents unlimited IOPS. Default value is 0.",
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "IOPS limit. Valid values are 0 or integers greater than 10. 0 represents unlimited IOPS. Default value is 0.",
						},
						"latest": schema.BoolAttribute{
							Description:         "IOPS limit. Valid values are 0 or integers greater than 10. 0 represents unlimited IOPS. Default value is 0.",
							Optional:            true,
							Computed:            true,
							MarkdownDescription: "IOPS limit. Valid values are 0 or integers greater than 10. 0 represents unlimited IOPS. Default value is 0.",
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

	_, err3 := r.gatewayClient.UploadPackages(plan.FilePath.ValueString())
	if err3 != nil {
		resp.Diagnostics.AddError(
			"Error getting with file path: "+plan.FilePath.ValueString(),
			"unexpected error: "+err3.Error(),
		)
		return
	}

	res, err3 := r.gatewayClient.GetPackgeDetails()
	if err3 != nil {
		resp.Diagnostics.AddError(
			"Error getting pacckage details:",
			"unexpected error: "+err3.Error(),
		)
		return
	}

	// Set refreshed state
	data, dgs := updateUploadPackageState(res, plan)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

func updateUploadPackageState(packageDetails []*goscaleio_types.PackageParam, plan UploadPackageModel) (UploadPackageModel, diag.Diagnostics) {
	var state UploadPackageModel
	state = plan
	var diags diag.Diagnostics

	PackageAttrTypes := getPackageType()

	PackageElemType := types.ObjectType{
		AttrTypes: PackageAttrTypes,
	}

	objectPackages := []attr.Value{}
	for _, vol := range packageDetails {
		objVal, dgs := getPackageValue(vol)
		diags = append(diags, dgs...)
		objectPackages = append(objectPackages, objVal)
		// state.Name = types.StringValue(vol.MappedSdcInfo[0].SdcName)
		// state.ID = types.StringValue(vol.MappedSdcInfo[0].SdcID)
	}
	setVal, dgs := types.SetValue(PackageElemType, objectPackages)
	diags = append(diags, dgs...)
	state.PackageParams = setVal

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
func getPackageValue(packageParam *goscaleio_types.PackageParam) (basetypes.ObjectValue, diag.Diagnostics) {
	return types.ObjectValue(getPackageType(), map[string]attr.Value{
		"file_name":        types.StringValue(packageParam.Filename),
		"operating_system": types.StringValue(packageParam.OperatingSystem),
		"linux_flavour":    types.StringValue(packageParam.LinuxFlavour),
		"version":          types.StringValue(packageParam.Version),
		"label":            types.StringValue(packageParam.Label),
		"type":             types.StringValue(packageParam.Type),
		"sio_patch_number": types.Int64Value(int64(packageParam.SioPatchNumber)),
		"size":             types.Int64Value(int64(packageParam.Size)),
		"latest":           types.BoolValue(bool(packageParam.Latest)),
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

	res, err3 := r.gatewayClient.GetPackgeDetails()
	if err3 != nil {
		resp.Diagnostics.AddError(
			"Error getting with file path: "+state.FilePath.ValueString(),
			"unexpected error: "+err3.Error(),
		)
		return
	}

	// Set refreshed state
	data, dgs := updateUploadPackageState(res, state)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *uploadPackageResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	var plan UploadPackageModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// uploadPackageParam := &goscaleio_types.UploadPackageParam{
	// 	FilePath: plan.FilePath.ValueString(),
	// }

	_, err3 := r.gatewayClient.UploadPackages(plan.FilePath.ValueString())
	if err3 != nil {
		resp.Diagnostics.AddError(
			"Error getting with file path: "+plan.FilePath.ValueString(),
			"unexpected error: "+err3.Error(),
		)
		return
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

	packageName := state.FilePath.ValueString()[strings.LastIndex(state.FilePath.ValueString(), "/")+1:]

	_, err := r.gatewayClient.DeletePackge(packageName)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error removing device with ID: "+packageName,
			"unexpected error: "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource
func (r *uploadPackageResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}
