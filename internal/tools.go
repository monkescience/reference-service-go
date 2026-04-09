//go:build tools

package internal

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
	_ "github.com/sqlc-dev/sqlc/cmd/sqlc"
)

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f ../sqlc.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../openapi/reference-api.oapi-codegen.server.yaml ../openapi/reference-api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../openapi/pokeapi.oapi-codegen.client.yaml ../openapi/pokeapi.yaml
