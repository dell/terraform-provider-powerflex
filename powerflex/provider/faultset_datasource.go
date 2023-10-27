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

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ datasource.DataSource              = &faultSetDataSource{}
	_ datasource.DataSourceWithConfigure = &faultSetDataSource{}
)

// FaultSetDataSource returns the FaultSet data source
func FaultSetDataSource() datasource.DataSource {
	return &faultSetDataSource{}
}

type faultSetDataSource struct {
	client *goscaleio.Client
	system *goscaleio.System
}

func (d *faultSetDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_faultset"
}

func (d *faultSetDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = FaultSetDataSourceSchema
}

func (d *faultSetDataSource) Configure(_ context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
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

// GetSdsDetails fetches the SDS details associated with fault set
func (d *faultSetDataSource) GetSdsDetails(id string) ([]scaleiotypes.Sds, error) {
	sdsDetails, err := d.system.GetAllFaultSetsSds(id)
	if err != nil {
		return nil, err
	}
	return sdsDetails, nil
}

// Read refreshes the Terraform state with the latest data.
func (d *faultSetDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	tflog.Info(ctx, "Started fault set data source read method")
	var (
		state          models.FaultSetDataSourceModel
		faultSets      []scaleiotypes.FaultSet
		err            error
		faultSetsModel []models.FaultSet
	)

	diags := req.Config.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Fetch Fault set details if IDs are provided
	if !state.FaultSetIDs.IsNull() {
		faultSetIDs := make([]string, 0)
		diags.Append(state.FaultSetIDs.ElementsAs(ctx, &faultSetIDs, true)...)

		for _, faultSetID := range faultSetIDs {
			faultSet, err := d.system.GetFaultSetByID(faultSetID)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting faultset details using id %v", faultSetID), err.Error(),
				)
				return
			}
			sdsDetails, err := d.GetSdsDetails(faultSet.ID)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting SDS details connected to faultset details using id %v", faultSet.ID), err.Error(),
				)
				return
			}

			var sdsStateModels []models.SdsDataModel
			for _, sds := range sdsDetails {
				sdsState := getSdsState(&sds)
				sdsStateModels = append(sdsStateModels, sdsState)
			}
			faultSetsModel = append(faultSetsModel, helper.GetAllFaultSetState(*faultSet, sdsStateModels))
		}
	} else if !state.FaultSetNames.IsNull() {
		// Fetch Fault set details if names are provided
		faultSetNames := make([]string, 0)
		diags.Append(state.FaultSetNames.ElementsAs(ctx, &faultSetNames, true)...)

		for _, faultSetName := range faultSetNames {
			faultSet, err := d.system.GetFaultSetByName(faultSetName)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting faultset details using name %v", faultSetName), err.Error(),
				)
				return
			}
			sdsDetails, err := d.GetSdsDetails(faultSet.ID)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting SDS details connected to faultset details using id %v", faultSet.ID), err.Error(),
				)
				return
			}

			var sdsStateModels []models.SdsDataModel
			for _, sds := range sdsDetails {
				sdsState := getSdsState(&sds)
				sdsStateModels = append(sdsStateModels, sdsState)
			}
			faultSetsModel = append(faultSetsModel, helper.GetAllFaultSetState(*faultSet, sdsStateModels))
		}
	} else {
		// Fetch Fault set details for all the fault sets
		faultSets, err = d.system.GetAllFaultSets()
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting vTree details", err.Error(),
			)
			return
		}

		for _, faultSet := range faultSets {
			sdsDetails, err := d.GetSdsDetails(faultSet.ID)
			if err != nil {
				resp.Diagnostics.AddError(
					fmt.Sprintf("Error in getting SDS details connected to faultset details using id %v", faultSet.ID), err.Error(),
				)
				return
			}

			var sdsStateModels []models.SdsDataModel
			for _, sds := range sdsDetails {
				sdsState := getSdsState(&sds)
				sdsStateModels = append(sdsStateModels, sdsState)
			}
			faultSetsModel = append(faultSetsModel, helper.GetAllFaultSetState(faultSet, sdsStateModels))
		}
	}

	state.FaultSetDetails = faultSetsModel
	state.ID = types.StringValue("placeholder")
	diags = resp.State.Set(ctx, state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// FaultSetDataSourceSchema defines the schema for Fault Set datasource
var FaultSetDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing fault set from PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	MarkdownDescription: "This datasource is used to query the existing fault set from PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder attribute.",
			MarkdownDescription: "Placeholder attribute.",
			Computed:            true,
		},
		"faultset_ids": schema.SetAttribute{
			Description:         "List of fault set IDs",
			MarkdownDescription: "List of fault set IDs",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.ConflictsWith(
					path.MatchRoot("faultset_names"),
				),
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
		},
		"faultset_names": schema.SetAttribute{
			Description:         "List of fault set names",
			MarkdownDescription: "List of fault set names",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
			},
		},
		"faultset_details": schema.SetNestedAttribute{
			Description:         "Fault set details",
			MarkdownDescription: "Fault set details",
			Computed:            true,
			NestedObject: schema.NestedAttributeObject{
				Attributes: map[string]schema.Attribute{
					"protection_domain_id": schema.StringAttribute{
						MarkdownDescription: "Protection Domain ID",
						Description:         "Protection Domain ID",
						Computed:            true,
					},
					"name": schema.StringAttribute{
						MarkdownDescription: "Fault set name",
						Description:         "Fault set name",
						Computed:            true,
					},
					"id": schema.StringAttribute{
						MarkdownDescription: "Fault set ID",
						Description:         "Fault set ID",
						Computed:            true,
					},
					"links": schema.ListNestedAttribute{
						MarkdownDescription: "Specifies the links asscociated with fault set",
						Description:         "Specifies the links asscociated with fault set",
						Computed:            true,
						NestedObject:        schema.NestedAttributeObject{Attributes: FaultSetLinksSchema()},
					},
					"sds_details": schema.ListNestedAttribute{
						Description:         "List of fetched SDS.",
						MarkdownDescription: "List of fetched SDS.",
						Computed:            true,
						NestedObject: schema.NestedAttributeObject{
							Attributes: map[string]schema.Attribute{
								"id": schema.StringAttribute{
									Description:         "SDS ID.",
									MarkdownDescription: "SDS ID.",
									Computed:            true,
								},
								"name": schema.StringAttribute{
									Description:         "SDS name.",
									MarkdownDescription: "SDS name.",
									Computed:            true,
								},
								"ip_list": schema.ListNestedAttribute{
									Description:         "List of IPs associated with SDS.",
									MarkdownDescription: "List of IPs associated with SDS.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"ip": schema.StringAttribute{
												Description:         "SDS IP.",
												MarkdownDescription: "SDS IP.",
												Computed:            true,
											},
											"role": schema.StringAttribute{
												Description:         "SDS IP role.",
												MarkdownDescription: "SDS IP role.",
												Computed:            true,
											},
										},
									},
								},
								"port": schema.Int64Attribute{
									Description:         "SDS port.",
									MarkdownDescription: "SDS port.",
									Computed:            true,
								},
								"sds_state": schema.StringAttribute{
									Description:         "SDS state.",
									MarkdownDescription: "SDS state.",
									Computed:            true,
								},
								"membership_state": schema.StringAttribute{
									Description:         "Membership state.",
									MarkdownDescription: "Membership state.",
									Computed:            true,
								},
								"mdm_connection_state": schema.StringAttribute{
									Description:         "MDM connection state.",
									MarkdownDescription: "MDM connection state.",
									Computed:            true,
								},
								"drl_mode": schema.StringAttribute{
									Description:         "DRL mode.",
									MarkdownDescription: "DRL mode.",
									Computed:            true,
								},
								"rmcache_enabled": schema.BoolAttribute{
									Description:         "Whether RM cache is enabled or not.",
									MarkdownDescription: "Whether RM cache is enabled or not.",
									Computed:            true,
								},
								"rmcache_size": schema.Int64Attribute{
									Description:         "Indicates the size of Read RAM Cache on the specified SDS in KB.",
									MarkdownDescription: "Indicates the size of Read RAM Cache on the specified SDS in KB.",
									Computed:            true,
								},
								"rmcache_frozen": schema.BoolAttribute{
									Description:         "Indicates whether the Read RAM Cache is currently temporarily not in use.",
									MarkdownDescription: "Indicates whether the Read RAM Cache is currently temporarily not in use.",
									Computed:            true,
								},
								"on_vmware": schema.BoolAttribute{
									Description:         "Presence on VMware.",
									MarkdownDescription: "Presence on VMware.",
									Computed:            true,
								},
								"faultset_id": schema.StringAttribute{
									Description:         "Fault set ID.",
									MarkdownDescription: "Fault set ID.",
									Computed:            true,
								},
								"num_io_buffers": schema.Int64Attribute{
									Description:         "Number of IO buffers.",
									MarkdownDescription: "Number of IO buffers.",
									Computed:            true,
								},
								"rmcache_memory_allocation_state": schema.StringAttribute{
									Description:         "Indicates the state of the memory allocation process. Can be one of 'in progress' and 'done'.",
									MarkdownDescription: "Indicates the state of the memory allocation process. Can be one of `in progress` and `done`.",
									Computed:            true,
								},
								"performance_profile": schema.StringAttribute{
									Description:         "Performance profile.",
									MarkdownDescription: "Performance profile.",
									Computed:            true,
								},
								"software_version_info": schema.StringAttribute{
									Description:         "Software version information.",
									MarkdownDescription: "Software version information.",
									Computed:            true,
								},
								"configured_drl_mode": schema.StringAttribute{
									Description:         "Configured DRL mode.",
									MarkdownDescription: "Configured DRL mode.",
									Computed:            true,
								},
								"rfcache_enabled": schema.BoolAttribute{
									Description:         "Whether RF cache is enabled or not.",
									MarkdownDescription: "Whether RF cache is enabled or not.",
									Computed:            true,
								},
								"maintenance_state": schema.StringAttribute{
									Description:         "Maintenance state.",
									MarkdownDescription: "Maintenance state.",
									Computed:            true,
								},
								"maintenance_type": schema.StringAttribute{
									Description:         "Maintenance type.",
									MarkdownDescription: "Maintenance type.",
									Computed:            true,
								},
								"rfcache_error_low_resources": schema.BoolAttribute{
									Description:         "RF cache error for low resources.",
									MarkdownDescription: "RF cache error for low resources.",
									Computed:            true,
								},
								"rfcache_error_api_version_mismatch": schema.BoolAttribute{
									Description:         "RF cache error for API version mismatch.",
									MarkdownDescription: "RF cache error for API version mismatch.",
									Computed:            true,
								},
								"rfcache_error_inconsistent_cache_configuration": schema.BoolAttribute{
									Description:         "RF cache error for inconsistent cache configuration.",
									MarkdownDescription: "RF cache error for inconsistent cache configuration.",
									Computed:            true,
								},
								"rfcache_error_inconsistent_source_configuration": schema.BoolAttribute{
									Description:         "RF cache error for inconsistent source configuration.",
									MarkdownDescription: "RF cache error for inconsistent source configuration.",
									Computed:            true,
								},
								"rfcache_error_invalid_driver_path": schema.BoolAttribute{
									Description:         "RF cache error for invalid driver path.",
									MarkdownDescription: "RF cache error for invalid driver path.",
									Computed:            true,
								},
								"rfcache_error_device_does_not_exist": schema.BoolAttribute{
									Description:         "RF cache error for device does not exist.",
									MarkdownDescription: "RF cache error for device does not exist.",
									Computed:            true,
								},
								"authentication_error": schema.StringAttribute{
									Description:         "Authentication error.",
									MarkdownDescription: "Authentication error.",
									Computed:            true,
								},
								"fgl_num_concurrent_writes": schema.Int64Attribute{
									Description:         "FGL concurrent writes.",
									MarkdownDescription: "FGL concurrent writes.",
									Computed:            true,
								},
								"fgl_metadata_cache_state": schema.StringAttribute{
									Description:         "FGL metadata cache state.",
									MarkdownDescription: "FGL metadata cache state.",
									Computed:            true,
								},
								"fgl_metadata_cache_size": schema.Int64Attribute{
									Description:         "FGL metadata cache size.",
									MarkdownDescription: "FGL metadata cache size.",
									Computed:            true,
								},
								"num_restarts": schema.Int64Attribute{
									Description:         "Number of restarts.",
									MarkdownDescription: "Number of restarts.",
									Computed:            true,
								},
								"last_upgrade_time": schema.Int64Attribute{
									Description:         "Last time SDS was upgraded.",
									MarkdownDescription: "Last time SDS was upgraded.",
									Computed:            true,
								},
								"sds_decoupled": schema.SingleNestedAttribute{
									Description:         "SDS decoupled windows.",
									MarkdownDescription: "SDS decoupled windows.",
									Computed:            true,
									Attributes:          getSdsAllWindowParamsSchema(),
								},
								"sds_configuration_failure": schema.SingleNestedAttribute{
									Description:         "SDS configuration failure windows.",
									MarkdownDescription: "SDS configuration failure windows.",
									Computed:            true,
									Attributes:          getSdsAllWindowParamsSchema(),
								},
								"sds_receive_buffer_allocation_failures": schema.SingleNestedAttribute{
									Description:         "SDS receive buffer allocation failure windows.",
									MarkdownDescription: "SDS receive buffer allocation failure windows.",
									Computed:            true,
									Attributes:          getSdsAllWindowParamsSchema(),
								},
								"certificate_info": schema.SingleNestedAttribute{
									Description:         "Certificate Information.",
									MarkdownDescription: "Certificate Information.",
									Computed:            true,
									Attributes:          getCertificateInfoSchema(),
								},
								"raid_controllers": schema.ListNestedAttribute{
									Description:         "RAID controllers information.",
									MarkdownDescription: "RAID controllers information.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: getRaidControllersSchema(),
									},
								},
								"links": schema.ListNestedAttribute{
									Description:         "Specifies the links asscociated with SDS.",
									MarkdownDescription: "Specifies the links asscociated with SDS.",
									Computed:            true,
									NestedObject: schema.NestedAttributeObject{
										Attributes: map[string]schema.Attribute{
											"rel": schema.StringAttribute{
												Description:         "Specifies the relationship with the SDS.",
												MarkdownDescription: "Specifies the relationship with the SDS.",
												Computed:            true,
											},
											"href": schema.StringAttribute{
												Description:         "Specifies the exact path to fetch the details.",
												MarkdownDescription: "Specifies the exact path to fetch the details.",
												Computed:            true,
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	},
}

// FaultSetLinksSchema specifies the schema for fault set links
func FaultSetLinksSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"rel": schema.StringAttribute{
			MarkdownDescription: "Specifies the relationship with the fault set",
			Description:         "Specifies the relationship with the fault set",
			Computed:            true,
		},
		"href": schema.StringAttribute{
			MarkdownDescription: "Specifies the exact path to fetch the details",
			Description:         "Specifies the exact path to fetch the details",
			Computed:            true,
		},
	}
}
