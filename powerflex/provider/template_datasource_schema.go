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
	"terraform-provider-powerflex/powerflex/helper"
	"terraform-provider-powerflex/powerflex/models"

	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// TemplateDataSourceSchema defines the schema for template datasource
var TemplateDataSourceSchema schema.Schema = schema.Schema{
	Description:         "This datasource is used to query the existing templates from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	MarkdownDescription: "This datasource is used to query the existing templates from the PowerFlex array. The information fetched from this datasource can be used for getting the details / for further processing in resource block.",
	Attributes: map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "Placeholder attribute.",
			MarkdownDescription: "Placeholder attribute.",
			Computed:            true,
		},
		"template_details": schema.SetNestedAttribute{
			Description:         "Template details",
			MarkdownDescription: "Template details",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: TemplateDetailSchema()},
		},
	},
	Blocks: map[string]schema.Block{
		"filter": schema.SingleNestedBlock{
			Attributes: helper.GenerateSchemaAttributes(helper.TypeToMap(models.TemplateFilter{})),
		},
	},
}

// TemplateDetailSchema is a function that returns the schema for TemplateDetails
func TemplateDetailSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the template.",
			Description:         "Template ID",
			Computed:            true,
		},
		"template_name": schema.StringAttribute{
			MarkdownDescription: "The name of the template.",
			Description:         "Template Name",
			Computed:            true,
		},
		"template_description": schema.StringAttribute{
			MarkdownDescription: "The description of the template.",
			Description:         "Template Description",
			Computed:            true,
		},
		"template_type": schema.StringAttribute{
			MarkdownDescription: "The type/category of the template.",
			Description:         "Template Type",
			Computed:            true,
		},
		"template_version": schema.StringAttribute{
			MarkdownDescription: "The version of the template.",
			Description:         "Template Version",
			Computed:            true,
		},
		"original_template_id": schema.StringAttribute{
			MarkdownDescription: "The ID of the original template if this is a derived template.",
			Description:         "Original Template ID",
			Computed:            true,
		},
		"template_valid": schema.SingleNestedAttribute{
			MarkdownDescription: "Details about the validity of the template.",
			Description:         "Template Validity",
			Computed:            true,
			Attributes:          TemplateValidSchema(),
		},
		"template_locked": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the template is locked or not.",
			Description:         "Template Lock Status",
			Computed:            true,
		},
		"in_configuration": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the template is part of the current configuration.",
			Description:         "In Configuration",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "The date when the template was created.",
			Description:         "Creation Date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The user who created the template.",
			Description:         "Created By",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "The date when the template was last updated.",
			Description:         "Last Update Date",
			Computed:            true,
		},
		"last_deployed_date": schema.StringAttribute{
			MarkdownDescription: "The date when the template was last deployed.",
			Description:         "Last Deployed Date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The user who last updated the template.",
			Description:         "Last Updated By",
			Computed:            true,
		},
		"manage_firmware": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether firmware is managed by the template.",
			Description:         "Manage Firmware",
			Computed:            true,
		},
		"use_default_catalog": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the default catalog is used for the template.",
			Description:         "Use Default Catalog",
			Computed:            true,
		},
		"firmware_repository": schema.SingleNestedAttribute{
			MarkdownDescription: "Details about the firmware repository used by the template.",
			Description:         "Firmware Repository",
			Computed:            true,
			Attributes:          FirmwareRepositorySchema(),
		},
		"license_repository": schema.SingleNestedAttribute{
			MarkdownDescription: "Details about the license repository used by the template.",
			Description:         "License Repository",
			Computed:            true,
			Attributes:          LicenseRepositorySchema(),
		},
		"assigned_users": schema.ListNestedAttribute{
			MarkdownDescription: "List of users assigned to the template.",
			Description:         "Assigned Users",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: AssignedUsersSchema()},
		},
		"all_users_allowed": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether all users are allowed for the template.",
			Description:         "All Users Allowed",
			Computed:            true,
		},
		"category": schema.StringAttribute{
			MarkdownDescription: "The category to which the template belongs.",
			Description:         "Template Category",
			Computed:            true,
		},
		"components": schema.ListNestedAttribute{
			MarkdownDescription: "List of components included in the template.",
			Description:         "Template Components",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ComponentsSchema()},
		},
		"configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "Details about the configuration settings of the template.",
			Description:         "Template Configuration",
			Computed:            true,
			Attributes:          ConfigurationDetailsSchema(),
		},
		"server_count": schema.Int64Attribute{
			MarkdownDescription: "The count of servers associated with the template.",
			Description:         "Server Count",
			Computed:            true,
		},
		"storage_count": schema.Int64Attribute{
			MarkdownDescription: "The count of storage devices associated with the template.",
			Description:         "Storage Count",
			Computed:            true,
		},
		"cluster_count": schema.Int64Attribute{
			MarkdownDescription: "The count of clusters associated with the template.",
			Description:         "Cluster Count",
			Computed:            true,
		},
		"service_count": schema.Int64Attribute{
			MarkdownDescription: "The count of services associated with the template.",
			Description:         "Service Count",
			Computed:            true,
		},
		"switch_count": schema.Int64Attribute{
			MarkdownDescription: "The count of switches associated with the template.",
			Description:         "Switch Count",
			Computed:            true,
		},
		"vm_count": schema.Int64Attribute{
			MarkdownDescription: "The count of virtual machines associated with the template.",
			Description:         "Virtual Machine Count",
			Computed:            true,
		},
		"sdnas_count": schema.Int64Attribute{
			MarkdownDescription: "The count of software-defined network appliances associated with the template.",
			Description:         "SDNAs Count",
			Computed:            true,
		},
		"brownfield_template_type": schema.StringAttribute{
			MarkdownDescription: "The type of template for brownfield deployments.",
			Description:         "Brownfield Template Type",
			Computed:            true,
		},
		"networks": schema.ListNestedAttribute{
			MarkdownDescription: "List of networks associated with the template.",
			Description:         "Networks",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: NetworksSchema()},
		},
		"draft": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the template is in draft mode.",
			Description:         "Draft Mode",
			Computed:            true,
		},
	}
}

// MessagesSchema is a function that returns the schema for Messages
func MessagesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the message.",
			Description:         "Message ID",
			Computed:            true,
		},
		"message_code": schema.StringAttribute{
			MarkdownDescription: "The code associated with the message.",
			Description:         "Message Code",
			Computed:            true,
		},
		"message_bundle": schema.StringAttribute{
			MarkdownDescription: "The bundle or group to which the message belongs.",
			Description:         "Message Bundle",
			Computed:            true,
		},
		"severity": schema.StringAttribute{
			MarkdownDescription: "The severity level of the message (e.g., INFO, WARNING, ERROR).",
			Description:         "Message Severity",
			Computed:            true,
		},
		"category": schema.StringAttribute{
			MarkdownDescription: "The category or type of the message.",
			Description:         "Message Category",
			Computed:            true,
		},
		"display_message": schema.StringAttribute{
			MarkdownDescription: "The message to be displayed or shown.",
			Description:         "Display Message",
			Computed:            true,
		},
		"response_action": schema.StringAttribute{
			MarkdownDescription: "The action to be taken in response to the message.",
			Description:         "Response Action",
			Computed:            true,
		},
		"detailed_message": schema.StringAttribute{
			MarkdownDescription: "A detailed version or description of the message.",
			Description:         "Detailed Message",
			Computed:            true,
		},
		"correlation_id": schema.StringAttribute{
			MarkdownDescription: "The identifier used to correlate related messages.",
			Description:         "Correlation ID",
			Computed:            true,
		},
		"agent_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the agent associated with the message.",
			Description:         "Agent ID",
			Computed:            true,
		},
		"time_stamp": schema.StringAttribute{
			MarkdownDescription: "The timestamp indicating when the message was generated.",
			Description:         "Timestamp",
			Computed:            true,
		},
		"sequence_number": schema.Int64Attribute{
			MarkdownDescription: "The sequence number of the message in a series.",
			Description:         "Sequence Number",
			Computed:            true,
		},
	}
}

// TemplateValidSchema is a function that returns the schema for TemplateValid
func TemplateValidSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"valid": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the template is valid.",
			Description:         "Template Validity",
			Computed:            true,
		},
		"messages": schema.ListNestedAttribute{
			MarkdownDescription: "List of messages associated with the template validity.",
			Description:         "Validation Messages",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: MessagesSchema()},
		},
	}
}

// SoftwareComponentsSchema is a function that returns the schema for SoftwareComponents
func SoftwareComponentsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the software component.",
			Description:         "Software Component ID",
			Computed:            true,
		},
		"package_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the package to which the component belongs.",
			Description:         "Package ID",
			Computed:            true,
		},
		"dell_version": schema.StringAttribute{
			MarkdownDescription: "The version of the component according to Dell standards.",
			Description:         "Dell Version",
			Computed:            true,
		},
		"vendor_version": schema.StringAttribute{
			MarkdownDescription: "The version of the component according to the vendor's standards.",
			Description:         "Vendor Version",
			Computed:            true,
		},
		"component_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the component.",
			Description:         "Component ID",
			Computed:            true,
		},
		"device_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the device associated with the component.",
			Description:         "Device ID",
			Computed:            true,
		},
		"sub_device_id": schema.StringAttribute{
			MarkdownDescription: "The sub-identifier of the device associated with the component.",
			Description:         "Sub-Device ID",
			Computed:            true,
		},
		"vendor_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the vendor associated with the component.",
			Description:         "Vendor ID",
			Computed:            true,
		},
		"sub_vendor_id": schema.StringAttribute{
			MarkdownDescription: "The sub-identifier of the vendor associated with the component.",
			Description:         "Sub-Vendor ID",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "The date when the component was created.",
			Description:         "Creation Date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The user who created the component.",
			Description:         "Created By",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "The date when the component was last updated.",
			Description:         "Last Update Date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The user who last updated the component.",
			Description:         "Last Updated By",
			Computed:            true,
		},
		"path": schema.StringAttribute{
			MarkdownDescription: "The path where the component is stored.",
			Description:         "Component Path",
			Computed:            true,
		},
		"hash_md_5": schema.StringAttribute{
			MarkdownDescription: "The MD5 hash value of the component.",
			Description:         "MD5 Hash",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the software component.",
			Description:         "Component Name",
			Computed:            true,
		},
		"category": schema.StringAttribute{
			MarkdownDescription: "The category to which the component belongs.",
			Description:         "Component Category",
			Computed:            true,
		},
		"component_type": schema.StringAttribute{
			MarkdownDescription: "The type of the component.",
			Description:         "Component Type",
			Computed:            true,
		},
		"operating_system": schema.StringAttribute{
			MarkdownDescription: "The operating system associated with the component.",
			Description:         "Operating System",
			Computed:            true,
		},
		"system_ids": schema.ListAttribute{
			MarkdownDescription: "List of system IDs associated with the component.",
			Description:         "System IDs",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"custom": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the component is custom or not.",
			Description:         "Custom Component",
			Computed:            true,
		},
		"needs_attention": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the component needs attention.",
			Description:         "Needs Attention",
			Computed:            true,
		},
		"ignore": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the component should be ignored.",
			Description:         "Ignore Component",
			Computed:            true,
		},
		"original_version": schema.StringAttribute{
			MarkdownDescription: "The original version of the component.",
			Description:         "Original Version",
			Computed:            true,
		},
		"original_component_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the original component.",
			Description:         "Original Component ID",
			Computed:            true,
		},
		"firmware_repo_name": schema.StringAttribute{
			MarkdownDescription: "The name of the firmware repository associated with the component.",
			Description:         "Firmware Repository Name",
			Computed:            true,
		},
	}
}

// SoftwareBundlesSchema is a function that returns the schema for SoftwareBundles
func SoftwareBundlesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the software bundle.",
			Description:         "Software Bundle ID",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the software bundle.",
			Description:         "Bundle Name",
			Computed:            true,
		},
		"version": schema.StringAttribute{
			MarkdownDescription: "The version of the software bundle.",
			Description:         "Bundle Version",
			Computed:            true,
		},
		"bundle_date": schema.StringAttribute{
			MarkdownDescription: "The date when the software bundle was created.",
			Description:         "Bundle Creation Date",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "The date when the software bundle was initially created.",
			Description:         "Creation Date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The user who initially created the software bundle.",
			Description:         "Created By",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "The date when the software bundle was last updated.",
			Description:         "Last Update Date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The user who last updated the software bundle.",
			Description:         "Last Updated By",
			Computed:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "A brief description of the software bundle.",
			Description:         "Bundle Description",
			Computed:            true,
		},
		"user_bundle": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the software bundle is a user-specific bundle.",
			Description:         "User-Specific Bundle",
			Computed:            true,
		},
		"user_bundle_path": schema.StringAttribute{
			MarkdownDescription: "The path associated with the user-specific software bundle.",
			Description:         "User Bundle Path",
			Computed:            true,
		},
		"user_bundle_hash_md_5": schema.StringAttribute{
			MarkdownDescription: "The MD5 hash value of the user-specific software bundle.",
			Description:         "User Bundle MD5 Hash",
			Computed:            true,
		},
		"device_type": schema.StringAttribute{
			MarkdownDescription: "The type of device associated with the software bundle.",
			Description:         "Device Type",
			Computed:            true,
		},
		"device_model": schema.StringAttribute{
			MarkdownDescription: "The model of the device associated with the software bundle.",
			Description:         "Device Model",
			Computed:            true,
		},
		"criticality": schema.StringAttribute{
			MarkdownDescription: "The criticality level of the software bundle.",
			Description:         "Criticality Level",
			Computed:            true,
		},
		"fw_repository_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the firmware repository associated with the software bundle.",
			Description:         "Firmware Repository ID",
			Computed:            true,
		},
		"bundle_type": schema.StringAttribute{
			MarkdownDescription: "The type of the software bundle.",
			Description:         "Bundle Type",
			Computed:            true,
		},
		"custom": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the software bundle is custom.",
			Description:         "Custom Bundle",
			Computed:            true,
		},
		"needs_attention": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the software bundle needs attention.",
			Description:         "Needs Attention",
			Computed:            true,
		},
		"software_components": schema.ListNestedAttribute{
			MarkdownDescription: "List of software components associated with the software bundle.",
			Description:         "Software Components",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: SoftwareComponentsSchema()},
		},
	}
}

// DeploymentValidSchema is a function that returns the schema for DeploymentValid
func DeploymentValidSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"valid": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment is valid.",
			Description:         "Deployment Validity",
			Computed:            true,
		},
		"messages": schema.ListNestedAttribute{
			MarkdownDescription: "List of messages related to the deployment.",
			Description:         "Deployment Messages",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: MessagesSchema()},
		},
	}
}

// DeploymentDeviceSchema is a function that returns the schema for DeploymentDevice
func DeploymentDeviceSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"ref_id": schema.StringAttribute{
			MarkdownDescription: "The reference ID associated with the deployment device.",
			Description:         "Reference ID",
			Computed:            true,
		},
		"ref_type": schema.StringAttribute{
			MarkdownDescription: "The reference type associated with the deployment device.",
			Description:         "Reference Type",
			Computed:            true,
		},
		"log_dump": schema.StringAttribute{
			MarkdownDescription: "The log dump information associated with the deployment device.",
			Description:         "Log Dump",
			Computed:            true,
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "The status of the deployment device.",
			Description:         "Device Status",
			Computed:            true,
		},
		"status_end_time": schema.StringAttribute{
			MarkdownDescription: "The end time of the status for the deployment device.",
			Description:         "Status End Time",
			Computed:            true,
		},
		"status_start_time": schema.StringAttribute{
			MarkdownDescription: "The start time of the status for the deployment device.",
			Description:         "Status Start Time",
			Computed:            true,
		},
		"device_health": schema.StringAttribute{
			MarkdownDescription: "The health status of the deployment device.",
			Description:         "Device Health",
			Computed:            true,
		},
		"health_message": schema.StringAttribute{
			MarkdownDescription: "The health message associated with the deployment device.",
			Description:         "Health Message",
			Computed:            true,
		},
		"compliant_state": schema.StringAttribute{
			MarkdownDescription: "The compliant state of the deployment device.",
			Description:         "Compliant State",
			Computed:            true,
		},
		"brownfield_status": schema.StringAttribute{
			MarkdownDescription: "The brownfield status of the deployment device.",
			Description:         "Brownfield Status",
			Computed:            true,
		},
		"device_type": schema.StringAttribute{
			MarkdownDescription: "The type of device associated with the deployment device.",
			Description:         "Device Type",
			Computed:            true,
		},
		"device_group_name": schema.StringAttribute{
			MarkdownDescription: "The name of the device group associated with the deployment device.",
			Description:         "Device Group Name",
			Computed:            true,
		},
		"ip_address": schema.StringAttribute{
			MarkdownDescription: "The IP address of the deployment device.",
			Description:         "IP Address",
			Computed:            true,
		},
		"current_ip_address": schema.StringAttribute{
			MarkdownDescription: "The current IP address of the deployment device.",
			Description:         "Current IP Address",
			Computed:            true,
		},
		"service_tag": schema.StringAttribute{
			MarkdownDescription: "The service tag associated with the deployment device.",
			Description:         "Service Tag",
			Computed:            true,
		},
		"component_id": schema.StringAttribute{
			MarkdownDescription: "The component ID associated with the deployment device.",
			Description:         "Component ID",
			Computed:            true,
		},
		"status_message": schema.StringAttribute{
			MarkdownDescription: "The status message associated with the deployment device.",
			Description:         "Status Message",
			Computed:            true,
		},
		"model": schema.StringAttribute{
			MarkdownDescription: "The model of the deployment device.",
			Description:         "Device Model",
			Computed:            true,
		},
		"cloud_link": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment device has a cloud link.",
			Description:         "Cloud Link",
			Computed:            true,
		},
		"das_cache": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment device has Direct-Attached Storage (DAS) cache.",
			Description:         "DAS Cache",
			Computed:            true,
		},
		"device_state": schema.StringAttribute{
			MarkdownDescription: "The state of the deployment device.",
			Description:         "Device State",
			Computed:            true,
		},
		"puppet_cert_name": schema.StringAttribute{
			MarkdownDescription: "The Puppet certificate name associated with the deployment device.",
			Description:         "Puppet Certificate Name",
			Computed:            true,
		},
		"brownfield": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment device is associated with a brownfield deployment.",
			Description:         "Brownfield Deployment",
			Computed:            true,
		},
	}
}

// VmsSchema is a function that returns the schema for Vms
func VmsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"certificate_name": schema.StringAttribute{
			MarkdownDescription: "The certificate name associated with the virtual machine (VM).",
			Description:         "Certificate Name",
			Computed:            true,
		},
		"vm_model": schema.StringAttribute{
			MarkdownDescription: "The model of the virtual machine (VM).",
			Description:         "VM Model",
			Computed:            true,
		},
		"vm_ipaddress": schema.StringAttribute{
			MarkdownDescription: "The IP address of the virtual machine (VM).",
			Description:         "VM IP Address",
			Computed:            true,
		},
		"vm_manufacturer": schema.StringAttribute{
			MarkdownDescription: "The manufacturer of the virtual machine (VM).",
			Description:         "VM Manufacturer",
			Computed:            true,
		},
		"vm_service_tag": schema.StringAttribute{
			MarkdownDescription: "The service tag associated with the virtual machine (VM).",
			Description:         "VM Service Tag",
			Computed:            true,
		},
	}
}

// LicenseRepositorySchema is a function that returns the schema for LicenseRepository
func LicenseRepositorySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the license repository.",
			Description:         "License Repository ID",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the license repository.",
			Description:         "License Repository Name",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "The type of the license repository.",
			Description:         "License Repository Type",
			Computed:            true,
		},
		"disk_location": schema.StringAttribute{
			MarkdownDescription: "The disk location of the license repository.",
			Description:         "License Repository Disk Location",
			Computed:            true,
		},
		"filename": schema.StringAttribute{
			MarkdownDescription: "The filename associated with the license repository.",
			Description:         "License Repository Filename",
			Computed:            true,
		},
		"state": schema.StringAttribute{
			MarkdownDescription: "The state of the license repository.",
			Description:         "License Repository State",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "The date when the license repository was created.",
			Description:         "License Repository Created Date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The user who created the license repository.",
			Description:         "License Repository Created By",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "The date when the license repository was last updated.",
			Description:         "License Repository Updated Date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The user who last updated the license repository.",
			Description:         "License Repository Updated By",
			Computed:            true,
		},
		"license_data": schema.StringAttribute{
			MarkdownDescription: "The license data associated with the license repository.",
			Description:         "License Repository Data",
			Computed:            true,
		},
	}
}

// AssignedUsersSchema is a function that returns the schema for AssignedUsers
func AssignedUsersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"user_seq_id": schema.Int64Attribute{
			MarkdownDescription: "The sequential ID of the assigned user.",
			Description:         "User Sequential ID",
			Computed:            true,
		},
		"user_name": schema.StringAttribute{
			MarkdownDescription: "The username of the assigned user.",
			Description:         "Username",
			Computed:            true,
		},
		"password": schema.StringAttribute{
			MarkdownDescription: "The password associated with the assigned user.",
			Description:         "User Password",
			Computed:            true,
		},
		"update_password": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the user password needs to be updated.",
			Description:         "Update User Password",
			Computed:            true,
		},
		"domain_name": schema.StringAttribute{
			MarkdownDescription: "The domain name of the assigned user.",
			Description:         "User Domain Name",
			Computed:            true,
		},
		"group_dn": schema.StringAttribute{
			MarkdownDescription: "The distinguished name (DN) of the group associated with the assigned user.",
			Description:         "Group DN",
			Computed:            true,
		},
		"group_name": schema.StringAttribute{
			MarkdownDescription: "The name of the group associated with the assigned user.",
			Description:         "Group Name",
			Computed:            true,
		},
		"first_name": schema.StringAttribute{
			MarkdownDescription: "The first name of the assigned user.",
			Description:         "User First Name",
			Computed:            true,
		},
		"last_name": schema.StringAttribute{
			MarkdownDescription: "The last name of the assigned user.",
			Description:         "User Last Name",
			Computed:            true,
		},
		"email": schema.StringAttribute{
			MarkdownDescription: "The email address of the assigned user.",
			Description:         "User Email",
			Computed:            true,
		},
		"phone_number": schema.StringAttribute{
			MarkdownDescription: "The phone number of the assigned user.",
			Description:         "User Phone Number",
			Computed:            true,
		},
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the assigned user is enabled.",
			Description:         "User Enabled",
			Computed:            true,
		},
		"system_user": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the assigned user is a system user.",
			Description:         "System User",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "The date when the assigned user was created.",
			Description:         "User Created Date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The user who created the assigned user.",
			Description:         "User Created By",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "The date when the assigned user was last updated.",
			Description:         "User Updated Date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The user who last updated the assigned user.",
			Description:         "User Updated By",
			Computed:            true,
		},
		"role": schema.StringAttribute{
			MarkdownDescription: "The role associated with the assigned user.",
			Description:         "User Role",
			Computed:            true,
		},
		"user_preference": schema.StringAttribute{
			MarkdownDescription: "The preferences of the assigned user.",
			Description:         "User Preferences",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the assigned user.",
			Description:         "User ID",
			Computed:            true,
		},
		"roles": schema.ListAttribute{
			MarkdownDescription: "The roles associated with the assigned user.",
			Description:         "User Roles",
			Computed:            true,
			ElementType:         types.StringType,
		},
	}
}

// JobDetailsSchema is a function that returns the schema for JobDetails
func JobDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"level": schema.StringAttribute{
			MarkdownDescription: "The log level of the job.",
			Description:         "Log Level",
			Computed:            true,
		},
		"message": schema.StringAttribute{
			MarkdownDescription: "The log message of the job.",
			Description:         "Log Message",
			Computed:            true,
		},
		"timestamp": schema.StringAttribute{
			MarkdownDescription: "The timestamp of the job execution.",
			Description:         "Timestamp",
			Computed:            true,
		},
		"execution_id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the job execution.",
			Description:         "Execution ID",
			Computed:            true,
		},
		"component_id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the component associated with the job.",
			Description:         "Component ID",
			Computed:            true,
		},
	}
}

// DeploymentValidationResponseSchema is a function that returns the schema for DeploymentValidationResponse
func DeploymentValidationResponseSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"nodes": schema.Int64Attribute{
			MarkdownDescription: "The number of nodes in the deployment.",
			Description:         "Number of Nodes",
			Computed:            true,
		},
		"storage_pools": schema.Int64Attribute{
			MarkdownDescription: "The number of storage pools in the deployment.",
			Description:         "Number of Storage Pools",
			Computed:            true,
		},
		"drives_per_storage_pool": schema.Int64Attribute{
			MarkdownDescription: "The number of drives per storage pool in the deployment.",
			Description:         "Drives per Storage Pool",
			Computed:            true,
		},
		"max_scalability": schema.Int64Attribute{
			MarkdownDescription: "The maximum scalability of the deployment.",
			Description:         "Maximum Scalability",
			Computed:            true,
		},
		"virtual_machines": schema.Int64Attribute{
			MarkdownDescription: "The number of virtual machines in the deployment.",
			Description:         "Number of Virtual Machines",
			Computed:            true,
		},
		"number_of_service_volumes": schema.Int64Attribute{
			MarkdownDescription: "The number of service volumes in the deployment.",
			Description:         "Number of Service Volumes",
			Computed:            true,
		},
		"can_deploy": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment can be executed.",
			Description:         "Can Deploy",
			Computed:            true,
		},
		"warning_messages": schema.ListAttribute{
			MarkdownDescription: "A list of warning messages associated with the deployment validation.",
			Description:         "Warning Messages",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"storage_pool_disk_type": schema.ListAttribute{
			MarkdownDescription: "The disk types associated with each storage pool in the deployment.",
			Description:         "Storage Pool Disk Types",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"hostnames": schema.ListAttribute{
			MarkdownDescription: "A list of hostnames associated with the deployment.",
			Description:         "Hostnames",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"new_node_disk_types": schema.ListAttribute{
			MarkdownDescription: "The disk types associated with new nodes in the deployment.",
			Description:         "New Node Disk Types",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"no_of_fault_sets": schema.Int64Attribute{
			MarkdownDescription: "The number of fault sets in the deployment.",
			Description:         "Number of Fault Sets",
			Computed:            true,
		},
		"nodes_per_fault_set": schema.Int64Attribute{
			MarkdownDescription: "The number of nodes per fault set in the deployment.",
			Description:         "Nodes per Fault Set",
			Computed:            true,
		},
		"protection_domain": schema.StringAttribute{
			MarkdownDescription: "The protection domain associated with the deployment.",
			Description:         "Protection Domain",
			Computed:            true,
		},
		"disk_type_mismatch": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether there is a disk type mismatch in the deployment.",
			Description:         "Disk Type Mismatch",
			Computed:            true,
		},
	}
}

// DeploymentsSchema is a function that returns the schema for Deployments
func DeploymentsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the deployment.",
			Description:         "ID",
			Computed:            true,
		},
		"deployment_name": schema.StringAttribute{
			MarkdownDescription: "The name of the deployment.",
			Description:         "Deployment Name",
			Computed:            true,
		},
		"deployment_description": schema.StringAttribute{
			MarkdownDescription: "The description of the deployment.",
			Description:         "Deployment Description",
			Computed:            true,
		},
		"deployment_valid": schema.SingleNestedAttribute{
			MarkdownDescription: "Details about the validity of the deployment.",
			Description:         "Deployment Validity",
			Computed:            true,
			Attributes:          DeploymentValidSchema(),
		},
		"retry": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment should be retried.",
			Description:         "Retry",
			Computed:            true,
		},
		"teardown": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment should be torn down.",
			Description:         "Teardown",
			Computed:            true,
		},
		"teardown_after_cancel": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether teardown should occur after canceling the deployment.",
			Description:         "Teardown After Cancel",
			Computed:            true,
		},
		"remove_service": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the associated service should be removed.",
			Description:         "Remove Service",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "The date when the deployment was created.",
			Description:         "Created Date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The user who created the deployment.",
			Description:         "Created By",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "The date when the deployment was last updated.",
			Description:         "Updated Date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The user who last updated the deployment.",
			Description:         "Updated By",
			Computed:            true,
		},
		"deployment_scheduled_date": schema.StringAttribute{
			MarkdownDescription: "The scheduled date for the deployment.",
			Description:         "Deployment Scheduled Date",
			Computed:            true,
		},
		"deployment_started_date": schema.StringAttribute{
			MarkdownDescription: "The date when the deployment started.",
			Description:         "Deployment Started Date",
			Computed:            true,
		},
		"deployment_finished_date": schema.StringAttribute{
			MarkdownDescription: "The date when the deployment finished.",
			Description:         "Deployment Finished Date",
			Computed:            true,
		},
		"schedule_date": schema.StringAttribute{
			MarkdownDescription: "The date when the deployment is scheduled.",
			Description:         "Schedule Date",
			Computed:            true,
		},
		"status": schema.StringAttribute{
			MarkdownDescription: "The status of the deployment.",
			Description:         "Status",
			Computed:            true,
		},
		"compliant": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment is compliant.",
			Description:         "Compliant",
			Computed:            true,
		},
		"deployment_device": schema.ListNestedAttribute{
			MarkdownDescription: "List of devices associated with the deployment.",
			Description:         "Deployment Devices",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DeploymentDeviceSchema()},
		},
		"vms": schema.ListNestedAttribute{
			MarkdownDescription: "List of virtual machines associated with the deployment.",
			Description:         "Virtual Machines",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: VmsSchema()},
		},
		"update_server_firmware": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to update server firmware during the deployment.",
			Description:         "Update Server Firmware",
			Computed:            true,
		},
		"use_default_catalog": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to use the default catalog for the deployment.",
			Description:         "Use Default Catalog",
			Computed:            true,
		},
		"firmware_repository_id": schema.StringAttribute{
			MarkdownDescription: "The ID of the firmware repository associated with the deployment.",
			Description:         "Firmware Repository ID",
			Computed:            true,
		},
		"license_repository": schema.SingleNestedAttribute{
			MarkdownDescription: "Details about the license repository associated with the deployment.",
			Description:         "License Repository",
			Computed:            true,
			Attributes:          LicenseRepositorySchema(),
		},
		"license_repository_id": schema.StringAttribute{
			MarkdownDescription: "The ID of the license repository associated with the deployment.",
			Description:         "License Repository ID",
			Computed:            true,
		},
		"individual_teardown": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to perform individual teardown for the deployment.",
			Description:         "Individual Teardown",
			Computed:            true,
		},
		"deployment_health_status_type": schema.StringAttribute{
			MarkdownDescription: "The type of health status associated with the deployment.",
			Description:         "Deployment Health Status Type",
			Computed:            true,
		},
		"assigned_users": schema.ListNestedAttribute{
			MarkdownDescription: "List of users assigned to the deployment.",
			Description:         "Assigned Users",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: AssignedUsersSchema()},
		},
		"all_users_allowed": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether all users are allowed for the deployment.",
			Description:         "All Users Allowed",
			Computed:            true,
		},
		"owner": schema.StringAttribute{
			MarkdownDescription: "The owner of the deployment.",
			Description:         "Owner",
			Computed:            true,
		},
		"no_op": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment is a no-op (no operation).",
			Description:         "No Operation",
			Computed:            true,
		},
		"firmware_init": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether firmware initialization is performed during deployment.",
			Description:         "Firmware Initialization",
			Computed:            true,
		},
		"disruptive_firmware": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether disruptive firmware actions are allowed.",
			Description:         "Disruptive Firmware",
			Computed:            true,
		},
		"preconfigure_svm": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to preconfigure SVM (Storage Virtual Machine).",
			Description:         "Preconfigure SVM",
			Computed:            true,
		},
		"preconfigure_svm_and_update": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to preconfigure SVM and perform an update.",
			Description:         "Preconfigure SVM and Update",
			Computed:            true,
		},
		"services_deployed": schema.StringAttribute{
			MarkdownDescription: "Details about the services deployed during the deployment.",
			Description:         "Services Deployed",
			Computed:            true,
		},
		"precalculated_device_health": schema.StringAttribute{
			MarkdownDescription: "The precalculated health of devices associated with the deployment.",
			Description:         "Precalculated Device Health",
			Computed:            true,
		},
		"lifecycle_mode_reasons": schema.ListAttribute{
			MarkdownDescription: "List of reasons for the lifecycle mode of the deployment.",
			Description:         "Lifecycle Mode Reasons",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"job_details": schema.ListNestedAttribute{
			MarkdownDescription: "List of job details associated with the deployment.",
			Description:         "Job Details",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: JobDetailsSchema()},
		},
		"number_of_deployments": schema.Int64Attribute{
			MarkdownDescription: "The total number of deployments.",
			Description:         "Number of Deployments",
			Computed:            true,
		},
		"operation_type": schema.StringAttribute{
			MarkdownDescription: "The type of operation associated with the deployment.",
			Description:         "Operation Type",
			Computed:            true,
		},
		"operation_status": schema.StringAttribute{
			MarkdownDescription: "The status of the operation associated with the deployment.",
			Description:         "Operation Status",
			Computed:            true,
		},
		"operation_data": schema.StringAttribute{
			MarkdownDescription: "Additional data associated with the operation.",
			Description:         "Operation Data",
			Computed:            true,
		},
		"deployment_validation_response": schema.SingleNestedAttribute{
			MarkdownDescription: "Details about the validation response for the deployment.",
			Description:         "Deployment Validation Response",
			Computed:            true,
			Attributes:          DeploymentValidationResponseSchema(),
		},
		"current_step_count": schema.StringAttribute{
			MarkdownDescription: "The current step count during deployment.",
			Description:         "Current Step Count",
			Computed:            true,
		},
		"total_num_of_steps": schema.StringAttribute{
			MarkdownDescription: "The total number of steps involved in the deployment.",
			Description:         "Total Number of Steps",
			Computed:            true,
		},
		"current_step_message": schema.StringAttribute{
			MarkdownDescription: "The message associated with the current step during deployment.",
			Description:         "Current Step Message",
			Computed:            true,
		},
		"custom_image": schema.StringAttribute{
			MarkdownDescription: "The custom image used for deployment.",
			Description:         "Custom Image",
			Computed:            true,
		},
		"original_deployment_id": schema.StringAttribute{
			MarkdownDescription: "The ID of the original deployment.",
			Description:         "Original Deployment ID",
			Computed:            true,
		},
		"current_batch_count": schema.StringAttribute{
			MarkdownDescription: "The current batch count during deployment.",
			Description:         "Current Batch Count",
			Computed:            true,
		},
		"total_batch_count": schema.StringAttribute{
			MarkdownDescription: "The total number of batches involved in the deployment.",
			Description:         "Total Batch Count",
			Computed:            true,
		},
		"brownfield": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment involves brownfield operations.",
			Description:         "Brownfield",
			Computed:            true,
		},
		"scale_up": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment involves scaling up.",
			Description:         "Scale Up",
			Computed:            true,
		},
		"lifecycle_mode": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the deployment is in lifecycle mode.",
			Description:         "Lifecycle Mode",
			Computed:            true,
		},
		"overall_device_health": schema.StringAttribute{
			MarkdownDescription: "The overall health status of the device in the deployment.",
			Description:         "Overall Device Health",
			Computed:            true,
		},
		"vds": schema.BoolAttribute{
			MarkdownDescription: "Specifies whether the deployment involves Virtual Desktop Infrastructure (VDI) configuration.",
			Description:         "VDI Configuration",
			Computed:            true,
		},
		"template_valid": schema.BoolAttribute{
			MarkdownDescription: "Indicates if the deployment template is valid.",
			Description:         "Template Validity",
			Computed:            true,
		},
		"configuration_change": schema.BoolAttribute{
			MarkdownDescription: "Specifies whether there has been a change in the deployment configuration.",
			Description:         "Configuration Change",
			Computed:            true,
		},
		"can_migratev_clsv_ms": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether migration of cluster virtual machines is allowed.",
			Description:         "Cluster VM Migration Allowed",
			Computed:            true,
		},
	}
}

// FirmwareRepositorySchema is a function that returns the schema for FirmwareRepository
func FirmwareRepositorySchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the firmware repository.",
			Description:         "ID",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the firmware repository.",
			Description:         "Name",
			Computed:            true,
		},
		"source_location": schema.StringAttribute{
			MarkdownDescription: "The location of the source for the firmware repository.",
			Description:         "Source Location",
			Computed:            true,
		},
		"source_type": schema.StringAttribute{
			MarkdownDescription: "The type of the source for the firmware repository.",
			Description:         "Source Type",
			Computed:            true,
		},
		"disk_location": schema.StringAttribute{
			MarkdownDescription: "The location on disk where the firmware repository is stored.",
			Description:         "Disk Location",
			Computed:            true,
		},
		"filename": schema.StringAttribute{
			MarkdownDescription: "The filename of the firmware repository.",
			Description:         "Filename",
			Computed:            true,
		},
		"md_5_hash": schema.StringAttribute{
			MarkdownDescription: "The MD5 hash of the firmware repository.",
			Description:         "MD5 Hash",
			Computed:            true,
		},
		"username": schema.StringAttribute{
			MarkdownDescription: "The username associated with the firmware repository.",
			Description:         "Username",
			Computed:            true,
		},
		"password": schema.StringAttribute{
			MarkdownDescription: "The password associated with the firmware repository.",
			Description:         "Password",
			Computed:            true,
		},
		"download_status": schema.StringAttribute{
			MarkdownDescription: "The download status of the firmware repository.",
			Description:         "Download Status",
			Computed:            true,
		},
		"created_date": schema.StringAttribute{
			MarkdownDescription: "The date when the firmware repository was created.",
			Description:         "Created Date",
			Computed:            true,
		},
		"created_by": schema.StringAttribute{
			MarkdownDescription: "The user who created the firmware repository.",
			Description:         "Created By",
			Computed:            true,
		},
		"updated_date": schema.StringAttribute{
			MarkdownDescription: "The date when the firmware repository was last updated.",
			Description:         "Updated Date",
			Computed:            true,
		},
		"updated_by": schema.StringAttribute{
			MarkdownDescription: "The user who last updated the firmware repository.",
			Description:         "Updated By",
			Computed:            true,
		},
		"default_catalog": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the firmware repository is the default catalog.",
			Description:         "Default Catalog",
			Computed:            true,
		},
		"embedded": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the firmware repository is embedded.",
			Description:         "Embedded",
			Computed:            true,
		},
		"state": schema.StringAttribute{
			MarkdownDescription: "The state of the firmware repository.",
			Description:         "State",
			Computed:            true,
		},
		"software_components": schema.ListNestedAttribute{
			MarkdownDescription: "List of software components associated with the firmware repository.",
			Description:         "Software Components",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: SoftwareComponentsSchema()},
		},
		"software_bundles": schema.ListNestedAttribute{
			MarkdownDescription: "List of software bundles associated with the firmware repository.",
			Description:         "Software Bundles",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: SoftwareBundlesSchema()},
		},
		"deployments": schema.ListNestedAttribute{
			MarkdownDescription: "List of deployments associated with the firmware repository.",
			Description:         "Deployments",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DeploymentsSchema()},
		},
		"bundle_count": schema.Int64Attribute{
			MarkdownDescription: "The count of software bundles in the firmware repository.",
			Description:         "Bundle Count",
			Computed:            true,
		},
		"component_count": schema.Int64Attribute{
			MarkdownDescription: "The count of software components in the firmware repository.",
			Description:         "Component Count",
			Computed:            true,
		},
		"user_bundle_count": schema.Int64Attribute{
			MarkdownDescription: "The count of user-specific software bundles in the firmware repository.",
			Description:         "User Bundle Count",
			Computed:            true,
		},
		"minimal": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the firmware repository is minimal.",
			Description:         "Minimal",
			Computed:            true,
		},
		"download_progress": schema.Int64Attribute{
			MarkdownDescription: "The progress of the download for the firmware repository.",
			Description:         "Download Progress",
			Computed:            true,
		},
		"extract_progress": schema.Int64Attribute{
			MarkdownDescription: "The progress of the extraction for the firmware repository.",
			Description:         "Extract Progress",
			Computed:            true,
		},
		"file_size_in_gigabytes": schema.Int64Attribute{
			MarkdownDescription: "The size of the firmware repository file in gigabytes.",
			Description:         "File Size (in Gigabytes)",
			Computed:            true,
		},
		"signed_key_source_location": schema.StringAttribute{
			MarkdownDescription: "The source location of the signed key associated with the firmware repository.",
			Description:         "Signed Key Source Location",
			Computed:            true,
		},
		"signature": schema.StringAttribute{
			MarkdownDescription: "The signature of the firmware repository.",
			Description:         "Signature",
			Computed:            true,
		},
		"custom": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the firmware repository is custom.",
			Description:         "Custom",
			Computed:            true,
		},
		"needs_attention": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the firmware repository needs attention.",
			Description:         "Needs Attention",
			Computed:            true,
		},
		"job_id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the job associated with the firmware repository.",
			Description:         "Job ID",
			Computed:            true,
		},
		"rcmapproved": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the firmware repository is RCM approved.",
			Description:         "RCM Approved",
			Computed:            true,
		},
	}
}

// ComponentValidSchema is a function that returns the schema for ComponentValid
func ComponentValidSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"valid": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the component is valid.",
			Description:         "Validity",
			Computed:            true,
		},
		"messages": schema.ListNestedAttribute{
			MarkdownDescription: "List of messages associated with the component validity.",
			Description:         "Messages",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: MessagesSchema()},
		},
	}
}

// ConfigurationDetailsSchema is a function that returns the schema for ConfigurationDetails
func ConfigurationDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "Unique identifier for the configuration.",
			Description:         "Configuration ID",
			Computed:            true,
		},
		"disktype": schema.StringAttribute{
			MarkdownDescription: "Type of disk in the configuration.",
			Description:         "Disk Type",
			Computed:            true,
		},
		"comparator": schema.StringAttribute{
			MarkdownDescription: "Comparator used in the configuration.",
			Description:         "Comparator",
			Computed:            true,
		},
		"numberofdisks": schema.Int64Attribute{
			MarkdownDescription: "Number of disks in the configuration.",
			Description:         "Number of Disks",
			Computed:            true,
		},
		"raidlevel": schema.StringAttribute{
			MarkdownDescription: "RAID level of the configuration.",
			Description:         "RAID Level",
			Computed:            true,
		},
		"virtual_disk_fqdd": schema.StringAttribute{
			MarkdownDescription: "Fully Qualified Device Descriptor (FQDD) of the virtual disk in the configuration.",
			Description:         "Virtual Disk FQDD",
			Computed:            true,
		},
		"controller_fqdd": schema.StringAttribute{
			MarkdownDescription: "Fully Qualified Device Descriptor (FQDD) of the controller in the configuration.",
			Description:         "Controller FQDD",
			Computed:            true,
		},
		"categories": schema.ListNestedAttribute{
			MarkdownDescription: "List of categories associated with the configuration.",
			Description:         "Categories",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: CategoriesSchema()},
		},
	}
}

// ResourcesSchema is a function that returns the schema for Resources
func ResourcesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"guid": schema.StringAttribute{
			MarkdownDescription: "The globally unique identifier (GUID) for the resources.",
			Description:         "The globally unique identifier (GUID) for the resources.",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "The identifier for the resources.",
			Description:         "The identifier for the resources.",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "The display name for the resources.",
			Description:         "The display name for the resources.",
			Computed:            true,
		},
	}
}

// ComponentsSchema is a function that returns the schema for Components
func ComponentsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the component.",
			Description:         "The unique identifier for the component.",
			Computed:            true,
		},
		"component_id": schema.StringAttribute{
			MarkdownDescription: "The identifier for the component.",
			Description:         "The identifier for the component.",
			Computed:            true,
		},
		"identifier": schema.StringAttribute{
			MarkdownDescription: "The identifier for the component.",
			Description:         "The identifier for the component.",
			Computed:            true,
		},
		"component_valid": schema.SingleNestedAttribute{
			MarkdownDescription: "Information about the validity of the component.",
			Description:         "Information about the validity of the component.",
			Computed:            true,
			Attributes:          ComponentValidSchema(),
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the component.",
			Description:         "The name of the component.",
			Computed:            true,
		},
		"help_text": schema.StringAttribute{
			MarkdownDescription: "Help text associated with the component.",
			Description:         "Help text associated with the component.",
			Computed:            true,
		},
		"cloned_from_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the component from which this component is cloned.",
			Description:         "The identifier of the component from which this component is cloned.",
			Computed:            true,
		},
		"teardown": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the component should be torn down.",
			Description:         "Indicates whether the component should be torn down.",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "The type of the component.",
			Description:         "The type of the component.",
			Computed:            true,
		},
		"sub_type": schema.StringAttribute{
			MarkdownDescription: "The sub-type of the component.",
			Description:         "The sub-type of the component.",
			Computed:            true,
		},
		"related_components": schema.MapAttribute{
			MarkdownDescription: "Related components associated with this component.",
			Description:         "Related components associated with this component.",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"resources": schema.ListNestedAttribute{
			MarkdownDescription: "List of resources associated with the component.",
			Description:         "List of resources associated with the component.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ResourcesSchema()},
		},
		"brownfield": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the component is brownfield.",
			Description:         "Indicates whether the component is brownfield.",
			Computed:            true,
		},
		"puppet_cert_name": schema.StringAttribute{
			MarkdownDescription: "The Puppet certificate name associated with the component.",
			Description:         "The Puppet certificate name associated with the component.",
			Computed:            true,
		},
		"os_puppet_cert_name": schema.StringAttribute{
			MarkdownDescription: "The OS Puppet certificate name associated with the component.",
			Description:         "The OS Puppet certificate name associated with the component.",
			Computed:            true,
		},
		"management_ip_address": schema.StringAttribute{
			MarkdownDescription: "The management IP address of the component.",
			Description:         "The management IP address of the component.",
			Computed:            true,
		},
		"serial_number": schema.StringAttribute{
			MarkdownDescription: "The serial number of the component.",
			Description:         "The serial number of the component.",
			Computed:            true,
		},
		"asm_guid": schema.StringAttribute{
			MarkdownDescription: "The ASM GUID (Global Unique Identifier) associated with the component.",
			Description:         "The ASM GUID (Global Unique Identifier) associated with the component.",
			Computed:            true,
		},
		"cloned": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the component is cloned.",
			Description:         "Indicates whether the component is cloned.",
			Computed:            true,
		},
		"config_file": schema.StringAttribute{
			MarkdownDescription: "The configuration file associated with the component.",
			Description:         "The configuration file associated with the component.",
			Computed:            true,
		},
		"manage_firmware": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether firmware is managed for the component.",
			Description:         "Indicates whether firmware is managed for the component.",
			Computed:            true,
		},
		"instances": schema.Int64Attribute{
			MarkdownDescription: "The number of instances of the component.",
			Description:         "The number of instances of the component.",
			Computed:            true,
		},
		"ref_id": schema.StringAttribute{
			MarkdownDescription: "The reference identifier associated with the component.",
			Description:         "The reference identifier associated with the component.",
			Computed:            true,
		},
		"cloned_from_asm_guid": schema.StringAttribute{
			MarkdownDescription: "The ASM GUID from which the component is cloned.",
			Description:         "The ASM GUID from which the component is cloned.",
			Computed:            true,
		},
		"changed": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the component has changed.",
			Description:         "Indicates whether the component has changed.",
			Computed:            true,
		},
		"ip": schema.StringAttribute{
			MarkdownDescription: "The IP address associated with the component.",
			Description:         "The IP address associated with the component.",
			Computed:            true,
		},
	}
}

// IPRangeSchema is a function that returns the schema for IPRange
func IPRangeSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the IP range.",
			Description:         "The unique identifier for the IP range.",
			Computed:            true,
		},
		"starting_ip": schema.StringAttribute{
			MarkdownDescription: "The starting IP address of the range.",
			Description:         "The starting IP address of the range.",
			Computed:            true,
		},
		"ending_ip": schema.StringAttribute{
			MarkdownDescription: "The ending IP address of the range.",
			Description:         "The ending IP address of the range.",
			Computed:            true,
		},
		"role": schema.StringAttribute{
			MarkdownDescription: "The role associated with the IP range.",
			Description:         "The role associated with the IP range.",
			Computed:            true,
		},
	}
}

// StaticRouteSchema is a function that returns the schema for StaticRoute
func StaticRouteSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"static_route_source_network_id": schema.StringAttribute{
			MarkdownDescription: "The ID of the source network for the static route.",
			Description:         "The ID of the source network for the static route.",
			Computed:            true,
		},
		"static_route_destination_network_id": schema.StringAttribute{
			MarkdownDescription: "The ID of the destination network for the static route.",
			Description:         "The ID of the destination network for the static route.",
			Computed:            true,
		},
		"static_route_gateway": schema.StringAttribute{
			MarkdownDescription: "The gateway for the static route.",
			Description:         "The gateway for the static route.",
			Computed:            true,
		},
		"subnet_mask": schema.StringAttribute{
			MarkdownDescription: "The subnet mask for the static route.",
			Description:         "The subnet mask for the static route.",
			Computed:            true,
		},
		"destination_ip_address": schema.StringAttribute{
			MarkdownDescription: "The IP address of the destination for the static route.",
			Description:         "The IP address of the destination for the static route.",
			Computed:            true,
		},
	}
}

// StaticNetworkConfigurationSchema is a function that returns the schema for StaticNetworkConfiguration
func StaticNetworkConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"gateway": schema.StringAttribute{
			MarkdownDescription: "The gateway for the static network configuration.",
			Description:         "The gateway for the static network configuration.",
			Computed:            true,
		},
		"subnet": schema.StringAttribute{
			MarkdownDescription: "The subnet for the static network configuration.",
			Description:         "The subnet for the static network configuration.",
			Computed:            true,
		},
		"primary_dns": schema.StringAttribute{
			MarkdownDescription: "The primary DNS server for the static network configuration.",
			Description:         "The primary DNS server for the static network configuration.",
			Computed:            true,
		},
		"secondary_dns": schema.StringAttribute{
			MarkdownDescription: "The secondary DNS server for the static network configuration.",
			Description:         "The secondary DNS server for the static network configuration.",
			Computed:            true,
		},
		"dns_suffix": schema.StringAttribute{
			MarkdownDescription: "The DNS suffix for the static network configuration.",
			Description:         "The DNS suffix for the static network configuration.",
			Computed:            true,
		},
		"ip_range": schema.ListNestedAttribute{
			MarkdownDescription: "List of IP ranges associated with the static network configuration.",
			Description:         "List of IP ranges associated with the static network configuration.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: IPRangeSchema()},
		},
		"ip_address": schema.StringAttribute{
			MarkdownDescription: "The IP address associated with the static network configuration.",
			Description:         "The IP address associated with the static network configuration.",
			Computed:            true,
		},
		"static_route": schema.ListNestedAttribute{
			MarkdownDescription: "List of static routes associated with the static network configuration.",
			Description:         "List of static routes associated with the static network configuration.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: StaticRouteSchema()},
		},
	}
}

// NetworksSchema is a function that returns the schema for Networks
func NetworksSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the network.",
			Description:         "The unique identifier for the network.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the network.",
			Description:         "The name of the network.",
			Computed:            true,
		},
		"description": schema.StringAttribute{
			MarkdownDescription: "The description of the network.",
			Description:         "The description of the network.",
			Computed:            true,
		},
		"vlan_id": schema.Int64Attribute{
			MarkdownDescription: "The VLAN ID associated with the network.",
			Description:         "The VLAN ID associated with the network.",
			Computed:            true,
		},
		"static_network_configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "The static network configuration settings.",
			Description:         "The static network configuration settings.",
			Computed:            true,
			Attributes:          StaticNetworkConfigurationSchema(),
		},
		"destination_ip_address": schema.StringAttribute{
			MarkdownDescription: "The destination IP address for the network.",
			Description:         "The destination IP address for the network.",
			Computed:            true,
		},
		"static": schema.BoolAttribute{
			MarkdownDescription: "Boolean indicating if the network is static.",
			Description:         "Boolean indicating if the network is static.",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "The type of the network.",
			Description:         "The type of the network.",
			Computed:            true,
		},
	}
}

// CategoriesSchema is a function that returns the schema for Categories
func CategoriesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the category.",
			Description:         "The unique identifier for the category.",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "The display name of the category.",
			Description:         "The display name of the category.",
			Computed:            true,
		},
		"device_type": schema.StringAttribute{
			MarkdownDescription: "The type of device associated with the category.",
			Description:         "The type of device associated with the category.",
			Computed:            true,
		},
	}
}
