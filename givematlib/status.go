package givematlib

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type StatusCache struct {
	cache map[string][]string
}

func NewStatusCache() StatusCache {
	return StatusCache{
		cache: make(map[string][]string),
	}
}

func (sc StatusCache) ReadLearnableStatus(language string) ([]string, error) {
	if val, ok := sc.cache[language]; ok {
		return val, nil
	} else {
		learnables, err := ReadLearnableStatus(language)
		if err != nil {
			return nil, err
		}

		sc.cache[language] = learnables
		return learnables, nil
	}
}

func SaveLearnableStatus(
	name string,
	language string,
	learnables []string,
) error {
	filename := fmt.Sprintf("%s.%s", name, language)

	file, err := InDataDir("status", filename)
	if err != nil {
		return err
	}

	f, err := os.Create(file)
	if err != nil {
		return err
	}
	defer f.Close()

	writer := bufio.NewWriter(f)
	defer writer.Flush()

	for _, learnable := range learnables {
		for _, learnableWord := range strings.Fields(learnable) {
			_, err := writer.WriteString(fmt.Sprintf("%s\n", learnableWord))
			if err != nil {
				return err
			}
		}
	}

	return nil
}

func ReadLearnableStatus(language string) ([]string, error) {
	statusDir, err := InDataDir("status")
	if err != nil {
		return nil, err
	}

	files, err := ioutil.ReadDir(statusDir)
	if err != nil {
		return nil, err
	}

	var learnables []string
	for _, file := range files {
		if strings.HasSuffix(file.Name(), fmt.Sprintf(".%s", language)) {
			singleFileLearnables, err := ReadLearnableStatusFile(file.Name())
			if err != nil {
				return nil, err
			}

			learnables = append(learnables, singleFileLearnables...)
		}
	}
	return learnables, nil
}

func ReadLearnableStatusFile(filename string) ([]string, error) {
	var lines []string
	file, err := InDataDir("status", filename)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
