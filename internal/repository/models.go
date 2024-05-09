package repository

type GameRow struct {
	ID        int64
	Number    string
	APIKey    string
	PlayerUID int
	StartTime int64
	TickRate  int
}

type GameRowCreate struct {
	Number       string
	APIKey       string
	PlayerUID    int   `json:"player_uid"`
	StartTimeRaw int64 `json:"start_time"`
	TickRate     int   `json:"tick_rate"`
}

type SnapshotRow struct {
	ID   int64
	Path string
}
