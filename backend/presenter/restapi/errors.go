package restapi

import "github.com/labstack/echo/v4"

var (
	Err400InvalidReqContentType = echo.NewHTTPError(400, "invalid Content-Type")
	Err400EmptyReqBody          = echo.NewHTTPError(400, "empty body")
)
