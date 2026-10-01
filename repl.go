package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"strings"

	"github.com/Agron123/pokedex-cli/internal/pokedexapi"
)

type cliCommand struct {
	name        string
	description string
	callback    func(c *config, args []string) error
}

type config struct {
	next      *string
	previous  *string
	apiClient pokedexapi.Client
	pokedex   map[string]pokedexapi.Pokemon
}

func getCommands() map[string]cliCommand {
	return map[string]cliCommand{
		"exit": {
			name:        "exit",
			description: "Exit the Pokedex",
			callback:    commandExit,
		},
		"help": {
			name:        "help",
			description: "Displays a help message",
			callback:    commandHelp,
		},
		"map": {
			name:        "map",
			description: "Prints 20 locations",
			callback:    commandMap,
		},
		"mapb": {
			name:        "mapb",
			description: "Prints 20 previous locations",
			callback:    commandMapb,
		},
		"explore": {
			name:        "explore",
			description: "Prints pokemons in area",
			callback:    commandExplore,
		},
		"catch": {
			name:        "catch",
			description: "Catches a pokemon and adds it to the pokedex",
			callback:    commandCatch,
		},
		"inspect": {
			name:        "inspect",
			description: "Prints pokemon stats",
			callback:    commandInspect,
		},
		"pokedex": {
			name:        "pokedex",
			description: "Prints all pokemons in your pokedex",
			callback:    commandPokedex,
		},
	}
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	words := strings.Fields(lowerText)
	return words
}

func commandExit(c *config, args []string) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config, args []string) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	commands := getCommands()

	for _, value := range commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}

func commandMap(c *config, args []string) error {
	var url string

	if c.next == nil {
		url = locationAreaURL
	} else {
		url = *c.next
	}

	page, err := c.apiClient.GetLocations(url)
	if err != nil {
		return err
	}

	for _, location := range page.Results {
		fmt.Println(location.Name)
	}

	c.next = page.Next
	c.previous = page.Previous

	return nil
}

func commandMapb(c *config, args []string) error {
	var url string

	if c.previous == nil {
		fmt.Println("You are on the first page")
		return nil
	}

	url = *c.previous

	page, err := c.apiClient.GetLocations(url)
	if err != nil {
		return err
	}

	for _, location := range page.Results {
		fmt.Println(location.Name)
	}

	c.next = page.Next
	c.previous = page.Previous

	return nil
}

func commandExplore(c *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Usage: explore <location-area>")
		return nil
	}

	areaName := args[0]

	url := locationAreaURL + "/" + areaName
	locationArea, err := c.apiClient.GetLocationArea(url)
	if err != nil {
		return err
	}

	for _, pokemonEncounter := range locationArea.PokemonEncounters {
		fmt.Println(pokemonEncounter.Pokemon.Name)
	}
	return nil

}

func commandCatch(c *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Usage: catch <pokemon-name>")
		return nil
	}

	pokemonName := args[0]

	fmt.Printf("Throwing a Pokeball at %s...\n", pokemonName)

	url := pokemonBaseURL + "/" + pokemonName
	pokemonInfo, err := c.apiClient.GetPokemonInfo(url)
	if err != nil {
		return err
	}

	randomInt := rand.Intn(100)
	if randomInt > pokemonInfo.BaseExperience/5 {
		c.pokedex[pokemonInfo.Name] = pokemonInfo
		fmt.Printf("%s was caught!\n", pokemonInfo.Name)
	} else {
		fmt.Printf("%s escaped!\n", pokemonInfo.Name)
	}
	return nil
}

func commandInspect(c *config, args []string) error {
	if len(args) == 0 {
		fmt.Println("Usage: inspect <pokemon-name>")
		return nil
	}

	pokemonName := args[0]
	pokemon, exists := c.pokedex[pokemonName]
	if !exists {
		fmt.Printf("You have not caught %s\n", pokemonName)
		return nil
	}

	fmt.Printf("Name: %s\n", pokemon.Name)
	fmt.Printf("Height: %d\n", pokemon.Height)
	fmt.Printf("Weight: %d\n", pokemon.Weight)
	fmt.Println("Stats:")
	for _, stat := range pokemon.Stats {
		fmt.Printf("-%s: %d\n", stat.Stat.Name, stat.BaseStat)
	}

	fmt.Println("Types:")
	for _, types := range pokemon.Types {
		fmt.Printf("-%s\n", types.Types.Name)
	}
	return nil
}

func commandPokedex(c *config, args []string) error {
	fmt.Println("Your pokedex:")

	for key := range c.pokedex {
		fmt.Printf("-%s\n", key)
	}

	return nil
}

func StartRepl() {
	reader := bufio.NewScanner(os.Stdin)
	allowedCommands := getCommands()
	client := pokedexapi.NewClient(timeout, interval)
	cfg := config{apiClient: client,
		pokedex: make(map[string]pokedexapi.Pokemon)}

	for {
		fmt.Print("Pokedex > ")
		reader.Scan()

		words := cleanInput(reader.Text())
		if len(words) == 0 {
			continue
		}

		command := words[0]
		c, ok := allowedCommands[command]
		if !ok {
			fmt.Println("Unknown command")
			continue
		}
		err := c.callback(&cfg, words[1:])
		if err != nil {
			fmt.Println(err)
		}

	}
}
