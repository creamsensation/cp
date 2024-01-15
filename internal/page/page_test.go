package page

import (
	"testing"
	
	"github.com/stretchr/testify/assert"
)

func TestPage(t *testing.T) {
	p := New()
	title := "Test title"
	keywords := "a,b,c"
	description := "Test description"
	metaName, metaContent := "test name", "test content"
	
	t.Run(
		"set", func(t *testing.T) {
			p.Set().Title(title).
				Keywords(keywords).
				Description(description).
				Meta(metaName, metaContent)
		},
	)
	
	t.Run(
		"get", func(t *testing.T) {
			assert.Equal(t, title, p.Get().Title())
			assert.Equal(t, keywords, p.Get().Keywords())
			assert.Equal(t, description, p.Get().Description())
			assert.Equal(t, 1, len(p.Get().Metas()))
			assert.Equal(t, metaName, p.Get().Metas()[0][0])
			assert.Equal(t, metaContent, p.Get().Metas()[0][1])
		},
	)
}
