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
					path.MatchRoot("template_names"),
				),
			},
		},
		"template_names": schema.SetAttribute{
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

// RelatedComponentsSchema is a function that returns the schema for RelatedComponents
func RelatedComponentsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"additional_prop_1": schema.StringAttribute{
			MarkdownDescription: "Additional property 1 for related components.",
			Description:         "Additional Property 1",
			Computed:            true,
		},
		"additional_prop_2": schema.StringAttribute{
			MarkdownDescription: "Additional property 2 for related components.",
			Description:         "Additional Property 2",
			Computed:            true,
		},
		"additional_prop_3": schema.StringAttribute{
			MarkdownDescription: "Additional property 3 for related components.",
			Description:         "Additional Property 3",
			Computed:            true,
		},
	}
}

// DependenciesDetailsSchema is a function that returns the schema for DependenciesDetails
func DependenciesDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the dependency.",
			Description:         "ID",
			Computed:            true,
		},
		"dependency_target": schema.StringAttribute{
			MarkdownDescription: "The target of the dependency.",
			Description:         "Dependency Target",
			Computed:            true,
		},
		"dependency_value": schema.StringAttribute{
			MarkdownDescription: "The value of the dependency.",
			Description:         "Dependency Value",
			Computed:            true,
		},
	}
}

// NetworkIPAddressListSchema is a function that returns the schema for NetworkIPAddressList
func NetworkIPAddressListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the network IP address.",
			Description:         "ID",
			Computed:            true,
		},
		"ip_address": schema.StringAttribute{
			MarkdownDescription: "The IP address in the network IP address list.",
			Description:         "IP Address",
			Computed:            true,
		},
	}
}

// PartitionsSchema is a function that returns the schema for Partitions
func PartitionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the partition.",
			Description:         "ID",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the partition.",
			Description:         "Name",
			Computed:            true,
		},
		"networks": schema.ListAttribute{
			MarkdownDescription: "List of networks associated with the partition.",
			Description:         "Networks",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"network_ip_address_list": schema.ListNestedAttribute{
			MarkdownDescription: "List of network IP addresses associated with the partition.",
			Description:         "Network IP Address List",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: NetworkIPAddressListSchema()},
		},
		"minimum": schema.Int64Attribute{
			MarkdownDescription: "The minimum value associated with the partition.",
			Description:         "Minimum",
			Computed:            true,
		},
		"maximum": schema.Int64Attribute{
			MarkdownDescription: "The maximum value associated with the partition.",
			Description:         "Maximum",
			Computed:            true,
		},
		"lan_mac_address": schema.StringAttribute{
			MarkdownDescription: "The LAN MAC address associated with the partition.",
			Description:         "LAN MAC Address",
			Computed:            true,
		},
		"iscsi_mac_address": schema.StringAttribute{
			MarkdownDescription: "The iSCSI MAC address associated with the partition.",
			Description:         "iSCSI MAC Address",
			Computed:            true,
		},
		"iscsi_iqn": schema.StringAttribute{
			MarkdownDescription: "The iSCSI IQN associated with the partition.",
			Description:         "iSCSI IQN",
			Computed:            true,
		},
		"wwnn": schema.StringAttribute{
			MarkdownDescription: "The WWNN associated with the partition.",
			Description:         "WWNN",
			Computed:            true,
		},
		"wwpn": schema.StringAttribute{
			MarkdownDescription: "The WWPN associated with the partition.",
			Description:         "WWPN",
			Computed:            true,
		},
		"fqdd": schema.StringAttribute{
			MarkdownDescription: "The FQDD (Fully Qualified Device Descriptor) associated with the partition.",
			Description:         "FQDD",
			Computed:            true,
		},
		"mirrored_port": schema.StringAttribute{
			MarkdownDescription: "The mirrored port associated with the partition.",
			Description:         "Mirrored Port",
			Computed:            true,
		},
		"mac_address": schema.StringAttribute{
			MarkdownDescription: "The MAC address associated with the partition.",
			Description:         "MAC Address",
			Computed:            true,
		},
		"port_no": schema.Int64Attribute{
			MarkdownDescription: "The port number associated with the partition.",
			Description:         "Port Number",
			Computed:            true,
		},
		"partition_no": schema.Int64Attribute{
			MarkdownDescription: "The partition number.",
			Description:         "Partition Number",
			Computed:            true,
		},
		"partition_index": schema.Int64Attribute{
			MarkdownDescription: "The partition index.",
			Description:         "Partition Index",
			Computed:            true,
		},
	}
}

// InterfacesSchema is a function that returns the schema for Interfaces
func InterfacesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the interface.",
			Description:         "ID",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the interface.",
			Description:         "Name",
			Computed:            true,
		},
		"partitioned": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the interface is partitioned.",
			Description:         "Partitioned",
			Computed:            true,
		},
		"partitions": schema.ListNestedAttribute{
			MarkdownDescription: "List of partitions associated with the interface.",
			Description:         "Partitions",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: PartitionsSchema()},
		},
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the interface is enabled.",
			Description:         "Enabled",
			Computed:            true,
		},
		"redundancy": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the interface has redundancy.",
			Description:         "Redundancy",
			Computed:            true,
		},
		"nictype": schema.StringAttribute{
			MarkdownDescription: "The type of the network interface.",
			Description:         "NIC Type",
			Computed:            true,
		},
		"fqdd": schema.StringAttribute{
			MarkdownDescription: "The FQDD (Fully Qualified Device Descriptor) of the interface.",
			Description:         "FQDD",
			Computed:            true,
		},
		"max_partitions": schema.Int64Attribute{
			MarkdownDescription: "The maximum number of partitions allowed for the interface.",
			Description:         "Maximum Partitions",
			Computed:            true,
		},
		"all_networks": schema.ListAttribute{
			MarkdownDescription: "List of all networks associated with the interface.",
			Description:         "All Networks",
			Computed:            true,
			ElementType:         types.StringType,
		},
	}
}

// InterfacesDetailsSchema is a function that returns the schema for InterfacesDetails
func InterfacesDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the interface.",
			Description:         "ID",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the interface.",
			Description:         "Name",
			Computed:            true,
		},
		"redundancy": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the interface has redundancy.",
			Description:         "Redundancy",
			Computed:            true,
		},
		"enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the interface is enabled.",
			Description:         "Enabled",
			Computed:            true,
		},
		"partitioned": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the interface is partitioned.",
			Description:         "Partitioned",
			Computed:            true,
		},
		"interfaces": schema.ListNestedAttribute{
			MarkdownDescription: "List of interfaces associated with this interface.",
			Description:         "Interfaces",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: InterfacesSchema()},
		},
		"nictype": schema.StringAttribute{
			MarkdownDescription: "The type of the network interface.",
			Description:         "NIC Type",
			Computed:            true,
		},
		"fabrictype": schema.StringAttribute{
			MarkdownDescription: "The fabric type associated with the interface.",
			Description:         "Fabric Type",
			Computed:            true,
		},
		"max_partitions": schema.Int64Attribute{
			MarkdownDescription: "The maximum number of partitions allowed for the interface.",
			Description:         "Maximum Partitions",
			Computed:            true,
		},
		"nports": schema.Int64Attribute{
			MarkdownDescription: "The number of ports associated with the interface.",
			Description:         "Number of Ports",
			Computed:            true,
		},
		"card_index": schema.Int64Attribute{
			MarkdownDescription: "The card index associated with the interface.",
			Description:         "Card Index",
			Computed:            true,
		},
		"nictype_source": schema.StringAttribute{
			MarkdownDescription: "The source of the NIC type information.",
			Description:         "NIC Type Source",
			Computed:            true,
		},
	}
}

// NetworkConfigurationSchema is a function that returns the schema for NetworkConfiguration
func NetworkConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the network configuration.",
			Description:         "ID",
			Computed:            true,
		},
		"interfaces": schema.ListNestedAttribute{
			MarkdownDescription: "List of interfaces associated with the network configuration.",
			Description:         "Interfaces",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: InterfacesDetailsSchema()},
		},
		"software_only": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the network configuration is for software only.",
			Description:         "Software Only",
			Computed:            true,
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

// VirtualDisksSchema is a function that returns the schema for VirtualDisks
func VirtualDisksSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"physical_disks": schema.ListAttribute{
			MarkdownDescription: "List of physical disks associated with the virtual disk.",
			Description:         "Physical Disks",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"virtual_disk_fqdd": schema.StringAttribute{
			MarkdownDescription: "Fully Qualified Device Descriptor (FQDD) of the virtual disk.",
			Description:         "Virtual Disk FQDD",
			Computed:            true,
		},
		"raid_level": schema.StringAttribute{
			MarkdownDescription: "RAID level of the virtual disk.",
			Description:         "RAID Level",
			Computed:            true,
		},
		"roll_up_status": schema.StringAttribute{
			MarkdownDescription: "Roll-up status of the virtual disk.",
			Description:         "Roll-up Status",
			Computed:            true,
		},
		"controller": schema.StringAttribute{
			MarkdownDescription: "Controller associated with the virtual disk.",
			Description:         "Controller",
			Computed:            true,
		},
		"controller_product_name": schema.StringAttribute{
			MarkdownDescription: "Product name of the controller associated with the virtual disk.",
			Description:         "Controller Product Name",
			Computed:            true,
		},
		"configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration details of the virtual disk.",
			Description:         "Configuration",
			Computed:            true,
			Attributes:          ConfigurationDetailsSchema(),
		},
		"media_type": schema.StringAttribute{
			MarkdownDescription: "Media type of the virtual disk.",
			Description:         "Media Type",
			Computed:            true,
		},
		"encryption_type": schema.StringAttribute{
			MarkdownDescription: "Type of encryption used for the virtual disk.",
			Description:         "Encryption Type",
			Computed:            true,
		},
	}
}

// ExternalVirtualDisksSchema is a function that returns the schema for ExternalVirtualDisks
func ExternalVirtualDisksSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"physical_disks": schema.ListAttribute{
			MarkdownDescription: "List of physical disks associated with the external virtual disks.",
			Description:         "Physical Disks",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"virtual_disk_fqdd": schema.StringAttribute{
			MarkdownDescription: "Fully Qualified Device Descriptor (FQDD) of the external virtual disk.",
			Description:         "Virtual Disk FQDD",
			Computed:            true,
		},
		"raid_level": schema.StringAttribute{
			MarkdownDescription: "RAID level of the external virtual disk.",
			Description:         "RAID Level",
			Computed:            true,
		},
		"roll_up_status": schema.StringAttribute{
			MarkdownDescription: "Roll-up status of the external virtual disk.",
			Description:         "Roll-up Status",
			Computed:            true,
		},
		"controller": schema.StringAttribute{
			MarkdownDescription: "Controller associated with the external virtual disk.",
			Description:         "Controller",
			Computed:            true,
		},
		"controller_product_name": schema.StringAttribute{
			MarkdownDescription: "Product name of the controller associated with the external virtual disk.",
			Description:         "Controller Product Name",
			Computed:            true,
		},
		"configuration": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration details of the external virtual disk.",
			Description:         "Configuration",
			Computed:            true,
			Attributes:          ConfigurationDetailsSchema(),
		},
		"media_type": schema.StringAttribute{
			MarkdownDescription: "Media type of the external virtual disk.",
			Description:         "Media Type",
			Computed:            true,
		},
		"encryption_type": schema.StringAttribute{
			MarkdownDescription: "Encryption type of the external virtual disk.",
			Description:         "Encryption Type",
			Computed:            true,
		},
	}
}

// SizeToDiskMapSchema is a function that returns the schema for SizeToDiskMap
func SizeToDiskMapSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"additional_prop_1": schema.Int64Attribute{
			MarkdownDescription: "Additional property 1 in the size-to-disk map.",
			Description:         "Additional Property 1",
			Computed:            true,
		},
		"additional_prop_2": schema.Int64Attribute{
			MarkdownDescription: "Additional property 2 in the size-to-disk map.",
			Description:         "Additional Property 2",
			Computed:            true,
		},
		"additional_prop_3": schema.Int64Attribute{
			MarkdownDescription: "Additional property 3 in the size-to-disk map.",
			Description:         "Additional Property 3",
			Computed:            true,
		},
	}
}

// RaidConfigurationSchema is a function that returns the schema for RaidConfiguration
func RaidConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"virtual_disks": schema.ListNestedAttribute{
			MarkdownDescription: "List of virtual disks in the RAID configuration.",
			Description:         "Virtual Disks",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: VirtualDisksSchema()},
		},
		"external_virtual_disks": schema.ListNestedAttribute{
			MarkdownDescription: "List of external virtual disks in the RAID configuration.",
			Description:         "External Virtual Disks",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ExternalVirtualDisksSchema()},
		},
		"hdd_hot_spares": schema.ListAttribute{
			MarkdownDescription: "List of HDD hot spares in the RAID configuration.",
			Description:         "HDD Hot Spares",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"ssd_hot_spares": schema.ListAttribute{
			MarkdownDescription: "List of SSD hot spares in the RAID configuration.",
			Description:         "SSD Hot Spares",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"external_hdd_hot_spares": schema.ListAttribute{
			MarkdownDescription: "List of external HDD hot spares in the RAID configuration.",
			Description:         "External HDD Hot Spares",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"external_ssd_hot_spares": schema.ListAttribute{
			MarkdownDescription: "List of external SSD hot spares in the RAID configuration.",
			Description:         "External SSD Hot Spares",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"size_to_disk_map": schema.SingleNestedAttribute{
			MarkdownDescription: "Mapping of size to disks in the RAID configuration.",
			Description:         "Size to Disk Map",
			Computed:            true,
			Attributes:          SizeToDiskMapSchema(),
		},
	}
}

// AttributesSchema is a function that returns the schema for Attributes
func AttributesSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"additional_prop_1": schema.StringAttribute{
			MarkdownDescription: "Additional property 1.",
			Description:         "Additional Property 1",
			Computed:            true,
		},
		"additional_prop_2": schema.StringAttribute{
			MarkdownDescription: "Additional property 2.",
			Description:         "Additional Property 2",
			Computed:            true,
		},
		"additional_prop_3": schema.StringAttribute{
			MarkdownDescription: "Additional property 3.",
			Description:         "Additional Property 3",
			Computed:            true,
		},
	}
}

// OptionsDetailsSchema is a function that returns the schema for OptionsDetails
func OptionsDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "ID of the option.",
			Description:         "Option ID",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "Name of the option.",
			Description:         "Option Name",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "Value of the option.",
			Description:         "Option Value",
			Computed:            true,
		},
		"dependencies": schema.ListNestedAttribute{
			MarkdownDescription: "List of dependencies for the option.",
			Description:         "Option Dependencies",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DependenciesDetailsSchema()},
		},
		"attributes": schema.SingleNestedAttribute{
			MarkdownDescription: "Attributes associated with the option.",
			Description:         "Option Attributes",
			Computed:            true,
			Attributes:          AttributesSchema(),
		},
	}
}

// ScaleIOStoragePoolDisksSchema is a function that returns the schema for ScaleIOStoragePoolDisks
func ScaleIOStoragePoolDisksSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "ID of the protection domain.",
			Description:         "Protection Domain ID",
			Computed:            true,
		},
		"protection_domain_name": schema.StringAttribute{
			MarkdownDescription: "Name of the protection domain.",
			Description:         "Protection Domain Name",
			Computed:            true,
		},
		"storage_pool_id": schema.StringAttribute{
			MarkdownDescription: "ID of the storage pool.",
			Description:         "Storage Pool ID",
			Computed:            true,
		},
		"storage_pool_name": schema.StringAttribute{
			MarkdownDescription: "Name of the storage pool.",
			Description:         "Storage Pool Name",
			Computed:            true,
		},
		"disk_type": schema.StringAttribute{
			MarkdownDescription: "Type of the disk.",
			Description:         "Disk Type",
			Computed:            true,
		},
		"physical_disk_fqdds": schema.ListAttribute{
			MarkdownDescription: "List of fully qualified domain names (FQDDs) of physical disks.",
			Description:         "Physical Disk FQDDs",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"virtual_disk_fqdds": schema.ListAttribute{
			MarkdownDescription: "List of fully qualified domain names (FQDDs) of virtual disks.",
			Description:         "Virtual Disk FQDDs",
			Computed:            true,
			ElementType:         types.StringType,
		},
		"software_only_disks": schema.ListAttribute{
			MarkdownDescription: "List of software-only disks.",
			Description:         "Software-Only Disks",
			Computed:            true,
			ElementType:         types.StringType,
		},
	}
}

// ScaleIODiskConfigurationSchema is a function that returns the schema for ScaleIODiskConfiguration
func ScaleIODiskConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"scale_io_storage_pool_disks": schema.ListNestedAttribute{
			MarkdownDescription: "List of ScaleIO storage pool disks configuration.",
			Description:         "ScaleIO Storage Pool Disks Configuration",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: ScaleIOStoragePoolDisksSchema()},
		},
	}
}

// ShortWindowSchema is a function that returns the schema for ShortWindow
func ShortWindowSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"threshold": schema.Int64Attribute{
			MarkdownDescription: "Threshold value for the short window.",
			Description:         "Short Window Threshold",
			Computed:            true,
		},
		"window_size_in_sec": schema.Int64Attribute{
			MarkdownDescription: "Size of the window in seconds for the short window.",
			Description:         "Short Window Size (in seconds)",
			Computed:            true,
		},
	}
}

// MediumWindowSchema is a function that returns the schema for MediumWindow
func MediumWindowSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"threshold": schema.Int64Attribute{
			MarkdownDescription: "Threshold value for the medium window.",
			Description:         "Medium Window Threshold",
			Computed:            true,
		},
		"window_size_in_sec": schema.Int64Attribute{
			MarkdownDescription: "Size of the window in seconds for the medium window.",
			Description:         "Medium Window Size (in seconds)",
			Computed:            true,
		},
	}
}

// LongWindowSchema is a function that returns the schema for LongWindow
func LongWindowSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"threshold": schema.Int64Attribute{
			MarkdownDescription: "Threshold value for the long window.",
			Description:         "Threshold for Long Window",
			Computed:            true,
		},
		"window_size_in_sec": schema.Int64Attribute{
			MarkdownDescription: "Window size in seconds for the long window.",
			Description:         "Long Window Size (in seconds)",
			Computed:            true,
		},
	}
}

// SdsDecoupledCounterParametersSchema is a function that returns the schema for SdsDecoupledCounterParameters
func SdsDecoupledCounterParametersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for short window counters.",
			Description:         "Configuration for short window counters.",
			Computed:            true,
			Attributes:          ShortWindowSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for medium window counters.",
			Description:         "Configuration for medium window counters.",
			Computed:            true,
			Attributes:          MediumWindowSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for long window counters.",
			Description:         "Configuration for long window counters.",
			Computed:            true,
			Attributes:          LongWindowSchema(),
		},
	}
}

// SdsConfigurationFailureCounterParametersSchema is a function that returns the schema for SdsConfigurationFailureCounterParameters
func SdsConfigurationFailureCounterParametersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for short window counters.",
			Description:         "Configuration for short window counters.",
			Computed:            true,
			Attributes:          ShortWindowSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for medium window counters.",
			Description:         "Configuration for medium window counters.",
			Computed:            true,
			Attributes:          MediumWindowSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for long window counters.",
			Description:         "Configuration for long window counters.",
			Computed:            true,
			Attributes:          LongWindowSchema(),
		},
	}
}

// MdmSdsCounterParametersSchema is a function that returns the schema for MdmSdsCounterParameters
func MdmSdsCounterParametersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for short window counters.",
			Description:         "Configuration for short window counters.",
			Computed:            true,
			Attributes:          ShortWindowSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for medium window counters.",
			Description:         "Configuration for medium window counters.",
			Computed:            true,
			Attributes:          MediumWindowSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for long window counters.",
			Description:         "Configuration for long window counters.",
			Computed:            true,
			Attributes:          LongWindowSchema(),
		},
	}
}

// SdsSdsCounterParametersSchema is a function that returns the schema for SdsSdsCounterParameters
func SdsSdsCounterParametersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for short window counters.",
			Description:         "Configuration for short window counters.",
			Computed:            true,
			Attributes:          ShortWindowSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for medium window counters.",
			Description:         "Configuration for medium window counters.",
			Computed:            true,
			Attributes:          MediumWindowSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for long window counters.",
			Description:         "Configuration for long window counters.",
			Computed:            true,
			Attributes:          LongWindowSchema(),
		},
	}
}

// SdsReceiveBufferAllocationFailuresCounterParametersSchema is a function that returns the schema for SdsReceiveBufferAllocationFailuresCounterParameters
func SdsReceiveBufferAllocationFailuresCounterParametersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"short_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for short window counters.",
			Description:         "Configuration for short window counters.",
			Computed:            true,
			Attributes:          ShortWindowSchema(),
		},
		"medium_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for medium window counters.",
			Description:         "Configuration for medium window counters.",
			Computed:            true,
			Attributes:          MediumWindowSchema(),
		},
		"long_window": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for long window counters.",
			Description:         "Configuration for long window counters.",
			Computed:            true,
			Attributes:          LongWindowSchema(),
		},
	}
}

// GeneralSchema is a function that returns the schema for General
func GeneralSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the system.",
			Description:         "The identifier of the system.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the system.",
			Description:         "The name of the system.",
			Computed:            true,
		},
		"system_id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the system.",
			Description:         "The unique identifier of the system.",
			Computed:            true,
		},
		"protection_domain_state": schema.StringAttribute{
			MarkdownDescription: "The state of the protection domain.",
			Description:         "The state of the protection domain.",
			Computed:            true,
		},
		"rebuild_network_throttling_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "The network throttling for rebuild operations in kilobits per second.",
			Description:         "The network throttling for rebuild operations in kilobits per second.",
			Computed:            true,
		},
		"rebalance_network_throttling_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "The network throttling for rebalance operations in kilobits per second.",
			Description:         "The network throttling for rebalance operations in kilobits per second.",
			Computed:            true,
		},
		"overall_io_network_throttling_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "The overall I/O network throttling in kilobits per second.",
			Description:         "The overall I/O network throttling in kilobits per second.",
			Computed:            true,
		},
		"sds_decoupled_counter_parameters": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for SDS decoupled counters.",
			Description:         "Configuration for SDS decoupled counters.",
			Computed:            true,
			Attributes:          SdsDecoupledCounterParametersSchema(),
		},
		"sds_configuration_failure_counter_parameters": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for SDS configuration failure counters.",
			Description:         "Configuration for SDS configuration failure counters.",
			Computed:            true,
			Attributes:          SdsConfigurationFailureCounterParametersSchema(),
		},
		"mdm_sds_counter_parameters": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for MDM SDS counters.",
			Description:         "Configuration for MDM SDS counters.",
			Computed:            true,
			Attributes:          MdmSdsCounterParametersSchema(),
		},
		"sds_sds_counter_parameters": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for SDS SDS counters.",
			Description:         "Configuration for SDS SDS counters.",
			Computed:            true,
			Attributes:          SdsSdsCounterParametersSchema(),
		},
		"rfcache_opertional_mode": schema.StringAttribute{
			MarkdownDescription: "The operational mode of Rfcache.",
			Description:         "The operational mode of Rfcache.",
			Computed:            true,
		},
		"rfcache_page_size_kb": schema.Int64Attribute{
			MarkdownDescription: "The page size of Rfcache in kilobytes.",
			Description:         "The page size of Rfcache in kilobytes.",
			Computed:            true,
		},
		"rfcache_max_io_size_kb": schema.Int64Attribute{
			MarkdownDescription: "The maximum I/O size of Rfcache in kilobytes.",
			Description:         "The maximum I/O size of Rfcache in kilobytes.",
			Computed:            true,
		},
		"sds_receive_buffer_allocation_failures_counter_parameters": schema.SingleNestedAttribute{
			MarkdownDescription: "Configuration for SDS receive buffer allocation failures counters.",
			Description:         "Configuration for SDS receive buffer allocation failures counters.",
			Computed:            true,
			Attributes:          SdsReceiveBufferAllocationFailuresCounterParametersSchema(),
		},
		"rebuild_network_throttling_enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether rebuild network throttling is enabled.",
			Description:         "Indicates whether rebuild network throttling is enabled.",
			Computed:            true,
		},
		"rebalance_network_throttling_enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether rebalance network throttling is enabled.",
			Description:         "Indicates whether rebalance network throttling is enabled.",
			Computed:            true,
		},
		"overall_io_network_throttling_enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether overall I/O network throttling is enabled.",
			Description:         "Indicates whether overall I/O network throttling is enabled.",
			Computed:            true,
		},
		"rfcache_enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether Rfcache is enabled.",
			Description:         "Indicates whether Rfcache is enabled.",
			Computed:            true,
		},
	}
}

// StatisticsDetailsSchema is a function that returns the schema for StatisticsDetails
func StatisticsDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"num_of_devices": schema.Int64Attribute{
			MarkdownDescription: "The number of devices in the system.",
			Description:         "The number of devices in the system.",
			Computed:            true,
		},
		"unused_capacity_in_kb": schema.Int64Attribute{
			MarkdownDescription: "The amount of unused capacity in kilobytes.",
			Description:         "The amount of unused capacity in kilobytes.",
			Computed:            true,
		},
		"num_of_volumes": schema.Int64Attribute{
			MarkdownDescription: "The number of volumes in the system.",
			Description:         "The number of volumes in the system.",
			Computed:            true,
		},
		"num_of_mapped_to_all_volumes": schema.Int64Attribute{
			MarkdownDescription: "The number of devices mapped to all volumes.",
			Description:         "The number of devices mapped to all volumes.",
			Computed:            true,
		},
		"capacity_available_for_volume_allocation_in_kb": schema.Int64Attribute{
			MarkdownDescription: "The capacity available for volume allocation in kilobytes.",
			Description:         "The capacity available for volume allocation in kilobytes.",
			Computed:            true,
		},
		"volume_allocation_limit_in_kb": schema.Int64Attribute{
			MarkdownDescription: "The volume allocation limit in kilobytes.",
			Description:         "The volume allocation limit in kilobytes.",
			Computed:            true,
		},
		"capacity_limit_in_kb": schema.Int64Attribute{
			MarkdownDescription: "The capacity limit in kilobytes.",
			Description:         "The capacity limit in kilobytes.",
			Computed:            true,
		},
		"num_of_unmapped_volumes": schema.Int64Attribute{
			MarkdownDescription: "The number of unmapped volumes in the system.",
			Description:         "The number of unmapped volumes in the system.",
			Computed:            true,
		},
		"spare_capacity_in_kb": schema.Int64Attribute{
			MarkdownDescription: "The spare capacity in kilobytes.",
			Description:         "The spare capacity in kilobytes.",
			Computed:            true,
		},
		"capacity_in_use_in_kb": schema.Int64Attribute{
			MarkdownDescription: "The capacity in use in kilobytes.",
			Description:         "The capacity in use in kilobytes.",
			Computed:            true,
		},
		"max_capacity_in_kb": schema.Int64Attribute{
			MarkdownDescription: "The maximum capacity in kilobytes.",
			Description:         "The maximum capacity in kilobytes.",
			Computed:            true,
		},
		"num_of_sds": schema.Int64Attribute{
			MarkdownDescription: "The number of SDS (Software-Defined Storage) instances.",
			Description:         "The number of SDS (Software-Defined Storage) instances.",
			Computed:            true,
		},
		"num_of_storage_pools": schema.Int64Attribute{
			MarkdownDescription: "The number of storage pools in the system.",
			Description:         "The number of storage pools in the system.",
			Computed:            true,
		},
		"num_of_fault_sets": schema.Int64Attribute{
			MarkdownDescription: "The number of fault sets in the system.",
			Description:         "The number of fault sets in the system.",
			Computed:            true,
		},
		"thin_capacity_in_use_in_kb": schema.Int64Attribute{
			MarkdownDescription: "The thin capacity in use in kilobytes.",
			Description:         "The thin capacity in use in kilobytes.",
			Computed:            true,
		},
		"thick_capacity_in_use_in_kb": schema.Int64Attribute{
			MarkdownDescription: "The thick capacity in use in kilobytes.",
			Description:         "The thick capacity in use in kilobytes.",
			Computed:            true,
		},
	}
}

// DiskListSchema is a function that returns the schema for DiskList
func DiskListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			Description:         "The unique identifier for the disk.",
			MarkdownDescription: "The unique identifier for the disk.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			Description:         "The name of the disk.",
			MarkdownDescription: "The name of the disk.",
			Computed:            true,
		},
		"error_state": schema.StringAttribute{
			Description:         "The error state of the disk.",
			MarkdownDescription: "The error state of the disk.",
			Computed:            true,
		},
		"sds_id": schema.StringAttribute{
			Description:         "The SDS (Software-Defined Storage) identifier for the disk.",
			MarkdownDescription: "The SDS (Software-Defined Storage) identifier for the disk.",
			Computed:            true,
		},
		"device_state": schema.StringAttribute{
			Description:         "The current state of the disk device.",
			MarkdownDescription: "The current state of the disk device.",
			Computed:            true,
		},
		"capacity_limit_in_kb": schema.Int64Attribute{
			Description:         "The capacity limit of the disk in kilobytes.",
			MarkdownDescription: "The capacity limit of the disk in kilobytes.",
			Computed:            true,
		},
		"max_capacity_in_kb": schema.Int64Attribute{
			Description:         "The maximum capacity of the disk in kilobytes.",
			MarkdownDescription: "The maximum capacity of the disk in kilobytes.",
			Computed:            true,
		},
		"storage_pool_id": schema.StringAttribute{
			Description:         "The identifier of the storage pool associated with the disk.",
			MarkdownDescription: "The identifier of the storage pool associated with the disk.",
			Computed:            true,
		},
		"device_current_path_name": schema.StringAttribute{
			Description:         "The current path name of the disk device.",
			MarkdownDescription: "The current path name of the disk device.",
			Computed:            true,
		},
		"device_original_path_name": schema.StringAttribute{
			Description:         "The original path name of the disk device.",
			MarkdownDescription: "The original path name of the disk device.",
			Computed:            true,
		},
		"serial_number": schema.StringAttribute{
			Description:         "The serial number of the disk.",
			MarkdownDescription: "The serial number of the disk.",
			Computed:            true,
		},
		"vendor_name": schema.StringAttribute{
			Description:         "The name of the disk vendor.",
			MarkdownDescription: "The name of the disk vendor.",
			Computed:            true,
		},
		"model_name": schema.StringAttribute{
			Description:         "The model name of the disk.",
			MarkdownDescription: "The model name of the disk.",
			Computed:            true,
		},
	}
}

// MappedSdcInfoDetailsSchema is a function that returns the schema for MappedSdcInfoDetails
func MappedSdcInfoDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"sdc_ip": schema.StringAttribute{
			MarkdownDescription: "The IP address of the Storage Data Controller (SDC).",
			Description:         "The IP address of the Storage Data Controller (SDC).",
			Computed:            true,
		},
		"sdc_id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier of the Storage Data Controller (SDC).",
			Description:         "The unique identifier of the Storage Data Controller (SDC).",
			Computed:            true,
		},
		"limit_bw_in_mbps": schema.Int64Attribute{
			MarkdownDescription: "The bandwidth limit in megabits per second (Mbps).",
			Description:         "The bandwidth limit in megabits per second (Mbps).",
			Computed:            true,
		},
		"limit_iops": schema.Int64Attribute{
			MarkdownDescription: "The IOPS (Input/Output Operations Per Second) limit.",
			Description:         "The IOPS (Input/Output Operations Per Second) limit.",
			Computed:            true,
		},
	}
}

// VolumeListSchema is a function that returns the schema for VolumeList
func VolumeListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the volume.",
			Description:         "The unique identifier for the volume.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the volume.",
			Description:         "The name of the volume.",
			Computed:            true,
		},
		"volume_type": schema.StringAttribute{
			MarkdownDescription: "The type of the volume.",
			Description:         "The type of the volume.",
			Computed:            true,
		},
		"storage_pool_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the storage pool associated with the volume.",
			Description:         "The identifier of the storage pool associated with the volume.",
			Computed:            true,
		},
		"data_layout": schema.StringAttribute{
			MarkdownDescription: "The data layout of the volume.",
			Description:         "The data layout of the volume.",
			Computed:            true,
		},
		"compression_method": schema.StringAttribute{
			MarkdownDescription: "The compression method used for the volume.",
			Description:         "The compression method used for the volume.",
			Computed:            true,
		},
		"size_in_kb": schema.Int64Attribute{
			MarkdownDescription: "The size of the volume in kilobytes.",
			Description:         "The size of the volume in kilobytes.",
			Computed:            true,
		},
		"mapped_sdc_info": schema.ListNestedAttribute{
			MarkdownDescription: "Information about the Storage Data Controllers (SDCs) mapped to the volume.",
			Description:         "Information about the Storage Data Controllers (SDCs) mapped to the volume.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: MappedSdcInfoDetailsSchema()},
		},
		"volume_class": schema.StringAttribute{
			MarkdownDescription: "The class or category of the volume.",
			Description:         "The class or category of the volume.",
			Computed:            true,
		},
	}
}

// StoragePoolListSchema is a function that returns the schema for StoragePoolList
func StoragePoolListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the storage pool.",
			Description:         "The unique identifier for the storage pool.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the storage pool.",
			Description:         "The name of the storage pool.",
			Computed:            true,
		},
		"rebuild_io_priority_policy": schema.StringAttribute{
			MarkdownDescription: "The rebuild I/O priority policy for the storage pool.",
			Description:         "The rebuild I/O priority policy for the storage pool.",
			Computed:            true,
		},
		"rebalance_io_priority_policy": schema.StringAttribute{
			MarkdownDescription: "The rebalance I/O priority policy for the storage pool.",
			Description:         "The rebalance I/O priority policy for the storage pool.",
			Computed:            true,
		},
		"rebuild_io_priority_num_of_concurrent_ios_per_device": schema.Int64Attribute{
			MarkdownDescription: "The number of concurrent I/O operations per device during rebuild.",
			Description:         "The number of concurrent I/O operations per device during rebuild.",
			Computed:            true,
		},
		"rebalance_io_priority_num_of_concurrent_ios_per_device": schema.Int64Attribute{
			MarkdownDescription: "The number of concurrent I/O operations per device during rebalance.",
			Description:         "The number of concurrent I/O operations per device during rebalance.",
			Computed:            true,
		},
		"rebuild_io_priority_bw_limit_per_device_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "The bandwidth limit per device in kilobits per second (kbps) during rebuild.",
			Description:         "The bandwidth limit per device in kilobits per second (kbps) during rebuild.",
			Computed:            true,
		},
		"rebalance_io_priority_bw_limit_per_device_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "The bandwidth limit per device in kilobits per second (kbps) during rebalance.",
			Description:         "The bandwidth limit per device in kilobits per second (kbps) during rebalance.",
			Computed:            true,
		},
		"rebuild_io_priority_app_iops_per_device_threshold": schema.StringAttribute{
			MarkdownDescription: "The application IOPS per device threshold during rebuild.",
			Description:         "The application IOPS per device threshold during rebuild.",
			Computed:            true,
		},
		"rebalance_io_priority_app_iops_per_device_threshold": schema.StringAttribute{
			MarkdownDescription: "The application IOPS per device threshold during rebalance.",
			Description:         "The application IOPS per device threshold during rebalance.",
			Computed:            true,
		},
		"rebuild_io_priority_app_bw_per_device_threshold_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "The application bandwidth per device threshold in kilobits per second (kbps) during rebuild.",
			Description:         "The application bandwidth per device threshold in kilobits per second (kbps) during rebuild.",
			Computed:            true,
		},
		"rebalance_io_priority_app_bw_per_device_threshold_in_kbps": schema.Int64Attribute{
			MarkdownDescription: "The application bandwidth per device threshold in kilobits per second (kbps) during rebalance.",
			Description:         "The application bandwidth per device threshold in kilobits per second (kbps) during rebalance.",
			Computed:            true,
		},
		"rebuild_io_priority_quiet_period_in_msec": schema.Int64Attribute{
			MarkdownDescription: "The quiet period in milliseconds during rebuild.",
			Description:         "The quiet period in milliseconds during rebuild.",
			Computed:            true,
		},
		"rebalance_io_priority_quiet_period_in_msec": schema.Int64Attribute{
			MarkdownDescription: "The quiet period in milliseconds during rebalance.",
			Description:         "The quiet period in milliseconds during rebalance.",
			Computed:            true,
		},
		"zero_padding_enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether zero padding is enabled.",
			Description:         "Indicates whether zero padding is enabled.",
			Computed:            true,
		},
		"background_scanner_mode": schema.StringAttribute{
			MarkdownDescription: "The mode of the background scanner.",
			Description:         "The mode of the background scanner.",
			Computed:            true,
		},
		"background_scanner_bw_limit_k_bps": schema.Int64Attribute{
			MarkdownDescription: "The bandwidth limit for the background scanner in kilobits per second (kbps).",
			Description:         "The bandwidth limit for the background scanner in kilobits per second (kbps).",
			Computed:            true,
		},
		"use_rmcache": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to use read-modify cache (RMCache).",
			Description:         "Indicates whether to use read-modify cache (RMCache).",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "The identifier for the protection domain.",
			Description:         "The identifier for the protection domain.",
			Computed:            true,
		},
		"sp_class": schema.StringAttribute{
			MarkdownDescription: "The class of the storage pool.",
			Description:         "The class of the storage pool.",
			Computed:            true,
		},
		"use_rfcache": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether to use read flash cache (RFCache).",
			Description:         "Indicates whether to use read flash cache (RFCache).",
			Computed:            true,
		},
		"spare_percentage": schema.Int64Attribute{
			MarkdownDescription: "The percentage of spare capacity in the storage pool.",
			Description:         "The percentage of spare capacity in the storage pool.",
			Computed:            true,
		},
		"rmcache_write_handling_mode": schema.StringAttribute{
			MarkdownDescription: "The write handling mode for the read-modify cache (RMCache).",
			Description:         "The write handling mode for the read-modify cache (RMCache).",
			Computed:            true,
		},
		"checksum_enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether checksum is enabled.",
			Description:         "Indicates whether checksum is enabled.",
			Computed:            true,
		},
		"rebuild_enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether rebuild is enabled.",
			Description:         "Indicates whether rebuild is enabled.",
			Computed:            true,
		},
		"rebalance_enabled": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether rebalance is enabled.",
			Description:         "Indicates whether rebalance is enabled.",
			Computed:            true,
		},
		"num_of_parallel_rebuild_rebalance_jobs_per_device": schema.Int64Attribute{
			MarkdownDescription: "The number of parallel rebuild or rebalance jobs per device.",
			Description:         "The number of parallel rebuild or rebalance jobs per device.",
			Computed:            true,
		},
		"capacity_alert_high_threshold": schema.Int64Attribute{
			MarkdownDescription: "The high threshold for capacity alerts.",
			Description:         "The high threshold for capacity alerts.",
			Computed:            true,
		},
		"capacity_alert_critical_threshold": schema.Int64Attribute{
			MarkdownDescription: "The critical threshold for capacity alerts.",
			Description:         "The critical threshold for capacity alerts.",
			Computed:            true,
		},
		"statistics": schema.SingleNestedAttribute{
			MarkdownDescription: "Statistics related to the storage pool.",
			Description:         "Statistics related to the storage pool.",
			Computed:            true,
			Attributes:          StatisticsDetailsSchema(),
		},
		"data_layout": schema.StringAttribute{
			MarkdownDescription: "The data layout used by the storage pool.",
			Description:         "The data layout used by the storage pool.",
			Computed:            true,
		},
		"replication_capacity_max_ratio": schema.StringAttribute{
			MarkdownDescription: "The maximum ratio of replication capacity.",
			Description:         "The maximum ratio of replication capacity.",
			Computed:            true,
		},
		"media_type": schema.StringAttribute{
			MarkdownDescription: "The type of media used by the storage pool.",
			Description:         "The type of media used by the storage pool.",
			Computed:            true,
		},
		"disk_list": schema.ListNestedAttribute{
			MarkdownDescription: "The list of disks associated with the storage pool.",
			Description:         "The list of disks associated with the storage pool.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DiskListSchema()},
		},
		"volume_list": schema.ListNestedAttribute{
			MarkdownDescription: "The list of volumes associated with the storage pool.",
			Description:         "The list of volumes associated with the storage pool.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: VolumeListSchema()},
		},
		"fgl_accp_id": schema.StringAttribute{
			MarkdownDescription: "The identifier for FGL (Fast Global Lock) acceleration.",
			Description:         "The identifier for FGL (Fast Global Lock) acceleration.",
			Computed:            true,
		},
	}
}

// IPListSchema is a function that returns the schema for IPList
func IPListSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"ip": schema.StringAttribute{
			MarkdownDescription: "The IP address.",
			Description:         "The IP address.",
			Computed:            true,
		},
		"role": schema.StringAttribute{
			MarkdownDescription: "The role associated with the IP address.",
			Description:         "The role associated with the IP address.",
			Computed:            true,
		},
	}
}

// SdsListDetailsSchema is a function that returns the schema for SdsListDetails
func SdsListDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the SDS (Software Defined Storage).",
			Description:         "The unique identifier for the SDS (Software Defined Storage).",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the SDS.",
			Description:         "The name of the SDS.",
			Computed:            true,
		},
		"port": schema.Int64Attribute{
			MarkdownDescription: "The port number used by the SDS.",
			Description:         "The port number used by the SDS.",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "The identifier for the protection domain associated with the SDS.",
			Description:         "The identifier for the protection domain associated with the SDS.",
			Computed:            true,
		},
		"fault_set_id": schema.StringAttribute{
			MarkdownDescription: "The identifier for the fault set associated with the SDS.",
			Description:         "The identifier for the fault set associated with the SDS.",
			Computed:            true,
		},
		"software_version_info": schema.StringAttribute{
			MarkdownDescription: "Information about the software version of the SDS.",
			Description:         "Information about the software version of the SDS.",
			Computed:            true,
		},
		"sds_state": schema.StringAttribute{
			MarkdownDescription: "The current state of the SDS.",
			Description:         "The current state of the SDS.",
			Computed:            true,
		},
		"membership_state": schema.StringAttribute{
			MarkdownDescription: "The membership state of the SDS.",
			Description:         "The membership state of the SDS.",
			Computed:            true,
		},
		"mdm_connection_state": schema.StringAttribute{
			MarkdownDescription: "The connection state of the SDS with the MDM (Metadata Manager).",
			Description:         "The connection state of the SDS with the MDM (Metadata Manager).",
			Computed:            true,
		},
		"drl_mode": schema.StringAttribute{
			MarkdownDescription: "The mode of Dynamic Resynchronization Lock (DRL) for the SDS.",
			Description:         "The mode of Dynamic Resynchronization Lock (DRL) for the SDS.",
			Computed:            true,
		},
		"maintenance_state": schema.StringAttribute{
			MarkdownDescription: "The maintenance state of the SDS.",
			Description:         "The maintenance state of the SDS.",
			Computed:            true,
		},
		"perf_profile": schema.StringAttribute{
			MarkdownDescription: "The performance profile of the SDS.",
			Description:         "The performance profile of the SDS.",
			Computed:            true,
		},
		"on_vm_ware": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the SDS is running on VMware.",
			Description:         "Indicates whether the SDS is running on VMware.",
			Computed:            true,
		},
		"ip_list": schema.ListNestedAttribute{
			MarkdownDescription: "The list of IP addresses associated with the SDS.",
			Description:         "The list of IP addresses associated with the SDS.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: IPListSchema()},
		},
	}
}

// SdrListDetailsSchema is a function that returns the schema for SdrListDetails
func SdrListDetailsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the SDR (Software Defined Router).",
			Description:         "The unique identifier for the SDR (Software Defined Router).",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the SDR.",
			Description:         "The name of the SDR.",
			Computed:            true,
		},
		"port": schema.Int64Attribute{
			MarkdownDescription: "The port number used by the SDR.",
			Description:         "The port number used by the SDR.",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "The identifier for the protection domain associated with the SDR.",
			Description:         "The identifier for the protection domain associated with the SDR.",
			Computed:            true,
		},
		"software_version_info": schema.StringAttribute{
			MarkdownDescription: "Information about the software version of the SDR.",
			Description:         "Information about the software version of the SDR.",
			Computed:            true,
		},
		"sdr_state": schema.StringAttribute{
			MarkdownDescription: "The current state of the SDR.",
			Description:         "The current state of the SDR.",
			Computed:            true,
		},
		"membership_state": schema.StringAttribute{
			MarkdownDescription: "The membership state of the SDR.",
			Description:         "The membership state of the SDR.",
			Computed:            true,
		},
		"mdm_connection_state": schema.StringAttribute{
			MarkdownDescription: "The connection state of the SDR with the MDM (Metadata Manager).",
			Description:         "The connection state of the SDR with the MDM (Metadata Manager).",
			Computed:            true,
		},
		"maintenance_state": schema.StringAttribute{
			MarkdownDescription: "The maintenance state of the SDR.",
			Description:         "The maintenance state of the SDR.",
			Computed:            true,
		},
		"perf_profile": schema.StringAttribute{
			MarkdownDescription: "The performance profile of the SDR.",
			Description:         "The performance profile of the SDR.",
			Computed:            true,
		},
		"ip_list": schema.ListNestedAttribute{
			MarkdownDescription: "The list of IP addresses associated with the SDR.",
			Description:         "The list of IP addresses associated with the SDR.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: IPListSchema()},
		},
	}
}

// AccelerationPoolSchema is a function that returns the schema for AccelerationPool
func AccelerationPoolSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the acceleration pool.",
			Description:         "The unique identifier for the acceleration pool.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the acceleration pool.",
			Description:         "The name of the acceleration pool.",
			Computed:            true,
		},
		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "The identifier for the protection domain associated with the acceleration pool.",
			Description:         "The identifier for the protection domain associated with the acceleration pool.",
			Computed:            true,
		},
		"media_type": schema.StringAttribute{
			MarkdownDescription: "The type of media used by the acceleration pool.",
			Description:         "The type of media used by the acceleration pool.",
			Computed:            true,
		},
		"rfcache": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the acceleration pool uses rfcache.",
			Description:         "Indicates whether the acceleration pool uses rfcache.",
			Computed:            true,
		},
	}
}

// ProtectionDomainSettingsSchema is a function that returns the schema for ProtectionDomainSettings
func ProtectionDomainSettingsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"general": schema.SingleNestedAttribute{
			MarkdownDescription: "General settings for the protection domain.",
			Description:         "General settings for the protection domain.",
			Computed:            true,
			Attributes:          GeneralSchema(),
		},
		"statistics": schema.SingleNestedAttribute{
			MarkdownDescription: "Statistics details for the protection domain.",
			Description:         "Statistics details for the protection domain.",
			Computed:            true,
			Attributes:          StatisticsDetailsSchema(),
		},
		"storage_pool_list": schema.ListNestedAttribute{
			MarkdownDescription: "List of storage pools associated with the protection domain.",
			Description:         "List of storage pools associated with the protection domain.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: StoragePoolListSchema()},
		},
		"sds_list": schema.ListNestedAttribute{
			MarkdownDescription: "List of SDS (Software Defined Storage) details associated with the protection domain.",
			Description:         "List of SDS (Software Defined Storage) details associated with the protection domain.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: SdsListDetailsSchema()},
		},
		"sdr_list": schema.ListNestedAttribute{
			MarkdownDescription: "List of SDR (Software Defined Router) details associated with the protection domain.",
			Description:         "List of SDR (Software Defined Router) details associated with the protection domain.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: SdrListDetailsSchema()},
		},
		"acceleration_pool": schema.ListNestedAttribute{
			MarkdownDescription: "List of acceleration pools associated with the protection domain.",
			Description:         "List of acceleration pools associated with the protection domain.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: AccelerationPoolSchema()},
		},
	}
}

// FaultSetSettingsSchema is a function that returns the schema for FaultSetSettings
func FaultSetSettingsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"protection_domain_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the protection domain associated with the fault set.",
			Description:         "The identifier of the protection domain associated with the fault set.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the fault set.",
			Description:         "The name of the fault set.",
			Computed:            true,
		},
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the fault set.",
			Description:         "The unique identifier for the fault set.",
			Computed:            true,
		},
	}
}

// DatacenterSchema is a function that returns the schema for Datacenter
func DatacenterSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"vcenter_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the vCenter associated with the datacenter.",
			Description:         "The identifier of the vCenter associated with the datacenter.",
			Computed:            true,
		},
		"datacenter_id": schema.StringAttribute{
			MarkdownDescription: "The identifier of the datacenter.",
			Description:         "The identifier of the datacenter.",
			Computed:            true,
		},
		"datacenter_name": schema.StringAttribute{
			MarkdownDescription: "The name of the datacenter.",
			Description:         "The name of the datacenter.",
			Computed:            true,
		},
	}
}

// PortGroupOptionsSchema is a function that returns the schema for PortGroupOptions
func PortGroupOptionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the port group option.",
			Description:         "The unique identifier for the port group option.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the port group option.",
			Description:         "The name of the port group option.",
			Computed:            true,
		},
	}
}

// PortGroupsSchema is a function that returns the schema for PortGroups
func PortGroupsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the port group.",
			Description:         "The unique identifier for the port group.",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "The display name of the port group.",
			Description:         "The display name of the port group.",
			Computed:            true,
		},
		"vlan": schema.Int64Attribute{
			MarkdownDescription: "The VLAN (Virtual Local Area Network) associated with the port group.",
			Description:         "The VLAN (Virtual Local Area Network) associated with the port group.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the port group.",
			Description:         "The name of the port group.",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "The value of the port group.",
			Description:         "The value of the port group.",
			Computed:            true,
		},
		"port_group_options": schema.ListNestedAttribute{
			MarkdownDescription: "List of options associated with the port group.",
			Description:         "List of options associated with the port group.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: PortGroupOptionsSchema()},
		},
	}
}

// VdsSettingsSchema is a function that returns the schema for VdsSettings
func VdsSettingsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the virtual distributed switch (vDS) setting.",
			Description:         "The unique identifier for the virtual distributed switch (vDS) setting.",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "The display name of the vDS setting.",
			Description:         "The display name of the vDS setting.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the vDS setting.",
			Description:         "The name of the vDS setting.",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "The value of the vDS setting.",
			Description:         "The value of the vDS setting.",
			Computed:            true,
		},
		"port_groups": schema.ListNestedAttribute{
			MarkdownDescription: "List of port groups associated with the vDS setting.",
			Description:         "List of port groups associated with the vDS setting.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: PortGroupsSchema()},
		},
	}
}

// VdsNetworkMtuSizeConfigurationSchema is a function that returns the schema for VdsNetworkMtuSizeConfiguration
func VdsNetworkMtuSizeConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the network MTU size configuration.",
			Description:         "The unique identifier for the network MTU size configuration.",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "The value of the network MTU size configuration.",
			Description:         "The value of the network MTU size configuration.",
			Computed:            true,
		},
	}
}

// VdsNetworkMTUSizeConfigurationSchema is a function that returns the schema for VdsNetworkMTUSizeConfiguration
func VdsNetworkMTUSizeConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the network MTU size configuration.",
			Description:         "The unique identifier for the network MTU size configuration.",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "The value of the network MTU size configuration.",
			Description:         "The value of the network MTU size configuration.",
			Computed:            true,
		},
	}
}

// VdsConfigurationSchema is a function that returns the schema for VdsConfiguration
func VdsConfigurationSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"datacenter": schema.SingleNestedAttribute{
			MarkdownDescription: "The datacenter associated with the virtual distributed switch (vDS) configuration.",
			Description:         "The datacenter associated with the virtual distributed switch (vDS) configuration.",
			Computed:            true,
			Attributes:          DatacenterSchema(),
		},
		"port_group_option": schema.StringAttribute{
			MarkdownDescription: "The option for port group associated with the vDS configuration.",
			Description:         "The option for port group associated with the vDS configuration.",
			Computed:            true,
		},
		"port_group_creation_option": schema.StringAttribute{
			MarkdownDescription: "The option for port group creation associated with the vDS configuration.",
			Description:         "The option for port group creation associated with the vDS configuration.",
			Computed:            true,
		},
		"vds_settings": schema.ListNestedAttribute{
			MarkdownDescription: "List of virtual distributed switch (vDS) settings.",
			Description:         "List of virtual distributed switch (vDS) settings.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: VdsSettingsSchema()},
		},
		"vds_network_mtu_size_configuration": schema.ListNestedAttribute{
			MarkdownDescription: "List of network MTU size configurations associated with the vDS configuration.",
			Description:         "List of network MTU size configurations associated with the vDS configuration.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: VdsNetworkMtuSizeConfigurationSchema()},
		},
	}
}

// NodeSelectionSchema is a function that returns the schema for NodeSelection
func NodeSelectionSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the node selection.",
			Description:         "The unique identifier for the node selection.",
			Computed:            true,
		},
		"service_tag": schema.StringAttribute{
			MarkdownDescription: "The service tag associated with the node selection.",
			Description:         "The service tag associated with the node selection.",
			Computed:            true,
		},
		"mgmt_ip_address": schema.StringAttribute{
			MarkdownDescription: "The management IP address of the node selection.",
			Description:         "The management IP address of the node selection.",
			Computed:            true,
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
		"related_components": schema.SingleNestedAttribute{
			MarkdownDescription: "Related components associated with this component.",
			Description:         "Related components associated with this component.",
			Computed:            true,
			Attributes:          RelatedComponentsSchema(),
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

// OptionsSchema is a function that returns the schema for Options
func OptionsSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the options.",
			Description:         "The unique identifier for the options.",
			Computed:            true,
		},
		"name": schema.StringAttribute{
			MarkdownDescription: "The name of the options.",
			Description:         "The name of the options.",
			Computed:            true,
		},
		"dependencies": schema.ListNestedAttribute{
			MarkdownDescription: "List of dependencies associated with the options.",
			Description:         "List of dependencies associated with the options.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DependenciesDetailsSchema()},
		},
		"attributes": schema.SingleNestedAttribute{
			MarkdownDescription: "Attributes associated with the options.",
			Description:         "Attributes associated with the options.",
			Computed:            true,
			Attributes:          AttributesSchema(),
		},
	}
}

// ParametersSchema is a function that returns the schema for Parameters
func ParametersSchema() map[string]schema.Attribute {
	return map[string]schema.Attribute{
		"id": schema.StringAttribute{
			MarkdownDescription: "The unique identifier for the parameter.",
			Description:         "The unique identifier for the parameter.",
			Computed:            true,
		},
		"value": schema.StringAttribute{
			MarkdownDescription: "The value associated with the parameter.",
			Description:         "The value associated with the parameter.",
			Computed:            true,
		},
		"display_name": schema.StringAttribute{
			MarkdownDescription: "The display name of the parameter.",
			Description:         "The display name of the parameter.",
			Computed:            true,
		},
		"type": schema.StringAttribute{
			MarkdownDescription: "The type of the parameter.",
			Description:         "The type of the parameter.",
			Computed:            true,
		},
		"tool_tip": schema.StringAttribute{
			MarkdownDescription: "Tooltip information for the parameter.",
			Description:         "Tooltip information for the parameter.",
			Computed:            true,
		},
		"required": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the parameter is required.",
			Description:         "Indicates whether the parameter is required.",
			Computed:            true,
		},
		"hide_from_template": schema.BoolAttribute{
			MarkdownDescription: "Specifies if the parameter should be hidden from the template.",
			Description:         "Specifies if the parameter should be hidden from the template.",
			Computed:            true,
		},
		"device_type": schema.StringAttribute{
			MarkdownDescription: "The type of device associated with the parameter.",
			Description:         "The type of device associated with the parameter.",
			Computed:            true,
		},
		"dependencies": schema.ListNestedAttribute{
			MarkdownDescription: "List of dependencies associated with the parameter.",
			Description:         "List of dependencies associated with the parameter.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: DependenciesDetailsSchema()},
		},
		"group": schema.StringAttribute{
			MarkdownDescription: "The group to which the parameter belongs.",
			Description:         "The group to which the parameter belongs.",
			Computed:            true,
		},
		"read_only": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether the parameter is read-only.",
			Description:         "Indicates whether the parameter is read-only.",
			Computed:            true,
		},
		"generated": schema.BoolAttribute{
			MarkdownDescription: "Specifies if the parameter is generated.",
			Description:         "Specifies if the parameter is generated.",
			Computed:            true,
		},
		"info_icon": schema.BoolAttribute{
			MarkdownDescription: "Indicates whether an information icon is associated with the parameter.",
			Description:         "Indicates whether an information icon is associated with the parameter.",
			Computed:            true,
		},
		"step": schema.Int64Attribute{
			MarkdownDescription: "The step value for the parameter.",
			Description:         "The step value for the parameter.",
			Computed:            true,
		},
		"max_length": schema.Int64Attribute{
			MarkdownDescription: "The maximum length allowed for the parameter.",
			Description:         "The maximum length allowed for the parameter.",
			Computed:            true,
		},
		"min": schema.Int64Attribute{
			MarkdownDescription: "The minimum value allowed for the parameter.",
			Description:         "The minimum value allowed for the parameter.",
			Computed:            true,
		},
		"max": schema.Int64Attribute{
			MarkdownDescription: "The maximum value allowed for the parameter.",
			Description:         "The maximum value allowed for the parameter.",
			Computed:            true,
		},
		"networks": schema.ListNestedAttribute{
			MarkdownDescription: "List of networks associated with the parameter.",
			Description:         "List of networks associated with the parameter.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: NetworksSchema()},
		},
		"options": schema.ListNestedAttribute{
			MarkdownDescription: "List of options associated with the parameter.",
			Description:         "List of options associated with the parameter.",
			Computed:            true,
			NestedObject:        schema.NestedAttributeObject{Attributes: OptionsSchema()},
		},
		"options_sortable": schema.BoolAttribute{
			MarkdownDescription: "Specifies if options are sortable for the parameter.",
			Description:         "Specifies if options are sortable for the parameter.",
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
