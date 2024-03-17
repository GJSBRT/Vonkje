package metrics

var powerMeterMetrics = []Metric{
	{
		Namespace: "power_meter",
		Name: "status",
		Help: "The status of the power meter",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "power_meter",
		Name: "phase_voltage",
		Help: "The phase voltage",
		Fields: []string{
			"inverter",
			"phase",
		},
	},
	{
		Namespace: "power_meter",
		Name: "phase_current",
		Help: "The phase current",
		Fields: []string{
			"inverter",
			"phase",
		},
	},
	{
		Namespace: "power_meter",
		Name: "active_power",
		Help: "The active power",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "power_meter",
		Name: "reactive_power",
		Help: "The reactive power",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "power_meter",
		Name: "power_factor",
		Help: "The power factor",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "power_meter",
		Name: "frequency",
		Help: "The frequency",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "power_meter",
		Name: "positive_active_electricity",
		Help: "The positive active electricity",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "power_meter",
		Name: "reverse_active_power",
		Help: "The reverse active power",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "power_meter",
		Name: "accumulated_reactive_power",
		Help: "The accumulated reactive power",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "power_meter",
		Name: "line_voltage",
		Help: "The line voltage",
		Fields: []string{
			"inverter",
			"line",
		},
	},
	{
		Namespace: "power_meter",
		Name: "phase_active_power",
		Help: "The phase active power",
		Fields: []string{
			"inverter",
			"phase",
		},
	},
}
