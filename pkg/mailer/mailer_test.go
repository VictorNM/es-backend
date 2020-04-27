package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const tpl = `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>{{.Title}}</title>
	</head>
	<body>
		{{range .Items}}<div>{{ . }}</div>{{else}}<div><strong>no rows</strong></div>{{end}}
	</body>
</html>`

type data struct {
	Title string
	Items []string
}

func TestParseTemplate(t *testing.T) {
	msg, err := parseTemplate(tpl, data{
		Title: "My page",
		Items: []string{
			"My photos",
			"My blog",
		},
	})

	assert.NoError(t, err)

	res := `<!DOCTYPE html>
<html>
	<head>
		<meta charset="UTF-8">
		<title>My page</title>
	</head>
	<body>
		<div>My photos</div><div>My blog</div>
	</body>
</html>`
	assert.Equal(t, res, string(msg))
}
