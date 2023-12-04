package route

type Config interface {
}

type PathValue interface{}

type config[T any] struct {
	configType int
	value      T
}

type pathConfig struct {
	method []string
	path   any
}

const (
	ConfigName = iota
	ConfigPath
	ConfigHandler
	ConfigLocalize
	ConfigLayout
)

func CreateConfig[T any](configType int, value T) Config {
	return config[T]{
		configType: configType,
		value:      value,
	}
}

func CreateMethodConfig(methods ...Config) Config {
	return config[pathConfig]{
		configType: ConfigPath,
		value: pathConfig{
			method: processMethods(methods...),
		},
	}
}

func CreatePathConfig(path PathValue) Config {
	return config[pathConfig]{
		configType: ConfigPath,
		value: pathConfig{
			path: path,
		},
	}
}

func CreatePath(method string, path ...PathValue) Config {
	var p PathValue
	if len(path) > 0 {
		p = path[0]
	}
	return config[pathConfig]{
		configType: ConfigPath,
		value: pathConfig{
			method: []string{method},
			path:   p,
		},
	}
}
