package main

import (
	"fmt"
	"math/rand"
	"time"

	"go.uber.org/zap"
)

// Function that creates specific measurements as per PQA device
func pqaMultipleMessages(logger *zap.Logger, deviceID string) []PQA_Message {
	//Measurements are based on values that are monitored on device and defined in keys
	//Keys mapping are listed below
	keys := []string{"U_L1",
		"U_L2",
		"U_L3",
		"P_L1",
		"P_L2",
		"P_L3",
		"P_L2",
		"S_Tot",
		"P_Tot",
		"Q_Tot",
		"E_P_Imp",
		"E_Q_Imp",
		"PF1",
		"PF2",
		"PF3",
		"Freq_avg"}

	msgs := []PQA_Message{}
	for _, key := range keys {

		msg := PQA_Message{
			DeviceID:    deviceID,
			SlaveID:     "1",
			MeasureTime: fmt.Sprint(time.Now().Unix()),
			Date:        string(time.Now().Format("02/01/2006 15:04:05")),
			Key:         key,
			ModbusAdd:   fmt.Sprint(rand.Intn(1001) + 3000),
			Value:       rand.Float64() * 100,
		}
		msgs = append(msgs, msg)
	}

	return msgs
}
