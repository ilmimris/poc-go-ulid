package registry

type ServiceRegistry struct {
	// contains filtered or unexported fields

}

type OptServiceRegistry func(registry *ServiceRegistry)

func NewServiceRegistry(opts ...OptServiceRegistry) *ServiceRegistry {
	r := &ServiceRegistry{}
	for _, opt := range opts {
		opt(r)
	}
	return r
}
