package pokedexapi

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Location struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type Page struct {
	Next     *string    `json:"next"`
	Previous *string    `json:"previous"`
	Results  []Location `json:"results"`
}

func GetLocations(url string) (Page, error) {
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Page{}, err
	}

	res, err := client.Do(req)
	if err != nil {
		return Page{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Page{}, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	page := Page{}
	decoder := json.NewDecoder(res.Body)
	if err := decoder.Decode(&page); err != nil {
		return Page{}, err
	}

	return page, nil

}
