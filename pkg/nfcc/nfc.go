package nfcc

import (
	"fmt"

	"github.com/clausecker/nfc/v2"
)

var (
	// These settings works with the ACR122U. Your milage may vary with
	// other devices.
	m = nfc.Modulation{Type: nfc.ISO14443a, BaudRate: nfc.Nbr106}
	// Use an empty string to select first device libnfc sees
	//devstr = ""
)

func GetCard(pnd nfc.Device) ([10]byte, error) {
	for {
		targets, err := pnd.InitiatorListPassiveTargets(m)
		if err != nil {
			return [10]byte{}, fmt.Errorf("failed to list nfc targets: %w", err)
		}

		for _, t := range targets {
			if card, ok := t.(*nfc.ISO14443aTarget); ok {
				return card.UID, nil
			}
		}
	}
}
