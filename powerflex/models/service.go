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

package models

import (
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// ServiceResourceModel is the tfsdk model of ServiceResponse
type ServiceResourceModel struct {
	ID                    types.String `tfsdk:"id"`
	DeploymentName        types.String `tfsdk:"deployment_name"`
	DeploymentDescription types.String `tfsdk:"deployment_description"`
	TemplateID            types.String `tfsdk:"template_id"`
	FirmwareID            types.String `tfsdk:"firmware_id"`
	Nodes                 types.Int64  `tfsdk:"nodes"`
	Status                types.String `tfsdk:"status"`
	Compliant             types.Bool   `tfsdk:"compliant"`
	TemplateName          types.String `tfsdk:"template_name"`
	DeploymentTimeout     types.Int64  `tfsdk:"deployment_timeout"`
}
