package model_test

import (
	"nphud/pkg/np/model"
	"testing"
	"time"
)

func TestScanningData_GetNextProductionTime(t *testing.T) {
	scanningData := model.ScanningData{
		TickRate:          60,
		ProductionRate:    24,
		ProductionCounter: 4,
	}

	tolerance := 2 * time.Second
	expectedTime := time.Now().Add(time.Duration(1200) * time.Minute)
	actualTime := scanningData.GetNextProductionTime()

	timeDiff := expectedTime.Sub(actualTime)
	if timeDiff < -tolerance || timeDiff > tolerance {
		t.Errorf("Expected time to be within %v of %v, but got %v", tolerance, expectedTime, actualTime)
	}
}
