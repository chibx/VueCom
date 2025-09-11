package server

import (
	"fmt"
	"strings"
)

func extractRouteParts(route string) []string {
	var routeLength = len(route)
	var routeParts []string

	fmt.Println(route, routeLength)
	if routeLength == 1 {
		// It's just "/"
		routeParts = []string{""}
	} else {
		var hasTrailingSlash = string(route[routeLength-1]) == "/"

		if hasTrailingSlash {
			route = route[:len(route)-1]
		}
		// First index is gonna be ""
		routeParts = strings.Split(route, "/")[1:]
	}

	return routeParts
}
