package route

func CreateFlatBuilders(builders ...*Builder) []*Builder {
	result := make([]*Builder, 0)
	for _, b := range builders {
		if len(b.Subroutes) > 0 {
			result = append(result, CreateFlatBuilders(b.Subroutes...)...)
			b.Subroutes = nil
		}
		result = append(result, b)
	}
	return result
}

func IsLocalized(configs ...Config) bool {
	for _, item := range configs {
		switch c := item.(type) {
		case config[pathConfig]:
			switch c.value.path.(type) {
			case map[string]any:
				return true
			}
		}
	}
	return false
}
