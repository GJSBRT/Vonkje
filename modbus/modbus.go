package modbus

import (
	"context"
	"fmt"
	"time"

	"gijs.eu/vonkje/utils"

	"github.com/sirupsen/logrus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/simonvetter/modbus"
)

type Inverter struct {
	Name string `mapstructure:"name"`
	UnitId uint8 `mapstructure:"unit-id"`
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
	Inverters []Inverter `mapstructure:"inverters"`
}

type Connection struct {
	config ConnectionConfig
	client *modbus.ModbusClient
}

type Config struct {
	Run bool `mapstructure:"run"`
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
	if !m.config.Run {
		m.logger.Warn("Modbus metrics collector is disabled")
		return
	}

	m.logger.Info("Starting modbus metrics collector")

	m.updateMetrics()

	ticker := time.NewTicker(time.Duration(m.config.ReadMetricsInterval) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-m.ctx.Done():
			m.logger.Info("Stopping modbus metrics collector")
			return
		case <-ticker.C:
			go m.updateMetrics()
		}
	}
}

func (m *Modbus) updateMetrics() {
	for _, connection := range m.connections {
		for _, inverter := range connection.config.Inverters {
			err := m.updateSun2000Metrics(connection, inverter)
			if err != nil {
				m.errChannel <- err
			}

			if inverter.Luna2000 {
				err = m.updateLuna2000Metrics(connection, inverter)
				if err != nil {
					m.errChannel <- err
				}
			}

			if inverter.PowerMeter {
				err = m.updatePowerMeterMetrics(connection, inverter)
				if err != nil {
					m.errChannel <- err
				}
			}
		}
	}
}

func (m *Modbus) updateLuna2000Metrics(connection *Connection, inverter Inverter) error {
	m.logger.Infof("Updating luna2000 metrics for %s", inverter.Name)
	
	err := connection.client.SetUnitId(inverter.UnitId)
	if err != nil {
		return err
	}

	runningStatus, err := connection.client.ReadRegister(MODBUS_BATTERY_1_RUNNING_STATUS, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	runningStatusGauge.With(prometheus.Labels{"inverter": inverter.Name, "battery": "1"}).Set(float64(runningStatus))

	chargingStatus, err := connection.client.ReadUint32(MODBUS_BATTERY_1_CHARGING_STATUS, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if chargingStatus > 999999 {
		chargingStatusGauge.With(prometheus.Labels{"inverter": inverter.Name, "battery": "1"}).Set(0)
	} else {
		chargingStatusGauge.With(prometheus.Labels{"inverter": inverter.Name, "battery": "1"}).Set(float64(chargingStatus))
	}

	busVoltage, err := connection.client.ReadRegister(MODBUS_BATTERY_1_BUS_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	busVoltageGauge.With(prometheus.Labels{"inverter": inverter.Name, "battery": "1"}).Set(float64(busVoltage) / 10)

	batteryCapacity, err := connection.client.ReadRegister(MODBUS_BATTERY_1_CAPACITY, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	batteryCapacityGauge.With(prometheus.Labels{"inverter": inverter.Name, "battery": "1"}).Set(float64(batteryCapacity) / 10)

	totalCharge, err := connection.client.ReadUint32(MODBUS_BATTERY_1_TOTAL_CHARGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	totalChargeGauge.With(prometheus.Labels{"inverter": inverter.Name, "battery": "1"}).Set(float64(totalCharge) / 100)

	totalDischarge, err := connection.client.ReadUint32(MODBUS_BATTERY_1_TOTAL_DISCHARGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	totalDischargeGauge.With(prometheus.Labels{"inverter": inverter.Name, "battery": "1"}).Set(float64(totalDischarge) / 100)

	return nil
}

func (m *Modbus) updatePowerMeterMetrics(connection *Connection, inverter Inverter) error {
	m.logger.Infof("Updating power meter metrics for %s", inverter.Name)
	
	err := connection.client.SetUnitId(inverter.UnitId)
	if err != nil {
		return err
	}

	status, err := connection.client.ReadRegister(MODBUS_POWER_METER_STATUS, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterStatusGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(status))

	powerMeterPhaseAVoltage, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_A_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterPhaseVoltageGauge.With(prometheus.Labels{"inverter": inverter.Name, "phase": "A"}).Set(float64(powerMeterPhaseAVoltage) / 10)

	var powerMeterPhaseACurrentResult uint32
	powerMeterPhaseACurrent, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_A_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterPhaseACurrent > 999999 {
		powerMeterPhaseACurrentBytes, err := connection.client.ReadBytes(MODBUS_POWER_METER_PHASE_A_CURRENT, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		powerMeterPhaseACurrentResult = utils.ConvertTooLargeNumber(powerMeterPhaseACurrentBytes)
	} else {
		powerMeterPhaseACurrentResult = powerMeterPhaseACurrent
	}
	powerMeterPhaseCurrentGauge.With(prometheus.Labels{"inverter": inverter.Name, "phase": "A"}).Set(float64(powerMeterPhaseACurrentResult) / 100)

	powerMeterPhaseBVoltage, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_B_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterPhaseVoltageGauge.With(prometheus.Labels{"inverter": inverter.Name, "phase": "B"}).Set(float64(powerMeterPhaseBVoltage) / 10)

	var powerMeterPhaseBCurrentResult uint32
	powerMeterPhaseBCurrent, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_B_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterPhaseBCurrent > 999999 {
		powerMeterPhaseBCurrentBytes, err := connection.client.ReadBytes(MODBUS_POWER_METER_PHASE_B_CURRENT, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		powerMeterPhaseBCurrentResult = utils.ConvertTooLargeNumber(powerMeterPhaseBCurrentBytes)
	} else {
		powerMeterPhaseBCurrentResult = powerMeterPhaseBCurrent
	}
	powerMeterPhaseCurrentGauge.With(prometheus.Labels{"inverter": inverter.Name, "phase": "B"}).Set(float64(powerMeterPhaseBCurrentResult) / 100)

	powerMeterPhaseCVoltage, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_C_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterPhaseVoltageGauge.With(prometheus.Labels{"inverter": inverter.Name, "phase": "C"}).Set(float64(powerMeterPhaseCVoltage) / 10)

	var powerMeterPhaseCCurrentResult uint32
	powerMeterPhaseCCurrent, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_C_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterPhaseCCurrent > 999999 {
		powerMeterPhaseCCurrentBytes, err := connection.client.ReadBytes(MODBUS_POWER_METER_PHASE_C_CURRENT, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		powerMeterPhaseCCurrentResult = utils.ConvertTooLargeNumber(powerMeterPhaseCCurrentBytes)
	} else {
		powerMeterPhaseCCurrentResult = powerMeterPhaseCCurrent
	}
	powerMeterPhaseCurrentGauge.With(prometheus.Labels{"inverter": inverter.Name, "phase": "C"}).Set(float64(powerMeterPhaseCCurrentResult) / 100)

	var powerMeterActivePowerResult uint32
	powerMeterActivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterActivePower > 999999 {
		powerMeterActivePowerBytes, err := connection.client.ReadBytes(MODBUS_POWER_METER_ACTIVE_POWER, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		powerMeterActivePowerResult = utils.ConvertTooLargeNumber(powerMeterActivePowerBytes)
	} else {
		powerMeterActivePowerResult = powerMeterActivePower
	}
	powerMeterActivePowerGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(powerMeterActivePowerResult) / 100)

	powerMeterReactivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_REACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterReactivePowerGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(powerMeterReactivePower))

	powerMeterPowerFactor, err := connection.client.ReadRegister(MODBUS_POWER_METER_POWER_FACTOR, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterPowerFactorGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(powerMeterPowerFactor) / 1000)

	powerMeterFrequency, err := connection.client.ReadRegister(MODBUS_POWER_METER_FREQUENCY, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterFrequencyGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(powerMeterFrequency) / 100)

	positiveActiveElectricity, err := connection.client.ReadUint32(MODBUS_POWER_METER_POSITIVE_ACTIVE_ELECTRICITY, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterPositiveActiveElectricityGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(positiveActiveElectricity) / 100)
	
	reverseActivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_REVERSE_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterReverseActivePowerGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(reverseActivePower) / 100)

	accumulatedReactivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_ACCUMULATED_REACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterAccumulatedReactivePowerGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(accumulatedReactivePower) / 100)

	abLineVoltage, err := connection.client.ReadUint32(MODBUS_POWER_METER_AB_LINE_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterLineVoltageGauge.With(prometheus.Labels{"inverter": inverter.Name, "line": "AB"}).Set(float64(abLineVoltage) / 10)

	bcLineVoltage, err := connection.client.ReadUint32(MODBUS_POWER_METER_BC_LINE_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterLineVoltageGauge.With(prometheus.Labels{"inverter": inverter.Name, "line": "BC"}).Set(float64(bcLineVoltage) / 10)

	caLineVoltage, err := connection.client.ReadUint32(MODBUS_POWER_METER_CA_LINE_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterLineVoltageGauge.With(prometheus.Labels{"inverter": inverter.Name, "line": "CA"}).Set(float64(caLineVoltage) / 10)

	var phaseAActivePowerResult uint32
	powerMeterPhaseAActivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_A_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterPhaseAActivePower > 999999 {
		powerMeterPhaseAActivePowerBytes, err := connection.client.ReadBytes(MODBUS_POWER_METER_PHASE_A_ACTIVE_POWER, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		phaseAActivePowerResult = utils.ConvertTooLargeNumber(powerMeterPhaseAActivePowerBytes)
	} else {
		phaseAActivePowerResult = powerMeterPhaseAActivePower
	}
	powerMeterPhaseActivePowerGauge.With(prometheus.Labels{"inverter": inverter.Name, "phase": "A"}).Set(float64(phaseAActivePowerResult) / 100)

	var phaseBActivePowerResult uint32
	powerMeterPhaseBActivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_B_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterPhaseBActivePower > 999999 {
		powerMeterPhaseBActivePowerBytes, err := connection.client.ReadBytes(MODBUS_POWER_METER_PHASE_B_ACTIVE_POWER, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		phaseBActivePowerResult = utils.ConvertTooLargeNumber(powerMeterPhaseBActivePowerBytes)
	} else {
		phaseBActivePowerResult = powerMeterPhaseBActivePower
	}
	powerMeterPhaseActivePowerGauge.With(prometheus.Labels{"inverter": inverter.Name, "phase": "B"}).Set(float64(phaseBActivePowerResult) / 100)

	var phaseCActivePowerResult uint32
	powerMeterPhaseCActivePower, err := connection.client.ReadUint32(MODBUS_POWER_METER_PHASE_C_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterPhaseCActivePower > 999999 {
		powerMeterPhaseCActivePowerBytes, err := connection.client.ReadBytes(MODBUS_POWER_METER_PHASE_C_ACTIVE_POWER, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		phaseCActivePowerResult = utils.ConvertTooLargeNumber(powerMeterPhaseCActivePowerBytes)
	} else {
		phaseCActivePowerResult = powerMeterPhaseCActivePower
	}
	powerMeterPhaseActivePowerGauge.With(prometheus.Labels{"inverter": inverter.Name, "phase": "C"}).Set(float64(phaseCActivePowerResult) / 100)

	powerMeterModelResult, err := connection.client.ReadRegister(MODBUS_POWER_METER_MODEL_RESULT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	powerMeterModelResultGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(powerMeterModelResult))

	return nil
}

func (m *Modbus) updateSun2000Metrics(connection *Connection, inverter Inverter) error {
	m.logger.Infof("Updating sun2000 metrics for %s", inverter.Name)
	
	err := connection.client.SetUnitId(inverter.UnitId)
	if err != nil {
		return err
	}

	// string 1
	pv1Voltage, err := connection.client.ReadRegister(MODBUS_INVERTER_PV1_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvVoltage.With(prometheus.Labels{"inverter": inverter.Name, "string": "1"}).Set(float64(pv1Voltage) / 10)

	pv1Current, err := connection.client.ReadRegister(MODBUS_INVERTER_PV1_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvCurrent.With(prometheus.Labels{"inverter": inverter.Name, "string": "1"}).Set(float64(pv1Current) / 100)

	// string 2
	pv2Voltage, err := connection.client.ReadRegister(MODBUS_INVERTER_PV2_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvVoltage.With(prometheus.Labels{"inverter": inverter.Name, "string": "2"}).Set(float64(pv2Voltage) / 10)

	pv2Current, err := connection.client.ReadRegister(MODBUS_INVERTER_PV2_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvCurrent.With(prometheus.Labels{"inverter": inverter.Name, "string": "2"}).Set(float64(pv2Current) / 100)


	// phase A
	phaseAVoltage, err := connection.client.ReadRegister(MODBUS_INVERTER_PHASE_A_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseVoltage.With(prometheus.Labels{"inverter": inverter.Name, "phase": "A"}).Set(float64(phaseAVoltage) / 10)

	phaseACurrent, err := connection.client.ReadUint32(MODBUS_INVERTER_PHASE_A_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseCurrent.With(prometheus.Labels{"inverter": inverter.Name, "phase": "A"}).Set(float64(phaseACurrent) / 1000)

	// phase B
	phaseBVoltage, err := connection.client.ReadRegister(MODBUS_INVERTER_PHASE_B_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseVoltage.With(prometheus.Labels{"inverter": inverter.Name, "phase": "B"}).Set(float64(phaseBVoltage) / 10)

	phaseBCurrent, err := connection.client.ReadUint32(MODBUS_INVERTER_PHASE_B_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseCurrent.With(prometheus.Labels{"inverter": inverter.Name, "phase": "B"}).Set(float64(phaseBCurrent) / 1000)

	// phase C
	phaseCVoltage, err := connection.client.ReadRegister(MODBUS_INVERTER_PHASE_C_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseVoltage.With(prometheus.Labels{"inverter": inverter.Name, "phase": "C"}).Set(float64(phaseCVoltage) / 10)

	phaseCCurrent, err := connection.client.ReadUint32(MODBUS_INVERTER_PHASE_C_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseCurrent.With(prometheus.Labels{"inverter": inverter.Name, "phase": "C"}).Set(float64(phaseCCurrent) / 1000)


	// other
	inputPower, err := connection.client.ReadUint32(MODBUS_INVERTER_INPUT_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	sun2000InputPowerGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(inputPower) / 1000)

	stateOne, err := connection.client.ReadUint32(MODBUS_INVERTER_STATE_1, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	sun2000StateGauge.With(prometheus.Labels{"inverter": inverter.Name, "state": "1"}).Set(float64(stateOne))

	inverterDeviceStatus, err := connection.client.ReadUint32(MODBUS_INVERTER_DEVICE_STATUS, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	sun2000DeviceStatusGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(inverterDeviceStatus))

	activePower, err := connection.client.ReadUint32(MODBUS_INVERTER_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	activePowerGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(activePower) / 1000)

	reactivePower, err := connection.client.ReadUint32(MODBUS_INVERTER_REACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	reactivePowerGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(reactivePower) / 1000)

	powerFactor, err := connection.client.ReadRegister(MODBUS_INVERTER_POWER_FACTOR, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	powerFactorGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(powerFactor) / 1000)

	gridFrequency, err := connection.client.ReadRegister(MODBUS_INVERTER_FREQUENCY, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	gridFrequencyGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(gridFrequency) / 100)

	inverterEfficiency, err := connection.client.ReadRegister(MODBUS_INVERTER_FREQUENCY, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	inverterEfficiencyGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(inverterEfficiency) / 100)

	cabinetTemperature, err := connection.client.ReadRegister(MODBUS_INVERTER_CABINET_TEMPERATURE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	cabinetTemperatureGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(cabinetTemperature) / 10)

	isulationResistance, err := connection.client.ReadRegister(MODBUS_INVERTER_INSULATION_RESISTANCE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	isulationResistanceGauge.With(prometheus.Labels{"inverter": inverter.Name}).Set(float64(isulationResistance) / 10)

	return nil
}
