package modbus

type RegisterType uint8

type Register struct {
	Namespace 	string
	Name 		string
	Fields 		map[string]string
	Address 	uint16
	Unit   		string
	Gain   		float64
	Quantity 	uint16
	Type   		RegisterType
	Writeable 	bool
}

const (
	RegisterTypeUint16 RegisterType = iota
	RegisterTypeUint32
	RegisterTypeInt16
	RegisterTypeInt32
)

const (
	MODBUS_STATE_BATTERY_1_FORCIBLE_CHARGE_DISCHARGE_STOP uint16 = iota
	MODBUS_STATE_BATTERY_1_FORCIBLE_CHARGE_DISCHARGE_CHARGE
	MODBUS_STATE_BATTERY_1_FORCIBLE_CHARGE_DISCHARGE_DISCHARGE
)

var (
	sun2000Registers = map[string]Register{
		"pv_voltage_string_1": 		Register{Namespace: "sun2000",	Name: "pv_voltage",				Fields: map[string]string{"string": "1"},	Address: 32016, 	Unit: "V", 		Gain: 10, 		Quantity: 1, 	Type: RegisterTypeInt16, 	Writeable: false},
		"pv_current_string_1": 		Register{Namespace: "sun2000",	Name: "pv_current",				Fields: map[string]string{"string": "1"},	Address: 32017, 	Unit: "A", 		Gain: 100, 		Quantity: 1, 	Type: RegisterTypeInt16, 	Writeable: false},
		"pv_voltage_string_2": 		Register{Namespace: "sun2000",	Name: "pv_voltage",				Fields: map[string]string{"string": "2"},	Address: 32018, 	Unit: "V", 		Gain: 10, 		Quantity: 1, 	Type: RegisterTypeInt16, 	Writeable: false},
		"pv_current_string_2": 		Register{Namespace: "sun2000",	Name: "pv_current",				Fields: map[string]string{"string": "2"},	Address: 32019, 	Unit: "A",		Gain: 100,		Quantity: 1, 	Type: RegisterTypeInt16,	Writeable: false},
		"device_status": 			Register{Namespace: "sun2000",	Name: "device_status",			Fields: map[string]string{},				Address: 32089,		Unit: "",		Gain: 1,		Quantity: 1,	Type: RegisterTypeUint16,	Writeable: false},
		"input_power": 				Register{Namespace: "sun2000",	Name: "input_power",			Fields: map[string]string{},				Address: 32064,		Unit: "kW",		Gain: 1000,		Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"phase_voltage_phase_a": 	Register{Namespace: "sun2000",	Name: "phase_voltage",			Fields: map[string]string{"phase": "A"},	Address: 32069,		Unit: "V",		Gain: 10, 		Quantity: 1,	Type: RegisterTypeUint16,	Writeable: false},
		"phase_voltage_phase_b": 	Register{Namespace: "sun2000",	Name: "phase_voltage",			Fields: map[string]string{"phase": "B"},	Address: 32070,		Unit: "V",		Gain: 10, 		Quantity: 1,	Type: RegisterTypeUint16,	Writeable: false},
		"phase_voltage_phase_c": 	Register{Namespace: "sun2000",	Name: "phase_voltage",			Fields: map[string]string{"phase": "C"},	Address: 32071,		Unit: "V",		Gain: 10, 		Quantity: 1,	Type: RegisterTypeUint16,	Writeable: false},
		"phase_current_phase_a": 	Register{Namespace: "sun2000",	Name: "phase_current",			Fields: map[string]string{"phase": "A"},	Address: 32072,		Unit: "A",		Gain: 1000, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"phase_current_phase_b": 	Register{Namespace: "sun2000",	Name: "phase_current",			Fields: map[string]string{"phase": "B"},	Address: 32074,		Unit: "A",		Gain: 1000, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"phase_current_phase_c": 	Register{Namespace: "sun2000",	Name: "phase_current",			Fields: map[string]string{"phase": "C"},	Address: 32076,		Unit: "A",		Gain: 1000, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"active_power": 			Register{Namespace: "sun2000",	Name: "active_power",			Fields: map[string]string{},				Address: 32080,		Unit: "kW",		Gain: 1000, 	Quantity: 2,	Type: RegisterTypeInt32, 	Writeable: false},
		"reactive_power": 			Register{Namespace: "sun2000",	Name: "reactive_power",			Fields: map[string]string{},				Address: 32082,		Unit: "kVar",	Gain: 1000,		Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"power_factor": 			Register{Namespace: "sun2000",	Name: "power_factor",			Fields: map[string]string{},				Address: 32084, 	Unit: "",		Gain: 1000,		Quantity: 1,	Type: RegisterTypeInt16,	Writeable: false},
		"grid_frequency": 			Register{Namespace: "sun2000",	Name: "grid_frequency",			Fields: map[string]string{},				Address: 32085, 	Unit: "Hz",		Gain: 100,		Quantity: 100,	Type: RegisterTypeUint16,	Writeable: false},
		"inverter_efficiency": 		Register{Namespace: "sun2000",	Name: "inverter_efficiency",	Fields: map[string]string{},				Address: 32086, 	Unit: "%",		Gain: 100,		Quantity: 1,	Type: RegisterTypeUint16,	Writeable: false},
		"cabinet_temperature": 		Register{Namespace: "sun2000",	Name: "cabinet_temperature",	Fields: map[string]string{},				Address: 32087, 	Unit: "°C",		Gain: 10,		Quantity: 1,	Type: RegisterTypeInt16,	Writeable: false},
		"isulation_resistance": 	Register{Namespace: "sun2000",	Name: "isulation_resistance",	Fields: map[string]string{},				Address: 32088, 	Unit: "MΩ",		Gain: 1000,		Quantity: 1,	Type: RegisterTypeUint16,	Writeable: false},	
	}

	luna2000Registers = map[string]Register{
		"running_status_battery_1": 			Register{Namespace: "luna2000",	Name: "running_status",				Fields: map[string]string{"battery": "1"},	Address: 37000, Unit: "",		Gain: 1,	Quantity: 1,	Type: RegisterTypeUint16,	Writeable: false},
		"charging_status_battery_1": 			Register{Namespace: "luna2000",	Name: "charging_status",			Fields: map[string]string{"battery": "1"},	Address: 37001,	Unit: "W",		Gain: 1,	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"bus_voltage_battery_1": 				Register{Namespace: "luna2000",	Name: "bus_voltage",				Fields: map[string]string{"battery": "1"},	Address: 37003,	Unit: "V",		Gain: 10,	Quantity: 1,	Type: RegisterTypeUint16,	Writeable: false},
		"battery_capacity_battery_1": 			Register{Namespace: "luna2000",	Name: "battery_capacity",			Fields: map[string]string{"battery": "1"},	Address: 37004,	Unit: "%",		Gain: 10,	Quantity: 1,	Type: RegisterTypeUint16,	Writeable: false},
		"total_charge_battery_1": 				Register{Namespace: "luna2000",	Name: "total_charge",				Fields: map[string]string{"battery": "1"},	Address: 37066,	Unit: "kWh",	Gain: 100,	Quantity: 2,	Type: RegisterTypeUint32,	Writeable: false},
		"total_discharge_battery_1": 			Register{Namespace: "luna2000",	Name: "total_discharge",			Fields: map[string]string{"battery": "1"},	Address: 37068,	Unit: "kWh",	Gain: 100,	Quantity: 2,	Type: RegisterTypeUint32,	Writeable: false},
		"forcible_charge_discharge_battery_1": 	Register{Namespace: "luna2000",	Name: "forcible_charge_discharge",	Fields: map[string]string{"battery": "1"},	Address: 47100,	Unit: "",		Gain: 1,	Quantity: 1,	Type: RegisterTypeUint16,	Writeable: true},
		"forcible_charge_power_battery_1": 		Register{Namespace: "luna2000",	Name: "forcible_charge_power",		Fields: map[string]string{"battery": "1"},	Address: 47247,	Unit: "kW",		Gain: 1000,	Quantity: 2,	Type: RegisterTypeUint32,	Writeable: true},
		"forcible_discharge_power_battery_1": 	Register{Namespace: "luna2000",	Name: "forcible_discharge_power",	Fields: map[string]string{"battery": "1"},	Address: 47249,	Unit: "kW",		Gain: 1000,	Quantity: 2,	Type: RegisterTypeUint32,	Writeable: true},
	}

	powerMeterRegisters = map[string]Register{
		"status": 						Register{Namespace: "power_meter",	Name: "status",							Fields: map[string]string{},				Address: 37100,	Unit: "",		Gain: 1,	Quantity: 1,	Type: RegisterTypeUint16,	Writeable: false},
		"phase_voltage_phase_a": 		Register{Namespace: "power_meter",	Name: "phase_voltage",					Fields: map[string]string{"phase": "A"},	Address: 37101,	Unit: "V",		Gain: 10, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"phase_voltage_phase_b": 		Register{Namespace: "power_meter",	Name: "phase_voltage",					Fields: map[string]string{"phase": "B"},	Address: 37103,	Unit: "V",		Gain: 10, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"phase_voltage_phase_c": 		Register{Namespace: "power_meter",	Name: "phase_voltage",					Fields: map[string]string{"phase": "C"},	Address: 37105,	Unit: "V",		Gain: 10, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"phase_current_phase_a": 		Register{Namespace: "power_meter",	Name: "phase_current",					Fields: map[string]string{"phase": "A"},	Address: 37107,	Unit: "A",		Gain: 100, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"phase_current_phase_b": 		Register{Namespace: "power_meter",	Name: "phase_current",					Fields: map[string]string{"phase": "B"},	Address: 37109,	Unit: "A",		Gain: 100, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"phase_current_phase_c": 		Register{Namespace: "power_meter",	Name: "phase_current",					Fields: map[string]string{"phase": "C"},	Address: 37111,	Unit: "A",		Gain: 100, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"active_power": 				Register{Namespace: "power_meter",	Name: "active_power",					Fields: map[string]string{},				Address: 37113, Unit: "W",		Gain: 1, 	Quantity: 2, 	Type: RegisterTypeInt32, 	Writeable: false},
		"reactive_power": 				Register{Namespace: "power_meter",	Name: "reactive_power",					Fields: map[string]string{},				Address: 37115, Unit: "Var",	Gain: 1, 	Quantity: 2, 	Type: RegisterTypeInt32,	Writeable: false},
		"power_factor": 				Register{Namespace: "power_meter",	Name: "power_factor",					Fields: map[string]string{},				Address: 37117, Unit: "",		Gain: 1000,	Quantity: 1,	Type: RegisterTypeInt16, 	Writeable: false},
		"frequency": 					Register{Namespace: "power_meter",	Name: "frequency",						Fields: map[string]string{},				Address: 37118,	Unit: "Hz",		Gain: 100,	Quantity: 1,	Type: RegisterTypeInt16,	Writeable: false},
		"positive_active_electricity": 	Register{Namespace: "power_meter",	Name: "positive_active_electricity",	Fields: map[string]string{},				Address: 37119,	Unit: "kWh",	Gain: 100, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"reverse_active_power": 		Register{Namespace: "power_meter",	Name: "reverse_active_power",			Fields: map[string]string{},				Address: 37121, Unit: "kWh", 	Gain: 100,	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"accumulated_reactive_power": 	Register{Namespace: "power_meter",	Name: "accumulated_reactive_power",		Fields: map[string]string{},				Address: 37123,	Unit: "kVarh",	Gain: 100,	Quantity: 2,	Type: RegisterTypeInt32, 	Writeable: false},
		"line_voltage_line_ab": 		Register{Namespace: "power_meter",	Name: "line_voltage",					Fields: map[string]string{"line": "AB"},	Address: 37126,	Unit: "V", 		Gain: 10, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"line_voltage_line_bc": 		Register{Namespace: "power_meter",	Name: "line_voltage",					Fields: map[string]string{"line": "BC"},	Address: 37128,	Unit: "V", 		Gain: 10, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"line_voltage_line_ca": 		Register{Namespace: "power_meter",	Name: "line_voltage",					Fields: map[string]string{"line": "CA"},	Address: 37130,	Unit: "V", 		Gain: 10, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"phase_active_power_phase_a": 	Register{Namespace: "power_meter",	Name: "phase_active_power",				Fields: map[string]string{"phase": "A"},	Address: 37132,	Unit: "W", 		Gain: 1, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"phase_active_power_phase_b": 	Register{Namespace: "power_meter",	Name: "phase_active_power",				Fields: map[string]string{"phase": "B"},	Address: 37134,	Unit: "W", 		Gain: 1, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
		"phase_active_power_phase_c": 	Register{Namespace: "power_meter",	Name: "phase_active_power",				Fields: map[string]string{"phase": "C"},	Address: 37136,	Unit: "W", 		Gain: 1, 	Quantity: 2,	Type: RegisterTypeInt32,	Writeable: false},
	}
)
