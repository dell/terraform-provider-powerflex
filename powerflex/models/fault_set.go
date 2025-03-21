/*
Copyright (c) 2023-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

// FaultSetResourceModel defines struct fault set resource
type FaultSetResourceModel struct {
	ID                 types.String `tfsdk:"id"`
	ProtectionDomainID types.String `tfsdk:"protection_domain_id"`
	Name               types.String `tfsdk:"name"`
}

// FaultSetDataSourceModel maps the struct to FaultSet data source schema
type FaultSetDataSourceModel struct {
	FaultSetFilter  *FaultSetFilter `tfsdk:"filter"`
	FaultSetDetails []FaultSet      `tfsdk:"fault_set_details"`
	ID              types.String    `tfsdk:"id"`
}

// FaultSetFilter defines the filter for fault set
type FaultSetFilter struct {
	ProtectionDomainID []types.String `tfsdk:"protection_domain_id"`
	Name               []types.String `tfsdk:"name"`
	ID                 []types.String `tfsdk:"id"`
}

// FaultSet maps the struct to FaultSet schema
type FaultSet struct {
	ProtectionDomainID string         `tfsdk:"protection_domain_id"`
	Name               string         `tfsdk:"name"`
	ID                 string         `tfsdk:"id"`
	Links              []*LinkModel   `tfsdk:"links"`
	SdsDetails         []SdsDataModel `tfsdk:"sds_details"`
}
