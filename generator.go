package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	minArgsRequired     = 3
	serverContainerName = "server"
	serverImage         = "server:latest"
	clientImage         = "client:latest"
	networkName         = "testing_net"
	debugLogLevel       = "DEBUG"
)

const serverTemplate = `  server:
    container_name: %s
    image: %s
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=%s
    networks:
      - %s
`

const clientTemplate = `  client1:
    container_name: client1
    image: %s
    entrypoint: /client
    environment:
      - CLI_ID=1
      - CLI_LOG_LEVEL=%s
    networks:
      - %s
    depends_on:
      - %s
`

const networkTemplate = `networks:
  %s:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
`

func serverContent() string {
	return fmt.Sprintf(serverTemplate, serverContainerName, serverImage, debugLogLevel, networkName)
}

func clientContent() string {
	return fmt.Sprintf(clientTemplate, clientImage, debugLogLevel, networkName, serverContainerName)
}

func network() string {
	return fmt.Sprintf(networkTemplate, networkName)
}

func services() string {
	var parts []string
	servicesHeader := "services:"
	parts = append(parts, servicesHeader, serverContent(), clientContent())
	return strings.Join(parts, "\n")
}

func fileContent() string {
	var parts []string
	fileHeader := "name: tp0"
	parts = append(parts, fileHeader, services(), network())
	return strings.Join(parts, "\n")
}

func main() {
	if len(os.Args) < minArgsRequired {
		fmt.Println("Bad invocation. Usage: go run generator.go <docker-compose-file-name> <amount-of-clients>")
		os.Exit(1)
	}

	fileName := os.Args[1]

	err := os.WriteFile(fileName, []byte(fileContent()), 0644)
	if err != nil {
		fmt.Printf("Error to write file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Your file %s was generated sucessfully.\n", fileName)
}
