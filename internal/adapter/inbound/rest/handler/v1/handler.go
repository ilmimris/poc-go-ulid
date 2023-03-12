package handlerv1

import "github.com/ilmimris/poc-go-ulid/internal/adapter/inbound/rest/registry"

type Handler struct {
	serviceRegistry *registry.ServiceRegistry
}

type OptHandler struct {
	ServiceRegistry *registry.ServiceRegistry
}

func New(opt OptHandler) *Handler {
	return &Handler{
		serviceRegistry: opt.ServiceRegistry,
	}
}

func (h *Handler) GetServiceRegistry() *registry.ServiceRegistry {
	return h.serviceRegistry
}
