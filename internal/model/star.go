package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Star struct {
	// TODO: Complete Star struct
	UID              int     `json:"uid"`
	Name             string  `json:"n"`
	PlayerUID        int     `json:"puid"`
	Visible          bool    `json:"v"`
	Strength         int     `json:"st"` // Number of ships on the star
	NaturalResources int     `json:"nr"`
	ResourceLevel    int     `json:"r"` // Resource level of the star including terraforming bonus
	EconomyLevel     int     `json:"e"`
	IndustryLevel    int     `json:"i"`
	ScienceLevel     int     `json:"s"`
	WarpGate         bool    `json:"ga"`
	X                float32 `json:"x"`
	Y                float32 `json:"y"`
	C                float64 `json:"c"` // Where ships/tick is not a whole number, the amount currently produced
}

func (s *Star) UnmarshalJSON(data []byte) error {
	type Alias Star
	aux := struct {
		Visible  string `json:"v"`
		WarpGate int    `json:"ga"`
		X        string `json:"x"`
		Y        string `json:"y"`
		*Alias
	}{
		Alias: (*Alias)(s),
	}

	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}

	x, err := strconv.ParseFloat(aux.X, 32)
	if err != nil {
		return fmt.Errorf("failed to parse 'x' as float: %w", err)
	}
	y, err := strconv.ParseFloat(aux.Y, 32)
	if err != nil {
		return fmt.Errorf("failed to parse 'y' as float: %w", err)
	}

	s.X = float32(x)
	s.Y = float32(y)
	s.Visible = aux.Visible == "1"
	s.WarpGate = aux.WarpGate == 1
	return nil
}
