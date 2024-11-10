package processing

import (
	"fmt"
	"strings"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/en"
	"github.com/kljensen/snowball"
)

const engLang = "english"

type Processing interface {
	Lemming(text string) ([]string, error)
	Stemming(text string) ([]string, error)
}

// структура обработчика текста перед хранением,
// можно хранить в RAM
type MyProcessing struct {
	StopWords map[string]bool
}

func (p *MyProcessing) Lemming(text string) ([]string, error) {
	words := strings.Fields(text)

	lemmatizer, err := golem.New(en.New())
	if err != nil {
		return nil, fmt.Errorf("golem.New: %w", err)
	}

	lemmatizedWords := make([]string, len(words))
	for i, word := range words {
		word = strings.ToLower(word)
		lemmatizedWords[i] = lemmatizer.Lemma(word)
	}

	var finalDocument []string
	for _, word := range lemmatizedWords {
		if !p.StopWords[word] {
			finalDocument = append(finalDocument, word)
		}
	}

	return finalDocument, nil
}

func (p *MyProcessing) Stemming(text string) ([]string, error) {
	words := strings.Fields(text)
	var stemmedWords []string

	for _, word := range words {
		word = strings.ToLower(word)
		stemmedWord, err := snowball.Stem(word, engLang, true)
		if err != nil {
			return nil, err
		}
		stemmedWords = append(stemmedWords, stemmedWord)
	}

	var result []string
	for _, word := range stemmedWords {
		if !p.StopWords[word] {
			result = append(result, word)
		}
	}

	return result, nil
}

func NewMyProcessing(stopWords []string) *MyProcessing {
	proc := MyProcessing{}
	proc.StopWords = make(map[string]bool, 0)
	for _, word := range stopWords {
		word = strings.ToLower(word)
		proc.StopWords[word] = true
	}
	return &proc
}

func (p *MyProcessing) UpdateStopWords(words ...string) {
	for _, value := range words {
		value = strings.ToLower(value)
		p.StopWords[value] = true
	}
}

func (p *MyProcessing) DeleteStopWords(words ...string) {
	for _, value := range words {
		value = strings.ToLower(value)
		delete(p.StopWords, value)
	}
}
