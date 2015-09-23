package senders

import "github.com/cloudfoundry-incubator/notifications/v2/collections"

type SendersListResponse struct {
	Senders []SenderResponse         `json:"senders"`
	Links   SendersListResponseLinks `json:"_links"`
}

type SendersListResponseLinks struct {
	Self Link `json:"self"`
}

func NewSendersListResponse(senderList []collections.Sender) SendersListResponse {
	senderResponseList := []SenderResponse{}

	for _, sender := range senderList {
		senderResponseList = append(senderResponseList, NewSenderResponse(sender))
	}

	return SendersListResponse{
		Senders: senderResponseList,
		Links: SendersListResponseLinks{
			Self: Link{"/senders"},
		},
	}
}
