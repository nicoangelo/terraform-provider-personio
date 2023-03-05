package adapter

import (
	"fmt"
	"strings"
	"time"

	personio "github.com/giantswarm/personio-go/v1"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

// convertAnyPersonioAttrToString inspects the underlying type returned from the Personio API
// and uses an appropriate conversion mechanism to convert it to a string.
// Conventions:
//   - integer, decimal: fmt.Sprintf()
//   - standard, multiline, link, list: direct conversion to string
//   - date: fmt.Sprintf() in UTC
func convertAnyPersonioAttrToString(v personio.Attribute) types.String {
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
		return convertPersonioAttrToString(v)
	case "date":
		return convertPersonioAttrToDateString(v)
	}
	return types.StringNull()
}

func convertPersonioTagsToStrings(v personio.Attribute) []string {
	if v.Value == nil {
		return []string{}
	}
	strValue, ok := v.Value.(string)
	if ok {
		return strings.Split(strValue, ",")
	}
	return []string{}
}

// convertPersonioAttrToString converts any standard, multiline, link, list API value
// to a Terraform string value. If the value is null, types.StringNull is returned.
func convertPersonioAttrToString(v personio.Attribute) types.String {
	if v.Value == nil {
		return types.StringNull()
	}
	strVal, ok := v.Value.(string)
	if ok {
		return types.StringValue(strVal)
	}
	return types.StringNull()
}

// convertPersonioAttrToString converts a integer API value
// to a Terraform Int64 value. If the value is null, types.Int64Null is returned.
func convertPersonioAttrToInt(v personio.Attribute) types.Int64 {
	if v.Value == nil {
		return types.Int64Null()
	}
	intVal, ok := v.Value.(int64)
	if ok {
		return types.Int64Value(intVal)
	}
	return types.Int64Null()
}

// convertPersonioAttrToString converts a decimal API value
// to a Terraform Float64 value. If the value is null, types.Float64Null is returned.
func convertPersonioAttrToFloat(v personio.Attribute) types.Float64 {
	if v.Value == nil {
		return types.Float64Null()
	}
	decVal, ok := v.Value.(float64)
	if ok {
		return types.Float64Value(decVal)
	}
	return types.Float64Null()
}

// convertPersonioAttrToDateString converts a time API value
// to a Terraform String value in UTC timezone. If the value is null, types.StringNull is returned.
func convertPersonioAttrToDateString(v personio.Attribute) types.String {
	if v.Value == nil {
		return types.StringNull()
	}
	timeVal := v.GetTimeValue()
	if timeVal != nil {
		return types.StringValue((*timeVal).UTC().Format(time.RFC3339))
	}
	return types.StringNull()
}
