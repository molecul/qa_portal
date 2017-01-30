package web

import (
	"html/template"
)

func initTemplatesFuncs() map[string]interface{} {
	return map[string]interface{}{
		"isNil": func(v interface{}) bool {
			return v == nil
		},
		"isNotNil": func(v interface{}) bool {
			return v != nil
		},
		"HTML": func(v string) template.HTML {
			return template.HTML(v)
		},
	}
}
