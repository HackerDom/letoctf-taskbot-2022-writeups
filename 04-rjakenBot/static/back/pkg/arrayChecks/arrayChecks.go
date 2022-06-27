package arrayChecks

func Contains[T comparable](arr []T, elem T) bool {
	for _, val := range arr {
		if val == elem {
			return true
		}
	}

	return false
}
