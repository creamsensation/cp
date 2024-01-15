package cp

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

type testModule struct{}

func (r *testModule) Name() string {
	return "test"
}

func (r *testModule) Controllers() Controllers {
	return Controllers{
		&testController{},
	}
}

func TestModule(t *testing.T) {
	c := &core{}
	c.router = createRouter(c)
	r := createModule(c, &testModule{})
	r.run()
	assert.Equal(t, 1, len(c.router.builders))
}
