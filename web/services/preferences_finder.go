package services

import "github.com/cloudfoundry-incubator/notifications/models"

type PreferencesFinder struct {
    preferencesRepo        models.PreferencesRepoInterface
    globalUnsubscribesRepo models.GlobalUnsubscribesRepoInterface
}

type PreferencesFinderInterface interface {
    Find(string) (PreferencesBuilder, error)
}

func NewPreferencesFinder(preferencesRepo models.PreferencesRepoInterface, globalUnsubscribesRepo models.GlobalUnsubscribesRepoInterface) *PreferencesFinder {
    return &PreferencesFinder{
        preferencesRepo:        preferencesRepo,
        globalUnsubscribesRepo: globalUnsubscribesRepo,
    }
}

func (finder PreferencesFinder) Find(userGUID string) (PreferencesBuilder, error) {
    conn := models.Database().Connection()
    builder := NewPreferencesBuilder()

    globallyUnsubscribed, err := finder.globalUnsubscribesRepo.Get(conn, userGUID)
    if err != nil {
        return builder, err
    }

    preferences, err := finder.preferencesRepo.FindNonCriticalPreferences(conn, userGUID)
    if err != nil {
        return builder, err
    }

    builder.GlobalUnsubscribe = globallyUnsubscribed
    for _, preference := range preferences {
        builder.Add(preference)
    }

    return builder, nil
}
