package model

import (
	"encoding/json"
	"fmt"
	"strconv"
)

type Fleet struct {
	UID       int            `json:"uid"`
	Looping   bool           `json:"l"`
	Orders    []CarrierOrder `json:"o"`
	Name      string         `json:"n"`
	PlayerUID int            `json:"puid"`
	Warping   bool           `json:"w"`
	X         float32        `json:"x"`
	Y         float32        `json:"y"`
	Strength  int            `json:"st"` // Number of ships
	LastX     float32        `json:"lx"`
	LastY     float32        `json:"ly"`
}

func (f *Fleet) UnmarshalJSON(data []byte) error {
	type Alias Fleet
	aux := struct {
		Looping int     `json:"l"`
		Orders  [][]int `json:"o"`
		Warping int     `json:"w"`
		X       string  `json:"x"`
		Y       string  `json:"y"`
		LastX   string  `json:"lx"`
		LastY   string  `json:"ly"`
		*Alias
	}{
		Alias: (*Alias)(f),
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
	lx, err := strconv.ParseFloat(aux.LastX, 32)
	if err != nil {
		return fmt.Errorf("failed to parse 'x' as float: %w", err)
	}
	ly, err := strconv.ParseFloat(aux.LastY, 32)
	if err != nil {
		return fmt.Errorf("failed to parse 'y' as float: %w", err)
	}

	// Process carrier orders
	for _, o := range aux.Orders {
		delay := o[0]
		uid := o[1]
		orderTypeID := o[2]
		numShips := o[3]

		carrierOrder := CarrierOrder{
			Delay:         delay,
			UID:           uid,
			OrderTypeID:   CarrierOrderType(orderTypeID),
			NumberOfShips: numShips,
		}

		f.Orders = append(f.Orders, carrierOrder)
	}

	f.Looping = aux.Looping == 1
	f.Warping = aux.Warping == 1
	f.X = float32(x)
	f.Y = float32(y)
	f.LastX = float32(lx)
	f.LastY = float32(ly)
	return nil
}
