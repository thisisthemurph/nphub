package np

const apiBaseURL string = "https://np.ironhelmet.com/api"

type NeptunesPrideGame struct {
	Number string
	APIKey string

	api api
}

func New(gameNumber, apiKey string) NeptunesPrideGame {
	return NeptunesPrideGame{
		Number: gameNumber,
		APIKey: apiKey,
		api:    newAPI(gameNumber, apiKey),
	}
}

// TakeSnapshot returns the JSON representation of the current game in bytes.
func (np NeptunesPrideGame) TakeSnapshot() ([]byte, error) {
	return np.api.GetCurrentSnapshot()
}
