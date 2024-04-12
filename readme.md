# Vonkje
Vonkje is a small service which manages communication, metrics, decision making and control of a solar plant. The goal of this application is to provide insight and to improve the efficiency of a solar plant.

## Features
- Metrics collection of devices
- Controlling state of devices
- Collecting power prices from suppliers

## Supported Devices
- Huawei Sun2000 and connected peripherals like Luna2000 battery and power meter.

## Runtime Dependencies
- **Grafana** for visualisation
- **Victoria Metrics** for metrics storage

## Development Dependencies
- **Golang** 1.22 or higher
- **Docker** for Victoria Metrics and Grafana

## Images
Overview:
![Grafana](./docs/images/overview.png)
Solar:
![Solar](./docs/images/solar.png)
Battery:
![Battery](./docs/images/battery.png)
Power Meter:
![Power Meter](./docs/images/power-meter.png)
