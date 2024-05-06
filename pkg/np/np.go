package np

const apiBaseURL string = "https://np.ironhelmet.com/api"

type NeptunesPrideGame struct {
	Number string

	api api
}

func New(gameNumber, apiKey string) NeptunesPrideGame {
	return NeptunesPrideGame{
		Number: gameNumber,
		api:    NewAPI(gameNumber, apiKey),
	}
}

// GetCurrentSnapshot returns the JSON representation of the current game in bytes.
func (np NeptunesPrideGame) GetCurrentSnapshot() ([]byte, error) {
	return np.api.GetCurrentSnapshot()
}
