package services

type PreferencesFinder struct {
	preferencesRepo        PreferencesRepo
	globalUnsubscribesRepo GlobalUnsubscribesRepo
}

func NewPreferencesFinder(preferencesRepo PreferencesRepo, globalUnsubscribesRepo GlobalUnsubscribesRepo) *PreferencesFinder {
	return &PreferencesFinder{
		preferencesRepo:        preferencesRepo,
		globalUnsubscribesRepo: globalUnsubscribesRepo,
	}
}

func (finder PreferencesFinder) Find(database DatabaseInterface, userGUID string) (PreferencesBuilder, error) {
	conn := database.Connection()
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
