package getresource

import (
	"context"
	"time"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

var (
	_ resource.Resource                = &sdcResource{}
	_ resource.ResourceWithConfigure   = &sdcResource{}
	_ resource.ResourceWithImportState = &sdcResource{}
)

func SDCResource() resource.Resource {
	return &sdcResource{}
}

type sdcResource struct {
	client *goscaleio.Client
}

type sdcResourceModel struct {
	Sdcs        []sdcModel   `tfsdk:"sdcs"`
	ID          types.String `tfsdk:"sdcid"`
	SystemID    types.String `tfsdk:"systemid"`
	Name        types.String `tfsdk:"name"`
	LastUpdated types.String `tfsdk:"last_updated"`
}

// sdcModel - MODEL for SDC data returned by goscaleio.
type sdcModel struct {
	ID                 types.String   `tfsdk:"id"`
	SystemID           types.String   `tfsdk:"systemid"`
	SdcIP              types.String   `tfsdk:"sdcip"`
	SdcApproved        types.Bool     `tfsdk:"sdcapproved"`
	OnVMWare           types.Bool     `tfsdk:"onvmware"`
	SdcGUID            types.String   `tfsdk:"sdcguid"`
	MdmConnectionState types.String   `tfsdk:"mdmconnectionstate"`
	Name               types.String   `tfsdk:"name"`
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
	// Retrieve values from plan
	var plan sdcResourceModel
	diags := req.Plan.Get(ctx, &plan)
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
	// Get current state
	var state sdcResourceModel
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

func (r *sdcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var plan sdcResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	plan.LastUpdated = types.StringValue(time.Now().Format(time.RFC850))

	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *sdcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
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
