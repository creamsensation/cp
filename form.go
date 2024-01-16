package cp

type Form interface {
	Errors(errors map[string]string) Form
}

type formManager struct {
	errors map[string]string
}

func (f *formManager) Errors(errors map[string]string) Form {
	f.errors = errors
	return f
}
