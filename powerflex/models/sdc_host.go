package models

import "github.com/hashicorp/terraform-plugin-framework/types"

type SdcHostModel struct {
	ID     types.String `tfsdk:"id"`
	Remote types.Object `tfsdk:"remote"` // SdcHostRemoteModel
	Host   types.String `tfsdk:"hostname"`
	Pkg    types.String `tfsdk:"package_base64"`
	// OS     types.String `tfsdk:"os_family"`

	// optional
	Name   types.String `tfsdk:"name"`
	MdmIPs types.List   `tfsdk:"mdm_ips"` // list(string)
}

type SdcHostRemoteModel struct {
	User       string  `tfsdk:"user"`
	Password   *string `tfsdk:"password"`
	PrivateKey *string `tfsdk:"private_key"`
}
