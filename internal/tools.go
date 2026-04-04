//go:build tools

package internal

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../openapi/import-api.oapi-codegen.server.yaml ../openapi/import-api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../openapi/import-api.oapi-codegen.client.yaml ../openapi/import-api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../openapi/pokemon-api.oapi-codegen.server.yaml ../openapi/pokemon-api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../openapi/pokeapi.oapi-codegen.client.yaml ../openapi/pokeapi.yaml
