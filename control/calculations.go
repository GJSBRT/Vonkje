package control

import (
	"gijs.eu/vonkje/metrics"
)

// calculateHomeLoad calculates the home load based on the active power of the inverter and power meter.
// If the sun is out, the home load is calculated based on the active power of the inverter.
// Many hours have been spent on trying to find a calculation for this any many more will be spent :/.
func calculateHomeLoad() (float64, error) {
	var homeLoad float64

	inverterInputPower, err := metrics.GetMetricLastEntrySum("sun2000", "active_power")
	if err != nil {
		return 0, err
	}
	inverterInputPower = inverterInputPower * 1000

	powerMeterActivePower, err := metrics.GetMetricLastEntrySum("power_meter", "active_power")
	if err != nil {
		return 0, err
	}

	if inverterInputPower > powerMeterActivePower {
		homeLoad = inverterInputPower - powerMeterActivePower
	} else {
		homeLoad = powerMeterActivePower - inverterInputPower
	}

	return homeLoad, nil
}
