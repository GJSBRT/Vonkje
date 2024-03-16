package metrics

var sun2000Metrics = []Metric{
	{
		Namespace: "sun2000",
		Name: "pv_voltage",
		Help: "The total amount of voltage",
		Fields: []string{
			"inverter",
			"string",
		},
	},
	{
		Namespace: "sun2000",
		Name: "pv_current",
		Help: "The total amount of current",
		Fields: []string{
			"inverter",
			"string",
		},
	},
	{
		Namespace: "sun2000",
		Name: "phase_voltage",
		Help: "The total amount of voltage",
		Fields: []string{
			"inverter",
			"phase",
		},
	},
	{
		Namespace: "sun2000",
		Name: "phase_current",
		Help: "The total amount of current",
		Fields: []string{
			"inverter",
			"phase",
		},
	},
	{
		Namespace: "sun2000",
		Name: "input_power",
		Help: "The total amount of input power",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "sun2000",
		Name: "active_power",
		Help: "The total amount of active power",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "sun2000",
		Name: "state",
		Help: "The state of the sun2000 inverter",
		Fields: []string{
			"inverter",
			"state",
		},
	},
	{
		Namespace: "sun2000",
		Name: "device_status",
		Help: "The device status of the inverter",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "sun2000",
		Name: "reactive_power",
		Help: "The total amount of reactive power",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "sun2000",
		Name: "power_factor",
		Help: "The power factor",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "sun2000",
		Name: "grid_frequency",
		Help: "The grid frequency",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "sun2000",
		Name: "inverter_efficiency",
		Help: "The inverter efficiency",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "sun2000",
		Name: "cabinet_temperature",
		Help: "The cabinet temperature",
		Fields: []string{
			"inverter",
		},
	},
	{
		Namespace: "sun2000",
		Name: "isulation_resistance",
		Help: "The isulation resistance",
		Fields: []string{
			"inverter",
		},
	},
}
