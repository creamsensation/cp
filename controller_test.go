package cp

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

type testController struct{}

func (r *testController) Name() string {
	return "test"
}

func (r *testController) Routes() Routes {
	return Routes{
		Route(Get("test")),
	}
}

func TestController(t *testing.T) {
	c := &core{}
	c.router = createRouter(c)
	r := createController(c, &testController{})
	r.run()
	assert.Equal(t, 1, len(c.router.builders))
}
