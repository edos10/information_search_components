package main

import (
	"fmt"
	"log"

	"github.com/aaaton/golem/v4"
	"github.com/aaaton/golem/v4/dicts/ru"
)

func main() {
	lemmatizer, err := golem.New(ru.New())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(lemmatizer.Lemma("Кошек"))
}
