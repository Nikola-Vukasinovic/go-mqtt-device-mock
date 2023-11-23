# Simple Go app for simulating multiple IoT telemetry devices

## Motivation

Little round trip for the motivation for this app. When I started working on my own services that will be used in IoT projects mainly using Kubernetes clusters biggest issue was not on the device side since we already had a lot of devices deployed in the field and a lot of real-life hardware somewhere on the table sending data as part of their testing. As the pace picks up in the project and when you really get in the domain where you ask yourself how will your broker perform at tens of thousands of messages and as a matter a fact whole backend infrastructure these kind of testing apps are more than welcome.

Go proves itself as an ideal candidate for these kinds of tasks involving concurrency.

## Prerequisites

Go version 1.20 or later

App relies on Eclipse Paho MQTT Go client [MQTT](https://pkg.go.dev/github.com/eclipse/paho.mqtt.golang) for communication with MQTT broker. For the testing purposes this implementation of v3.1/v3.1.1. performed for now rather nicely with high performance observed when using concurency on my local development machine.

## Environment

Environment variables that **need to be set** with env or .sh script are:

```
MQTT_BROKER_HOST
MQTT_BROKER_PORT
MQTT_BROKER_SUB_TOPIC
MQTT_QOS
DEVICES_NUM
SEND_PERIOD
```

My example vars:

env.sh

```
MQTT_BROKER_HOST="localhost"
MQTT_BROKER_PORT="1883"
MQTT_BROKER_SUB_TOPIC="devices/telemetry"
MQTT_QOS=0
DEVICES_NUM=50
SEND_PERIOD=10
```

Depending on the scenario sometimes I prefer to use .env file but with go lately i use shell scripts that I standardize across the infrastructure stack.

Rember to run:

```
export env.sh
```

## Sidenote

You might notice that *MQTT_BROKER_HOST* is set to *localhost*. This is because I used this app in local development to test development IoT infrastructure in my local Kubernetes cluster. Or to be more honest k3s cluster instead of K8 cluster. Very interesting revelation and suggestion to have a [look](https://k3s.io/) for local development.

Another sidenote is that remember when using clusters on local development machine you will have to deal with port forwarding and in case kubectl it is something along the lines:

```
kubectl port-forward mqtt-broker-id port:port
```

This will enable to forward traffic sent to ports 1883 or 8883 to be redirected to cluster port. 

## Running the example

Since this is running example of simple app allow that certain module versions issues could arise so best practise would be to do:

```
go mod tidy
```

before 

```
go build .
```

Should the errors persist you can try to clean module cache:

```
go clean -modcache
go clean -cache
go clean -i
```

## **Simulated device data**

Currently I'm simulating data that is being sent from Teltonika [TRB145](https://teltonika-networks.com/products/gateways/trb145) that is used as IoT gateway that reads over Modbus [RTU](https://www.modbustools.com/modbus.html) information from Power Quality Analyser. We one uses power quality analyser word it typicall referes to electrical values monitoring in order to detect anomalies and monitor electrical energy consumption and it's quality. Quality is usually checked vs some applied norm on the market like in the EU region EN 50160.

Simulated incoming data from the devices has following fields in standard JSON format:

| device_id | slave_id | measure_time | full_date | modbus_add | value | var_name |
| --------- | -------- | ------------ | --------- | ---------- | ----- | -------- |

With the amazing power of Go this is marshalled into [PQA_Message](https://github.com/Nikola-Vukasinovic/go-mqtt-device-mock/blob/main/mqtt.go).

One can appricate the elegance of the Go in this area and one can easily expand this to more complex formats.

### Planned works

It is in plan to add more message formats for devices that are out in the market and being used in Industrial IoT projects.
