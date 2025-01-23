package main

import (
	"log"

	"github.com/sefaphlvn/bigbang/cmd"
)

var (
	Version             string // Bigbang version
	ControlPlaneVersion string // Control plane version
)

func main() {
	log.Printf("Bigbang Version: %s", Version)
	log.Printf("Control Plane Version: %s", ControlPlaneVersion)

	// go suubar.Start()
	cmd.Execute()
}
