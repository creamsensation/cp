package cp

import (
	"strings"
	"testing"
	
	"github.com/stretchr/testify/assert"
	
	"github.com/creamsensation/gox"
)

func TestEmail(t *testing.T) {
	e := &email{}
	title := "Test title"
	subject := "Test subject"
	from := "from@test.com"
	to := "to@test.com"
	body := gox.Text("test")
	testEmail := e.Title(title).
		Subject(subject).
		From(from).
		To(to).
		Body(body)
	r := testEmail.String()
	assert.True(t, strings.Contains(r, "From: "+from))
	assert.True(t, strings.Contains(r, "To: "+to))
	assert.True(t, strings.Contains(r, "Subject: "+subject))
	assert.True(t, strings.Contains(r, gox.Render(body)))
}
