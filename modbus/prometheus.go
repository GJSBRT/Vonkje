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
	sun2000InputPowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "input_power",
			Help: "The total amount of input power",
		},
		[]string{
			"connection",
		},
	)
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


	// Power meter
	powerMeterStatusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "status",
			Help: "The status of the power meter",
		},
		[]string{
			"connection",
		},
	)
	powerMeterPhaseVoltageGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "phase_voltage",
			Help: "The phase voltage",
		},
		[]string{
			"connection",
			"phase",
		},
	)
	powerMeterPhaseCurrentGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "phase_current",
			Help: "The phase current",
		},
		[]string{
			"connection",
			"phase",
		},
	)
	powerMeterActivePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "active_power",
			Help: "The active power",
		},
		[]string{
			"connection",
		},
	)
	powerMeterReactivePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "reactive_power",
			Help: "The reactive power",
		},
		[]string{
			"connection",
		},
	)
	powerMeterPowerFactorGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "power_factor",
			Help: "The power factor",
		},
		[]string{
			"connection",
		},
	)
	powerMeterFrequencyGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "frequency",
			Help: "The frequency",
		},
		[]string{
			"connection",
		},
	)
	powerMeterPositiveActiveElectricityGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "positive_active_electricity",
			Help: "The positive active electricity",
		},
		[]string{
			"connection",
		},
	)
	powerMeterReverseActivePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "reverse_active_power",
			Help: "The reverse active power",
		},
		[]string{
			"connection",
		},
	)
	powerMeterAccumulatedReactivePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "accumulated_reactive_power",
			Help: "The accumulated reactive power",
		},
		[]string{
			"connection",
		},
	)
	powerMeterLineVoltageGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "line_voltage",
			Help: "The line voltage",
		},
		[]string{
			"connection",
			"line",
		},
	)
	powerMeterPhaseActivePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "phase_active_power",
			Help: "The phase active power",
		},
		[]string{
			"connection",
			"phase",
		},
	)
	powerMeterModelResultGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "model_result",
			Help: "The power meter model result",
		},
		[]string{
			"connection",
		},
	)
)