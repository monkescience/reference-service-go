//go:build tools
// +build tools

package internal

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../openapi/status-api.oapi-codegen.models.yaml ../openapi/status-api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../openapi/status-api.oapi-codegen.server.yaml ../openapi/status-api.yaml
