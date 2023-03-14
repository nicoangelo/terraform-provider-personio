package utils

import (
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
)

func ReplaceAttribute(attrs map[string]schema.Attribute, key string, newAttr schema.Attribute) map[string]schema.Attribute {
	attrs[key] = newAttr
	return attrs
}
