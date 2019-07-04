package bean

import (
	"github.com/Xuanwo/beancollect/transform"
	"github.com/Xuanwo/beancollect/types"
	"github.com/sirupsen/logrus"
	"os"
	"sort"
	"text/template"
)

const content = `{{ .Time.Format "2006-01-02" }} {{ .Flag }} "{{ .Payee }}" "{{ .Narration }}"
    {{ index .Accounts 0 }} {{ .Amount }} {{ .Currency }}
{{ if gt (len .Accounts) 1 }}    {{ index .Accounts 1 }}
{{ end }}
`

var tmpl = template.Must(template.New("bean").Parse(content))

// Generate will generate the transactions into bean.
func Generate(c *types.Config, t *types.Transactions) {
	sort.Sort(t)

	for _, v := range c.Rules {
		transform.Execute(v, t)
	}

	for _, v := range *t {
		err := tmpl.Execute(os.Stdout, v)
		if err != nil {
			logrus.Fatalf("Template execute failed for %v", err)
		}
	}
}
