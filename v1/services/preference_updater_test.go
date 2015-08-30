package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PreferenceUpdater", func() {
	Describe("Execute", func() {
		var (
			doorOpen                   models.Unsubscribe
			barking                    models.Unsubscribe
			unsubscribesRepo           *mocks.UnsubscribesRepo
			kindsRepo                  *mocks.KindsRepo
			fakeGlobalUnsubscribesRepo *mocks.GlobalUnsubscribesRepo
			conn                       *mocks.Connection
			updater                    services.PreferenceUpdater
		)

		BeforeEach(func() {
			conn = mocks.NewConnection()
			unsubscribesRepo = mocks.NewUnsubscribesRepo()
			kindsRepo = mocks.NewKindsRepo()
			fakeGlobalUnsubscribesRepo = mocks.NewGlobalUnsubscribesRepo()
			updater = services.NewPreferenceUpdater(fakeGlobalUnsubscribesRepo, unsubscribesRepo, kindsRepo)
		})

		Context("when globally unsubscribing", func() {
			It("inserts a record into the global unsubscribes repo", func() {
				updater.Execute(conn, []models.Preference{}, true, "user-guid")
				Expect(fakeGlobalUnsubscribesRepo.SetCall.Receives.Unsubscribed).To(BeTrue())

				updater.Execute(conn, []models.Preference{}, false, "user-guid")
				Expect(fakeGlobalUnsubscribesRepo.SetCall.Receives.Unsubscribed).To(BeFalse())
			})

			Context("when the global unsubscribe repo errors", func() {
				It("returns the error", func() {
					fakeGlobalUnsubscribesRepo.SetCall.Returns.Error = errors.New("global unsubscribe db error")

					err := updater.Execute(conn, []models.Preference{}, true, "user-guid")
					Expect(err).To(MatchError(errors.New("global unsubscribe db error")))
				})
			})
		})

		Context("When unsubscribing from existing kinds of existing clients", func() {
			BeforeEach(func() {
				doorOpen = models.Unsubscribe{
					UserID:   "the-user",
					ClientID: "raptors",
					KindID:   "door-open",
				}

				barking = models.Unsubscribe{
					UserID:   "the-user",
					ClientID: "dogs",
					KindID:   "barking",
				}

				kindsRepo.Create(conn, models.Kind{
					ID:       "door-open",
					ClientID: "raptors",
				})

				kindsRepo.Create(conn, models.Kind{
					ID:       "barking",
					ClientID: "dogs",
				})

			})

			It("Adds New Unsubscribes to the unsubscribes Repo", func() {
				updater.Execute(conn, []models.Preference{
					{
						ClientID: "raptors",
						KindID:   "door-open",
						Email:    false,
					},

					{
						ClientID: "dogs",
						KindID:   "barking",
						Email:    false,
					},
				}, false, "the-user")

				unsubscribed, err := unsubscribesRepo.Get(conn, "the-user", "raptors", "door-open")
				Expect(err).NotTo(HaveOccurred())
				Expect(unsubscribed).To(BeTrue())

				unsubscribed, err = unsubscribesRepo.Get(conn, "the-user", "dogs", "barking")
				Expect(err).NotTo(HaveOccurred())
				Expect(unsubscribed).To(BeTrue())
			})

			It("does not add resubscriptions to the unsubscribes Repo", func() {
				updater.Execute(conn, []models.Preference{
					{
						ClientID: "dogs",
						KindID:   "barking",
						Email:    true,
					},
				}, false, "the-user")

				unsubscribed, err := unsubscribesRepo.Get(conn, "the-user", "dogs", "barking")
				Expect(err).NotTo(HaveOccurred())
				Expect(unsubscribed).To(BeFalse())
			})

			It("removes unsubscribes when they are resubscribed", func() {
				err := unsubscribesRepo.Set(conn, "my-user", "raptors", "door-open", true)
				Expect(err).NotTo(HaveOccurred())

				err = updater.Execute(conn, []models.Preference{
					{
						ClientID: "raptors",
						KindID:   "door-open",
						Email:    true,
					},
				}, false, "my-user")
				Expect(err).NotTo(HaveOccurred())

				unsubscribed, err := unsubscribesRepo.Get(conn, "my-user", "raptors", "door-open")
				Expect(err).NotTo(HaveOccurred())
				Expect(unsubscribed).To(BeFalse())
			})
		})

		Context("when unsubscribing from missing client", func() {
			var hungry models.Preference
			var boo models.Preference

			BeforeEach(func() {
				hungry = models.Preference{
					ClientID: "raptors",
					KindID:   "hungry",
				}

				boo = models.Preference{
					ClientID: "ghosts",
					KindID:   "boo",
				}

				kindsRepo.Create(conn, models.Kind{
					ID:       "hungry",
					ClientID: "raptors",
				})

				kindsRepo.Create(conn, models.Kind{
					ID:       "boo",
					ClientID: "missing-client",
				})

			})

			It("should return a MissingKindOrClientError", func() {
				err := updater.Execute(conn, []models.Preference{hungry, boo}, false, "the-user")

				Expect(err).To(Equal(services.MissingKindOrClientError("The kind 'boo' cannot be found for client 'ghosts'")))
			})
		})

		Context("when unsubscribing from a missing kind", func() {
			var hungry models.Preference
			var dead models.Preference

			BeforeEach(func() {
				hungry = models.Preference{
					ClientID: "raptors",
					KindID:   "hungry",
				}

				dead = models.Preference{
					ClientID: "raptors",
					KindID:   "dead",
				}

				kindsRepo.Create(conn, models.Kind{
					ID:       "hungry",
					ClientID: "raptors",
				})

			})

			It("should return a MissingKindOrClientError", func() {

				err := updater.Execute(conn, []models.Preference{hungry, dead}, false, "the-user")

				Expect(err).To(Equal(services.MissingKindOrClientError("The kind 'dead' cannot be found for client 'raptors'")))
			})
		})

		Context("when unsubscribing from a critical kind", func() {
			var hungry models.Preference
			var barking models.Preference

			BeforeEach(func() {
				hungry = models.Preference{
					ClientID: "raptors",
					KindID:   "hungry",
				}

				barking = models.Preference{
					ClientID: "dogs",
					KindID:   "barking",
				}

				kindsRepo.Create(conn, models.Kind{
					ClientID: "raptors",
					ID:       "hungry",
					Critical: true,
				})

				kindsRepo.Create(conn, models.Kind{
					ClientID: "dogs",
					ID:       "barking",
				})

			})

			It("should return a CriticalKindError", func() {
				err := updater.Execute(conn, []models.Preference{barking, hungry}, false, "the-user")

				Expect(err).To(Equal(services.CriticalKindError("The kind 'hungry' for the 'raptors' client is critical and cannot be unsubscribed from")))
			})
		})
	})
})
