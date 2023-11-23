package main

import (
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"go.uber.org/zap"
)

type PQA_Message struct {
	DeviceID    string  `json:"device_id"`
	SlaveID     string  `json:"slave_id"`
	MeasureTime string  `json:"measure_time"`
	Date        string  `json:"full_date"`
	ModbusAdd   string  `json:"modbus_add"`
	Value       float64 `json:"value"`
	Key         string  `json:"var_name"`
}

var connectHandler = func(logger *zap.Logger, deviceID string) MQTT.OnConnectHandler {
	return func(client MQTT.Client) {
		logger.Info("Connected to broker", zap.String("deviceID:", deviceID))
	}
}

var connectLostHandler = func(logger *zap.Logger, deviceID string) MQTT.ConnectionLostHandler {
	return func(client MQTT.Client, err error) {
		//TODO Check connection lost handling
		logger.Info("Connection to broker lost", zap.String("deviceID:", deviceID))
	}
}
