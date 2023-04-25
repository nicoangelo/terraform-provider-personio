package formatter

import "github.com/hashicorp/terraform-plugin-framework/types"

type AttributeFormatter interface {
	Format(input string) string
}

type FormatterCollection struct {
	formatters []Formatter
}

type Formatter struct {
	formatter    AttributeFormatter
	attributeKey string
}

func (fc *FormatterCollection) FromConfig(cfg []FormatterConfig) {
	fc.formatters = []Formatter{}

	for _, v := range cfg {
		var f AttributeFormatter
		if v.PhoneNumber != nil {
			pnf := &PhoneNumberFormatter{}
			pnf.Configure(v.PhoneNumber)
			f = pnf
		}
		if f != nil {
			fc.formatters = append(fc.formatters, Formatter{
				formatter:    f,
				attributeKey: v.AttributeKey.ValueString(),
			})
		}
	}
}

func (fc *FormatterCollection) FormatAll(attrs map[string]types.String) {
	for _, v := range fc.formatters {
		attr, ok := attrs[v.attributeKey]
		if !ok || attr.IsNull() || attr.IsUnknown() {
			continue
		}

		attrs[v.attributeKey] = types.StringValue(v.formatter.Format(attr.ValueString()))
	}
}

type FormatterConfig struct {
	AttributeKey types.String       `tfsdk:"attribute"`
	PhoneNumber  *PhoneNumberConfig `tfsdk:"phonenumber"`
}

type PhoneNumberConfig struct {
	DefaultRegion types.String `tfsdk:"default_region"`
	Format        types.String `tfsdk:"format"`
}
