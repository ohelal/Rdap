package docs

import (
	"encoding/json"
	"github.com/swaggo/swag"
)

type APIDocGenerator struct {
	Title       string
	Description string
	Version     string
	BasePath    string
}

func NewAPIDocGenerator(title, description, version, basePath string) *APIDocGenerator {
	return &APIDocGenerator{
		Title:       title,
		Description: description,
		Version:     version,
		BasePath:    basePath,
	}
}

func (g *APIDocGenerator) GenerateSwaggerDocs() *swag.Spec {
	return &swag.Spec{
		InfoProps: swag.InfoProps{
			Title:       g.Title,
			Description: g.Description,
			Version:     g.Version,
		},
		BasePath: g.BasePath,
		Schemes:  []string{"http", "https"},
		Consumes: []string{"application/json"},
		Produces: []string{"application/json"},
	}
}
