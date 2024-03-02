package modbus

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	// Strings
	pvVoltage = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "pv_voltage",
			Help: "The total amount of voltage",
		},
		[]string{
			"string",
		},
	)
	pvCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "pv_current",
			Help: "The total amount of current",
		},
		[]string{
			"string",
		},
	)

	// Phases
	phaseVoltage = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "phase_voltage",
			Help: "The total amount of voltage",
		},
		[]string{
			"phase",
		},
	)
	phaseCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "phase_current",
			Help: "The total amount of current",
		},
		[]string{
			"phase",
		},
	)

	// other
	activePowerGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "active_power",
			Help: "The total amount of active power",
		},
	)
	reactivePowerGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "reactive_power",
			Help: "The total amount of reactive power",
		},
	)
	powerFactorGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "power_factor",
			Help: "The power factor",
		},
	)
	gridFrequencyGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "grid_frequency",
			Help: "The grid frequency",
		},
	)
	inverterEfficiencyGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "inverter_efficiency",
			Help: "The inverter efficiency",
		},
	)
	cabinetTemperatureGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "cabinet_temperature",
			Help: "The cabinet temperature",
		},
	)
	isulationResistanceGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "isulation_resistance",
			Help: "The isulation resistance",
		},
	)

	// Battery
	runningStatusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "running_status",
			Help: "The running status",
		},
		[]string{
			"battery",
		},
	)
	chargingStatusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "charging_status",
			Help: "The charging status",
		},
		[]string{
			"battery",
		},
	)
	busVoltageGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "bus_voltage",
			Help: "The bus voltage",
		},
		[]string{
			"battery",
		},
	)
	batteryCapacityGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "battery_capacity",
			Help: "The battery capacity",
		},
		[]string{
			"battery",
		},
	)
	totalChargeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "total_charge",
			Help: "The total charge",
		},
		[]string{
			"battery",
		},
	)
	totalDischargeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "total_discharge",
			Help: "The total discharge",
		},
		[]string{
			"battery",
		},
	)
)