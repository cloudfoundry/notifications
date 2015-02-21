package fakes

type ServiceInstances struct {
	store map[string]ServiceInstance
}

func NewServiceInstances() *ServiceInstances {
	return &ServiceInstances{
		store: make(map[string]ServiceInstance),
	}
}

func (instances *ServiceInstances) Add(instance ServiceInstance) {
	instances.store[instance.GUID] = instance
}

func (instances *ServiceInstances) Get(guid string) (ServiceInstance, bool) {
	instance, ok := instances.store[guid]
	return instance, ok
}
