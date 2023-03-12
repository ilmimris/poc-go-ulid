package group

import (
	"github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/interfaces"
	v1 "github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/rest/handler/v1"
)

func InitRouterHealth(rest interfaces.Rest) {
	root := rest.GetRouter()
	root.Group("/")

	// path healthcheck
	InitHealthCheck(root)
}

func InitRouterMetrics(rest interfaces.Rest) {
	root := rest.GetRouter()
	root.Group("/")

	InitMetrics(root)
}

func InitRouterV1(rest interfaces.Rest, h v1.Handler) {
	// path root
	root := rest.GetRouter()
	root.Group("/")

	// define path v1
	// define all API under path v1
	_ = InitV1Group(root, h)
}
