/*
Copyright (c) 2022-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	goscaleio_types "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

var (
	_ resource.Resource                = &hostResource{}
	_ resource.ResourceWithConfigure   = &hostResource{}
	_ resource.ResourceWithImportState = &hostResource{}
)

// HostResource - function to return resource interface
func HostResource() resource.Resource {
	return &hostResource{}
}

// hostResource - struct to define Host Resource
type hostResource struct {
	client *goscaleio.Client
}

// Metadata - function to return metadata for Host Resource.
func (r *hostResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_host"
}

// Schema - function to return Schema for Host Resource.
func (r *hostResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = HostReourceSchema
}

// Configure - function to return Configuration for Host Resource.
func (r *hostResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client != nil {
		r.client = req.ProviderData.(*powerflexProvider).client
	} else {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)

		return
	}
}

// ModifyPlan modify resource plan attribute value
func (r *hostResource) ModifyPlan(ctx context.Context, req resource.ModifyPlanRequest, resp *resource.ModifyPlanResponse) {
	if req.Plan.Raw.IsNull() {
		return
	}
	// var plan models.HostResourceModel
	// diags := req.Plan.Get(ctx, &plan)
	// resp.Diagnostics.Append(diags...)

	// planHost := make([]models.HostResourceModel, 0)
	// diags.Append(plan.HostDetails.ElementsAs(ctx, &planHost, true)...)

	// for index, host := range planHost {
	// 	if !host.ID.IsNull() {
	// 		system, err := helper.GetFirstSystem(r.client)

	// 		if err != nil {
	// 			resp.Diagnostics.AddError(
	// 				"Error in getting system instance on the PowerFlex cluster",
	// 				err.Error(),
	// 			)
	// 			return
	// 		}

	// 		hostData, err := system.GetSdcByID(host.ID.ValueString())

	// 		if err != nil {
	// 			resp.Diagnostics.AddError(
	// 				"[Read] Unable to Find HOST by ID:"+host.ID.ValueString(),
	// 				err.Error(),
	// 			)
	// 			return
	// 		}
	// 		planHost[index].IP = types.StringValue(hostData.Sdc.SdcIP)
	// 	}
	// }

	// hostList, dgs := helper.GetSdcsValue(planHost)
	// diags.Append(dgs...)
	// plan.HostDetails = hostList

	// diags = resp.Plan.Set(ctx, &plan)
	// resp.Diagnostics.Append(diags...)
}

// Create - function to Create for Host Resource.
func (r *hostResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Create")

	var plan models.HostResourceModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	credential := []models.CredentialModel{}
	diags = plan.Credential.ElementsAs(ctx, &credential, true)
	resp.Diagnostics.Append(diags...)

	hostDetailList := []models.HostDetailModel{}
	diags = plan.HostDetails.ElementsAs(ctx, &hostDetailList, true)
	resp.Diagnostics.Append(diags...)

	system, err := helper.GetFirstSystem(r.client)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}

	//Checking that SDC exist or not
	sdcData, _ := system.FindSdc("SdcIP", plan.IP.ValueString())
	if sdcData != nil {
		resp.Diagnostics.AddError(
			"SDC Host is already Connected with PowerFlex cluster",
			"SDC Host is already Connected with PowerFlex cluster",
		)
		return
	}

	mdmIP := []string{}
	if !plan.MdmIPs.IsNull() && len(plan.MdmIPs.Elements()) > 0 {
		diags = plan.MdmIPs.ElementsAs(ctx, &mdmIP, true)
	} else {
		mdmDetails, err := system.GetMDMClusterDetails()

		mdmIP = GetMdmIPList(mdmDetails)
		if err != nil {
			resp.Diagnostics.AddError(
				"Error in getting MDM Details on the PowerFlex cluster",
				err.Error(),
			)
			return
		}
	}

	//var chnagedSDCs []models.HostResourceModel

	if plan.OSFamily.ValueString() == "windows" {
		hostDetailList, _ = helper.HostWindowsOperations(ctx, *r.client, plan, mdmIP, credential[0], system)
	} else if plan.OSFamily.ValueString() == "linux" {

	} else if plan.OSFamily.ValueString() == "esxi" {

	}

	if len(hostDetailList) > 0 {

		data, dgs := helper.UpdateHostState(hostDetailList, plan)
		resp.Diagnostics.Append(dgs...)

		diags = resp.State.Set(ctx, data)
		resp.Diagnostics.Append(diags...)

		tflog.Info(ctx, "Host Details updated to state file successfully")

		// return
	}

	resp.Diagnostics.AddError(
		"Please provide valid inputs",
		"Please provide valid inputs",
	)
}

func GetMdmIPList(mdmDetails *goscaleio_types.MdmCluster) []string {
	var ipmap []string

	for index := range mdmDetails.PrimaryMDM.IPs {
		ipmap = append(ipmap, mdmDetails.PrimaryMDM.IPs[index])
	}

	for _, mdm := range mdmDetails.SecondaryMDM {
		for index := range mdm.IPs {
			ipmap = append(ipmap, mdm.IPs[index])
		}
	}

	return ipmap
}

// Read - function to Read for Host Resource.
func (r *hostResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Read")
	// var state models.HostResourceModel
	// diags := req.State.Get(ctx, &state)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }

	// sdcStateList := []models.SDCStateDataModel{}
	// diags = state.HostDetails.ElementsAs(ctx, &sdcStateList, true)
	// resp.Diagnostics.Append(diags...)

	// system, err := helper.GetFirstSystem(r.client)
	// if err != nil {
	// 	resp.Diagnostics.AddError(
	// 		"Error in getting system instance on the PowerFlex cluster",
	// 		err.Error(),
	// 	)
	// 	return
	// }

	// var chnagedSDCs []models.SDCStateDataModel

	// //For handling the import case
	// if state.ID.ValueString() != "" && state.ID.ValueString() != "placeholder" {

	// 	for _, id := range strings.Split(state.ID.ValueString(), ",") {
	// 		sdcData, err := system.GetSdcByID(id)

	// 		if err != nil {
	// 			resp.Diagnostics.AddError(
	// 				"[Import] Unable to Find SDC by ID:"+id,
	// 				err.Error(),
	// 			)
	// 			return
	// 		}

	// 		if sdcData != nil {
	// 			changedSDCDetail := helper.GetSDCState(*sdcData.Sdc)

	// 			chnagedSDCs = append(chnagedSDCs, changedSDCDetail)
	// 		}
	// 	}
	// } else if len(sdcStateList) > 0 {

	// 	//For handling the multiple sdc_details update
	// 	for _, sdc := range sdcStateList {
	// 		var sdcData *goscaleio.Sdc

	// 		if sdc.SDCID.ValueString() != "" {
	// 			sdcData, err = system.GetSdcByID(sdc.SDCID.ValueString())

	// 			if err != nil {
	// 				resp.Diagnostics.AddError(
	// 					"[Read] Unable to Find SDC by ID:"+sdc.SDCID.ValueString(),
	// 					err.Error(),
	// 				)
	// 			}
	// 		} else if sdc.IP.ValueString() != "" {
	// 			sdcData, err = system.FindSdc("SdcIP", sdc.IP.ValueString())

	// 			if err != nil {
	// 				resp.Diagnostics.AddError(
	// 					"[Read] Unable to Find SDC by IP:"+sdc.IP.ValueString(),
	// 					err.Error(),
	// 				)
	// 			}
	// 		} else if sdc.SDCName.ValueString() != "" {
	// 			sdcData, err = system.FindSdc("Name", sdc.SDCName.ValueString())

	// 			if err != nil {
	// 				resp.Diagnostics.AddError(
	// 					"[Read] Unable to Find SDC by Name:"+sdc.SDCName.ValueString(),
	// 					err.Error(),
	// 				)
	// 			}
	// 		}

	// 		if sdcData != nil {
	// 			changedSDCDetail := helper.GetSDCState(*sdcData.Sdc)

	// 			chnagedSDCs = append(chnagedSDCs, changedSDCDetail)
	// 		}
	// 	}
	// } else if state.ID.ValueString() == "" {
	// 	resp.Diagnostics.AddError("[Read] Please provide valid Host ID", "Please provide valid Host ID")
	// 	return
	// }

	// data, dgs := helper.UpdateState(chnagedSDCs, state)
	// resp.Diagnostics.Append(dgs...)

	// diags = resp.State.Set(ctx, data)
	// resp.Diagnostics.Append(diags...)
	// if resp.Diagnostics.HasError() {
	// 	return
	// }
}

// Update - function to Update for Host Resource.
func (r *hostResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Update")
	// Retrieve values from plan

	tflog.Info(ctx, "SDC Details deleted from state file successfully")
}

// Delete - function to Delete for Host Resource.
func (r *hostResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Delete")
	var state models.HostResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	hostDetailList := []models.HostResourceModel{}
	diags = state.HostDetails.ElementsAs(ctx, &hostDetailList, true)
	resp.Diagnostics.Append(diags...)

	system, err := helper.GetFirstSystem(r.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}

	for _, host := range hostDetailList {

		err := system.DeleteSdc(host.ID.ValueString())

		if err != nil {
			resp.Diagnostics.AddError(
				"[Delete] Unable to Delete Host by ID:"+host.ID.ValueString(),
				err.Error(),
			)
			return
		}
	}

	resp.State.RemoveResource(ctx)

}

// ImportState - function to ImportState for Host Resource.
func (r *hostResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] ImportState :-- "+helper.PrettyJSON(req))
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// HostInstallationOperations function for the Host Installation Operation
func (r *hostResource) HostInstallationOperations(ctx context.Context, plan models.HostResourceModel, system *goscaleio.System, sdcDetails []models.HostDetailModel, sdcStateDetails []models.SDCStateDataModel) (dia diag.Diagnostics) {

	// if helper.CheckForExpansion(sdcDetails, sdcStateDetails) {
	// 	parsecsvRespose, parseCSVError := helper.ParseCSVOperation(ctx, sdcDetails, r.gatewayClient)

	// 	if parseCSVError != nil {
	// 		dia.AddError(
	// 			"Error while Parsing CSV",
	// 			"unexpected error: "+parseCSVError.Error(),
	// 		)
	// 		return
	// 	}

	// 	// to make gateway available for installation
	// 	queueOperationError := helper.ResetInstallerQueue(r.gatewayClient)
	// 	if queueOperationError != nil {
	// 		dia.AddError(
	// 			"Error Clearing Queue",
	// 			"unexpected error: "+queueOperationError.Error(),
	// 		)
	// 		return
	// 	}

	// 	tflog.Info(ctx, "Gateway Installer changed to idle phase before initiating process")

	// 	mdmIP, mdmIPError := helper.GetMDMIP(ctx, sdcDetails)
	// 	if mdmIPError != nil {
	// 		dia.AddError(
	// 			"Error while Getting MDM IP",
	// 			"unexpected error: "+mdmIPError.Error(),
	// 		)
	// 		return
	// 	}

	// 	tflog.Info(ctx, "CSV File parsed ssucessfully")

	// 	// Vaidate the MDM credentials
	// 	validateMDMResponse, validateMDMError := helper.ValidateMDMOperation(ctx, plan, r.gatewayClient, mdmIP)
	// 	if validateMDMError != nil {
	// 		dia.AddError(
	// 			"Error While Validating MDM Details",
	// 			"unexpected error: "+validateMDMResponse.Message,
	// 		)
	// 		return
	// 	}

	// 	if validateMDMResponse.StatusCode == 200 {

	// 		tflog.Info(ctx, "MDM Details validated successfully")

	// 		if !helper.CheckForNewSDCIPs(strings.Split(parsecsvRespose.Message, ","), strings.Split(validateMDMResponse.Data, ",")) {
	// 			installationError := helper.InstallationOperations(ctx, plan, r.gatewayClient, parsecsvRespose)

	// 			if installationError != nil {
	// 				dia.AddError(
	// 					"Error in Installation Process",
	// 					"unexpected error: "+installationError.Error(),
	// 				)
	// 				return
	// 			}
	// 		}
	// 	} else if validateMDMResponse.StatusCode != 200 {
	// 		dia.AddError(
	// 			"Error While Validating MDM Credentials",
	// 			"unexpected error: "+validateMDMResponse.Message+" & Status Code: "+strconv.Itoa(validateMDMResponse.StatusCode),
	// 		)
	// 		return
	// 	}
	// }

	return
}

// UpdateSDCNamdPerfProfileOperations function for Update Name and Performance Profile of SDC
// func (r *hostResource) UpdateSDCNamdPerfProfileOperations(ctx context.Context, sdcDetailList []models.HostDetailModel, system *goscaleio.System, chnagedSDCs *[]models.SDCStateDataModel) (dia diag.Diagnostics) {

// 	for _, sdc := range sdcDetailList {

// 		if sdc.HostName.ValueString() != "" && sdc.HostID.ValueString() != "" {

// 			nameExist, _ := helper.CheckForSDCName(system, sdc)

// 			if !nameExist {
// 				nameChng, err := system.ChangeSdcName(sdc.HostID.ValueString(), sdc.HostName.ValueString())

// 				if err != nil {
// 					dia.AddError(
// 						"[Create] Unable to Change Name Powerflex Host by ID:"+sdc.HostID.ValueString()+" Name:"+sdc.HostName.ValueString(),
// 						err.Error(),
// 					)
// 				}

// 				tflog.Debug(ctx, fmt.Sprintf("[POWERFLEX] Name Change Result: %s  Host ID: %s", helper.PrettyJSON(nameChng), sdc.HostID))
// 			}
// 		}

// 		if sdc.PerformanceProfile.ValueString() != "" && sdc.HostID.ValueString() != "" {
// 			perProfile, err := system.ChangeSdcPerfProfile(sdc.HostID.ValueString(), sdc.PerformanceProfile.ValueString())

// 			if err != nil {
// 				dia.AddError(
// 					"[Create] Unable to Change Performance Profile Powerflex Host by ID:"+sdc.HostID.ValueString(),
// 					err.Error(),
// 				)
// 			}

// 			tflog.Debug(ctx, fmt.Sprintf("[POWERFLEX] Performance Profile Change Result: %s  Host ID: %s", helper.PrettyJSON(perProfile), sdc.HostID))
// 		}

// 		finalSDC, err := system.GetSdcByID(sdc.HostID.ValueString())
// 		if err != nil {
// 			dia.AddError(
// 				"Unable to Read Changed Host",
// 				err.Error(),
// 			)
// 			return
// 		}

// 		if finalSDC != nil {
// 			changedSDCDetail := helper.GetSDCState(*finalSDC.Sdc)
// 			*chnagedSDCs = append(*chnagedSDCs, changedSDCDetail)
// 		}
// 	}

// 	return
// }
