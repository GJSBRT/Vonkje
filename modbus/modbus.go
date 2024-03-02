package modbus

import (
	"context"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/simonvetter/modbus"
)

type ConnectionPeripheral struct {
	PowerMeter bool `mapstructure:"power-meter"`
	Luna2000 bool `mapstructure:"luna2000"`
}

type ConnectionConfig struct {
	Name string `mapstructure:"name"`
	IP string `mapstructure:"ip"`
	Port uint `mapstructure:"port"`
	Protocol string `mapstructure:"protocol"`
	Baudrate uint `mapstructure:"baudrate"`
	DataBits uint `mapstructure:"data-bits"`
	StopBits uint `mapstructure:"stop-bits"`
	Timeout uint `mapstructure:"timeout"`
	UnitId uint8 `mapstructure:"unit-id"`
	Peripherals ConnectionPeripheral `mapstructure:"peripherals"`
}

type Connection struct {
	config ConnectionConfig
	client *modbus.ModbusClient
}

type Config struct {
	ReadMetricsInterval uint `mapstructure:"read-metrics-interval"`
	Connections []ConnectionConfig `mapstructure:"connections"`
}

type Modbus struct {
	config Config
	errChannel chan error
	ctx context.Context
	logger *logrus.Logger
	connections map[string]*Connection
}

func New(
	config Config, 
	errChannel chan error,
	ctx context.Context,
	logger *logrus.Logger,
) (*Modbus, error) {
	m := &Modbus{
		connections: make(map[string]*Connection),
	}

	for _, connectionConfig := range config.Connections {
		client, err := modbus.NewClient(&modbus.ClientConfiguration{
			URL: fmt.Sprintf("%s://%s:%d", connectionConfig.Protocol, connectionConfig.IP, connectionConfig.Port),
			Speed: connectionConfig.Baudrate,
			DataBits: connectionConfig.DataBits,
			StopBits: connectionConfig.StopBits,
			Parity: modbus.PARITY_NONE,
			Timeout: time.Duration(connectionConfig.Timeout) * time.Second,
		})
		if err != nil {
			return nil, err
		}
	
		err = client.SetUnitId(connectionConfig.UnitId)
		if err != nil {
			return nil, err
		}

		err = client.Open()
		if err != nil {
			return nil, err
		}

		logger.Infof("Connected to %s://%s:%d", connectionConfig.Protocol, connectionConfig.IP, connectionConfig.Port)
		m.connections[fmt.Sprintf("%s:%d", connectionConfig.IP, connectionConfig.Port)] = &Connection{
			config: connectionConfig,
			client: client,
		}
	}

	m.config = config
	m.errChannel = errChannel
	m.ctx = ctx
	m.logger = logger

	return m, nil
}

func (m *Modbus) Close() {
	for _, connection := range m.connections {
		connection.client.Close()
	}
}

func (m *Modbus) Start() {
	m.logger.Info("Starting modbus metrics collector")

	for _, connection := range m.connections {
		go m.updateMetrics(connection)
	}

	ticker := time.NewTicker(time.Duration(m.config.ReadMetricsInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			m.logger.Info("Stopping modbus metrics collector")
			return
		case <-ticker.C:
			for _, connection := range m.connections {
				go m.updateMetrics(connection)
			}
		}
	}
}

func (m *Modbus) updateMetrics(connection *Connection) {
	m.logger.Infof("Updating metrics for %s", connection.config.Name)

	// string 1
	pv1Voltage, err := connection.client.ReadRegister(MODBUS_INVERTER_PV1_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvVoltage.With(prometheus.Labels{"connection": connection.config.Name, "string": "1"}).Set(float64(pv1Voltage) / 10)

	pv1Current, err := connection.client.ReadRegister(MODBUS_INVERTER_PV1_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvCurrent.With(prometheus.Labels{"connection": connection.config.Name, "string": "1"}).Set(float64(pv1Current) / 100)

	// string 2
	pv2Voltage, err := connection.client.ReadRegister(MODBUS_INVERTER_PV2_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvVoltage.With(prometheus.Labels{"connection": connection.config.Name, "string": "2"}).Set(float64(pv2Voltage) / 10)

	pv2Current, err := connection.client.ReadRegister(MODBUS_INVERTER_PV2_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvCurrent.With(prometheus.Labels{"connection": connection.config.Name, "string": "2"}).Set(float64(pv2Current) / 100)


	// phase A
	phaseAVoltage, err := connection.client.ReadRegister(MODBUS_INVERTER_PHASE_A_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseVoltage.With(prometheus.Labels{"connection": connection.config.Name, "phase": "A"}).Set(float64(phaseAVoltage) / 10)

	phaseACurrent, err := connection.client.ReadUint32(MODBUS_INVERTER_PHASE_A_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseCurrent.With(prometheus.Labels{"connection": connection.config.Name, "phase": "A"}).Set(float64(phaseACurrent) / 1000)

	// phase B
	phaseBVoltage, err := connection.client.ReadRegister(MODBUS_INVERTER_PHASE_B_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseVoltage.With(prometheus.Labels{"connection": connection.config.Name, "phase": "B"}).Set(float64(phaseBVoltage) / 10)

	phaseBCurrent, err := connection.client.ReadUint32(MODBUS_INVERTER_PHASE_B_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseCurrent.With(prometheus.Labels{"connection": connection.config.Name, "phase": "B"}).Set(float64(phaseBCurrent) / 1000)

	// phase C
	phaseCVoltage, err := connection.client.ReadRegister(MODBUS_INVERTER_PHASE_C_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseVoltage.With(prometheus.Labels{"connection": connection.config.Name, "phase": "C"}).Set(float64(phaseCVoltage) / 10)

	phaseCCurrent, err := connection.client.ReadUint32(MODBUS_INVERTER_PHASE_C_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseCurrent.With(prometheus.Labels{"connection": connection.config.Name, "phase": "C"}).Set(float64(phaseCCurrent) / 1000)


	// other
	activePower, err := connection.client.ReadUint32(MODBUS_INVERTER_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	activePowerGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(activePower) / 1000)

	reactivePower, err := connection.client.ReadUint32(MODBUS_INVERTER_REACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	reactivePowerGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(reactivePower) / 1000)

	powerFactor, err := connection.client.ReadRegister(MODBUS_INVERTER_POWER_FACTOR, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	powerFactorGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(powerFactor) / 1000)

	gridFrequency, err := connection.client.ReadRegister(MODBUS_INVERTER_FREQUENCY, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	gridFrequencyGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(gridFrequency) / 100)

	inverterEfficiency, err := connection.client.ReadRegister(MODBUS_INVERTER_FREQUENCY, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	inverterEfficiencyGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(inverterEfficiency) / 100)

	cabinetTemperature, err := connection.client.ReadRegister(MODBUS_INVERTER_CABINET_TEMPERATURE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	cabinetTemperatureGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(cabinetTemperature) / 10)

	isulationResistance, err := connection.client.ReadRegister(MODBUS_INVERTER_INSULATION_RESISTANCE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	isulationResistanceGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(isulationResistance) / 10)


	// Battery
	if connection.config.Peripherals.Luna2000 {
		runningStatus, err := connection.client.ReadRegister(MODBUS_BATTERY_1_RUNNING_STATUS, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		runningStatusGauge.With(prometheus.Labels{"connection": connection.config.Name, "battery": "1"}).Set(float64(runningStatus))

		chargingStatus, err := connection.client.ReadUint32(MODBUS_BATTERY_1_CHARGING_STATUS, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		chargingStatusGauge.With(prometheus.Labels{"connection": connection.config.Name, "battery": "1"}).Set(float64(chargingStatus))

		busVoltage, err := connection.client.ReadRegister(MODBUS_BATTERY_1_BUS_VOLTAGE, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		busVoltageGauge.With(prometheus.Labels{"connection": connection.config.Name, "battery": "1"}).Set(float64(busVoltage) / 10)

		batteryCapacity, err := connection.client.ReadRegister(MODBUS_BATTERY_1_CAPACITY, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		batteryCapacityGauge.With(prometheus.Labels{"connection": connection.config.Name, "battery": "1"}).Set(float64(batteryCapacity) / 10)

		totalCharge, err := connection.client.ReadUint32(MODBUS_BATTERY_1_TOTAL_CHARGE, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		totalChargeGauge.With(prometheus.Labels{"connection": connection.config.Name, "battery": "1"}).Set(float64(totalCharge) / 100)

		totalDischarge, err := connection.client.ReadUint32(MODBUS_BATTERY_1_TOTAL_DISCHARGE, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		totalDischargeGauge.With(prometheus.Labels{"connection": connection.config.Name, "battery": "1"}).Set(float64(totalDischarge) / 100)
	}


	// Power meter
	if connection.config.Peripherals.PowerMeter {
		status, err := connection.client.ReadRegister(MODBUS_POWER_METER_STATUS, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterStatusGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(status))

		powerMeterPhaseAVoltage, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_A_VOLTAGE, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterPhaseVoltageGauge.With(prometheus.Labels{"connection": connection.config.Name, "phase": "A"}).Set(float64(powerMeterPhaseAVoltage) / 10)

		powerMeterPhaseACurrent, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_A_CURRENT, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterPhaseCurrentGauge.With(prometheus.Labels{"connection": connection.config.Name, "phase": "A"}).Set(float64(powerMeterPhaseACurrent) / 1000)

		powerMeterPhaseBVoltage, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_B_VOLTAGE, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterPhaseVoltageGauge.With(prometheus.Labels{"connection": connection.config.Name, "phase": "B"}).Set(float64(powerMeterPhaseBVoltage) / 10)

		powerMeterPhaseBCurrent, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_B_CURRENT, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterPhaseCurrentGauge.With(prometheus.Labels{"connection": connection.config.Name, "phase": "B"}).Set(float64(powerMeterPhaseBCurrent) / 1000)

		powerMeterPhaseCVoltage, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_C_VOLTAGE, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterPhaseVoltageGauge.With(prometheus.Labels{"connection": connection.config.Name, "phase": "C"}).Set(float64(powerMeterPhaseCVoltage) / 10)

		powerMeterPhaseCCurrent, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_C_CURRENT, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterPhaseCurrentGauge.With(prometheus.Labels{"connection": connection.config.Name, "phase": "C"}).Set(float64(powerMeterPhaseCCurrent) / 1000)

		powerMeterActivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_ACTIVE_POWER, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterActivePowerGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(powerMeterActivePower) / 1000)

		powerMeterReactivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_REACTIVE_POWER, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterReactivePowerGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(powerMeterReactivePower) / 1000)

		powerMeterPowerFactor, err := connection.client.ReadRegister(MODBUS_POWER_METER_POWER_FACTOR, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterPowerFactorGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(powerMeterPowerFactor) / 1000)

		powerMeterFrequency, err := connection.client.ReadRegister(MODBUS_POWER_METER_FREQUENCY, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterFrequencyGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(powerMeterFrequency) / 100)

		positiveActiveElectricity, err := connection.client.ReadUint32(MODBUS_POWER_METER_POSITIVE_ACTIVE_ELECTRICITY, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterPositiveActiveElectricityGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(positiveActiveElectricity) / 1000)
		
		reverseActivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_REVERSE_ACTIVE_POWER, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterReverseActivePowerGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(reverseActivePower) / 1000)

		accumulatedReactivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_ACCUMULATED_REACTIVE_POWER, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterAccumulatedReactivePowerGauge.With(prometheus.Labels{"connection": connection.config.Name}).Set(float64(accumulatedReactivePower) / 1000)

		abLineVoltage, err := connection.client.ReadUint32(MODBUS_POWER_METER_AB_LINE_VOLTAGE, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterLineVoltageGauge.With(prometheus.Labels{"connection": connection.config.Name, "line": "AB"}).Set(float64(abLineVoltage) / 10)

		bcLineVoltage, err := connection.client.ReadUint32(MODBUS_POWER_METER_BC_LINE_VOLTAGE, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterLineVoltageGauge.With(prometheus.Labels{"connection": connection.config.Name, "line": "BC"}).Set(float64(bcLineVoltage) / 10)

		caLineVoltage, err := connection.client.ReadUint32(MODBUS_POWER_METER_CA_LINE_VOLTAGE, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterLineVoltageGauge.With(prometheus.Labels{"connection": connection.config.Name, "line": "CA"}).Set(float64(caLineVoltage) / 10)

		phaseAActivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_A_ACTIVE_POWER, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterPhaseActivePowerGauge.With(prometheus.Labels{"connection": connection.config.Name, "phase": "A"}).Set(float64(phaseAActivePower) / 1000)

		phaseBActivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_B_ACTIVE_POWER, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterPhaseActivePowerGauge.With(prometheus.Labels{"connection": connection.config.Name, "phase": "B"}).Set(float64(phaseBActivePower) / 1000)

		phaseCActivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_C_ACTIVE_POWER, modbus.HOLDING_REGISTER)
		if err != nil {
			panic(err)
		}
		powerMeterPhaseActivePowerGauge.With(prometheus.Labels{"connection": connection.config.Name, "phase": "C"}).Set(float64(phaseCActivePower) / 1000)
	}
}
