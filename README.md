# AQI calculation package

This is based on this document from the EPA:
    https://nepis.epa.gov/Exe/ZyPDF.cgi/P100W5UG.PDF?Dockey=P100W5UG.PDF

Note that additional work is needed to average the Particulate Matter (PM)
values before conversion, this package only implements the calcuation, not data
gathering or averaging.

The AQI conversion is essentially a table of ranges with the concentration
linerly scaled to map the the AQI output value
The concentration input may be averaged over some timespan, eg. NowCast

Which is described here:
    https://www.airnow.gov/faqs/how-nowcast-algorithm-used-report/


## Basic usage

Call the `Lookup` function with `PM2.5` or `PM10` and a `float64` PM
value. It will return the AQI value and the name of the range it falls
into.

```go
package main

import (
	"fmt"

	"github.com/bcl/aqi"
)

func main() {
	var pm2p5 float64
	var pm10 float64

	pm2p5 = 17.4
	pm10 = 3

	aqi_2p5, name_2p5, _ := aqi.Lookup("PM2.5", pm2p5)
	aqi_10, name_10, _ := aqi.Lookup("PM10", pm10)

	fmt.Printf("PM2.5 AQI = %d (%s)\n", aqi_2p5, name_2p5)
	fmt.Printf("PM10 AQI = %d (%s)\n", aqi_10, name_10)
}
```
