package powerflex

import (
	"context"
	"fmt"

	scaleiotypes "github.com/dell/goscaleio/types/v1"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SDSResourceModel maps the resource schema data.
type sdsResourceModel struct {
	ID                   types.String `tfsdk:"id"`
	Name                 types.String `tfsdk:"name"`
	ProtectionDomainID   types.String `tfsdk:"protection_domain_id"`
	ProtectionDomainName types.String `tfsdk:"protection_domain_name"`
	IPList               types.List   `tfsdk:"ip_list"`
	// []sdsIPModel
	Port                         types.Int64  `tfsdk:"port"`
	SdsState                     types.String `tfsdk:"sds_state"`
	MembershipState              types.String `tfsdk:"membership_state"`
	MdmConnectionState           types.String `tfsdk:"mdm_connection_state"`
	DrlMode                      types.String `tfsdk:"drl_mode"`
	RmcacheEnabled               types.Bool   `tfsdk:"rmcache_enabled"`
	RmcacheSizeInKb              types.Int64  `tfsdk:"rmcache_size_in_kb"`
	RmcacheFrozen                types.Bool   `tfsdk:"rmcache_frozen"`
	IsOnVMware                   types.Bool   `tfsdk:"is_on_vmware"`
	FaultSetID                   types.String `tfsdk:"fault_set_id"`
	NumOfIoBuffers               types.Int64  `tfsdk:"num_of_io_buffers"`
	RmcacheMemoryAllocationState types.String `tfsdk:"rmcache_memory_allocation_state"`
}

// SDS IP object
type sdsIPModel struct {
	IP   types.String `tfsdk:"ip"`
	Role types.String `tfsdk:"role"`
}

var (
	_ resource.Resource                = &sdsResource{}
	_ resource.ResourceWithConfigure   = &sdsResource{}
	_ resource.ResourceWithImportState = &sdsResource{}
)

// NewSDSResource is a helper function to simplify the provider implementation.
func NewSDSResource() resource.Resource {
	return &sdsResource{}
}

// sdsResource is the resource implementation.
type sdsResource struct {
	client *goscaleio.Client
}

func (r *sdsResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sds"
}

func (r *sdsResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SDSResourceSchema
}

func (r *sdsResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*goscaleio.Client)
}

// func getIPList(ctx context.Context, sds sdsResourceModel) []scaleiotypes.SdsIP {
// 	iplist := []scaleiotypes.SdsIP{}
// 	var ipModellist []sdsIPModel
// 	sds.IPList.ElementsAs(ctx, &ipModellist, false)
// 	for _, v := range ipModellist {
// 		sdsIp := scaleiotypes.SdsIP{
// 			IP:   v.IP.ValueString(),
// 			Role: v.Role.ValueString(),
// 		}
// 		iplist = append(iplist, sdsIp)
// 	}
// 	return iplist
// }

// func marshallIPList(ctx context.Context, ips []scaleiotypes.SdsIP, sds *sdsResourceModel) {
// 	ipModellist := []sdsIPModel{}
// 	for _, ip := range ips {
// 		ipModellist = append(ipModellist, sdsIPModel{IP: types.StringValue(ip.IP), Role: types.StringValue(ip.Role)})
// 	}
// 	sds.IPList = basetypes.NewListValueMust(sds.IPList.Type(), ipModellist)
// }

// Create creates the resource and sets the initial Terraform state.
func (r *sdsResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan sdsResourceModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pdm, err := getNewProtectionDomainEx(r.client, plan.ProtectionDomainID.ValueString(), plan.ProtectionDomainName.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			err.Error(),
		)
		return
	}

	sdsName := plan.Name.ValueString()
	iplist := []scaleiotypes.SdsIP{}
	var ipModellist []sdsIPModel
	plan.IPList.ElementsAs(ctx, &ipModellist, false)
	for _, v := range ipModellist {
		sdsIp := scaleiotypes.SdsIP{
			IP:   v.IP.ValueString(),
			Role: v.Role.ValueString(),
		}
		iplist = append(iplist, sdsIp)
	}

	// iplist := getIPList(ctx, plan)

	// Create SDS
	sdsID, err2 := pdm.CreateSdsWithIPRole(sdsName, iplist)
	if err2 != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create SDS with name %s and IP list %v", sdsName, iplist),
			err2.Error(),
		)
		return
	}

	// Get created SDS
	rsp, err3 := pdm.FindSds("ID", sdsID)
	if err3 != nil {
		resp.Diagnostics.AddError(
			"Error getting SDS after creation",
			err3.Error(),
		)
		return
	}

	// Set refreshed state
	state := updateSdsState(rsp, plan)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Read refreshes the Terraform state with the latest data.
func (r *sdsResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state sdsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pdm, err := getNewProtectionDomainEx(r.client, state.ProtectionDomainID.ValueString(), state.ProtectionDomainName.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			err.Error(),
		)
		return
	}

	// Get SDS
	var rsp *scaleiotypes.Sds
	if rsp, err = pdm.FindSds("ID", state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			"Could not get SDS",
			err.Error(),
		)
		return
	}

	// Set refreshed state
	state = updateSdsState(rsp, state)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sdsResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan sdsResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	var state sdsResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pdm, err := getNewProtectionDomainEx(r.client, state.ProtectionDomainID.ValueString(), state.ProtectionDomainName.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			err.Error(),
		)
		return
	}

	// Check if there difference between plan and state
	if plan.Name.ValueString() != state.Name.ValueString() {
		err := pdm.SetSdsName(state.ID.ValueString(), plan.Name.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Could not rename SDS",
				err.Error(),
			)

			return
		}
	}

	// Find updated SDS
	rsp, err := pdm.FindSds("ID", state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting SDS after updation",
			err.Error(),
		)
		return
	}

	// Set refreshed state
	state = updateSdsState(rsp, plan)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sdsResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state sdsResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	pdm, err := getNewProtectionDomainEx(r.client, state.ProtectionDomainID.ValueString(), state.ProtectionDomainName.ValueString(), "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting Protection Domain",
			err.Error(),
		)
		return
	}

	// Find SDS
	sds, err := pdm.FindSds("ID", state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Find Powerflex SDS",
			err.Error(),
		)

		return
	}

	// Delete SDS
	err = pdm.DeleteSds(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Powerflex SDS",
			err.Error(),
		)

		return
	}

	// Set state
	state = updateSdsState(sds, state)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *sdsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func updateSdsState(sds *scaleiotypes.Sds, plan sdsResourceModel) sdsResourceModel {
	state := plan
	state.ID = types.StringValue(sds.ID)
	state.Name = types.StringValue(sds.Name)
	state.Port = types.Int64Value(int64(sds.Port))
	state.SdsState = types.StringValue(sds.SdsState)
	state.MembershipState = types.StringValue(sds.MembershipState)
	state.MdmConnectionState = types.StringValue(sds.MdmConnectionState)
	state.DrlMode = types.StringValue(sds.DrlMode)
	state.RmcacheEnabled = types.BoolValue(sds.RmcacheEnabled)
	state.RmcacheSizeInKb = types.Int64Value(int64(sds.RmcacheSizeInKb))
	state.RmcacheFrozen = types.BoolValue(sds.RmcacheFrozen)
	state.IsOnVMware = types.BoolValue(sds.IsOnVMware)
	state.FaultSetID = types.StringValue(sds.FaultSetID)
	state.NumOfIoBuffers = types.Int64Value(int64(sds.NumOfIoBuffers))
	state.RmcacheMemoryAllocationState = types.StringValue(sds.RmcacheMemoryAllocationState)

	return state
}
