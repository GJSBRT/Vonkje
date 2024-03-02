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
			"connection",
			"string",
		},
	)
	pvCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "pv_current",
			Help: "The total amount of current",
		},
		[]string{
			"connection",
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
			"connection",
			"phase",
		},
	)
	phaseCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "phase_current",
			Help: "The total amount of current",
		},
		[]string{
			"connection",
			"phase",
		},
	)

	// other
	activePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "active_power",
			Help: "The total amount of active power",
		},
		[]string{
			"connection",
		},
	)
	reactivePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "reactive_power",
			Help: "The total amount of reactive power",
		},
		[]string{
			"connection",
		},
	)
	powerFactorGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "power_factor",
			Help: "The power factor",
		},
		[]string{
			"connection",
		},
	)
	gridFrequencyGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "grid_frequency",
			Help: "The grid frequency",
		},
		[]string{
			"connection",
		},
	)
	inverterEfficiencyGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "connection_efficiency",
			Help: "The connection efficiency",
		},
		[]string{
			"connection",
		},
	)
	cabinetTemperatureGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "cabinet_temperature",
			Help: "The cabinet temperature",
		},
		[]string{
			"connection",
		},
	)
	isulationResistanceGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "isulation_resistance",
			Help: "The isulation resistance",
		},
		[]string{
			"connection",
		},
	)

	// Battery
	runningStatusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "running_status",
			Help: "The running status",
		},
		[]string{
			"connection",
			"battery",
		},
	)
	chargingStatusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "charging_status",
			Help: "The charging status",
		},
		[]string{
			"connection",
			"battery",
		},
	)
	busVoltageGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "bus_voltage",
			Help: "The bus voltage",
		},
		[]string{
			"connection",
			"battery",
		},
	)
	batteryCapacityGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "battery_capacity",
			Help: "The battery capacity",
		},
		[]string{
			"connection",
			"battery",
		},
	)
	totalChargeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "total_charge",
			Help: "The total charge",
		},
		[]string{
			"connection",
			"battery",
		},
	)
	totalDischargeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "total_discharge",
			Help: "The total discharge",
		},
		[]string{
			"connection",
			"battery",
		},
	)
)