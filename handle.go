package cp

type Handle interface {
	Hx() HxHandle
}

type handle struct {
	control *control
}

func (h handle) Hx() HxHandle {
	return hxHandle{control: h.control}
}
