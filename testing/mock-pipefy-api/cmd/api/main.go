package main

import "mock-pipefy-api/internal/server"

func main() {
	println("Mock Pipefy API is running...")

	defer println("Mock Pipefy API stopped.")

	if err := server.Init().Start(); err != nil {
		println("Error starting Mock Pipefy API:", err.Error())
	}
}
