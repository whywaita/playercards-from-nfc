package playercards

import "fmt"

type Card struct {
	Suit Suit `json:"suit"`
	Rank uint `json:"rank"`
}

type Suit int

const (
	SuitUnknown Suit = iota
	SuitSpade
	SuitHeart
	SuitDiamond
	SuitClub
)

func (s Suit) String() string {
	switch s {
	case SuitSpade:
		return "spades"
	case SuitHeart:
		return "hearts"
	case SuitDiamond:
		return "diamonds"
	case SuitClub:
		return "clubs"
	}

	return "unknown"
}

func (s Suit) StringShort() string {
	switch s {
	case SuitSpade:
		return "s"
	case SuitHeart:
		return "h"
	case SuitDiamond:
		return "d"
	case SuitClub:
		return "c"
	}

	return "u"
}

func getCardNumber() []string {
	return []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "T", "J", "Q", "K"}
}

func getSuit() []Suit {
	return []Suit{SuitSpade, SuitHeart, SuitDiamond, SuitClub}
}

func GeneratePlayerCardList() []string {
	numbers := getCardNumber()
	suits := getSuit()

	var playerCards []string

	for _, suit := range suits {
		for _, number := range numbers {
			playerCards = append(playerCards, fmt.Sprintf("%s%s", number, suit.StringShort()))
		}
	}

	return playerCards
}

func UnmarshalPlayerCard(in string) (Card, error) {
	var card Card

	if len(in) != 2 {
		return card, fmt.Errorf("invalid card length: %s", in)
	}

	switch in[1] {
	case 's':
		card.Suit = SuitSpade
	case 'h':
		card.Suit = SuitHeart
	case 'd':
		card.Suit = SuitDiamond
	case 'c':
		card.Suit = SuitClub
	default:
		return card, fmt.Errorf("invalid suit: %s", in)
	}

	switch in[0] {
	case 'A':
		card.Rank = 1
	case 'T':
		card.Rank = 10
	case 'J':
		card.Rank = 11
	case 'Q':
		card.Rank = 12
	case 'K':
		card.Rank = 13
	default:
		rank, err := fmt.Sscanf(in, "%d", &card.Rank)
		if err != nil {
			return card, fmt.Errorf("invalid rank: %s", in)
		}

		if rank != 1 {
			return card, fmt.Errorf("invalid rank: %s", in)
		}
	}

	return card, nil
}
