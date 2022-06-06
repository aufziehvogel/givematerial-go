package givematlib

import (
	"crypto/sha256"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Text struct {
	Title      string   `xml:"title"`
	Language   string   `xml:"language"`
	Fulltext   string   `xml:"fulltext"`
	Learnables []string `xml:"learnables>learnable"`
}

func (t *Text) Unknown(knownLearnables []string) (unknown []string) {
	knownMap := make(map[string]struct{})

	for _, item := range knownLearnables {
		knownMap[strings.ToLower(item)] = struct{}{}
	}

	for _, item := range t.Learnables {
		if _, ok := knownMap[strings.ToLower(item)]; !ok {
			unknown = append(unknown, item)
		}
	}
	return
}

func ListTexts() ([]string, error) {
	var texts []string

	textsDir, err := InDataDir("texts")
	if err != nil {
		return texts, err
	}

	files, err := ioutil.ReadDir(textsDir)
	if err != nil {
		return texts, err
	}

	for _, file := range files {
		texts = append(texts, file.Name())
	}
	return texts, nil
}

func LoadText(textID string) (Text, error) {
	var t Text
	xmlFilePath, err := InDataDir("texts", textID)
	if err != nil {
		return t, err
	}

	xmlFile, err := os.Open(xmlFilePath)
	if err != nil {
		return t, err
	}

	content, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return t, err
	}

	xml.Unmarshal(content, &t)
	return t, nil
}

func SaveText(text Text) error {
	textID := fmt.Sprintf("%x", sha256.Sum256([]byte(text.Title)))
	xmlFilePath, err := InDataDir("texts", textID)
	if err != nil {
		return err
	}

	xmlString, _ := xml.MarshalIndent(text, "", " ")
	return ioutil.WriteFile(xmlFilePath, xmlString, 0644)
}
