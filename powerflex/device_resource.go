package powerflex

import (
	"context"

	"github.com/dell/goscaleio"
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// DeviceModel defines the struct for device resource
type DeviceModel struct {
	ID                       types.String `tfsdk:"id"`
	Name                     types.String `tfsdk:"name"`
	DevicePath               types.String `tfsdk:"device_path"`
	ProtectionDomainName     types.String `tfsdk:"protection_domain_name"`
	ProtectionDomainID       types.String `tfsdk:"protection_domain_id"`
	StoragePoolName          types.String `tfsdk:"storage_pool_name"`
	StoragePoolID            types.String `tfsdk:"storage_pool_id"`
	SdsID                    types.String `tfsdk:"sds_id"`
	SdsName                  types.String `tfsdk:"sds_name"`
	MediaType                types.String `tfsdk:"media_type"`
	ExternalAccelerationType types.String `tfsdk:"external_acceleration_type"`
	DeviceCapacity           types.Int64  `tfsdk:"device_capacity"`
	DeviceCapacityInKB       types.Int64  `tfsdk:"device_capacity_in_kb"`
	DeviceState              types.String `tfsdk:"device_state"`
}

// NewDeviceResource is a helper function to simplify the provider implementation.
func NewDeviceResource() resource.Resource {
	return &deviceResource{}
}

// deviceResource is the resource implementation.
type deviceResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (r *deviceResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_device"
}

func (r *deviceResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource can be used to manage devices on a PowerFlex array.",
		MarkdownDescription: "This resource can be used to manage devices on a PowerFlex array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The ID of the device.",
				Computed:            true,
				MarkdownDescription: "The ID of the device.",
			},
			"name": schema.StringAttribute{
				Description:         "The name of the device.",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "The name of the device.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ConflictsWith(path.MatchRoot("id")),
				},
			},
			"device_path": schema.StringAttribute{
				Description:         "The path of the device.",
				Required:            true,
				MarkdownDescription: "The path of the device.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"storage_pool_id": schema.StringAttribute{
				Description:         "ID of the storage pool. Conflicts with 'storage_pool_name'. Cannot be updated.",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "ID of the storage pool. Conflicts with `storage_pool_name`. Cannot be updated.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(path.MatchRoot("storage_pool_name")),
				},
			},
			"storage_pool_name": schema.StringAttribute{
				Description:         "Name of the storage pool. Conflicts with 'storage_pool_id'. Cannot be updated.",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Name of the storage pool. Conflicts with `storage_pool_id`. Cannot be updated.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"protection_domain_id": schema.StringAttribute{
				Description:         "ID of the protection domain. Conflicts with 'protection_domain_name'. Cannot be updated.",
				MarkdownDescription: "ID of the protection domain. Conflicts with `protection_domain_name`. Cannot be updated.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ConflictsWith(path.MatchRoot("protection_domain_name")),
				},
			},
			"protection_domain_name": schema.StringAttribute{
				Description:         "Name of the protection domain. Conflicts with 'protection_domain_id'. Cannot be updated.",
				MarkdownDescription: "Name of the protection domain. Conflicts with `protection_domain_id`. Cannot be updated.",
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"sds_id": schema.StringAttribute{
				Description:         "ID of the SDS. Conflicts with 'sds_name'. Cannot be updated.",
				MarkdownDescription: "ID of the SDS. Conflicts with `sds_name`. Cannot be updated.",
				Computed:            true,
				Optional:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
					stringvalidator.ExactlyOneOf(path.MatchRoot("sds_name")),
				},
			},
			"sds_name": schema.StringAttribute{
				Description:         "Name of the SDS. Conflicts with 'sds_id'. Cannot be updated.",
				MarkdownDescription: "Name of the SDS. Conflicts with `sds_id`. Cannot be updated.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"media_type": schema.StringAttribute{
				Description:         "Media type of the device. Valid values are 'HDD', 'SSD'.",
				MarkdownDescription: "Media type of the device. Valid values are `HDD`, `SSD`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{stringvalidator.OneOf(
					"HDD",
					"SSD",
				)},
			},
			"external_acceleration_type": schema.StringAttribute{
				Description:         "External acceleration type of the device. Valid values are 'None', 'Read', 'Write', 'ReadAndWrite'.",
				MarkdownDescription: "External acceleration type of the device. Valid values are `None`, `Read`, `Write`, `ReadAndWrite`.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{stringvalidator.OneOf(
					"None",
					"Read",
					"Write",
					"ReadAndWrite",
				)},
			},
			"device_capacity": schema.Int64Attribute{
				Description:         "Capacity of the device in GB.",
				MarkdownDescription: "Capacity of the device in GB.",
				Optional:            true,
			},
			"device_capacity_in_kb": schema.Int64Attribute{
				Description:         "Capacity of the device in KB.",
				MarkdownDescription: "Capacity of the device in KB.",
				Computed:            true,
			},
			"device_state": schema.StringAttribute{
				Description:         "State of the device.",
				MarkdownDescription: "State of the device.",
				Computed:            true,
			},
		},
	}
}

func (r *deviceResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*goscaleio.Client)
	system, err := getFirstSystem(r.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}
	r.system = system
}

func getStoragePoolType(r *goscaleio.Client, storagePoolID string) (*goscaleio.StoragePool, error) {
	system, err := getFirstSystem(r)
	if err != nil {
		return nil, err
	}

	sp, err := system.GetStoragePoolByID(storagePoolID)
	if err != nil {
		return nil, err
	}

	sp1 := goscaleio.NewStoragePoolEx(r, sp)
	return sp1, nil
}

func (r *deviceResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config DeviceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	if !config.StoragePoolName.IsNull() {
		if config.ProtectionDomainID.IsNull() && config.ProtectionDomainName.IsNull() {
			resp.Diagnostics.AddAttributeError(
				path.Root("storage_pool_name"),
				"Please provide protection_domain_name or protection_domain_id with storage_pool_name.",
				"Please provide protection_domain_name or protection_domain_id with storage_pool_name.",
			)
		}
	}
}

// ModifyPlan modify resource plan attribute value
func (r *deviceResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	var (
		plan DeviceModel
		pd   *goscaleio.ProtectionDomain
		err  error
	)
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	if !plan.ProtectionDomainID.IsNull() || !plan.ProtectionDomainName.IsNull() {
		pd, err = getNewProtectionDomainEx(r.client, plan.ProtectionDomainID.ValueString(), plan.ProtectionDomainName.ValueString(), "")
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting protection domain details with ID: "+plan.ProtectionDomainID.ValueString()+" name: "+plan.ProtectionDomainName.ValueString(),
				err.Error(),
			)
			return
		}
	}

	if !plan.StoragePoolID.IsUnknown() {
		sp, err := r.system.GetStoragePoolByID(plan.StoragePoolID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting storage pool details with ID: "+plan.StoragePoolID.ValueString(),
				err.Error(),
			)
			return
		}
		plan.StoragePoolName = types.StringValue(sp.Name)
	} else if !plan.StoragePoolName.IsUnknown() {
		sp, err1 := pd.FindStoragePool("", plan.StoragePoolName.ValueString(), "")
		if err1 != nil {
			resp.Diagnostics.AddError(
				"Error in getting storage pool details with name: "+plan.ProtectionDomainName.ValueString(),
				err1.Error(),
			)
			return
		}
		plan.StoragePoolID = types.StringValue(sp.ID)
	}

	if !plan.SdsID.IsUnknown() {
		sds, err := r.system.GetSdsByID(plan.SdsID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting sds details with ID: "+plan.SdsID.ValueString(),
				err.Error(),
			)
			return
		}
		plan.SdsName = types.StringValue(sds.Name)
	} else if !plan.SdsName.IsUnknown() {
		sds, err := r.system.FindSds("Name", plan.SdsName.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting sds details with name: "+plan.SdsName.ValueString(),
				err.Error(),
			)
			return
		}
		plan.SdsID = types.StringValue(sds.ID)
	}

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Create creates the resource and sets the initial Terraform state.
func (r *deviceResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan DeviceModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceParam := &goscaleio_types.DeviceParam{
		Name:                     plan.Name.ValueString(),
		DeviceCurrentPathname:    plan.DevicePath.ValueString(),
		SdsID:                    plan.SdsID.ValueString(),
		StoragePoolID:            plan.StoragePoolID.ValueString(),
		MediaType:                plan.MediaType.ValueString(),
		ExternalAccelerationType: plan.ExternalAccelerationType.ValueString(),
	}

	sp, err := getStoragePoolType(r.client, plan.StoragePoolID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting storage pool instance with ID: "+plan.StoragePoolID.ValueString(),
			"unexpected error: "+err.Error(),
		)
		return
	}

	deviceID, err2 := sp.AttachDevice(deviceParam)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error adding device with path: "+plan.DevicePath.ValueString(),
			"unexpected error: "+err2.Error(),
		)
		return
	}

	deviceResponse, err3 := r.system.GetDevice(deviceID)
	if err3 != nil {
		resp.Diagnostics.AddError(
			"Error getting device with ID: "+deviceID,
			"unexpected error: "+err3.Error(),
		)
		return
	}

	// Set refreshed state
	state, dgs := updateDeviceState(deviceResponse, plan)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func updateDeviceState(deviceResponse *goscaleio_types.Device, plan DeviceModel) (DeviceModel, diag.Diagnostics) {
	state := plan
	var diags diag.Diagnostics

	state.ID = types.StringValue(deviceResponse.ID)
	if deviceResponse.Name == "" {
		state.Name = types.StringNull()
	} else {
		state.Name = types.StringValue(deviceResponse.Name)
	}

	state.MediaType = types.StringValue(deviceResponse.MediaType)
	state.ExternalAccelerationType = types.StringValue(deviceResponse.ExternalAccelerationType)
	state.DeviceCapacityInKB = types.Int64Value(int64(deviceResponse.CapacityLimitInKb))
	state.DeviceState = types.StringValue(deviceResponse.DeviceState)
	state.SdsID = types.StringValue(deviceResponse.SdsID)
	state.StoragePoolID = types.StringValue(deviceResponse.StoragePoolID)
	return state, diags
}

// Read refreshes the Terraform state with the latest data.
func (r *deviceResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state DeviceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	deviceResponse, err3 := r.system.GetDevice(state.ID.ValueString())
	if err3 != nil {
		resp.Diagnostics.AddError(
			"Error getting device with ID: "+state.ID.ValueString(),
			"unexpected error: "+err3.Error(),
		)
		return
	}

	// Set refreshed state
	state, dgs := updateDeviceState(deviceResponse, state)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *deviceResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *deviceResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state DeviceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sp, err := getStoragePoolType(r.client, state.StoragePoolID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting storage pool instance with ID: "+state.StoragePoolID.ValueString(),
			"unexpected error: "+err.Error(),
		)
		return
	}

	err = sp.RemoveDevice(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error removing device with ID: "+state.ID.ValueString(),
			"unexpected error: "+err.Error(),
		)
		return
	}

	resp.State.RemoveResource(ctx)
}

// ImportState imports the resource
func (r *deviceResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
}