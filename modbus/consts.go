package modbus

// uint16 -> ReadRegister
// int32 -> ReadUint32

const (
	// Sun2000 Inverter
	MODBUS_INVERTER_PV1_VOLTAGE = 32016	// int16 readonly
	MODBUS_INVERTER_PV1_CURRENT = 32017	// int16 readonly
	MODBUS_INVERTER_PV2_VOLTAGE = 32018	// int16 readonly
	MODBUS_INVERTER_PV2_CURRENT = 32019	// int16 readonly

	MODBUS_INVERTER_STATE_1 = 32000 // uint16 readonly
	MODBUS_INVERTER_DEVICE_STATUS = 32089 // uint16 readonly
	MODBUS_INVERTER_INPUT_POWER = 32064	// int32 readonly

	MODBUS_INVERTER_PHASE_A_VOLTAGE = 32069 // uint16 readonly
	MODBUS_INVERTER_PHASE_B_VOLTAGE = 32070 // uint16 readonly
	MODBUS_INVERTER_PHASE_C_VOLTAGE = 32071 // uint16 readonly
	MODBUS_INVERTER_PHASE_A_CURRENT = 32072	// int32 readonly
	MODBUS_INVERTER_PHASE_B_CURRENT = 32074 // int32 readonly
	MODBUS_INVERTER_PHASE_C_CURRENT = 32076	// int32 readonly
	MODBUS_INVERTER_ACTIVE_POWER = 32080	// int32 readonly
	MODBUS_INVERTER_REACTIVE_POWER = 32082	// int32 readonly
	MODBUS_INVERTER_POWER_FACTOR = 32084	// int16 readonly
	MODBUS_INVERTER_FREQUENCY = 32085 // uint16 readonly
	MODBUS_INVERTER_INVERTER_EFFICIENCY = 32086	// uint16 readonly
	MODBUS_INVERTER_CABINET_TEMPERATURE = 32087	// int16 readonly
	MODBUS_INVERTER_INSULATION_RESISTANCE = 32088 // uint16 readonly

	// Luna Battery
	MODBUS_BATTERY_1_RUNNING_STATUS = 37000 // uint16 readonly - 0: offline, 1: standby, 2: running, 3: fault, 4: sleep
	MODBUS_BATTERY_1_CHARGING_STATUS = 37001 // int32  readonly - >0: charging, <0: discharging
	MODBUS_BATTERY_1_BUS_VOLTAGE = 37003 // uint16 readonly
	MODBUS_BATTERY_1_CAPACITY = 37004 // uint16 readonly
	MODBUS_BATTERY_1_TOTAL_CHARGE = 37066 // int32 readonly
	MODBUS_BATTERY_1_TOTAL_DISCHARGE = 37068 // int32 readonly

	// Power meter
	MODBUS_POWER_METER_STATUS = 37100 // uint16 readonly - 0: offline, 1: normal
	MODBUS_POWER_METER_PHASE_A_VOLTAGE = 37101 // int32 readonly
	MODBUS_POWER_METER_PHASE_B_VOLTAGE = 37103 // int32 readonly
	MODBUS_POWER_METER_PHASE_C_VOLTAGE = 37105 // int32 readonly
	MODBUS_POWER_METER_PHASE_A_CURRENT = 37107 // int32 readonly
	MODBUS_POWER_METER_PHASE_B_CURRENT = 37109 // int32 readonly
	MODBUS_POWER_METER_PHASE_C_CURRENT = 37111 // int32 readonly
	MODBUS_POWER_METER_ACTIVE_POWER = 37113 // int32 readonly
	MODBUS_POWER_METER_REACTIVE_POWER = 37115 // int32 readonly
	MODBUS_POWER_METER_POWER_FACTOR = 37117 // int16 readonly
	MODBUS_POWER_METER_FREQUENCY = 37118 // uint16 readonly
	MODBUS_POWER_METER_POSITIVE_ACTIVE_ELECTRICITY = 37119 // int32 readonly
	MODBUS_POWER_METER_REVERSE_ACTIVE_POWER = 37121 // int32 readonly
	MODBUS_POWER_METER_ACCUMULATED_REACTIVE_POWER = 37123 // int32 readonly
	MODBUS_POWER_METER_AB_LINE_VOLTAGE = 37126 // int32 readonly
	MODBUS_POWER_METER_BC_LINE_VOLTAGE = 37128 // int32 readonly
	MODBUS_POWER_METER_CA_LINE_VOLTAGE = 37130 // int32 readonly
	MODBUS_POWER_METER_PHASE_A_ACTIVE_POWER = 37132 // int32 readonly
	MODBUS_POWER_METER_PHASE_B_ACTIVE_POWER = 37134 // int32 readonly
	MODBUS_POWER_METER_PHASE_C_ACTIVE_POWER = 37136 // int32 readonly
	MODBUS_POWER_METER_MODEL_RESULT = 37138 // uint16 readonly
)
