package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/kkpoon/aqicalc"
)

var (
	Version   = "unknown"
	Timestamp = "unknown"
)

func main() {
	versionFlag := flag.Bool("v", false, "print the version number and quit")

	var debugFlag bool
	flag.BoolVar(&debugFlag, "d", false, "enable debug messages")

	host := flag.String("h", "localhost", "mqtt host to connect to")
	port := flag.Int("p", 1883, "network port to connect to")

	clientId := flag.String("i", "esp_mqtt_aqi", "mqtt client id to use")
	user := flag.String("u", "", "mqtt user (optional)")
	password := flag.String("P", "", "mqtt password (optional)")
	retained := flag.Bool("r", false, "mqtt publish retained flag")

	pm25Topic := flag.String("a", "OpenAir/SDS011/PM2.5", "mqtt topic to subscribe to PM2.5 value")
	pm10Topic := flag.String("b", "OpenAir/SDS011/PM10", "mqtt topic to subscribe to PM10 value")
	aqiTopic := flag.String("q", "OpenAir/AQI", "mqtt topic to publish AQI value to")
	aqiPublishPeriod := flag.Duration("t", 30*time.Second, "AQI value publish period")

	flag.Parse()

	if *versionFlag {
		fmt.Printf("Build version: %s\n", Version)
		fmt.Printf("Build timestamp: %s\n", Timestamp)
		return
	}

	url := fmt.Sprintf("tcp://%s:%d", *host, *port)

	opts := mqtt.NewClientOptions().AddBroker(url)
	opts.SetAutoReconnect(true)
	opts.SetClientID(*clientId)
	opts.SetUsername(*user)
	opts.SetPassword(*password)

	client := mqtt.NewClient(opts)
	if token := client.Connect(); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	if debugFlag {
		log.Printf("connected to %s...", url)
	}

	signalCh := make(chan os.Signal, 1)
	signal.Notify(signalCh,
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	)

	pm25Value := -1.
	if token := client.Subscribe(*pm25Topic, 0, func(c mqtt.Client, m mqtt.Message) {
		if debugFlag {
			log.Printf("received PM2.5 value: %s", m.Payload())
		}
		if pm, err := strconv.ParseFloat(string(m.Payload()), 32); err == nil {
			pm25Value = pm
		}
	}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	pm10Value := -1.
	if token := client.Subscribe(*pm10Topic, 0, func(c mqtt.Client, m mqtt.Message) {
		if debugFlag {
			log.Printf("received PM10 value: %s", m.Payload())
		}
		if pm, err := strconv.ParseFloat(string(m.Payload()), 32); err == nil {
			pm10Value = pm
		}
	}); token.Wait() && token.Error() != nil {
		log.Fatal(token.Error())
	}

	for {
		select {
		case <-time.After(*aqiPublishPeriod):
			if pm25Value >= 0 && pm10Value >= 0 {
				aqiValue := aqicalc.CalculateAQI(aqicalc.Conc{PM25_24hr: pm25Value, PM10_24hr: pm10Value}).AQI
				if debugFlag {
					log.Printf("publish AQI value: %d", aqiValue)
				}
				if token := client.Publish(*aqiTopic, 0, *retained, strconv.Itoa(aqiValue)); token.Wait() && token.Error() != nil {
					log.Printf("publish error: %v", token.Error())
				}
			}
		case sig := <-signalCh:
			log.Printf("received %v signal, exiting...", sig)
			return
		}
	}
}
