package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/testing/mocks"
	"github.com/cloudfoundry-incubator/notifications/v1/models"
	"github.com/cloudfoundry-incubator/notifications/v1/services"

	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
)

var _ = Describe("PreferenceUpdater", func() {
	Describe("Update", func() {
		var (
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
				updater.Update(conn, []models.Preference{}, true, "user-guid")
				Expect(fakeGlobalUnsubscribesRepo.SetCall.Receives.Unsubscribed).To(BeTrue())

				updater.Update(conn, []models.Preference{}, false, "user-guid")
				Expect(fakeGlobalUnsubscribesRepo.SetCall.Receives.Unsubscribed).To(BeFalse())
			})

			Context("when the global unsubscribe repo errors", func() {
				It("returns the error", func() {
					fakeGlobalUnsubscribesRepo.SetCall.Returns.Error = errors.New("global unsubscribe db error")

					err := updater.Update(conn, []models.Preference{}, true, "user-guid")
					Expect(err).To(MatchError(errors.New("global unsubscribe db error")))
				})
			})
		})

		Context("When unsubscribing from existing kinds of existing clients", func() {
			BeforeEach(func() {

				kindsRepo.FindCall.Returns.Kinds = []models.Kind{
					{
						ID:       "door-open",
						ClientID: "raptors",
					},
					{
						ID:       "barking",
						ClientID: "dogs",
					},
				}
			})

			It("Adds New Unsubscribes to the unsubscribes Repo", func() {
				updater.Update(conn, []models.Preference{
					{
						ClientID: "raptors",
						KindID:   "door-open",
						Email:    false,
					},
				}, false, "the-user")

				Expect(unsubscribesRepo.SetCall.Receives.Connection).To(Equal(conn))
				Expect(unsubscribesRepo.SetCall.Receives.UserID).To(Equal("the-user"))
				Expect(unsubscribesRepo.SetCall.Receives.ClientID).To(Equal("raptors"))
				Expect(unsubscribesRepo.SetCall.Receives.KindID).To(Equal("door-open"))
				Expect(unsubscribesRepo.SetCall.Receives.Unsubscribe).To(BeTrue())
			})

			It("does not add resubscriptions to the unsubscribes Repo", func() {
				updater.Update(conn, []models.Preference{
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

				err = updater.Update(conn, []models.Preference{
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
			BeforeEach(func() {
				kindsRepo.FindCall.Returns.Kinds = []models.Kind{
					{
						ID:       "hungry",
						ClientID: "raptors",
					},
					{
						ID:       "boo",
						ClientID: "missing-client",
					},
				}

			})

			It("should return a MissingKindOrClientError", func() {
				preferences := []models.Preference{
					{
						ClientID: "ghosts",
						KindID:   "boo",
					},
					{
						ClientID: "raptors",
						KindID:   "hungry",
					},
				}
				kindsRepo.FindCall.Returns.Error = errors.New("something bad happened")

				err := updater.Update(conn, preferences, false, "the-user")
				Expect(err).To(MatchError(services.MissingKindOrClientError{Err: errors.New("The kind 'boo' cannot be found for client 'ghosts'")}))
			})
		})

		Context("when unsubscribing from a missing kind", func() {
			BeforeEach(func() {
				kindsRepo.FindCall.Returns.Kinds = []models.Kind{
					{
						ID:       "hungry",
						ClientID: "raptors",
					},
				}
			})

			It("should return a MissingKindOrClientError", func() {
				preferences := []models.Preference{
					{
						ClientID: "raptors",
						KindID:   "dead",
					},
					{
						ClientID: "raptors",
						KindID:   "hungry",
					},
				}
				kindsRepo.FindCall.Returns.Error = errors.New("something bad happened")

				err := updater.Update(conn, preferences, false, "the-user")
				Expect(err).To(Equal(services.MissingKindOrClientError{Err: errors.New("The kind 'dead' cannot be found for client 'raptors'")}))
			})
		})

		Context("when unsubscribing from a critical kind", func() {
			BeforeEach(func() {
				kindsRepo.FindCall.Returns.Kinds = []models.Kind{
					{
						ClientID: "raptors",
						ID:       "hungry",
						Critical: true,
					},
					{
						ClientID: "dogs",
						ID:       "barking",
					},
				}
			})

			It("should return a CriticalKindError", func() {
				preferences := []models.Preference{
					{
						ClientID: "raptors",
						KindID:   "hungry",
					},
					{
						ClientID: "dogs",
						KindID:   "barking",
					},
				}

				err := updater.Update(conn, preferences, false, "the-user")
				Expect(err).To(Equal(services.CriticalKindError{Err: errors.New("The kind 'hungry' for the 'raptors' client is critical and cannot be unsubscribed from")}))
			})
		})
	})
})
