package sdc

import (
	"context"

	"terraform-provider-powerflex/client"
	sdchelper "terraform-provider-powerflex/helper/sdc"

	"github.com/dell/goscaleio"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func getAllSdcs(ctx context.Context, system *goscaleio.System, sdcID string, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics
	sdcs, err := system.GetSdc()
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Goscaleio all SDC[system.GetSdc]",
			Detail:   err.Error(),
		})
		tflog.Error(ctx, "[PowerFlex][getAllSdcs] Unable to get Goscaleio all SDC[system.GetSdc]")
		tflog.Error(ctx, err.Error())
		return diags
	}
	tflog.Debug(ctx, "[PowerFlex][getAllSdcs] sdc - "+client.PrettyJSON(sdcs))

	resultMap := make([]map[string]interface{}, 0)
	for _, v := range sdcs {
		resultSdcMap := sdchelper.SdcToMap(v)
		resultMap = append(resultMap, resultSdcMap)
	}

	if err := d.Set("sdcs", resultMap); err != nil {
		return diag.FromErr(err)
	}
	return diags

}

func getSingleSdc(ctx context.Context, system *goscaleio.System, sdcID string, d *schema.ResourceData) diag.Diagnostics {
	var diags diag.Diagnostics

	tflog.Debug(ctx, "[PowerFlex][getSingleSdc] sdcID - "+sdcID)
	singleSdc, err := system.FindSdc("ID", sdcID)
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Goscaleio single SDC[system.FindSdc]",
			Detail:   err.Error(),
		})
		tflog.Error(ctx, "[PowerFlex][getSingleSdc] Unable to get Goscaleio single SDC[system.FindSdc]")
		tflog.Error(ctx, err.Error())
		return diags
	}
	tflog.Debug(ctx, "[PowerFlex][getSingleSdc] sdc - "+client.PrettyJSON(singleSdc))

	resultMap := make([]map[string]interface{}, 0)
	resultSdcMap := sdchelper.SdcToMap(*singleSdc.Sdc)

	tflog.Debug(ctx, "[PowerFlex][Anshuman] sdc - "+client.PrettyJSON(resultSdcMap))
	resultMap = append(resultMap, resultSdcMap)

	if err := d.Set("sdcs", resultMap); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
