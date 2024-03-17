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
	MinimumSolarOverProduction int `mapstructure:"minimum-solar-over-production"`
	OverDischargePercentage int `mapstructure:"over-discharge-percentage"`
	MinimumBatteryCapacity int `mapstructure:"minimum-battery-capacity"`
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

	c.logger.Infof("Waiting %d seconds before starting control loop to collect metrics", viper.GetInt("modbus.read-metrics-interval"))
	time.Sleep(time.Duration(viper.GetInt("modbus.read-metrics-interval")) * time.Second)

	c.logger.Info("Starting control loop")

	ticker := time.NewTicker(time.Duration(viper.GetInt("modbus.read-metrics-interval")) * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-c.ctx.Done():
			c.logger.Info("Stopping control loop")
			return
		case <-ticker.C:
			c.logger.Debug("Control loop tick")
			metrics.SetMetricValue("control", "action", map[string]string{"action": "charge_batteries"}, 0)
			metrics.SetMetricValue("control", "action", map[string]string{"action": "discharge_battery"}, 0)
			metrics.SetMetricValue("control", "action", map[string]string{"action": "pull_from_grid"}, 0)

			// 1. Get current home energy consumption
			inverterInputPower, err := metrics.GetMetricLastEntrySum("sun2000", "active_power")
			if err != nil {
				c.errChannel <- err
				continue
			}

			powerMeterActivePower, err := metrics.GetMetricLastEntrySum("power_meter", "active_power")
			if err != nil {
				c.errChannel <- err
				continue
			}

			avgHomeLoad := (inverterInputPower * 1000) - powerMeterActivePower
			if avgHomeLoad < 0 {
				avgHomeLoad = 0
			}

			// 2. Get current solar production
			avgSolarIn, err := metrics.GetMetricLastEntryAverage("sun2000", "input_power")
			if err != nil {
				c.errChannel <- err
				continue
			}
			avgSolarIn = math.Floor(avgSolarIn * 1000)
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
			var overProductionWatts int
			if avgSolarIn > avgHomeLoad {
				overProduction = math.Ceil((avgSolarIn - avgHomeLoad) / avgSolarIn * 100)
				overProductionWatts = int(math.Floor(avgSolarIn - avgHomeLoad))
			} else {
				overProduction = 0
				overProductionWatts = 0
			}
			metrics.SetMetricValue("control", "over_production", map[string]string{}, overProduction)
			c.logger.WithFields(logrus.Fields{"percentage": overProduction, "watts": overProductionWatts}).Info("Over production")

			// 4. if solar over production is more than x%, charge battery
			if overProduction > float64(c.config.MinimumSolarOverProduction) {
				metrics.SetMetricValue("control", "action", map[string]string{"action": "charge_batteries"}, 1)

				// charge batteries with 20% less than over production
				batteryChargeWatts := uint(math.Floor(float64(overProductionWatts) * 0.80))

				for _, battery := range batteries {
					if battery.capacity < 100 {
						c.logger.WithFields(logrus.Fields{"inverter": battery.inverter, "battery": battery.battery, "capacity": battery.capacity, "watts": batteryChargeWatts}).Info("Battery is not fully charged, starting charge")
						err := c.modbus.ChangeBatteryForceCharge(battery.inverter, battery.battery, modbus.MODBUS_STATE_BATTERY_FORCIBLE_CHARGE_DISCHARGE_CHARGE, batteryChargeWatts)
						if err != nil {
							c.errChannel <- err
							continue
						}
					} else {
						c.logger.WithFields(logrus.Fields{"inverter": battery.inverter, "battery": battery.battery, "watts": batteryChargeWatts}).Info("Battery is fully charged, stopping charge")
						err := c.modbus.ChangeBatteryForceCharge(battery.inverter, battery.battery, modbus.MODBUS_STATE_BATTERY_FORCIBLE_CHARGE_DISCHARGE_STOP, 0)
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
			wattsRequired := math.Ceil(avgHomeLoad - avgSolarIn) * ((100 + float64(c.config.OverDischargePercentage)) / 100)
			if avgSolarIn < avgHomeLoad {
				if wattsRequired > 0 {
					metrics.SetMetricValue("control", "action", map[string]string{"action": "discharge_battery"}, 1)
				} else {
					metrics.SetMetricValue("control", "action", map[string]string{"action": "discharge_battery"}, 0)
				}

				for _, battery := range batteries {
					if wattsRequired > 0 {
						if battery.capacity > float64(c.config.MinimumBatteryCapacity) {
							var useWatts uint
							if wattsRequired > 5000 {
								useWatts = 5000
								wattsRequired -= 5000
							} else {
								useWatts = uint(wattsRequired)
								wattsRequired = 0
							}

							c.logger.WithFields(logrus.Fields{"inverter": battery.inverter, "battery": battery.battery, "capacity": battery.capacity, "watts": useWatts}).Info("Discharging battery")

							err := c.modbus.ChangeBatteryForceCharge(battery.inverter, battery.battery, modbus.MODBUS_STATE_BATTERY_FORCIBLE_CHARGE_DISCHARGE_DISCHARGE, useWatts)
							if err != nil {
								c.errChannel <- err
							}
						}
					} else {
						c.logger.WithFields(logrus.Fields{"inverter": battery.inverter, "battery": battery.battery}).Info("Battery is not required, stopping discharge")

						err := c.modbus.ChangeBatteryForceCharge(battery.inverter, battery.battery, modbus.MODBUS_STATE_BATTERY_FORCIBLE_CHARGE_DISCHARGE_STOP, 0)
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
