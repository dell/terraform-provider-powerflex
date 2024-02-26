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

package helper

import (
	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetTemplateState converts scaleiotypes.Template to models.Template
func GetTemplateState(input scaleiotypes.TemplateDetails) models.TemplateModel {
	return models.TemplateModel{
		ID:                     types.StringValue(input.ID),
		TemplateName:           types.StringValue(input.TemplateName),
		TemplateDescription:    types.StringValue(input.TemplateDescription),
		TemplateType:           types.StringValue(input.TemplateType),
		TemplateVersion:        types.StringValue(input.TemplateVersion),
		OriginalTemplateID:     types.StringValue(input.OriginalTemplateID),
		TemplateValid:          GetTemplateValid(input.TemplateValid),
		TemplateLocked:         types.BoolValue(input.TemplateLocked),
		InConfiguration:        types.BoolValue(input.InConfiguration),
		CreatedDate:            types.StringValue(input.CreatedDate),
		CreatedBy:              types.StringValue(input.CreatedBy),
		UpdatedDate:            types.StringValue(input.UpdatedDate),
		LastDeployedDate:       types.StringValue(input.LastDeployedDate),
		UpdatedBy:              types.StringValue(input.UpdatedBy),
		ManageFirmware:         types.BoolValue(input.ManageFirmware),
		UseDefaultCatalog:      types.BoolValue(input.UseDefaultCatalog),
		FirmwareRepository:     GetFirmwareRepository(input.FirmwareRepository),
		LicenseRepository:      GetLicenseRepository(input.LicenseRepository),
		AssignedUsers:          GetAssignedUsersList(input.AssignedUsers),
		AllUsersAllowed:        types.BoolValue(input.AllUsersAllowed),
		Category:               types.StringValue(input.Category),
		Components:             GetComponentsList(input.Components),
		Configuration:          GetConfigurationDetails(input.Configuration),
		ServerCount:            types.Int64Value(int64(input.ServerCount)),
		StorageCount:           types.Int64Value(int64(input.StorageCount)),
		ClusterCount:           types.Int64Value(int64(input.ClusterCount)),
		ServiceCount:           types.Int64Value(int64(input.ServiceCount)),
		SwitchCount:            types.Int64Value(int64(input.SwitchCount)),
		VMCount:                types.Int64Value(int64(input.VMCount)),
		SdnasCount:             types.Int64Value(int64(input.SdnasCount)),
		BrownfieldTemplateType: types.StringValue(input.BrownfieldTemplateType),
		Networks:               GetNetworksList(input.Networks),
		Draft:                  types.BoolValue(input.Draft),
	}
}

// GetAssignedUsersList converts list of scaleiotypes.AssignedUsers to list of models.AssignedUsers
func GetAssignedUsersList(inputs []scaleiotypes.AssignedUsers) []models.AssignedUsers {
	out := make([]models.AssignedUsers, 0)
	for _, input := range inputs {
		out = append(out, GetAssignedUsers(input))
	}
	return out
}

// GetComponentsList converts list of scaleiotypes.Components to list of models.Components
func GetComponentsList(inputs []scaleiotypes.Components) []models.Components {
	out := make([]models.Components, 0)
	for _, input := range inputs {
		out = append(out, GetComponents(input))
	}
	return out
}

// GetNetworksList converts list of scaleiotypes.Networks to list of models.Networks
func GetNetworksList(inputs []scaleiotypes.Networks) []models.Networks {
	out := make([]models.Networks, 0)
	for _, input := range inputs {
		out = append(out, GetNetworks(input))
	}
	return out
}

// GetMessages converts scaleiotypes.Messages to models.Messages
func GetMessages(input scaleiotypes.Messages) models.Messages {
	return models.Messages{
		ID:              types.StringValue(input.ID),
		MessageCode:     types.StringValue(input.MessageCode),
		MessageBundle:   types.StringValue(input.MessageBundle),
		Severity:        types.StringValue(input.Severity),
		Category:        types.StringValue(input.Category),
		DisplayMessage:  types.StringValue(input.DisplayMessage),
		ResponseAction:  types.StringValue(input.ResponseAction),
		DetailedMessage: types.StringValue(input.DetailedMessage),
		CorrelationID:   types.StringValue(input.CorrelationID),
		AgentID:         types.StringValue(input.AgentID),
		TimeStamp:       types.StringValue(input.TimeStamp),
		SequenceNumber:  types.Int64Value(int64(input.SequenceNumber)),
	}
}

// GetTemplateValid converts scaleiotypes.TemplateValid to models.TemplateValid
func GetTemplateValid(input scaleiotypes.TemplateValid) models.TemplateValid {
	return models.TemplateValid{
		Valid:    types.BoolValue(input.Valid),
		Messages: GetMessagesList(input.Messages),
	}
}

// GetMessagesList converts list of scaleiotypes.Messages to list of models.Messages
func GetMessagesList(inputs []scaleiotypes.Messages) []models.Messages {
	out := make([]models.Messages, 0)
	for _, input := range inputs {
		out = append(out, GetMessages(input))
	}
	return out
}

// GetSoftwareComponents converts scaleiotypes.SoftwareComponents to models.SoftwareComponents
func GetSoftwareComponents(input scaleiotypes.SoftwareComponents) models.SoftwareComponents {
	return models.SoftwareComponents{
		ID:                  types.StringValue(input.ID),
		PackageID:           types.StringValue(input.PackageID),
		DellVersion:         types.StringValue(input.DellVersion),
		VendorVersion:       types.StringValue(input.VendorVersion),
		ComponentID:         types.StringValue(input.ComponentID),
		DeviceID:            types.StringValue(input.DeviceID),
		SubDeviceID:         types.StringValue(input.SubDeviceID),
		VendorID:            types.StringValue(input.VendorID),
		SubVendorID:         types.StringValue(input.SubVendorID),
		CreatedDate:         types.StringValue(input.CreatedDate),
		CreatedBy:           types.StringValue(input.CreatedBy),
		UpdatedDate:         types.StringValue(input.UpdatedDate),
		UpdatedBy:           types.StringValue(input.UpdatedBy),
		Path:                types.StringValue(input.Path),
		HashMd5:             types.StringValue(input.HashMd5),
		Name:                types.StringValue(input.Name),
		Category:            types.StringValue(input.Category),
		ComponentType:       types.StringValue(input.ComponentType),
		OperatingSystem:     types.StringValue(input.OperatingSystem),
		SystemIDs:           GeStringList(input.SystemIDs),
		Custom:              types.BoolValue(input.Custom),
		NeedsAttention:      types.BoolValue(input.NeedsAttention),
		Ignore:              types.BoolValue(input.Ignore),
		OriginalVersion:     types.StringValue(input.OriginalVersion),
		OriginalComponentID: types.StringValue(input.OriginalComponentID),
		FirmwareRepoName:    types.StringValue(input.FirmwareRepoName),
	}
}

// GeStringList converts list of string to list of String
func GeStringList(inputs []string) []types.String {
	out := make([]types.String, 0)
	for _, input := range inputs {
		out = append(out, types.StringValue(input))
	}
	return out
}

// GetSoftwareBundles converts scaleiotypes.SoftwareBundles to models.SoftwareBundles
func GetSoftwareBundles(input scaleiotypes.SoftwareBundles) models.SoftwareBundles {
	return models.SoftwareBundles{
		ID:                 types.StringValue(input.ID),
		Name:               types.StringValue(input.Name),
		Version:            types.StringValue(input.Version),
		BundleDate:         types.StringValue(input.BundleDate),
		CreatedDate:        types.StringValue(input.CreatedDate),
		CreatedBy:          types.StringValue(input.CreatedBy),
		UpdatedDate:        types.StringValue(input.UpdatedDate),
		UpdatedBy:          types.StringValue(input.UpdatedBy),
		Description:        types.StringValue(input.Description),
		UserBundle:         types.BoolValue(input.UserBundle),
		UserBundlePath:     types.StringValue(input.UserBundlePath),
		UserBundleHashMd5:  types.StringValue(input.UserBundleHashMd5),
		DeviceType:         types.StringValue(input.DeviceType),
		DeviceModel:        types.StringValue(input.DeviceModel),
		Criticality:        types.StringValue(input.Criticality),
		FwRepositoryID:     types.StringValue(input.FwRepositoryID),
		BundleType:         types.StringValue(input.BundleType),
		Custom:             types.BoolValue(input.Custom),
		NeedsAttention:     types.BoolValue(input.NeedsAttention),
		SoftwareComponents: GetSoftwareComponentsList(input.SoftwareComponents),
	}
}

// GetSoftwareComponentsList converts list of scaleiotypes.SoftwareComponents to list of models.SoftwareComponents
func GetSoftwareComponentsList(inputs []scaleiotypes.SoftwareComponents) []models.SoftwareComponents {
	out := make([]models.SoftwareComponents, 0)
	for _, input := range inputs {
		out = append(out, GetSoftwareComponents(input))
	}
	return out
}

// GetDeploymentValid converts scaleiotypes.DeploymentValid to models.DeploymentValid
func GetDeploymentValid(input scaleiotypes.DeploymentValid) models.DeploymentValid {
	return models.DeploymentValid{
		Valid:    types.BoolValue(input.Valid),
		Messages: GetMessagesList(input.Messages),
	}
}

// GetDeploymentDevice converts scaleiotypes.DeploymentDevice to models.DeploymentDevice
func GetDeploymentDevice(input scaleiotypes.DeploymentDevice) models.DeploymentDevice {
	return models.DeploymentDevice{
		RefID:            types.StringValue(input.RefID),
		RefType:          types.StringValue(input.RefType),
		LogDump:          types.StringValue(input.LogDump),
		Status:           types.StringValue(input.Status),
		StatusEndTime:    types.StringValue(input.StatusEndTime),
		StatusStartTime:  types.StringValue(input.StatusStartTime),
		DeviceHealth:     types.StringValue(input.DeviceHealth),
		HealthMessage:    types.StringValue(input.HealthMessage),
		CompliantState:   types.StringValue(input.CompliantState),
		BrownfieldStatus: types.StringValue(input.BrownfieldStatus),
		DeviceType:       types.StringValue(input.DeviceType),
		DeviceGroupName:  types.StringValue(input.DeviceGroupName),
		IPAddress:        types.StringValue(input.IPAddress),
		CurrentIPAddress: types.StringValue(input.CurrentIPAddress),
		ServiceTag:       types.StringValue(input.ServiceTag),
		ComponentID:      types.StringValue(input.ComponentID),
		StatusMessage:    types.StringValue(input.StatusMessage),
		Model:            types.StringValue(input.Model),
		CloudLink:        types.BoolValue(input.CloudLink),
		DasCache:         types.BoolValue(input.DasCache),
		DeviceState:      types.StringValue(input.DeviceState),
		PuppetCertName:   types.StringValue(input.PuppetCertName),
		Brownfield:       types.BoolValue(input.Brownfield),
	}
}

// GetVms converts scaleiotypes.Vms to models.Vms
func GetVms(input scaleiotypes.Vms) models.Vms {
	return models.Vms{
		CertificateName: types.StringValue(input.CertificateName),
		VMModel:         types.StringValue(input.VMModel),
		VMIpaddress:     types.StringValue(input.VMIpaddress),
		VMManufacturer:  types.StringValue(input.VMManufacturer),
		VMServiceTag:    types.StringValue(input.VMServiceTag),
	}
}

// GetLicenseRepository converts scaleiotypes.LicenseRepository to models.LicenseRepository
func GetLicenseRepository(input scaleiotypes.LicenseRepository) models.LicenseRepository {
	return models.LicenseRepository{
		ID:           types.StringValue(input.ID),
		Name:         types.StringValue(input.Name),
		Type:         types.StringValue(input.Type),
		DiskLocation: types.StringValue(input.DiskLocation),
		Filename:     types.StringValue(input.Filename),
		State:        types.StringValue(input.State),
		CreatedDate:  types.StringValue(input.CreatedDate),
		CreatedBy:    types.StringValue(input.CreatedBy),
		UpdatedDate:  types.StringValue(input.UpdatedDate),
		UpdatedBy:    types.StringValue(input.UpdatedBy),
		LicenseData:  types.StringValue(input.LicenseData),
	}
}

// GetAssignedUsers converts scaleiotypes.AssignedUsers to models.AssignedUsers
func GetAssignedUsers(input scaleiotypes.AssignedUsers) models.AssignedUsers {
	return models.AssignedUsers{
		UserSeqID:      types.Int64Value(int64(input.UserSeqID)),
		UserName:       types.StringValue(input.UserName),
		Password:       types.StringValue(input.Password),
		UpdatePassword: types.BoolValue(input.UpdatePassword),
		DomainName:     types.StringValue(input.DomainName),
		GroupDN:        types.StringValue(input.GroupDN),
		GroupName:      types.StringValue(input.GroupName),
		FirstName:      types.StringValue(input.FirstName),
		LastName:       types.StringValue(input.LastName),
		Email:          types.StringValue(input.Email),
		PhoneNumber:    types.StringValue(input.PhoneNumber),
		Enabled:        types.BoolValue(input.Enabled),
		SystemUser:     types.BoolValue(input.SystemUser),
		CreatedDate:    types.StringValue(input.CreatedDate),
		CreatedBy:      types.StringValue(input.CreatedBy),
		UpdatedDate:    types.StringValue(input.UpdatedDate),
		UpdatedBy:      types.StringValue(input.UpdatedBy),
		Role:           types.StringValue(input.Role),
		UserPreference: types.StringValue(input.UserPreference),
		ID:             types.StringValue(input.ID),
		Roles:          GeStringList(input.Roles),
	}
}

// GetJobDetails converts scaleiotypes.JobDetails to models.JobDetails
func GetJobDetails(input scaleiotypes.JobDetails) models.JobDetails {
	return models.JobDetails{
		Level:       types.StringValue(input.Level),
		Message:     types.StringValue(input.Message),
		Timestamp:   types.StringValue(input.Timestamp),
		ExecutionID: types.StringValue(input.ExecutionID),
		ComponentID: types.StringValue(input.ComponentID),
	}
}

// GetDeploymentValidationResponse converts scaleiotypes.DeploymentValidationResponse to models.DeploymentValidationResponse
func GetDeploymentValidationResponse(input scaleiotypes.DeploymentValidationResponse) models.DeploymentValidationResponse {
	return models.DeploymentValidationResponse{
		Nodes:                  types.Int64Value(int64(input.Nodes)),
		StoragePools:           types.Int64Value(int64(input.StoragePools)),
		DrivesPerStoragePool:   types.Int64Value(int64(input.DrivesPerStoragePool)),
		MaxScalability:         types.Int64Value(int64(input.MaxScalability)),
		VirtualMachines:        types.Int64Value(int64(input.VirtualMachines)),
		NumberOfServiceVolumes: types.Int64Value(int64(input.NumberOfServiceVolumes)),
		CanDeploy:              types.BoolValue(input.CanDeploy),
		WarningMessages:        GeStringList(input.WarningMessages),
		StoragePoolDiskType:    GeStringList(input.StoragePoolDiskType),
		Hostnames:              GeStringList(input.Hostnames),
		NewNodeDiskTypes:       GeStringList(input.NewNodeDiskTypes),
		NoOfFaultSets:          types.Int64Value(int64(input.NoOfFaultSets)),
		NodesPerFaultSet:       types.Int64Value(int64(input.NodesPerFaultSet)),
		ProtectionDomain:       types.StringValue(input.ProtectionDomain),
		DiskTypeMismatch:       types.BoolValue(input.DiskTypeMismatch),
	}
}

// GetDeployments converts scaleiotypes.Deployments to models.Deployments
func GetDeployments(input scaleiotypes.Deployments) models.Deployments {
	return models.Deployments{
		ID:                           types.StringValue(input.ID),
		DeploymentName:               types.StringValue(input.DeploymentName),
		DeploymentDescription:        types.StringValue(input.DeploymentDescription),
		DeploymentValid:              GetDeploymentValid(input.DeploymentValid),
		Retry:                        types.BoolValue(input.Retry),
		Teardown:                     types.BoolValue(input.Teardown),
		TeardownAfterCancel:          types.BoolValue(input.TeardownAfterCancel),
		RemoveService:                types.BoolValue(input.RemoveService),
		CreatedDate:                  types.StringValue(input.CreatedDate),
		CreatedBy:                    types.StringValue(input.CreatedBy),
		UpdatedDate:                  types.StringValue(input.UpdatedDate),
		UpdatedBy:                    types.StringValue(input.UpdatedBy),
		DeploymentScheduledDate:      types.StringValue(input.DeploymentScheduledDate),
		DeploymentStartedDate:        types.StringValue(input.DeploymentStartedDate),
		DeploymentFinishedDate:       types.StringValue(input.DeploymentFinishedDate),
		ScheduleDate:                 types.StringValue(input.ScheduleDate),
		Status:                       types.StringValue(input.Status),
		Compliant:                    types.BoolValue(input.Compliant),
		DeploymentDevice:             GetDeploymentDeviceList(input.DeploymentDevice),
		Vms:                          GetVmsList(input.Vms),
		UpdateServerFirmware:         types.BoolValue(input.UpdateServerFirmware),
		UseDefaultCatalog:            types.BoolValue(input.UseDefaultCatalog),
		FirmwareRepositoryID:         types.StringValue(input.FirmwareRepositoryID),
		LicenseRepository:            GetLicenseRepository(input.LicenseRepository),
		LicenseRepositoryID:          types.StringValue(input.LicenseRepositoryID),
		IndividualTeardown:           types.BoolValue(input.IndividualTeardown),
		DeploymentHealthStatusType:   types.StringValue(input.DeploymentHealthStatusType),
		AssignedUsers:                GetAssignedUsersList(input.AssignedUsers),
		AllUsersAllowed:              types.BoolValue(input.AllUsersAllowed),
		Owner:                        types.StringValue(input.Owner),
		NoOp:                         types.BoolValue(input.NoOp),
		FirmwareInit:                 types.BoolValue(input.FirmwareInit),
		DisruptiveFirmware:           types.BoolValue(input.DisruptiveFirmware),
		PreconfigureSVM:              types.BoolValue(input.PreconfigureSVM),
		PreconfigureSVMAndUpdate:     types.BoolValue(input.PreconfigureSVMAndUpdate),
		ServicesDeployed:             types.StringValue(input.ServicesDeployed),
		PrecalculatedDeviceHealth:    types.StringValue(input.PrecalculatedDeviceHealth),
		LifecycleModeReasons:         GeStringList(input.LifecycleModeReasons),
		JobDetails:                   GetJobDetailsList(input.JobDetails),
		NumberOfDeployments:          types.Int64Value(int64(input.NumberOfDeployments)),
		OperationType:                types.StringValue(input.OperationType),
		OperationStatus:              types.StringValue(input.OperationStatus),
		OperationData:                types.StringValue(input.OperationData),
		DeploymentValidationResponse: GetDeploymentValidationResponse(input.DeploymentValidationResponse),
		CurrentStepCount:             types.StringValue(input.CurrentStepCount),
		TotalNumOfSteps:              types.StringValue(input.TotalNumOfSteps),
		CurrentStepMessage:           types.StringValue(input.CurrentStepMessage),
		CustomImage:                  types.StringValue(input.CustomImage),
		OriginalDeploymentID:         types.StringValue(input.OriginalDeploymentID),
		CurrentBatchCount:            types.StringValue(input.CurrentBatchCount),
		TotalBatchCount:              types.StringValue(input.TotalBatchCount),
		Brownfield:                   types.BoolValue(input.Brownfield),
		ScaleUp:                      types.BoolValue(input.ScaleUp),
		LifecycleMode:                types.BoolValue(input.LifecycleMode),
		OverallDeviceHealth:          types.StringValue(input.OverallDeviceHealth),
		Vds:                          types.BoolValue(input.Vds),
		TemplateValid:                types.BoolValue(input.TemplateValid),
		ConfigurationChange:          types.BoolValue(input.ConfigurationChange),
		CanMigratevCLSVMs:            types.BoolValue(input.CanMigratevCLSVMs),
	}
}

// GetDeploymentDeviceList converts list of scaleiotypes.DeploymentDevice to list of models.DeploymentDevice
func GetDeploymentDeviceList(inputs []scaleiotypes.DeploymentDevice) []models.DeploymentDevice {
	out := make([]models.DeploymentDevice, 0)
	for _, input := range inputs {
		out = append(out, GetDeploymentDevice(input))
	}
	return out
}

// GetVmsList converts list of scaleiotypes.Vms to list of models.Vms
func GetVmsList(inputs []scaleiotypes.Vms) []models.Vms {
	out := make([]models.Vms, 0)
	for _, input := range inputs {
		out = append(out, GetVms(input))
	}
	return out
}

// GetJobDetailsList converts list of scaleiotypes.JobDetails to list of models.JobDetails
func GetJobDetailsList(inputs []scaleiotypes.JobDetails) []models.JobDetails {
	out := make([]models.JobDetails, 0)
	for _, input := range inputs {
		out = append(out, GetJobDetails(input))
	}
	return out
}

// GetFirmwareRepository converts scaleiotypes.FirmwareRepository to models.FirmwareRepository
func GetFirmwareRepository(input scaleiotypes.FirmwareRepository) models.FirmwareRepository {
	return models.FirmwareRepository{
		ID:                      types.StringValue(input.ID),
		Name:                    types.StringValue(input.Name),
		SourceLocation:          types.StringValue(input.SourceLocation),
		SourceType:              types.StringValue(input.SourceType),
		DiskLocation:            types.StringValue(input.DiskLocation),
		Filename:                types.StringValue(input.Filename),
		Md5Hash:                 types.StringValue(input.Md5Hash),
		Username:                types.StringValue(input.Username),
		Password:                types.StringValue(input.Password),
		DownloadStatus:          types.StringValue(input.DownloadStatus),
		CreatedDate:             types.StringValue(input.CreatedDate),
		CreatedBy:               types.StringValue(input.CreatedBy),
		UpdatedDate:             types.StringValue(input.UpdatedDate),
		UpdatedBy:               types.StringValue(input.UpdatedBy),
		DefaultCatalog:          types.BoolValue(input.DefaultCatalog),
		Embedded:                types.BoolValue(input.Embedded),
		State:                   types.StringValue(input.State),
		SoftwareComponents:      GetSoftwareComponentsList(input.SoftwareComponents),
		SoftwareBundles:         GetSoftwareBundlesList(input.SoftwareBundles),
		Deployments:             GetDeploymentsList(input.Deployments),
		BundleCount:             types.Int64Value(int64(input.BundleCount)),
		ComponentCount:          types.Int64Value(int64(input.ComponentCount)),
		UserBundleCount:         types.Int64Value(int64(input.UserBundleCount)),
		Minimal:                 types.BoolValue(input.Minimal),
		DownloadProgress:        types.Int64Value(int64(input.DownloadProgress)),
		ExtractProgress:         types.Int64Value(int64(input.ExtractProgress)),
		FileSizeInGigabytes:     types.Int64Value(int64(input.FileSizeInGigabytes)),
		SignedKeySourceLocation: types.StringValue(input.SignedKeySourceLocation),
		Signature:               types.StringValue(input.Signature),
		Custom:                  types.BoolValue(input.Custom),
		NeedsAttention:          types.BoolValue(input.NeedsAttention),
		JobID:                   types.StringValue(input.JobID),
		Rcmapproved:             types.BoolValue(input.Rcmapproved),
	}
}

// GetSoftwareBundlesList converts list of scaleiotypes.SoftwareBundles to list of models.SoftwareBundles
func GetSoftwareBundlesList(inputs []scaleiotypes.SoftwareBundles) []models.SoftwareBundles {
	out := make([]models.SoftwareBundles, 0)
	for _, input := range inputs {
		out = append(out, GetSoftwareBundles(input))
	}
	return out
}

// GetDeploymentsList converts list of scaleiotypes.Deployments to list of models.Deployments
func GetDeploymentsList(inputs []scaleiotypes.Deployments) []models.Deployments {
	out := make([]models.Deployments, 0)
	for _, input := range inputs {
		out = append(out, GetDeployments(input))
	}
	return out
}

// GetComponentValid converts scaleiotypes.ComponentValid to models.ComponentValid
func GetComponentValid(input scaleiotypes.ComponentValid) models.ComponentValid {
	return models.ComponentValid{
		Valid:    types.BoolValue(input.Valid),
		Messages: GetMessagesList(input.Messages),
	}
}

// GetRelatedComponents converts scaleiotypes.RelatedComponents to models.RelatedComponents
func GetRelatedComponents(input map[string]string) basetypes.MapValue {
	elements := make(map[string]attr.Value) 
	for key, val := range input {
		elements[key] = types.StringValue(val)
	}
	var setRelComp basetypes.MapValue
		setRelComp , _ = types.MapValue(types.StringType, elements) 

	 return setRelComp	
	}

// GetDependenciesDetails converts scaleiotypes.DependenciesDetails to models.DependenciesDetails
func GetDependenciesDetails(input scaleiotypes.DependenciesDetails) models.DependenciesDetails {
	return models.DependenciesDetails{
		ID:               types.StringValue(input.ID),
		DependencyTarget: types.StringValue(input.DependencyTarget),
		DependencyValue:  types.StringValue(input.DependencyValue),
	}
}

// GetNetworkIPAddressList converts scaleiotypes.NetworkIPAddressList to models.NetworkIPAddressList
func GetNetworkIPAddressList(input scaleiotypes.NetworkIPAddressList) models.NetworkIPAddressList {
	return models.NetworkIPAddressList{
		ID:        types.StringValue(input.ID),
		IPAddress: types.StringValue(input.IPAddress),
	}
}

// GetPartitions converts scaleiotypes.Partitions to models.Partitions
func GetPartitions(input scaleiotypes.Partitions) models.Partitions {
	return models.Partitions{
		ID:                   types.StringValue(input.ID),
		Name:                 types.StringValue(input.Name),
		Networks:             GeStringList(input.Networks),
		NetworkIPAddressList: GetNetworkIPAddressListList(input.NetworkIPAddressList),
		Minimum:              types.Int64Value(int64(input.Minimum)),
		Maximum:              types.Int64Value(int64(input.Maximum)),
		LanMacAddress:        types.StringValue(input.LanMacAddress),
		IscsiMacAddress:      types.StringValue(input.IscsiMacAddress),
		IscsiIQN:             types.StringValue(input.IscsiIQN),
		Wwnn:                 types.StringValue(input.Wwnn),
		Wwpn:                 types.StringValue(input.Wwpn),
		Fqdd:                 types.StringValue(input.Fqdd),
		MirroredPort:         types.StringValue(input.MirroredPort),
		MacAddress:           types.StringValue(input.MacAddress),
		PortNo:               types.Int64Value(int64(input.PortNo)),
		PartitionNo:          types.Int64Value(int64(input.PartitionNo)),
		PartitionIndex:       types.Int64Value(int64(input.PartitionIndex)),
	}
}

// GetNetworkIPAddressListList converts list of scaleiotypes.NetworkIPAddressList to list of models.NetworkIPAddressList
func GetNetworkIPAddressListList(inputs []scaleiotypes.NetworkIPAddressList) []models.NetworkIPAddressList {
	out := make([]models.NetworkIPAddressList, 0)
	for _, input := range inputs {
		out = append(out, GetNetworkIPAddressList(input))
	}
	return out
}

// GetInterfaces converts scaleiotypes.Interfaces to models.Interfaces
func GetInterfaces(input scaleiotypes.Interfaces) models.Interfaces {
	return models.Interfaces{
		ID:            types.StringValue(input.ID),
		Name:          types.StringValue(input.Name),
		Partitioned:   types.BoolValue(input.Partitioned),
		Partitions:    GetPartitionsList(input.Partitions),
		Enabled:       types.BoolValue(input.Enabled),
		Redundancy:    types.BoolValue(input.Redundancy),
		Nictype:       types.StringValue(input.Nictype),
		Fqdd:          types.StringValue(input.Fqdd),
		MaxPartitions: types.Int64Value(int64(input.MaxPartitions)),
		AllNetworks:   GeStringList(input.AllNetworks),
	}
}

// GetPartitionsList converts list of scaleiotypes.Partitions to list of models.Partitions
func GetPartitionsList(inputs []scaleiotypes.Partitions) []models.Partitions {
	out := make([]models.Partitions, 0)
	for _, input := range inputs {
		out = append(out, GetPartitions(input))
	}
	return out
}

// GetInterfacesDetails converts scaleiotypes.InterfacesDetails to models.InterfacesDetails
func GetInterfacesDetails(input scaleiotypes.InterfacesDetails) models.InterfacesDetails {
	return models.InterfacesDetails{
		ID:            types.StringValue(input.ID),
		Name:          types.StringValue(input.Name),
		Redundancy:    types.BoolValue(input.Redundancy),
		Enabled:       types.BoolValue(input.Enabled),
		Partitioned:   types.BoolValue(input.Partitioned),
		Interfaces:    GetInterfacesList(input.Interfaces),
		Nictype:       types.StringValue(input.Nictype),
		Fabrictype:    types.StringValue(input.Fabrictype),
		MaxPartitions: types.Int64Value(int64(input.MaxPartitions)),
		Nports:        types.Int64Value(int64(input.Nports)),
		CardIndex:     types.Int64Value(int64(input.CardIndex)),
		NictypeSource: types.StringValue(input.NictypeSource),
	}
}

// GetInterfacesList converts list of scaleiotypes.Interfaces to list of models.Interfaces
func GetInterfacesList(inputs []scaleiotypes.Interfaces) []models.Interfaces {
	out := make([]models.Interfaces, 0)
	for _, input := range inputs {
		out = append(out, GetInterfaces(input))
	}
	return out
}

// GetNetworkConfiguration converts scaleiotypes.NetworkConfiguration to models.NetworkConfiguration
func GetNetworkConfiguration(input scaleiotypes.NetworkConfiguration) models.NetworkConfiguration {
	return models.NetworkConfiguration{
		ID:           types.StringValue(input.ID),
		Interfaces:   GetInterfacesDetailsList(input.Interfaces),
		SoftwareOnly: types.BoolValue(input.SoftwareOnly),
	}
}

// GetInterfacesDetailsList converts list of scaleiotypes.InterfacesDetails to list of models.InterfacesDetails
func GetInterfacesDetailsList(inputs []scaleiotypes.InterfacesDetails) []models.InterfacesDetails {
	out := make([]models.InterfacesDetails, 0)
	for _, input := range inputs {
		out = append(out, GetInterfacesDetails(input))
	}
	return out
}

// GetConfigurationDetails converts scaleiotypes.ConfigurationDetails to models.ConfigurationDetails
func GetConfigurationDetails(input scaleiotypes.ConfigurationDetails) models.ConfigurationDetails {
	return models.ConfigurationDetails{
		ID:              types.StringValue(input.ID),
		Disktype:        types.StringValue(input.Disktype),
		Comparator:      types.StringValue(input.Comparator),
		Numberofdisks:   types.Int64Value(int64(input.Numberofdisks)),
		Raidlevel:       types.StringValue(input.Raidlevel),
		VirtualDiskFqdd: types.StringValue(input.VirtualDiskFqdd),
		ControllerFqdd:  types.StringValue(input.ControllerFqdd),
		Categories:      GetCategoriesList(input.Categories),
	}
}

// GetCategoriesList converts list of scaleiotypes.Categories to list of models.Categories
func GetCategoriesList(inputs []scaleiotypes.Categories) []models.Categories {
	out := make([]models.Categories, 0)
	for _, input := range inputs {
		out = append(out, GetCategories(input))
	}
	return out
}

// GetVirtualDisks converts scaleiotypes.VirtualDisks to models.VirtualDisks
func GetVirtualDisks(input scaleiotypes.VirtualDisks) models.VirtualDisks {
	return models.VirtualDisks{
		PhysicalDisks:         GeStringList(input.PhysicalDisks),
		VirtualDiskFqdd:       types.StringValue(input.VirtualDiskFqdd),
		RaidLevel:             types.StringValue(input.RaidLevel),
		RollUpStatus:          types.StringValue(input.RollUpStatus),
		Controller:            types.StringValue(input.Controller),
		ControllerProductName: types.StringValue(input.ControllerProductName),
		Configuration:         GetConfigurationDetails(input.Configuration),
		MediaType:             types.StringValue(input.MediaType),
		EncryptionType:        types.StringValue(input.EncryptionType),
	}
}

// GetExternalVirtualDisks converts scaleiotypes.ExternalVirtualDisks to models.ExternalVirtualDisks
func GetExternalVirtualDisks(input scaleiotypes.ExternalVirtualDisks) models.ExternalVirtualDisks {
	return models.ExternalVirtualDisks{
		PhysicalDisks:         GeStringList(input.PhysicalDisks),
		VirtualDiskFqdd:       types.StringValue(input.VirtualDiskFqdd),
		RaidLevel:             types.StringValue(input.RaidLevel),
		RollUpStatus:          types.StringValue(input.RollUpStatus),
		Controller:            types.StringValue(input.Controller),
		ControllerProductName: types.StringValue(input.ControllerProductName),
		Configuration:         GetConfigurationDetails(input.Configuration),
		MediaType:             types.StringValue(input.MediaType),
		EncryptionType:        types.StringValue(input.EncryptionType),
	}
}

// GetSizeToDiskMap converts scaleiotypes.SizeToDiskMap to models.SizeToDiskMap
func GetSizeToDiskMap(input map[string]int) basetypes.MapValue {
	elements := make(map[string]attr.Value) 
	for key, val := range input {
		elements[key] = types.Int64Value(int64(val))
	}
	var setSizeToDisk basetypes.MapValue
	setSizeToDisk , _ = types.MapValue(types.Int64Type, elements) 
		
	 return setSizeToDisk
}

// GetRaidConfiguration converts scaleiotypes.RaidConfiguration to models.RaidConfiguration
func GetRaidConfiguration(input scaleiotypes.RaidConfiguration) models.RaidConfiguration {
	return models.RaidConfiguration{
		VirtualDisks:         GetVirtualDisksList(input.VirtualDisks),
		ExternalVirtualDisks: GetExternalVirtualDisksList(input.ExternalVirtualDisks),
		HddHotSpares:         GeStringList(input.HddHotSpares),
		SsdHotSpares:         GeStringList(input.SsdHotSpares),
		ExternalHddHotSpares: GeStringList(input.ExternalHddHotSpares),
		ExternalSsdHotSpares: GeStringList(input.ExternalSsdHotSpares),
		SizeToDiskMap:        GetSizeToDiskMap(input.SizeToDiskMap),
	}
}

// GetVirtualDisksList converts list of scaleiotypes.VirtualDisks to list of models.VirtualDisks
func GetVirtualDisksList(inputs []scaleiotypes.VirtualDisks) []models.VirtualDisks {
	out := make([]models.VirtualDisks, 0)
	for _, input := range inputs {
		out = append(out, GetVirtualDisks(input))
	}
	return out
}

// GetExternalVirtualDisksList converts list of scaleiotypes.ExternalVirtualDisks to list of models.ExternalVirtualDisks
func GetExternalVirtualDisksList(inputs []scaleiotypes.ExternalVirtualDisks) []models.ExternalVirtualDisks {
	out := make([]models.ExternalVirtualDisks, 0)
	for _, input := range inputs {
		out = append(out, GetExternalVirtualDisks(input))
	}
	return out
}

// GetAttributes converts scaleiotypes.Attributes to models.Attributes
func GetAttributes(input map[string]string) basetypes.MapValue {
	elements := make(map[string]attr.Value) 
	for key, val := range input {
		elements[key] = types.StringValue(val)
	}
	var setAttr basetypes.MapValue
	setAttr , _ = types.MapValue(types.StringType, elements) 

	 return setAttr	
}

// GetOptionsDetails converts scaleiotypes.OptionsDetails to models.OptionsDetails
func GetOptionsDetails(input scaleiotypes.OptionsDetails) models.OptionsDetails {
	return models.OptionsDetails{
		ID:           types.StringValue(input.ID),
		Name:         types.StringValue(input.Name),
		Value:        types.StringValue(input.Value),
		Dependencies: GetDependenciesDetailsList(input.Dependencies),
		Attributes:   GetAttributes(input.Attributes),
	}
}

// GetDependenciesDetailsList converts list of scaleiotypes.DependenciesDetails to list of models.DependenciesDetails
func GetDependenciesDetailsList(inputs []scaleiotypes.DependenciesDetails) []models.DependenciesDetails {
	out := make([]models.DependenciesDetails, 0)
	for _, input := range inputs {
		out = append(out, GetDependenciesDetails(input))
	}
	return out
}

// GetScaleIOStoragePoolDisks converts scaleiotypes.ScaleIOStoragePoolDisks to models.ScaleIOStoragePoolDisks
func GetScaleIOStoragePoolDisks(input scaleiotypes.ScaleIOStoragePoolDisks) models.ScaleIOStoragePoolDisks {
	return models.ScaleIOStoragePoolDisks{
		ProtectionDomainID:   types.StringValue(input.ProtectionDomainID),
		ProtectionDomainName: types.StringValue(input.ProtectionDomainName),
		StoragePoolID:        types.StringValue(input.StoragePoolID),
		StoragePoolName:      types.StringValue(input.StoragePoolName),
		DiskType:             types.StringValue(input.DiskType),
		PhysicalDiskFqdds:    GeStringList(input.PhysicalDiskFqdds),
		VirtualDiskFqdds:     GeStringList(input.VirtualDiskFqdds),
		SoftwareOnlyDisks:    GeStringList(input.SoftwareOnlyDisks),
	}
}

// GetScaleIODiskConfiguration converts scaleiotypes.ScaleIODiskConfiguration to models.ScaleIODiskConfiguration
func GetScaleIODiskConfiguration(input scaleiotypes.ScaleIODiskConfiguration) models.ScaleIODiskConfiguration {
	return models.ScaleIODiskConfiguration{
		ScaleIOStoragePoolDisks: GetScaleIOStoragePoolDisksList(input.ScaleIOStoragePoolDisks),
	}
}

// GetScaleIOStoragePoolDisksList converts list of scaleiotypes.ScaleIOStoragePoolDisks to list of models.ScaleIOStoragePoolDisks
func GetScaleIOStoragePoolDisksList(inputs []scaleiotypes.ScaleIOStoragePoolDisks) []models.ScaleIOStoragePoolDisks {
	out := make([]models.ScaleIOStoragePoolDisks, 0)
	for _, input := range inputs {
		out = append(out, GetScaleIOStoragePoolDisks(input))
	}
	return out
}

// GetShortWindow converts scaleiotypes.ShortWindow to models.ShortWindow
func GetShortWindow(input scaleiotypes.ShortWindow) models.ShortWindow {
	return models.ShortWindow{
		Threshold:       types.Int64Value(int64(input.Threshold)),
		WindowSizeInSec: types.Int64Value(int64(input.WindowSizeInSec)),
	}
}

// GetMediumWindow converts scaleiotypes.MediumWindow to models.MediumWindow
func GetMediumWindow(input scaleiotypes.MediumWindow) models.MediumWindow {
	return models.MediumWindow{
		Threshold:       types.Int64Value(int64(input.Threshold)),
		WindowSizeInSec: types.Int64Value(int64(input.WindowSizeInSec)),
	}
}

// GetLongWindow converts scaleiotypes.LongWindow to models.LongWindow
func GetLongWindow(input scaleiotypes.LongWindow) models.LongWindow {
	return models.LongWindow{
		Threshold:       types.Int64Value(int64(input.Threshold)),
		WindowSizeInSec: types.Int64Value(int64(input.WindowSizeInSec)),
	}
}

// GetSdsDecoupledCounterParameters converts scaleiotypes.SdsDecoupledCounterParameters to models.SdsDecoupledCounterParameters
func GetSdsDecoupledCounterParameters(input scaleiotypes.SdsDecoupledCounterParameters) models.SdsDecoupledCounterParameters {
	return models.SdsDecoupledCounterParameters{
		ShortWindow:  GetShortWindow(input.ShortWindow),
		MediumWindow: GetMediumWindow(input.MediumWindow),
		LongWindow:   GetLongWindow(input.LongWindow),
	}
}

// GetSdsConfigurationFailureCounterParameters converts scaleiotypes.SdsConfigurationFailureCounterParameters to models.SdsConfigurationFailureCounterParameters
func GetSdsConfigurationFailureCounterParameters(input scaleiotypes.SdsConfigurationFailureCounterParameters) models.SdsConfigurationFailureCounterParameters {
	return models.SdsConfigurationFailureCounterParameters{
		ShortWindow:  GetShortWindow(input.ShortWindow),
		MediumWindow: GetMediumWindow(input.MediumWindow),
		LongWindow:   GetLongWindow(input.LongWindow),
	}
}

// GetMdmSdsCounterParameters converts scaleiotypes.MdmSdsCounterParameters to models.MdmSdsCounterParameters
func GetMdmSdsCounterParameters(input scaleiotypes.MdmSdsCounterParameters) models.MdmSdsCounterParameters {
	return models.MdmSdsCounterParameters{
		ShortWindow:  GetShortWindow(input.ShortWindow),
		MediumWindow: GetMediumWindow(input.MediumWindow),
		LongWindow:   GetLongWindow(input.LongWindow),
	}
}

// GetSdsSdsCounterParameters converts scaleiotypes.SdsSdsCounterParameters to models.SdsSdsCounterParameters
func GetSdsSdsCounterParameters(input scaleiotypes.SdsSdsCounterParameters) models.SdsSdsCounterParameters {
	return models.SdsSdsCounterParameters{
		ShortWindow:  GetShortWindow(input.ShortWindow),
		MediumWindow: GetMediumWindow(input.MediumWindow),
		LongWindow:   GetLongWindow(input.LongWindow),
	}
}

// GetSdsReceiveBufferAllocationFailuresCounterParameters converts scaleiotypes.SdsReceiveBufferAllocationFailuresCounterParameters to models.SdsReceiveBufferAllocationFailuresCounterParameters
func GetSdsReceiveBufferAllocationFailuresCounterParameters(input scaleiotypes.SdsReceiveBufferAllocationFailuresCounterParameters) models.SdsReceiveBufferAllocationFailuresCounterParameters {
	return models.SdsReceiveBufferAllocationFailuresCounterParameters{
		ShortWindow:  GetShortWindow(input.ShortWindow),
		MediumWindow: GetMediumWindow(input.MediumWindow),
		LongWindow:   GetLongWindow(input.LongWindow),
	}
}

// GetGeneral converts scaleiotypes.General to models.General
func GetGeneral(input scaleiotypes.General) models.General {
	return models.General{
		ID:                                       types.StringValue(input.ID),
		Name:                                     types.StringValue(input.Name),
		SystemID:                                 types.StringValue(input.SystemID),
		ProtectionDomainState:                    types.StringValue(input.ProtectionDomainState),
		RebuildNetworkThrottlingInKbps:           types.Int64Value(int64(input.RebuildNetworkThrottlingInKbps)),
		RebalanceNetworkThrottlingInKbps:         types.Int64Value(int64(input.RebalanceNetworkThrottlingInKbps)),
		OverallIoNetworkThrottlingInKbps:         types.Int64Value(int64(input.OverallIoNetworkThrottlingInKbps)),
		SdsDecoupledCounterParameters:            GetSdsDecoupledCounterParameters(input.SdsDecoupledCounterParameters),
		SdsConfigurationFailureCounterParameters: GetSdsConfigurationFailureCounterParameters(input.SdsConfigurationFailureCounterParameters),
		MdmSdsCounterParameters:                  GetMdmSdsCounterParameters(input.MdmSdsCounterParameters),
		SdsSdsCounterParameters:                  GetSdsSdsCounterParameters(input.SdsSdsCounterParameters),
		RfcacheOpertionalMode:                    types.StringValue(input.RfcacheOpertionalMode),
		RfcachePageSizeKb:                        types.Int64Value(int64(input.RfcachePageSizeKb)),
		RfcacheMaxIoSizeKb:                       types.Int64Value(int64(input.RfcacheMaxIoSizeKb)),
		SdsReceiveBufferAllocationFailuresCounterParameters: GetSdsReceiveBufferAllocationFailuresCounterParameters(input.SdsReceiveBufferAllocationFailuresCounterParameters),
		RebuildNetworkThrottlingEnabled:                     types.BoolValue(input.RebuildNetworkThrottlingEnabled),
		RebalanceNetworkThrottlingEnabled:                   types.BoolValue(input.RebalanceNetworkThrottlingEnabled),
		OverallIoNetworkThrottlingEnabled:                   types.BoolValue(input.OverallIoNetworkThrottlingEnabled),
		RfcacheEnabled:                                      types.BoolValue(input.RfcacheEnabled),
	}
}

// GetStatisticsDetails converts scaleiotypes.StatisticsDetails to models.StatisticsDetails
func GetStatisticsDetails(input scaleiotypes.StatisticsDetails) models.StatisticsDetails {
	return models.StatisticsDetails{
		NumOfDevices:                             types.Int64Value(int64(input.NumOfDevices)),
		UnusedCapacityInKb:                       types.Int64Value(int64(input.UnusedCapacityInKb)),
		NumOfVolumes:                             types.Int64Value(int64(input.NumOfVolumes)),
		NumOfMappedToAllVolumes:                  types.Int64Value(int64(input.NumOfMappedToAllVolumes)),
		CapacityAvailableForVolumeAllocationInKb: types.Int64Value(int64(input.CapacityAvailableForVolumeAllocationInKb)),
		VolumeAllocationLimitInKb:                types.Int64Value(int64(input.VolumeAllocationLimitInKb)),
		CapacityLimitInKb:                        types.Int64Value(int64(input.CapacityLimitInKb)),
		NumOfUnmappedVolumes:                     types.Int64Value(int64(input.NumOfUnmappedVolumes)),
		SpareCapacityInKb:                        types.Int64Value(int64(input.SpareCapacityInKb)),
		CapacityInUseInKb:                        types.Int64Value(int64(input.CapacityInUseInKb)),
		MaxCapacityInKb:                          types.Int64Value(int64(input.MaxCapacityInKb)),
		NumOfSds:                                 types.Int64Value(int64(input.NumOfSds)),
		NumOfStoragePools:                        types.Int64Value(int64(input.NumOfStoragePools)),
		NumOfFaultSets:                           types.Int64Value(int64(input.NumOfFaultSets)),
		ThinCapacityInUseInKb:                    types.Int64Value(int64(input.ThinCapacityInUseInKb)),
		ThickCapacityInUseInKb:                   types.Int64Value(int64(input.ThickCapacityInUseInKb)),
	}
}

// GetDiskList converts scaleiotypes.DiskList to models.DiskList
func GetDiskList(input scaleiotypes.DiskList) models.DiskList {
	return models.DiskList{
		ID:                     types.StringValue(input.ID),
		Name:                   types.StringValue(input.Name),
		ErrorState:             types.StringValue(input.ErrorState),
		SdsID:                  types.StringValue(input.SdsID),
		DeviceState:            types.StringValue(input.DeviceState),
		CapacityLimitInKb:      types.Int64Value(int64(input.CapacityLimitInKb)),
		MaxCapacityInKb:        types.Int64Value(int64(input.MaxCapacityInKb)),
		StoragePoolID:          types.StringValue(input.StoragePoolID),
		DeviceCurrentPathName:  types.StringValue(input.DeviceCurrentPathName),
		DeviceOriginalPathName: types.StringValue(input.DeviceOriginalPathName),
		SerialNumber:           types.StringValue(input.SerialNumber),
		VendorName:             types.StringValue(input.VendorName),
		ModelName:              types.StringValue(input.ModelName),
	}
}

// GetMappedSdcInfoDetails converts scaleiotypes.MappedSdcInfoDetails to models.MappedSdcInfoDetails
func GetMappedSdcInfoDetails(input scaleiotypes.MappedSdcInfoDetails) models.MappedSdcInfoDetails {
	return models.MappedSdcInfoDetails{
		SdcIP:         types.StringValue(input.SdcIP),
		SdcID:         types.StringValue(input.SdcID),
		LimitBwInMbps: types.Int64Value(int64(input.LimitBwInMbps)),
		LimitIops:     types.Int64Value(int64(input.LimitIops)),
	}
}

// GetVolumeList converts scaleiotypes.VolumeList to models.VolumeList
func GetVolumeList(input scaleiotypes.VolumeList) models.VolumeList {
	return models.VolumeList{
		ID:                types.StringValue(input.ID),
		Name:              types.StringValue(input.Name),
		VolumeType:        types.StringValue(input.VolumeType),
		StoragePoolID:     types.StringValue(input.StoragePoolID),
		DataLayout:        types.StringValue(input.DataLayout),
		CompressionMethod: types.StringValue(input.CompressionMethod),
		SizeInKb:          types.Int64Value(int64(input.SizeInKb)),
		MappedSdcInfo:     GetMappedSdcInfoDetailsList(input.MappedSdcInfo),
		VolumeClass:       types.StringValue(input.VolumeClass),
	}
}

// GetMappedSdcInfoDetailsList converts list of scaleiotypes.MappedSdcInfoDetails to list of models.MappedSdcInfoDetails
func GetMappedSdcInfoDetailsList(inputs []scaleiotypes.MappedSdcInfoDetails) []models.MappedSdcInfoDetails {
	out := make([]models.MappedSdcInfoDetails, 0)
	for _, input := range inputs {
		out = append(out, GetMappedSdcInfoDetails(input))
	}
	return out
}

// GetStoragePoolList converts scaleiotypes.StoragePoolList to models.StoragePoolList
func GetStoragePoolList(input scaleiotypes.StoragePoolList) models.StoragePoolList {
	return models.StoragePoolList{
		ID:                        types.StringValue(input.ID),
		Name:                      types.StringValue(input.Name),
		RebuildIoPriorityPolicy:   types.StringValue(input.RebuildIoPriorityPolicy),
		RebalanceIoPriorityPolicy: types.StringValue(input.RebalanceIoPriorityPolicy),
		RebuildIoPriorityNumOfConcurrentIosPerDevice:     types.Int64Value(int64(input.RebuildIoPriorityNumOfConcurrentIosPerDevice)),
		RebalanceIoPriorityNumOfConcurrentIosPerDevice:   types.Int64Value(int64(input.RebalanceIoPriorityNumOfConcurrentIosPerDevice)),
		RebuildIoPriorityBwLimitPerDeviceInKbps:          types.Int64Value(int64(input.RebuildIoPriorityBwLimitPerDeviceInKbps)),
		RebalanceIoPriorityBwLimitPerDeviceInKbps:        types.Int64Value(int64(input.RebalanceIoPriorityBwLimitPerDeviceInKbps)),
		RebuildIoPriorityAppIopsPerDeviceThreshold:       types.StringValue(input.RebuildIoPriorityAppIopsPerDeviceThreshold),
		RebalanceIoPriorityAppIopsPerDeviceThreshold:     types.StringValue(input.RebalanceIoPriorityAppIopsPerDeviceThreshold),
		RebuildIoPriorityAppBwPerDeviceThresholdInKbps:   types.Int64Value(int64(input.RebuildIoPriorityAppBwPerDeviceThresholdInKbps)),
		RebalanceIoPriorityAppBwPerDeviceThresholdInKbps: types.Int64Value(int64(input.RebalanceIoPriorityAppBwPerDeviceThresholdInKbps)),
		RebuildIoPriorityQuietPeriodInMsec:               types.Int64Value(int64(input.RebuildIoPriorityQuietPeriodInMsec)),
		RebalanceIoPriorityQuietPeriodInMsec:             types.Int64Value(int64(input.RebalanceIoPriorityQuietPeriodInMsec)),
		ZeroPaddingEnabled:                               types.BoolValue(input.ZeroPaddingEnabled),
		BackgroundScannerMode:                            types.StringValue(input.BackgroundScannerMode),
		BackgroundScannerBWLimitKBps:                     types.Int64Value(int64(input.BackgroundScannerBWLimitKBps)),
		UseRmcache:                                       types.BoolValue(input.UseRmcache),
		ProtectionDomainID:                               types.StringValue(input.ProtectionDomainID),
		SpClass:                                          types.StringValue(input.SpClass),
		UseRfcache:                                       types.BoolValue(input.UseRfcache),
		SparePercentage:                                  types.Int64Value(int64(input.SparePercentage)),
		RmcacheWriteHandlingMode:                         types.StringValue(input.RmcacheWriteHandlingMode),
		ChecksumEnabled:                                  types.BoolValue(input.ChecksumEnabled),
		RebuildEnabled:                                   types.BoolValue(input.RebuildEnabled),
		RebalanceEnabled:                                 types.BoolValue(input.RebalanceEnabled),
		NumOfParallelRebuildRebalanceJobsPerDevice:       types.Int64Value(int64(input.NumOfParallelRebuildRebalanceJobsPerDevice)),
		CapacityAlertHighThreshold:                       types.Int64Value(int64(input.CapacityAlertHighThreshold)),
		CapacityAlertCriticalThreshold:                   types.Int64Value(int64(input.CapacityAlertCriticalThreshold)),
		Statistics:                                       GetStatisticsDetails(input.Statistics),
		DataLayout:                                       types.StringValue(input.DataLayout),
		ReplicationCapacityMaxRatio:                      types.StringValue(input.ReplicationCapacityMaxRatio),
		MediaType:                                        types.StringValue(input.MediaType),
		DiskList:                                         GetDiskListList(input.DiskList),
		VolumeList:                                       GetVolumeListList(input.VolumeList),
		FglAccpID:                                        types.StringValue(input.FglAccpID),
	}
}

// GetDiskListList converts list of scaleiotypes.DiskList to list of models.DiskList
func GetDiskListList(inputs []scaleiotypes.DiskList) []models.DiskList {
	out := make([]models.DiskList, 0)
	for _, input := range inputs {
		out = append(out, GetDiskList(input))
	}
	return out
}

// GetVolumeListList converts list of scaleiotypes.VolumeList to list of models.VolumeList
func GetVolumeListList(inputs []scaleiotypes.VolumeList) []models.VolumeList {
	out := make([]models.VolumeList, 0)
	for _, input := range inputs {
		out = append(out, GetVolumeList(input))
	}
	return out
}

// GetIPList converts scaleiotypes.IPList to models.IPList
func GetIPList(input scaleiotypes.IPList) models.IPList {
	return models.IPList{
		IP:   types.StringValue(input.IP),
		Role: types.StringValue(input.Role),
	}
}

// GetSdsListDetails converts scaleiotypes.SdsListDetails to models.SdsListDetails
func GetSdsListDetails(input scaleiotypes.SdsListDetails) models.SdsListDetails {
	return models.SdsListDetails{
		ID:                  types.StringValue(input.ID),
		Name:                types.StringValue(input.Name),
		Port:                types.Int64Value(int64(input.Port)),
		ProtectionDomainID:  types.StringValue(input.ProtectionDomainID),
		FaultSetID:          types.StringValue(input.FaultSetID),
		SoftwareVersionInfo: types.StringValue(input.SoftwareVersionInfo),
		SdsState:            types.StringValue(input.SdsState),
		MembershipState:     types.StringValue(input.MembershipState),
		MdmConnectionState:  types.StringValue(input.MdmConnectionState),
		DrlMode:             types.StringValue(input.DrlMode),
		MaintenanceState:    types.StringValue(input.MaintenanceState),
		PerfProfile:         types.StringValue(input.PerfProfile),
		OnVMWare:            types.BoolValue(input.OnVMWare),
		IPList:              GetIPListList(input.IPList),
	}
}

// GetIPListList converts list of scaleiotypes.IPList to list of models.IPList
func GetIPListList(inputs []scaleiotypes.IPList) []models.IPList {
	out := make([]models.IPList, 0)
	for _, input := range inputs {
		out = append(out, GetIPList(input))
	}
	return out
}

// GetSdrListDetails converts scaleiotypes.SdrListDetails to models.SdrListDetails
func GetSdrListDetails(input scaleiotypes.SdrListDetails) models.SdrListDetails {
	return models.SdrListDetails{
		ID:                  types.StringValue(input.ID),
		Name:                types.StringValue(input.Name),
		Port:                types.Int64Value(int64(input.Port)),
		ProtectionDomainID:  types.StringValue(input.ProtectionDomainID),
		SoftwareVersionInfo: types.StringValue(input.SoftwareVersionInfo),
		SdrState:            types.StringValue(input.SdrState),
		MembershipState:     types.StringValue(input.MembershipState),
		MdmConnectionState:  types.StringValue(input.MdmConnectionState),
		MaintenanceState:    types.StringValue(input.MaintenanceState),
		PerfProfile:         types.StringValue(input.PerfProfile),
		IPList:              GetIPListList(input.IPList),
	}
}

// GetAccelerationPool converts scaleiotypes.AccelerationPool to models.AccelerationPool
func GetAccelerationPool(input scaleiotypes.AccelerationPool) models.AccelerationPool {
	return models.AccelerationPool{
		ID:                 types.StringValue(input.ID),
		Name:               types.StringValue(input.Name),
		ProtectionDomainID: types.StringValue(input.ProtectionDomainID),
		MediaType:          types.StringValue(input.MediaType),
		Rfcache:            types.BoolValue(input.Rfcache),
	}
}

// GetProtectionDomainSettings converts scaleiotypes.ProtectionDomainSettings to models.ProtectionDomainSettings
func GetProtectionDomainSettings(input scaleiotypes.ProtectionDomainSettings) models.ProtectionDomainSettings {
	return models.ProtectionDomainSettings{
		General:          GetGeneral(input.General),
		Statistics:       GetStatisticsDetails(input.Statistics),
		StoragePoolList:  GetStoragePoolListList(input.StoragePoolList),
		SdsList:          GetSdsListDetailsList(input.SdsList),
		SdrList:          GetSdrListDetailsList(input.SdrList),
		AccelerationPool: GetAccelerationPoolList(input.AccelerationPool),
	}
}

// GetStoragePoolListList converts list of scaleiotypes.StoragePoolList to list of models.StoragePoolList
func GetStoragePoolListList(inputs []scaleiotypes.StoragePoolList) []models.StoragePoolList {
	out := make([]models.StoragePoolList, 0)
	for _, input := range inputs {
		out = append(out, GetStoragePoolList(input))
	}
	return out
}

// GetSdsListDetailsList converts list of scaleiotypes.SdsListDetails to list of models.SdsListDetails
func GetSdsListDetailsList(inputs []scaleiotypes.SdsListDetails) []models.SdsListDetails {
	out := make([]models.SdsListDetails, 0)
	for _, input := range inputs {
		out = append(out, GetSdsListDetails(input))
	}
	return out
}

// GetSdrListDetailsList converts list of scaleiotypes.SdrListDetails to list of models.SdrListDetails
func GetSdrListDetailsList(inputs []scaleiotypes.SdrListDetails) []models.SdrListDetails {
	out := make([]models.SdrListDetails, 0)
	for _, input := range inputs {
		out = append(out, GetSdrListDetails(input))
	}
	return out
}

// GetAccelerationPoolList converts list of scaleiotypes.AccelerationPool to list of models.AccelerationPool
func GetAccelerationPoolList(inputs []scaleiotypes.AccelerationPool) []models.AccelerationPool {
	out := make([]models.AccelerationPool, 0)
	for _, input := range inputs {
		out = append(out, GetAccelerationPool(input))
	}
	return out
}

// GetFaultSetSettings converts scaleiotypes.FaultSetSettings to models.FaultSetSettings
func GetFaultSetSettings(input scaleiotypes.FaultSetSettings) models.FaultSetSettings {
	return models.FaultSetSettings{
		ProtectionDomainID: types.StringValue(input.ProtectionDomainID),
		Name:               types.StringValue(input.Name),
		ID:                 types.StringValue(input.ID),
	}
}

// GetDatacenter converts scaleiotypes.Datacenter to models.Datacenter
func GetDatacenter(input scaleiotypes.Datacenter) models.Datacenter {
	return models.Datacenter{
		VcenterID:      types.StringValue(input.VcenterID),
		DatacenterID:   types.StringValue(input.DatacenterID),
		DatacenterName: types.StringValue(input.DatacenterName),
	}
}

// GetPortGroupOptions converts scaleiotypes.PortGroupOptions to models.PortGroupOptions
func GetPortGroupOptions(input scaleiotypes.PortGroupOptions) models.PortGroupOptions {
	return models.PortGroupOptions{
		ID:   types.StringValue(input.ID),
		Name: types.StringValue(input.Name),
	}
}

// GetPortGroups converts scaleiotypes.PortGroups to models.PortGroups
func GetPortGroups(input scaleiotypes.PortGroups) models.PortGroups {
	return models.PortGroups{
		ID:               types.StringValue(input.ID),
		DisplayName:      types.StringValue(input.DisplayName),
		Vlan:             types.Int64Value(int64(input.Vlan)),
		Name:             types.StringValue(input.Name),
		Value:            types.StringValue(input.Value),
		PortGroupOptions: GetPortGroupOptionsList(input.PortGroupOptions),
	}
}

// GetPortGroupOptionsList converts list of scaleiotypes.PortGroupOptions to list of models.PortGroupOptions
func GetPortGroupOptionsList(inputs []scaleiotypes.PortGroupOptions) []models.PortGroupOptions {
	out := make([]models.PortGroupOptions, 0)
	for _, input := range inputs {
		out = append(out, GetPortGroupOptions(input))
	}
	return out
}

// GetVdsSettings converts scaleiotypes.VdsSettings to models.VdsSettings
func GetVdsSettings(input scaleiotypes.VdsSettings) models.VdsSettings {
	return models.VdsSettings{
		ID:          types.StringValue(input.ID),
		DisplayName: types.StringValue(input.DisplayName),
		Name:        types.StringValue(input.Name),
		Value:       types.StringValue(input.Value),
		PortGroups:  GetPortGroupsList(input.PortGroups),
	}
}

// GetPortGroupsList converts list of scaleiotypes.PortGroups to list of models.PortGroups
func GetPortGroupsList(inputs []scaleiotypes.PortGroups) []models.PortGroups {
	out := make([]models.PortGroups, 0)
	for _, input := range inputs {
		out = append(out, GetPortGroups(input))
	}
	return out
}

// GetVdsNetworkMtuSizeConfiguration converts scaleiotypes.VdsNetworkMtuSizeConfiguration to models.VdsNetworkMtuSizeConfiguration
func GetVdsNetworkMtuSizeConfiguration(input scaleiotypes.VdsNetworkMtuSizeConfiguration) models.VdsNetworkMtuSizeConfiguration {
	return models.VdsNetworkMtuSizeConfiguration{
		ID:    types.StringValue(input.ID),
		Value: types.StringValue(input.Value),
	}
}

// GetVdsNetworkMTUSizeConfiguration converts scaleiotypes.VdsNetworkMTUSizeConfiguration to models.VdsNetworkMTUSizeConfiguration
func GetVdsNetworkMTUSizeConfiguration(input scaleiotypes.VdsNetworkMTUSizeConfiguration) models.VdsNetworkMTUSizeConfiguration {
	return models.VdsNetworkMTUSizeConfiguration{
		ID:    types.StringValue(input.ID),
		Value: types.StringValue(input.Value),
	}
}

// GetVdsConfiguration converts scaleiotypes.VdsConfiguration to models.VdsConfiguration
func GetVdsConfiguration(input scaleiotypes.VdsConfiguration) models.VdsConfiguration {
	return models.VdsConfiguration{
		Datacenter:                     GetDatacenter(input.Datacenter),
		PortGroupOption:                types.StringValue(input.PortGroupOption),
		PortGroupCreationOption:        types.StringValue(input.PortGroupCreationOption),
		VdsSettings:                    GetVdsSettingsList(input.VdsSettings),
		VdsNetworkMtuSizeConfiguration: GetVdsNetworkMtuSizeConfigurationList(input.VdsNetworkMtuSizeConfiguration),
	}
}

// GetVdsSettingsList converts list of scaleiotypes.VdsSettings to list of models.VdsSettings
func GetVdsSettingsList(inputs []scaleiotypes.VdsSettings) []models.VdsSettings {
	out := make([]models.VdsSettings, 0)
	for _, input := range inputs {
		out = append(out, GetVdsSettings(input))
	}
	return out
}

// GetVdsNetworkMtuSizeConfigurationList converts list of scaleiotypes.VdsNetworkMtuSizeConfiguration to list of models.VdsNetworkMtuSizeConfiguration
func GetVdsNetworkMtuSizeConfigurationList(inputs []scaleiotypes.VdsNetworkMtuSizeConfiguration) []models.VdsNetworkMtuSizeConfiguration {
	out := make([]models.VdsNetworkMtuSizeConfiguration, 0)
	for _, input := range inputs {
		out = append(out, GetVdsNetworkMtuSizeConfiguration(input))
	}
	return out
}

// GetVdsNetworkMTUSizeConfigurationList converts list of scaleiotypes.VdsNetworkMTUSizeConfiguration to list of models.VdsNetworkMTUSizeConfiguration
func GetVdsNetworkMTUSizeConfigurationList(inputs []scaleiotypes.VdsNetworkMTUSizeConfiguration) []models.VdsNetworkMTUSizeConfiguration {
	out := make([]models.VdsNetworkMTUSizeConfiguration, 0)
	for _, input := range inputs {
		out = append(out, GetVdsNetworkMTUSizeConfiguration(input))
	}
	return out
}

// GetNodeSelection converts scaleiotypes.NodeSelection to models.NodeSelection
func GetNodeSelection(input scaleiotypes.NodeSelection) models.NodeSelection {
	return models.NodeSelection{
		ID:            types.StringValue(input.ID),
		ServiceTag:    types.StringValue(input.ServiceTag),
		MgmtIPAddress: types.StringValue(input.MgmtIPAddress),
	}
}

// GetParametersDetails converts scaleiotypes.ParametersDetails to models.ParametersDetails
func GetParametersDetails(input scaleiotypes.ParametersDetails) models.ParametersDetails {
	return models.ParametersDetails{
		GUID:                     types.StringValue(input.GUID),
		ID:                       types.StringValue(input.ID),
		Type:                     types.StringValue(input.Type),
		DisplayName:              types.StringValue(input.DisplayName),
		Value:                    types.StringValue(input.Value),
		ToolTip:                  types.StringValue(input.ToolTip),
		Required:                 types.BoolValue(input.Required),
		RequiredAtDeployment:     types.BoolValue(input.RequiredAtDeployment),
		HideFromTemplate:         types.BoolValue(input.HideFromTemplate),
		Dependencies:             GetDependenciesDetailsList(input.Dependencies),
		Group:                    types.StringValue(input.Group),
		ReadOnly:                 types.BoolValue(input.ReadOnly),
		Generated:                types.BoolValue(input.Generated),
		InfoIcon:                 types.BoolValue(input.InfoIcon),
		Step:                     types.Int64Value(int64(input.Step)),
		MaxLength:                types.Int64Value(int64(input.MaxLength)),
		Min:                      types.Int64Value(int64(input.Min)),
		Max:                      types.Int64Value(int64(input.Max)),
		NetworkIPAddressList:     GetNetworkIPAddressListList(input.NetworkIPAddressList),
		NetworkConfiguration:     GetNetworkConfiguration(input.NetworkConfiguration),
		RaidConfiguration:        GetRaidConfiguration(input.RaidConfiguration),
		Options:                  GetOptionsDetailsList(input.Options),
		OptionsSortable:          types.BoolValue(input.OptionsSortable),
		PreservedForDeployment:   types.BoolValue(input.PreservedForDeployment),
		ScaleIODiskConfiguration: GetScaleIODiskConfiguration(input.ScaleIODiskConfiguration),
		ProtectionDomainSettings: GetProtectionDomainSettingsList(input.ProtectionDomainSettings),
		FaultSetSettings:         GetFaultSetSettingsList(input.FaultSetSettings),
		Attributes:               GetAttributes(input.Attributes),
		VdsConfiguration:         GetVdsConfiguration(input.VdsConfiguration),
		NodeSelection:            GetNodeSelection(input.NodeSelection),
	}
}

// GetOptionsDetailsList converts list of scaleiotypes.OptionsDetails to list of models.OptionsDetails
func GetOptionsDetailsList(inputs []scaleiotypes.OptionsDetails) []models.OptionsDetails {
	out := make([]models.OptionsDetails, 0)
	for _, input := range inputs {
		out = append(out, GetOptionsDetails(input))
	}
	return out
}

// GetProtectionDomainSettingsList converts list of scaleiotypes.ProtectionDomainSettings to list of models.ProtectionDomainSettings
func GetProtectionDomainSettingsList(inputs []scaleiotypes.ProtectionDomainSettings) []models.ProtectionDomainSettings {
	out := make([]models.ProtectionDomainSettings, 0)
	for _, input := range inputs {
		out = append(out, GetProtectionDomainSettings(input))
	}
	return out
}

// GetFaultSetSettingsList converts list of scaleiotypes.FaultSetSettings to list of models.FaultSetSettings
func GetFaultSetSettingsList(inputs []scaleiotypes.FaultSetSettings) []models.FaultSetSettings {
	out := make([]models.FaultSetSettings, 0)
	for _, input := range inputs {
		out = append(out, GetFaultSetSettings(input))
	}
	return out
}

// GetOptionsList converts list of scaleiotypes.Options to list of models.Options
func GetOptionsList(inputs []scaleiotypes.Options) []models.Options {
	out := make([]models.Options, 0)
	for _, input := range inputs {
		out = append(out, GetOptions(input))
	}
	return out
}

// GetResources converts scaleiotypes.Resources to models.Resources
func GetResources(input scaleiotypes.Resources) models.Resources {
	return models.Resources{
		GUID:        types.StringValue(input.GUID),
		ID:          types.StringValue(input.ID),
		DisplayName: types.StringValue(input.DisplayName),
		// Parameters:    GetParametersDetailsList(input.Parameters),
		// ParametersMap: GetParametersMap(input.ParametersMap),
	}
}

// GetComponents converts scaleiotypes.Components to models.Components
func GetComponents(input scaleiotypes.Components) models.Components {
	return models.Components{
		ID:                  types.StringValue(input.ID),
		ComponentID:         types.StringValue(input.ComponentID),
		Identifier:          types.StringValue(input.Identifier),
		ComponentValid:      GetComponentValid(input.ComponentValid),
		Name:                types.StringValue(input.Name),
		HelpText:            types.StringValue(input.HelpText),
		ClonedFromID:        types.StringValue(input.ClonedFromID),
		Teardown:            types.BoolValue(input.Teardown),
		Type:                types.StringValue(input.Type),
		SubType:             types.StringValue(input.SubType),
		RelatedComponents:   GetRelatedComponents(input.RelatedComponents),
		Resources:           GetResourcesList(input.Resources),
		Brownfield:          types.BoolValue(input.Brownfield),
		PuppetCertName:      types.StringValue(input.PuppetCertName),
		OsPuppetCertName:    types.StringValue(input.OsPuppetCertName),
		ManagementIPAddress: types.StringValue(input.ManagementIPAddress),
		SerialNumber:        types.StringValue(input.SerialNumber),
		AsmGUID:             types.StringValue(input.AsmGUID),
		Cloned:              types.BoolValue(input.Cloned),
		ConfigFile:          types.StringValue(input.ConfigFile),
		ManageFirmware:      types.BoolValue(input.ManageFirmware),
		Instances:           types.Int64Value(int64(input.Instances)),
		RefID:               types.StringValue(input.RefID),
		ClonedFromAsmGUID:   types.StringValue(input.ClonedFromAsmGUID),
		Changed:             types.BoolValue(input.Changed),
		IP:                  types.StringValue(input.IP),
	}
}

// GetResourcesList converts list of scaleiotypes.Resources to list of models.Resources
func GetResourcesList(inputs []scaleiotypes.Resources) []models.Resources {
	out := make([]models.Resources, 0)
	for _, input := range inputs {
		out = append(out, GetResources(input))
	}
	return out
}

// GetIPRange converts scaleiotypes.IPRange to models.IPRange
func GetIPRange(input scaleiotypes.IPRange) models.IPRange {
	return models.IPRange{
		ID:         types.StringValue(input.ID),
		StartingIP: types.StringValue(input.StartingIP),
		EndingIP:   types.StringValue(input.EndingIP),
		Role:       types.StringValue(input.Role),
	}
}

// GetStaticRoute converts scaleiotypes.StaticRoute to models.StaticRoute
func GetStaticRoute(input scaleiotypes.StaticRoute) models.StaticRoute {
	return models.StaticRoute{
		StaticRouteSourceNetworkID:      types.StringValue(input.StaticRouteSourceNetworkID),
		StaticRouteDestinationNetworkID: types.StringValue(input.StaticRouteDestinationNetworkID),
		StaticRouteGateway:              types.StringValue(input.StaticRouteGateway),
		SubnetMask:                      types.StringValue(input.SubnetMask),
		DestinationIPAddress:            types.StringValue(input.DestinationIPAddress),
	}
}

// GetStaticNetworkConfiguration converts scaleiotypes.StaticNetworkConfiguration to models.StaticNetworkConfiguration
func GetStaticNetworkConfiguration(input scaleiotypes.StaticNetworkConfiguration) models.StaticNetworkConfiguration {
	return models.StaticNetworkConfiguration{
		Gateway:      types.StringValue(input.Gateway),
		Subnet:       types.StringValue(input.Subnet),
		PrimaryDNS:   types.StringValue(input.PrimaryDNS),
		SecondaryDNS: types.StringValue(input.SecondaryDNS),
		DNSSuffix:    types.StringValue(input.DNSSuffix),
		IPRange:      GetIPRangeList(input.IPRange),
		IPAddress:    types.StringValue(input.IPAddress),
		StaticRoute:  GetStaticRouteList(input.StaticRoute),
	}
}

// GetIPRangeList converts list of scaleiotypes.IPRange to list of models.IPRange
func GetIPRangeList(inputs []scaleiotypes.IPRange) []models.IPRange {
	out := make([]models.IPRange, 0)
	for _, input := range inputs {
		out = append(out, GetIPRange(input))
	}
	return out
}

// GetStaticRouteList converts list of scaleiotypes.StaticRoute to list of models.StaticRoute
func GetStaticRouteList(inputs []scaleiotypes.StaticRoute) []models.StaticRoute {
	out := make([]models.StaticRoute, 0)
	for _, input := range inputs {
		out = append(out, GetStaticRoute(input))
	}
	return out
}

// GetNetworks converts scaleiotypes.Networks to models.Networks
func GetNetworks(input scaleiotypes.Networks) models.Networks {
	return models.Networks{
		ID:                         types.StringValue(input.ID),
		Name:                       types.StringValue(input.Name),
		Description:                types.StringValue(input.Description),
		VlanID:                     types.Int64Value(int64(input.VlanID)),
		StaticNetworkConfiguration: GetStaticNetworkConfiguration(input.StaticNetworkConfiguration),
		DestinationIPAddress:       types.StringValue(input.DestinationIPAddress),
		Static:                     types.BoolValue(input.Static),
		Type:                       types.StringValue(input.Type),
	}
}

// GetOptions converts scaleiotypes.Options to models.Options
func GetOptions(input scaleiotypes.Options) models.Options {
	return models.Options{
		ID:           types.StringValue(input.ID),
		Name:         types.StringValue(input.Name),
		Dependencies: GetDependenciesDetailsList(input.Dependencies),
		Attributes:   GetAttributes(input.Attributes),
	}
}

// GetParameters converts scaleiotypes.Parameters to models.Parameters
func GetParameters(input scaleiotypes.Parameters) models.Parameters {
	return models.Parameters{
		ID:               types.StringValue(input.ID),
		Value:            types.StringValue(input.Value),
		DisplayName:      types.StringValue(input.DisplayName),
		Type:             types.StringValue(input.Type),
		ToolTip:          types.StringValue(input.ToolTip),
		Required:         types.BoolValue(input.Required),
		HideFromTemplate: types.BoolValue(input.HideFromTemplate),
		DeviceType:       types.StringValue(input.DeviceType),
		Dependencies:     GetDependenciesDetailsList(input.Dependencies),
		Group:            types.StringValue(input.Group),
		ReadOnly:         types.BoolValue(input.ReadOnly),
		Generated:        types.BoolValue(input.Generated),
		InfoIcon:         types.BoolValue(input.InfoIcon),
		Step:             types.Int64Value(int64(input.Step)),
		MaxLength:        types.Int64Value(int64(input.MaxLength)),
		Min:              types.Int64Value(int64(input.Min)),
		Max:              types.Int64Value(int64(input.Max)),
		Networks:         GetNetworksList(input.Networks),
		Options:          GetOptionsList(input.Options),
		OptionsSortable:  types.BoolValue(input.OptionsSortable),
	}
}

// GetCategories converts scaleiotypes.Categories to models.Categories
func GetCategories(input scaleiotypes.Categories) models.Categories {
	return models.Categories{
		ID:          types.StringValue(input.ID),
		DisplayName: types.StringValue(input.DisplayName),
		DeviceType:  types.StringValue(input.DeviceType),
	}
}
