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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// MdmResourceModel defines the struct for MDM resource
type MdmResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	PerformanceProfile types.String `tfsdk:"performance_profile"`
	ClusterMode        types.String `tfsdk:"cluster_mode"`
	PrimaryMdm         Mdm          `tfsdk:"primary_mdm"`
	SecondaryMdm       []Mdm        `tfsdk:"secondary_mdm"`
	TieBreakerMdm      []Mdm        `tfsdk:"tiebreaker_mdm"`
	StandByMdm         types.List   `tfsdk:"standby_mdm"`
}

// Mdm defines the struct for Primary, Secondary and TB MDM
type Mdm struct {
	ID            types.String `tfsdk:"id"`
	Name          types.String `tfsdk:"name"`
	Port          types.Int64  `tfsdk:"port"`
	Ips           types.Set    `tfsdk:"ips"`
	ManagementIps types.Set    `tfsdk:"management_ips"`
}

// StandByMdm defines the struct for StandBy MDM
type StandByMdm struct {
	ID                 types.String `tfsdk:"id"`
	Name               types.String `tfsdk:"name"`
	Port               types.Int64  `tfsdk:"port"`
	Ips                types.Set    `tfsdk:"ips"`
	ManagementIps      types.Set    `tfsdk:"management_ips"`
	Role               types.String `tfsdk:"role"`
	AllowAsymmetricIps types.Bool   `tfsdk:"allow_asymmetric_ips"`
}
