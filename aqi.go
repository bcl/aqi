// Package aqi contains functions used to convert concentration value to AQI
//
// This is based on this document from the EPA:
//
//	https://nepis.epa.gov/Exe/ZyPDF.cgi/P100W5UG.PDF?Dockey=P100W5UG.PDF
package aqi

import (
	"fmt"
	"math"
	"strings"
)

// The AQI conversion is essentially a table of ranges with the concentration
// linerly scaled to map the the AQI output value
// The concentration input may be averaged over some timespan, eg. NowCast
//
// Which is described here:
// https://www.airnow.gov/faqs/how-nowcast-algorithm-used-report/

// breakpoint contains 1 line of the table, mapping a concentration range to an AQI range
// Floating point values are terrible for using in comparisons, if you are not careful you
// will end up with unexpected holes. eg. when the result is something like 1.0004
// So the Lo and Hi values are * 10 and truncated to an int
type breakpoint struct {
	Lo    int
	Hi    int
	AQILo int
	AQIHi int
	Name  string
}

// pollutant is a table of breakpoints
type pollutant struct {
	Breakpoints []breakpoint
	Places      float64 // Decimal places in table (input values are * 10**Places)
}

// AQI looks up the concentration in the table and converts it to an AQI value
// It returns the value and the descriptive name of the range it falls into
func (p pollutant) AQI(Cpf float64) (int, string, error) {
	// Input data in the table is scaled by this
	inScale := math.Pow(10, p.Places)
	Cp := int(Cpf * inScale)
	// Find the correct entry
	for i := range p.Breakpoints {
		bp := p.Breakpoints[i]
		if Cp >= bp.Lo && Cp <= bp.Hi {
			// Calculate the AQI value (a linear scaling of the aqi min/max
			scale := inScale * float64(bp.AQIHi-bp.AQILo) / float64(bp.Hi-bp.Lo)
			value := scale * float64(Cp-bp.Lo) / inScale
			return int(math.Round(float64(bp.AQILo) + value)), bp.Name, nil
		}
	}

	return 0, "Unknown", fmt.Errorf("Concentration of %0.1f not found", float64(Cp)/inScale)
}

var (
	aqi = map[string]pollutant{
		"PM2.5": {
			Breakpoints: []breakpoint{
				{0, 120, 0, 50, "Good"},
				{120, 354, 51, 100, "Moderate"},
				{355, 554, 101, 150, "Unhealthy for Sensitive Groups"},
				{555, 1504, 151, 200, "Unhealthy"},
				{1505, 2504, 201, 300, "Very Unhealthy"},
				{2505, 3504, 301, 400, "Hazardous"},
				{3505, 5004, 401, 500, "Hazardous"},
			},
			Places: 1,
		},
		"PM10": {
			Breakpoints: []breakpoint{
				{0, 54, 0, 50, "Good"},
				{55, 154, 51, 100, "Moderate"},
				{155, 254, 101, 150, "Unhealthy for Sensitive Groups"},
				{255, 354, 151, 200, "Unhealthy"},
				{355, 424, 201, 300, "Very Unhealthy"},
				{425, 504, 301, 400, "Hazardous"},
				{505, 604, 401, 500, "Hazardous"},
			},
			Places: 0,
		},
	}
)

// Lookup returns the conversion of a concentration to an AQI value
// and the name of the range it falls into
func Lookup(name string, Cp float64) (int, string, error) {
	p, ok := aqi[name]
	if !ok {
		var keys []string
		for k := range aqi {
			keys = append(keys, k)
		}
		supported := strings.Join(keys, ", ")
		return 0, "Unknown", fmt.Errorf("Pollutant named %s not found in table. Supported pollutants are: %s", name, supported)
	}

	return p.AQI(Cp)
}
