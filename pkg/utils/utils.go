package utils

func IsStringSliceContain(strSlice []string, key string) bool {
	for _, val := range strSlice {
		if val == key {
			return true
		}
	}
	return false
}
