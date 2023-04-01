# WiFi Control Server
This is a simple Golang server that allows you to control your NETGEAR WiFi Router Nighthawk (R7800) router by turning the WiFi on or off. It serves a web page that displays the current WiFi status and provides buttons to enable or disable the WiFi.

## Prerequisites

- Go 1.16 or higher
- Docker (optional)

## Getting Started

1. Clone the repository:
```bash
git clone https://github.com/fpesce/wifi-control-server.git
cd wifi-control-server
```
2. Set the required environment variables:
```bash
export ROUTER_URL="http://your-router-url"
export USERNAME="your-username"
export PASSWORD="your-password"
```
3. Build and run the server:
```bash
go build
./wifi-control-server
```
4. Open a web browser and navigate to http://localhost:8080. You should see the WiFi Control web page.

## Docker
Alternatively, you can build and run the server using Docker.

1. Build the Docker image:
```bash
make docker-build
```
2. Run the Docker container:
```bash
make docker-run
```

## License
This project is licensed under the Apache License v2. See the LICENSE file for details.

## Contributing
Please feel free to submit issues or pull requests for any improvements or bug fixes.
