package book

import (
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/primalcs/words_of_wisdom/internal/interfaces"
)

var _ interfaces.Book = &wisdom{}

type wisdom struct {
	book []string
}

func NewWisdom() (*wisdom, error) {
	ba, err := os.ReadFile("./docs/words.txt")
	if err != nil {
		return nil, fmt.Errorf("could not read the book file")
	}
	b := strings.Split(string(ba), "\n")
	return &wisdom{book: b}, nil
}

func (b *wisdom) GetOne() string {
	n := rand.Intn(len(b.book))
	return b.book[n]
}
