package cp

import (
	"reflect"
	"testing"
	
	"github.com/stretchr/testify/assert"
)

type testService struct {
	Test string
}

type testServiceInterface interface{}

func (s *testService) Provide(control Control) Provider {
	return &testService{}
}

func (s *testService) GetValue() string {
	return "abc"
}

func TestDi(t *testing.T) {
	t.Run(
		"method", func(t *testing.T) {
			c := &core{
				deps: make(map[string]*Dependency),
			}
			c.Container(
				Register(&testService{}),
			)
			ctrl := &control{core: c}
			
			s := Provide[testService](ctrl)
			assert.Equal(t, "abc", s.GetValue())
		},
	)
	t.Run(
		"singleton", func(t *testing.T) {
			c := &core{
				deps: make(map[string]*Dependency),
			}
			c.Container(
				Register(&testService{Test: "test"}).Singleton(),
			)
			ctrl := &control{core: c}
			s := Provide[testService](ctrl)
			assert.Equal(t, "test", s.Test)
		},
	)
	t.Run(
		"custom name", func(t *testing.T) {
			c := &core{
				deps: make(map[string]*Dependency),
			}
			c.Container(
				Register(&testService{Test: "test 1"}).Singleton().Name(Interface[testServiceInterface]()),
			)
			ctrl := &control{core: c}
			s := Provide[testService](ctrl, Interface[testServiceInterface]())
			assert.Equal(t, "test 1", s.Test)
		},
	)
	t.Run(
		"autoinject", func(t *testing.T) {
			c := &core{
				deps: make(map[string]*Dependency),
			}
			c.Container()
			ctrl := &control{core: c}
			s := provide[testService](ctrl, reflect.TypeOf(&testService{})).(*testService)
			assert.Equal(t, "abc", s.GetValue())
		},
	)
}
