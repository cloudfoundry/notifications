package services_test

import (
	"errors"

	"github.com/cloudfoundry-incubator/notifications/fakes"
	"github.com/cloudfoundry-incubator/notifications/models"
	"github.com/cloudfoundry-incubator/notifications/web/services"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("PreferenceUpdater", func() {
	Describe("Execute", func() {
		var doorOpen models.Unsubscribe
		var barking models.Unsubscribe
		var unsubscribesRepo *fakes.UnsubscribesRepo
		var kindsRepo *fakes.KindsRepo
		var fakeGlobalUnsubscribesRepo *fakes.GlobalUnsubscribesRepo
		var conn *fakes.DBConn
		var updater services.PreferenceUpdater

		BeforeEach(func() {
			conn = fakes.NewDBConn()
			unsubscribesRepo = fakes.NewUnsubscribesRepo()
			kindsRepo = fakes.NewKindsRepo()
			fakeGlobalUnsubscribesRepo = fakes.NewGlobalUnsubscribesRepo()
			updater = services.NewPreferenceUpdater(fakeGlobalUnsubscribesRepo, unsubscribesRepo, kindsRepo)
		})

		Context("when globally unsubscribing", func() {
			It("inserts a record into the global unsubscribes repo", func() {
				userGUID := "user-guid"
				updater.Execute(conn, []models.Preference{}, true, userGUID)

				globallyUnsubscribed, err := fakeGlobalUnsubscribesRepo.Get(conn, userGUID)
				if err != nil {
					panic(err)
				}

				Expect(globallyUnsubscribed).To(BeTrue())

				updater.Execute(conn, []models.Preference{}, false, userGUID)

				globallyUnsubscribed, err = fakeGlobalUnsubscribesRepo.Get(conn, userGUID)
				if err != nil {
					panic(err)
				}

				Expect(globallyUnsubscribed).To(BeFalse())
			})

			Context("when the global unsubscribe repo errors", func() {
				It("returns the error", func() {
					fakeGlobalUnsubscribesRepo.SetError = errors.New("global unsubscribe db error")

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

				Expect(len(unsubscribesRepo.Unsubscribes)).To(Equal(2))
				Expect(unsubscribesRepo.Unsubscribes).To(ContainElement(doorOpen))
				Expect(unsubscribesRepo.Unsubscribes).To(ContainElement(barking))
			})

			It("does not insert duplicate unsubscribes", func() {
				_, err := unsubscribesRepo.Create(conn, models.Unsubscribe{
					UserID:   "my-user",
					ClientID: "raptors",
					KindID:   "door-open",
				})
				if err != nil {
					panic(err)
				}
				Expect(len(unsubscribesRepo.Unsubscribes)).To(Equal(1))

				err = updater.Execute(conn, []models.Preference{
					{
						ClientID: "raptors",
						KindID:   "door-open",
						Email:    false,
					},
				}, false, "my-user")

				Expect(err).To(BeNil())
				Expect(len(unsubscribesRepo.Unsubscribes)).To(Equal(1))
			})

			It("Does not add resubscriptions to the unsubscribes Repo", func() {
				updater.Execute(conn, []models.Preference{
					{
						ClientID: "dogs",
						KindID:   "barking",
						Email:    true,
					},
				}, false, "the-user")

				Expect(len(unsubscribesRepo.Unsubscribes)).To(Equal(0))
			})

			It("removes unsubscribes when they are resubscribed", func() {
				_, err := unsubscribesRepo.Create(conn, models.Unsubscribe{
					UserID:   "my-user",
					ClientID: "raptors",
					KindID:   "door-open",
				})
				if err != nil {
					panic(err)
				}
				Expect(len(unsubscribesRepo.Unsubscribes)).To(Equal(1))

				err = updater.Execute(conn, []models.Preference{
					{
						ClientID: "raptors",
						KindID:   "door-open",
						Email:    true,
					},
				}, false, "my-user")

				Expect(err).To(BeNil())
				Expect(len(unsubscribesRepo.Unsubscribes)).To(Equal(0))
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
