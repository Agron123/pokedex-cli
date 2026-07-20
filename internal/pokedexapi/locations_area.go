package pokedexapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type LocationArea struct {
	PokemonEncounters []PokemonEncounter `json:"pokemon_encounters"`
}

type PokemonEncounter struct {
	Pokemon Pokemon `json:"pokemon"`
}

type Pokemon struct {
	Name string `json:"name"`
}

func (c Client) GetLocationArea(url string) (LocationArea, error) {
	data, ok := c.cache.Get(url)
	if ok {
		locationArea := LocationArea{}
		reader := bytes.NewReader(data)
		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(&locationArea); err != nil {
			return LocationArea{}, err
		}

		return locationArea, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return LocationArea{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return LocationArea{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return LocationArea{}, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	locationArea := LocationArea{}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return LocationArea{}, err
	}

	c.cache.Add(url, body)
	reader := bytes.NewReader(body)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&locationArea); err != nil {
		return LocationArea{}, err
	}

	return locationArea, nil
}
