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
	"strconv"
	"strings"
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-framework/diag"
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

// sdcResource - struct to define sdc resource
type sdcResource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

// Metadata - function to return metadata for SDC resource.
func (r *sdcResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_sdc"
}

// Schema - function to return Schema for SDC resource.
func (r *sdcResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = SDCReourceSchema
}

// Configure - function to return Configuration for SDC resource.
func (r *sdcResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client != nil {
		r.client = req.ProviderData.(*powerflexProvider).client

		r.gatewayClient = req.ProviderData.(*powerflexProvider).gatewayClient
	} else {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)

		return
	}
}

// ModifyPlan modify resource plan attribute value
func (r *sdcResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	var plan models.SdcResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)

	planSdcs := make([]models.SDCDetailDataModel, 0)
	diags.Append(plan.SDCDetails.ElementsAs(ctx, &planSdcs, true)...)

	for index, sdc := range planSdcs {
		if !sdc.SDCID.IsNull() {
			system, err := helper.GetFirstSystem(r.client)

			if err != nil {
				resp.Diagnostics.AddError(
					"Error in getting system instance on the PowerFlex cluster",
					err.Error(),
				)
				return
			}

			sdcData, err := system.GetSdcByID(sdc.SDCID.ValueString())

			if err != nil {
				resp.Diagnostics.AddError(
					"[Read] Unable to Find SDC by ID:"+sdc.SDCID.ValueString(),
					err.Error(),
				)
				return
			}
			planSdcs[index].IP = types.StringValue(sdcData.Sdc.SdcIP)
		}
	}

	sdcList, dgs := helper.GetSdcsValue(planSdcs)
	diags.Append(dgs...)
	plan.SDCDetails = sdcList

	diags = resp.Plan.Set(ctx, &plan)
	resp.Diagnostics.Append(diags...)
}

// Create - function to Create for SDC resource.
func (r *sdcResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Create")

	var plan models.SdcResourceModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sdcDetailList := []models.SDCDetailDataModel{}
	diags = plan.SDCDetails.ElementsAs(ctx, &sdcDetailList, true)
	resp.Diagnostics.Append(diags...)

	system, err := helper.GetFirstSystem(r.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}

	var chnagedSDCs []models.SDCStateDataModel

	if len(sdcDetailList) > 0 {

		resp.Diagnostics.Append(r.SDCExpansionOperations(ctx, plan, system, sdcDetailList, nil)...)
		if resp.Diagnostics.HasError() {
			return
		}

		resp.Diagnostics.Append(r.UpdateSDCNamdPerfProfileOperations(ctx, sdcDetailList, system, &chnagedSDCs)...)

		data, dgs := helper.UpdateState(chnagedSDCs, plan)
		resp.Diagnostics.Append(dgs...)

		diags = resp.State.Set(ctx, data)
		resp.Diagnostics.Append(diags...)

		tflog.Info(ctx, "SDC Details updated to state file successfully")

		return
	}

	resp.Diagnostics.AddError(
		"Please provide valid inputs",
		"Please provide valid inputs",
	)
}

// Read - function to Read for SDC resource.
func (r *sdcResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Read")
	var state models.SdcResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sdcStateList := []models.SDCStateDataModel{}
	diags = state.SDCStateDetails.ElementsAs(ctx, &sdcStateList, true)
	resp.Diagnostics.Append(diags...)

	system, err := helper.GetFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}

	var chnagedSDCs []models.SDCStateDataModel

	//For handling the import case
	if state.ID.ValueString() != "" && state.ID.ValueString() != "placeholder" {

		for _, id := range strings.Split(state.ID.ValueString(), ",") {
			sdcData, err := system.GetSdcByID(id)

			if err != nil {
				resp.Diagnostics.AddError(
					"[Import] Unable to Find SDC by ID:"+id,
					err.Error(),
				)
				return
			}

			if sdcData != nil {
				changedSDCDetail := helper.GetSDCState(*sdcData.Sdc)

				chnagedSDCs = append(chnagedSDCs, changedSDCDetail)
			}
		}
	} else if len(sdcStateList) > 0 {

		//For handling the multiple sdc_details update
		for _, sdc := range sdcStateList {
			var sdcData *goscaleio.Sdc

			if sdc.SDCID.ValueString() != "" {
				sdcData, err = system.GetSdcByID(sdc.SDCID.ValueString())

				if err != nil {
					resp.Diagnostics.AddError(
						"[Read] Unable to Find SDC by ID:"+sdc.SDCID.ValueString(),
						err.Error(),
					)
				}
			} else if sdc.IP.ValueString() != "" {
				sdcData, err = system.FindSdc("SdcIP", sdc.IP.ValueString())

				if err != nil {
					resp.Diagnostics.AddError(
						"[Read] Unable to Find SDC by IP:"+sdc.IP.ValueString(),
						err.Error(),
					)
				}
			} else if sdc.SDCName.ValueString() != "" {
				sdcData, err = system.FindSdc("Name", sdc.SDCName.ValueString())

				if err != nil {
					resp.Diagnostics.AddError(
						"[Read] Unable to Find SDC by Name:"+sdc.SDCName.ValueString(),
						err.Error(),
					)
				}
			}

			if sdcData != nil {
				changedSDCDetail := helper.GetSDCState(*sdcData.Sdc)

				chnagedSDCs = append(chnagedSDCs, changedSDCDetail)
			}
		}
	} else if state.ID.ValueString() == "" {
		resp.Diagnostics.AddError("[Read] Please provide valid SDC ID", "Please provide valid SDC ID")
		return
	}

	data, dgs := helper.UpdateState(chnagedSDCs, state)
	resp.Diagnostics.Append(dgs...)

	diags = resp.State.Set(ctx, data)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

// Update - function to Update for SDC resource.
func (r *sdcResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Update")
	// Retrieve values from plan
	var plan models.SdcResourceModel
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state models.SdcResourceModel
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	system, err := helper.GetFirstSystem(r.client)

	planSdcDetailList := []models.SDCDetailDataModel{}
	diags = plan.SDCDetails.ElementsAs(ctx, &planSdcDetailList, true)
	resp.Diagnostics.Append(diags...)

	stateSdcDetailList := []models.SDCStateDataModel{}
	diags = state.SDCStateDetails.ElementsAs(ctx, &stateSdcDetailList, true)
	resp.Diagnostics.Append(diags...)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}

	var chnagedSDCs []models.SDCStateDataModel
	deletedSDC := helper.FindDeletedSDC(stateSdcDetailList, planSdcDetailList)

	if len(deletedSDC) > 0 {
		for _, sdc := range deletedSDC {
			err := system.DeleteSdc(sdc.SDCID.ValueString())

			if err != nil {
				resp.Diagnostics.AddError(
					"[Update] Unable to Delete SDC by ID:"+sdc.SDCID.ValueString(),
					err.Error(),
				)
				return
			}
		}
	}

	if len(planSdcDetailList) > 0 {

		resp.Diagnostics.Append(r.SDCExpansionOperations(ctx, plan, system, planSdcDetailList, stateSdcDetailList)...)
		if resp.Diagnostics.HasError() {

			//Handling the existing state file data
			for _, sdc := range planSdcDetailList {

				if strings.EqualFold(sdc.IsSdc.ValueString(), "Yes") {
					sdcData, _ := system.FindSdc("SdcIP", sdc.IP.ValueString())

					if sdcData != nil {
						changedSDCDetail := helper.GetSDCState(*sdcData.Sdc)

						chnagedSDCs = append(chnagedSDCs, changedSDCDetail)
					}
				}
			}

			data, dgs := helper.UpdateState(chnagedSDCs, plan)
			resp.Diagnostics.Append(dgs...)

			diags = resp.State.Set(ctx, data)
			resp.Diagnostics.Append(diags...)

			return
		}

		resp.Diagnostics.Append(r.UpdateSDCNamdPerfProfileOperations(ctx, planSdcDetailList, system, &chnagedSDCs)...)

		data, dgs := helper.UpdateState(chnagedSDCs, plan)
		resp.Diagnostics.Append(dgs...)

		diags = resp.State.Set(ctx, data)
		resp.Diagnostics.Append(diags...)

		tflog.Info(ctx, "SDC Details updated to state file successfully")

		return
	}

	// This code addresses the scenario when all SDCs needs to be deleted
	plan.SDCStateDetails = types.ListNull(types.ObjectType{
		AttrTypes: helper.GetSDCStateDetailType(),
	})
	diags = resp.State.Set(ctx, plan)
	resp.Diagnostics.Append(diags...)

	tflog.Info(ctx, "SDC Details deleted from state file successfully")

}

// Delete - function to Delete for SDC resource.
func (r *sdcResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Delete")
	var state models.SdcResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	sdcDetailList := []models.SDCStateDataModel{}
	diags = state.SDCStateDetails.ElementsAs(ctx, &sdcDetailList, true)
	resp.Diagnostics.Append(diags...)

	system, err := helper.GetFirstSystem(r.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}

	for _, sdc := range sdcDetailList {

		err := system.DeleteSdc(sdc.SDCID.ValueString())

		if err != nil {
			resp.Diagnostics.AddError(
				"[Delete] Unable to Delete SDC by ID:"+sdc.SDCID.ValueString(),
				err.Error(),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)

}

// ImportState - function to ImportState for SDC resource.
func (r *sdcResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] ImportState :-- "+helper.PrettyJSON(req))
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// SDCExpansionOperations function for the SDC Expansion Operation Like ParseCSV, Validate MDM and Installation
func (r *sdcResource) SDCExpansionOperations(ctx context.Context, plan models.SdcResourceModel, system *goscaleio.System, sdcDetails []models.SDCDetailDataModel, sdcStateDetails []models.SDCStateDataModel) (dia diag.Diagnostics) {

	if helper.CheckForExpansion(sdcDetails, sdcStateDetails) {
		parsecsvRespose, parseCSVError := helper.ParseCSVOperation(ctx, sdcDetails, r.gatewayClient)

		if parseCSVError != nil {
			dia.AddError(
				"Error while Parsing CSV",
				"unexpected error: "+parseCSVError.Error(),
			)
			return
		}

		// to make gateway available for installation
		queueOperationError := helper.ResetInstallerQueue(r.gatewayClient)
		if queueOperationError != nil {
			dia.AddError(
				"Error Clearing Queue",
				"unexpected error: "+queueOperationError.Error(),
			)
			return
		}

		tflog.Info(ctx, "Gateway Installer changed to idle phase before initiating process")

		mdmIP, mdmIPError := helper.GetMDMIP(ctx, sdcDetails)
		if mdmIPError != nil {
			dia.AddError(
				"Error while Getting MDM IP",
				"unexpected error: "+mdmIPError.Error(),
			)
			return
		}

		tflog.Info(ctx, "CSV File parsed ssucessfully")

		// Vaidate the MDM credentials
		validateMDMResponse, validateMDMError := helper.ValidateMDMOperation(ctx, plan, r.gatewayClient, mdmIP)
		if validateMDMError != nil {
			dia.AddError(
				"Error While Validating MDM Details",
				"unexpected error: "+validateMDMResponse.Message,
			)
			return
		}

		if validateMDMResponse.StatusCode == 200 {

			tflog.Info(ctx, "MDM Details validated successfully")

			if !helper.CheckForNewSDCIPs(strings.Split(parsecsvRespose.Message, ","), strings.Split(validateMDMResponse.Data, ",")) {
				installationError := helper.InstallationOperations(ctx, plan, r.gatewayClient, parsecsvRespose)

				if installationError != nil {
					dia.AddError(
						"Error in Installation Process",
						"unexpected error: "+installationError.Error(),
					)
					return
				}
			}
		} else if validateMDMResponse.StatusCode != 200 {
			dia.AddError(
				"Error While Validating MDM Credentials",
				"unexpected error: "+validateMDMResponse.Message+" & Status Code: "+strconv.Itoa(validateMDMResponse.StatusCode),
			)
			return
		}
	}

	return
}

// UpdateSDCNamdPerfProfileOperations function for Update Name and Performance Profile of SDC
func (r *sdcResource) UpdateSDCNamdPerfProfileOperations(ctx context.Context, sdcDetailList []models.SDCDetailDataModel, system *goscaleio.System, chnagedSDCs *[]models.SDCStateDataModel) (dia diag.Diagnostics) {

	for _, sdc := range sdcDetailList {
		var (
			sdcData *goscaleio.Sdc
			err     error
		)

		if strings.EqualFold(sdc.IsSdc.ValueString(), "Yes") && sdc.SDCName.ValueString() == "" && sdc.PerformanceProfile.ValueString() == "" {

			if !sdc.DataNetworkIP.IsNull() {
				sdcData, err = system.FindSdc("SdcIP", sdc.DataNetworkIP.ValueString())
				if err != nil {
					dia.AddError(
						"[Create] Unable to Find SDC by IP: "+sdc.DataNetworkIP.ValueString(),
						err.Error(),
					)
				}
			} else {
				sdcData, err = system.FindSdc("SdcIP", sdc.IP.ValueString())
				if err != nil {
					dia.AddError(
						"[Create] Unable to Find SDC by IP: "+sdc.IP.ValueString(),
						err.Error(),
					)
				}
			}

			if sdcData != nil {
				changedSDCDetail := helper.GetSDCState(*sdcData.Sdc)

				*chnagedSDCs = append(*chnagedSDCs, changedSDCDetail)
			}
		} else if strings.EqualFold(sdc.IsSdc.ValueString(), "Yes") && (sdc.SDCName.ValueString() != "" || sdc.PerformanceProfile.ValueString() != "") {
			if sdc.SDCID.ValueString() == "" && sdc.IP.ValueString() != "" {
				var sdcID string

				if !sdc.DataNetworkIP.IsNull() {
					sdcID, err = system.GetSdcIDByIP(sdc.DataNetworkIP.ValueString())
					if err != nil {
						dia.AddError(
							"[Create] Unable to Find SDC by IP:"+sdc.DataNetworkIP.ValueString(),
							err.Error(),
						)
					}
				} else {
					sdcID, err = system.GetSdcIDByIP(sdc.IP.ValueString())
					if err != nil {
						dia.AddError(
							"[Create] Unable to Find SDC by IP:"+sdc.IP.ValueString(),
							err.Error(),
						)
					}
				}

				sdc.SDCID = types.StringValue(sdcID)
			}

			if sdc.SDCID.ValueString() == "" && sdc.SDCName.ValueString() != "" {
				sdcData, err := system.FindSdc("Name", sdc.SDCName.ValueString())

				if err != nil {
					dia.AddError(
						"[Create] Unable to Find SDC by Name:"+sdc.SDCName.ValueString(),
						err.Error(),
					)
				}

				sdc.SDCID = types.StringValue(sdcData.Sdc.ID)
			}

			if sdc.SDCName.ValueString() != "" && sdc.SDCID.ValueString() != "" {

				nameExist, _ := helper.CheckForSDCName(system, sdc)

				if !nameExist {
					nameChng, err := system.ChangeSdcName(sdc.SDCID.ValueString(), sdc.SDCName.ValueString())

					if err != nil {
						dia.AddError(
							"[Create] Unable to Change Name Powerflex SDC by ID:"+sdc.SDCID.ValueString()+" Name:"+sdc.SDCName.ValueString(),
							err.Error(),
						)
					}

					tflog.Debug(ctx, fmt.Sprintf("[POWERFLEX] Name Change Result: %s  SDC ID: %s", helper.PrettyJSON(nameChng), sdc.SDCID))
				}
			}

			if sdc.PerformanceProfile.ValueString() != "" && sdc.SDCID.ValueString() != "" {
				perProfile, err := system.ChangeSdcPerfProfile(sdc.SDCID.ValueString(), sdc.PerformanceProfile.ValueString())

				if err != nil {
					dia.AddError(
						"[Create] Unable to Change Performance Profile Powerflex SDC by ID:"+sdc.SDCID.ValueString(),
						err.Error(),
					)
				}

				tflog.Debug(ctx, fmt.Sprintf("[POWERFLEX] Performance Profile Change Result: %s  SDC ID: %s", helper.PrettyJSON(perProfile), sdc.SDCID))
			}

			finalSDC, err := system.GetSdcByID(sdc.SDCID.ValueString())
			if err != nil {
				dia.AddError(
					"Unable to Read Changed SDC",
					err.Error(),
				)
				return
			}

			if finalSDC != nil {
				changedSDCDetail := helper.GetSDCState(*finalSDC.Sdc)
				*chnagedSDCs = append(*chnagedSDCs, changedSDCDetail)
			}
		}
	}

	return
}
