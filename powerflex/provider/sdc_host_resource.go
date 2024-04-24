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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &sdcHostResource{}
	_ resource.ResourceWithConfigure = &sdcHostResource{}
	// _ resource.ResourceWithImportState = &sdcHostResource{}
)

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

func (r *sdcHostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource is used to manage the Storage Data Servers entity of PowerFlex Array. We can Create, Update and Delete the SDC using this resource. We can also import an existing SDC from PowerFlex array.",
		MarkdownDescription: "This resource is used to manage the Storage Data Servers entity of PowerFlex Array. We can Create, Update and Delete the SDC using this resource. We can also import an existing SDC from PowerFlex array.",
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Description:         "The id of the SDC",
				Computed:            true,
				MarkdownDescription: "The id of the SDC",
			},
			"hostname": schema.StringAttribute{
				Description:         "IP address of the server to be coonfigured as SDC.",
				Required:            true,
				MarkdownDescription: "IP address of the server to be coonfigured as SDC.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			// "protection_domain_id": schema.StringAttribute{
			// 	Description: "ID of the Protection Domain under which the SDC will be created." +
			// 		" Conflicts with 'protection_domain_name'." +
			// 		" Cannot be updated.",
			// 	Optional: true,
			// 	Computed: true,
			// 	MarkdownDescription: "ID of the Protection Domain under which the SDC will be created." +
			// 		" Conflicts with `protection_domain_name`." +
			// 		" Cannot be updated.",
			// },
			"name": schema.StringAttribute{
				Description:         "Name of SDC.",
				Optional:            true,
				MarkdownDescription: "Name of SDC.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
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
						// Validators: []validator.String{
						// 	stringvalidator.AlsoRequires(path.MatchRelative().AtName("private_key")),
						// },
					},
					"private_key": schema.StringAttribute{
						Description:         "Remote Login private key of the SDC server.",
						MarkdownDescription: "Remote Login private key of the SDC server.",
						Optional:            true,
					},
				},
				// Validators: []validator.Object{
				// 	objectvalidator.ExactlyOneOf(
				// 		path.MatchRelative().AtName("password"),
				// 		path.MatchRelative().AtName("private_key"),
				// 	),
				// },
			},
			"package_base64": schema.StringAttribute{
				Description:         "Package to be installed on the SDC in its base64 format.",
				MarkdownDescription: "Package to be installed on the SDC in its base64 format.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"mdm_ips": schema.ListAttribute{
				Description:         "List of MDM IPs to be assigned to the SDC.",
				MarkdownDescription: "List of MDM IPs to be assigned to the SDC.",
				Optional:            true,
				Computed:            true,
				ElementType:         types.StringType,
				Validators: []validator.List{
					listvalidator.SizeAtLeast(1),
					listvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				},
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
	}, nil)
}

func (r *sdcHostResource) GetMdmIps(ctx context.Context, plan models.SdcHostModel) []string {
	var mdmIps []string
	plan.MdmIPs.ElementsAs(ctx, &mdmIps, true)
	return mdmIps
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

	sshP, err := r.GetSshProvisioner(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error connecting to host",
			err.Error(),
		)
		return
	}
	//////////////
	defer sshP.Close()

	// upload sw
	scpProv := client.NewScpProvisioner(sshP)
	err = scpProv.Upload(plan.Pkg.ValueString(), "/tmp/package.zip", "")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error uploading package",
			err.Error(),
		)
	}

	// install sw
	esxi := client.NewEsxCli(sshP)
	pkgInstallCmd := client.VibInstallCommand{
		ZipFile:  "/tmp/package.zip",
		SigCheck: true,
	}
	op, err := esxi.SoftwareInstall(pkgInstallCmd)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error installing package",
			err.Error()+"\n"+op,
		)
		return
	}

	// reboot
	err = sshP.RebootUnix()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rebooting",
			err.Error(),
		)
		return
	}

	// check sw
	tflog.Info(ctx, "Checking for installed sdc package")
	sdc, err := esxi.GetSoftwareByNameRegex(regexp.MustCompile(".*sdc.*"))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error checking for installed sdc package",
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Installed SDC package is %s", sdc))

	tflog.Info(ctx, "Setting scini module parameters")
	params := map[string]string{
		"IoctlIniGuidStr": "87254810-e3cc-4c0b-87a7-91d3476b9ec7",
		//"10.247.100.214,10.247.66.67"
		"IoctlMdmIPStr":          strings.Join(r.GetMdmIps(ctx, plan), ","),
		"bBlkDevIsPdlActive":     "1",
		"blkDevPdlTimeoutMillis": "60000",
	}
	_, err = esxi.SetModuleParameters("scini", params)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting module parameters",
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "Scini module parameters set")

	// reboot
	err = sshP.RebootUnix()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rebooting",
			err.Error(),
		)
		return
	}

	// load vmk modules
	tflog.Info(ctx, "Loading vmk modules")
	op, err = sshP.Run("vmkload_mod -l")
	if err != nil {
		resp.Diagnostics.AddError(
			"Error loading vmk modules",
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, "Finished loading vmk modules")
	tflog.Debug(ctx, op)

	plan.ID = types.StringValue("dummy")
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)
}

// Read refreshes the Terraform state with the latest data.
func (r *sdcHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	// var (
	// 	state models.SdcHostModel
	// 	rsp   scaleiotypes.Sds
	// 	err   error
	// )
	// diags := req.State.Get(ctx, &state)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// // Get SDC
	// if rsp, err = r.system.GetSdsByID(state.ID.ValueString()); err != nil {
	// 	resp.Diagnostics.AddError(
	// 		fmt.Sprintf("Could not get SDC by ID %s", state.ID.ValueString()),
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// // Set refreshed state
	// state, dgs := helper.UpdateSdsState(&rsp, state)
	// resp.Diagnostics.Append(dgs...)

	// diags = resp.State.Set(ctx, state)
	// resp.Diagnostics.Append(diags...)
	resp.State = req.State
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *sdcHostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan
	var (
		plan  models.SdcHostModel
		state models.SdcHostModel
		// err   error
	)

	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// check that plan is valid
	// resp.Diagnostics.Append(r.ValidatePlan(ctx, plan)...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// Retrieve values from state
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// TODO: check that any the stuff that cannot be updated are not changed
	// if !plan.ProtectionDomainID.IsUnknown() && plan.ProtectionDomainID.ValueString() != state.ProtectionDomainID.ValueString() {
	// 	resp.Diagnostics.AddError(
	// 		"Protection domain ID cannot be updated",
	// 		"Protection domain ID cannot be updated")
	// 	return
	// }

	if resp.Diagnostics.HasError() {
		return
	}

	// Set refreshed state
	// state, dgs := helper.UpdateSdsState(rsp, state)
	// resp.Diagnostics.Append(dgs...)

	// diags = resp.State.Set(ctx, state)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
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

	// Disconnect from PowerFlex
	tflog.Info(ctx, "Logging into host...")
	sshP, err := r.GetSshProvisioner(ctx, state)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error connecting to host",
			err.Error(),
		)
		return
	}
	defer sshP.Close()

	tflog.Info(ctx, "Checking for installed sdc package")
	esxi := client.NewEsxCli(sshP)
	sdc, err := esxi.GetSoftwareByNameRegex(regexp.MustCompile(".*sdc.*"))
	if err != nil {
		resp.Diagnostics.AddError(
			"Error checking for installed sdc package",
			err.Error(),
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("Installed SDC package is %v+", sdc))

	op, err := esxi.SoftwareRmv(sdc.Name)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error removing sdc package",
			err.Error()+"\n"+op,
		)
		return
	}
	tflog.Info(ctx, fmt.Sprintf("sdc package removed: %s", op))

	err = sshP.RebootUnix()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error rebooting host",
			err.Error(),
		)
		return
	}

	// remove state
	resp.State.RemoveResource(ctx)
}
