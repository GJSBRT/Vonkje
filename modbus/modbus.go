package modbus

import (
	"context"
	"fmt"
	"time"

	"gijs.eu/vonkje/utils"
	"gijs.eu/vonkje/metrics"

	"github.com/sirupsen/logrus"
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

func (m *Modbus) ChangeBatteryForceCharge(inverter string, battery string, state uint16, watts uint) error {
	inverterConfig, err := m.getInverterConfig(inverter)
	if err != nil {
		return err
	}

	connection, err := m.getConnection(inverter)
	if err != nil {
		return err
	}

	err = connection.client.SetUnitId(inverterConfig.UnitId)
	if err != nil {
		return err
	}

	if !inverterConfig.Luna2000 {
		return fmt.Errorf("Inverter %s does not have a luna2000 battery connected", inverter)
	}

	switch state {
	case MODBUS_STATE_BATTERY_1_FORCIBLE_CHARGE_DISCHARGE_CHARGE:
		err = connection.client.WriteUint32(MODBUS_ADDRESS_BATTERY_1_FORCIBLE_CHARGE_POWER, uint32(watts))
		if err != nil {
			return err
		}
	case MODBUS_STATE_BATTERY_1_FORCIBLE_CHARGE_DISCHARGE_DISCHARGE:
		err = connection.client.WriteUint32(MODBUS_ADDRESS_BATTERY_1_FORCIBLE_DISCHARGE_POWER, uint32(watts))
		if err != nil {
			return err
		}
	}

	err = connection.client.WriteRegister(MODBUS_ADDRESS_BATTERY_1_FORCIBLE_CHARGE_DISCHARGE, state)
	if err != nil {
		return err
	}

	return nil
}

func (m *Modbus) getInverterConfig(inverter string) (Inverter, error) {
	for _, connection := range m.connections {
		for _, inv := range connection.config.Inverters {
			if inv.Name == inverter {
				return inv, nil
			}
		}
	}

	return Inverter{}, fmt.Errorf("Inverter %s not found", inverter)
}

func (m *Modbus) getConnection(inverter string) (*Connection, error) {
	for _, connection := range m.connections {
		for _, inv := range connection.config.Inverters {
			if inv.Name == inverter {
				return connection, nil
			}
		}
	}

	return nil, fmt.Errorf("Inverter %s not found", inverter)
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

	runningStatus, err := connection.client.ReadRegister(MODBUS_ADDRESS_BATTERY_1_RUNNING_STATUS, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("luna2000", "running_status", map[string]string{"inverter": inverter.Name, "battery": "1"}, float64(runningStatus))

	chargingStatus, err := connection.client.ReadUint32(MODBUS_ADDRESS_BATTERY_1_CHARGING_STATUS, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if chargingStatus > 999999 {
		metrics.SetMetricValue("luna2000", "charging_status", map[string]string{"inverter": inverter.Name, "battery": "1"}, 0)
	} else {
		metrics.SetMetricValue("luna2000", "charging_status", map[string]string{"inverter": inverter.Name, "battery": "1"}, float64(chargingStatus))
	}

	busVoltage, err := connection.client.ReadRegister(MODBUS_ADDRESS_BATTERY_1_BUS_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("luna2000", "bus_voltage", map[string]string{"inverter": inverter.Name, "battery": "1"}, float64(busVoltage) / 10)

	batteryCapacity, err := connection.client.ReadRegister(MODBUS_ADDRESS_BATTERY_1_CAPACITY, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("luna2000", "battery_capacity", map[string]string{"inverter": inverter.Name, "battery": "1"}, float64(batteryCapacity) / 10)

	totalCharge, err := connection.client.ReadUint32(MODBUS_ADDRESS_BATTERY_1_TOTAL_CHARGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("luna2000", "total_charge", map[string]string{"inverter": inverter.Name, "battery": "1"}, float64(totalCharge) / 100)

	totalDischarge, err := connection.client.ReadUint32(MODBUS_ADDRESS_BATTERY_1_TOTAL_DISCHARGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("luna2000", "total_discharge", map[string]string{"inverter": inverter.Name, "battery": "1"}, float64(totalDischarge) / 100)

	return nil
}

func (m *Modbus) updatePowerMeterMetrics(connection *Connection, inverter Inverter) error {
	m.logger.Infof("Updating power meter metrics for %s", inverter.Name)
	
	err := connection.client.SetUnitId(inverter.UnitId)
	if err != nil {
		return err
	}

	status, err := connection.client.ReadRegister(MODBUS_ADDRESS_POWER_METER_STATUS, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "status", map[string]string{"inverter": inverter.Name}, float64(status))

	powerMeterPhaseAVoltage, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_PHASE_A_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "phase_voltage", map[string]string{"inverter": inverter.Name, "phase": "A"}, float64(powerMeterPhaseAVoltage) / 10)

	var powerMeterPhaseACurrentResult uint32
	powerMeterPhaseACurrent, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_PHASE_A_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterPhaseACurrent > 999999 {
		powerMeterPhaseACurrentBytes, err := connection.client.ReadBytes(MODBUS_ADDRESS_POWER_METER_PHASE_A_CURRENT, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		powerMeterPhaseACurrentResult = utils.ConvertTooLargeNumber(powerMeterPhaseACurrentBytes)
	} else {
		powerMeterPhaseACurrentResult = powerMeterPhaseACurrent
	}
	metrics.SetMetricValue("power_meter", "phase_current", map[string]string{"inverter": inverter.Name, "phase": "A"}, float64(powerMeterPhaseACurrentResult) / 100)

	powerMeterPhaseBVoltage, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_PHASE_B_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "phase_voltage", map[string]string{"inverter": inverter.Name, "phase": "B"}, float64(powerMeterPhaseBVoltage) / 10)

	var powerMeterPhaseBCurrentResult uint32
	powerMeterPhaseBCurrent, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_PHASE_B_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterPhaseBCurrent > 999999 {
		powerMeterPhaseBCurrentBytes, err := connection.client.ReadBytes(MODBUS_ADDRESS_POWER_METER_PHASE_B_CURRENT, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		powerMeterPhaseBCurrentResult = utils.ConvertTooLargeNumber(powerMeterPhaseBCurrentBytes)
	} else {
		powerMeterPhaseBCurrentResult = powerMeterPhaseBCurrent
	}
	metrics.SetMetricValue("power_meter", "phase_current", map[string]string{"inverter": inverter.Name, "phase": "B"}, float64(powerMeterPhaseBCurrentResult) / 100)

	powerMeterPhaseCVoltage, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_PHASE_C_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "phase_voltage", map[string]string{"inverter": inverter.Name, "phase": "C"}, float64(powerMeterPhaseCVoltage) / 10)

	var powerMeterPhaseCCurrentResult uint32
	powerMeterPhaseCCurrent, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_PHASE_C_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterPhaseCCurrent > 999999 {
		powerMeterPhaseCCurrentBytes, err := connection.client.ReadBytes(MODBUS_ADDRESS_POWER_METER_PHASE_C_CURRENT, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		powerMeterPhaseCCurrentResult = utils.ConvertTooLargeNumber(powerMeterPhaseCCurrentBytes)
	} else {
		powerMeterPhaseCCurrentResult = powerMeterPhaseCCurrent
	}
	metrics.SetMetricValue("power_meter", "phase_current", map[string]string{"inverter": inverter.Name, "phase": "C"}, float64(powerMeterPhaseCCurrentResult) / 100)

	var powerMeterActivePowerResult uint32
	powerMeterActivePower, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterActivePower > 999999 {
		powerMeterActivePowerBytes, err := connection.client.ReadBytes(MODBUS_ADDRESS_POWER_METER_ACTIVE_POWER, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		powerMeterActivePowerResult = utils.ConvertTooLargeNumber(powerMeterActivePowerBytes)
	} else {
		powerMeterActivePowerResult = powerMeterActivePower
	}
	metrics.SetMetricValue("power_meter", "active_power", map[string]string{"inverter": inverter.Name}, float64(powerMeterActivePowerResult) / 100)

	powerMeterReactivePower, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_REACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "reactive_power", map[string]string{"inverter": inverter.Name}, float64(powerMeterReactivePower) / 100)

	powerMeterPowerFactor, err := connection.client.ReadRegister(MODBUS_ADDRESS_POWER_METER_POWER_FACTOR, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "power_factor", map[string]string{"inverter": inverter.Name}, float64(powerMeterPowerFactor) / 1000)

	powerMeterFrequency, err := connection.client.ReadRegister(MODBUS_ADDRESS_POWER_METER_FREQUENCY, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "frequency", map[string]string{"inverter": inverter.Name}, float64(powerMeterFrequency) / 100)

	positiveActiveElectricity, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_POSITIVE_ACTIVE_ELECTRICITY, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "positive_active_electricity", map[string]string{"inverter": inverter.Name}, float64(positiveActiveElectricity) / 100)
	
	reverseActivePower, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_REVERSE_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "reverse_active_power", map[string]string{"inverter": inverter.Name}, float64(reverseActivePower) / 100)

	accumulatedReactivePower, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_ACCUMULATED_REACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "accumulated_reactive_power", map[string]string{"inverter": inverter.Name}, float64(accumulatedReactivePower) / 100)

	abLineVoltage, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_AB_LINE_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "line_voltage", map[string]string{"inverter": inverter.Name, "line": "AB"}, float64(abLineVoltage) / 10)

	bcLineVoltage, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_BC_LINE_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "line_voltage", map[string]string{"inverter": inverter.Name, "line": "BC"}, float64(bcLineVoltage) / 10)

	caLineVoltage, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_CA_LINE_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "line_voltage", map[string]string{"inverter": inverter.Name, "line": "CA"}, float64(caLineVoltage) / 10)

	var phaseAActivePowerResult uint32
	powerMeterPhaseAActivePower, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_PHASE_A_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterPhaseAActivePower > 999999 {
		powerMeterPhaseAActivePowerBytes, err := connection.client.ReadBytes(MODBUS_ADDRESS_POWER_METER_PHASE_A_ACTIVE_POWER, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		phaseAActivePowerResult = utils.ConvertTooLargeNumber(powerMeterPhaseAActivePowerBytes)
	} else {
		phaseAActivePowerResult = powerMeterPhaseAActivePower
	}
	metrics.SetMetricValue("power_meter", "phase_active_power", map[string]string{"inverter": inverter.Name, "phase": "A"}, float64(phaseAActivePowerResult) / 100)

	var phaseBActivePowerResult uint32
	powerMeterPhaseBActivePower, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_PHASE_B_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterPhaseBActivePower > 999999 {
		powerMeterPhaseBActivePowerBytes, err := connection.client.ReadBytes(MODBUS_ADDRESS_POWER_METER_PHASE_B_ACTIVE_POWER, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		phaseBActivePowerResult = utils.ConvertTooLargeNumber(powerMeterPhaseBActivePowerBytes)
	} else {
		phaseBActivePowerResult = powerMeterPhaseBActivePower
	}
	metrics.SetMetricValue("power_meter", "phase_active_power", map[string]string{"inverter": inverter.Name, "phase": "B"}, float64(phaseBActivePowerResult) / 100)

	var phaseCActivePowerResult uint32
	powerMeterPhaseCActivePower, err := connection.client.ReadUint32(MODBUS_ADDRESS_POWER_METER_PHASE_C_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	if powerMeterPhaseCActivePower > 999999 {
		powerMeterPhaseCActivePowerBytes, err := connection.client.ReadBytes(MODBUS_ADDRESS_POWER_METER_PHASE_C_ACTIVE_POWER, 4, modbus.HOLDING_REGISTER)
		if err != nil {
			return err
		}
		
		phaseCActivePowerResult = utils.ConvertTooLargeNumber(powerMeterPhaseCActivePowerBytes)
	} else {
		phaseCActivePowerResult = powerMeterPhaseCActivePower
	}
	metrics.SetMetricValue("power_meter", "phase_active_power", map[string]string{"inverter": inverter.Name, "phase": "C"}, float64(phaseCActivePowerResult) / 100)

	powerMeterModelResult, err := connection.client.ReadRegister(MODBUS_ADDRESS_POWER_METER_MODEL_RESULT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("power_meter", "model_result", map[string]string{"inverter": inverter.Name}, float64(powerMeterModelResult))

	return nil
}

func (m *Modbus) updateSun2000Metrics(connection *Connection, inverter Inverter) error {
	m.logger.Infof("Updating sun2000 metrics for %s", inverter.Name)
	
	err := connection.client.SetUnitId(inverter.UnitId)
	if err != nil {
		return err
	}

	// string 1
	pv1Voltage, err := connection.client.ReadRegister(MODBUS_ADDRESS_INVERTER_PV1_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "pv_voltage", map[string]string{"inverter": inverter.Name, "string": "1"}, float64(pv1Voltage) / 10)

	pv1Current, err := connection.client.ReadRegister(MODBUS_ADDRESS_INVERTER_PV1_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "pv_current", map[string]string{"inverter": inverter.Name, "string": "1"}, float64(pv1Current) / 100)

	// string 2
	pv2Voltage, err := connection.client.ReadRegister(MODBUS_ADDRESS_INVERTER_PV2_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "pv_voltage", map[string]string{"inverter": inverter.Name, "string": "2"}, float64(pv2Voltage) / 10)

	pv2Current, err := connection.client.ReadRegister(MODBUS_ADDRESS_INVERTER_PV2_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "pv_current", map[string]string{"inverter": inverter.Name, "string": "2"}, float64(pv2Current) / 100)


	// phase A
	phaseAVoltage, err := connection.client.ReadRegister(MODBUS_ADDRESS_INVERTER_PHASE_A_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "phase_voltage", map[string]string{"inverter": inverter.Name, "phase": "A"}, float64(phaseAVoltage) / 10)

	phaseACurrent, err := connection.client.ReadUint32(MODBUS_ADDRESS_INVERTER_PHASE_A_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "phase_current", map[string]string{"inverter": inverter.Name, "phase": "A"}, float64(phaseACurrent) / 1000)

	// phase B
	phaseBVoltage, err := connection.client.ReadRegister(MODBUS_ADDRESS_INVERTER_PHASE_B_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "phase_voltage", map[string]string{"inverter": inverter.Name, "phase": "B"}, float64(phaseBVoltage) / 10)

	phaseBCurrent, err := connection.client.ReadUint32(MODBUS_ADDRESS_INVERTER_PHASE_B_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "phase_current", map[string]string{"inverter": inverter.Name, "phase": "B"}, float64(phaseBCurrent) / 1000)

	// phase C
	phaseCVoltage, err := connection.client.ReadRegister(MODBUS_ADDRESS_INVERTER_PHASE_C_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "phase_voltage", map[string]string{"inverter": inverter.Name, "phase": "C"}, float64(phaseCVoltage) / 10)

	phaseCCurrent, err := connection.client.ReadUint32(MODBUS_ADDRESS_INVERTER_PHASE_C_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "phase_current", map[string]string{"inverter": inverter.Name, "phase": "C"}, float64(phaseCCurrent) / 1000)

	inputPower, err := connection.client.ReadUint32(MODBUS_ADDRESS_INVERTER_INPUT_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "input_power", map[string]string{"inverter": inverter.Name}, float64(inputPower) / 1000)

	stateOne, err := connection.client.ReadUint32(MODBUS_ADDRESS_INVERTER_STATE_1, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "state", map[string]string{"inverter": inverter.Name, "state": "1"}, float64(stateOne))

	inverterDeviceStatus, err := connection.client.ReadUint32(MODBUS_ADDRESS_INVERTER_DEVICE_STATUS, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "device_status", map[string]string{"inverter": inverter.Name}, float64(inverterDeviceStatus))

	activePower, err := connection.client.ReadUint32(MODBUS_ADDRESS_INVERTER_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "active_power", map[string]string{"inverter": inverter.Name}, float64(activePower) / 1000)

	reactivePower, err := connection.client.ReadUint32(MODBUS_ADDRESS_INVERTER_REACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "reactive_power", map[string]string{"inverter": inverter.Name}, float64(reactivePower) / 1000)

	powerFactor, err := connection.client.ReadRegister(MODBUS_ADDRESS_INVERTER_POWER_FACTOR, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "power_factor", map[string]string{"inverter": inverter.Name}, float64(powerFactor) / 1000)

	gridFrequency, err := connection.client.ReadRegister(MODBUS_ADDRESS_INVERTER_FREQUENCY, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "frequency", map[string]string{"inverter": inverter.Name}, float64(gridFrequency) / 100)

	inverterEfficiency, err := connection.client.ReadRegister(MODBUS_ADDRESS_INVERTER_FREQUENCY, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "efficiency", map[string]string{"inverter": inverter.Name}, float64(inverterEfficiency) / 100)

	cabinetTemperature, err := connection.client.ReadRegister(MODBUS_ADDRESS_INVERTER_CABINET_TEMPERATURE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "cabinet_temperature", map[string]string{"inverter": inverter.Name}, float64(cabinetTemperature) / 10)

	isulationResistance, err := connection.client.ReadRegister(MODBUS_ADDRESS_INVERTER_INSULATION_RESISTANCE, modbus.HOLDING_REGISTER)
	if err != nil {
		return err
	}
	metrics.SetMetricValue("sun2000", "isulation_resistance", map[string]string{"inverter": inverter.Name}, float64(isulationResistance) / 10)

	return nil
}
