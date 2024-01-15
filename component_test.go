package cp

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/cp/internal/route"
	"github.com/creamsensation/gox"
)

type testComponent struct {
	Component
}

func (c *testComponent) Name() string {
	return "test"
}

func (c *testComponent) Model() {
}

func (c *testComponent) Node() gox.Node {
	return gox.Text("test")
}

func TestCreatePrefixedComponentName(t *testing.T) {
	type testCase struct {
		*control
		name     string
		expected string
	}
	testCases := make([]testCase, 0)
	testCases = append(
		testCases, testCase{
			control: func() *control {
				c := &control{}
				c.component = &testComponent{Component: c}
				return c
			}(),
			name:     "NoControllerNoModule",
			expected: "test",
		},
	)
	testCases = append(
		testCases, testCase{
			control: func() *control {
				c := &control{}
				c.component = &testComponent{Component: c}
				c.route = route.Route{Controller: "controller"}
				return c
			}(),
			name:     "WithController",
			expected: "controller_test",
		},
	)
	testCases = append(
		testCases, testCase{
			control: func() *control {
				c := &control{}
				c.component = &testComponent{Component: c}
				c.route = route.Route{Module: "module"}
				return c
			}(),
			name:     "WithModule",
			expected: "module_test",
		},
	)
	testCases = append(
		testCases, testCase{
			control: func() *control {
				c := &control{}
				c.component = &testComponent{Component: c}
				c.route = route.Route{Module: "module", Controller: "controller"}
				return c
			}(),
			name:     "WithControllerAndModule",
			expected: "module_controller_test",
		},
	)
	for _, tc := range testCases {
		t.Run(
			tc.name, func(t *testing.T) {
				assert.Equal(t, tc.expected, createPrefixedComponentName(tc.control))
			},
		)
	}
}
