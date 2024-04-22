package template

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/giantswarm/microerror"
	"github.com/rs/zerolog/log"
)

func Execute(tmplFile string, data any) (string, error) {
	log.Debug().Msg("Creating file from template")
	tmpl, err := template.New(tmplFile).Parse(tmplFile)
	if err != nil {
		return "", fmt.Errorf("failed to parse template file %s: %w", tmplFile, microerror.Mask(err))
	}

	b := &bytes.Buffer{}
	err = tmpl.Execute(b, data)
	if err != nil {
		return "", fmt.Errorf("failed to execute template %s: %w", tmplFile, microerror.Mask(err))
	}

	return b.String(), nil
}
