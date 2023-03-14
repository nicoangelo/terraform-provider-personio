package adapter

import (
	"fmt"
	"math/big"
	"strings"
	"time"

	personio "github.com/giantswarm/personio-go/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// convertAnyAttrToString inspects the underlying type returned from the Personio API
// and uses an appropriate conversion mechanism to convert it to a string.
// Conventions:
//   - integer, decimal: fmt.Sprintf()
//   - standard, multiline, link, list: direct conversion to string
//   - date: fmt.Sprintf() in UTC
func convertAnyAttrToString(v personio.Attribute) types.String {
	// bail early
	if v.Value == nil {
		return types.StringNull()
	}
	switch v.Type {
	case "integer":
		intVal, ok := v.Value.(int64)
		if ok {
			return types.StringValue(fmt.Sprint(intVal))
		}
	case "decimal":
		decVal, ok := v.Value.(float64)
		if ok {
			return types.StringValue(fmt.Sprint(decVal))
		}
	case "standard", "multiline", "link", "list":
		return convertAttrToString(v)
	case "date":
		return convertAttrToDateString(v)
	}
	return types.StringNull()
}

func convertTagsToStrings(v personio.Attribute) (res []types.String) {
	if v.Value == nil {
		return []types.String{}
	}
	strValue, ok := v.Value.(string)
	if ok {
		res = make([]types.String, 0)
		for _, v := range strings.Split(strValue, ",") {
			res = append(res, types.StringValue(v))
		}
		return res
	}
	return []types.String{}
}

// convertAttrToString converts any standard, multiline, link, list API value
// to a Terraform string value. If the value is null, types.StringNull is returned.
func convertAttrToString(v personio.Attribute) types.String {
	if v.Value == nil {
		return types.StringNull()
	}
	strVal, ok := v.Value.(string)
	if ok {
		return types.StringValue(strVal)
	}
	return types.StringNull()
}

// convertAttrToNumber converts an integer API value
// to a Terraform Int64 value. If the value is null, types.Int64Null is returned.
func convertAttrToNumber(v personio.Attribute) types.Number {
	if v.Value == nil {
		return types.NumberNull()
	}
	intVal, ok := v.Value.(float64)
	if ok {
		return types.NumberValue(big.NewFloat(intVal))
	}
	return types.NumberNull()
}

// convertAttrToFloat converts a decimal API value
// to a Terraform Float64 value. If the value is null, types.Float64Null is returned.
func convertAttrToFloat(v personio.Attribute) types.Float64 {
	if v.Value == nil {
		return types.Float64Null()
	}
	decVal, ok := v.Value.(float64)
	if ok {
		return types.Float64Value(decVal)
	}
	return types.Float64Null()
}

// convertAttrToDateString converts a time API value
// to a Terraform String value in UTC timezone. If the value is null, types.StringNull is returned.
func convertAttrToDateString(v personio.Attribute) types.String {
	if v.Value == nil {
		return types.StringNull()
	}
	timeVal := v.GetTimeValue()
	if timeVal != nil {
		return types.StringValue((*timeVal).UTC().Format(time.RFC3339))
	}
	return types.StringNull()
}

// convertMapItemToString converts a specific attribute of a nested map API value (e.g. supervisor)
// to a Terraform String value. If the value is null, types.StringNull is returned.
func convertMapItemToString(v personio.Attribute, itemKey string) types.String {
	if v.Value == nil {
		return types.StringNull()
	}
	mapVal := v.GetMapValue()
	strVal, ok := mapVal[itemKey].(string)
	if ok {
		return types.StringValue(strVal)
	}
	return types.StringNull()
}

// convertMapItemToInt converts a specific attribute of a nested map API value (e.g. supervisor)
// to a Terraform number value. If the value is null, types.Int64Null is returned.
func convertMapItemToInt(v personio.Attribute, itemKey string) types.Int64 {
	if v.Value == nil {
		return types.Int64Null()
	}
	mapVal := v.GetMapValue()
	intVal, ok := mapVal[itemKey].(float64)
	if ok {
		return types.Int64Value(int64(intVal))
	}
	return types.Int64Null()
}

// convertNestedMapItemToString converts a specific attribute of a nested map API value (e.g. supervisor)
// to a Terraform String value. If the value is null, types.StringNull is returned.
func convertNestedMapItemToString(v personio.Attribute, itemKey string) types.String {
	if v.Value == nil {
		return types.StringNull()
	}
	mapVal := v.GetMapValue()[itemKey].(map[string]interface{})
	strVal, ok := mapVal["value"].(string)
	if ok {
		return types.StringValue(strVal)
	}
	return types.StringNull()
}

// convertNestedMapItemToInt converts a specific attribute of a nested map API value (e.g. supervisor)
// to a Terraform number value. If the value is null, types.Float64Null is returned.
func convertNestedMapItemToInt(v personio.Attribute, itemKey string) types.Int64 {
	if v.Value == nil {
		return types.Int64Null()
	}
	mapVal := v.GetMapValue()[itemKey].(map[string]interface{})
	// nested numbers are stored as float
	floatVal, ok := mapVal["value"].(float64)
	if ok {
		return types.Int64Value(int64(floatVal))
	}
	return types.Int64Null()
}
