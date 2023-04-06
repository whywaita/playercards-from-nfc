package server

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/clausecker/nfc/v2"
	"github.com/spf13/cobra"
	"github.com/whywaita/playercards-from-nfc/pkg/playercards"
)

func NewServerCmd(pnd nfc.Device) *cobra.Command {
	return &cobra.Command{
		Use:   "server",
		Short: "Start HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			ch := make(chan []string, 2)

			configPath := playercards.ReadConfigPath(args)
			cardConfigs, err := playercards.LoadConfig(configPath)
			if err != nil {
				return fmt.Errorf("playercards.LoadConfig(%s): %w", configPath, err)
			}

			go func() {
				log.Printf("Start loading cards...")
				if err := playercards.LoadCardsWithChannel(ctx, pnd, *cardConfigs, 2, ch); err != nil {
					log.Printf("playercards.LoadCardsWithChannel(ctx): %v", err)
				}
			}()

			m := NewMux(ch)
			if err := http.ListenAndServe(":8080", m); err != nil {
				return fmt.Errorf("failed to start server: %w", err)
			}

			return nil
		},
	}
}

type SendCard struct {
	Suit string `json:"suit"`
	Rank uint   `json:"rank"`
}

var (
	now []string
)

// NewMux create a new http.ServeMux.
func NewMux(ch chan []string) *http.ServeMux {
	m := http.NewServeMux()

	m.HandleFunc("/sse", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")

		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		go func() {
			for {
				v := <-ch
				log.Printf("receive card in server: %s", v)
				now = v
			}
		}()

		for {
			select {
			case <-ticker.C:
				//log.Printf("Start ticker now: %s", now)
				if len(now) != 2 {
					continue
				}

				c1, err := playercards.UnmarshalPlayerCard(now[0])
				if err != nil {
					log.Printf("UnmarshalPlayerCard(%s): %v", now[0], err)
					fmt.Fprint(w, "event: error")
					return
				}
				c2, err := playercards.UnmarshalPlayerCard(now[1])
				if err != nil {
					log.Printf("UnmarshalPlayerCard(%s): %v", now[1], err)
					fmt.Fprint(w, "event: error")
					return
				}

				cards := []SendCard{
					{
						Suit: c1.Suit.String(),
						Rank: c1.Rank,
					},
					{
						Suit: c2.Suit.String(),
						Rank: c2.Rank,
					},
				}
				b, err := json.Marshal(cards)
				if err != nil {
					log.Printf("json.Marshal(%v): %v", cards, err)
					fmt.Fprint(w, "event: error")
					return
				}
				log.Printf("send cards: %s", string(b))
				fmt.Fprintf(w, "data: %s\n\n", string(b))
				w.(http.Flusher).Flush()
			case <-r.Context().Done():
				return
			}

		}
	})

	return m
}
