package powerflex

import (
	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
)

// VolumeMarkdownDescription add notes for resource
const VolumeMarkdownDescription = `This resource can be used to manage volumes on a PowerFlex array.
Note: Volume creation or update is not atomic. In case of partially completed operations, terraform can mark the resource as tainted.
One can manually remove the taint and try applying the configuration (after making necessary adjustments).
Warning: If the taint is not removed, terraform will destroy and recreate the resource.
`

// VolumeResourceSchema variable to define schema for the volume resource
var VolumeResourceSchema schema.Schema = schema.Schema{
	Description:         "This resource can be used to manage volumes on a PowerFlex array.",
	MarkdownDescription: VolumeMarkdownDescription,
	Attributes: map[string]schema.Attribute{
		"name": schema.StringAttribute{
			Description:         "The name of the volume.",
			Required:            true,
			MarkdownDescription: "The name of the volume.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
			},
		},
		"storage_pool_id": schema.StringAttribute{
			Description: "ID of the Storage Pool under which the volume will be created." +
				" Must be provided if and only if 'storage_pool_name' is not provided." +
				" Cannot be updated.",
			Optional: true,
			Computed: true,
			MarkdownDescription: "ID of the Storage Pool under which the volume will be created." +
				" Must be provided if and only if `storage_pool_name` is not provided." +
				" Cannot be updated.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRoot("storage_pool_name")),
			},
		},
		"storage_pool_name": schema.StringAttribute{
			Description: "Name of the Storage Pool under which the volume will be created." +
				" Must be provided if and only if 'storage_pool_id' is not provided." +
				" Cannot be updated.",
			Optional: true,
			Computed: true,
			MarkdownDescription: "Name of the Storage Pool under which the volume will be created." +
				" Must be provided if and only if `storage_pool_id` is not provided." +
				" Cannot be updated.",
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRoot("storage_pool_id")),
			},
		},
		"protection_domain_id": schema.StringAttribute{
			Description: "ID of the Protection Domain under which the volume will be created." +
				" Must be provided if and only if 'protection_domain_name' is not provided." +
				" Cannot be updated.",
			MarkdownDescription: "ID of the Protection Domain under which the volume will be created." +
				" Must be provided if and only if `protection_domain_name` is not provided." +
				" Cannot be updated.",
			Computed: true,
			Optional: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRoot("protection_domain_name")),
			},
		},
		"protection_domain_name": schema.StringAttribute{
			Description: "Name of the Protection Domain under which the volume will be created." +
				" Must be provided if and only if 'protection_domain_id' is not provided." +
				" Cannot be updated.",
			MarkdownDescription: "Name of the Protection Domain under which the volume will be created." +
				" Must be provided if and only if `protection_domain_id` is not provided." +
				" Cannot be updated.",
			Optional: true,
			Computed: true,
			Validators: []validator.String{
				stringvalidator.LengthAtLeast(1),
				stringvalidator.ExactlyOneOf(path.MatchRoot("protection_domain_id")),
			},
		},
		"size": schema.Int64Attribute{
			Description:         "Size of the volume. The unit of size is defined by 'capacity_unit'. The storage capacity of a volume cannot be decreased.",
			Required:            true,
			MarkdownDescription: "Size of the volume. The unit of size is defined by `capacity_unit`. The storage capacity of a volume cannot be decreased.",
		},
		"capacity_unit": schema.StringAttribute{
			Description:         "Unit of capacity of the volume. Must be one of 'GB' and 'TB'. Default value is 'GB'.",
			Computed:            true,
			Optional:            true,
			MarkdownDescription: "Unit of capacity of the volume. Must be one of `GB` and `TB`. Default value is `GB`.",
			Validators: []validator.String{stringvalidator.OneOf(
				"GB",
				"TB",
			)},
			PlanModifiers: []planmodifier.String{
				stringDefault("GB"),
			},
		},
		"volume_type": schema.StringAttribute{
			Description:         "Volume type. Valid values are 'ThickProvisioned' and 'ThinProvisioned'. Default value is 'ThinProvisioned'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Volume type. Valid values are `ThickProvisioned` and `ThinProvisioned`. Default value is `ThinProvisioned`.",
			Validators: []validator.String{stringvalidator.OneOf(
				"ThickProvisioned",
				"ThinProvisioned",
			)},
			PlanModifiers: []planmodifier.String{
				stringDefault("ThinProvisioned"),
			},
		},
		"use_rm_cache": schema.BoolAttribute{
			Description:         "use rm cache",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "use rm cache",
		},
		"compression_method": schema.StringAttribute{
			Description:         "Compression Method of the volume. Valid values are 'None' and 'Normal'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Compression Method of the volume. Valid values are `None` and `Normal`.",
			Validators: []validator.String{stringvalidator.OneOf(
				"None",
				"Normal",
			)},
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"id": schema.StringAttribute{
			Description:         "The ID of the volume.",
			Computed:            true,
			MarkdownDescription: "The ID of the volume.",
			PlanModifiers: []planmodifier.String{
				stringplanmodifier.UseStateForUnknown(),
			},
		},
		"size_in_kb": schema.Int64Attribute{
			Description:         "Size in KB",
			Computed:            true,
			MarkdownDescription: "Size in KB",
		},
		"access_mode": schema.StringAttribute{
			Description:         "The Access mode of the volume. Valid values are 'ReadOnly' and 'ReadWrite'. Default value is 'ReadOnly'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "The Access mode of the volume. Valid values are `ReadOnly` and `ReadWrite`. Default value is `ReadOnly`.",
			Validators: []validator.String{stringvalidator.OneOf(
				"ReadOnly",
				"ReadWrite",
			)},
			PlanModifiers: []planmodifier.String{
				stringDefault("ReadOnly"),
			},
		},
		"remove_mode": schema.StringAttribute{
			Description:         "Remove mode of the volume. Valid values are 'ONLY_ME' and 'INCLUDING_DESCENDANTS'. Default value is 'ONLY_ME'.",
			Optional:            true,
			Computed:            true,
			MarkdownDescription: "Remove mode of the volume. Valid values are `ONLY_ME` and `INCLUDING_DESCENDANTS`. Default value is `ONLY_ME`.",
			Validators: []validator.String{stringvalidator.OneOf(
				"ONLY_ME",
				"INCLUDING_DESCENDANTS",
			)},
			PlanModifiers: []planmodifier.String{
				stringDefault("ONLY_ME"),
			},
		},
		"sdc_list": sdcListSchema,
	},
}
