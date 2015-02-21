package rainmaker_test

import (
	"github.com/pivotal-cf-experimental/rainmaker"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ServiceInstancesService", func() {
	var config rainmaker.Config
	var service *rainmaker.ServiceInstancesService
	var token string

	BeforeEach(func() {
		config = rainmaker.Config{
			Host: fakeCloudController.URL(),
		}
		service = rainmaker.NewServiceInstancesService(config)
		token = "TOKEN"
	})

	Describe("Create/Get", func() {
		It("allows the user to create and then fetch a service instance", func() {
			name := "my-service-instance"
			planGUID := "service-plan-guid"
			spaceGUID := "space-guid"
			instance, err := service.Create(name, planGUID, spaceGUID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(instance).To(BeAssignableToTypeOf(rainmaker.ServiceInstance{}))
			Expect(instance.GUID).NotTo(BeEmpty())
			Expect(instance.Name).To(Equal(name))
			Expect(instance.PlanGUID).To(Equal(planGUID))
			Expect(instance.SpaceGUID).To(Equal(spaceGUID))

			service = rainmaker.NewServiceInstancesService(config)
			fetchedInstance, err := service.Get(instance.GUID, token)
			Expect(err).NotTo(HaveOccurred())
			Expect(fetchedInstance).To(Equal(instance))
		})
	})
})
