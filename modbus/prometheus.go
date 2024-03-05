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
			"inverter",
			"string",
		},
	)
	pvCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "pv_current",
			Help: "The total amount of current",
		},
		[]string{
			"inverter",
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
			"inverter",
			"phase",
		},
	)
	phaseCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "phase_current",
			Help: "The total amount of current",
		},
		[]string{
			"inverter",
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
			"inverter",
		},
	)
	activePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "active_power",
			Help: "The total amount of active power",
		},
		[]string{
			"inverter",
		},
	)
	sun2000StateGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "state",
			Help: "The state of the sun2000 inverter",
		},
		[]string{
			"inverter",
			"state",
		},
	)
	sun2000DeviceStatusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "device_status",
			Help: "The device status of the inverter",
		},
		[]string{
			"inverter",
		},
	)
	reactivePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "reactive_power",
			Help: "The total amount of reactive power",
		},
		[]string{
			"inverter",
		},
	)
	powerFactorGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "power_factor",
			Help: "The power factor",
		},
		[]string{
			"inverter",
		},
	)
	gridFrequencyGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "grid_frequency",
			Help: "The grid frequency",
		},
		[]string{
			"inverter",
		},
	)
	inverterEfficiencyGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "connection_efficiency",
			Help: "The connection efficiency",
		},
		[]string{
			"inverter",
		},
	)
	cabinetTemperatureGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "cabinet_temperature",
			Help: "The cabinet temperature",
		},
		[]string{
			"inverter",
		},
	)
	isulationResistanceGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "isulation_resistance",
			Help: "The isulation resistance",
		},
		[]string{
			"inverter",
		},
	)

	// Battery
	runningStatusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "running_status",
			Help: "The running status",
		},
		[]string{
			"inverter",
			"battery",
		},
	)
	chargingStatusGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "charging_status",
			Help: "The charging status",
		},
		[]string{
			"inverter",
			"battery",
		},
	)
	busVoltageGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "bus_voltage",
			Help: "The bus voltage",
		},
		[]string{
			"inverter",
			"battery",
		},
	)
	batteryCapacityGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "battery_capacity",
			Help: "The battery capacity",
		},
		[]string{
			"inverter",
			"battery",
		},
	)
	totalChargeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "total_charge",
			Help: "The total charge",
		},
		[]string{
			"inverter",
			"battery",
		},
	)
	totalDischargeGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "luna2000",
			Name: "total_discharge",
			Help: "The total discharge",
		},
		[]string{
			"inverter",
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
			"inverter",
		},
	)
	powerMeterPhaseVoltageGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "phase_voltage",
			Help: "The phase voltage",
		},
		[]string{
			"inverter",
			"phase",
		},
	)
	powerMeterPhaseCurrentGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "phase_current",
			Help: "The phase current",
		},
		[]string{
			"inverter",
			"phase",
		},
	)
	powerMeterActivePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "active_power",
			Help: "The active power",
		},
		[]string{
			"inverter",
		},
	)
	powerMeterReactivePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "reactive_power",
			Help: "The reactive power",
		},
		[]string{
			"inverter",
		},
	)
	powerMeterPowerFactorGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "power_factor",
			Help: "The power factor",
		},
		[]string{
			"inverter",
		},
	)
	powerMeterFrequencyGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "frequency",
			Help: "The frequency",
		},
		[]string{
			"inverter",
		},
	)
	powerMeterPositiveActiveElectricityGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "positive_active_electricity",
			Help: "The positive active electricity",
		},
		[]string{
			"inverter",
		},
	)
	powerMeterReverseActivePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "reverse_active_power",
			Help: "The reverse active power",
		},
		[]string{
			"inverter",
		},
	)
	powerMeterAccumulatedReactivePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "accumulated_reactive_power",
			Help: "The accumulated reactive power",
		},
		[]string{
			"inverter",
		},
	)
	powerMeterLineVoltageGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "line_voltage",
			Help: "The line voltage",
		},
		[]string{
			"inverter",
			"line",
		},
	)
	powerMeterPhaseActivePowerGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "phase_active_power",
			Help: "The phase active power",
		},
		[]string{
			"inverter",
			"phase",
		},
	)
	powerMeterModelResultGauge = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "power_meter",
			Name: "model_result",
			Help: "The power meter model result",
		},
		[]string{
			"inverter",
		},
	)
)