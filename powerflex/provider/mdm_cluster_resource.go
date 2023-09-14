/*
Copyright (c) 2023 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"fmt"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"
	"time"

	"github.com/dell/goscaleio"
	scaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/listvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/int64planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/setplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource              = &mdmClusterResource{}
	_ resource.ResourceWithConfigure = &mdmClusterResource{}
)

// NewMdmClusterResource returns the resource for MDM
func NewMdmClusterResource() resource.Resource {
	return &mdmClusterResource{}
}

type mdmClusterResource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (d *mdmClusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_mdm_cluster"
}

func (d *mdmClusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = MdmClusterResourceSchema
}

// MdmClusterResourceSchema defines the schema for Mdm resource
var MdmClusterResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource can be used to manage MDM cluster on a PowerFlex array. Supports adding or removing standby MDMs, migrate from 3-node to 5-node cluster or vice-versa, changing MDM ownership, changing performance profile, and renaming MDMs.",
	MarkdownDescription: "This resource can be used to manage MDM cluster on a PowerFlex array. Supports adding or removing standby MDMs, migrate from 3-node to 5-node cluster or vice-versa, changing MDM ownership, changing performance profile, and renaming MDMs.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Unique identifier of the MDM cluster.",
			MarkdownDescription: "Unique identifier of the MDM cluster.",
			Computed:            true,
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"performance_profile": schema.StringAttribute{
			Description:         "Performance profile of the MDM cluster. Accepted values are 'Compact' and 'HighPerformance'.",
			MarkdownDescription: "Performance profile of the MDM cluster. Accepted values are `Compact` and `HighPerformance`.",
			Optional:            true,
			Computed:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(scaleio_types.PerformanceProfileCompact),
					string(scaleio_types.PerformanceProfileHigh),
				),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"cluster_mode": schema.StringAttribute{
			Description:         "Mode of the MDM cluster. Accepted values are 'ThreeNodes' and 'FiveNodes'.",
			MarkdownDescription: "Mode of the MDM cluster. Accepted values are `ThreeNodes` and `FiveNodes`.",
			Required:            true,
			Validators: []validator.String{
				stringvalidator.OneOf(
					string(scaleio_types.ThreeNodesClusterMode),
					string(scaleio_types.FiveNodesClusterMode),
				),
			},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"primary_mdm": schema.SingleNestedAttribute{
			Required:            true,
			Description:         "Primary MDM details.",
			MarkdownDescription: "Primary MDM details.",
			Attributes: map[string]schema.Attribute{
				"id": schema.StringAttribute{
					Description:         "ID of the primary MDM.",
					MarkdownDescription: "ID of the primary MDM.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
						stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("ips")),
					},
					PlanModifiers: []planmodifier.String{
						stringplanmodifier.UseStateForUnknown(),
					},
				},
				"name": schema.StringAttribute{
					Description:         "Name of the the primary MDM.",
					MarkdownDescription: "Name of the the primary MDM.",
					Optional:            true,
					Computed:            true,
					Validators: []validator.String{
						stringvalidator.LengthAtLeast(1),
					},
				},
				"port": schema.Int64Attribute{
					Description:         "Port of the primary MDM.",
					MarkdownDescription: "Port of the primary MDM.",
					Computed:            true,
					PlanModifiers: []planmodifier.Int64{
						int64planmodifier.UseStateForUnknown(),
					},
				},
				"ips": schema.SetAttribute{
					Description:         "The Ips of the primary MDM.",
					MarkdownDescription: "The Ips of the primary MDM.",
					ElementType:         types.StringType,
					Optional:            true,
					Computed:            true,
					Validators: []validator.Set{
						setvalidator.SizeBetween(1, 4),
						setvalidator.ValueStringsAre(stringvalidator.LengthBetween(7, 15)),
					},
				},
				"management_ips": schema.SetAttribute{
					Description:         "The management ips of the primary MDM.",
					MarkdownDescription: "The management ips of the primary MDM.",
					ElementType:         types.StringType,
					Computed:            true,
					PlanModifiers: []planmodifier.Set{
						setplanmodifier.UseStateForUnknown(),
					},
				},
			},
		},
		"secondary_mdm": schema.ListNestedAttribute{
			Required:            true,
			Description:         "Secondary MDM details.",
			MarkdownDescription: "Secondary MDM details.",
			Validators: []validator.List{
				listvalidator.SizeBetween(1, 2),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         "ID of the secondary MDM.",
						MarkdownDescription: "ID of the secondary MDM.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("ips")),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Description:         "Name of the the secondary MDM.",
						MarkdownDescription: "Name of the the secondary MDM.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
					},
					"port": schema.Int64Attribute{
						Description:         "Port of the secondary MDM.",
						MarkdownDescription: "Port of the secondary MDM.",
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"ips": schema.SetAttribute{
						Description:         "The Ips of the secondary MDM.",
						MarkdownDescription: "The Ips of the secondary MDM.",
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						Validators: []validator.Set{
							setvalidator.SizeBetween(1, 4),
							setvalidator.ValueStringsAre(stringvalidator.LengthBetween(7, 15)),
						},
					},
					"management_ips": schema.SetAttribute{
						Description:         "The management ips of the secondary MDM.",
						MarkdownDescription: "The management ips of the secondary MDM.",
						ElementType:         types.StringType,
						Computed:            true,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
		"tiebreaker_mdm": schema.ListNestedAttribute{
			Required:            true,
			Description:         "TieBreaker MDM details.",
			MarkdownDescription: "TieBreaker MDM details.",
			Validators: []validator.List{
				listvalidator.SizeBetween(1, 2),
			},
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         "ID of the tiebreaker MDM.",
						MarkdownDescription: "ID of the tiebreaker MDM.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
							stringvalidator.ExactlyOneOf(path.MatchRelative().AtParent().AtName("ips")),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Description:         "Name of the the tiebreaker MDM.",
						MarkdownDescription: "Name of the the tiebreaker MDM.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"port": schema.Int64Attribute{
						Description:         "Port of the tiebreaker MDM.",
						MarkdownDescription: "Port of the tiebreaker MDM.",
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"ips": schema.SetAttribute{
						Description:         "The Ips of the tiebreaker MDM.",
						MarkdownDescription: "The Ips of the tiebreaker MDM.",
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
						Validators: []validator.Set{
							setvalidator.SizeBetween(1, 4),
							setvalidator.ValueStringsAre(stringvalidator.LengthBetween(7, 15)),
						},
					},
					"management_ips": schema.SetAttribute{
						Description:         "The management ips of the tiebreaker MDM.",
						MarkdownDescription: "The management ips of the tiebreaker MDM.",
						ElementType:         types.StringType,
						Computed:            true,
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
					},
				},
			},
		},
		"standby_mdm": schema.ListNestedAttribute{
			Optional:            true,
			Computed:            true,
			Description:         "StandBy MDM details. StandBy MDM can be added/removed/promoted to manager/tiebreaker role.",
			MarkdownDescription: "StandBy MDM details. StandBy MDM can be added/removed/promoted to manager/tiebreaker role.",
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"id": schema.StringAttribute{
						Description:         "ID of the standby MDM.",
						MarkdownDescription: "ID of the standby MDM.",
						Computed:            true,
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"name": schema.StringAttribute{
						Description:         "Name of the the standby MDM.",
						MarkdownDescription: "Name of the the standby MDM.",
						Optional:            true,
						Computed:            true,
						Validators: []validator.String{
							stringvalidator.LengthAtLeast(1),
						},
						PlanModifiers: []planmodifier.String{
							stringplanmodifier.UseStateForUnknown(),
						},
					},
					"port": schema.Int64Attribute{
						Description:         "Port of the standby MDM. Cannot be updated.",
						MarkdownDescription: "Port of the standby MDM. Cannot be updated.",
						Optional:            true,
						Computed:            true,
						PlanModifiers: []planmodifier.Int64{
							int64planmodifier.UseStateForUnknown(),
						},
					},
					"ips": schema.SetAttribute{
						Description:         "The Ips of the standby MDM. Cannot be updated.",
						MarkdownDescription: "The Ips of the standby MDM. Cannot be updated.",
						ElementType:         types.StringType,
						Required:            true,
						Validators: []validator.Set{
							setvalidator.SizeBetween(1, 4),
							setvalidator.ValueStringsAre(stringvalidator.LengthBetween(7, 15)),
						},
					},
					"management_ips": schema.SetAttribute{
						Description:         "The management ips of the standby MDM. Cannot be updated.",
						MarkdownDescription: "The management ips of the standby MDM. Cannot be updated.",
						ElementType:         types.StringType,
						Optional:            true,
						Computed:            true,
						Validators: []validator.Set{
							setvalidator.SizeBetween(1, 4),
							setvalidator.ValueStringsAre(stringvalidator.LengthBetween(7, 15)),
						},
						PlanModifiers: []planmodifier.Set{
							setplanmodifier.UseStateForUnknown(),
						},
					},
					"role": schema.StringAttribute{
						Required:            true,
						Description:         "Role of the standby mdm. Accepted values are 'Manager' and 'TieBreaker'. Cannot be updated.",
						MarkdownDescription: "Role of the standby mdm. Accepted values are `Manager` and `TieBreaker`. Cannot be updated.",
						Validators: []validator.String{
							stringvalidator.OneOf(
								scaleio_types.Manager,
								scaleio_types.TieBreaker,
							),
						},
					},
					"allow_asymmetric_ips": schema.BoolAttribute{
						Optional:            true,
						Computed:            true,
						Description:         "Allow the added MDM to have a different number of IPs from the primary MDM. Cannot be updated.",
						MarkdownDescription: "Allow the added MDM to have a different number of IPs from the primary MDM. Cannot be updated.",
						PlanModifiers: []planmodifier.Bool{
							helper.BoolDefault(false),
						},
					},
				},
			},
		},
	},
}

func (d *mdmClusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client == nil {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)
		return
	}

	d.client = req.ProviderData.(*powerflexProvider).client

	system, err := helper.GetFirstSystem(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex System",
			err.Error(),
		)
		return
	}
	d.system = system
}

func (d *mdmClusterResource) ValidateConfig(ctx context.Context, ipmap map[string]string, plan *models.MdmResourceModel) diag.Diagnostics {
	var diags diag.Diagnostics

	// Check if valid IP is provided for Primary MDM
	if !plan.PrimaryMdm.Ips.IsNull() {
		configIps := make([]string, 0)
		diags.Append(plan.PrimaryMdm.Ips.ElementsAs(ctx, &configIps, true)...)

		if len(configIps) > 0 {
			if _, ok := ipmap[configIps[0]]; !ok {
				diags.AddAttributeError(
					path.Root("primary_mdm"),
					"Please enter valid IP for primary MDM",
					"Please enter valid IP for primary MDM",
				)
			}
		}
	}

	// Check if valid IP is provided for Secondary MDM
	for _, mdm := range plan.SecondaryMdm {
		if !mdm.Ips.IsNull() {
			configIps := make([]string, 0)
			diags.Append(mdm.Ips.ElementsAs(ctx, &configIps, true)...)
			if len(configIps) > 0 {
				if _, ok := ipmap[configIps[0]]; !ok {
					diags.AddAttributeError(
						path.Root("secondary_mdm"),
						"Please enter valid IP for secondary MDM",
						"Please enter valid IP for secondary MDM",
					)
				}
			}
		}
	}

	// Check if valid IP is provided for TB MDM
	for _, mdm := range plan.TieBreakerMdm {
		if !mdm.Ips.IsNull() {
			configIps := make([]string, 0)
			diags.Append(mdm.Ips.ElementsAs(ctx, &configIps, true)...)
			if len(configIps) > 0 {
				if _, ok := ipmap[configIps[0]]; !ok {
					diags.AddAttributeError(
						path.Root("tiebreaker_mdm"),
						"Please enter valid IP for tiebreaker MDM",
						"Please enter valid IP for tiebreaker MDM",
					)
				}
			}
		}
	}
	return diags
}

// PopulateMdmID populates MDM IDs based on IP
func (d *mdmClusterResource) PopulateMdmID(ctx context.Context, plan *models.MdmResourceModel, state *models.MdmResourceModel, flag bool) diag.Diagnostics {
	var dgs, diags diag.Diagnostics

	mdmDetails, err := d.system.GetMDMClusterDetails()
	if err != nil {
		diags.AddError(
			"Error getting MDM cluster details: ",
			"Error getting MDM cluster details: "+err.Error(),
		)
		return diags
	}

	if flag {
		state, dgs = helper.UpdateMdmClusterState(ctx, mdmDetails, plan, d.system.System.PerformanceProfile)
		diags.Append(dgs...)
	}

	// Get the IP with IDs as value
	ipmap := helper.GetMdmIPMap(mdmDetails)

	diags = d.ValidateConfig(ctx, ipmap, plan)
	if diags.HasError() {
		return diags
	}

	// Populate Primary MDM ID if IP is provided
	if plan.PrimaryMdm.ID.IsUnknown() {
		planIps := make([]string, 0)
		diags.Append(plan.PrimaryMdm.Ips.ElementsAs(ctx, &planIps, true)...)
		if val, ok := ipmap[planIps[0]]; ok {
			plan.PrimaryMdm.ID = types.StringValue(val)
		}
	}

	// Populate Secondary MDM IDs if IP is provided
	for index, mdm := range plan.SecondaryMdm {
		if mdm.ID.IsUnknown() {
			planIps := make([]string, 0)
			diags.Append(mdm.Ips.ElementsAs(ctx, &planIps, true)...)
			if val, ok := ipmap[planIps[0]]; ok {
				plan.SecondaryMdm[index].ID = types.StringValue(val)
			}
		}
	}

	// Populate TB MDM IDs if IP is provided
	for index, mdm := range plan.TieBreakerMdm {
		if mdm.ID.IsUnknown() {
			planIps := make([]string, 0)
			diags.Append(mdm.Ips.ElementsAs(ctx, &planIps, true)...)
			if val, ok := ipmap[planIps[0]]; ok {
				plan.TieBreakerMdm[index].ID = types.StringValue(val)
			}
		}
	}

	planStandby := make([]models.StandByMdm, 0)
	diags.Append(plan.StandByMdm.ElementsAs(ctx, &planStandby, true)...)

	// Populate StandBy MDM IDs if IP is provided
	for index, mdm := range planStandby {
		if mdm.ID.IsUnknown() {
			planIps := make([]string, 0)
			diags.Append(mdm.Ips.ElementsAs(ctx, &planIps, true)...)
			if val, ok := ipmap[planIps[0]]; ok {
				planStandby[index].ID = types.StringValue(val)
			}
		}
	}

	if len(planStandby) > 0 {
		standbyList, dgs := helper.GetStandByMdmSetValue(planStandby)
		diags.Append(dgs...)
		plan.StandByMdm = standbyList
	}

	if plan.PrimaryMdm.ID.ValueString() != state.PrimaryMdm.ID.ValueString() {
		for index, mdm := range state.SecondaryMdm {
			if mdm.ID.ValueString() == plan.PrimaryMdm.ID.ValueString() {
				if plan.PrimaryMdm.Name.IsUnknown() {
					plan.PrimaryMdm.Name = mdm.Name
				}

				if plan.PrimaryMdm.Ips.IsUnknown() {
					plan.PrimaryMdm.Ips = mdm.Ips
				}

				plan.PrimaryMdm.ManagementIps = mdm.ManagementIps
				plan.PrimaryMdm.Port = mdm.Port

				if plan.SecondaryMdm[index].Name.IsUnknown() {
					plan.SecondaryMdm[index].Name = state.PrimaryMdm.Name
				}

				if plan.SecondaryMdm[index].Ips.IsUnknown() {
					plan.SecondaryMdm[index].Ips = state.PrimaryMdm.Ips
				}

				plan.SecondaryMdm[index].ManagementIps = state.PrimaryMdm.ManagementIps
				plan.SecondaryMdm[index].Port = state.PrimaryMdm.Port
				break
			}
		}
	} else {
		if plan.PrimaryMdm.Name.IsUnknown() {
			plan.PrimaryMdm.Name = state.PrimaryMdm.Name
		}

		if plan.PrimaryMdm.Ips.IsUnknown() {
			plan.PrimaryMdm.Ips = state.PrimaryMdm.Ips
		}

		for index, mdm := range plan.SecondaryMdm {
			if index < len(state.SecondaryMdm) {
				if mdm.Name.IsUnknown() {
					plan.SecondaryMdm[index].Name = state.SecondaryMdm[index].Name
				}

				if mdm.Ips.IsUnknown() {
					plan.SecondaryMdm[index].Ips = state.SecondaryMdm[index].Ips
				}
			}
		}
	}
	return diags
}

// ModifyPlan modify resource plan attribute value
func (d *mdmClusterResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}

	var stateFlag bool
	if req.State.Raw.IsNull() {
		stateFlag = true
	}

	var plan models.MdmResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Reading state for update operation
	var state models.MdmResourceModel
	if !stateFlag {
		diags = req.State.Get(ctx, &state)
		resp.Diagnostics.Append(diags...)
	}

	resp.Diagnostics.Append(d.PopulateMdmID(ctx, &plan, &state, stateFlag)...)

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

func (d *mdmClusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Info(ctx, "In create operation")

	// Retrieve values from plan
	var plan models.MdmResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = d.PopulateMdmID(ctx, &plan, nil, true)
	if diags.HasError() {
		return
	}

	mdmDetails, err := d.system.GetMDMClusterDetails()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting MDM cluster details: ",
			"Error getting MDM cluster details: "+err.Error(),
		)
		return
	}

	state, dgs := helper.UpdateMdmClusterState(ctx, mdmDetails, &plan, d.system.System.PerformanceProfile)
	diags = append(diags, dgs...)

	// Perform the update operations based on the provided config
	resp.Diagnostics.Append(d.UpdateMdmClusterResource(ctx, plan, *state)...)

	mdmDetails, err = d.system.GetMDMClusterDetails()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting MDM cluster details: ",
			"Error getting MDM cluster details: "+err.Error(),
		)
		return
	}

	system, err := helper.GetFirstSystem(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex System",
			err.Error(),
		)
		return
	}

	state, dgs = helper.UpdateMdmClusterState(ctx, mdmDetails, &plan, system.System.PerformanceProfile)
	diags = append(diags, dgs...)
	resp.Diagnostics.Append(diags...)

	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
}

func (d *mdmClusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Info(ctx, "In read operation")
	// Retrieve values from state
	var state models.MdmResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	system, err := helper.GetFirstSystem(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex System",
			err.Error(),
		)
		return
	}

	mdmDetails, err := d.system.GetMDMClusterDetails()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting MDM cluster details: ",
			"Error getting MDM cluster details: "+err.Error(),
		)
		return
	}

	// Set refreshed state
	state1, dgs := helper.UpdateMdmClusterState(ctx, mdmDetails, &state, system.System.PerformanceProfile)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state1)
	resp.Diagnostics.Append(diags...)
}

func (d *mdmClusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Info(ctx, "In update operation")
	// Retrieve values from plan
	var plan models.MdmResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	// Retrieve values from state
	var state models.MdmResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	diags = d.PopulateMdmID(ctx, &plan, &state, false)
	if diags.HasError() {
		return
	}

	resp.Diagnostics.Append(d.UpdateMdmClusterResource(ctx, plan, state)...)

	system, err := helper.GetFirstSystem(d.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Unable to Read Powerflex System",
			err.Error(),
		)
		return
	}

	mdmDetails, err := d.system.GetMDMClusterDetails()
	if err != nil {
		resp.Diagnostics.AddError(
			"Error getting MDM cluster details: ",
			"Error getting MDM cluster details: "+err.Error(),
		)
		return
	}

	// Set refreshed state
	state1, dgs := helper.UpdateMdmClusterState(ctx, mdmDetails, &plan, system.System.PerformanceProfile)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, state1)
	resp.Diagnostics.Append(diags...)

}

func (d *mdmClusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Info(ctx, "In delete operation")
	// Retrieve values from state
	var state models.MdmResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
	resp.State.RemoveResource(ctx)
}

// UpdateMdmClusterResource performs the update operations
func (d *mdmClusterResource) UpdateMdmClusterResource(ctx context.Context, plan, state models.MdmResourceModel) diag.Diagnostics {
	var dia diag.Diagnostics

	// Modify performance profile if it differs in plan and state
	if !plan.PerformanceProfile.IsUnknown() && d.system.System.PerformanceProfile != plan.PerformanceProfile.ValueString() {
		err := d.system.ModifyPerformanceProfileMdmCluster(plan.PerformanceProfile.ValueString())
		if err != nil {
			dia.AddError("Could not change performance profile of MDM cluster.", err.Error())
		}
	}

	dia.Append(d.RenameMdms(ctx, plan, state)...)
	dia.Append(d.ChangeMdmOwnerShip(plan, state)...)
	dia.Append(d.SwitchClusterMode(plan, state)...)

	planMdmList := []models.StandByMdm{}
	stateMdmList := []models.StandByMdm{}

	// Populate planMdmList with standby MDMs defined in plan
	diags := plan.StandByMdm.ElementsAs(ctx, &planMdmList, true)
	dia.Append(diags...)

	// Populate stateMdmList with standby MDMs stored in state
	diags = state.StandByMdm.ElementsAs(ctx, &stateMdmList, true)
	dia.Append(diags...)

	dia.Append(d.AddStandByMdm(ctx, plan, state, planMdmList, stateMdmList)...)
	dia.Append(d.RemoveStandByMdm(ctx, plan, state, planMdmList, stateMdmList)...)

	return dia
}

// RenameMdms renames the MDMs
func (d *mdmClusterResource) RenameMdms(ctx context.Context, plan, state models.MdmResourceModel) diag.Diagnostics {
	var dia diag.Diagnostics
	renamePayload := make(map[string]string, 0)
	planMdm := make(map[string]string, 0)
	stateMdm := make(map[string]string, 0)
	existingNames := make(map[string]string, 0)

	for _, mdm := range state.SecondaryMdm {
		existingNames[mdm.Name.ValueString()] = mdm.Name.ValueString()
	}

	if !plan.PrimaryMdm.Name.IsUnknown() && plan.PrimaryMdm.Name.ValueString() != state.PrimaryMdm.Name.ValueString() && helper.CheckforExistingName(plan, existingNames) {
		renamePayload[plan.PrimaryMdm.ID.ValueString()] = plan.PrimaryMdm.Name.ValueString()
	}

	for _, mdm := range plan.SecondaryMdm {
		if !mdm.Name.IsUnknown() && mdm.Name.ValueString() != state.PrimaryMdm.Name.ValueString() {
			planMdm[mdm.ID.ValueString()] = mdm.Name.ValueString()
		}
	}

	for _, mdm := range state.SecondaryMdm {
		stateMdm[mdm.ID.ValueString()] = mdm.Name.ValueString()
	}

	for _, mdm := range plan.TieBreakerMdm {
		if !mdm.Name.IsUnknown() {
			planMdm[mdm.ID.ValueString()] = mdm.Name.ValueString()
		}
	}

	for _, mdm := range state.TieBreakerMdm {
		stateMdm[mdm.ID.ValueString()] = mdm.Name.ValueString()
	}

	planSbMdmList := []models.StandByMdm{}
	stateSbMdmList := []models.StandByMdm{}

	// Populate planMdmList with standby MDMs defined in plan
	diags := plan.StandByMdm.ElementsAs(ctx, &planSbMdmList, true)
	dia.Append(diags...)
	for _, mdm := range planSbMdmList {
		if !mdm.ID.IsUnknown() && !mdm.Name.IsUnknown() {
			planMdm[mdm.ID.ValueString()] = mdm.Name.ValueString()
		}
	}

	// Populate stateMdmList with standby MDMs stored in state
	diags = state.StandByMdm.ElementsAs(ctx, &stateSbMdmList, true)
	dia.Append(diags...)
	for _, mdm := range stateSbMdmList {
		stateMdm[mdm.ID.ValueString()] = mdm.Name.ValueString()
	}

	// Check the MDM names for all MDMs in plan and state if it differs
	for id, name := range planMdm {
		if val, ok := stateMdm[id]; ok && name != val && name != "" {
			renamePayload[id] = name
		}
	}

	// Rename MDMs
	for id, name := range renamePayload {
		payload := scaleio_types.RenameMdm{
			ID:      id,
			NewName: name,
		}
		err := d.system.RenameMdm(&payload)
		if err != nil {
			dia.AddError(
				fmt.Sprintf("Could not rename the MDM with ID %v and name %v", id, name), err.Error())
		}
	}
	return dia
}

// ChangeMdmOwnerShip modifies the primary MDM
func (d *mdmClusterResource) ChangeMdmOwnerShip(plan, state models.MdmResourceModel) diag.Diagnostics {
	var dia diag.Diagnostics
	if plan.PrimaryMdm.ID.ValueString() != state.PrimaryMdm.ID.ValueString() {
		err := d.system.ChangeMdmOwnerShip(plan.PrimaryMdm.ID.ValueString())
		if err != nil {
			dia.AddError(
				fmt.Sprintf("Could not change MDM ownership with ID %v.", plan.PrimaryMdm.ID.ValueString()), err.Error())
		}
		time.Sleep(15 * time.Second)
	}
	return dia
}

// SwitchClusterMode modifies the MDM cluster mode
func (d *mdmClusterResource) SwitchClusterMode(plan, state models.MdmResourceModel) diag.Diagnostics {
	var dia diag.Diagnostics
	if plan.ClusterMode.ValueString() != state.ClusterMode.ValueString() {
		planSecondaryMdm := make([]string, 0)
		stateSecondaryMdm := make([]string, 0)
		planTbMdm := make([]string, 0)
		stateTbMdm := make([]string, 0)
		for _, mdm := range plan.SecondaryMdm {
			planSecondaryMdm = append(planSecondaryMdm, mdm.ID.ValueString())
		}

		for _, mdm := range state.SecondaryMdm {
			stateSecondaryMdm = append(stateSecondaryMdm, mdm.ID.ValueString())
		}

		for _, mdm := range plan.TieBreakerMdm {
			planTbMdm = append(planTbMdm, mdm.ID.ValueString())
		}

		for _, mdm := range state.TieBreakerMdm {
			stateTbMdm = append(stateTbMdm, mdm.ID.ValueString())
		}

		// Exapnd the MDM cluster from 3 nodes to 5 nodes
		if plan.ClusterMode.ValueString() == scaleio_types.FiveNodesClusterMode {
			addSecondary := helper.Difference(planSecondaryMdm, stateSecondaryMdm)
			addTb := helper.Difference(planTbMdm, stateTbMdm)
			payload := scaleio_types.SwitchClusterMode{
				Mode:             scaleio_types.FiveNodesClusterMode,
				AddSecondaryMdms: addSecondary,
				AddTBMdms:        addTb,
			}
			err := d.system.SwitchClusterMode(&payload)
			if err != nil {
				dia.AddError("Could not expand the MDM cluster ", err.Error())
			}
		} else {
			// Reduce the MDM cluster from 5 nodes to 3 nodes
			removeSecondary := helper.Difference(stateSecondaryMdm, planSecondaryMdm)
			removeTb := helper.Difference(stateTbMdm, planTbMdm)
			payload := scaleio_types.SwitchClusterMode{
				Mode:                scaleio_types.ThreeNodesClusterMode,
				RemoveSecondaryMdms: removeSecondary,
				RemoveTBMdms:        removeTb,
			}
			err := d.system.SwitchClusterMode(&payload)
			if err != nil {
				dia.AddError("Could not reduce the MDM cluster ", err.Error())
			}
		}
	}
	return dia
}

// AddStandByMdm adds standby MDMs
func (d *mdmClusterResource) AddStandByMdm(ctx context.Context, plan models.MdmResourceModel, state models.MdmResourceModel, planMdmList, stateMdmList []models.StandByMdm) diag.Diagnostics {
	var dia diag.Diagnostics

	addStandby, diags1 := helper.StandByMdmDifference(ctx, planMdmList, stateMdmList)
	dia.Append(diags1...)

	addStandby, diags1 = helper.CheckForSwitchCluster(ctx, addStandby, state.SecondaryMdm, state.TieBreakerMdm)
	dia.Append(diags1...)

	// Add standby MDMs
	for _, mdm := range addStandby {
		ips := make([]string, 0)
		mgmtIps := make([]string, 0)
		dia.Append(mdm.Ips.ElementsAs(ctx, &ips, true)...)
		dia.Append(mdm.ManagementIps.ElementsAs(ctx, &mgmtIps, true)...)

		payload := scaleio_types.StandByMdm{
			Name:          mdm.Name.ValueString(),
			IPs:           ips,
			Role:          mdm.Role.ValueString(),
			ManagementIPs: mgmtIps,
		}

		// Allow asymmetric Ips
		if mdm.AllowAsymmetricIps.ValueBool() {
			payload.AllowAsymmetricIps = "true"
		} else {
			payload.AllowAsymmetricIps = "false"
		}

		_, err := d.system.AddStandByMdm(&payload)
		if err != nil {
			dia.AddError(
				fmt.Sprintf("Could not add standby MDM with IP %v ", ips[0]), err.Error())
		}
	}
	return dia
}

// RemoveStandByMdm removes standby MDMs
func (d *mdmClusterResource) RemoveStandByMdm(ctx context.Context, plan models.MdmResourceModel, state models.MdmResourceModel, planMdmList, stateMdmList []models.StandByMdm) diag.Diagnostics {
	var dia diag.Diagnostics

	deleteStandby, diags := helper.StandByMdmDifference(ctx, stateMdmList, planMdmList)
	dia.Append(diags...)
	deleteStandby, diags = helper.CheckForSwitchCluster(ctx, deleteStandby, plan.SecondaryMdm, plan.TieBreakerMdm)
	dia.Append(diags...)

	for _, mdm := range deleteStandby {
		err := d.system.RemoveStandByMdm(mdm.ID.ValueString())
		if err != nil {
			dia.AddError(
				fmt.Sprintf("Could not remove standby MDM with ID %v ", mdm.ID.ValueString()), err.Error())
		}
	}
	return dia
}
