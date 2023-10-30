package helper

import (
	"terraform-provider-powerflex/powerflex/models"

	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// UpdateFaultSetState updates the State for Fault set Resource
func UpdateFaultSetState(faultset *scaleiotypes.FaultSet, plan models.FaultSetResourceModel) models.FaultSetResourceModel {
	state := plan
	state.ProtectionDomainID = types.StringValue(faultset.ProtectionDomainId)
	state.ID = types.StringValue(faultset.ID)
	state.Name = types.StringValue(faultset.Name)
	return state
}
