package datasources

import (
	"context"
	"strconv"

	"terraform-provider-powerflex/client/helper"
	schemastructures "terraform-provider-powerflex/powerflex/schema-structures"
	"time"

	"github.com/AnshumanPradipPatil1506/goscaleio"
	"github.com/hashicorp/terraform-plugin-log/tflog"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

func DataSourceSdcs() *schema.Resource {
	return &schema.Resource{
		ReadContext: dataSourceSdcsRead,
		Schema:      schemastructures.SDC_DATA_RESOURCE_SCHEMA,
	}
}

func dataSourceSdcsRead(ctx context.Context, d *schema.ResourceData, m interface{}) diag.Diagnostics {
	var diags diag.Diagnostics

	systemID := d.Get("systemid").(string)
	sdcID := d.Get("id").(string)

	goscaleioClient := helper.GoscaleioClient

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
		diags = getAllSdcs(system, sdcID, ctx, d)
	} else {
		diags = getSingleSdc(system, sdcID, ctx, d)
	}

	d.SetId(strconv.FormatInt(time.Now().Unix(), 10))

	return diags
}

func getAllSdcs(system *goscaleio.System, sdcID string, ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
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
	tflog.Debug(ctx, "[PowerFlex][getAllSdcs] sdc - "+helper.PrettyJson(sdcs))

	resultMap := make([]map[string]interface{}, 0)
	for _, v := range sdcs {
		resultSdcMap := schemastructures.SdcToMap(v)
		resultMap = append(resultMap, resultSdcMap)
	}

	if err := d.Set("sdcs", resultMap); err != nil {
		return diag.FromErr(err)
	}
	return diags

}

func getSingleSdc(system *goscaleio.System, sdcID string, ctx context.Context, d *schema.ResourceData) diag.Diagnostics {
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
	tflog.Debug(ctx, "[PowerFlex][getSingleSdc] sdc - "+helper.PrettyJson(singleSdc))

	resultMap := make([]map[string]interface{}, 0)
	resultSdcMap := schemastructures.SdcToMap(*singleSdc.Sdc)

	tflog.Debug(ctx, "[PowerFlex][Anshuman] sdc - "+helper.PrettyJson(resultSdcMap))
	resultMap = append(resultMap, resultSdcMap)

	if err := d.Set("sdcs", resultMap); err != nil {
		return diag.FromErr(err)
	}
	return diags
}
