package notify

import "regexp"

var kindIDFormat = regexp.MustCompile(`^[0-9a-zA-Z_\-.]+$`)

type EmailValidator struct{}

func (validator EmailValidator) Validate(notify *NotifyParams) bool {
	notify.Errors = []string{}

	if notify.To == "" {
		notify.Errors = append(notify.Errors, `"to" is a required field`)
	}

	if notify.To == InvalidEmail {
		notify.Errors = append(notify.Errors, `"to" is improperly formatted`)
	}

	if missingTextOrHTMLFields(notify) {
		notify.Errors = append(notify.Errors, `"text" or "html" fields must be supplied`)
	}

	return len(notify.Errors) == 0
}

type GUIDValidator struct{}

func (validator GUIDValidator) Validate(notify *NotifyParams) bool {
	notify.Errors = []string{}

	validator.checkKindIDField(notify)

	if missingTextOrHTMLFields(notify) {
		notify.Errors = append(notify.Errors, `"text" or "html" fields must be supplied`)
	}

	if validator.invalidRoleField(notify.Role) {
		notify.Errors = append(notify.Errors, `"role" must be "OrgManager", "OrgAuditor", "BillingManager" or unset`)
	}

	return len(notify.Errors) == 0
}

func missingTextOrHTMLFields(notify *NotifyParams) bool {
	return notify.Text == "" && notify.ParsedHTML.BodyContent == ""
}

func (validator GUIDValidator) invalidRoleField(roleName string) bool {
	if roleName == "" {
		return false
	}

	for _, role := range validOrganizationRoles {
		if roleName == role {
			return false
		}
	}
	return true
}

func (validator GUIDValidator) checkKindIDField(notify *NotifyParams) {
	if notify.KindID == "" {
		notify.Errors = append(notify.Errors, `"kind_id" is a required field`)
	} else {
		if !kindIDFormat.MatchString(notify.KindID) {
			notify.Errors = append(notify.Errors, `"kind_id" is improperly formatted`)
		}
	}
}
