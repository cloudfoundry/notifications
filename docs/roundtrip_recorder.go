package docs

import (
	"fmt"
	"net/http"
)

type RoundTripRecorder struct {
	RoundTrips map[string]RoundTrip
}

func NewRoundTripRecorder() *RoundTripRecorder {
	return &RoundTripRecorder{
		RoundTrips: make(map[string]RoundTrip),
	}
}

func (r *RoundTripRecorder) Record(key string, request *http.Request, response *http.Response) error {
	if _, present := r.RoundTrips[key]; present {
		return fmt.Errorf("new roundtrip %q conflicts with existing roundtrip", key)
	}

	r.RoundTrips[key] = RoundTrip{
		Request:  request,
		Response: response,
	}

	return nil
}
