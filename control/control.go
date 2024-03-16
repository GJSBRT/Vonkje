package control

import (
	"time"
	"math"
	"context"

	"gijs.eu/vonkje/modbus"
	"gijs.eu/vonkje/metrics"
	"gijs.eu/vonkje/packages/victoria_metrics"

	"github.com/spf13/viper"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Run bool `mapstructure:"run"`
	LoopInterval int `mapstructure:"loop-interval"`
	MetricsAvgPeriod int `mapstructure:"metrics-avg-period"`
	MinimumSolarOverProduction int `mapstructure:"minimum-solar-over-production"`
}

type Control struct {
	config Config
	errChannel chan error
	ctx context.Context
	logger *logrus.Logger
	victoriaMetrics *victoria_metrics.VictoriaMetrics
	modbus *modbus.Modbus
}

func New(
	config Config,
	errChannel chan error,
	ctx context.Context,
	logger *logrus.Logger,
	victoriaMetrics *victoria_metrics.VictoriaMetrics,
	modbus *modbus.Modbus,
) *Control {
	return &Control{
		config: config,
		errChannel: errChannel,
		ctx: ctx,
		logger: logger,
		victoriaMetrics: victoriaMetrics,
		modbus: modbus,
	}
}

type batteryState struct {
	inverter string
	battery string
	capacity float64
}

func (c *Control) Start() {
	if !c.config.Run {
		c.logger.Warn("Control loop is disabled")
		return
	}

	c.logger.Infof("Waiting %d minutes before starting control loop to collect metrics", c.config.MetricsAvgPeriod)
	time.Sleep(time.Duration(c.config.MetricsAvgPeriod) * time.Minute)

	c.logger.Info("Starting control loop")

	ticker := time.NewTicker(time.Duration(c.config.LoopInterval) * time.Second)
	defer ticker.Stop()

	entriesWanted := uint(c.config.MetricsAvgPeriod * (60 / viper.GetInt("modbus.read-metrics-interval")))

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Info("Stopping control loop")
			return
		case <-ticker.C:
			c.logger.Debug("Control loop tick")

			// 1. Get current home energy consumption
			avgPhaseVoltage, err := metrics.GetMetricAverage("power_meter", "phase_voltage", entriesWanted)
			if err != nil {
				c.errChannel <- err
				continue
			}

			avgPhaseCurrent, err := metrics.GetMetricAverage("power_meter", "phase_current", entriesWanted)
			if err != nil {
				c.errChannel <- err
				continue
			}

			avgHomeLoad := math.Ceil(avgPhaseVoltage * avgPhaseCurrent)

			// 2. Get current solar production
			avgSolarIn, err := metrics.GetMetricAverageSum("sun2000", "input_power", entriesWanted)
			if err != nil {
				c.errChannel <- err
				continue
			}
			avgSolarIn = math.Ceil(avgSolarIn * 1000)
			c.logger.WithFields(logrus.Fields{"avgSolarIn": avgSolarIn, "avgHomeLoad": avgHomeLoad}).Info("Solar production and home load")

			// 3. Get current battery capacities
			batteryMetricValues, err := metrics.GetMetricValues("luna2000", "battery_capacity")
			if err != nil {
				c.errChannel <- err
				continue
			}
			batteries := []batteryState{}
			for _, batteryMetricValue := range batteryMetricValues {
				batteries = append(batteries, batteryState{
					inverter: batteryMetricValue.Fields["inverter"],
					battery: batteryMetricValue.Fields["battery"],
					capacity: batteryMetricValue.Values[len(batteryMetricValue.Values) - 1],
				})
			}

			// Get over production in percentage
			var overProduction float64
			if avgSolarIn > avgHomeLoad {
				overProduction = math.Ceil((avgSolarIn - avgHomeLoad) / avgSolarIn * 100)
			} else {
				overProduction = 0
			}
			metrics.SetMetricValue("control", "over_production", map[string]string{}, overProduction)
			c.logger.WithFields(logrus.Fields{"overProduction": overProduction}).Info("Over production")

			// 4. if solar over production is more than x%, charge battery
			if overProduction > float64(c.config.MinimumSolarOverProduction) {
				metrics.SetMetricValue("control", "action", map[string]string{"action": "charge_batteries"}, 1)

				for _, battery := range batteries {
					if battery.capacity < 100 {
						c.logger.WithFields(logrus.Fields{"inverter": battery.inverter, "battery": battery.battery, "capacity": battery.capacity}).Info("Battery is not fully charged, starting charge")
						err := c.modbus.ChangeBatteryForceCharge(battery.inverter, battery.battery, modbus.MODBUS_STATE_BATTERY_1_FORCIBLE_CHARGE_DISCHARGE_CHARGE)
						if err != nil {
							c.errChannel <- err
							continue
						}
					} else {
						c.logger.WithFields(logrus.Fields{"inverter": battery.inverter, "battery": battery.battery}).Info("Battery is fully charged, stopping charge")
						err := c.modbus.ChangeBatteryForceCharge(battery.inverter, battery.battery, modbus.MODBUS_STATE_BATTERY_1_FORCIBLE_CHARGE_DISCHARGE_STOP)
						if err != nil {
							c.errChannel <- err
							continue
						}
					}
				}

				continue
			} else {
				metrics.SetMetricValue("control", "action", map[string]string{"action": "charge_batteries"}, 0)
			}

			// 5. if solar production < home energy consumption && battery capacity > 5%, discharge battery
			wattsRequired := math.Ceil(avgHomeLoad - avgSolarIn)
			batteriesRequired := math.Ceil(wattsRequired / 1500)
			if avgSolarIn < avgHomeLoad {
				if batteriesRequired > 0 {
					metrics.SetMetricValue("control", "action", map[string]string{"action": "discharge_battery"}, 1)
				} else {
					metrics.SetMetricValue("control", "action", map[string]string{"action": "discharge_battery"}, 1)
				}

				for _, battery := range batteries {
					if batteriesRequired > 0 {
						if battery.capacity > 5 {
							batteriesRequired--
							wattsRequired -= 1500

							c.logger.WithFields(logrus.Fields{"inverter": battery.inverter, "battery": battery.battery, "capacity": battery.capacity}).Info("Discharging battery")

							err := c.modbus.ChangeBatteryForceCharge(battery.inverter, battery.battery, modbus.MODBUS_STATE_BATTERY_1_FORCIBLE_CHARGE_DISCHARGE_DISCHARGE)
							if err != nil {
								c.errChannel <- err
							}

							continue
						}
					} else {
						c.logger.WithFields(logrus.Fields{"inverter": battery.inverter, "battery": battery.battery}).Info("Battery is not required, stopping discharge")

						err := c.modbus.ChangeBatteryForceCharge(battery.inverter, battery.battery, modbus.MODBUS_STATE_BATTERY_1_FORCIBLE_CHARGE_DISCHARGE_STOP)
						if err != nil {
							c.errChannel <- err
						}
					}
				}

				if wattsRequired > 0 {
					metrics.SetMetricValue("control", "action", map[string]string{"action": "pull_from_grid"}, 1)
					c.logger.WithFields(logrus.Fields{"wattsRequired": wattsRequired}).Info("Pulling watts from grid")
				} else {
					metrics.SetMetricValue("control", "action", map[string]string{"action": "pull_from_grid"}, 0)
				}
			}
		}
	}
}
