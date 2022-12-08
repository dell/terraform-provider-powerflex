package sdc

import (
	"context"
	"strconv"

	"terraform-provider-powerflex/client"
	sdchelper "terraform-provider-powerflex/helper/sdc"
	sdcmodels "terraform-provider-powerflex/models/sdc"
	"time"

	"github.com/AnshumanPradipPatil1506/goscaleio"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

// DataSourceSdcs function to carry sdc resource read schema
func DataSourceSdcs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSdcsRead,
		Schema:      sdcmodels.SDCReadOnlyModel,
	}
}

func dataSourceSdcsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	systemID := d.Get("systemid").(string)
	sdcID := d.Get("id").(string)

	goscaleioClient := client.GoscaleioClient

	if systemID == "" {
		tflog.Debug(ctx, "[PowerFlex][dataSourceSdcsRead] Empty value passed for systemid so getting first system available.")
		allSystems, err := goscaleioClient.GetSystems()
		if err != nil {
			diags = append(diags, diag.Diagnostic{
				Severity: diag.Error,
				Summary:  "Unable to get Goscaleio systems[goscaleioClient.GetSystems]",
				Detail:   err.Error(),
			})
			tflog.Error(ctx, "[PowerFlex][dataSourceSdcsRead] Unable to get Goscaleio systems[goscaleioClient.C.GetSystems]")
			tflog.Error(ctx, err.Error())
			return diags
		}
		systemID = allSystems[0].ID
	}

	system, err := goscaleioClient.FindSystem(systemID, "", "")
	if err != nil {
		diags = append(diags, diag.Diagnostic{
			Severity: diag.Error,
			Summary:  "Unable to get Goscaleio single System[goscaleioClient.FindSystem]",
			Detail:   err.Error(),
		})
		tflog.Error(ctx, "[PowerFlex][dataSourceSdcsRead] Unable to get Goscaleio single System[goscaleioClient.FindSystem]")
		tflog.Error(ctx, err.Error())
		return diags
	}

	if sdcID == "" {
		diags = getAllSdcs(ctx, system, sdcID, d)
	} else {
		diags = getSingleSdc(ctx, system, sdcID, d)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

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
