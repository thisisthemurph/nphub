package np

import "nphud/pkg/np/model"

type APIResponse struct {
	ScanningData model.ScanningData `json:"scanning_data"`
}

type APIErrorResponse struct {
	ErrorMessage string `json:"error"`
}

func (e APIErrorResponse) Error() string {
	return e.ErrorMessage
}
