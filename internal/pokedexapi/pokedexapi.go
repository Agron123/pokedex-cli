package pokedexapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
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

func (c Client) GetLocations(url string) (Page, error) {

	data, ok := c.cache.Get(url)
	if ok {
		page := Page{}
		reader := bytes.NewReader(data)
		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(&page); err != nil {
			return Page{}, err
		}

		return page, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Page{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Page{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Page{}, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	page := Page{}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Page{}, err
	}

	c.cache.Add(url, body)
	reader := bytes.NewReader(body)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&page); err != nil {
		return Page{}, err
	}

	return page, nil

}
