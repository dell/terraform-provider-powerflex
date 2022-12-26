package powerflex

import (
	"context"
	"strconv"

	"github.com/dell/goscaleio"
	pftypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// Ensure the implementation satisfies the expected interfaces.
var (
	_ resource.Resource                = &volumeResource{}
	_ resource.ResourceWithConfigure   = &volumeResource{}
	_ resource.ResourceWithImportState = &volumeResource{}
)

// NewVolumeResource is a helper function to simplify the provider implementation.
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
	if plan.Size.ValueInt64()%8 != 0 {
		resp.Diagnostics.AddError(
			"Error: Size Must be in granularity of 8GB",
			"Could not assign volume with size. sizeInGb ("+strconv.FormatInt(plan.Size.ValueInt64(), 10)+") must be a positive number in granularity of 8 GB.",
		)
		return
	}
	VSIKB, err := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error: Invalid Capacity unit :"+plan.CapacityUnit.String(),
			err.Error(),
		)
		return
	}
	volumeCreate := &pftypes.VolumeParam{
		ProtectionDomainID: plan.ProtectionDomainID.ValueString(),
		StoragePoolID:      plan.StoragePoolID.ValueString(),
		UseRmCache:         strconv.FormatBool(plan.UseRmCache.ValueBool()),
		VolumeType:         plan.VolumeType.ValueString(),
		VolumeSizeInKb:     strconv.FormatInt(VSIKB, 10),
		Name:               plan.Name.ValueString(),
	}
	spr, _ := getStoragePoolInstance(r.client, volumeCreate.StoragePoolID, volumeCreate.ProtectionDomainID)
	volCreateResponse, err1 := spr.CreateVolume(volumeCreate)
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error creating volume",
			"Could not create volume, unexpected error: "+err1.Error(),
		)
		return
	}
	volsResponse, err2 := spr.GetVolume("", volCreateResponse.ID, "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume after creation",
			"Could not get volume, unexpected error: "+err2.Error(),
		)
		return
	}
	vol := volsResponse[0]
	vr := goscaleio.NewVolume(r.client)
	vr.Volume = vol
	msids := []string{}
	diags = plan.MapSdcsID.ElementsAs(ctx, &msids, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}
	for _, msid := range msids {
		// Add mapped SDC
		pfmvsp := pftypes.MapVolumeSdcParam{
			SdcID:                 msid,
			AllowMultipleMappings: "true",
		}
		err3 := vr.MapVolumeSdc(&pfmvsp)
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error Mapping Volume to SDCs",
				"Could not map volume to scs with id: "+msid+", unexpected error: "+err3.Error(),
			)
			return
		}
	}
	if plan.LockedAutoSnapshot.ValueBool() {
		err := vr.LockAutoSnapshot()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Locking Auto Snapshots",
				"Could not lock auto snapshots, unexpected error: "+err.Error(),
			)
		}
	}
	volsResponse, err2 = spr.GetVolume("", volCreateResponse.ID, "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume after mapping the sdcs",
			"Could not get volume after mapping the sdcs, unexpected error: "+err2.Error(),
		)
		return
	}
	vol = volsResponse[0]
	state := VolumeTerraformState(vol, plan)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
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
	spr, err1 := getStoragePoolInstance(r.client, state.StoragePoolID.ValueString(), state.ProtectionDomainID.ValueString())
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error getting storage pool",
			"Could not get storage pool, unexpected err: "+err1.Error(),
		)
		return
	}
	volsResponse, err2 := spr.GetVolume("", state.ID.ValueString(), "", "", false)
	if err2 != nil {
		resp.Diagnostics.AddError(
			"Error getting volume",
			"Could not get volume, unexpected error: "+err2.Error(),
		)
		return
	}
	vol := volsResponse[0]
	state = VolumeTerraformState(vol, state)
	// Set refreshed state
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *volumeResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values
	var plan volumeResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state
	var state volumeResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	VSIKB, _ := convertToKB(plan.CapacityUnit.ValueString(), plan.Size.ValueInt64())
	plan.VolumeSizeInKb = types.StringValue(strconv.FormatInt(VSIKB, 10))

	spr, err1 := getStoragePoolInstance(r.client, state.StoragePoolID.ValueString(), state.ProtectionDomainID.ValueString())
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error getting storage pool",
			"Could not get storage pool, unexpected err: "+err1.Error(),
		)
		return
	}
	volsplan, _ := spr.GetVolume("", state.ID.ValueString(), "", "", false)
	volresource := goscaleio.NewVolume(r.client)
	volresource.Volume = volsplan[0]

	// updating the name of volume if there is change in plan
	if plan.Name.ValueString() != state.Name.ValueString() {
		err_rename := volresource.SetVolumeName(plan.Name.ValueString())
		if err_rename != nil {
			resp.Diagnostics.AddError(
				"Error renaming the volume -> "+plan.Name.ValueString()+" : "+state.Name.ValueString(),
				"Could not rename the volume, unexpected error:"+err_rename.Error(),
			)
			return
		}
	}

	// updating the size of the volume if there is change in plan
	if plan.VolumeSizeInKb.ValueString() != state.VolumeSizeInKb.ValueString() {
		sizeInGb, _ := strconv.Atoi(strconv.FormatInt(VSIKB, 10))
		sizeInGb = sizeInGb / 1048576
		sizeInGB := strconv.FormatInt(int64(sizeInGb), 10)
		// sizeInGb = ((sizeInGb / 8) + 1) * 8
		// newSizeIn8Gb := strconv.FormatInt(int64(sizeInGb), 10)
		err3 := volresource.SetVolumeSize(sizeInGB)
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error setting volume size -> "+plan.VolumeSizeInKb.ValueString()+":"+state.VolumeSizeInKb.ValueString(),
				"Could not set new volume size -> "+sizeInGB+", unexpected err: "+err3.Error(),
			)
			return
		}
	}
	planSdcIds := []string{}
	stateSdcIds := []string{}
	diags = plan.MapSdcsID.ElementsAs(ctx, &planSdcIds, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}
	diags = state.MapSdcsID.ElementsAs(ctx, &stateSdcIds, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}
	mapSdcIds := Difference(planSdcIds, stateSdcIds)
	unmapSdcIds := Difference(stateSdcIds, planSdcIds)

	for _, msi := range mapSdcIds {
		pfmvsp := pftypes.MapVolumeSdcParam{
			SdcID:                 msi,
			AllowMultipleMappings: "true",
		}
		err3 := volresource.MapVolumeSdc(&pfmvsp)
		if err3 != nil {
			resp.Diagnostics.AddError(
				"Error Mapping Volume to SDCs",
				"Could map volume to scs with id: "+msi+", unexpected error: "+err3.Error(),
			)
			return
		}
	}

	for _, usi := range unmapSdcIds {
		err4 := volresource.UnmapVolumeSdc(
			&pftypes.UnmapVolumeSdcParam{
				SdcID: usi,
			},
		)
		if err4 != nil {
			resp.Diagnostics.AddError(
				"Error Unmapping Volume to SDCs",
				"Could Unmap volume to scs with id: "+usi+", unexpected error: "+err4.Error(),
			)
			return
		}
	}
	if plan.LockedAutoSnapshot.ValueBool() && !state.LockedAutoSnapshot.ValueBool() {
		err := volresource.LockAutoSnapshot()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Locking Auto Snapshots",
				"Could not lock auto snapshots, unexpected error: "+err.Error(),
			)
		}
	}
	if !plan.LockedAutoSnapshot.ValueBool() && state.LockedAutoSnapshot.ValueBool() {
		err := volresource.UnlockAutoSnapshot()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Unlocking Auto Snapshots",
				"Could not unlock auto snapshots, unexpected error: "+err.Error(),
			)
		}
	}
	vols, _ := spr.GetVolume("", state.ID.ValueString(), "", "", false)
	state = VolumeTerraformState(vols[0], plan)
	// Set refreshed state
	diags = resp.State.Set(ctx, state)
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
	spr, err1 := getStoragePoolInstance(r.client, state.StoragePoolID.ValueString(), state.ProtectionDomainID.ValueString())
	if err1 != nil {
		resp.Diagnostics.AddError(
			"Error getting storage pool",
			"Could not get storage pool, unexpected err: "+err1.Error(),
		)
		return
	}
	volsplan, _ := spr.GetVolume("", state.ID.ValueString(), "", "", false)
	volresource := goscaleio.NewVolume(r.client)
	volresource.Volume = volsplan[0]
	sdcsToUnmap := []string{}
	diags = state.MapSdcsID.ElementsAs(ctx, &sdcsToUnmap, true)
	if diags.HasError() {
		resp.Diagnostics.Append(diags...)
	}
	for _, stu := range sdcsToUnmap {
		err := volresource.UnmapVolumeSdc(
			&pftypes.UnmapVolumeSdcParam{
				SdcID: stu,
			},
		)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error Unmapping Volume to SDCs",
				"Couldn't unmap volume to scs with id: "+stu+", unexpected error: "+err.Error(),
			)
			return
		}
	}
	err := volresource.RemoveVolume("")
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

func (r *volumeResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
