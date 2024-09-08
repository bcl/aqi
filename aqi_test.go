package aqi

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPM2p5(t *testing.T) {
	// Test the endpoints of each section of the PM2.5 AQI table and some intermediate values
	data := []struct {
		in  float64
		out int
	}{
		{0, 0}, {6, 25}, {12, 50}, {12.1, 51}, {25, 78}, {35.4, 100}, {35.5, 101}, {42, 117}, {55.4, 150},
		{55.5, 151}, {66, 156}, {150.4, 200}, {150.5, 201}, {200, 250}, {223, 273}, {250.4, 300},
		{250.5, 301}, {301, 351}, {350.4, 400}, {350.5, 401}, {413, 442}, {478, 485}, {500.4, 500},
	}

	for i := range data {
		t.Logf("PM2.5 %0.1f ug/m3", data[i].in)
		v, _, err := Lookup("PM2.5", data[i].in)
		assert.Nil(t, err)
		if err != nil {
			continue
		}
		assert.Equal(t, data[i].out, v)
	}
}

func TestPM10(t *testing.T) {
	// Test the endpoints of each section of the PM10 AQI table
	data := []struct {
		in  float64
		out int
	}{
		{0, 0}, {27, 25}, {54, 50}, {55, 51}, {98, 72}, {154, 100}, {155, 101}, {203, 125}, {254, 150},
		{255, 151}, {313, 180}, {354, 200}, {355, 201}, {383, 241}, {424, 300},
		{425, 301}, {454, 337}, {504, 400}, {505, 401}, {589, 485}, {604, 500},
	}

	for i := range data {
		t.Logf("PM10 %0.0f ug/m3", data[i].in)
		v, _, err := Lookup("PM10", data[i].in)
		assert.Nil(t, err)
		if err != nil {
			continue
		}
		assert.Equal(t, data[i].out, v)
	}
}

func TestUnknownPollutant(t *testing.T) {
	_, _, err := Lookup("PCB", 11.1)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "not found in table")
}

func TestUnknownPM2p5Concentration(t *testing.T) {
	_, _, err := Lookup("PM2.5", 1000)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestUnknownPM10Concentration(t *testing.T) {
	_, _, err := Lookup("PM10", 1000)
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "not found")
}
