/*
Copyright (c) 2024 Dell Inc., or its subsidiaries. All Rights Reserved.

Licensed under the Mozilla Public License Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://mozilla.org/MPL/2.0/

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package provider

import (
	"context"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework-validators/resourcevalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/booldefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/listplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &sdcHostResource{}
	_ resource.ResourceWithConfigure   = &sdcHostResource{}
	_ resource.ResourceWithImportState = &sdcHostResource{}
)

// NewSDCHostResource is a helper function to simplify the provider implementation.
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
		resourcevalidator.Conflicting(
			remotePath.AtName("password"),
			remotePath.AtName("certificate"),
		),
	}
}

func (r *sdcHostResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	// Retrieve values from plan
	var cfg models.SdcHostModel

	diags := req.Config.Get(ctx, &cfg)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// for esxi, make sure that esxi block is provided (not null)
	// Note: cfg.OS cannot be null (required field)
	if !cfg.OS.IsUnknown() && cfg.OS.ValueString() == "esxi" && cfg.Esxi.IsNull() {
		resp.Diagnostics.AddAttributeError(
			path.Root("esxi"),
			"Esxi block is required for esxi SDC",
			"",
		)
	}

	if !cfg.OS.IsUnknown() && cfg.OS.ValueString() == "windows" && !cfg.Remote.IsNull() {
		var remote models.SdcHostRemoteModel
		cfg.Remote.As(ctx, &remote, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})

		if remote.Password == nil {
			resp.Diagnostics.AddAttributeError(
				path.Root("remote").AtName("password"),
				"Password is required for Windows SDC",
				"",
			)
		}
	}
}

func (r *sdcHostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description:         "This resource is used to manage the SDC entity of the PowerFlex Array. We can Create, Update and Delete the SDC using this resource. We can also import an existing SDC from the PowerFlex array.",
		MarkdownDescription: "This resource is used to manage the SDC entity of the PowerFlex Array. We can Create, Update and Delete the SDC using this resource. We can also import an existing SDC from the PowerFlex array.",
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
				Description:         "IP address of the server to be configured as SDC.",
				Required:            true,
				MarkdownDescription: "IP address of the server to be configured as SDC.",
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"os_family": schema.StringAttribute{
				Description: "Operating System family of the SDC." +
					" Accepted values are 'linux', 'windows' and 'esxi'." +
					" Cannot be changed once set.",
				Required: true,
				MarkdownDescription: "Operating System family of the SDC." +
					" Accepted values are 'linux', 'windows' and 'esxi'." +
					" Cannot be changed once set.",
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
				Description: "Performance profile of the SDC." +
					" Accepted values are 'HighPerformance' and 'Compact'.",
				MarkdownDescription: "Performance profile of the SDC." +
					" Accepted values are 'HighPerformance' and 'Compact'.",
				Optional: true,
				Computed: true,
				Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
					"HighPerformance",
					"Compact",
				)},
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.UseStateForUnknown(),
				},
			},
			"use_remote_path": schema.BoolAttribute{
				Description:         "Use path on remote server where SDC is installed. Defaults to `false`.",
				MarkdownDescription: "Use path on remote server where SDC is installed. Defaults to `false`.",
				Optional:            true,
				Computed:            true,
				Default:             booldefault.StaticBool(false),
			},
			"remote": schema.SingleNestedAttribute{
				Description:         "Remote login details of the SDC.",
				MarkdownDescription: "Remote login details of the SDC.",
				Required:            true,
				Attributes: map[string]schema.Attribute{
					"port": schema.StringAttribute{
						Description:         "Remote Login port of the SDC server. Defaults to `22`.",
						MarkdownDescription: "Remote Login port of the SDC server. Defaults to `22`.",
						Optional:            true,
						Computed:            true,
						Default:             stringdefault.StaticString("22"),
					},
					"user": schema.StringAttribute{
						Description:         "Remote Login username of the SDC server.",
						MarkdownDescription: "Remote Login username of the SDC server.",
						Required:            true,
					},
					"password": schema.StringAttribute{
						Description:         "Remote Login password of the SDC server.",
						MarkdownDescription: "Remote Login password of the SDC server.",
						Optional:            true,
						Sensitive:           true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"private_key": schema.StringAttribute{
						Description: "Remote Login private key of the SDC server." +
							" Corresponds to the IdentityFile field of OpenSSH.",
						MarkdownDescription: "Remote Login private key of the SDC server." +
							" Corresponds to the IdentityFile field of OpenSSH.",
						Optional: true,
					},
					"certificate": schema.StringAttribute{
						Description: "Remote Login certificate issued by a CA to the remote login user." +
							" Must be used with `private_key` and the private key must match the certificate.",
						MarkdownDescription: "Remote Login certificate issued by a CA to the remote login user." +
							" Must be used with `private_key` and the private key must match the certificate.",
						Optional: true,
					},
					"host_key": schema.StringAttribute{
						Description: "Remote Login host key of the SDC server." +
							" Corresponds to the UserKnownHostsFile field of OpenSSH.",
						MarkdownDescription: "Remote Login host key of the SDC server." +
							" Corresponds to the UserKnownHostsFile field of OpenSSH.",
						Optional: true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"dir": schema.StringAttribute{
						Description: "Directory on the SDC server to upload packages to for Unix." +
							" Defaults to `/tmp` on Unix.",
						MarkdownDescription: "Directory on the SDC server to upload packages to for Unix." +
							" Defaults to `/tmp` on Unix.",
						Optional: true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
				},
			},
			"esxi": schema.SingleNestedAttribute{
				Description:         "Details of the SDC host if the `os_family` is `esxi`.",
				MarkdownDescription: "Details of the SDC host if the `os_family` is `esxi`.",
				Optional:            true,
				Attributes: map[string]schema.Attribute{
					"guid": schema.StringAttribute{
						Description:         "GUID of the SDC.",
						MarkdownDescription: "GUID of the SDC.",
						Required:            true,
					},
					"verify_vib_signature": schema.BoolAttribute{
						Description:         "Whether to verify the VIB signature or not. Defaults to `true`.",
						MarkdownDescription: "Whether to verify the VIB signature or not. Defaults to `true`.",
						Optional:            true,
						Computed:            true,
						Default:             booldefault.StaticBool(true),
					},
				},
			},
			"linux_drv_cfg": schema.StringAttribute{
				Description:         "Path to the drv_cfg for linux, defaults to /opt/emc/scaleio/sdc/bin/",
				MarkdownDescription: "Path to the drv_cfg for linux, defaults to /opt/emc/scaleio/sdc/bin/",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("/opt/emc/scaleio/sdc/bin/"),
			},
			"windows_drv_cfg": schema.StringAttribute{
				Description:         "Path to the drv_cfg.exe for windows, defaults to C:\\Program Files\\EMC\\scaleio\\sdc\\bin\\",
				MarkdownDescription: "Path to the drv_cfg.exe config for windows, defaults to C:\\Program Files\\EMC\\scaleio\\sdc\\bin\\",
				Optional:            true,
				Computed:            true,
				Default:             stringdefault.StaticString("C:\\Program Files\\EMC\\scaleio\\sdc\\bin\\"),
			},
			"package_path": schema.StringAttribute{
				Description:         "Full path (on local machine) of the package to be installed on the SDC.",
				MarkdownDescription: "Full path (on local machine) of the package to be installed on the SDC.",
				Required:            true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"clusters_mdm_ips": schema.ListAttribute{
				Description: "List of MDM IPs (primary,secondary or list of virtual IPs) seperated by cluster, to be assigned to the SDC." +
					"Each string in the list is a set of Mdm Ips related to a specific cluster. These Ips should be seperated by comma I.E. ['x.x.x.x,y.y.y.y', 'z.z.z.z,a.a.a.a']. ",
				MarkdownDescription: "List of MDM IPs (primary,secondary or list of virtual IPs) seperated by cluster, to be assigned to the SDC." +
					"Each string in the list is a set of Mdm Ips related to a specific cluster. These Ips should be seperated by comma I.E. ['x.x.x.x,y.y.y.y', 'z.z.z.z,a.a.a.a']. ",
				Optional:    true,
				Computed:    true,
				ElementType: types.StringType,
				PlanModifiers: []planmodifier.List{
					listplanmodifier.UseStateForUnknown(),
				},
			},
			"guid": schema.StringAttribute{
				Description:         "GUID of the HOST",
				MarkdownDescription: "GUID of the HOST",
				Computed:            true,
			},
			"on_vmware": schema.BoolAttribute{
				Description:         "Is Host on VMware",
				MarkdownDescription: "Is Host on VMware",
				Computed:            true,
			},
			"is_approved": schema.BoolAttribute{
				Description:         "Is Host Approved",
				MarkdownDescription: "Is Host Approved",
				Computed:            true,
			},
			"system_id": schema.StringAttribute{
				Description:         "System ID of the Host",
				MarkdownDescription: "System ID of the Host",
				Computed:            true,
			},
			"mdm_connection_state": schema.StringAttribute{
				Description:         "MDM Connection State",
				MarkdownDescription: "MDM Connection State",
				Computed:            true,
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

// ModifyPlan handles plan modification
func (r *sdcHostResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	// If the entire plan is null, the resource is planned for destruction. Dont do anything.
	if req.Plan.Raw.IsNull() {
		return
	}

	// Retrieve values from plan
	var plan models.SdcHostModel
	resp.Diagnostics.Append(req.Plan.Get(ctx, &plan)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if resource is getting created
	if req.State.Raw.IsNull() {
		return
	}
	// Only update scenarios from here on

	// Retrieve values from state
	var state models.SdcHostModel
	resp.Diagnostics.Append(req.State.Get(ctx, &state)...)
	if resp.Diagnostics.HasError() {
		return
	}

	// if esxi block is getting updated, throw error
	if helper.Known(state.Esxi, plan.Esxi) && !state.Esxi.Equal(plan.Esxi) {
		resp.Diagnostics.AddAttributeError(
			path.Root("esxi"),
			"ESXi SDC details cannot be updated",
			"",
		)
	}

	// if IP is getting updated, throw error
	if helper.Known(state.Host, plan.Host) && !state.Host.Equal(plan.Host) {
		resp.Diagnostics.AddAttributeError(
			path.Root("ip"),
			"SDC IP cannot be updated through this resource",
			"Please update the SDC IP by other means and refresh state by running terraform apply -refresh-only",
		)
	}
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

	//Checking that SDC exist or not
	allSdcs, _ := r.system.GetSdc()
	for _, sdcData := range allSdcs {
		// check if IP exists
		if sdcData.SdcIP == plan.Host.ValueString() {
			resp.Diagnostics.AddAttributeError(
				path.Root("ip"),
				"SDC Host with given IP already exists",
				"",
			)
		}
		// check if name exists
		if helper.Known(plan.Name) && sdcData.Name == plan.Name.ValueString() {
			resp.Diagnostics.AddAttributeError(
				path.Root("name"),
				"SDC Host with given name already exists",
				"",
			)
		}
	}
	if resp.Diagnostics.HasError() {
		return
	}

	resHelper := helper.SdcHostResource{
		System: r.system,
	}

	// install software
	if plan.OS.ValueString() == "esxi" {
		resp.Diagnostics.Append(resHelper.CreateEsxi(ctx, plan)...)
	} else if plan.OS.ValueString() == "windows" {
		resp.Diagnostics.Append(resHelper.CreateWindows(ctx, plan)...)
	} else if plan.OS.ValueString() == "linux" {
		resp.Diagnostics.Append(resHelper.LinuxOp(ctx, plan, true)...)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	// read unconfigured SDC state after installation
	currState, err := resHelper.ReadSDCHost(ctx, plan)
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
	err = resHelper.SetSDCParams(ctx, plan, currState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting SDC parameters",
			err.Error()+
				"\n\nThis is a minor issue which does not require resource recreation (ie. SDC re-installation) to resolve. Terraform, by default, will mark this resource as tainted "+
				"and recreate it on the next apply. Please untaint this resource manually prior to applying again to prevent unnecessary SDC re-installations.",
		)
		return
	}

	// read final state of SDC and set state
	state, err := resHelper.ReadSDCHost(ctx, plan)
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

// Read refreshes the Terraform state with the latest data.
func (r *sdcHostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from state
	var state models.SdcHostModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	resHelper := helper.SdcHostResource{
		System: r.system,
	}

	newState, err := resHelper.ReadSDCHost(ctx, state)
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

	// Check that anything that cannot be updated are not changed
	// unupdateable fields: os_family, package
	if !currState.OS.IsNull() && !plan.OS.Equal(currState.OS) {
		resp.Diagnostics.AddError("Error updating SDC", "OS cannot be changed")
	}
	if !currState.Pkg.IsNull() && !plan.Pkg.Equal(currState.Pkg) {
		resp.Diagnostics.AddError("Error updating SDC", "package cannot be changed")
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resHelper := helper.SdcHostResource{
		System: r.system,
	}

	// Only run this update if the mdms need to be updated
	if !plan.MdmIPs.Equal(currState.MdmIPs) {

		// if the mdms need to be updated do it for the specific OS
		if plan.OS.ValueString() == "esxi" {
			resp.Diagnostics.Append(resHelper.UpdateEsxiMdms(ctx, plan)...)
		} else if plan.OS.ValueString() == "windows" {
			resp.Diagnostics.Append(resHelper.UpdateWindowsMdms(ctx, plan)...)
		} else if plan.OS.ValueString() == "linux" {
			resp.Diagnostics.Append(resHelper.UpdateLinuxMdms(ctx, plan)...)
		}
		if resp.Diagnostics.HasError() {
			return
		}
	}
	// configure SDC via API
	err := resHelper.SetSDCParams(ctx, plan, currState)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error setting SDC parameters",
			err.Error(),
		)
		return
	}

	// read final state of SDC and set state
	state, err := resHelper.ReadSDCHost(ctx, plan)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error reading SDC state",
			err.Error(),
		)
		return
	}
	// This is needed when doing an import then update
	state.Host = currState.Host

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
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

	resHelper := helper.SdcHostResource{
		System: r.system,
	}

	// remove software
	if state.OS.ValueString() == "esxi" {
		resp.Diagnostics.Append(resHelper.DeleteEsxi(ctx, state)...)
	} else if state.OS.ValueString() == "windows" {
		resp.Diagnostics.Append(resHelper.DeleteWindows(ctx, state)...)
	} else if state.OS.ValueString() == "linux" {
		resp.Diagnostics.Append(resHelper.LinuxOp(ctx, state, false)...)
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
