package main

import (
	"log"

	"github.com/sefaphlvn/bigbang/cmd"
	"github.com/sefaphlvn/bigbang/pkg/version"
)

func main() {
	log.Printf("Envoy Version: %s", version.GetVersion())
	// go suubar.Start()
	cmd.Execute()
}
