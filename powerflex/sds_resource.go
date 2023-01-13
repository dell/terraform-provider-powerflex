package powerflex

import (
	"context"
	"fmt"

	scaleiotypes "github.com/dell/goscaleio/types/v1"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// SDSResourceModel maps the resource schema data.
type sdsResourceModel struct {
	ID                           types.String `tfsdk:"id"`
	Name                         types.String `tfsdk:"name"`
	ProtectionDomainID           types.String `tfsdk:"protection_domain_id"`
	ProtectionDomainName         types.String `tfsdk:"protection_domain_name"`
	IPList                       types.List   `tfsdk:"ip_list"`
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
	iplist := []*scaleiotypes.SdsIP{}
	var ipModellist []sdsIPModel
	plan.IPList.ElementsAs(ctx, &ipModellist, false)
	for _, v := range ipModellist {
		sdsIp := scaleiotypes.SdsIP{
			IP:   v.IP.ValueString(),
			Role: v.Role.ValueString(),
		}
		iplist = append(iplist, &sdsIp)
	}

	params := scaleiotypes.Sds{
		Name:   sdsName,
		IPList: iplist,
	}
	if !plan.RmcacheEnabled.IsNull() {
		params.RmcacheEnabled = plan.RmcacheEnabled.ValueBool()
	}
	if !plan.RmcacheSizeInKb.IsNull() {
		params.RmcacheSizeInKb = int(plan.RmcacheSizeInKb.ValueInt64())
	}
	if !plan.DrlMode.IsNull() {
		params.DrlMode = plan.DrlMode.ValueString()
	}
	if !plan.FaultSetID.IsNull() {
		params.FaultSetID = plan.FaultSetID.ValueString()
	}
	if !plan.Port.IsNull() {
		params.Port = int(plan.Port.ValueInt64())
	}
	// this is still not working for whatever reason
	// if !plan.NumOfIoBuffers.IsNull() {
	// 	params.NumOfIoBuffers = int(plan.NumOfIoBuffers.ValueInt64())
	// }
	sdsID, err2 := pdm.CreateSdsWithParams(&params)
	if err2 != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not create SDS with name %s and IP list %v and niobuff %d", sdsName, iplist, params.NumOfIoBuffers),
			err2.Error(),
		)
		return
	}

	// Get created SDS
	// rsp, err3 := pdm.FindSds("ID", sdsID)
	sdss, err3 := pdm.GetSds()
	// if err != nil {
	// 	return nil, err
	// }
	if err3 != nil {
		resp.Diagnostics.AddError(
			"Error getting all SDS after creation",
			err3.Error(),
		)
		return
	}
	var rsp *scaleiotypes.Sds
	found := false
	for _, sds := range sdss {
		if sds.ID == sdsID {
			rsp = &sds
			found = true
			break
		}
	}
	if !found {
		resp.Diagnostics.AddError(
			"No matching SDS",
			fmt.Sprintf("The SDS ID: %s", sdsID),
		)
		return
	}
	// resp.Diagnostics.AddError("1st Dummy", fmt.Sprintf("Sds created with IP:%s, role:%s", rsp.IPList[0].SdsIP.IP, rsp.IPList[0].SdsIP.Role))

	// Set refreshed state
	state, dgs := updateSdsState(rsp, plan)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
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
	state, dgs := updateSdsState(rsp, state)
	resp.Diagnostics.Append(dgs...)

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
	state, dgs := updateSdsState(rsp, plan)
	resp.Diagnostics.Append(dgs...)

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

	// Delete SDS
	err = pdm.DeleteSds(state.ID.ValueString())
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to delete Powerflex SDS",
			err.Error(),
		)

		return
	}

	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)

}

func (r *sdsResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	// Retrieve import ID and save to id attribute
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func updateSdsState(sds *scaleiotypes.Sds, plan sdsResourceModel) (sdsResourceModel, diag.Diagnostics) {
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

	IPAttrTypes := map[string]attr.Type{
		"ip":   types.StringType,
		"role": types.StringType,
	}
	IPElemType := types.ObjectType{
		AttrTypes: IPAttrTypes,
	}

	objectIPs := []attr.Value{}
	var diags diag.Diagnostics
	for _, ip := range sds.IPList {
		obj := map[string]attr.Value{
			"ip":   types.StringValue(ip.IP),
			"role": types.StringValue(ip.Role),
		}
		objVal, dgs := types.ObjectValue(IPAttrTypes, obj)
		diags = append(diags, dgs...)
		objectIPs = append(objectIPs, objVal)
	}
	setVal, dgs := types.ListValue(IPElemType, objectIPs)
	diags = append(diags, dgs...)
	state.IPList = setVal

	return state, diags
}
