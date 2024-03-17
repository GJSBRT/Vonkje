package modbus

import (
	"fmt"
	"time"
	"context"

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

		err = client.SetEncoding(modbus.BIG_ENDIAN, modbus.HIGH_WORD_FIRST)
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

	if !inverterConfig.Luna2000 {
		return fmt.Errorf("Inverter %s does not have a luna2000 battery connected", inverter)
	}

	connection, err := m.getConnection(inverter)
	if err != nil {
		return err
	}

	err = connection.client.SetUnitId(inverterConfig.UnitId)
	if err != nil {
		return err
	}

	switch state {
	case MODBUS_STATE_BATTERY_1_FORCIBLE_CHARGE_DISCHARGE_CHARGE:
		err = connection.client.WriteUint32(luna2000Registers["forcible_charge_power_battery_1"].Address, uint32(watts))
		if err != nil {
			return err
		}
	case MODBUS_STATE_BATTERY_1_FORCIBLE_CHARGE_DISCHARGE_DISCHARGE:
		err = connection.client.WriteUint32(luna2000Registers["forcible_discharge_power_battery_1"].Address, uint32(watts))
		if err != nil {
			return err
		}
	}

	err = connection.client.WriteRegister(luna2000Registers["forcible_charge_discharge_battery_1"].Address, state)
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
			err := m.updateMetricsRegisters(connection, inverter, sun2000Registers)
			if err != nil {
				m.errChannel <- err
			}

			if inverter.Luna2000 {
				err := m.updateMetricsRegisters(connection, inverter, luna2000Registers)
				if err != nil {
					m.errChannel <- err
				}
			}

			if inverter.PowerMeter {
				err := m.updateMetricsRegisters(connection, inverter, powerMeterRegisters)
				if err != nil {
					m.errChannel <- err
				}
			}
		}
	}
}

func (m *Modbus) updateMetricsRegisters(connection *Connection, inverter Inverter, registers map[string]Register) error {
	err := connection.client.SetUnitId(inverter.UnitId)
	if err != nil {
		return err
	}

	for _, register := range registers {
		var result int

		switch register.Type {
		case RegisterTypeUint16:
			var reg uint16
			reg, err = connection.client.ReadRegister(register.Address, modbus.HOLDING_REGISTER)
			if err != nil { 
				return err
			}

			result = int(reg)
		case RegisterTypeUint32:
			var reg uint32
			reg, err = connection.client.ReadUint32(register.Address, modbus.HOLDING_REGISTER)
			if err != nil {
				return err
			}

			result = int(reg)
		case RegisterTypeInt16:
			var res int16
			reg, err := connection.client.ReadRegister(register.Address, modbus.HOLDING_REGISTER)
			if err != nil {
				return err
			}

			res = int16(reg)
			result = int(res)
		case RegisterTypeInt32:
			var res int32
			reg, err := connection.client.ReadUint32(register.Address, modbus.HOLDING_REGISTER)
			if err != nil {
				return err
			}

			res = int32(reg)
			result = int(res)
		}

		fields := map[string]string{"inverter": inverter.Name}
		for k, v := range register.Fields {
			fields[k] = v
		}

		metrics.SetMetricValue(register.Namespace, register.Name, fields, float64(result) / register.Gain)
	}

	return nil
}
