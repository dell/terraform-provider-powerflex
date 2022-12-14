package helper

import (
	goscaleiotypes "github.com/dell/goscaleio/types/v1"
)

// SdcToMap function for Convert returned sdc object(from goscaleio) to terraform output schema
func SdcToMap(s goscaleiotypes.Sdc) map[string]interface{} {
	resultSDC := make(map[string]interface{})
	resultSDC["id"] = s.ID
	resultSDC["name"] = s.Name
	resultSDC["sdcguid"] = s.SdcGUID
	resultSDC["sdcapproved"] = s.SdcApproved
	resultSDC["onvmware"] = s.OnVMWare
	resultSDC["systemid"] = s.SystemID
	resultSDC["sdcip"] = s.SdcIP
	resultSDC["mdmconnectionstate"] = s.MdmConnectionState
	return resultSDC
}

// NameChangedSdcToMap function to convert returned sdc object(from goscaleio) [after name changes] to terraform output schema
func NameChangedSdcToMap(s goscaleiotypes.Sdc) map[string]interface{} {
	resultSDC := make(map[string]interface{})
	resultSDC["id"] = s.ID
	resultSDC["name"] = s.Name
	return resultSDC
}
