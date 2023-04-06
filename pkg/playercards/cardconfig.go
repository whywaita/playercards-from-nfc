package playercards

import (
	"encoding/hex"
	"fmt"
	"log"
	"os"

	"github.com/clausecker/nfc/v2"
	"github.com/goccy/go-yaml"
	"github.com/whywaita/playercards-from-nfc/pkg/nfcc"
)

type CardConfig struct {
	CardIDs map[string]string `yaml:"card_ids"` // key: uid value: card
}

func GenerateCardConfig(pnd nfc.Device) (*CardConfig, error) {
	cards := GeneratePlayerCardList()

	cardIDs := map[string]string{}

	for _, card := range cards {
		log.Printf("Please read %s\n", card)
		cardUID, err := nfcc.GetCard(pnd)
		if err != nil {
			return nil, fmt.Errorf("nfcc.GetCard(pnd): %w", err)
		}
		cus := hex.EncodeToString(cardUID[:])
		if v, ok := cardIDs[cus]; ok {
			return nil, fmt.Errorf("found same uid (%s)", v)
		}

		cardIDs[cus] = card
	}

	cc := &CardConfig{
		CardIDs: cardIDs,
	}

	return cc, nil
}

func LoadConfig(p string) (*CardConfig, error) {
	b, err := os.ReadFile(p)
	if err != nil {
		return nil, fmt.Errorf("os.ReadFile(%s): %w", p, err)
	}

	cc, err := UnmarshalCardConfig(b)
	if err != nil {
		return nil, fmt.Errorf("UnmarshalCardConfig(%s): %w", b, err)
	}

	return cc, nil
}

func UnmarshalCardConfig(in []byte) (*CardConfig, error) {
	var cc CardConfig
	if err := yaml.Unmarshal(in, &cc); err != nil {
		return nil, fmt.Errorf("yaml.Unmarshal(); %w", err)
	}

	return &cc, nil
}

func LoadPlayerCard(in [10]byte, cc CardConfig) (string, error) {
	str := hex.EncodeToString(in[:])

	v, ok := cc.CardIDs[str]
	if !ok {
		return "", fmt.Errorf("unknown card")
	}

	return v, nil
}
