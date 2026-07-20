package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Agron123/pokedexcli/internal/pokedexapi"
)

const locationAreaURL = "https://pokeapi.co/api/v2/location-area"
const interval = 5 * time.Minute
const timeout = 10 * time.Second

type cliCommand struct {
	name        string
	description string
	callback    func(c *config) error
}

type config struct {
	next      *string
	previous  *string
	apiClient pokedexapi.Client
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
	}
}

func cleanInput(text string) []string {
	lowerText := strings.ToLower(text)
	words := strings.Fields(lowerText)
	return words
}

func commandExit(c *config) error {
	fmt.Println("Closing the Pokedex... Goodbye!")
	os.Exit(0)
	return nil
}

func commandHelp(c *config) error {
	fmt.Println("Welcome to the Pokedex!\nUsage:")
	commands := getCommands()

	for _, value := range commands {
		fmt.Printf("%s: %s\n", value.name, value.description)
	}
	return nil
}

func commandMap(c *config) error {
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

func commandMapb(c *config) error {
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

func StartRepl() {
	reader := bufio.NewScanner(os.Stdin)
	allowedCommands := getCommands()
	client := pokedexapi.NewClient(timeout, interval)
	cfg := config{apiClient: client}

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
		err := c.callback(&cfg)
		if err != nil {
			fmt.Println(err)
		}

	}
}
