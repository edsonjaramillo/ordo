package cli

func filterCompletedArgs(items []string, args []string, argStart int) []string {
	if len(items) == 0 || len(args) <= argStart {
		return items
	}

	selected := make(map[string]struct{}, len(args)-argStart)
	for _, arg := range args[argStart:] {
		if arg == "" {
			continue
		}
		selected[arg] = struct{}{}
	}
	if len(selected) == 0 {
		return items
	}

	filtered := make([]string, 0, len(items))
	for _, item := range items {
		if _, ok := selected[item]; ok {
			continue
		}
		filtered = append(filtered, item)
	}
	return filtered
}
