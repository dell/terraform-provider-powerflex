package volume

import (
	"context"
	"strconv"

	// "time"

	"terraform-provider-powerflex/helper"

	"github.com/dell/goscaleio"
	pftypes "github.com/dell/goscaleio/types/v1"

	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"

	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	// "github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &volumeResource{}
	_ resource.ResourceWithConfigure   = &volumeResource{}
	_ resource.ResourceWithImportState = &volumeResource{}
)

// NewvolumeResource is a helper function to simplify the provider implementation.
func NewVolumeResource() resource.Resource {
	return &volumeResource{}
}

// volumeResource is the resource implementation.
type volumeResource struct {
	client *goscaleio.Client
}

// Metadata returns the resource type name.
func (r *volumeResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_volume"
}

// Schema defines the schema for the resource.
func (r *volumeResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = VolumeResourceSchema
}

// Configure adds the provider configured client to the data source.
func (r *volumeResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*goscaleio.Client)
}

// Create creates the resource and sets the initial Terraform state.
func (r *volumeResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan volumeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	volumeCreate := &pftypes.VolumeParam{
		ProtectionDomainID: plan.ProtectionDomainID.ValueString(),
		StoragePoolID:      plan.StoragePoolID.ValueString(),
		UseRmCache:         strconv.FormatBool(plan.UseRmCache.ValueBool()),
		VolumeType:         plan.VolumeType.ValueString(),
		VolumeSizeInKb:     plan.VolumeSizeInKb.ValueString(),
		Name:               plan.Name.ValueString(),
	}
	getSystems, _ := r.client.GetSystems()
	sr := goscaleio.NewSystem(r.client)
	sr.System = getSystems[0]
	getProtectionDomains, _ := sr.GetProtectionDomain("")
	tflog.Info(ctx, "2. [POWERFLEX] volume Resource State"+helper.PrettyJSON((getSystems[0])))
	for _, protectionDomain := range getProtectionDomains {
		pdr := goscaleio.NewProtectionDomain(r.client)
		pdr.ProtectionDomain = protectionDomain
		tflog.Info(ctx, "hello"+pdr.ProtectionDomain.ID+" "+plan.ProtectionDomainID.ValueString())
		if pdr.ProtectionDomain.ID == plan.ProtectionDomainID.ValueString() {
			getStoragePool, _ := pdr.GetStoragePool("")
			tflog.Info(ctx, "selected"+pdr.ProtectionDomain.ID+" "+plan.ProtectionDomainID.ValueString())
			for _, sp := range getStoragePool {
				spr := goscaleio.NewStoragePool(r.client)
				spr.StoragePool = sp
				tflog.Info(ctx, spr.StoragePool.ID+" "+plan.StoragePoolID.ValueString())
				if spr.StoragePool.ID == plan.StoragePoolID.ValueString() {
					tflog.Info(ctx, "selected : "+spr.StoragePool.ID+" "+plan.StoragePoolID.ValueString())
					volCreateResponse, err1 := spr.CreateVolume(volumeCreate)
					if err1 != nil {
						resp.Diagnostics.AddError(
							"Error creating volume",
							"Could not create volume, unexpected error: "+err1.Error(),
						)
						return
					}
					// plan.ID = types.StringValue(volCreateResponse.ID)
					volsResponse, err2 := spr.GetVolume("", volCreateResponse.ID, "", "", false)
					if err2 != nil {
						resp.Diagnostics.AddError(
							"Error getting volume after creation",
							"Could not get volume, unexpected error: "+err2.Error(),
						)
						return
					}
					tflog.Info(ctx, "[Volume] volume Resource State"+helper.PrettyJSON((volsResponse[0])))
					vol := volsResponse[0]
					spi := types.StringValue(vol.StoragePoolID)
					tflog.Info(ctx, "[Volume-SPI] volume Resource State"+spi.ValueString())

					plan = updateVolumeState(vol, plan)
					diags = resp.State.Set(ctx, plan)
					resp.Diagnostics.Append(diags...)
					if resp.Diagnostics.HasError() {
						return
					}
				}
			}
		}
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *volumeResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Get current state
	var state volumeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *volumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan volumeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *volumeResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state volumeResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *volumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
