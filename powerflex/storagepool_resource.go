package powerflex

import (
	"context"
	"terraform-provider-powerflex/helper"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &storagepoolResource{}
	_ resource.ResourceWithConfigure   = &storagepoolResource{}
	_ resource.ResourceWithImportState = &storagepoolResource{}
)

// StoragepoolResource - function to return resource interface
func StoragepoolResource() resource.Resource {
	return &storagepoolResource{}
}

type storagepoolResource struct {
	client *goscaleio.Client
}

type storagepoolResourceModel struct {
	LastUpdated        types.String `tfsdk:"last_updated"`
	ID                 types.String `tfsdk:"id"`
	SystemID           types.String `tfsdk:"systemid"`
	ProtectionDomainID types.String `tfsdk:"protection_domain_id"`
	Name               types.String `tfsdk:"name"`
	MediaType          types.String `tfsdk:"media_type"`
	UseRmcache         types.Bool   `tfsdk:"use_rmcache"`
	UseRfcache         types.Bool   `tfsdk:"use_rfcache"`
	Links              types.List   `tfsdk:"links"`
}

// Link - MODEL for Storagepool Links data returned by goscaleio.
type Link struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

func (r *storagepoolResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_storagepool"
}

func (r *storagepoolResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = StoragepoolReourceSchema
}

func (r *storagepoolResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*goscaleio.Client)
}

// Function used to Create Storagepool Resource
func (r *storagepoolResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "Create storagepool")
	// Retrieve values from plan
	var plan storagepoolResourceModel
	var pd *scaleiotypes.ProtectionDomain

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	s := goscaleio.NewSystem(r.client)
	systems, err := r.client.GetSystems()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error",
			"Could not get Systems, unexpected error: "+err.Error(),
		)
		return
	}
	s.System = systems[0]

	pd, err = s.FindProtectionDomain(plan.ProtectionDomainID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find protection domain",
			err.Error(),
		)
		return
	}
	p1 := goscaleio.NewProtectionDomain(r.client)
	p1.ProtectionDomain = pd

	sp, err := p1.CreateStoragePool(plan.Name.ValueString(), plan.MediaType.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error creating Storage Pool",
			"Could not create Storage Pool, unexpected error: "+err.Error(),
		)
		return
	}
	spResponse, err := p1.FindStoragePool(sp, "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Storagepool after creation",
			"Could not get Storagepool, unexpected error: "+err.Error(),
		)
		return
	}
	state := updateStoragepoolState(spResponse, plan)
	tflog.Debug(ctx, "Create Storagepool :-- "+helper.PrettyJSON(sp))
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Read Storagepool Resource
func (r *storagepoolResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "Read Storagepool")
	// Get current state
	var state storagepoolResourceModel
	var pd *scaleiotypes.ProtectionDomain

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	s := goscaleio.NewSystem(r.client)
	systems, err := r.client.GetSystems()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error",
			"Could not get Systems, unexpected error: "+err.Error(),
		)
		return
	}
	s.System = systems[0]

	pd, err = s.FindProtectionDomain(state.ProtectionDomainID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find protection domain",
			err.Error(),
		)
		return
	}
	p1 := goscaleio.NewProtectionDomain(r.client)
	p1.ProtectionDomain = pd
	spr, err := p1.FindStoragePool(state.ID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex Storagepool",
			err.Error(),
		)
		return
	}
	spResponse := updateStoragepoolState(spr, state)
	tflog.Debug(ctx, "Read Storagepool :-- "+helper.PrettyJSON(spr))
	diags = resp.State.Set(ctx, spResponse)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Update Storagepool Resource
func (r *storagepoolResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "Update Storagepool")
	// Retrieve values from plan
	var plan storagepoolResourceModel
	var pd *scaleiotypes.ProtectionDomain

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	//Get Current State
	var state storagepoolResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	s := goscaleio.NewSystem(r.client)
	systems, err := r.client.GetSystems()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error",
			"Could not get Systems, unexpected error: "+err.Error(),
		)
		return
	}
	s.System = systems[0]

	pd, err = s.FindProtectionDomain(plan.ProtectionDomainID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find protection domain",
			err.Error(),
		)
		return
	}
	p1 := goscaleio.NewProtectionDomain(r.client)
	p1.ProtectionDomain = pd

	if plan.Name.ValueString() != state.Name.ValueString() {
		_, err := p1.ModifyStoragePoolName(state.ID.ValueString(), plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting while updating Storagepool",
				"Could not get Storagepool, unexpected error: "+err.Error(),
			)
			return
		}
	}

	if plan.MediaType.ValueString() != state.MediaType.ValueString() {
		_, err := p1.ModifyStoragePoolMedia(state.ID.ValueString(), plan.MediaType.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error getting while updating Storagepool",
				"Could not get Storagepool, unexpected error: "+err.Error(),
			)
			return
		}
	}

	spResponse, err := p1.FindStoragePool(state.ID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting while updating Storagepool",
			"Could not get Storagepool, unexpected error: "+err.Error(),
		)
		return
	}

	state1 := updateStoragepoolState(spResponse, plan)
	tflog.Debug(ctx, "Update Storagepool :-- "+helper.PrettyJSON(spResponse))
	diags = resp.State.Set(ctx, state1)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Function used to Delete Storagepool Resource
func (r *storagepoolResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "Delete Storagepool")
	// Retrieve values from state
	var state storagepoolResourceModel
	var pd *scaleiotypes.ProtectionDomain

	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	stp, err := r.client.FindSystem(state.SystemID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex Systems Storagepool Create",
			err.Error(),
		)
		return
	}
	pd, err = stp.FindProtectionDomain(state.ProtectionDomainID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to find protection domain",
			err.Error(),
		)
		return
	}
	p1 := goscaleio.NewProtectionDomain(r.client)
	p1.ProtectionDomain = pd
	spr, err := r.client.FindStoragePool(state.ID.ValueString(), "", "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Find Powerflex Storagepool",
			err.Error(),
		)
		return
	}
	err = p1.DeleteStoragePool(state.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error Deleting Storagepool",
			"Couldn't Delete Storagepool "+err.Error(),
		)
		return
	}
	tflog.Debug(ctx, "Delete Storagepool :-- "+helper.PrettyJSON(spr))
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

// Function used to ImportState for Storagepool Resource
func (r *storagepoolResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Function to update the State for Storagepool Resource
func updateStoragepoolState(storagepool *scaleiotypes.StoragePool, plan storagepoolResourceModel) (state storagepoolResourceModel) {
	state.ProtectionDomainID = plan.ProtectionDomainID
	state.ID = types.StringValue(storagepool.ID)
	state.Name = types.StringValue(storagepool.Name)
	state.MediaType = types.StringValue(storagepool.MediaType)
	state.UseRmcache = types.BoolValue(storagepool.UseRmcache)
	state.UseRfcache = types.BoolValue(storagepool.UseRfcache)

	linkAttrTypes := map[string]attr.Type{
		"rel":  types.StringType,
		"href": types.StringType,
	}
	linkElemType := types.ObjectType{
		AttrTypes: linkAttrTypes,
	}

	objectLinks := []attr.Value{}
	for _, link := range storagepool.Links {
		obj := map[string]attr.Value{
			"rel":  types.StringValue(link.Rel),
			"href": types.StringValue(link.HREF),
		}
		objVal, _ := types.ObjectValue(linkAttrTypes, obj)
		objectLinks = append(objectLinks, objVal)
	}
	listVal, _ := types.ListValue(linkElemType, objectLinks)
	state.Links = listVal
	return state
}
