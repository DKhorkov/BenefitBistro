package structures

import (
	"net/http"
)

type RouteHandlersStructure struct {
	HomePage func(http.ResponseWriter, *http.Request)
}

type RouteHandlersNamesStructure struct {
	HomePage string
}

type ServerParametersStructure struct {
	Host string
	Port string
}

type URLPathsStructure struct {
	HomePage string
	StaticFiles string
}
