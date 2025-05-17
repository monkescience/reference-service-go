//go:build tools
// +build tools

package internal

import (
	_ "github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen"
)

//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../openapi/oapi-codegen.models.yaml ../openapi/v1_order_api.yaml
//go:generate go run github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen --config=../openapi/oapi-codegen.server.yaml ../openapi/v1_order_api.yaml
