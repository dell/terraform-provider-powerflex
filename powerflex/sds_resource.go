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
	IPList                       types.Set    `tfsdk:"ip_list"`
	Port                         types.Int64  `tfsdk:"port"`
	SdsState                     types.String `tfsdk:"sds_state"`
	MembershipState              types.String `tfsdk:"membership_state"`
	MdmConnectionState           types.String `tfsdk:"mdm_connection_state"`
	DrlMode                      types.String `tfsdk:"drl_mode"`
	RmcacheEnabled               types.Bool   `tfsdk:"rmcache_enabled"`
	RmcacheSizeInMB              types.Int64  `tfsdk:"rmcache_size_in_mb"`
	RfcacheEnabled               types.Bool   `tfsdk:"rfcache_enabled"`
	RmcacheFrozen                types.Bool   `tfsdk:"rmcache_frozen"`
	IsOnVMware                   types.Bool   `tfsdk:"is_on_vmware"`
	FaultSetID                   types.String `tfsdk:"fault_set_id"`
	NumOfIoBuffers               types.Int64  `tfsdk:"num_of_io_buffers"`
	RmcacheMemoryAllocationState types.String `tfsdk:"rmcache_memory_allocation_state"`
	PerformanceProfile           types.String `tfsdk:"performance_profile"`
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

// Conversion of list of IPs from tf model to go type
func (sds *sdsResourceModel) getIPList(ctx context.Context) []*scaleiotypes.SdsIP {
	iplist := []*scaleiotypes.SdsIP{}
	var ipModellist []sdsIPModel
	sds.IPList.ElementsAs(ctx, &ipModellist, false)
	for _, v := range ipModellist {
		sdsIP := scaleiotypes.SdsIP{
			IP:   v.IP.ValueString(),
			Role: v.Role.ValueString(),
		}
		iplist = append(iplist, &sdsIP)
	}
	return iplist
}

// Get difference between sets of IP in state and plan
func sdsIPListDiff(ctx context.Context, plan, state *sdsResourceModel) (toAdd, toRmv, changed, common []*scaleiotypes.SdsIP) {
	plist, slist := plan.getIPList(ctx), state.getIPList(ctx)
	type ipObj struct {
		pip *scaleiotypes.SdsIP
		sip *scaleiotypes.SdsIP
	}
	vmap := make(map[string]*ipObj)
	for _, pip := range plist {
		vmap[pip.IP] = &ipObj{pip, nil}
	}
	for _, sip := range slist {
		if mip, ok := vmap[sip.IP]; ok {
			mip.sip = sip
		} else {
			vmap[sip.IP] = &ipObj{nil, sip}
		}
	}
	toAdd, toRmv, common, changed = make([]*scaleiotypes.SdsIP, 0), make([]*scaleiotypes.SdsIP, 0),
		make([]*scaleiotypes.SdsIP, 0), make([]*scaleiotypes.SdsIP, 0)
	for _, mip := range vmap {
		if mip.sip != nil {
			if mip.pip != nil {
				if mip.pip.Role == mip.sip.Role {
					common = append(common, mip.pip)
				} else {
					changed = append(changed, mip.pip)
				}
			} else {
				toRmv = append(toRmv, mip.sip)
			}
		} else {
			toAdd = append(toAdd, mip.pip)
		}
	}
	return toAdd, toRmv, changed, common
}

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

	// set the protection domain name in the plan so that it gets propagated to the state
	plan.ProtectionDomainName = types.StringValue(pdm.ProtectionDomain.Name)

	sdsName := plan.Name.ValueString()
	iplist := plan.getIPList(ctx)

	params := scaleiotypes.Sds{
		Name:   sdsName,
		IPList: iplist,
	}
	if !plan.RmcacheEnabled.IsUnknown() {
		params.RmcacheEnabled = plan.RmcacheEnabled.ValueBool()
	}
	if !plan.RmcacheSizeInMB.IsUnknown() {
		params.RmcacheSizeInKb = int(plan.RmcacheSizeInMB.ValueInt64()) * 1024
	}
	if !plan.DrlMode.IsUnknown() {
		params.DrlMode = plan.DrlMode.ValueString()
	}
	if !plan.Port.IsUnknown() {
		params.Port = int(plan.Port.ValueInt64())
	}
	sdsID, err2 := pdm.CreateSdsWithParams(&params)
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
	state, dgs := updateSdsState(rsp, plan)
	resp.Diagnostics.Append(dgs...)
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)

	if !plan.RfcacheEnabled.IsUnknown() {
		rfCacheEnabled := plan.RfcacheEnabled.ValueBool()
		err := pdm.SetSdsRfCache(sdsID, rfCacheEnabled)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set SDS Rf Cache settings to %t", rfCacheEnabled),
				err.Error(),
			)
		}
	}

	if !plan.PerformanceProfile.IsUnknown() {
		perfprof := plan.PerformanceProfile.ValueString()
		err := pdm.SetSdsPerformanceProfile(sdsID, perfprof)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set SDS Performance Profile settings to %s", perfprof),
				err.Error(),
			)
		}
	}

	// Get updated SDS
	rsp, err4 := pdm.FindSds("ID", sdsID)
	if err4 != nil {
		resp.Diagnostics.AddError(
			"Error getting SDS after setting Rf cache and Performance Profile",
			err4.Error(),
		)
		return
	}
	// Set refreshed state
	state, dgs = updateSdsState(rsp, plan)
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

	// Get the system on the PowerFlex cluster
	system, err := getFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}

	// Get SDS
	var rsp scaleiotypes.Sds
	if rsp, err = system.GetSdsByID(state.ID.ValueString()); err != nil {
		resp.Diagnostics.AddError(
			fmt.Sprintf("Could not get SDS by ID %s", state.ID.ValueString()),
			err.Error(),
		)
		return
	}

	// when SDS is imported, protection domain name is not known and this causes a non empty plan
	if state.ProtectionDomainName.IsNull() {
		protectionDomain, err := system.FindProtectionDomain(rsp.ProtectionDomainID, "", "")
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Unable to read name of protection domain of ID %s for SDS %s", rsp.ProtectionDomainID, rsp.Name),
				err.Error(),
			)
		} else {
			state.ProtectionDomainName = types.StringValue(protectionDomain.Name)
		}
	}

	// Set refreshed state
	state, dgs := updateSdsState(&rsp, state)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
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
		}
	}

	// Check if there are updates in ip lists
	// Stop updating IPs if one IP updation fails
	func() {
		toAdd, toRmv, changed, _ := sdsIPListDiff(ctx, &plan, &state)
		for _, ip := range toAdd {
			err := pdm.AddSdSIP(state.ID.ValueString(), ip.IP, ip.Role)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error adding IP %s to SDS with role %s", ip.IP, ip.Role),
					err.Error(),
				)
				return
			}
		}
		for _, ip := range changed {
			err := pdm.SetSDSIPRole(state.ID.ValueString(), ip.IP, ip.Role)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error updating IP %s role to %s in SDS", ip.IP, ip.Role),
					err.Error(),
				)
				return
			}
		}
		for _, ip := range toRmv {
			err := pdm.RemoveSDSIP(state.ID.ValueString(), ip.IP)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error removing IP %s with role %s from SDS", ip.IP, ip.Role),
					err.Error(),
				)
				return
			}
		}
	}()

	// check if there is change in sds port
	if !plan.Port.IsUnknown() && !state.Port.Equal(plan.Port) {
		port := plan.Port.ValueInt64()
		err := pdm.SetSdsPort(state.ID.ValueString(), int(port))
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not change SDS port to %d", port),
				err.Error(),
			)
		}
	}

	// check if there is change in sds drl mode
	if !plan.DrlMode.IsUnknown() && !state.DrlMode.Equal(plan.DrlMode) {
		drlMode := plan.DrlMode.ValueString()
		err := pdm.SetSdsDrlMode(state.ID.ValueString(), drlMode)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not change SDS DRL Mode to %s", drlMode),
				err.Error(),
			)
		}
	}

	// check if there is change in sds rmcache
	if !plan.RmcacheEnabled.IsUnknown() && !state.RmcacheEnabled.Equal(plan.RmcacheEnabled) {
		rmCacheEnabled := plan.RmcacheEnabled.ValueBool()
		err := pdm.SetSdsRmCache(state.ID.ValueString(), rmCacheEnabled)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not change SDS Read Ram Cache settings to %t", rmCacheEnabled),
				err.Error(),
			)
		}
	}
	if !plan.RmcacheSizeInMB.IsUnknown() && !state.RmcacheSizeInMB.Equal(plan.RmcacheSizeInMB) {
		rmCacheSize := plan.RmcacheSizeInMB.ValueInt64()
		err := pdm.SetSdsRmCacheSize(state.ID.ValueString(), int(rmCacheSize))
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not change SDS Read Ram Cache size to %d", rmCacheSize),
				err.Error(),
			)
		}
	}

	// check if there is change in sds rfcache
	if !plan.RfcacheEnabled.IsUnknown() && !state.RfcacheEnabled.Equal(plan.RfcacheEnabled) {
		rfCacheEnabled := plan.RfcacheEnabled.ValueBool()
		err := pdm.SetSdsRfCache(state.ID.ValueString(), rfCacheEnabled)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not change SDS Rf Cache settings to %t", rfCacheEnabled),
				err.Error(),
			)
		}
	}

	// Check if performance profile has been changed
	if !plan.PerformanceProfile.IsUnknown() && !state.PerformanceProfile.Equal(plan.PerformanceProfile) {
		perfprof := plan.PerformanceProfile.ValueString()
		err := pdm.SetSdsPerformanceProfile(state.ID.ValueString(), perfprof)
		if err != nil {
			resp.Diagnostics.AddError(
				fmt.Sprintf("Could not set SDS Performance Profile settings to %s", perfprof),
				err.Error(),
			)
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
	state, dgs := updateSdsState(rsp, state)
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
	state.ProtectionDomainID = types.StringValue(sds.ProtectionDomainID)
	state.Port = types.Int64Value(int64(sds.Port))
	state.SdsState = types.StringValue(sds.SdsState)
	state.MembershipState = types.StringValue(sds.MembershipState)
	state.MdmConnectionState = types.StringValue(sds.MdmConnectionState)
	state.DrlMode = types.StringValue(sds.DrlMode)
	state.RmcacheEnabled = types.BoolValue(sds.RmcacheEnabled)
	state.RmcacheSizeInMB = types.Int64Value(int64(sds.RmcacheSizeInKb) / 1024)
	state.RfcacheEnabled = types.BoolValue(sds.RfcacheEnabled)
	state.RmcacheFrozen = types.BoolValue(sds.RmcacheFrozen)
	state.IsOnVMware = types.BoolValue(sds.IsOnVMware)
	state.FaultSetID = types.StringValue(sds.FaultSetID)
	state.NumOfIoBuffers = types.Int64Value(int64(sds.NumOfIoBuffers))
	state.RmcacheMemoryAllocationState = types.StringValue(sds.RmcacheMemoryAllocationState)
	state.PerformanceProfile = types.StringValue(sds.PerformanceProfile)

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
	setVal, dgs := types.SetValue(IPElemType, objectIPs)
	diags = append(diags, dgs...)
	state.IPList = setVal

	return state, diags
}
