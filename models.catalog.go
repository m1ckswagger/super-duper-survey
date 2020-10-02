package main

import (
	"bufio"
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

type Catalog struct {
	Questions []Question `json:"questions"`
}
type Question struct {
	ID       int      `json:"id"`
	Question string   `json:"question"`
	Options  []string `json:"options"`
}

func NewCatalog(in, out string, options []string) *Catalog {
	inFile, err := os.OpenFile(in, os.O_RDONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer inFile.Close()
	scanner := bufio.NewScanner(inFile)
	var i int
	cata := new(Catalog)
	for scanner.Scan() {
		question := scanner.Text()
		i++
		cata.Questions = append(cata.Questions, Question{ID: i, Question: question, Options: options})
	}
	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	content, _ := json.MarshalIndent(cata, "", " ")
	ioutil.WriteFile(out, content, 0644)
	return cata
}
