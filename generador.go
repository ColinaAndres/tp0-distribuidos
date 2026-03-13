package main

import (
	"fmt"
	"os"
)

const (
	minArgsRequired = 3
)

const serverTemplate = `  server:
    container_name: server
    image: server:latest
    entrypoint: python3 /main.py
    environment:
      - PYTHONUNBUFFERED=1
      - LOGGING_LEVEL=DEBUG
    networks:
      - testing_net
`

const clientTemplate = `  client1:
    container_name: client1
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID=1
      - CLI_LOG_LEVEL=DEBUG
    networks:
      - testing_net
    depends_on:
      - server
`

const networkTemplate = `networks:
  testing_net:
    ipam:
      driver: default
      config:
        - subnet: 172.25.125.0/24
`

func serverContent() string {
	return serverTemplate
}

func clientContent() string {
	return clientTemplate
}

func network() string {
	return networkTemplate
}

func services() string {
	return fmt.Sprintf("services:\n%s\n%s\n",
		serverContent(),
		clientContent(),
	)
}

func fileContent() string {
	filecontent := "name: tp0\n"
	filecontent += services()
	filecontent += network()
	return filecontent
}

func main() {
	if len(os.Args) < minArgsRequired {
		fmt.Println("Bad invocation. Usage: go run generator.go <docker-compose-file-name> <amount-of-clients>")
		os.Exit(1)
	}

	fileName := os.Args[1]

	err := os.WriteFile(fileName, []byte(fileContent()), 0644)
	if err != nil {
		fmt.Printf("Error al crear el archivo: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Archivo %s generado con éxito.\n", fileName)
}
