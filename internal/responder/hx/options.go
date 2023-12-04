package hx

import (
	. "hx"
)

type Options interface {
	Append() Options
	Delete() Options
	Modifier(modifier string) Options
	ScrollBottom() Options
	Swap(swap string) Options
}

type options struct {
	modifier string
	swap     string
	target   string
}

func (o *options) Append() Options {
	o.swap = SwapBeforeEnd
	return o
}

func (o *options) Delete() Options {
	o.swap = SwapDelete
	return o
}

func (o *options) Modifier(modifier string) Options {
	o.modifier = modifier
	return o
}

func (o *options) ScrollBottom() Options {
	o.modifier = EventScrollBottom
	return o
}

func (o *options) Swap(swap string) Options {
	o.swap = swap
	return o
}
