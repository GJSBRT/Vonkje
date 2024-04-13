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
	MinimumBatteryCapacity int `mapstructure:"minimum-battery-capacity"`
	BatteryChargePercentage int `mapstructure:"battery-charge-percentage"`
	BatteryDischargePercentage int `mapstructure:"battery-discharge-percentage"`
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

			// 1. Get power meter active power
			powerMeterActivePower, err := metrics.GetMetricLastEntrySum("power_meter", "active_power")
			if err != nil {
				c.errChannel <- err
				continue
			}
			var overUsage int
			var overProduction int
			if powerMeterActivePower < 0 {
				overUsage = int(math.Abs(powerMeterActivePower))
			} else {
				overProduction = int(powerMeterActivePower)
			}

			// 2. Get current solar production sum
			avgSolarIn, err := metrics.GetMetricLastEntrySum("sun2000", "input_power")
			if err != nil {
				c.errChannel <- err
				continue
			}
			avgSolarIn = math.Floor(avgSolarIn * 1000) // convert to watts

			c.logger.WithFields(logrus.Fields{"avgSolarIn": avgSolarIn, "powerMeterActivePower": powerMeterActivePower}).Info("Solar production and home load")

			// Get battery (dis)charge watts
			batteryWatts, err := metrics.GetMetricLastEntrySum("luna2000", "charging_status")
			if err != nil {
				c.errChannel <- err
				continue
			}
			var batteryChargeWatts int
			var batteryDischargeWatts int
			if batteryWatts > 0 {
				batteryChargeWatts = int(batteryWatts)
			} else {
				batteryDischargeWatts = int(math.Abs(batteryWatts))
			}

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
			var batteriesWithCapacity int
			var batteriesFull int
			for _, battery := range batteries {
				if battery.capacity > float64(c.config.MinimumBatteryCapacity) {
					batteriesWithCapacity++
				}

				if battery.capacity == 100 {
					batteriesFull++
				}
			}

			c.logger.WithFields(logrus.Fields{"batteriesWithCapacity": batteriesWithCapacity, "batteriesFull": batteriesFull}).Info("Battery capacities")

			// Are we pulling power from the grid?
			if overUsage > 0 {
				c.logger.WithFields(logrus.Fields{"overUsage": overUsage}).Info("Overusage detected")

				// If batteries are charging with more than the overusage we should dial back the charging power
				if batteryChargeWatts > overUsage {
					newChargeWatts := int(float64(batteryChargeWatts - overUsage) * (float64(c.config.BatteryChargePercentage) / 100)) // Add buffer to prevent charging to much
					newChargeWattsPerBattery := newChargeWatts / len(batteries)

					for _, battery := range batteries {
						err := c.modbus.ChangeBatteryForceCharge(battery.inverter, battery.battery, modbus.MODBUS_STATE_BATTERY_FORCIBLE_CHARGE_DISCHARGE_CHARGE, uint(newChargeWattsPerBattery))
						if err != nil {
							c.errChannel <- err
						}
					}

					c.logger.WithFields(logrus.Fields{"overUsage": overUsage, "batteryChargeWatts": batteryChargeWatts, "newChargeWatts": newChargeWatts}).Info("Dialing back battery charge")
					metrics.SetMetricValue("control", "action", map[string]string{"action": "charge_batteries"}, 1)
					
					// Overusage has been compensated. No further actions is required.
					continue
				}

				// If no batteries have capacity, pull from grid :(
				if batteriesWithCapacity == 0 {
					metrics.SetMetricValue("control", "action", map[string]string{"action": "pull_from_grid"}, 1)
					c.logger.WithFields(logrus.Fields{"overUsage": overUsage}).Info("Pulling watts from grid")
					continue
				}

				newDischargeWatts := int(float64(batteryDischargeWatts - overUsage) * (float64(c.config.BatteryDischargePercentage) / 100))
				newDischargeWattsPerBattery := newDischargeWatts / batteriesWithCapacity

				for _, battery := range batteries {
					err := c.modbus.ChangeBatteryForceCharge(battery.inverter, battery.battery, modbus.MODBUS_STATE_BATTERY_FORCIBLE_CHARGE_DISCHARGE_DISCHARGE, uint(newDischargeWattsPerBattery))
					if err != nil {
						c.errChannel <- err
					}
				}

				c.logger.WithFields(logrus.Fields{"overUsage": overUsage, "batteryDischargeWatts": batteryDischargeWatts, "newDischargeWatts": newDischargeWatts}).Info("Discharging batteries to compensate for overusage")
				metrics.SetMetricValue("control", "action", map[string]string{"action": "discharge_battery"}, 1)

				// Overusage has been compensated. No further actions is required.
				continue
			} else {
				availableWatts := int(float64(overProduction) * (float64(c.config.BatteryChargePercentage) / 100))
				c.logger.WithFields(logrus.Fields{"overProduction":overProduction, "availableWatts": availableWatts}).Info("Overproduction and available watts for charging batteries")

				if availableWatts > 0 {
					if (len(batteries) - batteriesFull) == 0 {
						c.logger.WithFields(logrus.Fields{"availableWatts": availableWatts}).Info("No batteries available to charge")
						continue
					}

					perBatteryChargeWatts := availableWatts / (len(batteries) - batteriesFull)
					for _, battery := range batteries {
						if battery.capacity < 100 {
							err := c.modbus.ChangeBatteryForceCharge(battery.inverter, battery.battery, modbus.MODBUS_STATE_BATTERY_FORCIBLE_CHARGE_DISCHARGE_CHARGE, uint(perBatteryChargeWatts))
							if err != nil {
								c.errChannel <- err
							}
						}
					}

					c.logger.WithFields(logrus.Fields{"availableWatts": availableWatts}).Info("Charging batteries")
					metrics.SetMetricValue("control", "action", map[string]string{"action": "charge_batteries"}, 1)
					continue
				}
			}
		}
	}
}
