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

package helper

import (
	"context"
	"fmt"

	"bytes"
	"encoding/json"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// GetFirstSystem - finds available first system and returns it.
func GetFirstSystem(rc *goscaleio.Client) (*goscaleio.System, error) {
	allSystems, err := rc.GetSystems()
	if err != nil {
		return nil, fmt.Errorf("Error in goscaleio GetSystems")
	}
	if numSys := len((allSystems)); numSys == 0 {
		return nil, fmt.Errorf("no systems found")
	} else if numSys > 1 {
		return nil, fmt.Errorf("more than one system found")
	}
	system, err := rc.FindSystem(allSystems[0].ID, "", "")
	if err != nil {
		return nil, fmt.Errorf("Error in goscaleio FindSystem")
	}
	return system, nil
}

// PrettyJSON - function for logging json readable output.
func PrettyJSON(data interface{}) string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return ""
	}
	return buffer.String()
}

// GetNewProtectionDomainEx function to get Protection Domain
func GetNewProtectionDomainEx(c *goscaleio.Client, pdID string, pdName string, href string) (*goscaleio.ProtectionDomain, error) {
	system, err := GetFirstSystem(c)
	if err != nil {
		return nil, err
	}
	pdr := goscaleio.NewProtectionDomainEx(c, &scaleiotypes.ProtectionDomain{})
	if pdID != "" {
		protectionDomain, err := system.FindProtectionDomain(pdID, "", "")
		pdr.ProtectionDomain = protectionDomain
		if err != nil {
			return nil, err
		}
	} else {
		protectionDomain, err := system.FindProtectionDomain("", pdName, "")
		pdr.ProtectionDomain = protectionDomain
		if err != nil {
			return nil, err
		}
	}
	return pdr, nil
}

// GetStoragePoolType returns storage pool type
func GetStoragePoolType(r *goscaleio.Client, storagePoolID string) (*goscaleio.StoragePool, error) {
	system, err := GetFirstSystem(r)
	if err != nil {
		return nil, err
	}

	sp, err := system.GetStoragePoolByID(storagePoolID)
	if err != nil {
		return nil, err
	}

	sp1 := goscaleio.NewStoragePoolEx(r, sp)
	return sp1, nil
}

// GetSdcType function returns SDC type
func GetSdcType(c *goscaleio.Client, sdcID string) (*goscaleio.Sdc, error) {
	system, err := GetFirstSystem(c)
	if err != nil {
		return nil, err
	}
	return system.GetSdcByID(sdcID)
}

// GetVolumeType function returns volume type
func GetVolumeType(c *goscaleio.Client, volID string) (*goscaleio.Volume, error) {
	volumes, err := c.GetVolume("", volID, "", "", false)
	if err != nil {
		return nil, err
	}

	volume := volumes[0]
	volType := goscaleio.NewVolume(c)
	volType.Volume = volume
	return volType, nil
}

// StringDefaultModifier is a plan modifier that sets a default value for a
// types.StringType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type StringDefaultModifier struct {
	Default string
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m StringDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m StringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
}

// PlanModifyString runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m StringDefaultModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the value is unknown or known, do not set default value.
	if req.PlanValue.IsNull() {
		resp.PlanValue = types.StringValue(m.Default)
	}
	if req.PlanValue.IsUnknown() {
		resp.PlanValue = types.StringValue(m.Default)
	}
}

// StringDefault sets default value fot string attributes
func StringDefault(defaultValue string) planmodifier.String {
	return StringDefaultModifier{
		Default: defaultValue,
	}
}

// boolDefaultModifier is a plan modifier that sets a default value for a
// types.BoolType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type boolDefaultModifier struct {
	Default bool
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m boolDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %t", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m boolDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%t`", m.Default)
}

// PlanModifyBool runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m boolDefaultModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// If the value is unknown or known, do not set default value.
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		resp.PlanValue = types.BoolValue(m.Default)
	}
}

// BoolDefault sets default value fot string attributes
func BoolDefault(defaultValue bool) planmodifier.Bool {
	return boolDefaultModifier{
		Default: defaultValue,
	}
}

// ResetInstallerQueue function for the Abort, Clear and Move To Idle Execution
func ResetInstallerQueue(gatewayClient *goscaleio.GatewayClient) error {

	_, err := gatewayClient.AbortOperation()

	if err != nil {
		return fmt.Errorf("Error while Aborting Operation is %s", err.Error())
	}
	_, err = gatewayClient.ClearQueueCommand()

	if err != nil {
		return fmt.Errorf("Error while Clearing Queue is %s", err.Error())
	}

	_, err = gatewayClient.MoveToIdlePhase()

	if err != nil {
		return fmt.Errorf("Error while Move to Ideal Phase is %s", err.Error())
	}

	return nil
}

// CompareStringSlice Compare string slices. return true if the length and elements are same.
func CompareStringSlice(plan, state []string) bool {
	if len(plan) != len(state) {
		return false
	}

	itemAppearsTimes := make(map[string]int, len(plan))
	for _, i := range plan {
		itemAppearsTimes[i]++
	}

	for _, i := range state {
		if _, ok := itemAppearsTimes[i]; !ok {
			return false
		}

		itemAppearsTimes[i]--
		if itemAppearsTimes[i] == 0 {
			delete(itemAppearsTimes, i)
		}
	}
	return len(itemAppearsTimes) == 0
}

// CompareInt64Slice compares two slices of int64 and returns true if the length and elements are same.
func CompareInt64Slice(plan, state []int64) bool {

	if len(plan) != len(state) {
		return false
	}

	for i := range plan {
		if plan[i] != state[i] {
			return false
		}
	}
	return true
}
