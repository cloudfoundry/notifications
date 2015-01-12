package support

type PreferencesService struct {
	client *Client
}

func (service PreferencesService) User(userGUID string) UserPreferencesService {
	return UserPreferencesService{
		client:   service.client,
		userGUID: userGUID,
	}
}
