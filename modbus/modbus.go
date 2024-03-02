package modbus

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/simonvetter/modbus"
)

type Config struct {
	IP string `mapstructure:"ip"`
	Port uint `mapstructure:"port"`
	Baudrate uint `mapstructure:"baudrate"`
	DataBits uint `mapstructure:"data-bits"`
	StopBits uint `mapstructure:"stop-bits"`
	Timeout uint `mapstructure:"timeout"`
	ReadMetricsInterval uint `mapstructure:"read-metrics-interval"`
}

type Modbus struct {
	config Config
	errChannel chan error
	ctx context.Context
	logger *logrus.Logger
	client *modbus.ModbusClient
}

func New(
	config Config, 
	errChannel chan error,
	ctx context.Context,
	logger *logrus.Logger,
) (*Modbus, error) {
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
		URL: fmt.Sprintf("rtuovertcp://%s:%d", config.IP, config.Port),
		Speed: config.Baudrate,
		DataBits: config.DataBits,
		StopBits: config.StopBits,
		Parity: modbus.PARITY_NONE,
		Timeout: time.Duration(config.Timeout) * time.Second,
	})
	if err != nil {
		return nil, err
	}

	err = client.Open()
	if err != nil {
		return nil, err
	}

	return &Modbus{
		config: config,
		errChannel: errChannel,
		ctx: ctx,
		logger: logger,
		client: client,
	}, nil
}

func (m *Modbus) Close() {
	m.client.Close()
}

func (m *Modbus) Start() {
	m.logger.Info("Starting modbus metrics collector")

	ticker := time.NewTicker(time.Duration(m.config.ReadMetricsInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			m.logger.Info("Stopping modbus metrics collector")
			return
		case <-ticker.C:
			m.updateMetrics()
		}
	}
}

func (m *Modbus) updateMetrics() {
	// string 1
	pv1Voltage, err := m.client.ReadRegister(MODBUS_INVERTER_PV1_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvVoltage.With(prometheus.Labels{"string": "1"}).Set(float64(pv1Voltage) / 10)

	pv1Current, err := m.client.ReadRegister(MODBUS_INVERTER_PV1_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvCurrent.With(prometheus.Labels{"string": "1"}).Set(float64(pv1Current) / 100)

	// string 2
	pv2Voltage, err := m.client.ReadRegister(MODBUS_INVERTER_PV2_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvVoltage.With(prometheus.Labels{"string": "2"}).Set(float64(pv2Voltage) / 10)

	pv2Current, err := m.client.ReadRegister(MODBUS_INVERTER_PV2_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvCurrent.With(prometheus.Labels{"string": "2"}).Set(float64(pv2Current) / 100)

	// string 3
	pv3Voltage, err := m.client.ReadRegister(MODBUS_INVERTER_PV3_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvVoltage.With(prometheus.Labels{"string": "3"}).Set(float64(pv3Voltage) / 10)

	pv3Current, err := m.client.ReadRegister(MODBUS_INVERTER_PV3_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvCurrent.With(prometheus.Labels{"string": "3"}).Set(float64(pv3Current) / 100)

	// string 4
	pv4Voltage, err := m.client.ReadRegister(MODBUS_INVERTER_PV4_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvVoltage.With(prometheus.Labels{"string": "4"}).Set(float64(pv4Voltage) / 10)

	pv4Current, err := m.client.ReadRegister(MODBUS_INVERTER_PV4_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvCurrent.With(prometheus.Labels{"string": "4"}).Set(float64(pv4Current) / 100)


	// phase A
	phaseAVoltage, err := m.client.ReadRegister(MODBUS_INVERTER_PHASE_A_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseVoltage.With(prometheus.Labels{"phase": "A"}).Set(float64(phaseAVoltage) / 10)

	phaseACurrent, err := m.client.ReadUint32(MODBUS_INVERTER_PHASE_A_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseCurrent.With(prometheus.Labels{"phase": "A"}).Set(float64(phaseACurrent) / 1000)

	// phase B
	phaseBVoltage, err := m.client.ReadRegister(MODBUS_INVERTER_PHASE_B_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseVoltage.With(prometheus.Labels{"phase": "B"}).Set(float64(phaseBVoltage) / 10)

	phaseBCurrent, err := m.client.ReadUint32(MODBUS_INVERTER_PHASE_B_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseCurrent.With(prometheus.Labels{"phase": "B"}).Set(float64(phaseBCurrent) / 1000)

	// phase C
	phaseCVoltage, err := m.client.ReadRegister(MODBUS_INVERTER_PHASE_C_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseVoltage.With(prometheus.Labels{"phase": "C"}).Set(float64(phaseCVoltage) / 10)

	phaseCCurrent, err := m.client.ReadUint32(MODBUS_INVERTER_PHASE_C_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseCurrent.With(prometheus.Labels{"phase": "C"}).Set(float64(phaseCCurrent) / 1000)


	// other
	activePower, err := m.client.ReadUint32(MODBUS_INVERTER_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	activePowerGauge.Set(float64(activePower) / 1000)

	reactivePower, err := m.client.ReadUint32(MODBUS_INVERTER_REACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	reactivePowerGauge.Set(float64(reactivePower) / 1000)

	powerFactor, err := m.client.ReadRegister(MODBUS_INVERTER_POWER_FACTOR, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	powerFactorGauge.Set(float64(powerFactor) / 1000)

	gridFrequency, err := m.client.ReadRegister(MODBUS_INVERTER_FREQUENCY, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	gridFrequencyGauge.Set(float64(gridFrequency) / 100)

	inverterEfficiency, err := m.client.ReadRegister(MODBUS_INVERTER_FREQUENCY, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	inverterEfficiencyGauge.Set(float64(inverterEfficiency) / 100)

	cabinetTemperature, err := m.client.ReadRegister(MODBUS_INVERTER_CABINET_TEMPERATURE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	cabinetTemperatureGauge.Set(float64(cabinetTemperature) / 10)

	isulationResistance, err := m.client.ReadRegister(MODBUS_INVERTER_INSULATION_RESISTANCE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	isulationResistanceGauge.Set(float64(isulationResistance) / 10)


	// Battery
	runningStatus, err := m.client.ReadRegister(MODBUS_BATTERY_1_RUNNING_STATUS, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	runningStatusGauge.With(prometheus.Labels{"battery": "1"}).Set(float64(runningStatus))

	chargingStatus, err := m.client.ReadUint32(MODBUS_BATTERY_1_CHARGING_STATUS, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	chargingStatusGauge.With(prometheus.Labels{"battery": "1"}).Set(float64(chargingStatus))

	busVoltage, err := m.client.ReadRegister(MODBUS_BATTERY_1_BUS_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	busVoltageGauge.With(prometheus.Labels{"battery": "1"}).Set(float64(busVoltage) / 10)

	batteryCapacity, err := m.client.ReadRegister(MODBUS_BATTERY_1_CAPACITY, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	batteryCapacityGauge.With(prometheus.Labels{"battery": "1"}).Set(float64(batteryCapacity) / 10)

	totalCharge, err := m.client.ReadUint32(MODBUS_BATTERY_1_TOTAL_CHARGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	totalChargeGauge.With(prometheus.Labels{"battery": "1"}).Set(float64(totalCharge) / 100)

	totalDischarge, err := m.client.ReadUint32(MODBUS_BATTERY_1_TOTAL_DISCHARGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	totalDischargeGauge.With(prometheus.Labels{"battery": "1"}).Set(float64(totalDischarge) / 100)
}
