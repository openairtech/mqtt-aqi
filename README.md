# OpenAir-MQTT-AQI

Subscribes to MQTT PM2.5/PM10 value topics and publish calculated US AQI value.

## Build

You can build **OpenAir-MQTT-AQI** binary using `make build` or `make build-static`. Go 1.22.4 or newer is needed for the build.

## Usage

```
Usage: openair-mqtt-aqi [options]

Options:
  -P string
    	mqtt password (optional)
  -a string
    	mqtt topic to subscribe to PM2.5 value (default "OpenAir/SDS011/PM2.5")
  -b string
    	mqtt topic to subscribe to PM10 value (default "OpenAir/SDS011/PM10")
  -d	enable debug messages
  -h string
    	mqtt host to connect to (default "localhost")
  -i string
    	mqtt client id to use (default "openair_mqtt_aqi")
  -p int
    	network port to connect to (default 1883)
  -q string
    	mqtt topic to publish AQI value to (default "OpenAir/AQI")
  -r	mqtt publish retained flag
  -t duration
    	AQI value publish period (default 30s)
  -u string
    	mqtt user (optional)
  -v	print the version number and quit
```

## License

OpenAir-MQTT-AQI is released under the Apache 2.0 license. See [LICENSE.txt](https://github.com/openairtech/mqtt-aqi/blob/master/LICENSE.txt)
