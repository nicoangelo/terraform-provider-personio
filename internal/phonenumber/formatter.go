package phonenumber

import "github.com/nyaruka/phonenumbers"

// ParseAndFormatPhonenumber parses the input i and interprets
// it as a phone number in the given region.
// The number is returned formatted in national format,
// in case of errors the original input is returned.
func ParseAndFormatPhonenumber(i string, region string) string {
	pn, err := phonenumbers.Parse(i, region)
	if err == nil {
		return phonenumbers.Format(pn, phonenumbers.NATIONAL)
	}
	return i
}
