package metrics

var controlMetrics = []Metric{
	{
		Namespace: "control",
		Name: "action",
		Help: "The what action is currently being executed",
		Fields: []string{
			"action",
		},
	},
	{
		Namespace: "control",
		Name: "over_production",
		Help: "The a percentage of solar over production",
		Fields: []string{},
	},
}
