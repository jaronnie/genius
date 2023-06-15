package genius

func removeDuplicateElement(slices []string) []string {
	result := make([]string, 0, len(slices))
	temp := map[string]struct{}{}
	for _, item := range slices {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}
