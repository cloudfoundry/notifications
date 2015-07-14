package domain

import "github.com/nu7hatch/gouuid"

func generateID() string {
	guid, err := uuid.NewV4()
	if err != nil {
		panic(err)
	}

	return guid.String()
}
