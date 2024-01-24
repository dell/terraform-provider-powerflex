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

package provider

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TemplateDataSourceSchema defines the schema for template datasource
var TemplateDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing templates from PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	MarkdownDescription: "This datasource is used to query the existing templates from PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder attribute.",
			MarkdownDescription: "Placeholder attribute.",
			Computed:            true,
		},
		"template_ids": schema.SetAttribute{
			Description:         "List of template IDs",
			MarkdownDescription: "List of template IDs",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				setvalidator.ConflictsWith(
					path.MatchRoot("names"),
				),
			},
		},
		"names": schema.SetAttribute{
			Description:         "List of template names",
			MarkdownDescription: "List of template names",
			Optional:            true,
			ElementType:         types.StringType,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
				setvalidator.ValueStringsAre(stringvalidator.LengthAtLeast(1)),
				setvalidator.ConflictsWith(
					path.MatchRoot("template_ids"),
				),
			},
		},
		"template_details": schema.SetNestedAttribute{
			Description:         "Template details",
			MarkdownDescription: "Template details",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: TemplateDetailSchema()},
		},
	},
}

// TemplateDetailSchema is a function that returns the schema for TemplateDetails
func TemplateDetailSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"template_name": schema.StringAttribute{
			MarkdownDescription: "template_name",
			Description:         "template_name",
			Computed:            true,
		},
		"template_description": schema.StringAttribute{
			MarkdownDescription: "template_description",
			Description:         "template_description",
			Computed:            true,
		},
		"template_type": schema.StringAttribute{
			MarkdownDescription: "template_type",
			Description:         "template_type",
			Computed:            true,
		},
		"template_version": schema.StringAttribute{
			MarkdownDescription: "template_version",
			Description:         "template_version",
			Computed:            true,
		},
		"original_template_id": schema.StringAttribute{
			MarkdownDescription: "original_template_id",
			Description:         "original_template_id",
			Computed:            true,
		},
		"template_valid": schema.SingleNestedAttribute{
			MarkdownDescription: "template_valid",
			Description:         "template_valid",
			Computed:            true,
			Attributes:          TemplateValidSchema(),
		},
		"template_locked": schema.BoolAttribute{
			MarkdownDescription: "template_locked",
			Description:         "template_locked",
			Computed:            true,
		},
		"in_configuration": schema.BoolAttribute{
			MarkdownDescription: "in_configuration",
			Description:         "in_configuration",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "created_date",
			Description:         "created_date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "created_by",
			Description:         "created_by",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "updated_date",
			Description:         "updated_date",
			Computed:            true,
		},
		"last_deployed_date": schema.StringAttribute{
			MarkdownDescription: "last_deployed_date",
			Description:         "last_deployed_date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "updated_by",
			Description:         "updated_by",
			Computed:            true,
		},
		"manage_firmware": schema.BoolAttribute{
			MarkdownDescription: "manage_firmware",
			Description:         "manage_firmware",
			Computed:            true,
		},
		"use_default_catalog": schema.BoolAttribute{
			MarkdownDescription: "use_default_catalog",
			Description:         "use_default_catalog",
			Computed:            true,
		},
		"firmware_repository": schema.SingleNestedAttribute{
			MarkdownDescription: "firmware_repository",
			Description:         "firmware_repository",
			Computed:            true,
			Attributes:          FirmwareRepositorySchema(),
		},
		"license_repository": schema.SingleNestedAttribute{
			MarkdownDescription: "license_repository",
			Description:         "license_repository",
			Computed:            true,
			Attributes:          LicenseRepositorySchema(),
		},
		"assigned_users": schema.ListNestedAttribute{
			MarkdownDescription: "assigned_users",
			Description:         "assigned_users",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: AssignedUsersSchema()},
		},
		"all_users_allowed": schema.BoolAttribute{
			MarkdownDescription: "all_users_allowed",
			Description:         "all_users_allowed",
			Computed:            true,
		},
		"category": schema.StringAttribute{
			MarkdownDescription: "category",
			Description:         "category",
			Computed:            true,
		},
		"components": schema.ListNestedAttribute{
			MarkdownDescription: "components",
			Description:         "components",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ComponentsSchema()},
		},
		"configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "configuration",
			Description:         "configuration",
			Computed:            true,
			Attributes:          ConfigurationDetailsSchema(),
		},
		"server_count": schema.Int64Attribute{
			MarkdownDescription: "server_count",
			Description:         "server_count",
			Computed:            true,
		},
		"storage_count": schema.Int64Attribute{
			MarkdownDescription: "storage_count",
			Description:         "storage_count",
			Computed:            true,
		},
		"cluster_count": schema.Int64Attribute{
			MarkdownDescription: "cluster_count",
			Description:         "cluster_count",
			Computed:            true,
		},
		"service_count": schema.Int64Attribute{
			MarkdownDescription: "service_count",
			Description:         "service_count",
			Computed:            true,
		},
		"switch_count": schema.Int64Attribute{
			MarkdownDescription: "switch_count",
			Description:         "switch_count",
			Computed:            true,
		},
		"vm_count": schema.Int64Attribute{
			MarkdownDescription: "vm_count",
			Description:         "vm_count",
			Computed:            true,
		},
		"sdnas_count": schema.Int64Attribute{
			MarkdownDescription: "sdnas_count",
			Description:         "sdnas_count",
			Computed:            true,
		},
		"brownfield_template_type": schema.StringAttribute{
			MarkdownDescription: "brownfield_template_type",
			Description:         "brownfield_template_type",
			Computed:            true,
		},
		"networks": schema.ListNestedAttribute{
			MarkdownDescription: "networks",
			Description:         "networks",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: NetworksSchema()},
		},
		"draft": schema.BoolAttribute{
			MarkdownDescription: "draft",
			Description:         "draft",
			Computed:            true,
		},
	}
}

// MessagesSchema is a function that returns the schema for Messages
func MessagesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"message_code": schema.StringAttribute{
			MarkdownDescription: "message_code",
			Description:         "message_code",
			Computed:            true,
		},
		"message_bundle": schema.StringAttribute{
			MarkdownDescription: "message_bundle",
			Description:         "message_bundle",
			Computed:            true,
		},
		"severity": schema.StringAttribute{
			MarkdownDescription: "severity",
			Description:         "severity",
			Computed:            true,
		},
		"category": schema.StringAttribute{
			MarkdownDescription: "category",
			Description:         "category",
			Computed:            true,
		},
		"display_message": schema.StringAttribute{
			MarkdownDescription: "display_message",
			Description:         "display_message",
			Computed:            true,
		},
		"response_action": schema.StringAttribute{
			MarkdownDescription: "response_action",
			Description:         "response_action",
			Computed:            true,
		},
		"detailed_message": schema.StringAttribute{
			MarkdownDescription: "detailed_message",
			Description:         "detailed_message",
			Computed:            true,
		},
		"correlation_id": schema.StringAttribute{
			MarkdownDescription: "correlation_id",
			Description:         "correlation_id",
			Computed:            true,
		},
		"agent_id": schema.StringAttribute{
			MarkdownDescription: "agent_id",
			Description:         "agent_id",
			Computed:            true,
		},
		"time_stamp": schema.StringAttribute{
			MarkdownDescription: "time_stamp",
			Description:         "time_stamp",
			Computed:            true,
		},
		"sequence_number": schema.Int64Attribute{
			MarkdownDescription: "sequence_number",
			Description:         "sequence_number",
			Computed:            true,
		},
	}
}

// TemplateValidSchema is a function that returns the schema for TemplateValid
func TemplateValidSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"valid": schema.BoolAttribute{
			MarkdownDescription: "valid",
			Description:         "valid",
			Computed:            true,
		},
		"messages": schema.ListNestedAttribute{
			MarkdownDescription: "messages",
			Description:         "messages",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: MessagesSchema()},
		},
	}
}

// SoftwareComponentsSchema is a function that returns the schema for SoftwareComponents
func SoftwareComponentsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"package_id": schema.StringAttribute{
			MarkdownDescription: "package_id",
			Description:         "package_id",
			Computed:            true,
		},
		"dell_version": schema.StringAttribute{
			MarkdownDescription: "dell_version",
			Description:         "dell_version",
			Computed:            true,
		},
		"vendor_version": schema.StringAttribute{
			MarkdownDescription: "vendor_version",
			Description:         "vendor_version",
			Computed:            true,
		},
		"component_id": schema.StringAttribute{
			MarkdownDescription: "component_id",
			Description:         "component_id",
			Computed:            true,
		},
		"device_id": schema.StringAttribute{
			MarkdownDescription: "device_id",
			Description:         "device_id",
			Computed:            true,
		},
		"sub_device_id": schema.StringAttribute{
			MarkdownDescription: "sub_device_id",
			Description:         "sub_device_id",
			Computed:            true,
		},
		"vendor_id": schema.StringAttribute{
			MarkdownDescription: "vendor_id",
			Description:         "vendor_id",
			Computed:            true,
		},
		"sub_vendor_id": schema.StringAttribute{
			MarkdownDescription: "sub_vendor_id",
			Description:         "sub_vendor_id",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "created_date",
			Description:         "created_date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "created_by",
			Description:         "created_by",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "updated_date",
			Description:         "updated_date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "updated_by",
			Description:         "updated_by",
			Computed:            true,
		},
		"path": schema.StringAttribute{
			MarkdownDescription: "path",
			Description:         "path",
			Computed:            true,
		},
		"hash_md_5": schema.StringAttribute{
			MarkdownDescription: "hash_md_5",
			Description:         "hash_md_5",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"category": schema.StringAttribute{
			MarkdownDescription: "category",
			Description:         "category",
			Computed:            true,
		},
		"component_type": schema.StringAttribute{
			MarkdownDescription: "component_type",
			Description:         "component_type",
			Computed:            true,
		},
		"operating_system": schema.StringAttribute{
			MarkdownDescription: "operating_system",
			Description:         "operating_system",
			Computed:            true,
		},
		"system_i_ds": schema.ListAttribute{
			MarkdownDescription: "system_i_ds",
			Description:         "system_i_ds",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"custom": schema.BoolAttribute{
			MarkdownDescription: "custom",
			Description:         "custom",
			Computed:            true,
		},
		"needs_attention": schema.BoolAttribute{
			MarkdownDescription: "needs_attention",
			Description:         "needs_attention",
			Computed:            true,
		},
		"ignore": schema.BoolAttribute{
			MarkdownDescription: "ignore",
			Description:         "ignore",
			Computed:            true,
		},
		"original_version": schema.StringAttribute{
			MarkdownDescription: "original_version",
			Description:         "original_version",
			Computed:            true,
		},
		"original_component_id": schema.StringAttribute{
			MarkdownDescription: "original_component_id",
			Description:         "original_component_id",
			Computed:            true,
		},
		"firmware_repo_name": schema.StringAttribute{
			MarkdownDescription: "firmware_repo_name",
			Description:         "firmware_repo_name",
			Computed:            true,
		},
	}
}

// SoftwareBundlesSchema is a function that returns the schema for SoftwareBundles
func SoftwareBundlesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"version": schema.StringAttribute{
			MarkdownDescription: "version",
			Description:         "version",
			Computed:            true,
		},
		"bundle_date": schema.StringAttribute{
			MarkdownDescription: "bundle_date",
			Description:         "bundle_date",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "created_date",
			Description:         "created_date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "created_by",
			Description:         "created_by",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "updated_date",
			Description:         "updated_date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "updated_by",
			Description:         "updated_by",
			Computed:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "description",
			Description:         "description",
			Computed:            true,
		},
		"user_bundle": schema.BoolAttribute{
			MarkdownDescription: "user_bundle",
			Description:         "user_bundle",
			Computed:            true,
		},
		"user_bundle_path": schema.StringAttribute{
			MarkdownDescription: "user_bundle_path",
			Description:         "user_bundle_path",
			Computed:            true,
		},
		"user_bundle_hash_md_5": schema.StringAttribute{
			MarkdownDescription: "user_bundle_hash_md_5",
			Description:         "user_bundle_hash_md_5",
			Computed:            true,
		},
		"device_type": schema.StringAttribute{
			MarkdownDescription: "device_type",
			Description:         "device_type",
			Computed:            true,
		},
		"device_model": schema.StringAttribute{
			MarkdownDescription: "device_model",
			Description:         "device_model",
			Computed:            true,
		},
		"criticality": schema.StringAttribute{
			MarkdownDescription: "criticality",
			Description:         "criticality",
			Computed:            true,
		},
		"fw_repository_id": schema.StringAttribute{
			MarkdownDescription: "fw_repository_id",
			Description:         "fw_repository_id",
			Computed:            true,
		},
		"bundle_type": schema.StringAttribute{
			MarkdownDescription: "bundle_type",
			Description:         "bundle_type",
			Computed:            true,
		},
		"custom": schema.BoolAttribute{
			MarkdownDescription: "custom",
			Description:         "custom",
			Computed:            true,
		},
		"needs_attention": schema.BoolAttribute{
			MarkdownDescription: "needs_attention",
			Description:         "needs_attention",
			Computed:            true,
		},
		"software_components": schema.ListNestedAttribute{
			MarkdownDescription: "software_components",
			Description:         "software_components",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: SoftwareComponentsSchema()},
		},
	}
}

// DeploymentValidSchema is a function that returns the schema for DeploymentValid
func DeploymentValidSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"valid": schema.BoolAttribute{
			MarkdownDescription: "valid",
			Description:         "valid",
			Computed:            true,
		},
		"messages": schema.ListNestedAttribute{
			MarkdownDescription: "messages",
			Description:         "messages",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: MessagesSchema()},
		},
	}
}

// DeploymentDeviceSchema is a function that returns the schema for DeploymentDevice
func DeploymentDeviceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"ref_id": schema.StringAttribute{
			MarkdownDescription: "ref_id",
			Description:         "ref_id",
			Computed:            true,
		},
		"ref_type": schema.StringAttribute{
			MarkdownDescription: "ref_type",
			Description:         "ref_type",
			Computed:            true,
		},
		"log_dump": schema.StringAttribute{
			MarkdownDescription: "log_dump",
			Description:         "log_dump",
			Computed:            true,
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "status",
			Description:         "status",
			Computed:            true,
		},
		"status_end_time": schema.StringAttribute{
			MarkdownDescription: "status_end_time",
			Description:         "status_end_time",
			Computed:            true,
		},
		"status_start_time": schema.StringAttribute{
			MarkdownDescription: "status_start_time",
			Description:         "status_start_time",
			Computed:            true,
		},
		"device_health": schema.StringAttribute{
			MarkdownDescription: "device_health",
			Description:         "device_health",
			Computed:            true,
		},
		"health_message": schema.StringAttribute{
			MarkdownDescription: "health_message",
			Description:         "health_message",
			Computed:            true,
		},
		"compliant_state": schema.StringAttribute{
			MarkdownDescription: "compliant_state",
			Description:         "compliant_state",
			Computed:            true,
		},
		"brownfield_status": schema.StringAttribute{
			MarkdownDescription: "brownfield_status",
			Description:         "brownfield_status",
			Computed:            true,
		},
		"device_type": schema.StringAttribute{
			MarkdownDescription: "device_type",
			Description:         "device_type",
			Computed:            true,
		},
		"device_group_name": schema.StringAttribute{
			MarkdownDescription: "device_group_name",
			Description:         "device_group_name",
			Computed:            true,
		},
		"ip_address": schema.StringAttribute{
			MarkdownDescription: "ip_address",
			Description:         "ip_address",
			Computed:            true,
		},
		"current_ip_address": schema.StringAttribute{
			MarkdownDescription: "current_ip_address",
			Description:         "current_ip_address",
			Computed:            true,
		},
		"service_tag": schema.StringAttribute{
			MarkdownDescription: "service_tag",
			Description:         "service_tag",
			Computed:            true,
		},
		"component_id": schema.StringAttribute{
			MarkdownDescription: "component_id",
			Description:         "component_id",
			Computed:            true,
		},
		"status_message": schema.StringAttribute{
			MarkdownDescription: "status_message",
			Description:         "status_message",
			Computed:            true,
		},
		"model": schema.StringAttribute{
			MarkdownDescription: "model",
			Description:         "model",
			Computed:            true,
		},
		"cloud_link": schema.BoolAttribute{
			MarkdownDescription: "cloud_link",
			Description:         "cloud_link",
			Computed:            true,
		},
		"das_cache": schema.BoolAttribute{
			MarkdownDescription: "das_cache",
			Description:         "das_cache",
			Computed:            true,
		},
		"device_state": schema.StringAttribute{
			MarkdownDescription: "device_state",
			Description:         "device_state",
			Computed:            true,
		},
		"puppet_cert_name": schema.StringAttribute{
			MarkdownDescription: "puppet_cert_name",
			Description:         "puppet_cert_name",
			Computed:            true,
		},
		"brownfield": schema.BoolAttribute{
			MarkdownDescription: "brownfield",
			Description:         "brownfield",
			Computed:            true,
		},
	}
}

// VmsSchema is a function that returns the schema for Vms
func VmsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"certificate_name": schema.StringAttribute{
			MarkdownDescription: "certificate_name",
			Description:         "certificate_name",
			Computed:            true,
		},
		"vm_model": schema.StringAttribute{
			MarkdownDescription: "vm_model",
			Description:         "vm_model",
			Computed:            true,
		},
		"vm_ipaddress": schema.StringAttribute{
			MarkdownDescription: "vm_ipaddress",
			Description:         "vm_ipaddress",
			Computed:            true,
		},
		"vm_manufacturer": schema.StringAttribute{
			MarkdownDescription: "vm_manufacturer",
			Description:         "vm_manufacturer",
			Computed:            true,
		},
		"vm_service_tag": schema.StringAttribute{
			MarkdownDescription: "vm_service_tag",
			Description:         "vm_service_tag",
			Computed:            true,
		},
	}
}

// LicenseRepositorySchema is a function that returns the schema for LicenseRepository
func LicenseRepositorySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "type",
			Description:         "type",
			Computed:            true,
		},
		"disk_location": schema.StringAttribute{
			MarkdownDescription: "disk_location",
			Description:         "disk_location",
			Computed:            true,
		},
		"filename": schema.StringAttribute{
			MarkdownDescription: "filename",
			Description:         "filename",
			Computed:            true,
		},
		"state": schema.StringAttribute{
			MarkdownDescription: "state",
			Description:         "state",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "created_date",
			Description:         "created_date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "created_by",
			Description:         "created_by",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "updated_date",
			Description:         "updated_date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "updated_by",
			Description:         "updated_by",
			Computed:            true,
		},
		"license_data": schema.StringAttribute{
			MarkdownDescription: "license_data",
			Description:         "license_data",
			Computed:            true,
		},
	}
}

// AssignedUsersSchema is a function that returns the schema for AssignedUsers
func AssignedUsersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"user_seq_id": schema.Int64Attribute{
			MarkdownDescription: "user_seq_id",
			Description:         "user_seq_id",
			Computed:            true,
		},
		"user_name": schema.StringAttribute{
			MarkdownDescription: "user_name",
			Description:         "user_name",
			Computed:            true,
		},
		"password": schema.StringAttribute{
			MarkdownDescription: "password",
			Description:         "password",
			Computed:            true,
		},
		"update_password": schema.BoolAttribute{
			MarkdownDescription: "update_password",
			Description:         "update_password",
			Computed:            true,
		},
		"domain_name": schema.StringAttribute{
			MarkdownDescription: "domain_name",
			Description:         "domain_name",
			Computed:            true,
		},
		"group_dn": schema.StringAttribute{
			MarkdownDescription: "group_dn",
			Description:         "group_dn",
			Computed:            true,
		},
		"group_name": schema.StringAttribute{
			MarkdownDescription: "group_name",
			Description:         "group_name",
			Computed:            true,
		},
		"first_name": schema.StringAttribute{
			MarkdownDescription: "first_name",
			Description:         "first_name",
			Computed:            true,
		},
		"last_name": schema.StringAttribute{
			MarkdownDescription: "last_name",
			Description:         "last_name",
			Computed:            true,
		},
		"email": schema.StringAttribute{
			MarkdownDescription: "email",
			Description:         "email",
			Computed:            true,
		},
		"phone_number": schema.StringAttribute{
			MarkdownDescription: "phone_number",
			Description:         "phone_number",
			Computed:            true,
		},
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "enabled",
			Description:         "enabled",
			Computed:            true,
		},
		"system_user": schema.BoolAttribute{
			MarkdownDescription: "system_user",
			Description:         "system_user",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "created_date",
			Description:         "created_date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "created_by",
			Description:         "created_by",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "updated_date",
			Description:         "updated_date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "updated_by",
			Description:         "updated_by",
			Computed:            true,
		},
		"role": schema.StringAttribute{
			MarkdownDescription: "role",
			Description:         "role",
			Computed:            true,
		},
		"user_preference": schema.StringAttribute{
			MarkdownDescription: "user_preference",
			Description:         "user_preference",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"roles": schema.ListAttribute{
			MarkdownDescription: "roles",
			Description:         "roles",
			Computed:            true,
			ElementType:         types.StringType,
		},
	}
}

// JobDetailsSchema is a function that returns the schema for JobDetails
func JobDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"level": schema.StringAttribute{
			MarkdownDescription: "level",
			Description:         "level",
			Computed:            true,
		},
		"message": schema.StringAttribute{
			MarkdownDescription: "message",
			Description:         "message",
			Computed:            true,
		},
		"timestamp": schema.StringAttribute{
			MarkdownDescription: "timestamp",
			Description:         "timestamp",
			Computed:            true,
		},
		"execution_id": schema.StringAttribute{
			MarkdownDescription: "execution_id",
			Description:         "execution_id",
			Computed:            true,
		},
		"component_id": schema.StringAttribute{
			MarkdownDescription: "component_id",
			Description:         "component_id",
			Computed:            true,
		},
	}
}

// DeploymentValidationResponseSchema is a function that returns the schema for DeploymentValidationResponse
func DeploymentValidationResponseSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"nodes": schema.Int64Attribute{
			MarkdownDescription: "nodes",
			Description:         "nodes",
			Computed:            true,
		},
		"storage_pools": schema.Int64Attribute{
			MarkdownDescription: "storage_pools",
			Description:         "storage_pools",
			Computed:            true,
		},
		"drives_per_storage_pool": schema.Int64Attribute{
			MarkdownDescription: "drives_per_storage_pool",
			Description:         "drives_per_storage_pool",
			Computed:            true,
		},
		"max_scalability": schema.Int64Attribute{
			MarkdownDescription: "max_scalability",
			Description:         "max_scalability",
			Computed:            true,
		},
		"virtual_machines": schema.Int64Attribute{
			MarkdownDescription: "virtual_machines",
			Description:         "virtual_machines",
			Computed:            true,
		},
		"number_of_service_volumes": schema.Int64Attribute{
			MarkdownDescription: "number_of_service_volumes",
			Description:         "number_of_service_volumes",
			Computed:            true,
		},
		"can_deploy": schema.BoolAttribute{
			MarkdownDescription: "can_deploy",
			Description:         "can_deploy",
			Computed:            true,
		},
		"warning_messages": schema.ListAttribute{
			MarkdownDescription: "warning_messages",
			Description:         "warning_messages",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"storage_pool_disk_type": schema.ListAttribute{
			MarkdownDescription: "storage_pool_disk_type",
			Description:         "storage_pool_disk_type",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"hostnames": schema.ListAttribute{
			MarkdownDescription: "hostnames",
			Description:         "hostnames",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"new_node_disk_types": schema.ListAttribute{
			MarkdownDescription: "new_node_disk_types",
			Description:         "new_node_disk_types",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"no_of_fault_sets": schema.Int64Attribute{
			MarkdownDescription: "no_of_fault_sets",
			Description:         "no_of_fault_sets",
			Computed:            true,
		},
		"nodes_per_fault_set": schema.Int64Attribute{
			MarkdownDescription: "nodes_per_fault_set",
			Description:         "nodes_per_fault_set",
			Computed:            true,
		},
		"protection_domain": schema.StringAttribute{
			MarkdownDescription: "protection_domain",
			Description:         "protection_domain",
			Computed:            true,
		},
		"disk_type_mismatch": schema.BoolAttribute{
			MarkdownDescription: "disk_type_mismatch",
			Description:         "disk_type_mismatch",
			Computed:            true,
		},
	}
}

// DeploymentsSchema is a function that returns the schema for Deployments
func DeploymentsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"deployment_name": schema.StringAttribute{
			MarkdownDescription: "deployment_name",
			Description:         "deployment_name",
			Computed:            true,
		},
		"deployment_description": schema.StringAttribute{
			MarkdownDescription: "deployment_description",
			Description:         "deployment_description",
			Computed:            true,
		},
		"deployment_valid": schema.SingleNestedAttribute{
			MarkdownDescription: "deployment_valid",
			Description:         "deployment_valid",
			Computed:            true,
			Attributes:          DeploymentValidSchema(),
		},
		"retry": schema.BoolAttribute{
			MarkdownDescription: "retry",
			Description:         "retry",
			Computed:            true,
		},
		"teardown": schema.BoolAttribute{
			MarkdownDescription: "teardown",
			Description:         "teardown",
			Computed:            true,
		},
		"teardown_after_cancel": schema.BoolAttribute{
			MarkdownDescription: "teardown_after_cancel",
			Description:         "teardown_after_cancel",
			Computed:            true,
		},
		"remove_service": schema.BoolAttribute{
			MarkdownDescription: "remove_service",
			Description:         "remove_service",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "created_date",
			Description:         "created_date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "created_by",
			Description:         "created_by",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "updated_date",
			Description:         "updated_date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "updated_by",
			Description:         "updated_by",
			Computed:            true,
		},
		"deployment_scheduled_date": schema.StringAttribute{
			MarkdownDescription: "deployment_scheduled_date",
			Description:         "deployment_scheduled_date",
			Computed:            true,
		},
		"deployment_started_date": schema.StringAttribute{
			MarkdownDescription: "deployment_started_date",
			Description:         "deployment_started_date",
			Computed:            true,
		},
		"deployment_finished_date": schema.StringAttribute{
			MarkdownDescription: "deployment_finished_date",
			Description:         "deployment_finished_date",
			Computed:            true,
		},
		"schedule_date": schema.StringAttribute{
			MarkdownDescription: "schedule_date",
			Description:         "schedule_date",
			Computed:            true,
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "status",
			Description:         "status",
			Computed:            true,
		},
		"compliant": schema.BoolAttribute{
			MarkdownDescription: "compliant",
			Description:         "compliant",
			Computed:            true,
		},
		"deployment_device": schema.ListNestedAttribute{
			MarkdownDescription: "deployment_device",
			Description:         "deployment_device",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DeploymentDeviceSchema()},
		},
		"vms": schema.ListNestedAttribute{
			MarkdownDescription: "vms",
			Description:         "vms",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: VmsSchema()},
		},
		"update_server_firmware": schema.BoolAttribute{
			MarkdownDescription: "update_server_firmware",
			Description:         "update_server_firmware",
			Computed:            true,
		},
		"use_default_catalog": schema.BoolAttribute{
			MarkdownDescription: "use_default_catalog",
			Description:         "use_default_catalog",
			Computed:            true,
		},
		"firmware_repository_id": schema.StringAttribute{
			MarkdownDescription: "firmware_repository_id",
			Description:         "firmware_repository_id",
			Computed:            true,
		},
		"license_repository": schema.SingleNestedAttribute{
			MarkdownDescription: "license_repository",
			Description:         "license_repository",
			Computed:            true,
			Attributes:          LicenseRepositorySchema(),
		},
		"license_repository_id": schema.StringAttribute{
			MarkdownDescription: "license_repository_id",
			Description:         "license_repository_id",
			Computed:            true,
		},
		"individual_teardown": schema.BoolAttribute{
			MarkdownDescription: "individual_teardown",
			Description:         "individual_teardown",
			Computed:            true,
		},
		"deployment_health_status_type": schema.StringAttribute{
			MarkdownDescription: "deployment_health_status_type",
			Description:         "deployment_health_status_type",
			Computed:            true,
		},
		"assigned_users": schema.ListNestedAttribute{
			MarkdownDescription: "assigned_users",
			Description:         "assigned_users",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: AssignedUsersSchema()},
		},
		"all_users_allowed": schema.BoolAttribute{
			MarkdownDescription: "all_users_allowed",
			Description:         "all_users_allowed",
			Computed:            true,
		},
		"owner": schema.StringAttribute{
			MarkdownDescription: "owner",
			Description:         "owner",
			Computed:            true,
		},
		"no_op": schema.BoolAttribute{
			MarkdownDescription: "no_op",
			Description:         "no_op",
			Computed:            true,
		},
		"firmware_init": schema.BoolAttribute{
			MarkdownDescription: "firmware_init",
			Description:         "firmware_init",
			Computed:            true,
		},
		"disruptive_firmware": schema.BoolAttribute{
			MarkdownDescription: "disruptive_firmware",
			Description:         "disruptive_firmware",
			Computed:            true,
		},
		"preconfigure_svm": schema.BoolAttribute{
			MarkdownDescription: "preconfigure_svm",
			Description:         "preconfigure_svm",
			Computed:            true,
		},
		"preconfigure_svm_and_update": schema.BoolAttribute{
			MarkdownDescription: "preconfigure_svm_and_update",
			Description:         "preconfigure_svm_and_update",
			Computed:            true,
		},
		"services_deployed": schema.StringAttribute{
			MarkdownDescription: "services_deployed",
			Description:         "services_deployed",
			Computed:            true,
		},
		"precalculated_device_health": schema.StringAttribute{
			MarkdownDescription: "precalculated_device_health",
			Description:         "precalculated_device_health",
			Computed:            true,
		},
		"lifecycle_mode_reasons": schema.ListAttribute{
			MarkdownDescription: "lifecycle_mode_reasons",
			Description:         "lifecycle_mode_reasons",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"job_details": schema.ListNestedAttribute{
			MarkdownDescription: "job_details",
			Description:         "job_details",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: JobDetailsSchema()},
		},
		"number_of_deployments": schema.Int64Attribute{
			MarkdownDescription: "number_of_deployments",
			Description:         "number_of_deployments",
			Computed:            true,
		},
		"operation_type": schema.StringAttribute{
			MarkdownDescription: "operation_type",
			Description:         "operation_type",
			Computed:            true,
		},
		"operation_status": schema.StringAttribute{
			MarkdownDescription: "operation_status",
			Description:         "operation_status",
			Computed:            true,
		},
		"operation_data": schema.StringAttribute{
			MarkdownDescription: "operation_data",
			Description:         "operation_data",
			Computed:            true,
		},
		"deployment_validation_response": schema.SingleNestedAttribute{
			MarkdownDescription: "deployment_validation_response",
			Description:         "deployment_validation_response",
			Computed:            true,
			Attributes:          DeploymentValidationResponseSchema(),
		},
		"current_step_count": schema.StringAttribute{
			MarkdownDescription: "current_step_count",
			Description:         "current_step_count",
			Computed:            true,
		},
		"total_num_of_steps": schema.StringAttribute{
			MarkdownDescription: "total_num_of_steps",
			Description:         "total_num_of_steps",
			Computed:            true,
		},
		"current_step_message": schema.StringAttribute{
			MarkdownDescription: "current_step_message",
			Description:         "current_step_message",
			Computed:            true,
		},
		"custom_image": schema.StringAttribute{
			MarkdownDescription: "custom_image",
			Description:         "custom_image",
			Computed:            true,
		},
		"original_deployment_id": schema.StringAttribute{
			MarkdownDescription: "original_deployment_id",
			Description:         "original_deployment_id",
			Computed:            true,
		},
		"current_batch_count": schema.StringAttribute{
			MarkdownDescription: "current_batch_count",
			Description:         "current_batch_count",
			Computed:            true,
		},
		"total_batch_count": schema.StringAttribute{
			MarkdownDescription: "total_batch_count",
			Description:         "total_batch_count",
			Computed:            true,
		},
		"brownfield": schema.BoolAttribute{
			MarkdownDescription: "brownfield",
			Description:         "brownfield",
			Computed:            true,
		},
		"scale_up": schema.BoolAttribute{
			MarkdownDescription: "scale_up",
			Description:         "scale_up",
			Computed:            true,
		},
		"lifecycle_mode": schema.BoolAttribute{
			MarkdownDescription: "lifecycle_mode",
			Description:         "lifecycle_mode",
			Computed:            true,
		},
		"overall_device_health": schema.StringAttribute{
			MarkdownDescription: "overall_device_health",
			Description:         "overall_device_health",
			Computed:            true,
		},
		"vds": schema.BoolAttribute{
			MarkdownDescription: "vds",
			Description:         "vds",
			Computed:            true,
		},
		"template_valid": schema.BoolAttribute{
			MarkdownDescription: "template_valid",
			Description:         "template_valid",
			Computed:            true,
		},
		"configuration_change": schema.BoolAttribute{
			MarkdownDescription: "configuration_change",
			Description:         "configuration_change",
			Computed:            true,
		},
		"can_migratev_clsv_ms": schema.BoolAttribute{
			MarkdownDescription: "can_migratev_clsv_ms",
			Description:         "can_migratev_clsv_ms",
			Computed:            true,
		},
	}
}

// FirmwareRepositorySchema is a function that returns the schema for FirmwareRepository
func FirmwareRepositorySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"source_location": schema.StringAttribute{
			MarkdownDescription: "source_location",
			Description:         "source_location",
			Computed:            true,
		},
		"source_type": schema.StringAttribute{
			MarkdownDescription: "source_type",
			Description:         "source_type",
			Computed:            true,
		},
		"disk_location": schema.StringAttribute{
			MarkdownDescription: "disk_location",
			Description:         "disk_location",
			Computed:            true,
		},
		"filename": schema.StringAttribute{
			MarkdownDescription: "filename",
			Description:         "filename",
			Computed:            true,
		},
		"md_5_hash": schema.StringAttribute{
			MarkdownDescription: "md_5_hash",
			Description:         "md_5_hash",
			Computed:            true,
		},
		"username": schema.StringAttribute{
			MarkdownDescription: "username",
			Description:         "username",
			Computed:            true,
		},
		"password": schema.StringAttribute{
			MarkdownDescription: "password",
			Description:         "password",
			Computed:            true,
		},
		"download_status": schema.StringAttribute{
			MarkdownDescription: "download_status",
			Description:         "download_status",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "created_date",
			Description:         "created_date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "created_by",
			Description:         "created_by",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "updated_date",
			Description:         "updated_date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "updated_by",
			Description:         "updated_by",
			Computed:            true,
		},
		"default_catalog": schema.BoolAttribute{
			MarkdownDescription: "default_catalog",
			Description:         "default_catalog",
			Computed:            true,
		},
		"embedded": schema.BoolAttribute{
			MarkdownDescription: "embedded",
			Description:         "embedded",
			Computed:            true,
		},
		"state": schema.StringAttribute{
			MarkdownDescription: "state",
			Description:         "state",
			Computed:            true,
		},
		"software_components": schema.ListNestedAttribute{
			MarkdownDescription: "software_components",
			Description:         "software_components",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: SoftwareComponentsSchema()},
		},
		"software_bundles": schema.ListNestedAttribute{
			MarkdownDescription: "software_bundles",
			Description:         "software_bundles",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: SoftwareBundlesSchema()},
		},
		"deployments": schema.ListNestedAttribute{
			MarkdownDescription: "deployments",
			Description:         "deployments",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DeploymentsSchema()},
		},
		"bundle_count": schema.Int64Attribute{
			MarkdownDescription: "bundle_count",
			Description:         "bundle_count",
			Computed:            true,
		},
		"component_count": schema.Int64Attribute{
			MarkdownDescription: "component_count",
			Description:         "component_count",
			Computed:            true,
		},
		"user_bundle_count": schema.Int64Attribute{
			MarkdownDescription: "user_bundle_count",
			Description:         "user_bundle_count",
			Computed:            true,
		},
		"minimal": schema.BoolAttribute{
			MarkdownDescription: "minimal",
			Description:         "minimal",
			Computed:            true,
		},
		"download_progress": schema.Int64Attribute{
			MarkdownDescription: "download_progress",
			Description:         "download_progress",
			Computed:            true,
		},
		"extract_progress": schema.Int64Attribute{
			MarkdownDescription: "extract_progress",
			Description:         "extract_progress",
			Computed:            true,
		},
		"file_size_in_gigabytes": schema.Int64Attribute{
			MarkdownDescription: "file_size_in_gigabytes",
			Description:         "file_size_in_gigabytes",
			Computed:            true,
		},
		"signed_key_source_location": schema.StringAttribute{
			MarkdownDescription: "signed_key_source_location",
			Description:         "signed_key_source_location",
			Computed:            true,
		},
		"signature": schema.StringAttribute{
			MarkdownDescription: "signature",
			Description:         "signature",
			Computed:            true,
		},
		"custom": schema.BoolAttribute{
			MarkdownDescription: "custom",
			Description:         "custom",
			Computed:            true,
		},
		"needs_attention": schema.BoolAttribute{
			MarkdownDescription: "needs_attention",
			Description:         "needs_attention",
			Computed:            true,
		},
		"job_id": schema.StringAttribute{
			MarkdownDescription: "job_id",
			Description:         "job_id",
			Computed:            true,
		},
		"rcmapproved": schema.BoolAttribute{
			MarkdownDescription: "rcmapproved",
			Description:         "rcmapproved",
			Computed:            true,
		},
	}
}

// ComponentValidSchema is a function that returns the schema for ComponentValid
func ComponentValidSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"valid": schema.BoolAttribute{
			MarkdownDescription: "valid",
			Description:         "valid",
			Computed:            true,
		},
		"messages": schema.ListNestedAttribute{
			MarkdownDescription: "messages",
			Description:         "messages",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: MessagesSchema()},
		},
	}
}

// RelatedComponentsSchema is a function that returns the schema for RelatedComponents
func RelatedComponentsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"additional_prop_1": schema.StringAttribute{
			MarkdownDescription: "additional_prop_1",
			Description:         "additional_prop_1",
			Computed:            true,
		},
		"additional_prop_2": schema.StringAttribute{
			MarkdownDescription: "additional_prop_2",
			Description:         "additional_prop_2",
			Computed:            true,
		},
		"additional_prop_3": schema.StringAttribute{
			MarkdownDescription: "additional_prop_3",
			Description:         "additional_prop_3",
			Computed:            true,
		},
	}
}

// DependenciesDetailsSchema is a function that returns the schema for DependenciesDetails
func DependenciesDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"dependency_target": schema.StringAttribute{
			MarkdownDescription: "dependency_target",
			Description:         "dependency_target",
			Computed:            true,
		},
		"dependency_value": schema.StringAttribute{
			MarkdownDescription: "dependency_value",
			Description:         "dependency_value",
			Computed:            true,
		},
	}
}

// NetworkIPAddressListSchema is a function that returns the schema for NetworkIPAddressList
func NetworkIPAddressListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"ip_address": schema.StringAttribute{
			MarkdownDescription: "ip_address",
			Description:         "ip_address",
			Computed:            true,
		},
	}
}

// PartitionsSchema is a function that returns the schema for Partitions
func PartitionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"networks": schema.ListAttribute{
			MarkdownDescription: "networks",
			Description:         "networks",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"network_ip_address_list": schema.ListNestedAttribute{
			MarkdownDescription: "network_ip_address_list",
			Description:         "network_ip_address_list",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: NetworkIPAddressListSchema()},
		},
		"minimum": schema.Int64Attribute{
			MarkdownDescription: "minimum",
			Description:         "minimum",
			Computed:            true,
		},
		"maximum": schema.Int64Attribute{
			MarkdownDescription: "maximum",
			Description:         "maximum",
			Computed:            true,
		},
		"lan_mac_address": schema.StringAttribute{
			MarkdownDescription: "lan_mac_address",
			Description:         "lan_mac_address",
			Computed:            true,
		},
		"iscsi_mac_address": schema.StringAttribute{
			MarkdownDescription: "iscsi_mac_address",
			Description:         "iscsi_mac_address",
			Computed:            true,
		},
		"iscsi_iqn": schema.StringAttribute{
			MarkdownDescription: "iscsi_iqn",
			Description:         "iscsi_iqn",
			Computed:            true,
		},
		"wwnn": schema.StringAttribute{
			MarkdownDescription: "wwnn",
			Description:         "wwnn",
			Computed:            true,
		},
		"wwpn": schema.StringAttribute{
			MarkdownDescription: "wwpn",
			Description:         "wwpn",
			Computed:            true,
		},
		"fqdd": schema.StringAttribute{
			MarkdownDescription: "fqdd",
			Description:         "fqdd",
			Computed:            true,
		},
		"mirrored_port": schema.StringAttribute{
			MarkdownDescription: "mirrored_port",
			Description:         "mirrored_port",
			Computed:            true,
		},
		"mac_address": schema.StringAttribute{
			MarkdownDescription: "mac_address",
			Description:         "mac_address",
			Computed:            true,
		},
		"port_no": schema.Int64Attribute{
			MarkdownDescription: "port_no",
			Description:         "port_no",
			Computed:            true,
		},
		"partition_no": schema.Int64Attribute{
			MarkdownDescription: "partition_no",
			Description:         "partition_no",
			Computed:            true,
		},
		"partition_index": schema.Int64Attribute{
			MarkdownDescription: "partition_index",
			Description:         "partition_index",
			Computed:            true,
		},
	}
}

// InterfacesSchema is a function that returns the schema for Interfaces
func InterfacesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"partitioned": schema.BoolAttribute{
			MarkdownDescription: "partitioned",
			Description:         "partitioned",
			Computed:            true,
		},
		"partitions": schema.ListNestedAttribute{
			MarkdownDescription: "partitions",
			Description:         "partitions",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: PartitionsSchema()},
		},
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "enabled",
			Description:         "enabled",
			Computed:            true,
		},
		"redundancy": schema.BoolAttribute{
			MarkdownDescription: "redundancy",
			Description:         "redundancy",
			Computed:            true,
		},
		"nictype": schema.StringAttribute{
			MarkdownDescription: "nictype",
			Description:         "nictype",
			Computed:            true,
		},
		"fqdd": schema.StringAttribute{
			MarkdownDescription: "fqdd",
			Description:         "fqdd",
			Computed:            true,
		},
		"max_partitions": schema.Int64Attribute{
			MarkdownDescription: "max_partitions",
			Description:         "max_partitions",
			Computed:            true,
		},
		"all_networks": schema.ListAttribute{
			MarkdownDescription: "all_networks",
			Description:         "all_networks",
			Computed:            true,
			ElementType:         types.StringType,
		},
	}
}

// InterfacesDetailsSchema is a function that returns the schema for InterfacesDetails
func InterfacesDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"redundancy": schema.BoolAttribute{
			MarkdownDescription: "redundancy",
			Description:         "redundancy",
			Computed:            true,
		},
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "enabled",
			Description:         "enabled",
			Computed:            true,
		},
		"partitioned": schema.BoolAttribute{
			MarkdownDescription: "partitioned",
			Description:         "partitioned",
			Computed:            true,
		},
		"interfaces": schema.ListNestedAttribute{
			MarkdownDescription: "interfaces",
			Description:         "interfaces",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: InterfacesSchema()},
		},
		"nictype": schema.StringAttribute{
			MarkdownDescription: "nictype",
			Description:         "nictype",
			Computed:            true,
		},
		"fabrictype": schema.StringAttribute{
			MarkdownDescription: "fabrictype",
			Description:         "fabrictype",
			Computed:            true,
		},
		"max_partitions": schema.Int64Attribute{
			MarkdownDescription: "max_partitions",
			Description:         "max_partitions",
			Computed:            true,
		},
		"nports": schema.Int64Attribute{
			MarkdownDescription: "nports",
			Description:         "nports",
			Computed:            true,
		},
		"card_index": schema.Int64Attribute{
			MarkdownDescription: "card_index",
			Description:         "card_index",
			Computed:            true,
		},
		"nictype_source": schema.StringAttribute{
			MarkdownDescription: "nictype_source",
			Description:         "nictype_source",
			Computed:            true,
		},
	}
}

// NetworkConfigurationSchema is a function that returns the schema for NetworkConfiguration
func NetworkConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"interfaces": schema.ListNestedAttribute{
			MarkdownDescription: "interfaces",
			Description:         "interfaces",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: InterfacesDetailsSchema()},
		},
		"software_only": schema.BoolAttribute{
			MarkdownDescription: "software_only",
			Description:         "software_only",
			Computed:            true,
		},
	}
}

// ConfigurationDetailsSchema is a function that returns the schema for ConfigurationDetails
func ConfigurationDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"disktype": schema.StringAttribute{
			MarkdownDescription: "disktype",
			Description:         "disktype",
			Computed:            true,
		},
		"comparator": schema.StringAttribute{
			MarkdownDescription: "comparator",
			Description:         "comparator",
			Computed:            true,
		},
		"numberofdisks": schema.Int64Attribute{
			MarkdownDescription: "numberofdisks",
			Description:         "numberofdisks",
			Computed:            true,
		},
		"raidlevel": schema.StringAttribute{
			MarkdownDescription: "raidlevel",
			Description:         "raidlevel",
			Computed:            true,
		},
		"virtual_disk_fqdd": schema.StringAttribute{
			MarkdownDescription: "virtual_disk_fqdd",
			Description:         "virtual_disk_fqdd",
			Computed:            true,
		},
		"controller_fqdd": schema.StringAttribute{
			MarkdownDescription: "controller_fqdd",
			Description:         "controller_fqdd",
			Computed:            true,
		},
		"categories": schema.ListNestedAttribute{
			MarkdownDescription: "categories",
			Description:         "categories",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: CategoriesSchema()},
		},
	}
}

// VirtualDisksSchema is a function that returns the schema for VirtualDisks
func VirtualDisksSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"physical_disks": schema.ListAttribute{
			MarkdownDescription: "physical_disks",
			Description:         "physical_disks",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"virtual_disk_fqdd": schema.StringAttribute{
			MarkdownDescription: "virtual_disk_fqdd",
			Description:         "virtual_disk_fqdd",
			Computed:            true,
		},
		"raid_level": schema.StringAttribute{
			MarkdownDescription: "raid_level",
			Description:         "raid_level",
			Computed:            true,
		},
		"roll_up_status": schema.StringAttribute{
			MarkdownDescription: "roll_up_status",
			Description:         "roll_up_status",
			Computed:            true,
		},
		"controller": schema.StringAttribute{
			MarkdownDescription: "controller",
			Description:         "controller",
			Computed:            true,
		},
		"controller_product_name": schema.StringAttribute{
			MarkdownDescription: "controller_product_name",
			Description:         "controller_product_name",
			Computed:            true,
		},
		"configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "configuration",
			Description:         "configuration",
			Computed:            true,
			Attributes:          ConfigurationDetailsSchema(),
		},
		"media_type": schema.StringAttribute{
			MarkdownDescription: "media_type",
			Description:         "media_type",
			Computed:            true,
		},
		"encryption_type": schema.StringAttribute{
			MarkdownDescription: "encryption_type",
			Description:         "encryption_type",
			Computed:            true,
		},
	}
}

// ExternalVirtualDisksSchema is a function that returns the schema for ExternalVirtualDisks
func ExternalVirtualDisksSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"physical_disks": schema.ListAttribute{
			MarkdownDescription: "physical_disks",
			Description:         "physical_disks",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"virtual_disk_fqdd": schema.StringAttribute{
			MarkdownDescription: "virtual_disk_fqdd",
			Description:         "virtual_disk_fqdd",
			Computed:            true,
		},
		"raid_level": schema.StringAttribute{
			MarkdownDescription: "raid_level",
			Description:         "raid_level",
			Computed:            true,
		},
		"roll_up_status": schema.StringAttribute{
			MarkdownDescription: "roll_up_status",
			Description:         "roll_up_status",
			Computed:            true,
		},
		"controller": schema.StringAttribute{
			MarkdownDescription: "controller",
			Description:         "controller",
			Computed:            true,
		},
		"controller_product_name": schema.StringAttribute{
			MarkdownDescription: "controller_product_name",
			Description:         "controller_product_name",
			Computed:            true,
		},
		"configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "configuration",
			Description:         "configuration",
			Computed:            true,
			Attributes:          ConfigurationDetailsSchema(),
		},
		"media_type": schema.StringAttribute{
			MarkdownDescription: "media_type",
			Description:         "media_type",
			Computed:            true,
		},
		"encryption_type": schema.StringAttribute{
			MarkdownDescription: "encryption_type",
			Description:         "encryption_type",
			Computed:            true,
		},
	}
}

// SizeToDiskMapSchema is a function that returns the schema for SizeToDiskMap
func SizeToDiskMapSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"additional_prop_1": schema.Int64Attribute{
			MarkdownDescription: "additional_prop_1",
			Description:         "additional_prop_1",
			Computed:            true,
		},
		"additional_prop_2": schema.Int64Attribute{
			MarkdownDescription: "additional_prop_2",
			Description:         "additional_prop_2",
			Computed:            true,
		},
		"additional_prop_3": schema.Int64Attribute{
			MarkdownDescription: "additional_prop_3",
			Description:         "additional_prop_3",
			Computed:            true,
		},
	}
}

// RaidConfigurationSchema is a function that returns the schema for RaidConfiguration
func RaidConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"virtual_disks": schema.ListNestedAttribute{
			MarkdownDescription: "virtual_disks",
			Description:         "virtual_disks",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: VirtualDisksSchema()},
		},
		"external_virtual_disks": schema.ListNestedAttribute{
			MarkdownDescription: "external_virtual_disks",
			Description:         "external_virtual_disks",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ExternalVirtualDisksSchema()},
		},
		"hdd_hot_spares": schema.ListAttribute{
			MarkdownDescription: "hdd_hot_spares",
			Description:         "hdd_hot_spares",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"ssd_hot_spares": schema.ListAttribute{
			MarkdownDescription: "ssd_hot_spares",
			Description:         "ssd_hot_spares",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"external_hdd_hot_spares": schema.ListAttribute{
			MarkdownDescription: "external_hdd_hot_spares",
			Description:         "external_hdd_hot_spares",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"external_ssd_hot_spares": schema.ListAttribute{
			MarkdownDescription: "external_ssd_hot_spares",
			Description:         "external_ssd_hot_spares",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"size_to_disk_map": schema.SingleNestedAttribute{
			MarkdownDescription: "size_to_disk_map",
			Description:         "size_to_disk_map",
			Computed:            true,
			Attributes:          SizeToDiskMapSchema(),
		},
	}
}

// AttributesSchema is a function that returns the schema for Attributes
func AttributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"additional_prop_1": schema.StringAttribute{
			MarkdownDescription: "additional_prop_1",
			Description:         "additional_prop_1",
			Computed:            true,
		},
		"additional_prop_2": schema.StringAttribute{
			MarkdownDescription: "additional_prop_2",
			Description:         "additional_prop_2",
			Computed:            true,
		},
		"additional_prop_3": schema.StringAttribute{
			MarkdownDescription: "additional_prop_3",
			Description:         "additional_prop_3",
			Computed:            true,
		},
	}
}

// OptionsDetailsSchema is a function that returns the schema for OptionsDetails
func OptionsDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "value",
			Description:         "value",
			Computed:            true,
		},
		"dependencies": schema.ListNestedAttribute{
			MarkdownDescription: "dependencies",
			Description:         "dependencies",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DependenciesDetailsSchema()},
		},
		"attributes": schema.SingleNestedAttribute{
			MarkdownDescription: "attributes",
			Description:         "attributes",
			Computed:            true,
			Attributes:          AttributesSchema(),
		},
	}
}

// ScaleIOStoragePoolDisksSchema is a function that returns the schema for ScaleIOStoragePoolDisks
func ScaleIOStoragePoolDisksSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "protection_domain_id",
			Description:         "protection_domain_id",
			Computed:            true,
		},
		"protection_domain_name": schema.StringAttribute{
			MarkdownDescription: "protection_domain_name",
			Description:         "protection_domain_name",
			Computed:            true,
		},
		"storage_pool_id": schema.StringAttribute{
			MarkdownDescription: "storage_pool_id",
			Description:         "storage_pool_id",
			Computed:            true,
		},
		"storage_pool_name": schema.StringAttribute{
			MarkdownDescription: "storage_pool_name",
			Description:         "storage_pool_name",
			Computed:            true,
		},
		"disk_type": schema.StringAttribute{
			MarkdownDescription: "disk_type",
			Description:         "disk_type",
			Computed:            true,
		},
		"physical_disk_fqdds": schema.ListAttribute{
			MarkdownDescription: "physical_disk_fqdds",
			Description:         "physical_disk_fqdds",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"virtual_disk_fqdds": schema.ListAttribute{
			MarkdownDescription: "virtual_disk_fqdds",
			Description:         "virtual_disk_fqdds",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"software_only_disks": schema.ListAttribute{
			MarkdownDescription: "software_only_disks",
			Description:         "software_only_disks",
			Computed:            true,
			ElementType:         types.StringType,
		},
	}
}

// ScaleIODiskConfigurationSchema is a function that returns the schema for ScaleIODiskConfiguration
func ScaleIODiskConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"scale_io_storage_pool_disks": schema.ListNestedAttribute{
			MarkdownDescription: "scale_io_storage_pool_disks",
			Description:         "scale_io_storage_pool_disks",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ScaleIOStoragePoolDisksSchema()},
		},
	}
}

// ShortWindowSchema is a function that returns the schema for ShortWindow
func ShortWindowSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"threshold": schema.Int64Attribute{
			MarkdownDescription: "threshold",
			Description:         "threshold",
			Computed:            true,
		},
		"window_size_in_sec": schema.Int64Attribute{
			MarkdownDescription: "window_size_in_sec",
			Description:         "window_size_in_sec",
			Computed:            true,
		},
	}
}

// MediumWindowSchema is a function that returns the schema for MediumWindow
func MediumWindowSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"threshold": schema.Int64Attribute{
			MarkdownDescription: "threshold",
			Description:         "threshold",
			Computed:            true,
		},
		"window_size_in_sec": schema.Int64Attribute{
			MarkdownDescription: "window_size_in_sec",
			Description:         "window_size_in_sec",
			Computed:            true,
		},
	}
}

// LongWindowSchema is a function that returns the schema for LongWindow
func LongWindowSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"threshold": schema.Int64Attribute{
			MarkdownDescription: "threshold",
			Description:         "threshold",
			Computed:            true,
		},
		"window_size_in_sec": schema.Int64Attribute{
			MarkdownDescription: "window_size_in_sec",
			Description:         "window_size_in_sec",
			Computed:            true,
		},
	}
}

// SdsDecoupledCounterParametersSchema is a function that returns the schema for SdsDecoupledCounterParameters
func SdsDecoupledCounterParametersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			MarkdownDescription: "short_window",
			Description:         "short_window",
			Computed:            true,
			Attributes:          ShortWindowSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			MarkdownDescription: "medium_window",
			Description:         "medium_window",
			Computed:            true,
			Attributes:          MediumWindowSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			MarkdownDescription: "long_window",
			Description:         "long_window",
			Computed:            true,
			Attributes:          LongWindowSchema(),
		},
	}
}

// SdsConfigurationFailureCounterParametersSchema is a function that returns the schema for SdsConfigurationFailureCounterParameters
func SdsConfigurationFailureCounterParametersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			MarkdownDescription: "short_window",
			Description:         "short_window",
			Computed:            true,
			Attributes:          ShortWindowSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			MarkdownDescription: "medium_window",
			Description:         "medium_window",
			Computed:            true,
			Attributes:          MediumWindowSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			MarkdownDescription: "long_window",
			Description:         "long_window",
			Computed:            true,
			Attributes:          LongWindowSchema(),
		},
	}
}

// MdmSdsCounterParametersSchema is a function that returns the schema for MdmSdsCounterParameters
func MdmSdsCounterParametersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			MarkdownDescription: "short_window",
			Description:         "short_window",
			Computed:            true,
			Attributes:          ShortWindowSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			MarkdownDescription: "medium_window",
			Description:         "medium_window",
			Computed:            true,
			Attributes:          MediumWindowSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			MarkdownDescription: "long_window",
			Description:         "long_window",
			Computed:            true,
			Attributes:          LongWindowSchema(),
		},
	}
}

// SdsSdsCounterParametersSchema is a function that returns the schema for SdsSdsCounterParameters
func SdsSdsCounterParametersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			MarkdownDescription: "short_window",
			Description:         "short_window",
			Computed:            true,
			Attributes:          ShortWindowSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			MarkdownDescription: "medium_window",
			Description:         "medium_window",
			Computed:            true,
			Attributes:          MediumWindowSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			MarkdownDescription: "long_window",
			Description:         "long_window",
			Computed:            true,
			Attributes:          LongWindowSchema(),
		},
	}
}

// SdsReceiveBufferAllocationFailuresCounterParametersSchema is a function that returns the schema for SdsReceiveBufferAllocationFailuresCounterParameters
func SdsReceiveBufferAllocationFailuresCounterParametersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			MarkdownDescription: "short_window",
			Description:         "short_window",
			Computed:            true,
			Attributes:          ShortWindowSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			MarkdownDescription: "medium_window",
			Description:         "medium_window",
			Computed:            true,
			Attributes:          MediumWindowSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			MarkdownDescription: "long_window",
			Description:         "long_window",
			Computed:            true,
			Attributes:          LongWindowSchema(),
		},
	}
}

// GeneralSchema is a function that returns the schema for General
func GeneralSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"system_id": schema.StringAttribute{
			MarkdownDescription: "system_id",
			Description:         "system_id",
			Computed:            true,
		},
		"protection_domain_state": schema.StringAttribute{
			MarkdownDescription: "protection_domain_state",
			Description:         "protection_domain_state",
			Computed:            true,
		},
		"rebuild_network_throttling_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "rebuild_network_throttling_in_kbps",
			Description:         "rebuild_network_throttling_in_kbps",
			Computed:            true,
		},
		"rebalance_network_throttling_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "rebalance_network_throttling_in_kbps",
			Description:         "rebalance_network_throttling_in_kbps",
			Computed:            true,
		},
		"overall_io_network_throttling_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "overall_io_network_throttling_in_kbps",
			Description:         "overall_io_network_throttling_in_kbps",
			Computed:            true,
		},
		"sds_decoupled_counter_parameters": schema.SingleNestedAttribute{
			MarkdownDescription: "sds_decoupled_counter_parameters",
			Description:         "sds_decoupled_counter_parameters",
			Computed:            true,
			Attributes:          SdsDecoupledCounterParametersSchema(),
		},
		"sds_configuration_failure_counter_parameters": schema.SingleNestedAttribute{
			MarkdownDescription: "sds_configuration_failure_counter_parameters",
			Description:         "sds_configuration_failure_counter_parameters",
			Computed:            true,
			Attributes:          SdsConfigurationFailureCounterParametersSchema(),
		},
		"mdm_sds_counter_parameters": schema.SingleNestedAttribute{
			MarkdownDescription: "mdm_sds_counter_parameters",
			Description:         "mdm_sds_counter_parameters",
			Computed:            true,
			Attributes:          MdmSdsCounterParametersSchema(),
		},
		"sds_sds_counter_parameters": schema.SingleNestedAttribute{
			MarkdownDescription: "sds_sds_counter_parameters",
			Description:         "sds_sds_counter_parameters",
			Computed:            true,
			Attributes:          SdsSdsCounterParametersSchema(),
		},
		"rfcache_opertional_mode": schema.StringAttribute{
			MarkdownDescription: "rfcache_opertional_mode",
			Description:         "rfcache_opertional_mode",
			Computed:            true,
		},
		"rfcache_page_size_kb": schema.Int64Attribute{
			MarkdownDescription: "rfcache_page_size_kb",
			Description:         "rfcache_page_size_kb",
			Computed:            true,
		},
		"rfcache_max_io_size_kb": schema.Int64Attribute{
			MarkdownDescription: "rfcache_max_io_size_kb",
			Description:         "rfcache_max_io_size_kb",
			Computed:            true,
		},
		"sds_receive_buffer_allocation_failures_counter_parameters": schema.SingleNestedAttribute{
			MarkdownDescription: "sds_receive_buffer_allocation_failures_counter_parameters",
			Description:         "sds_receive_buffer_allocation_failures_counter_parameters",
			Computed:            true,
			Attributes:          SdsReceiveBufferAllocationFailuresCounterParametersSchema(),
		},
		"rebuild_network_throttling_enabled": schema.BoolAttribute{
			MarkdownDescription: "rebuild_network_throttling_enabled",
			Description:         "rebuild_network_throttling_enabled",
			Computed:            true,
		},
		"rebalance_network_throttling_enabled": schema.BoolAttribute{
			MarkdownDescription: "rebalance_network_throttling_enabled",
			Description:         "rebalance_network_throttling_enabled",
			Computed:            true,
		},
		"overall_io_network_throttling_enabled": schema.BoolAttribute{
			MarkdownDescription: "overall_io_network_throttling_enabled",
			Description:         "overall_io_network_throttling_enabled",
			Computed:            true,
		},
		"rfcache_enabled": schema.BoolAttribute{
			MarkdownDescription: "rfcache_enabled",
			Description:         "rfcache_enabled",
			Computed:            true,
		},
	}
}

// StatisticsDetailsSchema is a function that returns the schema for StatisticsDetails
func StatisticsDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"num_of_devices": schema.Int64Attribute{
			MarkdownDescription: "num_of_devices",
			Description:         "num_of_devices",
			Computed:            true,
		},
		"unused_capacity_in_kb": schema.Int64Attribute{
			MarkdownDescription: "unused_capacity_in_kb",
			Description:         "unused_capacity_in_kb",
			Computed:            true,
		},
		"num_of_volumes": schema.Int64Attribute{
			MarkdownDescription: "num_of_volumes",
			Description:         "num_of_volumes",
			Computed:            true,
		},
		"num_of_mapped_to_all_volumes": schema.Int64Attribute{
			MarkdownDescription: "num_of_mapped_to_all_volumes",
			Description:         "num_of_mapped_to_all_volumes",
			Computed:            true,
		},
		"capacity_available_for_volume_allocation_in_kb": schema.Int64Attribute{
			MarkdownDescription: "capacity_available_for_volume_allocation_in_kb",
			Description:         "capacity_available_for_volume_allocation_in_kb",
			Computed:            true,
		},
		"volume_allocation_limit_in_kb": schema.Int64Attribute{
			MarkdownDescription: "volume_allocation_limit_in_kb",
			Description:         "volume_allocation_limit_in_kb",
			Computed:            true,
		},
		"capacity_limit_in_kb": schema.Int64Attribute{
			MarkdownDescription: "capacity_limit_in_kb",
			Description:         "capacity_limit_in_kb",
			Computed:            true,
		},
		"num_of_unmapped_volumes": schema.Int64Attribute{
			MarkdownDescription: "num_of_unmapped_volumes",
			Description:         "num_of_unmapped_volumes",
			Computed:            true,
		},
		"spare_capacity_in_kb": schema.Int64Attribute{
			MarkdownDescription: "spare_capacity_in_kb",
			Description:         "spare_capacity_in_kb",
			Computed:            true,
		},
		"capacity_in_use_in_kb": schema.Int64Attribute{
			MarkdownDescription: "capacity_in_use_in_kb",
			Description:         "capacity_in_use_in_kb",
			Computed:            true,
		},
		"max_capacity_in_kb": schema.Int64Attribute{
			MarkdownDescription: "max_capacity_in_kb",
			Description:         "max_capacity_in_kb",
			Computed:            true,
		},
		"num_of_sds": schema.Int64Attribute{
			MarkdownDescription: "num_of_sds",
			Description:         "num_of_sds",
			Computed:            true,
		},
		"num_of_storage_pools": schema.Int64Attribute{
			MarkdownDescription: "num_of_storage_pools",
			Description:         "num_of_storage_pools",
			Computed:            true,
		},
		"num_of_fault_sets": schema.Int64Attribute{
			MarkdownDescription: "num_of_fault_sets",
			Description:         "num_of_fault_sets",
			Computed:            true,
		},
		"thin_capacity_in_use_in_kb": schema.Int64Attribute{
			MarkdownDescription: "thin_capacity_in_use_in_kb",
			Description:         "thin_capacity_in_use_in_kb",
			Computed:            true,
		},
		"thick_capacity_in_use_in_kb": schema.Int64Attribute{
			MarkdownDescription: "thick_capacity_in_use_in_kb",
			Description:         "thick_capacity_in_use_in_kb",
			Computed:            true,
		},
	}
}

// DiskListSchema is a function that returns the schema for DiskList
func DiskListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"error_state": schema.StringAttribute{
			MarkdownDescription: "error_state",
			Description:         "error_state",
			Computed:            true,
		},
		"sds_id": schema.StringAttribute{
			MarkdownDescription: "sds_id",
			Description:         "sds_id",
			Computed:            true,
		},
		"device_state": schema.StringAttribute{
			MarkdownDescription: "device_state",
			Description:         "device_state",
			Computed:            true,
		},
		"capacity_limit_in_kb": schema.Int64Attribute{
			MarkdownDescription: "capacity_limit_in_kb",
			Description:         "capacity_limit_in_kb",
			Computed:            true,
		},
		"max_capacity_in_kb": schema.Int64Attribute{
			MarkdownDescription: "max_capacity_in_kb",
			Description:         "max_capacity_in_kb",
			Computed:            true,
		},
		"storage_pool_id": schema.StringAttribute{
			MarkdownDescription: "storage_pool_id",
			Description:         "storage_pool_id",
			Computed:            true,
		},
		"device_current_path_name": schema.StringAttribute{
			MarkdownDescription: "device_current_path_name",
			Description:         "device_current_path_name",
			Computed:            true,
		},
		"device_original_path_name": schema.StringAttribute{
			MarkdownDescription: "device_original_path_name",
			Description:         "device_original_path_name",
			Computed:            true,
		},
		"serial_number": schema.StringAttribute{
			MarkdownDescription: "serial_number",
			Description:         "serial_number",
			Computed:            true,
		},
		"vendor_name": schema.StringAttribute{
			MarkdownDescription: "vendor_name",
			Description:         "vendor_name",
			Computed:            true,
		},
		"model_name": schema.StringAttribute{
			MarkdownDescription: "model_name",
			Description:         "model_name",
			Computed:            true,
		},
	}
}

// MappedSdcInfoDetailsSchema is a function that returns the schema for MappedSdcInfoDetails
func MappedSdcInfoDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"sdc_ip": schema.StringAttribute{
			MarkdownDescription: "sdc_ip",
			Description:         "sdc_ip",
			Computed:            true,
		},
		"sdc_id": schema.StringAttribute{
			MarkdownDescription: "sdc_id",
			Description:         "sdc_id",
			Computed:            true,
		},
		"limit_bw_in_mbps": schema.Int64Attribute{
			MarkdownDescription: "limit_bw_in_mbps",
			Description:         "limit_bw_in_mbps",
			Computed:            true,
		},
		"limit_iops": schema.Int64Attribute{
			MarkdownDescription: "limit_iops",
			Description:         "limit_iops",
			Computed:            true,
		},
	}
}

// VolumeListSchema is a function that returns the schema for VolumeList
func VolumeListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"volume_type": schema.StringAttribute{
			MarkdownDescription: "volume_type",
			Description:         "volume_type",
			Computed:            true,
		},
		"storage_pool_id": schema.StringAttribute{
			MarkdownDescription: "storage_pool_id",
			Description:         "storage_pool_id",
			Computed:            true,
		},
		"data_layout": schema.StringAttribute{
			MarkdownDescription: "data_layout",
			Description:         "data_layout",
			Computed:            true,
		},
		"compression_method": schema.StringAttribute{
			MarkdownDescription: "compression_method",
			Description:         "compression_method",
			Computed:            true,
		},
		"size_in_kb": schema.Int64Attribute{
			MarkdownDescription: "size_in_kb",
			Description:         "size_in_kb",
			Computed:            true,
		},
		"mapped_sdc_info": schema.ListNestedAttribute{
			MarkdownDescription: "mapped_sdc_info",
			Description:         "mapped_sdc_info",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: MappedSdcInfoDetailsSchema()},
		},
		"volume_class": schema.StringAttribute{
			MarkdownDescription: "volume_class",
			Description:         "volume_class",
			Computed:            true,
		},
	}
}

// StoragePoolListSchema is a function that returns the schema for StoragePoolList
func StoragePoolListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"rebuild_io_priority_policy": schema.StringAttribute{
			MarkdownDescription: "rebuild_io_priority_policy",
			Description:         "rebuild_io_priority_policy",
			Computed:            true,
		},
		"rebalance_io_priority_policy": schema.StringAttribute{
			MarkdownDescription: "rebalance_io_priority_policy",
			Description:         "rebalance_io_priority_policy",
			Computed:            true,
		},
		"rebuild_io_priority_num_of_concurrent_ios_per_device": schema.Int64Attribute{
			MarkdownDescription: "rebuild_io_priority_num_of_concurrent_ios_per_device",
			Description:         "rebuild_io_priority_num_of_concurrent_ios_per_device",
			Computed:            true,
		},
		"rebalance_io_priority_num_of_concurrent_ios_per_device": schema.Int64Attribute{
			MarkdownDescription: "rebalance_io_priority_num_of_concurrent_ios_per_device",
			Description:         "rebalance_io_priority_num_of_concurrent_ios_per_device",
			Computed:            true,
		},
		"rebuild_io_priority_bw_limit_per_device_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "rebuild_io_priority_bw_limit_per_device_in_kbps",
			Description:         "rebuild_io_priority_bw_limit_per_device_in_kbps",
			Computed:            true,
		},
		"rebalance_io_priority_bw_limit_per_device_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "rebalance_io_priority_bw_limit_per_device_in_kbps",
			Description:         "rebalance_io_priority_bw_limit_per_device_in_kbps",
			Computed:            true,
		},
		"rebuild_io_priority_app_iops_per_device_threshold": schema.StringAttribute{
			MarkdownDescription: "rebuild_io_priority_app_iops_per_device_threshold",
			Description:         "rebuild_io_priority_app_iops_per_device_threshold",
			Computed:            true,
		},
		"rebalance_io_priority_app_iops_per_device_threshold": schema.StringAttribute{
			MarkdownDescription: "rebalance_io_priority_app_iops_per_device_threshold",
			Description:         "rebalance_io_priority_app_iops_per_device_threshold",
			Computed:            true,
		},
		"rebuild_io_priority_app_bw_per_device_threshold_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "rebuild_io_priority_app_bw_per_device_threshold_in_kbps",
			Description:         "rebuild_io_priority_app_bw_per_device_threshold_in_kbps",
			Computed:            true,
		},
		"rebalance_io_priority_app_bw_per_device_threshold_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "rebalance_io_priority_app_bw_per_device_threshold_in_kbps",
			Description:         "rebalance_io_priority_app_bw_per_device_threshold_in_kbps",
			Computed:            true,
		},
		"rebuild_io_priority_quiet_period_in_msec": schema.Int64Attribute{
			MarkdownDescription: "rebuild_io_priority_quiet_period_in_msec",
			Description:         "rebuild_io_priority_quiet_period_in_msec",
			Computed:            true,
		},
		"rebalance_io_priority_quiet_period_in_msec": schema.Int64Attribute{
			MarkdownDescription: "rebalance_io_priority_quiet_period_in_msec",
			Description:         "rebalance_io_priority_quiet_period_in_msec",
			Computed:            true,
		},
		"zero_padding_enabled": schema.BoolAttribute{
			MarkdownDescription: "zero_padding_enabled",
			Description:         "zero_padding_enabled",
			Computed:            true,
		},
		"background_scanner_mode": schema.StringAttribute{
			MarkdownDescription: "background_scanner_mode",
			Description:         "background_scanner_mode",
			Computed:            true,
		},
		"background_scanner_bw_limit_k_bps": schema.Int64Attribute{
			MarkdownDescription: "background_scanner_bw_limit_k_bps",
			Description:         "background_scanner_bw_limit_k_bps",
			Computed:            true,
		},
		"use_rmcache": schema.BoolAttribute{
			MarkdownDescription: "use_rmcache",
			Description:         "use_rmcache",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "protection_domain_id",
			Description:         "protection_domain_id",
			Computed:            true,
		},
		"sp_class": schema.StringAttribute{
			MarkdownDescription: "sp_class",
			Description:         "sp_class",
			Computed:            true,
		},
		"use_rfcache": schema.BoolAttribute{
			MarkdownDescription: "use_rfcache",
			Description:         "use_rfcache",
			Computed:            true,
		},
		"spare_percentage": schema.Int64Attribute{
			MarkdownDescription: "spare_percentage",
			Description:         "spare_percentage",
			Computed:            true,
		},
		"rmcache_write_handling_mode": schema.StringAttribute{
			MarkdownDescription: "rmcache_write_handling_mode",
			Description:         "rmcache_write_handling_mode",
			Computed:            true,
		},
		"checksum_enabled": schema.BoolAttribute{
			MarkdownDescription: "checksum_enabled",
			Description:         "checksum_enabled",
			Computed:            true,
		},
		"rebuild_enabled": schema.BoolAttribute{
			MarkdownDescription: "rebuild_enabled",
			Description:         "rebuild_enabled",
			Computed:            true,
		},
		"rebalance_enabled": schema.BoolAttribute{
			MarkdownDescription: "rebalance_enabled",
			Description:         "rebalance_enabled",
			Computed:            true,
		},
		"num_of_parallel_rebuild_rebalance_jobs_per_device": schema.Int64Attribute{
			MarkdownDescription: "num_of_parallel_rebuild_rebalance_jobs_per_device",
			Description:         "num_of_parallel_rebuild_rebalance_jobs_per_device",
			Computed:            true,
		},
		"capacity_alert_high_threshold": schema.Int64Attribute{
			MarkdownDescription: "capacity_alert_high_threshold",
			Description:         "capacity_alert_high_threshold",
			Computed:            true,
		},
		"capacity_alert_critical_threshold": schema.Int64Attribute{
			MarkdownDescription: "capacity_alert_critical_threshold",
			Description:         "capacity_alert_critical_threshold",
			Computed:            true,
		},
		"statistics": schema.SingleNestedAttribute{
			MarkdownDescription: "statistics",
			Description:         "statistics",
			Computed:            true,
			Attributes:          StatisticsDetailsSchema(),
		},
		"data_layout": schema.StringAttribute{
			MarkdownDescription: "data_layout",
			Description:         "data_layout",
			Computed:            true,
		},
		"replication_capacity_max_ratio": schema.StringAttribute{
			MarkdownDescription: "replication_capacity_max_ratio",
			Description:         "replication_capacity_max_ratio",
			Computed:            true,
		},
		"media_type": schema.StringAttribute{
			MarkdownDescription: "media_type",
			Description:         "media_type",
			Computed:            true,
		},
		"disk_list": schema.ListNestedAttribute{
			MarkdownDescription: "disk_list",
			Description:         "disk_list",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DiskListSchema()},
		},
		"volume_list": schema.ListNestedAttribute{
			MarkdownDescription: "volume_list",
			Description:         "volume_list",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: VolumeListSchema()},
		},
		"fgl_accp_id": schema.StringAttribute{
			MarkdownDescription: "fgl_accp_id",
			Description:         "fgl_accp_id",
			Computed:            true,
		},
	}
}

// IPListSchema is a function that returns the schema for IPList
func IPListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"ip": schema.StringAttribute{
			MarkdownDescription: "ip",
			Description:         "ip",
			Computed:            true,
		},
		"role": schema.StringAttribute{
			MarkdownDescription: "role",
			Description:         "role",
			Computed:            true,
		},
	}
}

// SdsListDetailsSchema is a function that returns the schema for SdsListDetails
func SdsListDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"port": schema.Int64Attribute{
			MarkdownDescription: "port",
			Description:         "port",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "protection_domain_id",
			Description:         "protection_domain_id",
			Computed:            true,
		},
		"fault_set_id": schema.StringAttribute{
			MarkdownDescription: "fault_set_id",
			Description:         "fault_set_id",
			Computed:            true,
		},
		"software_version_info": schema.StringAttribute{
			MarkdownDescription: "software_version_info",
			Description:         "software_version_info",
			Computed:            true,
		},
		"sds_state": schema.StringAttribute{
			MarkdownDescription: "sds_state",
			Description:         "sds_state",
			Computed:            true,
		},
		"membership_state": schema.StringAttribute{
			MarkdownDescription: "membership_state",
			Description:         "membership_state",
			Computed:            true,
		},
		"mdm_connection_state": schema.StringAttribute{
			MarkdownDescription: "mdm_connection_state",
			Description:         "mdm_connection_state",
			Computed:            true,
		},
		"drl_mode": schema.StringAttribute{
			MarkdownDescription: "drl_mode",
			Description:         "drl_mode",
			Computed:            true,
		},
		"maintenance_state": schema.StringAttribute{
			MarkdownDescription: "maintenance_state",
			Description:         "maintenance_state",
			Computed:            true,
		},
		"perf_profile": schema.StringAttribute{
			MarkdownDescription: "perf_profile",
			Description:         "perf_profile",
			Computed:            true,
		},
		"on_vm_ware": schema.BoolAttribute{
			MarkdownDescription: "on_vm_ware",
			Description:         "on_vm_ware",
			Computed:            true,
		},
		"ip_list": schema.ListNestedAttribute{
			MarkdownDescription: "ip_list",
			Description:         "ip_list",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: IPListSchema()},
		},
	}
}

// SdrListDetailsSchema is a function that returns the schema for SdrListDetails
func SdrListDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"port": schema.Int64Attribute{
			MarkdownDescription: "port",
			Description:         "port",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "protection_domain_id",
			Description:         "protection_domain_id",
			Computed:            true,
		},
		"software_version_info": schema.StringAttribute{
			MarkdownDescription: "software_version_info",
			Description:         "software_version_info",
			Computed:            true,
		},
		"sdr_state": schema.StringAttribute{
			MarkdownDescription: "sdr_state",
			Description:         "sdr_state",
			Computed:            true,
		},
		"membership_state": schema.StringAttribute{
			MarkdownDescription: "membership_state",
			Description:         "membership_state",
			Computed:            true,
		},
		"mdm_connection_state": schema.StringAttribute{
			MarkdownDescription: "mdm_connection_state",
			Description:         "mdm_connection_state",
			Computed:            true,
		},
		"maintenance_state": schema.StringAttribute{
			MarkdownDescription: "maintenance_state",
			Description:         "maintenance_state",
			Computed:            true,
		},
		"perf_profile": schema.StringAttribute{
			MarkdownDescription: "perf_profile",
			Description:         "perf_profile",
			Computed:            true,
		},
		"ip_list": schema.ListNestedAttribute{
			MarkdownDescription: "ip_list",
			Description:         "ip_list",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: IPListSchema()},
		},
	}
}

// AccelerationPoolSchema is a function that returns the schema for AccelerationPool
func AccelerationPoolSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "protection_domain_id",
			Description:         "protection_domain_id",
			Computed:            true,
		},
		"media_type": schema.StringAttribute{
			MarkdownDescription: "media_type",
			Description:         "media_type",
			Computed:            true,
		},
		"rfcache": schema.BoolAttribute{
			MarkdownDescription: "rfcache",
			Description:         "rfcache",
			Computed:            true,
		},
	}
}

// ProtectionDomainSettingsSchema is a function that returns the schema for ProtectionDomainSettings
func ProtectionDomainSettingsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"general": schema.SingleNestedAttribute{
			MarkdownDescription: "general",
			Description:         "general",
			Computed:            true,
			Attributes:          GeneralSchema(),
		},
		"statistics": schema.SingleNestedAttribute{
			MarkdownDescription: "statistics",
			Description:         "statistics",
			Computed:            true,
			Attributes:          StatisticsDetailsSchema(),
		},
		"storage_pool_list": schema.ListNestedAttribute{
			MarkdownDescription: "storage_pool_list",
			Description:         "storage_pool_list",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: StoragePoolListSchema()},
		},
		"sds_list": schema.ListNestedAttribute{
			MarkdownDescription: "sds_list",
			Description:         "sds_list",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: SdsListDetailsSchema()},
		},
		"sdr_list": schema.ListNestedAttribute{
			MarkdownDescription: "sdr_list",
			Description:         "sdr_list",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: SdrListDetailsSchema()},
		},
		"acceleration_pool": schema.ListNestedAttribute{
			MarkdownDescription: "acceleration_pool",
			Description:         "acceleration_pool",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: AccelerationPoolSchema()},
		},
	}
}

// FaultSetSettingsSchema is a function that returns the schema for FaultSetSettings
func FaultSetSettingsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "protection_domain_id",
			Description:         "protection_domain_id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
	}
}

// DatacenterSchema is a function that returns the schema for Datacenter
func DatacenterSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"vcenter_id": schema.StringAttribute{
			MarkdownDescription: "vcenter_id",
			Description:         "vcenter_id",
			Computed:            true,
		},
		"datacenter_id": schema.StringAttribute{
			MarkdownDescription: "datacenter_id",
			Description:         "datacenter_id",
			Computed:            true,
		},
		"datacenter_name": schema.StringAttribute{
			MarkdownDescription: "datacenter_name",
			Description:         "datacenter_name",
			Computed:            true,
		},
	}
}

// PortGroupOptionsSchema is a function that returns the schema for PortGroupOptions
func PortGroupOptionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
	}
}

// PortGroupsSchema is a function that returns the schema for PortGroups
func PortGroupsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "display_name",
			Description:         "display_name",
			Computed:            true,
		},
		"vlan": schema.Int64Attribute{
			MarkdownDescription: "vlan",
			Description:         "vlan",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "value",
			Description:         "value",
			Computed:            true,
		},
		"port_group_options": schema.ListNestedAttribute{
			MarkdownDescription: "port_group_options",
			Description:         "port_group_options",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: PortGroupOptionsSchema()},
		},
	}
}

// VdsSettingsSchema is a function that returns the schema for VdsSettings
func VdsSettingsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "display_name",
			Description:         "display_name",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "value",
			Description:         "value",
			Computed:            true,
		},
		"port_groups": schema.ListNestedAttribute{
			MarkdownDescription: "port_groups",
			Description:         "port_groups",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: PortGroupsSchema()},
		},
	}
}

// VdsNetworkMtuSizeConfigurationSchema is a function that returns the schema for VdsNetworkMtuSizeConfiguration
func VdsNetworkMtuSizeConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "value",
			Description:         "value",
			Computed:            true,
		},
	}
}

// VdsNetworkMTUSizeConfigurationSchema is a function that returns the schema for VdsNetworkMTUSizeConfiguration
func VdsNetworkMTUSizeConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "value",
			Description:         "value",
			Computed:            true,
		},
	}
}

// VdsConfigurationSchema is a function that returns the schema for VdsConfiguration
func VdsConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"datacenter": schema.SingleNestedAttribute{
			MarkdownDescription: "datacenter",
			Description:         "datacenter",
			Computed:            true,
			Attributes:          DatacenterSchema(),
		},
		"port_group_option": schema.StringAttribute{
			MarkdownDescription: "port_group_option",
			Description:         "port_group_option",
			Computed:            true,
		},
		"port_group_creation_option": schema.StringAttribute{
			MarkdownDescription: "port_group_creation_option",
			Description:         "port_group_creation_option",
			Computed:            true,
		},
		"vds_settings": schema.ListNestedAttribute{
			MarkdownDescription: "vds_settings",
			Description:         "vds_settings",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: VdsSettingsSchema()},
		},
		"vds_network_mtu_size_configuration": schema.ListNestedAttribute{
			MarkdownDescription: "vds_network_mtu_size_configuration",
			Description:         "vds_network_mtu_size_configuration",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: VdsNetworkMtuSizeConfigurationSchema()},
		},
	}
}

// NodeSelectionSchema is a function that returns the schema for NodeSelection
func NodeSelectionSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"service_tag": schema.StringAttribute{
			MarkdownDescription: "service_tag",
			Description:         "service_tag",
			Computed:            true,
		},
		"mgmt_ip_address": schema.StringAttribute{
			MarkdownDescription: "mgmt_ip_address",
			Description:         "mgmt_ip_address",
			Computed:            true,
		},
	}
}

// ParametersDetailsSchema is a function that returns the schema for ParametersDetails
func ParametersDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"guid": schema.StringAttribute{
			MarkdownDescription: "guid",
			Description:         "guid",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "type",
			Description:         "type",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "display_name",
			Description:         "display_name",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "value",
			Description:         "value",
			Computed:            true,
		},
		"tool_tip": schema.StringAttribute{
			MarkdownDescription: "tool_tip",
			Description:         "tool_tip",
			Computed:            true,
		},
		"required": schema.BoolAttribute{
			MarkdownDescription: "required",
			Description:         "required",
			Computed:            true,
		},
		"required_at_deployment": schema.BoolAttribute{
			MarkdownDescription: "required_at_deployment",
			Description:         "required_at_deployment",
			Computed:            true,
		},
		"hide_from_template": schema.BoolAttribute{
			MarkdownDescription: "hide_from_template",
			Description:         "hide_from_template",
			Computed:            true,
		},
		"dependencies": schema.ListNestedAttribute{
			MarkdownDescription: "dependencies",
			Description:         "dependencies",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DependenciesDetailsSchema()},
		},
		"group": schema.StringAttribute{
			MarkdownDescription: "group",
			Description:         "group",
			Computed:            true,
		},
		"read_only": schema.BoolAttribute{
			MarkdownDescription: "read_only",
			Description:         "read_only",
			Computed:            true,
		},
		"generated": schema.BoolAttribute{
			MarkdownDescription: "generated",
			Description:         "generated",
			Computed:            true,
		},
		"info_icon": schema.BoolAttribute{
			MarkdownDescription: "info_icon",
			Description:         "info_icon",
			Computed:            true,
		},
		"step": schema.Int64Attribute{
			MarkdownDescription: "step",
			Description:         "step",
			Computed:            true,
		},
		"max_length": schema.Int64Attribute{
			MarkdownDescription: "max_length",
			Description:         "max_length",
			Computed:            true,
		},
		"min": schema.Int64Attribute{
			MarkdownDescription: "min",
			Description:         "min",
			Computed:            true,
		},
		"max": schema.Int64Attribute{
			MarkdownDescription: "max",
			Description:         "max",
			Computed:            true,
		},
		"network_ip_address_list": schema.ListNestedAttribute{
			MarkdownDescription: "network_ip_address_list",
			Description:         "network_ip_address_list",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: NetworkIPAddressListSchema()},
		},
		"network_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "network_configuration",
			Description:         "network_configuration",
			Computed:            true,
			Attributes:          NetworkConfigurationSchema(),
		},
		"raid_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "raid_configuration",
			Description:         "raid_configuration",
			Computed:            true,
			Attributes:          RaidConfigurationSchema(),
		},
		"options": schema.ListNestedAttribute{
			MarkdownDescription: "options",
			Description:         "options",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: OptionsDetailsSchema()},
		},
		"options_sortable": schema.BoolAttribute{
			MarkdownDescription: "options_sortable",
			Description:         "options_sortable",
			Computed:            true,
		},
		"preserved_for_deployment": schema.BoolAttribute{
			MarkdownDescription: "preserved_for_deployment",
			Description:         "preserved_for_deployment",
			Computed:            true,
		},
		"scale_io_disk_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "scale_io_disk_configuration",
			Description:         "scale_io_disk_configuration",
			Computed:            true,
			Attributes:          ScaleIODiskConfigurationSchema(),
		},
		"protection_domain_settings": schema.ListNestedAttribute{
			MarkdownDescription: "protection_domain_settings",
			Description:         "protection_domain_settings",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ProtectionDomainSettingsSchema()},
		},
		"fault_set_settings": schema.ListNestedAttribute{
			MarkdownDescription: "fault_set_settings",
			Description:         "fault_set_settings",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: FaultSetSettingsSchema()},
		},
		"attributes": schema.SingleNestedAttribute{
			MarkdownDescription: "attributes",
			Description:         "attributes",
			Computed:            true,
			Attributes:          AttributesSchema(),
		},
		"vds_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "vds_configuration",
			Description:         "vds_configuration",
			Computed:            true,
			Attributes:          VdsConfigurationSchema(),
		},
		"node_selection": schema.SingleNestedAttribute{
			MarkdownDescription: "node_selection",
			Description:         "node_selection",
			Computed:            true,
			Attributes:          NodeSelectionSchema(),
		},
	}
}

// ResourcesSchema is a function that returns the schema for Resources
func ResourcesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"guid": schema.StringAttribute{
			MarkdownDescription: "guid",
			Description:         "guid",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "display_name",
			Description:         "display_name",
			Computed:            true,
		},
	}
}

// ComponentsSchema is a function that returns the schema for Components
func ComponentsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"component_id": schema.StringAttribute{
			MarkdownDescription: "component_id",
			Description:         "component_id",
			Computed:            true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "identifier",
			Description:         "identifier",
			Computed:            true,
		},
		"component_valid": schema.SingleNestedAttribute{
			MarkdownDescription: "component_valid",
			Description:         "component_valid",
			Computed:            true,
			Attributes:          ComponentValidSchema(),
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"help_text": schema.StringAttribute{
			MarkdownDescription: "help_text",
			Description:         "help_text",
			Computed:            true,
		},
		"cloned_from_id": schema.StringAttribute{
			MarkdownDescription: "cloned_from_id",
			Description:         "cloned_from_id",
			Computed:            true,
		},
		"teardown": schema.BoolAttribute{
			MarkdownDescription: "teardown",
			Description:         "teardown",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "type",
			Description:         "type",
			Computed:            true,
		},
		"sub_type": schema.StringAttribute{
			MarkdownDescription: "sub_type",
			Description:         "sub_type",
			Computed:            true,
		},
		"related_components": schema.SingleNestedAttribute{
			MarkdownDescription: "related_components",
			Description:         "related_components",
			Computed:            true,
			Attributes:          RelatedComponentsSchema(),
		},
		"resources": schema.ListNestedAttribute{
			MarkdownDescription: "resources",
			Description:         "resources",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ResourcesSchema()},
		},
		"brownfield": schema.BoolAttribute{
			MarkdownDescription: "brownfield",
			Description:         "brownfield",
			Computed:            true,
		},
		"puppet_cert_name": schema.StringAttribute{
			MarkdownDescription: "puppet_cert_name",
			Description:         "puppet_cert_name",
			Computed:            true,
		},
		"os_puppet_cert_name": schema.StringAttribute{
			MarkdownDescription: "os_puppet_cert_name",
			Description:         "os_puppet_cert_name",
			Computed:            true,
		},
		"management_ip_address": schema.StringAttribute{
			MarkdownDescription: "management_ip_address",
			Description:         "management_ip_address",
			Computed:            true,
		},
		"serial_number": schema.StringAttribute{
			MarkdownDescription: "serial_number",
			Description:         "serial_number",
			Computed:            true,
		},
		"asm_guid": schema.StringAttribute{
			MarkdownDescription: "asm_guid",
			Description:         "asm_guid",
			Computed:            true,
		},
		"cloned": schema.BoolAttribute{
			MarkdownDescription: "cloned",
			Description:         "cloned",
			Computed:            true,
		},
		"config_file": schema.StringAttribute{
			MarkdownDescription: "config_file",
			Description:         "config_file",
			Computed:            true,
		},
		"manage_firmware": schema.BoolAttribute{
			MarkdownDescription: "manage_firmware",
			Description:         "manage_firmware",
			Computed:            true,
		},
		"instances": schema.Int64Attribute{
			MarkdownDescription: "instances",
			Description:         "instances",
			Computed:            true,
		},
		"ref_id": schema.StringAttribute{
			MarkdownDescription: "ref_id",
			Description:         "ref_id",
			Computed:            true,
		},
		"cloned_from_asm_guid": schema.StringAttribute{
			MarkdownDescription: "cloned_from_asm_guid",
			Description:         "cloned_from_asm_guid",
			Computed:            true,
		},
		"changed": schema.BoolAttribute{
			MarkdownDescription: "changed",
			Description:         "changed",
			Computed:            true,
		},
		"ip": schema.StringAttribute{
			MarkdownDescription: "ip",
			Description:         "ip",
			Computed:            true,
		},
	}
}

// IPRangeSchema is a function that returns the schema for IPRange
func IPRangeSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"starting_ip": schema.StringAttribute{
			MarkdownDescription: "starting_ip",
			Description:         "starting_ip",
			Computed:            true,
		},
		"ending_ip": schema.StringAttribute{
			MarkdownDescription: "ending_ip",
			Description:         "ending_ip",
			Computed:            true,
		},
		"role": schema.StringAttribute{
			MarkdownDescription: "role",
			Description:         "role",
			Computed:            true,
		},
	}
}

// StaticRouteSchema is a function that returns the schema for StaticRoute
func StaticRouteSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"static_route_source_network_id": schema.StringAttribute{
			MarkdownDescription: "static_route_source_network_id",
			Description:         "static_route_source_network_id",
			Computed:            true,
		},
		"static_route_destination_network_id": schema.StringAttribute{
			MarkdownDescription: "static_route_destination_network_id",
			Description:         "static_route_destination_network_id",
			Computed:            true,
		},
		"static_route_gateway": schema.StringAttribute{
			MarkdownDescription: "static_route_gateway",
			Description:         "static_route_gateway",
			Computed:            true,
		},
		"subnet_mask": schema.StringAttribute{
			MarkdownDescription: "subnet_mask",
			Description:         "subnet_mask",
			Computed:            true,
		},
		"destination_ip_address": schema.StringAttribute{
			MarkdownDescription: "destination_ip_address",
			Description:         "destination_ip_address",
			Computed:            true,
		},
	}
}

// StaticNetworkConfigurationSchema is a function that returns the schema for StaticNetworkConfiguration
func StaticNetworkConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"gateway": schema.StringAttribute{
			MarkdownDescription: "gateway",
			Description:         "gateway",
			Computed:            true,
		},
		"subnet": schema.StringAttribute{
			MarkdownDescription: "subnet",
			Description:         "subnet",
			Computed:            true,
		},
		"primary_dns": schema.StringAttribute{
			MarkdownDescription: "primary_dns",
			Description:         "primary_dns",
			Computed:            true,
		},
		"secondary_dns": schema.StringAttribute{
			MarkdownDescription: "secondary_dns",
			Description:         "secondary_dns",
			Computed:            true,
		},
		"dns_suffix": schema.StringAttribute{
			MarkdownDescription: "dns_suffix",
			Description:         "dns_suffix",
			Computed:            true,
		},
		"ip_range": schema.ListNestedAttribute{
			MarkdownDescription: "ip_range",
			Description:         "ip_range",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: IPRangeSchema()},
		},
		"ip_address": schema.StringAttribute{
			MarkdownDescription: "ip_address",
			Description:         "ip_address",
			Computed:            true,
		},
		"static_route": schema.ListNestedAttribute{
			MarkdownDescription: "static_route",
			Description:         "static_route",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: StaticRouteSchema()},
		},
	}
}

// NetworksSchema is a function that returns the schema for Networks
func NetworksSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "description",
			Description:         "description",
			Computed:            true,
		},
		"vlan_id": schema.Int64Attribute{
			MarkdownDescription: "vlan_id",
			Description:         "vlan_id",
			Computed:            true,
		},
		"static_network_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "static_network_configuration",
			Description:         "static_network_configuration",
			Computed:            true,
			Attributes:          StaticNetworkConfigurationSchema(),
		},
		"destination_ip_address": schema.StringAttribute{
			MarkdownDescription: "destination_ip_address",
			Description:         "destination_ip_address",
			Computed:            true,
		},
		"static": schema.BoolAttribute{
			MarkdownDescription: "static",
			Description:         "static",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "type",
			Description:         "type",
			Computed:            true,
		},
	}
}

// OptionsSchema is a function that returns the schema for Options
func OptionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "name",
			Description:         "name",
			Computed:            true,
		},
		"dependencies": schema.ListNestedAttribute{
			MarkdownDescription: "dependencies",
			Description:         "dependencies",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DependenciesDetailsSchema()},
		},
		"attributes": schema.SingleNestedAttribute{
			MarkdownDescription: "attributes",
			Description:         "attributes",
			Computed:            true,
			Attributes:          AttributesSchema(),
		},
	}
}

// ParametersSchema is a function that returns the schema for Parameters
func ParametersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "value",
			Description:         "value",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "display_name",
			Description:         "display_name",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "type",
			Description:         "type",
			Computed:            true,
		},
		"tool_tip": schema.StringAttribute{
			MarkdownDescription: "tool_tip",
			Description:         "tool_tip",
			Computed:            true,
		},
		"required": schema.BoolAttribute{
			MarkdownDescription: "required",
			Description:         "required",
			Computed:            true,
		},
		"hide_from_template": schema.BoolAttribute{
			MarkdownDescription: "hide_from_template",
			Description:         "hide_from_template",
			Computed:            true,
		},
		"device_type": schema.StringAttribute{
			MarkdownDescription: "device_type",
			Description:         "device_type",
			Computed:            true,
		},
		"dependencies": schema.ListNestedAttribute{
			MarkdownDescription: "dependencies",
			Description:         "dependencies",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DependenciesDetailsSchema()},
		},
		"group": schema.StringAttribute{
			MarkdownDescription: "group",
			Description:         "group",
			Computed:            true,
		},
		"read_only": schema.BoolAttribute{
			MarkdownDescription: "read_only",
			Description:         "read_only",
			Computed:            true,
		},
		"generated": schema.BoolAttribute{
			MarkdownDescription: "generated",
			Description:         "generated",
			Computed:            true,
		},
		"info_icon": schema.BoolAttribute{
			MarkdownDescription: "info_icon",
			Description:         "info_icon",
			Computed:            true,
		},
		"step": schema.Int64Attribute{
			MarkdownDescription: "step",
			Description:         "step",
			Computed:            true,
		},
		"max_length": schema.Int64Attribute{
			MarkdownDescription: "max_length",
			Description:         "max_length",
			Computed:            true,
		},
		"min": schema.Int64Attribute{
			MarkdownDescription: "min",
			Description:         "min",
			Computed:            true,
		},
		"max": schema.Int64Attribute{
			MarkdownDescription: "max",
			Description:         "max",
			Computed:            true,
		},
		"networks": schema.ListNestedAttribute{
			MarkdownDescription: "networks",
			Description:         "networks",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: NetworksSchema()},
		},
		"options": schema.ListNestedAttribute{
			MarkdownDescription: "options",
			Description:         "options",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: OptionsSchema()},
		},
		"options_sortable": schema.BoolAttribute{
			MarkdownDescription: "options_sortable",
			Description:         "options_sortable",
			Computed:            true,
		},
	}
}

// CategoriesSchema is a function that returns the schema for Categories
func CategoriesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "id",
			Description:         "id",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "display_name",
			Description:         "display_name",
			Computed:            true,
		},
		"device_type": schema.StringAttribute{
			MarkdownDescription: "device_type",
			Description:         "device_type",
			Computed:            true,
		},
	}
}
