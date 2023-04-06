package playercards

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/clausecker/nfc/v2"
	"github.com/goccy/go-yaml"
	"github.com/spf13/cobra"
	"github.com/whywaita/playercards-from-nfc/pkg/nfcc"
)

func ReadConfigPath(args []string) string {
	configPath := "./config.yaml"
	if len(args) != 0 {
		if _, err := os.Stat(args[0]); err != nil {
			return configPath
		}
		configPath = args[0]
	}

	return configPath
}

func NewGenerateConfigCmd(pnd nfc.Device) *cobra.Command {
	return &cobra.Command{
		Use:   "generate",
		Short: "Generate config.yaml",
		Long:  "Generate config.yaml from deck (default: ./config.yaml)",
		RunE: func(cmd *cobra.Command, args []string) error {
			log.Println("Start preparing config.yaml...")
			configPath := ReadConfigPath(args)

			cardConfigs, err := GenerateCardConfig(pnd)
			if err != nil {
				return fmt.Errorf("playercards.GenerateCardConfig(&pnd): %w", err)
			}

			b, err := yaml.Marshal(cardConfigs)
			if err != nil {
				return fmt.Errorf("yaml.Marshal(%v): %w", cardConfigs, err)
			}
			if err := os.WriteFile(configPath, b, 0600); err != nil {
				return fmt.Errorf("os.WriteFile(%s, b, 0600): %w", configPath, err)
			}
			return nil
		},
	}
}

func NewLoadCardCmd(pnd nfc.Device) *cobra.Command {
	return &cobra.Command{
		Use:   "load",
		Short: "Load card from config",
		Long:  "Load card from config.yaml (default: ./config.yaml)",
		RunE: func(cmd *cobra.Command, args []string) error {
			configPath := ReadConfigPath(args)

			loaded, err := LoadCards(pnd, configPath)
			if err != nil {
				return fmt.Errorf("loadCards(pnd, %s): %w", configPath, err)
			}

			log.Printf("You load cards is %s\n", loaded)
			return nil
		},
	}
}

func LoadCards(pnd nfc.Device, configPath string) ([]string, error) {
	cardConfigs, err := LoadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("playercards.LoadConfig(%s): %w", configPath, err)
	}

	loaded, err := loadCards(pnd, *cardConfigs, 2)
	if err != nil {
		return nil, fmt.Errorf("loadCards(pnd, *cardConfigs, 2): %w", err)
	}

	return loaded, nil
}

func LoadCardsWithChannel(ctx context.Context, pnd nfc.Device, cc CardConfig, number int, ch chan []string) error {
	var loaded []string

	for {
		in, err := nfcc.GetCard(pnd)
		if err != nil {
			return fmt.Errorf("nfcc.GetCard(pnd): %w", err)
		}

		card, err := LoadPlayerCard(in, cc)
		if err != nil {
			return fmt.Errorf("playercards.LoadPlayerCard(in, cardConfigs): %w", err)
		}

		//log.Printf("found card: %s", card)

		if len(loaded) == 0 {
			loaded = append(loaded, card)
			//log.Printf("append: loaded %v, card: %v", loaded, card)
			continue
		}
		if len(loaded) != 0 {
			for _, v := range loaded {
				if card != v {
					loaded = append(loaded, card)
				} else {
					continue
				}
			}
		}

		if len(loaded) == number {
			log.Printf("loaded: %v, will send", loaded)
			ch <- loaded
			loaded = nil
			continue
		}
	}
}

func loadCards(pnd nfc.Device, cc CardConfig, number int) ([]string, error) {
	var loaded []string

finish:
	for {
		in, err := nfcc.GetCard(pnd)
		if err != nil {
			return nil, fmt.Errorf("nfcc.GetCard(pnd): %w", err)
		}

		card, err := LoadPlayerCard(in, cc)
		if err != nil {
			return nil, fmt.Errorf("playercards.LoadPlayerCard(in, cardConfigs): %w", err)
		}

		if len(loaded) != 0 {
			for _, v := range loaded {
				if card == v {
					continue
				}
			}
		}

		loaded = append(loaded, card)

		if len(loaded) == number {
			break finish
		}
	}

	if len(loaded) != number {
		return nil, fmt.Errorf("invalid length (loaded: %v)", loaded)
	}

	return loaded, nil
}
