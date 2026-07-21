package pokedexapi

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func (c Client) GetPokemonInfo(url string) (Pokemon, error) {
	data, ok := c.cache.Get(url)
	if ok {
		pokemon := Pokemon{}
		reader := bytes.NewReader(data)
		decoder := json.NewDecoder(reader)
		if err := decoder.Decode(&pokemon); err != nil {
			return Pokemon{}, err
		}

		return pokemon, nil
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return Pokemon{}, err
	}

	res, err := c.httpClient.Do(req)
	if err != nil {
		return Pokemon{}, err
	}

	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return Pokemon{}, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	pokemon := Pokemon{}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return Pokemon{}, err
	}

	c.cache.Add(url, body)
	reader := bytes.NewReader(body)
	decoder := json.NewDecoder(reader)
	if err := decoder.Decode(&pokemon); err != nil {
		return Pokemon{}, err
	}

	return pokemon, nil
}
