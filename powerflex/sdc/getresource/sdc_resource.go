package getresource

import (
	"context"
	"fmt"
	"terraform-provider-powerflex/helper"
	"time"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &sdcResource{}
	_ resource.ResourceWithConfigure   = &sdcResource{}
	_ resource.ResourceWithImportState = &sdcResource{}
)

// SDCResource - function to return resource interface
func SDCResource() resource.Resource {
	return &sdcResource{}
}

type sdcResource struct {
	client *goscaleio.Client
}

type sdcResourceModel struct {
	// Sdcs        []sdcModel   `tfsdk:"sdcs"`
	LastUpdated        types.String   `tfsdk:"last_updated"`
	SdcID              types.String   `tfsdk:"id"`
	SystemID           types.String   `tfsdk:"systemid"`
	Name               types.String   `tfsdk:"name"`
	SdcIP              types.String   `tfsdk:"sdcip"`
	SdcApproved        types.Bool     `tfsdk:"sdcapproved"`
	OnVMWare           types.Bool     `tfsdk:"onvmware"`
	SdcGUID            types.String   `tfsdk:"sdcguid"`
	MdmConnectionState types.String   `tfsdk:"mdmconnectionstate"`
	Links              []sdcLinkModel `tfsdk:"links"`
}

// sdcLinkModel - MODEL for SDC Links data returned by goscaleio.
type sdcLinkModel struct {
	Rel  types.String `tfsdk:"rel"`
	HREF types.String `tfsdk:"href"`
}

func (r *sdcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc"
}

func (r *sdcResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SDCReourceSchema
}

func (r *sdcResource) Configure(_ context.Context, req resource.ConfigureRequest, _ *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	r.client = req.ProviderData.(*goscaleio.Client)
}

func (r *sdcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "[ANSHU] Create")
	// Retrieve values from plan
	var plan sdcResourceModel
	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	system, err := r.client.FindSystem(plan.SystemID.ValueString(), "", "")

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex systems sdcs Create",
			err.Error(),
		)
		return
	}
	nameChng, err := system.ChangeSdcName(plan.SdcID.ValueString(), plan.Name.ValueString())

	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Change name Powerflex sdc",
			err.Error(),
		)
		return
	}

	sdcs, err := system.GetSdc()
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex sdcs",
			err.Error(),
		)
		return
	}

	finalSDC := findChangedSdc(sdcs, plan.SdcID.ValueString())
	plan = getSdcState(finalSDC)

	tflog.Debug(ctx, "[ANSHU] plan getSdcState plan"+helper.PrettyJSON(plan))
	tflog.Debug(ctx, "[ANSHU] nameChng Result :-- "+helper.PrettyJSON(nameChng))
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sdcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "[ANSHU] Read")
	// Get current state
	var state sdcResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	system, err := r.client.FindSystem(state.SystemID.ValueString(), "", "")
	singleSdc, err := system.FindSdc("id", state.SdcID.ValueString())
	fmt.Println(singleSdc)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex systems sdcs Read",
			err.Error(),
		)
		return
	}
	// state = getSdcState(*singleSdc)
	tflog.Debug(ctx, "[ANSHUM] state return"+helper.PrettyJSON(state))
	// Set refreshed state
	diags = resp.State.Set(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sdcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "[ANSHU] Update")
	// Retrieve values from plan
	var plan sdcResourceModel
	diags := req.Plan.Get(ctx, &plan)

	system, err := r.client.FindSystem(plan.SystemID.ValueString(), "", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex systems sdcs Update",
			err.Error(),
		)
		return
	}
	nameChng, err := system.ChangeSdcName(plan.SdcID.ValueString(), plan.Name.ValueString())

	fmt.Println(nameChng)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Change name Powerflex sdc",
			err.Error(),
		)
		return
	}
	// plan = getSdcState(*nameChng)
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	// Set state to fully populated data
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sdcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "[ANSHU] Delete")
	// Retrieve values from state
	var state sdcResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

}

func (r *sdcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// getSdcState - function to return all sdc result from goscaleio.
func getSdcState(sdc scaleiotypes.Sdc) (response sdcResourceModel) {
	// var basenameOpts []sdcModel = []sdcModel{}
	pln := sdcResourceModel{
		SdcID:              types.StringValue(sdc.ID),
		Name:               types.StringValue(sdc.Name),
		SdcGUID:            types.StringValue(sdc.SdcGUID),
		SdcApproved:        types.BoolValue(sdc.SdcApproved),
		OnVMWare:           types.BoolValue(sdc.OnVMWare),
		SystemID:           types.StringValue(sdc.SystemID),
		SdcIP:              types.StringValue(sdc.SdcIP),
		MdmConnectionState: types.StringValue(sdc.MdmConnectionState),
	}

	plnLinks := []sdcLinkModel{}

	for _, link := range sdc.Links {
		plnLinks = append(plnLinks, sdcLinkModel{
			Rel:  types.StringValue(link.Rel),
			HREF: types.StringValue(link.HREF),
		})
	}

	return pln
}

func findChangedSdc(sdcs []scaleiotypes.Sdc, id string) scaleiotypes.Sdc {
	var sdcReturnValue scaleiotypes.Sdc
	for _, sdcValue := range sdcs {

		if id == sdcValue.ID {
			sdcReturnValue = sdcValue
		}

	}
	return sdcReturnValue
}
