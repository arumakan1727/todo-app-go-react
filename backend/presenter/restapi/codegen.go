package restapi

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -config=.oapi-codegen.interfaces.yml ../../../api-spec/open-api.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -config=.oapi-codegen.schema.yml ../../../api-spec/open-api.yaml
//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -config=.oapi-codegen.client.yml ../../../api-spec/open-api.yaml
