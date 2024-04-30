package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type SdcHostModel struct {
	ID     types.String `tfsdk:"id"`
	Remote types.Object `tfsdk:"remote"` // SdcHostRemoteModel
	Host   types.String `tfsdk:"ip"`
	Pkg    types.String `tfsdk:"package_path"`
	OS     types.String `tfsdk:"os_family"`
	DrvCfg types.String `tfsdk:"drv_cfg_path"`

	// optional
	Name               types.String `tfsdk:"name"`
	PerformanceProfile types.String `tfsdk:"performance_profile"`
	MdmIPs             types.List   `tfsdk:"mdm_ips"` // list(string)

	// optional, os specific
	Esxi types.Object `tfsdk:"esxi"`
}

type SdcHostRemoteModel struct {
	User       string  `tfsdk:"user"`
	Password   *string `tfsdk:"password"`
	PrivateKey *string `tfsdk:"private_key"`
	Dir        *string `tfsdk:"dir"`
}

type SdcHostEsxiModel struct {
	Guid types.String `tfsdk:"guid"`
}
