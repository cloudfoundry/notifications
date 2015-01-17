package support

import "encoding/json"

type MessagesService struct {
	client *Client
}

type GetResponse struct {
	HTTPStatus int
	Status     string   `json:"status"`
	Errors     []string `json:"errors"`
}

func (m MessagesService) Get(token, messageGUID string) (GetResponse, error) {
	var responseStruct GetResponse

	status, body, err := m.client.makeRequest("GET", m.client.server.StatusPath(messageGUID), nil, token)
	responseStruct.HTTPStatus = status
	if err != nil {
		return responseStruct, err
	}

	err = json.NewDecoder(body).Decode(&responseStruct)
	return responseStruct, err

}
