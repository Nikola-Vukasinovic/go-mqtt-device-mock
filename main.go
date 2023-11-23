package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strconv"
	"sync"
	"time"

	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func main() {
	//Runtime vars
	var wg sync.WaitGroup
	var logger *zap.Logger
	var err error
	//Create context
	/*ctx, cancel := context.WithCancel(context.Background())
	defer cancel()*/
	// Configure the logger with a timestamp
	logConfig := zap.NewProductionConfig()
	logConfig.EncoderConfig.TimeKey = "time"
	logConfig.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder

	logger, err = logConfig.Build()

	if err != nil {
		panic("Failed to create logger")
	}
	defer logger.Sync() // Flush any buffered log entries

	//Env vars
	var broker string
	var topic string
	var port string
	var qos int
	var numDevices int
	var sendPeriod int

	//TODO: Add check is local or k8 and adjust broker address for dev/test/stage
	_, inCluster := os.LookupEnv("KUBERNETES_SERVICE_HOST")

	if inCluster {
		logger.Info("Running inside k8 cluster")
		broker = os.Getenv("MQTT_BROKER_HOST")
		port = os.Getenv("MQTT_BROKER_PORT")
		qos, _ = strconv.Atoi(os.Getenv("MQTT_QOS"))
		// Read the topic name from the environment variable
		topic = os.Getenv("MQTT_BROKER_SUB_TOPIC")
		if topic == "" {
			logger.Warn("MQTT_BROKER_SUB_TOPIC environment variable is not set or empty. Using default topic devices/telemetry.")
			topic = "devices/telemetry"
		}
		numDevices, _ = strconv.Atoi(os.Getenv("DEVICES_NUM"))
		sendPeriod, _ = strconv.Atoi(os.Getenv("SEND_PERIOD"))
	} else {
		if err := godotenv.Load(); err != nil {
			log.Fatalf("Error loading .env file: %v", err)
		}
		logger.Info("Running outside k8 cluster")
		broker = os.Getenv("MQTT_BROKER_HOST")
		port = os.Getenv("MQTT_BROKER_PORT")
		qos, _ = strconv.Atoi(os.Getenv("MQTT_QOS"))
		topic = "devices/telemetry"
		numDevices, _ = strconv.Atoi(os.Getenv("DEVICES_NUM"))
		sendPeriod, _ = strconv.Atoi(os.Getenv("SEND_PERIOD"))
	}

	//Simulate for DEVICES_NUM number of devices sending messages every SEND_PERIOD seconds
	for i := 0; i < numDevices; i++ {
		wg.Add(1)
		// Start a goroutine to handle the MQTT messages
		deviceID := fmt.Sprintf("mock_device_%d", i)
		go func(deviceID string) {
			defer wg.Done()
			opts := MQTT.NewClientOptions()
			opts.AddBroker(fmt.Sprintf("mqtt://%s:%s", broker, port))
			opts.SetClientID(deviceID)
			opts.OnConnect = connectHandler(logger, deviceID)
			opts.OnConnectionLost = connectLostHandler(logger, deviceID)
			client := MQTT.NewClient(opts)

			if token := client.Connect(); token.Wait() && token.Error() != nil {
				panic(token.Error())
			}
			//Start sending messages with SEND_PERIOD in seconds
			for {
				//Create JSON message for sending
				msg := pqaMultipleMessages(logger, deviceID)

				jsonMessage, err := json.Marshal(msg)
				if err != nil {
					fmt.Println("Error:", err)
					return
				}

				token := client.Publish(topic, byte(qos), false, jsonMessage)
				token.Wait()
				logger.Info("Sent message", zap.String("deviceID:", deviceID))
				time.Sleep(time.Duration(sendPeriod) * time.Second) // Simulate delay between messages
			}

		}(deviceID)

	}
	wg.Wait()
}
