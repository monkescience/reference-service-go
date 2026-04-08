//go:build tools

package internal

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../openapi/reference-api.oapi-codegen.server.yaml ../openapi/reference-api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../openapi/pokeapi.oapi-codegen.client.yaml ../openapi/pokeapi.yaml
