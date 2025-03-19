package tools

func StringSliceDiff(a, b []string) []string {
	set := make(map[string]bool)
	for _, v := range b {
		set[v] = true
	}

	var diff []string
	for _, v := range a {
		if !set[v] {
			diff = append(diff, v)
		}
	}
	return diff
}
