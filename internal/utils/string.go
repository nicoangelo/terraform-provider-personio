package utils

// coalesceEmpty returns the first non-empty value given
// in the arguments. If none is found, an empty string is returned.
func CoalesceEmpty(values ...string) string {
	for _, v := range values {
		if v != "" {
			return v
		}
	}
	return ""
}
