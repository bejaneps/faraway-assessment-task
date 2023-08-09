package quotes

import (
	_ "embed"
	"encoding/json"

	"github.com/bejaneps/faraway-assessment-task/internal/pkg/log"
)

//go:embed quotes.json
var quotesJSON []byte

type Quote struct {
	ID    int
	Quote string
}

var Seeds []Quote

func SeedsToStrings() []string {
	var strings []string

	for _, seed := range Seeds {
		strings = append(strings, seed.Quote)
	}

	return strings
}

func init() {
	if err := json.Unmarshal(quotesJSON, &Seeds); err != nil {
		log.Fatal(err.Error())
	}
}
