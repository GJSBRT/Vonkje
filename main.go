package main

import (
	"log"
	"fmt"
	"time"
	"syscall"
	"context"
	"os/signal"
	"net/http"

    "github.com/simonvetter/modbus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	// Strings
	pvVoltage = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "pv_voltage",
			Help: "The total amount of voltage",
		},
		[]string{
			"string",
		},
	)
	pvCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "pv_current",
			Help: "The total amount of current",
		},
		[]string{
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
			"phase",
		},
	)
	phaseCurrent = promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "phase_current",
			Help: "The total amount of current",
		},
		[]string{
			"phase",
		},
	)

	// other
	activePowerGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "active_power",
			Help: "The total amount of active power",
		},
	)
	reactivePowerGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "reactive_power",
			Help: "The total amount of reactive power",
		},
	)
	powerFactorGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "power_factor",
			Help: "The power factor",
		},
	)
	gridFrequencyGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "grid_frequency",
			Help: "The grid frequency",
		},
	)
	inverterEfficiencyGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "inverter_efficiency",
			Help: "The inverter efficiency",
		},
	)
	cabinetTemperatureGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "cabinet_temperature",
			Help: "The cabinet temperature",
		},
	)
	isulationResistanceGauge = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: "sun2000",
			Name: "isulation_resistance",
			Help: "The isulation resistance",
		},
	)
)

func main() {
	fmt.Println("wee")
	client, err := modbus.NewClient(&modbus.ClientConfiguration{
        URL:      	"rtuovertcp://81.172.220.193:10000",
		Speed: 		9600,
		DataBits: 	8,
		StopBits: 	1,
		Parity:   	modbus.PARITY_NONE,
        Timeout:  	5 * time.Second,
    })
	if err != nil {
		panic(err)
	}

	err = client.Open()
    if err != nil {
		panic(err)
    }

	log.Print("Connected to modbus tcp!")

	defer func() {
		client.Close()
		log.Print("Closed modbus tcp connection.")
	}()

	go func() {
		for {
			time.Sleep(30 * time.Second)

			updatemetrics(client)
		}
	}()

	server := &http.Server{
		Addr:    ":8080",
		Handler: nil,
	}

	go httpserv(server)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	<-ctx.Done()
	server.Shutdown(ctx)
}

func httpserv(server *http.Server) {
	log.Print("Running HTTP server on port 8080")
	http.Handle("/metrics", promhttp.Handler())
	server.ListenAndServe()
}

func updatemetrics(client *modbus.ModbusClient) {
	// string 1
    pv1Voltage, err := client.ReadRegister(MODBUS_PV1_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvVoltage.With(prometheus.Labels{"string": "1"}).Set(float64(pv1Voltage) / 10)

    pv1Current, err := client.ReadRegister(MODBUS_PV1_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvCurrent.With(prometheus.Labels{"string": "1"}).Set(float64(pv1Current) / 100)

	// string 2
    pv2Voltage, err := client.ReadRegister(MODBUS_PV2_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvVoltage.With(prometheus.Labels{"string": "2"}).Set(float64(pv2Voltage) / 10)

    pv2Current, err := client.ReadRegister(MODBUS_PV2_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvCurrent.With(prometheus.Labels{"string": "2"}).Set(float64(pv2Current) / 100)
	
	// string 3
    pv3Voltage, err := client.ReadRegister(MODBUS_PV3_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvVoltage.With(prometheus.Labels{"string": "3"}).Set(float64(pv3Voltage) / 10)

    pv3Current, err := client.ReadRegister(MODBUS_PV3_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvCurrent.With(prometheus.Labels{"string": "3"}).Set(float64(pv3Current) / 100)

	// string 4
    pv4Voltage, err := client.ReadRegister(MODBUS_PV4_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvVoltage.With(prometheus.Labels{"string": "4"}).Set(float64(pv4Voltage) / 10)

    pv4Current, err := client.ReadRegister(MODBUS_PV4_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	pvCurrent.With(prometheus.Labels{"string": "4"}).Set(float64(pv4Current) / 100)


	// phase A
	phaseAVoltage, err := client.ReadRegister(MODBUS_PHASE_A_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseVoltage.With(prometheus.Labels{"phase": "A"}).Set(float64(phaseAVoltage) / 10)

	phaseACurrent, err := client.ReadUint32(MODBUS_PHASE_A_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseCurrent.With(prometheus.Labels{"phase": "A"}).Set(float64(phaseACurrent) / 1000)

	// phase B
	phaseBVoltage, err := client.ReadRegister(MODBUS_PHASE_B_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseVoltage.With(prometheus.Labels{"phase": "B"}).Set(float64(phaseBVoltage) / 10)

	phaseBCurrent, err := client.ReadUint32(MODBUS_PHASE_B_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseCurrent.With(prometheus.Labels{"phase": "B"}).Set(float64(phaseBCurrent) / 1000)
	
	// phase C
	phaseCVoltage, err := client.ReadRegister(MODBUS_PHASE_C_VOLTAGE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseVoltage.With(prometheus.Labels{"phase": "C"}).Set(float64(phaseCVoltage) / 10)

	phaseCCurrent, err := client.ReadUint32(MODBUS_PHASE_C_CURRENT, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	phaseCurrent.With(prometheus.Labels{"phase": "C"}).Set(float64(phaseCCurrent) / 1000)


	// other
	activePower, err := client.ReadUint32(MODBUS_ACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	activePowerGauge.Set(float64(activePower) / 1000)

	reactivePower, err := client.ReadUint32(MODBUS_REACTIVE_POWER, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	reactivePowerGauge.Set(float64(reactivePower) / 1000)
	
	powerFactor, err := client.ReadRegister(MODBUS_POWER_FACTOR, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	powerFactorGauge.Set(float64(powerFactor) / 1000)
	
	gridFrequency, err := client.ReadRegister(MODBUS_FREQUENCY, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	gridFrequencyGauge.Set(float64(gridFrequency) / 100)
	
	inverterEfficiency, err := client.ReadRegister(MODBUS_FREQUENCY, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	inverterEfficiencyGauge.Set(float64(inverterEfficiency) / 100)
	
	cabinetTemperature, err := client.ReadRegister(MODBUS_CABINET_TEMPERATURE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	cabinetTemperatureGauge.Set(float64(cabinetTemperature) / 10)
	
	isulationResistance, err := client.ReadRegister(MODBUS_INSULATION_RESISTANCE, modbus.HOLDING_REGISTER)
	if err != nil {
		panic(err)
	}
	isulationResistanceGauge.Set(float64(isulationResistance) / 10)
}
