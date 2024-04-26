package provider

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"terraform-provider-powerflex/client"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &sdcHostResource{}
	_ resource.ResourceWithConfigure = &sdcHostResource{}
	// _ resource.ResourceWithImportState = &sdcHostResource{}
	_ client.Logger = &provisionerLogger{}
)

type provisionerLogger struct {
	ctx context.Context
}

func (l *provisionerLogger) Printf(format string, v ...any) {
	tflog.Info(l.ctx, fmt.Sprintf(format, v...))
}

func (l *provisionerLogger) Println(v ...any) {
	tflog.Info(l.ctx, fmt.Sprint(v...))
}

// NewSDCResource is a helper function to simplify the provider implementation.
func NewSDCHostResource() resource.Resource {
	return &sdcHostResource{}
}

// sdcHostResource is the resource implementation.
type sdcHostResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (r *sdcHostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc_host"
}

func (r *sdcHostResource) ConfigValidators(ctx context.Context) []resource.ConfigValidator {
	remotePath := path.MatchRoot("remote")
	return []resource.ConfigValidator{
		resourcevalidator.ExactlyOneOf(
			remotePath.AtName("password"),
			remotePath.AtName("private_key"),
		),
		// TODO: Add CA Cert validation
		// resourcevalidator.Conflicting(
		// 	remotePath.AtName("password"),
		// 	remotePath.AtName("ca_cert"),
		// ),
	}
}

func (r *sdcHostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource is used to manage the Storage Data Servers entity of PowerFlex Array. We can Create, Update and Delete the SDC using this resource. We can also import an existing SDC from PowerFlex array.",
		MarkdownDescription: "This resource is used to manage the Storage Data Servers entity of PowerFlex Array. We can Create, Update and Delete the SDC using this resource. We can also import an existing SDC from PowerFlex array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The id of the SDC",
				Computed:            true,
				MarkdownDescription: "The id of the SDC",
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"ip": schema.StringAttribute{
				Description:         "IP address of the server to be coonfigured as SDC.",
				Required:            true,
				MarkdownDescription: "IP address of the server to be coonfigured as SDC.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"os_family": schema.StringAttribute{
				Description:         "Operating System family of the SDC.",
				Required:            true,
				MarkdownDescription: "Operating System family of the SDC.",
				Validators: []validator.String{
					stringvalidator.OneOf("linux", "windows", "esxi"),
				},
			},
			"name": schema.StringAttribute{
				Description:         "Name of SDC.",
				Optional:            true,
				Computed:            true,
				MarkdownDescription: "Name of SDC.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"performance_profile": schema.StringAttribute{
				Description:         "Performance profile of the SDC.",
				MarkdownDescription: "Performance profile of the SDC.",
				Optional:            true,
				Computed:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"remote": schema.SingleNestedAttribute{
				Description:         "Remote login details of the SDC.",
				MarkdownDescription: "Remote login details of the SDC.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"user": schema.StringAttribute{
						Description:         "Remote Login username of the SDC server.",
						MarkdownDescription: "Remote Login username of the SDC server.",
						Required:            true,
					},
					"password": schema.StringAttribute{
						Description:         "Remote Login password of the SDC server.",
						MarkdownDescription: "Remote Login password of the SDC server.",
						Optional:            true,
					},
					"private_key": schema.StringAttribute{
						Description:         "Remote Login private key of the SDC server.",
						MarkdownDescription: "Remote Login private key of the SDC server.",
						Optional:            true,
					},
				},
			},
			"package_base64": schema.StringAttribute{
				Description:         "Package to be installed on the SDC in its base64 format.",
				MarkdownDescription: "Package to be installed on the SDC in its base64 format.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"drv_cfg_base64": schema.StringAttribute{
				Description:         "Driver Configuration file for the SDC in its base64 format.",
				MarkdownDescription: "Driver Configuration file for the SDC in its base64 format.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"mdm_ips": schema.ListAttribute{
				Description:         "List of MDM IPs to be assigned to the SDC.",
				MarkdownDescription: "List of MDM IPs to be assigned to the SDC.",
				// TODO: Make Optional + Computed
				Required:    true,
				ElementType: types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
				// TODO
				// PlanModifiers: []planmodifier.List{
				// 	listplanmodifier.UseStateForUnknown(),
				// },
			},
		},
	}
}

func (r *sdcHostResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	r.client = req.ProviderData.(*powerflexProvider).client

	// Get the system on the PowerFlex cluster
	system, err := helper.GetFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}
	r.system = system
}

func (r *sdcHostResource) GetSshProvisioner(ctx context.Context, plan models.SdcHostModel) (*client.SshProvisioner, error) {
	var remote models.SdcHostRemoteModel
	plan.Remote.As(ctx, &remote, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	return client.NewSshProvisioner(client.SshProvisionerConfig{
		IP:         plan.Host.ValueString(),
		Username:   remote.User,
		Password:   remote.Password,
		PrivateKey: remote.PrivateKey,
	}, &provisionerLogger{ctx: ctx})
}

func (r *sdcHostResource) GetMdmIps(ctx context.Context, plan models.SdcHostModel) []string {
	var mdmIps []string
	plan.MdmIPs.ElementsAs(ctx, &mdmIps, true)
	return mdmIps
}

// createEsxi creates an esxi SDC host
func (r *sdcHostResource) createEsxi(ctx context.Context, plan models.SdcHostModel) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics
	sshP, err := r.GetSshProvisioner(ctx, plan)
	if err != nil {
		respDiagnostics.AddError(
			"Error connecting to host",
			err.Error(),
		)
		return respDiagnostics
	}
	defer sshP.Close()

	// upload sw
	scpProv := client.NewScpProvisioner(sshP)
	err = scpProv.Upload(plan.Pkg.ValueString(), "/tmp/package.zip", "")
	if err != nil {
		respDiagnostics.AddError(
			"Error uploading package",
			err.Error(),
		)
		return respDiagnostics
	}

	// install sw
	esxi := client.NewEsxCli(sshP)
	pkgInstallCmd := client.VibInstallCommand{
		ZipFile:  "/tmp/package.zip",
		SigCheck: true,
	}
	op, err := esxi.SoftwareInstall(pkgInstallCmd)
	if err != nil {
		respDiagnostics.AddError(
			"Error installing package",
			err.Error()+"\n"+op,
		)
		return respDiagnostics
	}

	// reboot
	err = sshP.RebootUnix()
	if err != nil {
		respDiagnostics.AddError(
			"Error rebooting",
			err.Error(),
		)
		return respDiagnostics
	}

	// check sw
	tflog.Info(ctx, "Checking for installed sdc package")
	sdc, err := esxi.GetSoftwareByNameRegex(regexp.MustCompile(".*sdc.*"))
	if err != nil {
		respDiagnostics.AddError(
			"Error checking for installed sdc package",
			err.Error(),
		)
		return respDiagnostics
	}
	tflog.Info(ctx, fmt.Sprintf("Installed SDC package is %s", sdc))

	tflog.Info(ctx, "Setting scini module parameters")
	params := map[string]string{
		"IoctlIniGuidStr":        "87254810-e3cc-4c0b-87a7-91d3476b9ec7",
		"IoctlMdmIPStr":          strings.Join(r.GetMdmIps(ctx, plan), ","),
		"bBlkDevIsPdlActive":     "1",
		"blkDevPdlTimeoutMillis": "60000",
	}
	if op, err := esxi.SetModuleParameters("scini", params); err != nil {
		respDiagnostics.AddError(
			"Error setting module parameters",
			err.Error()+"\n"+op,
		)
		return respDiagnostics
	}
	tflog.Info(ctx, "Scini module parameters set")

	// reboot
	err = sshP.RebootUnix()
	if err != nil {
		respDiagnostics.AddError(
			"Error rebooting",
			err.Error(),
		)
		return respDiagnostics
	}

	// load esxi kernel modules
	tflog.Info(ctx, "Loading vmk modules")
	op, err = sshP.Run("vmkload_mod -l")
	if err != nil {
		respDiagnostics.AddError(
			"Error loading vmk modules",
			err.Error(),
		)
		return respDiagnostics
	}
	tflog.Info(ctx, "Finished loading vmk modules")
	tflog.Debug(ctx, op)

	// upload driver config
	// recreate scpProvisioner
	scpProv = client.NewScpProvisioner(sshP)
	tflog.Info(ctx, "Uploading driver config")
	err = scpProv.Upload(plan.DrvCfg.ValueString(), "/tmp/drv_cfg", "0755")
	if err != nil {
		respDiagnostics.AddError(
			"Error uploading package",
			err.Error(),
		)
		return respDiagnostics
	}
	// query mdms via drv cfg
	tflog.Info(ctx, "Querying mdm ips via drv cfg")
	op, err = sshP.Run("/tmp/drv_cfg --query_mdm")
	if err != nil {
		respDiagnostics.AddError(
			"Error querying mdm ips via drv cfg",
			err.Error()+"\n"+op,
		)
		return respDiagnostics
	}

	return respDiagnostics
}

// Create creates the resource and sets the initial Terraform state.
func (r *sdcHostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan
	var plan models.SdcHostModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// install software
	if plan.OS.ValueString() == "esxi" {
		resp.Diagnostics.Append(r.createEsxi(ctx, plan)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// read unconfigured SDC state after installation
	currState, err := r.readSDCHost(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading SDC state",
			err.Error(),
		)
		return
	}
	diags = resp.State.Set(ctx, currState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// configure SDC via API
	err = r.setSDCParams(ctx, plan, currState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting SDC parameters",
			err.Error(),
		)
		return
	}

	// read final state of SDC and set state
	state, err := r.readSDCHost(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading SDC state",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *sdcHostResource) setSDCParams(ctx context.Context, plan, state models.SdcHostModel) error {
	// set name
	if plan.Name.ValueString() != state.Name.ValueString() && !plan.Name.IsUnknown() {
		tflog.Info(ctx, "Setting SDC name")
		if _, err := r.system.ChangeSdcName(state.ID.ValueString(), plan.Name.ValueString()); err != nil {
			return fmt.Errorf("error setting SDC name: %w", err)
		}
		tflog.Info(ctx, "SDC name set")
	}

	// set Performance Profile
	if plan.PerformanceProfile.ValueString() != state.PerformanceProfile.ValueString() && !plan.PerformanceProfile.IsUnknown() {
		tflog.Info(ctx, "Setting SDC performance profile")
		if _, err := r.system.ChangeSdcPerfProfile(state.ID.ValueString(), plan.PerformanceProfile.ValueString()); err != nil {
			return fmt.Errorf("error setting SDC performance profile: %w", err)
		}
		tflog.Info(ctx, "SDC performance profile set")
	}

	return nil
}

func (r *sdcHostResource) readSDCHost(ctx context.Context, state models.SdcHostModel) (models.SdcHostModel, error) {
	// get SDC by IP
	tflog.Info(ctx, "Finding SDC by IP")
	sdcData, err := r.system.FindSdc("SdcIP", state.Host.ValueString())
	if err != nil {
		return state, fmt.Errorf("error finding SDC by IP %s: %w", state.Host.ValueString(), err)
	}
	tflog.Info(ctx, "Found SDC by IP")
	state.ID = types.StringValue(sdcData.Sdc.ID)
	state.PerformanceProfile = types.StringValue(sdcData.Sdc.PerfProfile)
	state.Name = types.StringValue(sdcData.Sdc.Name)
	return state, nil
}

// Read refreshes the Terraform state with the latest data.
func (r *sdcHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state models.SdcHostModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	newState, err := r.readSDCHost(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error refreshing SDC state",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, newState)
	resp.Diagnostics.Append(diags...)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sdcHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var (
		plan      models.SdcHostModel
		currState models.SdcHostModel
	)

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Retrieve values from state
	diags = req.State.Get(ctx, &currState)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: check that any the stuff that cannot be updated are not changed
	// unupdateable fields: os_family, mdm_ips, package, drv_cfg
	if !currState.OS.IsNull() && !plan.OS.Equal(currState.OS) {
		resp.Diagnostics.AddError("Error updating SDC", "OS cannot be changed")
	}
	if !currState.MdmIPs.IsNull() && !plan.MdmIPs.Equal(currState.MdmIPs) {
		resp.Diagnostics.AddError("Error updating SDC", "mdm_ips cannot be changed")
	}
	if !currState.Pkg.IsNull() && !plan.Pkg.Equal(currState.Pkg) {
		resp.Diagnostics.AddError("Error updating SDC", "package cannot be changed")
	}
	if !currState.DrvCfg.IsNull() && !plan.DrvCfg.Equal(currState.DrvCfg) {
		resp.Diagnostics.AddError("Error updating SDC", "drv_cfg cannot be changed")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// configure SDC via API
	err := r.setSDCParams(ctx, plan, currState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting SDC parameters",
			err.Error(),
		)
		return
	}

	// read final state of SDC and set state
	state, err := r.readSDCHost(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading SDC state",
			err.Error(),
		)
		return
	}

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (r *sdcHostResource) deleteEsxi(ctx context.Context, state models.SdcHostModel) diag.Diagnostics {
	var respDiagnostics diag.Diagnostics
	// Disconnect from PowerFlex
	tflog.Info(ctx, "Logging into host...")
	sshP, err := r.GetSshProvisioner(ctx, state)
	if err != nil {
		respDiagnostics.AddError(
			"Error connecting to host",
			err.Error(),
		)
		return respDiagnostics
	}
	defer sshP.Close()

	tflog.Info(ctx, "Checking for installed sdc package")
	esxi := client.NewEsxCli(sshP)
	sdc, err := esxi.GetSoftwareByNameRegex(regexp.MustCompile(".*sdc.*"))
	if err != nil {
		respDiagnostics.AddError(
			"Error checking for installed sdc package",
			err.Error(),
		)
		return respDiagnostics
	}
	tflog.Info(ctx, fmt.Sprintf("Installed SDC package is %v+", sdc))

	op, err := esxi.SoftwareRmv(sdc.Name)
	if err != nil {
		respDiagnostics.AddError(
			"Error removing sdc package",
			err.Error()+"\n"+op,
		)
		return respDiagnostics
	}
	tflog.Info(ctx, fmt.Sprintf("sdc package removed: %s", op))

	err = sshP.RebootUnix()
	if err != nil {
		respDiagnostics.AddError(
			"Error rebooting host",
			err.Error(),
		)
		return respDiagnostics
	}

	return respDiagnostics
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *sdcHostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.SdcHostModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// remove software
	if state.OS.ValueString() == "esxi" {
		resp.Diagnostics.Append(r.deleteEsxi(ctx, state)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// if name is configured, remove sdc via API
	if state.Name.ValueString() != "" {
		err := r.system.DeleteSdc(state.ID.ValueString())
		if err != nil {
			resp.Diagnostics.AddError(
				"Error deleting SDC",
				err.Error(),
			)
			return
		}
	}

	// remove state
	resp.State.RemoveResource(ctx)
}

// ImportState - function to ImportState for SDC resource.
func (r *sdcHostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] ImportState :-- "+helper.PrettyJSON(req))
	resource.ImportStatePassthroughID(ctx, path.Root("ip"), req, resp)
}
