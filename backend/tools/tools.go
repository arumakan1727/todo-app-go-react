//go:build tools

package tools

import (
	_ "github.com/deepmap/oapi-codegen/pkg/codegen"
	_ "github.com/k0kubun/pp/v3"
	_ "github.com/kyleconroy/sqlc/pkg/cli"
)
