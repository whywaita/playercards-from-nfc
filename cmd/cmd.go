package main

import (
	"fmt"
	"log"
	"os"

	nfc "github.com/clausecker/nfc/v2"
	"github.com/spf13/cobra"
	"github.com/whywaita/playercards-from-nfc/pkg/playercards"
	"github.com/whywaita/playercards-from-nfc/pkg/server"
)

func init() {
	log.Println("using libnfc", nfc.Version())

	p, err := nfc.Open("")
	if err != nil {
		log.Panicf("could not open device: %v", err)
	}

	if err := p.InitiatorInit(); err != nil {
		log.Panicf("could not init initiator: %v", err)
	}

	log.Println("opened device", p, p.Connection())
	pnd = p

	generateCmd := playercards.NewGenerateConfigCmd(pnd)
	rootCmd.AddCommand(generateCmd)

	loadCmd := playercards.NewLoadCardCmd(pnd)
	rootCmd.AddCommand(loadCmd)

	serverCmd := server.NewServerCmd(pnd)
	rootCmd.AddCommand(serverCmd)
}

func main() {
	defer pnd.Close()
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "playercards-from-nfc",
	Short: "playercards-from-nfc",
}

var (
	pnd = nfc.Device{}
)
