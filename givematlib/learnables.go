package givematlib

import (
	"bufio"
	"io"
)

type language string
type status struct {
	known []string
}

func getStatusFile(lang language) io.Reader {
	return nil
}

func readStatus(reader io.Reader) (status, error) {
	var lines []string
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return status{known: lines}, scanner.Err()
}

func saveStatus(writer io.Writer, status status) {

}
