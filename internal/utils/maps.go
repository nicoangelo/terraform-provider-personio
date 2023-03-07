package utils

func MergeMaps[TKey comparable, TValue interface{}](maps ...map[TKey]TValue) (res map[TKey]TValue) {
	res = make(map[TKey]TValue)
	for _, m := range maps {
		for k, v := range m {
			res[k] = v
		}
	}
	return res
}
