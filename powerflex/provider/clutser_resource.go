package provider

import (
	"context"
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
	_ resource.Resource                = &clusterResource{}
	_ resource.ResourceWithConfigure   = &clusterResource{}
	_ resource.ResourceWithImportState = &clusterResource{}
)

// NewClusterResource is a helper function to simplify the provider implementation.
func NewClusterResource() resource.Resource {
	return &clusterResource{}
}

// clusterResource is the resource implementation.
type clusterResource struct {
	client        *goscaleio.Client
	gatewayClient *goscaleio.GatewayClient
}

func (r *clusterResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_cluster"
}

func (r *clusterResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = ClusterReourceSchema
}

func (d *clusterResource) ValidateConfig(ctx context.Context, req resource.ValidateConfigRequest, resp *resource.ValidateConfigResponse) {
	var config models.ClusterResourceModel

	resp.Diagnostics.Append(req.Config.Get(ctx, &config)...)

	if resp.Diagnostics.HasError() {
		return
	}

	clusterInstallationDetailsDataModel := []models.ClusterModel{}
	diags := config.Cluster.ElementsAs(ctx, &clusterInstallationDetailsDataModel, true)
	resp.Diagnostics.Append(diags...)

	storagePoolDetailsDataModel := []models.StoragePoolDataModel{}
	diags = config.StoragePools.ElementsAs(ctx, &storagePoolDetailsDataModel, true)
	resp.Diagnostics.Append(diags...)

	sdrValidation := false

	if !config.Cluster.IsNull() && !config.StoragePools.IsNull() {

		sdrCheck := false

		for _, row := range clusterInstallationDetailsDataModel {
			if strings.EqualFold(row.IsSdr.ValueString(), "Yes") {
				sdrCheck = true
			}
		}

		if sdrCheck {
			if len(storagePoolDetailsDataModel) > 0 {
				for _, row := range storagePoolDetailsDataModel {
					if !row.ReplicationJournalCapacityPercentage.IsNull() {
						sdrValidation = true
					}
				}
			}
		} else {
			sdrValidation = true
		}
	}

	if !sdrValidation {
		resp.Diagnostics.AddAttributeError(
			path.Root("replication_journal_capacity_percentage"),
			"Please configure replication_journal_capacity_percentage for SDR.",
			"Please configure replication_journal_capacity_percentage for SDR.",
		)
	}
}

// Configure - function to return Configuration for SDC resource.
func (r *clusterResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	if req.ProviderData == nil {
		return
	}

	if req.ProviderData.(*powerflexProvider).client != nil {

		r.client = req.ProviderData.(*powerflexProvider).client
	}

	if req.ProviderData.(*powerflexProvider).gatewayClient != nil {

		r.gatewayClient = req.ProviderData.(*powerflexProvider).gatewayClient
	} else {
		resp.Diagnostics.AddError("Unable to Authenticate Goscaleio API Client", req.ProviderData.(*powerflexProvider).clientError)

		return
	}
}

// Create - function to Create for SDC resource.
func (r *clusterResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] Create")

	var plan models.ClusterResourceModel

	diags := req.Plan.Get(ctx, &plan)

	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	clusterInstallationDetailsDataModel := []models.ClusterModel{}
	diags = plan.Cluster.ElementsAs(ctx, &clusterInstallationDetailsDataModel, true)
	resp.Diagnostics.Append(diags...)

	storagePoolDetailsDataModel := []models.StoragePoolDataModel{}
	diags = plan.StoragePools.ElementsAs(ctx, &storagePoolDetailsDataModel, true)
	resp.Diagnostics.Append(diags...)

	if len(clusterInstallationDetailsDataModel) > 0 && len(storagePoolDetailsDataModel) > 0 {
		data, dgs := r.ClusterDeploymentOperations(ctx, plan, clusterInstallationDetailsDataModel, storagePoolDetailsDataModel)
		resp.Diagnostics.Append(dgs...)
		if resp.Diagnostics.HasError() {
			return
		}

		diags = resp.State.Set(ctx, data)
		resp.Diagnostics.Append(diags...)

		tflog.Info(ctx, "Cluster Details updated to state file successfully")

		return
	}

	resp.Diagnostics.AddError("[Create] Please provide valid Clustet and Storage Pool Details", "Please provide valid valid Clustet and Storage Pool Detail Details")

	return

}

// Read refreshes the Terraform state with the latest data.
func (r *clusterResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {

	tflog.Debug(ctx, "[POWERFLEX] Read")
	var state models.ClusterResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// to make gateway available for installation
	queueOperationError := helper.ResetInstallerQueue(r.gatewayClient)
	if queueOperationError != nil {
		resp.Diagnostics.AddError(
			"Error Clearing Queue",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}

	//For handling the import case
	if state.ID.ValueString() != "" && state.ID.ValueString() != "placeholder" {

		inputData := strings.Split(state.ID.ValueString(), ",")

		if len(inputData) == 3 {
			mdmIP := inputData[0]

			state.MdmPassword = types.StringValue(inputData[1])

			state.LiaPassword = types.StringValue(inputData[2])

			state.AllowNonSecureCommunicationWithMdm = types.BoolValue(true)

			state.AllowNonSecureCommunicationWithLia = types.BoolValue(true)

			data, diags := helper.UpdateClusterState(state, r.gatewayClient, mdmIP)
			resp.Diagnostics.Append(diags...)

			diags = resp.State.Set(ctx, data)
			resp.Diagnostics.Append(diags...)
			if resp.Diagnostics.HasError() {
				return
			}
		} else {
			resp.Diagnostics.AddError("[Read] Please provide valid Input Details", "Please provide valid Input Details")

			return
		}
	} else {

		mdmList := []models.MDMModel{}
		diags = state.MDMList.ElementsAs(ctx, &mdmList, true)
		resp.Diagnostics.Append(diags...)

		mdmIP, err := helper.GetMDMIPFromMDMList(mdmList)
		if err != nil {
			diags.AddError(
				"Error in Fecthing Primary MDM IP",
				"unexpected error: "+err.Error(),
			)
			return
		}

		data, diags := helper.UpdateClusterState(state, r.gatewayClient, mdmIP)
		resp.Diagnostics.Append(diags...)

		diags = resp.State.Set(ctx, data)
		resp.Diagnostics.Append(diags...)
		if resp.Diagnostics.HasError() {
			return
		}
	}

}

// ImportState - function to ImportState for SDC resource.
func (r *clusterResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	tflog.Debug(ctx, "[POWERFLEX] ImportState :-- "+helper.PrettyJSON(req))
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

// Update updates the resource and sets the updated Terraform state on success.
func (r *clusterResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {

	resp.Diagnostics.AddError("[Update] Update operation is not available.", "Update operation is not available.")

	return
}

// Delete deletes the resource and removes the Terraform state on success.
func (r *clusterResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from state
	var state models.ClusterResourceModel
	diags := req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	mdmList := []models.MDMModel{}
	diags = state.MDMList.ElementsAs(ctx, &mdmList, true)
	resp.Diagnostics.Append(diags...)

	// to make gateway available for installation
	queueOperationError := helper.ResetInstallerQueue(r.gatewayClient)
	if queueOperationError != nil {
		resp.Diagnostics.AddError(
			"Error Clearing Queue",
			"unexpected error: "+queueOperationError.Error(),
		)
		return
	}

	mdmIP, err := helper.GetMDMIPFromMDMList(mdmList)
	if err != nil {
		resp.Diagnostics.AddError(
			"Error in Fecthing Primary MDM IP",
			"unexpected error: "+err.Error(),
		)
		return
	}

	system, err := helper.GetFirstSystem(r.client)

	if err != nil {
		resp.Diagnostics.AddError(
			"Error in getting system instance on the PowerFlex cluster",
			err.Error(),
		)
		return
	}

	sdcs, err := system.GetSdc()

	for _, sdc := range sdcs {
		if sdc.MdmConnectionState == "Disconnected" {
			err := system.DeleteSdc(sdc.ID)

			if err != nil {
				resp.Diagnostics.AddError(
					"[Delete] Unable to Delete SDC by ID:"+sdc.ID,
					err.Error(),
				)
				return
			}
		}
	}

	clusteDetailResponse, error := helper.GetClusterDetails(state, r.gatewayClient, mdmIP, true)
	if error != nil {
		resp.Diagnostics.AddError(
			"Error in validating MDM IP",
			"unexpected error: "+error.Error(),
		)
		return
	}

	installationError := helper.ClusterUninstallationOperations(ctx, state, r.gatewayClient, clusteDetailResponse)
	if installationError != nil {
		resp.Diagnostics.AddError(
			"Error in Uninstallation Process",
			"unexpected error: "+installationError.Error(),
		)
	}

	if resp.Diagnostics.HasError() {
		return
	}

	resp.State.RemoveResource(ctx)
}

// ClusterDeploymentOperations function for the Cluster Deployment Operation Like ParseCSV,Installation and Validate Cluster
func (r *clusterResource) ClusterDeploymentOperations(ctx context.Context, plan models.ClusterResourceModel, clusterInstallationDetailsDataModel []models.ClusterModel, storagePoolDetailsDataModel []models.StoragePoolDataModel) (data models.ClusterResourceModel, dia diag.Diagnostics) {

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

	mdmIP, err := helper.GetMDMIPFromClusterDetails(clusterInstallationDetailsDataModel)
	if err != nil {
		dia.AddError(
			"Error in Fecthing Primary MDM IP Before Installation",
			"unexpected error: "+err.Error(),
		)
		return
	}

	_, clusterError := helper.GetClusterDetails(plan, r.gatewayClient, mdmIP, false)
	if clusterError == nil {
		dia.AddError(
			"Cluster Already Deployed",
			"Cluster already deployed for given inputs",
		)
		return
	} //TODO check else case

	parsecsvRespose, parseCSVError := helper.ParseClusterCSVOperation(ctx, r.gatewayClient, clusterInstallationDetailsDataModel, storagePoolDetailsDataModel)

	if parseCSVError != nil {
		dia.AddError(
			"Error while Parsing CSV",
			"unexpected error: "+parseCSVError.Error(),
		)
		return
	}

	installationError := helper.ClusterInstallationOperations(ctx, plan, r.gatewayClient, parsecsvRespose)

	if installationError != nil {
		dia.AddError(
			"Error in Installation Process",
			"unexpected error: "+installationError.Error(),
		)
		return
	}

	data, dia = helper.UpdateClusterState(plan, r.gatewayClient, mdmIP)

	return data, dia
}
