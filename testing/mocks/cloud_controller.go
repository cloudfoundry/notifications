package mocks

import "github.com/cloudfoundry-incubator/notifications/cf"

type CloudController struct {
	GetAuditorsByOrgGuidCall struct {
		Receives struct {
			OrgGUID string
			Token   string
		}
		Returns struct {
			Users []cf.CloudControllerUser
			Error error
		}
	}

	GetBillingManagersByOrgGuidCall struct {
		Receives struct {
			OrgGUID string
			Token   string
		}
		Returns struct {
			Users []cf.CloudControllerUser
			Error error
		}
	}

	GetManagersByOrgGuidCall struct {
		Receives struct {
			OrgGUID string
			Token   string
		}
		Returns struct {
			Users []cf.CloudControllerUser
			Error error
		}
	}

	GetUsersByOrgGuidCall struct {
		Receives struct {
			OrgGUID string
			Token   string
		}
		Returns struct {
			Users []cf.CloudControllerUser
			Error error
		}
	}

	GetUsersBySpaceGuidCall struct {
		Receives struct {
			SpaceGUID string
			Token     string
		}
		Returns struct {
			Users []cf.CloudControllerUser
			Error error
		}
	}

	LoadOrganizationCall struct {
		Receives struct {
			OrgGUID string
			Token   string
		}
		Returns struct {
			Organization cf.CloudControllerOrganization
			Error        error
		}
	}

	LoadSpaceCall struct {
		Receives struct {
			SpaceGUID string
			Token     string
		}
		Returns struct {
			Space cf.CloudControllerSpace
			Error error
		}
	}
}

func NewCloudController() *CloudController {
	return &CloudController{}
}

func (cc *CloudController) GetAuditorsByOrgGuid(orgGUID, token string) ([]cf.CloudControllerUser, error) {
	cc.GetAuditorsByOrgGuidCall.Receives.OrgGUID = orgGUID
	cc.GetAuditorsByOrgGuidCall.Receives.Token = token

	return cc.GetAuditorsByOrgGuidCall.Returns.Users, cc.GetAuditorsByOrgGuidCall.Returns.Error
}

func (cc *CloudController) GetBillingManagersByOrgGuid(orgGUID, token string) ([]cf.CloudControllerUser, error) {
	cc.GetBillingManagersByOrgGuidCall.Receives.OrgGUID = orgGUID
	cc.GetBillingManagersByOrgGuidCall.Receives.Token = token

	return cc.GetBillingManagersByOrgGuidCall.Returns.Users, cc.GetBillingManagersByOrgGuidCall.Returns.Error
}

func (cc *CloudController) GetManagersByOrgGuid(orgGUID, token string) ([]cf.CloudControllerUser, error) {
	cc.GetManagersByOrgGuidCall.Receives.OrgGUID = orgGUID
	cc.GetManagersByOrgGuidCall.Receives.Token = token

	return cc.GetManagersByOrgGuidCall.Returns.Users, cc.GetManagersByOrgGuidCall.Returns.Error
}

func (cc *CloudController) GetUsersByOrgGuid(orgGUID, token string) ([]cf.CloudControllerUser, error) {
	cc.GetUsersByOrgGuidCall.Receives.OrgGUID = orgGUID
	cc.GetUsersByOrgGuidCall.Receives.Token = token

	return cc.GetUsersByOrgGuidCall.Returns.Users, cc.GetUsersByOrgGuidCall.Returns.Error
}

func (cc *CloudController) GetUsersBySpaceGuid(spaceGUID, token string) ([]cf.CloudControllerUser, error) {
	cc.GetUsersBySpaceGuidCall.Receives.SpaceGUID = spaceGUID
	cc.GetUsersBySpaceGuidCall.Receives.Token = token

	return cc.GetUsersBySpaceGuidCall.Returns.Users, cc.GetUsersBySpaceGuidCall.Returns.Error
}

func (cc *CloudController) LoadOrganization(orgGUID, token string) (cf.CloudControllerOrganization, error) {
	cc.LoadOrganizationCall.Receives.OrgGUID = orgGUID
	cc.LoadOrganizationCall.Receives.Token = token

	return cc.LoadOrganizationCall.Returns.Organization, cc.LoadOrganizationCall.Returns.Error
}

func (cc *CloudController) LoadSpace(spaceGUID, token string) (cf.CloudControllerSpace, error) {
	cc.LoadSpaceCall.Receives.SpaceGUID = spaceGUID
	cc.LoadSpaceCall.Receives.Token = token

	return cc.LoadSpaceCall.Returns.Space, cc.LoadSpaceCall.Returns.Error
}
