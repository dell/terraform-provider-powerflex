package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type SdcHostModel struct {
	ID     types.String `tfsdk:"id"`
	Remote types.Object `tfsdk:"remote"` // SdcHostRemoteModel
	Host   types.String `tfsdk:"ip"`
	Pkg    types.String `tfsdk:"package_base64"`
	OS     types.String `tfsdk:"os_family"`
	DrvCfg types.String `tfsdk:"drv_cfg_base64"`

	// optional
	Name               types.String `tfsdk:"name"`
	PerformanceProfile types.String `tfsdk:"performance_profile"`
	MdmIPs             types.List   `tfsdk:"mdm_ips"` // list(string)
}

type SdcHostRemoteModel struct {
	User       string  `tfsdk:"user"`
	Password   *string `tfsdk:"password"`
	PrivateKey *string `tfsdk:"private_key"`
}
