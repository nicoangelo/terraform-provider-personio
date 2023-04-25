package formatter

import (
	"github.com/nyaruka/phonenumbers"
)

const (
	PHONENUMBER_DEFAULT_FORMAT = "INTERNATIONAL"
)

type PhoneNumberFormatter struct {
	DefaultRegion     string
	PhoneNumberFormat phonenumbers.PhoneNumberFormat
}

func (pnf *PhoneNumberFormatter) Format(input string) string {
	pn, err := phonenumbers.Parse(input, pnf.DefaultRegion)
	if err == nil {
		return phonenumbers.Format(pn, phonenumbers.INTERNATIONAL)
	}
	return input
}

func (pnf *PhoneNumberFormatter) Configure(cfg *PhoneNumberConfig) {
	regions := phonenumbers.GetSupportedRegions()

	if _, ok := regions[cfg.DefaultRegion.ValueString()]; ok {
		pnf.DefaultRegion = cfg.DefaultRegion.ValueString()
	} else {
		pnf.DefaultRegion = phonenumbers.UNKNOWN_REGION
	}

	if cfg.Format.IsNull() {
		pnf.PhoneNumberFormat = pnf.getFormatFromString(PHONENUMBER_DEFAULT_FORMAT)
	} else {
		pnf.PhoneNumberFormat = pnf.getFormatFromString(cfg.Format.ValueString())
	}
}

func (pnf *PhoneNumberFormatter) getFormatFromString(format string) phonenumbers.PhoneNumberFormat {
	switch format {
	case "E164":
		return phonenumbers.E164
	case "NATIONAL":
		return phonenumbers.NATIONAL
	case "RFC3966":
		return phonenumbers.RFC3966
	}
	return phonenumbers.INTERNATIONAL
}
