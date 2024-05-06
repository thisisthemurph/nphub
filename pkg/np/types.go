package np

type APIErrorResponse struct {
	ErrorMessage string `json:"error"`
}

func (e APIErrorResponse) Error() string {
	return e.ErrorMessage
}

type APIResponse struct {
	ScanningData ScanningData `json:"scanning_data"`
}

type ScanningData struct {
	GameOver int `json:"game_over"`
}
