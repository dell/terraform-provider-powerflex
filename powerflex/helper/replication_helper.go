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
	"context"
	"fmt"
	"reflect"
	"strings"
	sshClient "terraform-provider-powerflex/client"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/dell/goscaleio"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"
)

// GetPeerSystem GET peer system
func GetPeerSystem(client *goscaleio.Client, peerSystemID string) (*scaleiotypes.PeerMDM, error) {
	return client.GetPeerMDM(peerSystemID)
}

// MapPeerSystemResourceState map peer system state
func MapPeerSystemResourceState(peerSystem scaleiotypes.PeerMDM, state models.PeerMdmResourceModel) models.PeerMdmResourceModel {
	return models.PeerMdmResourceModel{
		ID:                        types.StringValue(peerSystem.ID),
		Name:                      state.Name,
		Port:                      state.Port,
		PeerSystemID:              state.PeerSystemID,
		SystemID:                  types.StringValue(peerSystem.SystemID),
		SoftwareVersionInfo:       types.StringValue(peerSystem.SoftwareVersionInfo),
		MembershipState:           types.StringValue(peerSystem.MembershipState),
		PerfProfile:               state.PerfProfile,
		NetworkType:               types.StringValue(peerSystem.NetworkType),
		CouplingRC:                types.StringValue(peerSystem.CouplingRC),
		IPList:                    state.IPList,
		AddCertificate:            state.AddCertificate,
		DestinationPrimaryMdmInfo: fillInPrimaryMdmInfo(state.DestinationPrimaryMdmInfo),
		SourcePrimaryMdmInfo:      fillInPrimaryMdmInfo(state.SourcePrimaryMdmInfo),
	}
}

// Fill in the Primary Mdm Info, if empty fill in the object with null values
func fillInPrimaryMdmInfo(info basetypes.ObjectValue) basetypes.ObjectValue {
	if info.IsUnknown() {
		val, _ := types.ObjectValue(map[string]attr.Type{
			"ip":                  types.StringType,
			"ssh_username":        types.StringType,
			"ssh_password":        types.StringType,
			"ssh_port":            types.StringType,
			"management_ip":       types.StringType,
			"management_username": types.StringType,
			"management_password": types.StringType,
		}, map[string]attr.Value{
			"ip":                  types.StringValue(""),
			"ssh_username":        types.StringValue(""),
			"ssh_password":        types.StringValue(""),
			"ssh_port":            types.StringValue(""),
			"management_ip":       types.StringValue(""),
			"management_username": types.StringValue(""),
			"management_password": types.StringValue(""),
		})
		return val
	}
	return info
}

// CreatePeerSystem POST peer system
func CreatePeerSystem(client *goscaleio.Client, plan models.PeerMdmResourceModel) (string, error) {
	var ipList []string
	for _, val := range plan.IPList {
		ipList = append(ipList, val.ValueString())
	}
	value, err := client.AddPeerMdm(&scaleiotypes.AddPeerMdm{
		PeerSystemID:  plan.PeerSystemID.ValueString(),
		Port:          fmt.Sprint(plan.Port.ValueInt64()),
		Name:          plan.Name.ValueString(),
		PeerSystemIps: ipList,
	})

	if err != nil {
		return "", err
	}
	return value.ID, err
}

// PeerSystemUpdate update the different values of the peer system
func PeerSystemUpdate(client *goscaleio.Client, state models.PeerMdmResourceModel, plan models.PeerMdmResourceModel) error {

	// Rename
	if state.Name.ValueString() != plan.Name.ValueString() {
		errRename := client.ModifyPeerMdmName(state.ID.ValueString(), &scaleiotypes.ModifyPeerMDMNameParam{
			NewName: plan.Name.ValueString(),
		})
		if errRename != nil {
			return errRename
		}
	}

	// Update IP List
	if !reflect.DeepEqual(state.IPList, plan.IPList) {
		var localIPList []string
		for _, ip := range plan.IPList {
			localIPList = append(localIPList, ip.ValueString())
		}

		errIP := client.ModifyPeerMdmIP(state.ID.ValueString(), localIPList)
		if errIP != nil {
			return errIP
		}
	}

	// Update Port
	if state.Port != plan.Port {
		errPort := client.ModifyPeerMdmPort(state.ID.ValueString(), &scaleiotypes.ModifyPeerMDMPortParam{
			NewPort: fmt.Sprint(plan.Port.ValueInt64()),
		})
		if errPort != nil {
			return errPort
		}
	}

	if state.PerfProfile != plan.PerfProfile {
		errPerf := client.ModifyPeerMdmPerformanceParameters(state.ID.ValueString(), &scaleiotypes.ModifyPeerMdmPerformanceParametersParam{
			NewPreformanceProfile: plan.PerfProfile.ValueString(),
		})
		if errPerf != nil {
			return errPerf
		}
	}
	return nil
}

// AddCertificate POST peer system
func AddCertificate(ctx context.Context, client *goscaleio.Client, plan models.PeerMdmResourceModel) error {
	var source models.PrimaryMdmInfo
	var destination models.PrimaryMdmInfo
	plan.DestinationPrimaryMdmInfo.As(ctx, &destination, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	plan.SourcePrimaryMdmInfo.As(ctx, &source, basetypes.ObjectAsOptions{UnhandledNullAsEmpty: true, UnhandledUnknownAsEmpty: true})
	// Create the destination ssh client
	destProv, createSSHClientDestErr := sshClient.NewSSHProvisioner(sshClient.SSHProvisionerConfig{
		Port:     destination.Port,
		IP:       destination.IP,
		Username: destination.Username,
		Password: destination.Password,
	}, &provisionerLogger{ctx: ctx})
	// Error creating destination ssh client
	if createSSHClientDestErr != nil {
		return createSSHClientDestErr
	}
	defer destProv.Close()

	// create the source ssh client
	sourceProv, createSSHClientSourceErr := sshClient.NewSSHProvisioner(sshClient.SSHProvisionerConfig{
		Port:     source.Port,
		IP:       source.IP,
		Username: source.Username,
		Password: source.Password,
	}, &provisionerLogger{ctx: ctx})
	// Error creating source ssh client
	if createSSHClientSourceErr != nil {
		return createSSHClientSourceErr
	}
	defer sourceProv.Close()

	// Grab the certificate from the  mdm
	// Commands to get the certificate
	// 1. Add local cert to login
	// 2. Login to mdm
	// 3. Extract root certificate
	// 4. Copy destination cert to source mdm

	// 1. Add local cert to login
	_, addLocalCertErr := destProv.Run("scli --add_certificate --certificate_file /opt/emc/scaleio/mdm/cfg/mgmt_ca.pem")
	if addLocalCertErr != nil {
		return addLocalCertErr
	}

	// 2. Login to mdm
	_, loginErr := destProv.Run(`scli --login --management_system_ip ` + destination.ManagementIP + ` --username ` + destination.ManagementUsername + ` --password ` + *destination.ManagementPassword)
	if loginErr != nil {
		return fmt.Errorf("Unable to login: %s", loginErr.Error())
	}

	// 3. Extract root cert
	_, extractRootCertErr := destProv.Run(`scli --extract_root_ca --certificate_file ` + destination.ManagementIP + `_root.pem`)
	if extractRootCertErr != nil {
		return fmt.Errorf("Unable to extract root certificate: %s", extractRootCertErr.Error())
	}

	// 4. Copy destination cert to source mdm Example: curl --insecure --user user:pass -T ip_root.pem sftp://ip/root/ip_root.pem
	_, scpRootCertErr := destProv.Run(`curl --insecure --user ` + source.Username + `:` + *source.Password + ` -T ` + destination.ManagementIP + `_root.pem sftp://` + source.IP + `/root/` + destination.ManagementIP + `_root.pem`)
	if scpRootCertErr != nil {
		return fmt.Errorf("Unable to copy from certificate from source to destination: %s", scpRootCertErr.Error())
	}

	// Go to Source MDM and add the new cert
	// 1. Add local cert to login
	// 2. Login to system
	// 3. Add the new cert from the destination as a trusted cert

	// 1. Add local cert to login
	_, addLocalCertSourceErr := sourceProv.Run("scli --add_certificate --certificate_file /opt/emc/scaleio/mdm/cfg/mgmt_ca.pem")
	if addLocalCertSourceErr != nil {
		return addLocalCertSourceErr
	}
	// 2. Login to mdm
	_, loginSourceErr := sourceProv.Run(`scli --login --management_system_ip ` + source.ManagementIP + ` --username ` + source.ManagementUsername + ` --password ` + *source.ManagementPassword)
	if loginSourceErr != nil {
		return fmt.Errorf("Unable to login: %s", loginSourceErr.Error())
	}
	// 3. Add the cert of the destination as a trusted cert
	_, addTrustedCertErr := sourceProv.Run(`scli --add_trusted_ca --certificate_file ` + destination.ManagementIP + `_root.pem --comment "Adding Peer Mdm ` + destination.ManagementIP + `_root.pem"`)

	// If the error state is equal to 7 that means the certificate is already trusted
	// In this case we skip the error and just continue
	if addTrustedCertErr != nil {
		err := fmt.Sprintf("%s", addTrustedCertErr)
		// If error string contains exit code 7 that means the certificate is already trusted and can be ignored
		if !strings.Contains(err, "7") {
			return fmt.Errorf("Unable to add trusted cert: %s", err)

		}
	}
	return nil
}

// GetReplicationPairs GET replication pairs
func GetReplicationPairs(client *goscaleio.Client) ([]scaleiotypes.ReplicationPair, error) {
	rps := []scaleiotypes.ReplicationPair{}

	// Get All Replication Pairs
	pairs, err := client.GetAllReplicationPairs()
	if err != nil {
		return nil, err
	}
	for _, val := range pairs {
		rps = append(rps, *val)
	}

	return rps, nil
}

// CreateReplicationPair POST replication pair
func CreateReplicationPair(client *goscaleio.Client, plan models.ReplicationPairResourceModel) (string, error) {
	rp := &scaleiotypes.QueryReplicationPair{
		Name:                          plan.Name.ValueString(),
		SourceVolumeID:                plan.SourceVolumeID.ValueString(),
		DestinationVolumeID:           plan.DestinationVolumeID.ValueString(),
		ReplicationConsistencyGroupID: plan.ReplicationConsistencyGroupID.ValueString(),
		// OnlineCopy is the only supported copy type for replication pair
		CopyType: "OnlineCopy",
	}
	res, err := client.CreateReplicationPair(rp)
	if err != nil {
		return "", err
	}
	return res.ID, err
}

// PauseReplicationPair Pause initial replication pair
func PauseReplicationPair(client *goscaleio.Client, id string) (*scaleiotypes.ReplicationPair, error) {
	return client.PausePairInitialCopy(id)
}

// ResumeReplicationPair Resume initial replication pair
func ResumeReplicationPair(client *goscaleio.Client, id string) (*scaleiotypes.ReplicationPair, error) {
	return client.ResumePairInitialCopy(id)
}

// GetSpecificReplicationPair GET a replication pair
func GetSpecificReplicationPair(client *goscaleio.Client, id string) (*scaleiotypes.ReplicationPair, error) {
	return client.GetReplicationPair(id)
}

// MapReplicationPairState map single replication pair state
func MapReplicationPairState(val scaleiotypes.ReplicationPair, state models.ReplicationPairResourceModel) models.ReplicationPairResourceModel {
	state.ID = types.StringValue(val.ID)
	state.Name = types.StringValue(val.Name)
	state.RemoteID = types.StringValue(val.RemoteID)
	state.UserRequestedPauseTransmitInitCopy = types.BoolValue(val.UserRequestedPauseTransmitInitCopy)
	state.RemoteCapacityInMB = types.Int64Value(int64(val.RemoteCapacityInMB))
	state.LocalVolumeID = types.StringValue(val.LocalVolumeID)
	state.RemoteVolumeID = types.StringValue(val.RemoteID)
	state.RemoteVolumeName = types.StringValue(val.RemoteVolumeName)
	state.ReplicationConsistencyGroupID = types.StringValue(val.ReplicationConsistencyGroupID)
	state.CopyType = types.StringValue(val.LifetimeState)
	state.LifetimeState = types.StringValue(val.CopyType)
	state.PeerSystemName = types.StringValue(val.LifetimeState)
	state.InitialCopyState = types.StringValue(val.InitialCopyState)
	state.InitialCopyPriority = types.Int64Value(int64(val.InitialCopyPriority))
	return state
}

// MapReplicationPairsState map replication pairs state
func MapReplicationPairsState(pairs []scaleiotypes.ReplicationPair, state models.ReplicationPairDataSourceModel) models.ReplicationPairDataSourceModel {
	mappedRps := []models.ReplicationPairModel{}
	for _, val := range pairs {
		temp := models.ReplicationPairModel{
			ID:                                 types.StringValue(val.ID),
			Name:                               types.StringValue(val.Name),
			RemoteID:                           types.StringValue(val.RemoteID),
			UserRequestedPauseTransmitInitCopy: types.BoolValue(val.UserRequestedPauseTransmitInitCopy),
			RemoteCapacityInMB:                 types.Int64Value(int64(val.RemoteCapacityInMB)),
			LocalVolumeID:                      types.StringValue(val.LocalVolumeID),
			RemoteVolumeID:                     types.StringValue(val.RemoteID),
			RemoteVolumeName:                   types.StringValue(val.RemoteVolumeName),
			ReplicationConsistencyGroupID:      types.StringValue(val.ReplicationConsistencyGroupID),
			CopyType:                           types.StringValue(val.LifetimeState),
			LifetimeState:                      types.StringValue(val.CopyType),
			PeerSystemName:                     types.StringValue(val.LifetimeState),
			InitialCopyState:                   types.StringValue(val.InitialCopyState),
			InitialCopyPriority:                types.Int64Value(int64(val.InitialCopyPriority)),
		}
		mappedRps = append(mappedRps, temp)
	}
	return models.ReplicationPairDataSourceModel{
		ID:                     types.StringValue("replication_pair_id"),
		ReplicationPairFilter:  state.ReplicationPairFilter,
		ReplicationPairDetails: mappedRps,
	}
}

// GetReplicationConsistancyGroups GET RCGs
func GetReplicationConsistancyGroups(client *goscaleio.Client) ([]scaleiotypes.ReplicationConsistencyGroup, error) {
	rps := []scaleiotypes.ReplicationConsistencyGroup{}

	// Get All RCGs
	rcgs, err := client.GetReplicationConsistencyGroups()
	if err != nil {
		return nil, err
	}
	for _, val := range rcgs {
		rps = append(rps, *val)
	}

	return rps, nil
}

// GetSpecificReplicationConsistencyGroup GET a specific RCG
func GetSpecificReplicationConsistencyGroup(client *goscaleio.Client, id string) (*scaleiotypes.ReplicationConsistencyGroup, error) {
	return client.GetReplicationConsistencyGroupByID(id)
}

// CreateReplicationConsistencyGroup POST replication consistency group
func CreateReplicationConsistencyGroup(client *goscaleio.Client, plan models.ReplicationConsistancyGroupModel) (string, error) {
	rcg := scaleiotypes.ReplicationConsistencyGroupCreatePayload{
		Name:                     plan.Name.ValueString(),
		RpoInSeconds:             plan.RpoInSeconds.String(),
		ProtectionDomainID:       plan.ProtectionDomainID.ValueString(),
		RemoteProtectionDomainID: plan.RemoteProtectionDomainID.ValueString(),
		DestinationSystemID:      plan.DestinationSystemID.ValueString(),
	}
	res, err := client.CreateReplicationConsistencyGroup(&rcg)
	if err != nil {
		return "", err
	}
	return res.ID, err
}

// MapReplicationConsistancyGroupsState map Replication Consistancy Groups state
func MapReplicationConsistancyGroupsState(rcgs []scaleiotypes.ReplicationConsistencyGroup, state models.ReplicationConsistancyGroupDataSourceModel) models.ReplicationConsistancyGroupDataSourceModel {
	mappedRps := []models.ReplicationConsistancyGroupModel{}
	for _, val := range rcgs {
		temp := models.ReplicationConsistancyGroupModel{
			ID:                          types.StringValue(val.ID),
			Name:                        types.StringValue(val.Name),
			RemoteID:                    types.StringValue(val.RemoteID),
			RpoInSeconds:                types.Int64Value(int64(val.RpoInSeconds)),
			ProtectionDomainID:          types.StringValue(val.ProtectionDomainID),
			RemoteProtectionDomainID:    types.StringValue(val.RemoteProtectionDomainID),
			DestinationSystemID:         types.StringValue(val.DestinationSystemID),
			PeerMdmID:                   types.StringValue(val.PeerMdmID),
			RemoteMdmID:                 types.StringValue(val.RemoteMdmID),
			ReplicationDirection:        types.StringValue(val.ReplicationDirection),
			CurrConsistMode:             types.StringValue(val.CurrConsistMode),
			FreezeState:                 types.StringValue(val.FreezeState),
			PauseMode:                   types.StringValue(val.PauseMode),
			LifetimeState:               types.StringValue(val.LifetimeState),
			SnapCreationInProgress:      types.BoolValue(val.SnapCreationInProgress),
			LastSnapGroupID:             types.StringValue(val.LastSnapGroupID),
			Type:                        types.StringValue(val.Type),
			DisasterRecoveryState:       types.StringValue(val.DisasterRecoveryState),
			RemoteDisasterRecoveryState: types.StringValue(val.RemoteDisasterRecoveryState),
			TargetVolumeAccessMode:      types.StringValue(val.TargetVolumeAccessMode),
			FailoverType:                types.StringValue(val.FailoverType),
			FailoverState:               types.StringValue(val.FailoverState),
			ActiveLocal:                 types.BoolValue(val.ActiveLocal),
			ActiveRemote:                types.BoolValue(val.ActiveRemote),
			AbstractState:               types.StringValue(val.AbstractState),
			Error:                       types.Int64Value(int64(val.Error)),
			LocalActivityState:          types.StringValue(val.LocalActivityState),
			RemoteActivityState:         types.StringValue(val.RemoteActivityState),
			InactiveReason:              types.Int64Value(int64(val.InactiveReason)),
		}
		mappedRps = append(mappedRps, temp)
	}
	return models.ReplicationConsistancyGroupDataSourceModel{
		ID:                                 types.StringValue("replication_consistancy_group_id"),
		ReplicationConsistancyGroupFilter:  state.ReplicationConsistancyGroupFilter,
		ReplicationConsistancyGroupDetails: mappedRps,
	}
}

// MapReplicationConsistancyGroupsResourceState map Replication Consistancy Groups state
func MapReplicationConsistancyGroupsResourceState(rcg scaleiotypes.ReplicationConsistencyGroup, state models.ReplicationConsistancyGroupModel) models.ReplicationConsistancyGroupModel {
	rcgMap := models.ReplicationConsistancyGroupModel{
		ID:                          types.StringValue(rcg.ID),
		Name:                        types.StringValue(rcg.Name),
		RemoteID:                    types.StringValue(rcg.RemoteID),
		RpoInSeconds:                state.RpoInSeconds,
		ProtectionDomainID:          types.StringValue(rcg.ProtectionDomainID),
		RemoteProtectionDomainID:    types.StringValue(rcg.RemoteProtectionDomainID),
		DestinationSystemID:         state.DestinationSystemID,
		PeerMdmID:                   types.StringValue(rcg.PeerMdmID),
		RemoteMdmID:                 types.StringValue(rcg.RemoteMdmID),
		ReplicationDirection:        types.StringValue(rcg.ReplicationDirection),
		CurrConsistMode:             state.CurrConsistMode,
		FreezeState:                 state.FreezeState,
		PauseMode:                   state.PauseMode,
		LifetimeState:               types.StringValue(rcg.LifetimeState),
		SnapCreationInProgress:      types.BoolValue(rcg.SnapCreationInProgress),
		LastSnapGroupID:             types.StringValue(rcg.LastSnapGroupID),
		Type:                        types.StringValue(rcg.Type),
		DisasterRecoveryState:       types.StringValue(rcg.DisasterRecoveryState),
		RemoteDisasterRecoveryState: types.StringValue(rcg.RemoteDisasterRecoveryState),
		TargetVolumeAccessMode:      state.TargetVolumeAccessMode,
		FailoverType:                types.StringValue(rcg.FailoverType),
		FailoverState:               types.StringValue(rcg.FailoverState),
		ActiveLocal:                 types.BoolValue(rcg.ActiveLocal),
		ActiveRemote:                types.BoolValue(rcg.ActiveRemote),
		AbstractState:               types.StringValue(rcg.AbstractState),
		Error:                       types.Int64Value(int64(rcg.Error)),
		LocalActivityState:          state.LocalActivityState,
		RemoteActivityState:         types.StringValue(rcg.RemoteActivityState),
		InactiveReason:              types.Int64Value(int64(rcg.InactiveReason)),
	}
	return rcgMap
}

// RCGUpdates Update the RCG
func RCGUpdates(client *goscaleio.Client, state models.ReplicationConsistancyGroupModel, plan models.ReplicationConsistancyGroupModel) error {
	rcgClient := goscaleio.NewReplicationConsistencyGroup(client)
	rcgClient.ReplicationConsistencyGroup.ID = state.ID.ValueString()
	// Update RPO
	if state.RpoInSeconds.ValueInt64() != plan.RpoInSeconds.ValueInt64() {
		rpoErr := rcgClient.SetRPOOnReplicationGroup(scaleiotypes.SetRPOReplicationConsistencyGroup{
			RpoInSeconds: fmt.Sprint(plan.RpoInSeconds.ValueInt64()),
		})
		if rpoErr != nil {
			return rpoErr
		}
	}

	// Update Activity State
	if state.LocalActivityState.ValueString() != plan.LocalActivityState.ValueString() {
		if plan.LocalActivityState.ValueString() == "Active" {
			activateErr := rcgClient.ExecuteActivateOnReplicationGroup()
			if activateErr != nil {
				return activateErr
			}
		} else {
			inactivateErr := rcgClient.ExecuteTerminateOnReplicationGroup()
			if inactivateErr != nil {
				return inactivateErr
			}
		}

	}

	// Update Access Mode
	if state.TargetVolumeAccessMode.ValueString() != plan.TargetVolumeAccessMode.ValueString() {
		vamErr := rcgClient.SetTargetVolumeAccessModeOnReplicationGroup(scaleiotypes.SetTargetVolumeAccessModeOnReplicationGroup{
			TargetVolumeAccessMode: plan.TargetVolumeAccessMode.ValueString(),
		})
		if vamErr != nil {
			return vamErr
		}
	}

	// Update Pause Mode
	if state.PauseMode.ValueString() != plan.PauseMode.ValueString() {
		if plan.PauseMode.ValueString() == "Pause" {
			pauseModeErr := rcgClient.ExecutePauseOnReplicationGroup()
			if pauseModeErr != nil {
				return pauseModeErr
			}
		} else {
			resumeModeErr := rcgClient.ExecuteResumeOnReplicationGroup()
			if resumeModeErr != nil {
				return resumeModeErr
			}
		}
	}

	// Update Freeze State
	if state.FreezeState.ValueString() != plan.FreezeState.ValueString() {
		if plan.FreezeState.ValueString() == "Frozen" {
			freezeErr := rcgClient.FreezeReplicationConsistencyGroup(state.ID.ValueString())
			if freezeErr != nil {
				return freezeErr
			}
		} else {
			unfreezeErr := rcgClient.UnfreezeReplicationConsistencyGroup()
			if unfreezeErr != nil {
				return unfreezeErr
			}
		}
	}

	// Update Consistency Mode
	if plan.CurrConsistMode.ValueString() != state.CurrConsistMode.ValueString() {
		if plan.CurrConsistMode.ValueString() == "Consistent" {
			consistentErr := rcgClient.ExecuteConsistentOnReplicationGroup()
			if consistentErr != nil {
				return consistentErr
			}
		} else {
			inconsistentErr := rcgClient.ExecuteInconsistentOnReplicationGroup()
			if inconsistentErr != nil {
				return inconsistentErr
			}
		}
	}

	// Update Name
	if state.Name.ValueString() != plan.Name.ValueString() {
		nameErr := rcgClient.SetNewNameOnReplicationGroup(scaleiotypes.SetNewNameOnReplicationGroup{
			NewName: plan.Name.ValueString(),
		})
		if nameErr != nil {
			return nameErr
		}
	}
	return nil
}
