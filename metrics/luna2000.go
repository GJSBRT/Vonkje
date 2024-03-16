package metrics

var luna2000Metrics = []Metric{
	{
		Namespace: "luna2000",
		Name: "running_status",
		Help: "The running status",
		Fields: []string{
			"inverter",
			"battery",
		},
	},
	{
		Namespace: "luna2000",
		Name: "charging_status",
		Help: "The charging status",
		Fields: []string{
			"inverter",
			"battery",
		},
	},
	{
		Namespace: "luna2000",
		Name: "bus_voltage",
		Help: "The bus voltage",
		Fields: []string{
			"inverter",
			"battery",
		},
	},
	{
		Namespace: "luna2000",
		Name: "battery_capacity",
		Help: "The battery capacity",
		Fields: []string{
			"inverter",
			"battery",
		},
	},
	{
		Namespace: "luna2000",
		Name: "total_charge",
		Help: "The total charge",
		Fields: []string{
			"inverter",
			"battery",
		},
	},
	{
		Namespace: "luna2000",
		Name: "total_discharge",
		Help: "The total discharge",
		Fields: []string{
			"inverter",
			"battery",
		},
	},
}
