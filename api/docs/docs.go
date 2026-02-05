package docs

import (
	_ "embed"
)

// OpenAPISpec contains the embedded OpenAPI specification.
//
//go:embed swagger.json
var OpenAPISpec []byte
