package main

import (
	"fmt"
	"log"
	"runtime/debug"
	"sort"
	"strings"

	"github.com/arumakan1727/todo-app-go-react/presenter/restapi"
	"github.com/labstack/echo/v4"
)

func filterMap[T any](a []*T, predicate func(*T) *T) []*T {
	res := make([]*T, 0, len(a))
	for i := range a {
		if v := predicate(a[i]); v != nil {
			res = append(res, v)
		}
	}
	return res
}

func getSelfModName() string {
	bi, ok := debug.ReadBuildInfo()
	if !ok {
		log.Fatalf("getSelfModName(): cannot read BuildInfo.")
	}
	return strings.TrimSuffix(bi.Path, "/cmd/api")
}

func CmdRoutes() {
	selfModName := getSelfModName()
	routes := filterMap(restapi.Routes(), func(r *echo.Route) *echo.Route {
		if name, ok := strings.CutPrefix(r.Name, selfModName); ok {
			return &echo.Route{
				Method: r.Method,
				Path:   r.Path,
				Name:   name,
			}
		}
		return nil
	})

	sort.Slice(routes, func(i, j int) bool {
		cmp := strings.Compare(routes[i].Path, routes[j].Path)
		if cmp != 0 {
			return cmp < 0
		}
		return strings.Compare(routes[i].Method, routes[j].Method) < 0
	})

	for _, r := range routes {
		fmt.Printf("%-8s %-24s %20s\n", r.Method, r.Path, r.Name)
	}
}
