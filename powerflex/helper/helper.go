/*
Copyright (c) 2022-2024 Dell Inc., or its subsidiaries. All Rights Reserved.

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

package helper

import (
	"context"
	"fmt"
	"math/big"
	"reflect"
	"strconv"
	"strings"

	"bytes"
	"encoding/json"

	"github.com/dell/goscaleio"
	scaleiotypes "github.com/dell/goscaleio/types/v1"
	"github.com/hashicorp/terraform-plugin-framework-validators/setvalidator"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-log/tflog"
)

// GetFirstSystem - finds available first system and returns it.
func GetFirstSystem(rc *goscaleio.Client) (*goscaleio.System, error) {
	allSystems, err := rc.GetSystems()
	if err != nil {
		return nil, fmt.Errorf("Error in goscaleio GetSystems %s", err.Error())
	}
	if numSys := len((allSystems)); numSys == 0 {
		return nil, fmt.Errorf("no systems found")
	} else if numSys > 1 {
		return nil, fmt.Errorf("more than one system found")
	}
	system, err := rc.FindSystem(allSystems[0].ID, "", "")
	if err != nil {
		return nil, fmt.Errorf("Error in goscaleio FindSystem")
	}
	return system, nil
}

// PrettyJSON - function for logging json readable output.
func PrettyJSON(data interface{}) string {
	buffer := new(bytes.Buffer)
	encoder := json.NewEncoder(buffer)
	encoder.SetIndent("", "\t")

	err := encoder.Encode(data)
	if err != nil {
		return ""
	}
	return buffer.String()
}

// GetNewProtectionDomainEx function to get Protection Domain
func GetNewProtectionDomainEx(c *goscaleio.Client, pdID string, pdName string, href string) (*goscaleio.ProtectionDomain, error) {
	system, err := GetFirstSystem(c)
	if err != nil {
		return nil, err
	}
	pdr := goscaleio.NewProtectionDomainEx(c, &scaleiotypes.ProtectionDomain{})
	if pdID != "" {
		protectionDomain, err := system.FindProtectionDomain(pdID, "", "")
		pdr.ProtectionDomain = protectionDomain
		if err != nil {
			return nil, err
		}
	} else {
		protectionDomain, err := system.FindProtectionDomain("", pdName, "")
		pdr.ProtectionDomain = protectionDomain
		if err != nil {
			return nil, err
		}
	}
	return pdr, nil
}

// ModifyStoragePoolName function to modify a storagepool name
func ModifyStoragePoolName(pd *goscaleio.ProtectionDomain, id string, newName string) (string, error) {
	return pd.ModifyStoragePoolName(id, newName)
}

// CreateFaultSet create a new fault set
func CreateFaultSet(pd *goscaleio.ProtectionDomain, param *scaleiotypes.FaultSetParam) (string, error) {
	return pd.CreateFaultSet(param)
}

// ModifyFaultSetName function to modify a fault set name
func ModifyFaultSetName(pd *goscaleio.ProtectionDomain, id string, newName string) error {
	return pd.ModifyFaultSetName(id, newName)
}

// GetStoragePoolType returns storage pool type
func GetStoragePoolType(r *goscaleio.Client, storagePoolID string) (*goscaleio.StoragePool, error) {
	system, err := GetFirstSystem(r)
	if err != nil {
		return nil, err
	}

	sp, err := system.GetStoragePoolByID(storagePoolID)
	if err != nil {
		return nil, err
	}

	sp1 := goscaleio.NewStoragePoolEx(r, sp)
	return sp1, nil
}

// GetSdcType function returns SDC type
func GetSdcType(c *goscaleio.Client, sdcID string) (*goscaleio.Sdc, error) {
	system, err := GetFirstSystem(c)
	if err != nil {
		return nil, err
	}
	return system.GetSdcByID(sdcID)
}

// GetVolumeType function returns volume type
func GetVolumeType(c *goscaleio.Client, volID string) (*goscaleio.Volume, error) {
	volumes, err := c.GetVolume("", volID, "", "", false)
	if err != nil {
		return nil, err
	}

	volume := volumes[0]
	volType := goscaleio.NewVolume(c)
	volType.Volume = volume
	return volType, nil
}

// StringDefaultModifier is a plan modifier that sets a default value for a
// types.StringType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type StringDefaultModifier struct {
	Default string
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m StringDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %s", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m StringDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%s`", m.Default)
}

// PlanModifyString runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m StringDefaultModifier) PlanModifyString(ctx context.Context, req planmodifier.StringRequest, resp *planmodifier.StringResponse) {
	// If the value is unknown or known, do not set default value.
	if req.PlanValue.IsNull() {
		resp.PlanValue = types.StringValue(m.Default)
	}
	if req.PlanValue.IsUnknown() {
		resp.PlanValue = types.StringValue(m.Default)
	}
}

// StringDefault sets default value fot string attributes
func StringDefault(defaultValue string) planmodifier.String {
	return StringDefaultModifier{
		Default: defaultValue,
	}
}

// boolDefaultModifier is a plan modifier that sets a default value for a
// types.BoolType attribute when it is not configured. The attribute must be
// marked as Optional and Computed. When setting the state during the resource
// Create, Read, or Update methods, this default value must also be included or
// the Terraform CLI will generate an error.
type boolDefaultModifier struct {
	Default bool
}

// Description returns a plain text description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m boolDefaultModifier) Description(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to %t", m.Default)
}

// MarkdownDescription returns a markdown formatted description of the validator's behavior, suitable for a practitioner to understand its impact.
func (m boolDefaultModifier) MarkdownDescription(ctx context.Context) string {
	return fmt.Sprintf("If value is not configured, defaults to `%t`", m.Default)
}

// PlanModifyBool runs the logic of the plan modifier.
// Access to the configuration, plan, and state is available in `req`, while
// `resp` contains fields for updating the planned value, triggering resource
// replacement, and returning diagnostics.
func (m boolDefaultModifier) PlanModifyBool(ctx context.Context, req planmodifier.BoolRequest, resp *planmodifier.BoolResponse) {
	// If the value is unknown or known, do not set default value.
	if req.PlanValue.IsNull() || req.PlanValue.IsUnknown() {
		resp.PlanValue = types.BoolValue(m.Default)
	}
}

// BoolDefault sets default value fot string attributes
func BoolDefault(defaultValue bool) planmodifier.Bool {
	return boolDefaultModifier{
		Default: defaultValue,
	}
}

// RenewInstallationCookie is used to renew the installation cookie, i.e. LEGACYGWCOOKIE.
// Using the same LEGACYGWCOOKIE ensures that the REST requests are sent to the same GW pod.
// That would help to get the correct response from the GW pod that stores installation packages.
func RenewInstallationCookie(gatewayClient *goscaleio.GatewayClient) error {
	return gatewayClient.RenewInstallationCookie(10)
}

// ResetInstallerQueue function for the Abort, Clear and Move To Idle Execution
func ResetInstallerQueue(gatewayClient *goscaleio.GatewayClient) error {

	_, err := gatewayClient.AbortOperation()

	if err != nil {
		return fmt.Errorf("Error while Aborting Operation is %s", err.Error())
	}
	_, err = gatewayClient.ClearQueueCommand()

	if err != nil {
		return fmt.Errorf("Error while Clearing Queue is %s", err.Error())
	}

	_, err = gatewayClient.MoveToIdlePhase()

	if err != nil {
		return fmt.Errorf("Error while Move to Ideal Phase is %s", err.Error())
	}

	return nil
}

// CompareStringSlice Compare string slices. return true if the length and elements are same.
func CompareStringSlice(plan, state []string) bool {
	if len(plan) != len(state) {
		return false
	}

	itemAppearsTimes := make(map[string]int, len(plan))
	for _, i := range plan {
		itemAppearsTimes[i]++
	}

	for _, i := range state {
		if _, ok := itemAppearsTimes[i]; !ok {
			return false
		}

		itemAppearsTimes[i]--
		if itemAppearsTimes[i] == 0 {
			delete(itemAppearsTimes, i)
		}
	}
	return len(itemAppearsTimes) == 0
}

// CompareInt64Slice compares two slices of int64 and returns true if the length and elements are same.
func CompareInt64Slice(plan, state []int64) bool {

	if len(plan) != len(state) {
		return false
	}

	for i := range plan {
		if plan[i] != state[i] {
			return false
		}
	}
	return true
}

// SetContains is a helper function to check if a Set contains a value
func SetContains(set types.Set, value string) bool {
	value = strings.Trim(value, `"`)
	for _, v := range set.Elements() {
		elementStr := strings.Trim(v.String(), `"`)
		if elementStr == value {
			return true
		}
	}
	return false
}

// GetDataSourceByValue is a helper function that gathers data based on all data gathered by the datasource.
//
// Parameters:
// - fields: The fields to filter the data.
// - allData: The data to be filtered.
//
// Returns:
// - []interface{}: The filtered data.
// - error: An error if any occurred.
func GetDataSourceByValue(fields interface{}, allData interface{}) ([]interface{}, error) {

	if isPointer(fields) || isPointer(allData) {
		return nil, fmt.Errorf("Pointers are not supported")
	}

	filteredArray := reflect.ValueOf(allData)
	fieldsArray := reflect.ValueOf(fields)
	var err error

	for j := 0; j < fieldsArray.NumField(); j++ {

		field := fieldsArray.Type().Field(j).Name
		fieldValue := fieldsArray.FieldByName(field)

		if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
			if fieldValue.IsNil() {
				continue
			}
		} else {
			if fieldValue.IsZero() {
				continue
			}
		}

		filteredArray, err = FilterByField(filteredArray, fieldValue, field)

		if err != nil {
			return nil, err
		}

	}

	allFilteredData := make([]interface{}, filteredArray.Len())
	for i := 0; i < filteredArray.Len(); i++ {
		allFilteredData[i] = filteredArray.Index(i).Interface()

	}

	return allFilteredData, nil

}

// FilterByField filters the array of data sources based on the provided field.
//
// Parameters:
// - dataSources: The array of data sources to filter.
// - fieldValue: The value to filter the data sources by.
// - field: The name of the field to filter by.
//
// Returns:
// - reflect.Value: The filtered array of data sources.
// - error: An error if any occurred.
func FilterByField(dataSources reflect.Value, fieldValue reflect.Value, field string) (reflect.Value, error) {
	filteredData := reflect.MakeSlice(dataSources.Type(), 0, dataSources.Len())

	for i := 0; i < dataSources.Len(); i++ {

		dataSource := dataSources.Index(i).Interface()

		dataSourceValue := reflect.ValueOf(dataSource)
		fieldValueInDataSource := dataSourceValue.FieldByName(field)

		if fieldValue.Kind() == reflect.Slice || fieldValue.Kind() == reflect.Array {
			for n := 0; n < fieldValue.Len(); n++ {

				interFieldValue, err := CheckAndConvertValue(fieldValue.Index(n))
				if err != nil {
					return reflect.Zero(nil), err
				}
				// if field is not found in the data source then break and continue
				if !fieldValueInDataSource.IsValid() || !interFieldValue.IsValid() {
					break
				}

				if fieldValueInDataSource.Interface() == interFieldValue.Interface() {
					filteredData = reflect.Append(filteredData, reflect.ValueOf(dataSource))
				}
			}
		} else {
			interFieldValue, err := CheckAndConvertValue(fieldValue)
			if err != nil {
				return reflect.Zero(nil), err
			}

			if fieldValueInDataSource.Interface() == interFieldValue.Interface() {
				filteredData = reflect.Append(filteredData, reflect.ValueOf(dataSource))
			}
		}
	}

	return filteredData, nil
}

// CheckAndConvertValue converts a reflect.Value to an attr.Type.
//
// It takes in a reflect.Value and checks its type. If the type is a
// types.StringType, it converts the input to a string, trims the quotes,
// and returns the resulting string. If the type is a types.Int64Type, it
// converts the input to an int, and returns the resulting int. If the
// type is a types.BoolType, it converts the input to a bool, and returns
// the resulting bool. If the type is none of the above, it returns an
// error.
//
// Returns:
// - reflect.Value: The converted value.
// - error: An error if the input type is not recognized.
func CheckAndConvertValue(input reflect.Value) (reflect.Value, error) {
	var valueRef reflect.Value
	switch ConvertType(input.Type()) {
	case types.StringType:
		value := fmt.Sprintf("%v", input)
		value = strings.Trim(value, "\"")
		valueRef = reflect.ValueOf(value)

		return valueRef, nil
	case types.Int64Type:
		value, err := strconv.Atoi(fmt.Sprintf("%v", input))
		if err != nil {
			return valueRef, nil
		}
		valueRef = reflect.ValueOf(value)

		return valueRef, nil
	case types.BoolType:
		value, err := strconv.ParseBool(fmt.Sprintf("%v", input))
		if err != nil {
			return valueRef, nil
		}
		valueRef = reflect.ValueOf(value)

		return valueRef, nil
	}

	return valueRef, fmt.Errorf("Value cannot be converted: %v", input)
}

// GenerateSchemaAttributes generates schema attributes based on a map of attribute names and respective types.
//
// The function takes a map of attributes, where each attribute is a map of attribute types and a boolean flag
// indicating whether the attribute is a set. The function iterates over the attributes and for each attribute,
// it generates a schema attribute using the SchemaAttributeGeneration function. The generated schema attributes
// are stored in a map, where the attribute name is the key.
//
// The function returns the generated schema attributes as a map, where the attribute name is the key and the
// corresponding schema attribute is the value.
func GenerateSchemaAttributes(attributes map[string]map[attr.Type]bool) map[string]schema.Attribute {
	schemaAttributes := make(map[string]schema.Attribute)
	for field, attrMap := range attributes {
		for attrType, isSet := range attrMap {
			schemaAttributes[field] = SchemaAttributeGeneration(field, attrType, isSet)
		}
	}
	tflog.Info(context.Background(), fmt.Sprintf("Generated Schema Attributes: %v", schemaAttributes))
	return schemaAttributes
}

// SchemaAttributeGeneration generates a schema attribute based on the type.
//
// It takes in a field name, attribute type, and a boolean flag indicating
// whether the attribute is a set. If the attribute is a set, it returns a
// schema.SetAttribute with the specified element type and validators. If the
// attribute is not a set, it returns a schema.Attribute of the specified type
// with the field name and description.
//
// Returns a schema.Attribute.
func SchemaAttributeGeneration(field string, attrType attr.Type, isSet bool) schema.Attribute {

	if isSet {
		return schema.SetAttribute{
			Description:         "List of " + field,
			MarkdownDescription: "List of " + field,
			ElementType:         attrType,
			Optional:            true,
			Validators: []validator.Set{
				setvalidator.SizeAtLeast(1),
			},
		}
	}
	switch attrType {
	case types.StringType:
		return schema.StringAttribute{
			Description:         "Value for" + field,
			MarkdownDescription: "Value for " + field,
			Optional:            true,
		}
	case types.Int64Type:
		return schema.Int64Attribute{
			Description:         "Value for " + field,
			MarkdownDescription: "Value for " + field,
			Optional:            true,
		}
	case types.BoolType:
		return schema.BoolAttribute{
			Description:         "Value for " + field,
			MarkdownDescription: "Value for " + field,
			Optional:            true,
		}
	}
	return nil
}

// TypeToMap converts any param into the specified type using its TFSDK tag.
//
// The function takes an interface{} parameter `t` and returns a map[string]map[attr.Type]bool.
// The map contains the field names as keys and a nested map as values.
// The nested map contains the converted type of the field as the key and a boolean value.
// The boolean value is true if the field is a slice or array, and false otherwise.
//
// Parameters:
// - t: The interface{} parameter to be converted.
//
// Returns:
// - map[string]map[attr.Type]bool: The converted map.
func TypeToMap(t interface{}) map[string]map[attr.Type]bool {
	r := reflect.TypeOf(t)
	m := make(map[string]map[attr.Type]bool)

	for i := 0; i < r.NumField(); i++ {
		field := r.Field(i)
		convType := ConvertType(field.Type)
		mTwo := make(map[attr.Type]bool)
		if convType == nil {
			continue
		} else if field.Type.Kind() == reflect.Slice || field.Type.Kind() == reflect.Array {
			mTwo[convType] = true
		} else {
			mTwo[convType] = false
		}
		m[field.Tag.Get("tfsdk")] = mTwo
	}

	return m
}

// ConvertType converts a reflect.Type to an attr.Type.
//
// It takes in a reflect.Type and checks its kind. If the type is a
// slice or array, it recursively calls itself with the element type.
// It then checks the name of the type and returns the corresponding
// attr.Type. If the type is none of the above, it returns nil.
//
// Parameters:
// - intialType: The reflect.Type to be converted.
//
// Returns:
// - attr.Type: The converted attr.Type.
func ConvertType(intialType reflect.Type) attr.Type {
	if intialType.Kind() == reflect.Slice || intialType.Kind() == reflect.Array {
		return ConvertType(intialType.Elem())
	}
	switch intialType.Name() {
	case "StringValue":
		return types.StringType
	case "Int64Value":
		return types.Int64Type
	case "BoolValue":
		return types.BoolType
	}

	return nil
}

// isPointer checks if the given value is a pointer.
//
// value: The value to check.
// Returns: A boolean indicating whether the value is a pointer.
func isPointer(value interface{}) bool {
	return reflect.ValueOf(value).Kind() == reflect.Ptr
}

// CopyFields copy the source of a struct to destination of struct with terraform types.
func CopyFields(ctx context.Context, source, destination interface{}) error {
	tflog.Debug(ctx, "Copy fields", map[string]interface{}{
		"source":      source,
		"destination": destination,
	})
	sourceValue := reflect.ValueOf(source)
	destinationValue := reflect.ValueOf(destination)

	// Check if destination is a pointer to a struct
	if destinationValue.Kind() != reflect.Ptr || destinationValue.Elem().Kind() != reflect.Struct {
		return fmt.Errorf("destination is not a pointer to a struct")
	}

	// if source is a pointer, use the Elem() method to get the value that the pointer points to
	if sourceValue.Kind() == reflect.Ptr {
		sourceValue = sourceValue.Elem()
	}

	if sourceValue.Kind() != reflect.Struct {
		return fmt.Errorf("source is not a struct")
	}

	// Get the type of the destination struct
	// destinationType := destinationValue.Elem().Type()
	for i := 0; i < sourceValue.NumField(); i++ {
		sourceFieldTag := getFieldJSONTag(sourceValue, i)

		tflog.Debug(ctx, "Converting source field", map[string]interface{}{
			"sourceFieldTag":  sourceFieldTag,
			"sourceFieldKind": sourceValue.Field(i).Kind().String(),
		})

		sourceField := sourceValue.Field(i)
		if sourceField.Kind() == reflect.Ptr {
			sourceField = sourceField.Elem()
		}
		if !sourceField.IsValid() {
			tflog.Error(ctx, "source field is not valid", map[string]interface{}{
				"sourceFieldTag": sourceFieldTag,
				"sourceField":    sourceField,
			})
			continue
		}

		destinationField := getFieldByTfTag(destinationValue.Elem(), sourceFieldTag)
		if destinationField.IsValid() && destinationField.CanSet() {

			tflog.Debug(ctx, "debugging source field", map[string]interface{}{
				"sourceField Interface": sourceField.Interface(),
			})
			// Convert the source value to the type of the destination field dynamically
			var destinationFieldValue attr.Value

			switch sourceField.Kind() {
			case reflect.String:
				destinationFieldValue = types.StringValue(sourceField.String())
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				destinationFieldValue = types.Int64Value(sourceField.Int())
			case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
				destinationFieldValue = types.Int64Value(sourceField.Int())
			case reflect.Float32, reflect.Float64:
				// destinationFieldValue = types.Float64Value(sourceField.Float())
				destinationFieldValue = types.NumberValue(big.NewFloat(sourceField.Float()))
			case reflect.Bool:
				destinationFieldValue = types.BoolValue(sourceField.Bool())
			case reflect.Array, reflect.Slice:
				if destinationField.Type().Kind() == reflect.Slice {
					arr := reflect.ValueOf(sourceField.Interface())
					slice := reflect.MakeSlice(destinationField.Type(), arr.Len(), arr.Cap())
					for index := 0; index < arr.Len(); index++ {
						value := arr.Index(index)
						v := slice.Index(index)
						switch v.Kind() {
						case reflect.Ptr:
							newDes := reflect.New(v.Type().Elem()).Interface()
							err := CopyFields(ctx, value.Interface(), newDes)
							if err != nil {
								return err
							}
							slice.Index(index).Set(reflect.ValueOf(newDes))
						case reflect.Struct:
							newDes := reflect.New(v.Type()).Interface()
							err := CopyFields(ctx, value.Interface(), newDes)
							if err != nil {
								return err
							}
							slice.Index(index).Set(reflect.ValueOf(newDes).Elem())
						}
					}
					destinationField.Set(slice)
				} else {
					destinationFieldValue = copySliceToTargetField(ctx, sourceField.Interface())
				}
			case reflect.Struct:
				// placeholder for improvement, need to consider both go struct and types.Object
				switch destinationField.Kind() {
				case reflect.Ptr:
					newDes := reflect.New(destinationField.Type().Elem()).Interface()
					err := CopyFields(ctx, sourceField.Interface(), newDes)
					if err != nil {
						return err
					}
					destinationField.Set(reflect.ValueOf(newDes))
				case reflect.Struct:
					newDes := reflect.New(destinationField.Type()).Interface()
					err := CopyFields(ctx, sourceField.Interface(), newDes)
					if err != nil {
						return err
					}
					destinationField.Set(reflect.ValueOf(newDes).Elem())
				}
				continue

			default:
				tflog.Error(ctx, "unsupported source field type", map[string]interface{}{
					"sourceField": sourceField,
				})
				continue
			}

			if destinationField.Type() == reflect.TypeOf(destinationFieldValue) {
				destinationField.Set(reflect.ValueOf(destinationFieldValue))
			}
		}
	}

	return nil
}

func getFieldJSONTag(sourceValue reflect.Value, i int) string {
	sourceFieldTag := sourceValue.Type().Field(i).Tag.Get("json")
	sourceFieldTag = strings.TrimSuffix(sourceFieldTag, ",omitempty")
	return sourceFieldTag
}

func getFieldByTfTag(destinationValue reflect.Value, tagValue string) reflect.Value {
	for j := 0; j < destinationValue.NumField(); j++ {
		field := destinationValue.Type().Field(j)
		if field.Tag.Get("tfsdk") == tagValue || field.Tag.Get("json") == tagValue {
			return destinationValue.Field(j)
		}
	}
	return reflect.Value{}
}

func copySliceToTargetField(ctx context.Context, fields interface{}) attr.Value {
	var objects []attr.Value
	attrTypeMap := make(map[string]attr.Type)

	// get the attrType for Object
	structElem := reflect.ValueOf(fields).Type().Elem()
	switch structElem.Kind() {
	case reflect.String:
		listValue, _ := types.ListValueFrom(ctx, types.StringType, fields)
		return listValue
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		listValue, _ := types.ListValueFrom(ctx, types.Int64Type, fields)
		return listValue
	case reflect.Float32, reflect.Float64:
		listValue, _ := types.ListValueFrom(ctx, types.Float64Type, fields)
		return listValue
	case reflect.Bool:
		listValue, _ := types.ListValueFrom(ctx, types.BoolType, fields)
		return listValue
	case reflect.Struct:
		for fieldIndex := 0; fieldIndex < structElem.NumField(); fieldIndex++ {
			field := structElem.Field(fieldIndex)
			tag := field.Tag.Get("json")
			tag = strings.TrimSuffix(tag, ",omitempty")
			fieldType := field.Type
			if fieldType.Kind() == reflect.Ptr {
				fieldType = fieldType.Elem()
			}

			switch fieldType.Kind() {
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				attrTypeMap[tag] = types.Int64Type
			case reflect.String:
				attrTypeMap[tag] = types.StringType
			case reflect.Float32, reflect.Float64:
				attrTypeMap[tag] = types.NumberType
			}
		}
		// iterate the slice
		arr := reflect.ValueOf(fields)
		for index := 0; index < arr.Len(); index++ {
			valueMap := make(map[string]attr.Value)
			// iterate the fields
			elem := arr.Index(index)
			for fieldIndex := 0; fieldIndex < elem.NumField(); fieldIndex++ {
				tag := elem.Type().Field(fieldIndex).Tag.Get("json")
				tag = strings.TrimSuffix(tag, ",omitempty")
				eleField := elem.Field(fieldIndex)
				eleFieldType := eleField.Type()
				if eleFieldType.Kind() == reflect.Ptr {
					eleFieldType = eleFieldType.Elem()
					eleField = eleField.Elem()
				}
				switch eleFieldType.Kind() {
				case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
					valueMap[tag] = types.Int64Value(eleField.Int())
				case reflect.String:
					valueMap[tag] = types.StringValue(eleField.String())
				case reflect.Float32, reflect.Float64:
					valueMap[tag] = types.NumberValue(big.NewFloat(eleField.Float()))
				}
			}
			object, _ := types.ObjectValue(attrTypeMap, valueMap)
			objects = append(objects, object)
		}
		listValue, _ := types.ListValue(types.ObjectType{AttrTypes: attrTypeMap}, objects)
		return listValue
	}
	return nil
}

// ConvertSlice converts a slice of interface{} to a slice of T.
func ConvertSlice[T any](input []interface{}) ([]T, error) {
	var result []T
	for _, item := range input {
		v, ok := item.(T)
		if !ok {
			return nil, fmt.Errorf("element %v cannot be converted to the target type", v)
		}
		result = append(result, v)

	}
	return result, nil
}
