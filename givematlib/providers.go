package givematlib

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type StatusProvider interface {
	FetchLearnables() []string
}

type wanikaniPages struct {
	NextURL string `json:"next_url"`
}

type wanikaniAssignmentsResult struct {
	Pages wanikaniPages `json:"pages"`
	Data  []struct {
		Data struct {
			SubjectId   int    `json:"subject_id"`
			SubjectType string `json:"subject_type"`
		} `json:"data"`
	} `json:"data"`
}

type wanikaniSubjectsResult struct {
	Pages wanikaniPages `json:"pages"`
	Data  []struct {
		Id     int    `json:"id"`
		Object string `json:"object"`
		Data   struct {
			Characters string `json:"characters"`
		} `json:"data"`
	} `json:"data"`
}

type WanikaniProvider struct {
	apiToken        string
	subjectsToCheck map[string]struct{}
	httpClient      *http.Client
}

func NewWanikaniProvider(apiToken string) *WanikaniProvider {
	subjectsToCheck := map[string]struct{}{
		"kanji":      struct{}{},
		"vocabulary": struct{}{},
	}

	return &WanikaniProvider{
		apiToken,
		subjectsToCheck,
		&http.Client{Timeout: 10 * time.Second},
	}
}

func (w WanikaniProvider) FetchLearnables() ([]string, error) {
	subjects, err := w.fetchSubjects()
	if err != nil {
		return nil, err
	}

	learnablesMap, err := w.fetchLearnables([]int{5, 6, 7, 8, 9}, subjects)
	if err != nil {
		return nil, err
	}

	learnablesSlice := make([]string, 0, len(learnablesMap))
	for learnable, _ := range learnablesMap {
		learnablesSlice = append(learnablesSlice, learnable)
	}
	return learnablesSlice, nil
}

func (w WanikaniProvider) fetchSubjects() (map[int]string, error) {
	subjects := make(map[int]string)

	nextURL := "https://api.wanikani.com/v2/subjects"

	for nextURL != "" {
		subjectsResponse := new(wanikaniSubjectsResult)
		err := w.fetchWanikani(nextURL, subjectsResponse)
		if err != nil {
			return nil, err
		}

		for _, data := range subjectsResponse.Data {
			if _, ok := w.subjectsToCheck[data.Object]; ok {
				subjects[data.Id] = data.Data.Characters
			}
		}

		nextURL = subjectsResponse.Pages.NextURL
	}

	return subjects, nil
}

func (w WanikaniProvider) fetchLearnables(
	levels []int,
	subjects map[int]string,
) (map[string]struct{}, error) {
	levelStrings := make([]string, len(levels))
	for i, level := range levels {
		levelStrings[i] = strconv.Itoa(level)
	}

	nextURL := fmt.Sprintf(
		"https://api.wanikani.com/v2/assignments?srs_stages=%s",
		strings.Join(levelStrings, ","),
	)

	learnables := make(map[string]struct{})

	for nextURL != "" {
		assignmentsResponse := new(wanikaniAssignmentsResult)
		err := w.fetchWanikani(nextURL, assignmentsResponse)
		if err != nil {
			return nil, err
		}

		for _, data := range assignmentsResponse.Data {
			if _, ok := w.subjectsToCheck[data.Data.SubjectType]; ok {
				learnables[subjects[data.Data.SubjectId]] = struct{}{}
			}
		}

		nextURL = assignmentsResponse.Pages.NextURL
	}

	return learnables, nil
}

func (w WanikaniProvider) fetchWanikani(url string, target interface{}) error {
	//log.Printf("Fetching Wanikani URL %q", url)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.Header.Set("Wanikani-Revision", "20170710")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", w.apiToken))

	r, err := w.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer r.Body.Close()

	return json.NewDecoder(r.Body).Decode(target)
}

type AnkiProvider struct {
	Decks    []string
	AnkiFile string
}

func (p AnkiProvider) FetchLearnables() ([]string, error) {
	db, err := sql.Open("sqlite3", p.AnkiFile)
	if err != nil {
		log.Println("Could not open database")
		return nil, err
	}
	defer db.Close()

	var words []string
	for _, deck := range p.Decks {
		row, err := db.Query(`SELECT substr(n.flds, 0, instr(n.flds, char(31)))
FROM cards c
JOIN notes n ON c.nid = n.id
JOIN decks d ON c.did = d.id
WHERE d.name = ? COLLATE NOCASE AND c.ivl != 0;`, deck)

		if err != nil {
			log.Println("Could not execute query")
			return nil, err
		}
		defer row.Close()

		for row.Next() {
			var word string
			if err := row.Scan(&word); err != nil {
				log.Println("Could not parse result")
				return nil, err
			}
			words = append(words, word)
		}
	}

	return words, nil
}
