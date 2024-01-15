package util

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestEscape(t *testing.T) {
	testStr := `a=<script>console.log("Test")</script>`
	assert.Equal(t, `a=&lt;script&gt;console.log(Test)&lt;/script&gt;`, Escape(testStr))
}
