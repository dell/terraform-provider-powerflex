package schemastructures

import (
	goscaleiotypes "github.com/AnshumanPradipPatil1506/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var SDCDataresourceSchema map[string]*schema.Schema = map[string]*schema.Schema{
	"id": {
		Type:        schema.TypeString,
		Description: "Enter ID of Powerflex SDC. [Default/empty will all sdc present in given system]",
		Required:    true,
		Sensitive:   true,
	},
	"systemid": {
		Type:        schema.TypeString,
		Description: "Enter System ID of Powerflex System. [Default/empty will be any first system in list]",
		Required:    true,
		Sensitive:   true,
	},
	"sdcs": {
		Type:     schema.TypeList,
		Computed: true,
		Elem: &schema.Resource{
			Schema: map[string]*schema.Schema{
				"id": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"name": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"sdcguid": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"sdcapproved": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"onvmware": {
					Type:     schema.TypeBool,
					Computed: true,
				},
				"systemid": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"sdcip": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"mdmconnectionstate": {
					Type:     schema.TypeString,
					Computed: true,
				},
				"links": {
					Type:     schema.TypeList,
					Computed: true,
					Elem: &schema.Resource{
						Schema: map[string]*schema.Schema{
							"rel": {
								Type:     schema.TypeString,
								Computed: true,
							},
							"href": {
								Type:     schema.TypeString,
								Computed: true,
							},
						},
					},
				},
			},
		},
	},
}

// Convert returned sdc object(from goscaleio) to terraform output schema
func SdcToMap(s goscaleiotypes.Sdc) map[string]interface{} {
	resultSDC := make(map[string]interface{})
	resultSDC["id"] = s.ID
	resultSDC["name"] = s.Name
	resultSDC["sdcguid"] = s.SdcGUID
	resultSDC["sdcapproved"] = s.SdcApproved
	resultSDC["onvmware"] = s.OnVMWare
	resultSDC["systemid"] = s.SystemID
	resultSDC["sdcip"] = s.SdcIP
	resultSDC["mdmconnectionstate"] = s.MdmConnectionState
	return resultSDC
}

// Convert returned sdc object(from goscaleio) [after name changes] to terraform output schema
func NameChangedSdcToMap(s goscaleiotypes.Sdc) map[string]interface{} {
	resultSDC := make(map[string]interface{})
	resultSDC["id"] = s.ID
	resultSDC["name"] = s.Name
	return resultSDC
}

// func StructToMap(val interface{}, ctx context.Context) map[string]interface{} {

// 	var data map[string]interface{} = make(map[string]interface{})
// 	varType := reflect.TypeOf(val)
// 	if varType.Kind() != reflect.Struct {
// 		fmt.Println("Not a struct")
// 		return nil
// 	}

// 	value := reflect.ValueOf(val)
// 	for i := 0; i < varType.NumField(); i++ {
// 		if !value.Field(i).CanInterface() {
// 			continue
// 		}
// 		fieldName := strings.ToLower(varType.Field(i).Name)

// 		if varType.Field(i).Type.Kind() != reflect.Struct {
// 			tflog.Debug(ctx, "[PowerFlex][Anshuman] --- ELSE BLOCK fieldName "+fieldName+" = "+varType.Field(i).Type.Kind().String())
// 			if varType.Field(i).Type.Kind() == reflect.Slice {
// 				tflog.Debug(ctx, "[PowerFlex][Anshuman] --- reflect.Slice ")
// 			}
// 			if varType.Field(i).Type.Kind() == reflect.Array {
// 				tflog.Debug(ctx, "[PowerFlex][Anshuman] --- reflect.Array ")
// 			}
// 			data[fieldName] = value.Field(i).Interface()
// 		} else {

// 			data[fieldName] = StructToMap(value.Field(i).Interface(), ctx)
// 		}

// 	}

// 	return data
// }
