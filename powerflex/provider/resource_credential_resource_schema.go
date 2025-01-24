/*
Copyright (c) 2025 Dell Inc., or its subsidiaries. All Rights Reserved.

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
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// ResourceCredentialResourceSchema defines the schema for ResourceCredential resource
var ResourceCredentialResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource is used to manage Resource Credential entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above. We can Create, Update and Delete the PowerFlex Resource Credentials using this resource. We can also Import an existing Resource Credentials from the PowerFlex array.",
	MarkdownDescription: "This resource is used to manage Resource Credential entity of the PowerFlex Array. This feature is only supported for PowerFlex 4.5 and above. We can Create, Update and Delete the PowerFlex Resource Credentials using this resource. We can also Import an existing Resource Credentials from the PowerFlex array.",
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description:         "Resource Credential name",
			MarkdownDescription: "Resource Credential name",
			Required:            true,
		},
		"username": schema.StringAttribute{
			Description:         "Resource Credential username",
			MarkdownDescription: "Resource Credential username",
			Required:            true,
		},
		"password": schema.StringAttribute{
			Description:         "Resource Credential password",
			MarkdownDescription: "Resource Credential password",
			Required:            true,
			Sensitive:           true,
		},
		"type": schema.StringAttribute{
			Description:         "Resource Credential type",
			MarkdownDescription: "Resource Credential type",
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"Node",
				"Switch",
				"vCenter",
				"ElementManager",
				"PowerflexGateway",
				"PresentationServer",
				"OSAdmin",
				"OSUser",
			)},
			Required: true,
		},
		"domain": schema.StringAttribute{
			Description:         "Resource Credential domain, used by types OSUser, ElementManager and vCenterCredential",
			MarkdownDescription: "Resource Credential domain, used by types OSUser, ElementManager and vCenterCredential",
			Optional:            true,
		},
		"snmp_v2_community_string": schema.StringAttribute{
			Description:         "Resource Credential snmp_v2_community_string, used by types Node, Switch and ElementManager. If unset it will default to public",
			MarkdownDescription: "Resource Credential snmp_v2_community_string, used by types Node, Switch and ElementManager. If unset it will default to public",
			Optional:            true,
		},
		"snmp_v3_security_level": schema.StringAttribute{
			Description:         "Resource Credential snmp_v3_security_level, used by types Node. Minimal requires only snmp_v3_security_name, Moderate requires snmp_v3_security_name and snmp_v3_md5_authentication_password, Maximal requires snmp_v3_security_name, snmp_v3_md5_authentication_password and snmp_v3_des_authentication_password",
			MarkdownDescription: "Resource Credential snmp_v3_security_level, used by types Node. Minimal requires only snmp_v3_security_name, Moderate requires snmp_v3_security_name and snmp_v3_md5_authentication_password, Maximal requires snmp_v3_security_name, snmp_v3_md5_authentication_password and snmp_v3_des_authentication_password",
			Optional:            true,
			Validators: []validator.String{stringvalidator.OneOfCaseInsensitive(
				"Minimal",
				"Moderate",
				"Maximal",
			)},
		},
		"snmp_v3_security_name": schema.StringAttribute{
			Description:         "Resource Credential snmp_v3_security_name, used by types Node",
			MarkdownDescription: "Resource Credential snmp_v3_security_name, used by types Node",
			Optional:            true,
		},
		"snmp_v3_md5_authentication_password": schema.StringAttribute{
			Description:         "Resource Credential snmp_v3_md5_authentication_password, used by types Node when snmp_v3_security_level is set to Moderate",
			MarkdownDescription: "Resource Credential snmp_v3_md5_authentication_password, used by types Node when snmp_v3_security_level is set to Moderate",
			Optional:            true,
			Sensitive:           true,
		},
		"snmp_v3_des_authentication_password": schema.StringAttribute{
			Description:         "Resource Credential snmp_v3_des_authentication_password, used by types Node when snmp_v3_security_level is set to Maximal",
			MarkdownDescription: "Resource Credential snmp_v3_des_authentication_password, used by types Node when snmp_v3_security_level is set to Maximal",
			Optional:            true,
			Sensitive:           true,
		},
		"ssh_private_key": schema.StringAttribute{
			Description:         "Resource Credential ssh_private_key, can be used by types Node, Switch, OSAdminCrednetial and OSUserCredential",
			MarkdownDescription: "Resource Credential ssh_private_key, can be used by types Node, Switch, OSAdminCrednetial and OSUserCredential",
			Optional:            true,
			Sensitive:           true,
		},
		"key_pair_name": schema.StringAttribute{
			Description:         "Resource Credential key_pair_name, can be used by types Node, Switch, OSAdminCrednetial and OSUserCredential. Is Required if ssh_private_key is set",
			MarkdownDescription: "Resource Credential key_pair_name, can be used by types Node, Switch, OSAdminCrednetial and OSUserCredential. Is Required if ssh_private_key is set",
			Optional:            true,
		},
		"os_username": schema.StringAttribute{
			Description:         "Resource Credential os_username, used by types PowerflexGateway",
			MarkdownDescription: "Resource Credential os_username, used by types PowerflexGateway",
			Optional:            true,
		},
		"os_password": schema.StringAttribute{
			Description:         "Resource Credential os_password, used by types PowerflexGateway",
			MarkdownDescription: "Resource Credential os_password, used by types PowerflexGateway",
			Optional:            true,
			Sensitive:           true,
		},
		"id": schema.StringAttribute{
			Description:         "Unique identifier of the resource credential instance.",
			MarkdownDescription: "Unique identifier of the resource credential instance.",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			Description:         "Resource Credential created_date",
			MarkdownDescription: "Resource Credential created_date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			Description:         "Resource Credential created_by",
			MarkdownDescription: "Resource Credential created_by",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			Description:         "Resource Credential updated_by",
			MarkdownDescription: "Resource Credential updated_by",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			Description:         "Resource Credential updated_date",
			MarkdownDescription: "Resource Credential updated_date",
			Computed:            true,
		},
	},
}
