package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// Country API struct - Basic Information.
type Country struct {
	CountryName  string            `json:"name"`
	Region       string            `json:"region"`
	Capital      string            `json:"capital"`
	Translations map[string]string `json:"translations"`
	Languages    []Language        `json:"languages"`
	Currencies   []Currency        `json:"Currencies"`
}

// Language struct - Some countries have more than one language.
type Language struct {
	Language string `json:"name"`
}

// Currency struct - Some countries have more than one currency.
type Currency struct {
	CurrencyName   string `json:"name"`
	CurrencySymbol string `json:"symbol"`
}

var (
	gameIndex int       = -1 // Correct answer index in game list.
	finalList []Country      // Game list.
)

const botToken = "BOT_TOKEN" // Remove on GitHub Version.

// Function to get API data: https://restcountries.eu/
func getCountriesAPI() []Country {
	var countries []Country

	resp, err := http.Get("https://restcountries.eu/rest/v2/all")

	if err != nil {
		log.Fatalln(err)
	}

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatalln(err)
	}

	json.Unmarshal(body, &countries)

	return countries
}

// Generates a random value according to a maximum value [0, max).
func getRandomNumber(max int) int {
	rand.Seed(time.Now().UTC().UnixNano())

	return rand.Intn(max)
}

// Choose a random country from the list.
func getRandomCountry(coutries []Country) Country {
	return coutries[getRandomNumber(len(coutries))]
}

// Checks if the name entered is a translation of the English name of a country.
func translationName(countries []Country, index int, countryName string) bool {
	for _, transalation := range countries[index].Translations {
		if strings.ToUpper(transalation) == strings.ToUpper(countryName) {
			return true
		}
	}

	return false
}

// CAN BE REPLACED BY BINARY SEARCH - WORK ON IT.
// Get the country index from the list.
func getCountryIndex(countries []Country, countryName string) int {
	for index := 0; index < len(countries); index++ {
		if strings.ToUpper(countries[index].CountryName) == strings.ToUpper(countryName) ||
			translationName(countries, index, countryName) {
			return index
		}
	}

	return -1
}

// Formats all information related to languages ​​and currencies.
func createStrings(selectedCountry Country) (languages, currencies string) {
	for i := 0; i < len(selectedCountry.Languages); i++ {
		languages += selectedCountry.Languages[i].Language + " | "
	}

	for i := 0; i < len(selectedCountry.Currencies); i++ {
		currencies += selectedCountry.Currencies[i].CurrencyName + " (" +
			selectedCountry.Currencies[i].CurrencySymbol + ") | "
	}

	return
}

// Returns information related to a country.
func getCountryInfo(countries []Country, countryName string) string {
	var selectedCountry Country

	if countryName == "" {
		selectedCountry = getRandomCountry(countries)
	} else {
		index := getCountryIndex(countries, countryName)

		if index == -1 {
			return `You must have entered the wrong country name. Try again or type /help for more information.`
		}

		selectedCountry = countries[index]
	}

	languages, currencies := createStrings(selectedCountry)

	return "Here are the informations about: " + selectedCountry.CountryName + "\nRegion: " +
		selectedCountry.Region +
		"\nCapital: " + selectedCountry.Capital + "\nOfficial languages: " +
		languages + "\nCurrencies: " + currencies
}

// Checks if a value is not yet in the list.
func notInList(randomCountry Country, gameList []Country) bool {
	for i := 0; i < len(gameList); i++ {
		if randomCountry.CountryName == gameList[i].CountryName {
			return false
		}
	}

	return true
}

// Remove an index from array.
func removeIndex(array []Country, index int) []Country {
	return append(array[:index], array[index+1:]...)
}

// Reset game variables.
func restartGame() {
	finalList = finalList[:0]
	gameIndex = -1
}

// Draw a game mode (currently two modes) and configure everything so that the user can play.
func countryGame(countries []Country) string {
	var gameList []Country
	var country Country
	var menssage string
	var flagControl bool

	restartGame()

	gameMode := getRandomNumber(2)

	country = getRandomCountry(countries)

	gameList = append(gameList, country) // Temp GameList.

	for i := 0; i < 3; i++ {
		randomCountry := getRandomCountry(countries)

		if notInList(randomCountry, gameList) {
			gameList = append(gameList, randomCountry)
		}
	}

	for len(gameList) > 0 {
		randomNumber := getRandomNumber(len(gameList))

		finalList = append(finalList, gameList[randomNumber])
		gameList = removeIndex(gameList, randomNumber)

		if randomNumber == 0 && !flagControl {
			gameIndex = len(finalList)
			flagControl = true
		}
	}

	menssage += "Okay, here we go!\n"

	if gameMode == 0 {
		menssage += fmt.Sprintf("What is the capital of %v:\n", country.CountryName)

		for i := 0; i < len(finalList); i++ {
			menssage += fmt.Sprintf("%v) %v\n", i+1, finalList[i].Capital)
		}
	} else {
		menssage += fmt.Sprintf("%v is the capital of which country:\n", country.Capital)

		for i := 0; i < len(finalList); i++ {
			menssage += fmt.Sprintf("%v) %v\n", i+1, finalList[i].CountryName)
		}
	}

	menssage += "Type: /answer number!\n"

	return menssage
}

// Compares user response to expected response.
func checkAnswer(userAnswer string) string {
	var menssage string
	userGuess, err := strconv.Atoi(userAnswer)

	if err != nil || gameIndex == -1 || userGuess < 1 || userGuess > 4 {
		return "Something went wrong, try again."
	}

	if userGuess == gameIndex {
		menssage += "Right!\nIt looks like we have an expert from the countries around here!\nWhy don't you try again?"
	} else {
		menssage += fmt.Sprintf("Unfortunately you were wrong :c\nThe right answer was (%v - %v) and not (%v - %v)",
			finalList[gameIndex-1].CountryName, finalList[gameIndex-1].Capital,
			finalList[userGuess-1].CountryName, finalList[userGuess-1].Capital)
	}

	/*
		IF A REMOVE THIS I CAN CREATE A NEW GAME MODE WHER THER QUESTIONS ARE ACCUMULATING.
	*/

	restartGame()

	return menssage
}

// Initial command.
func startBot() string {
	return `Hey, nice to meet you!
	I am the Geo Bot and I exist to help you to learn more about our world!
	If is your first time type /help to see all the commands possibilities
	I currently can't stand many commands, but my dad is working to improve me and so I can take even more knowledge to the world!`
}

// Help command.
func helpBot() string {
	return `The commands that I know are:
	/list -> List some available country names.
	/info "county name" -> Use this command if you want to know something about a country. If you dont type any country name I will get a random for you.
	/play -> Why not to practice what your learn. When you use this command I will ask random questions and you answer with /answer "your answer".`
}

// Standard response for unknown commands.
func defaultAnswer() string {
	return `I don't know this command :c. Why don't you try /start or /help?`
}

// List all countries available via the bot.
func listCountries(coutries []Country) string {
	var coutriesListNames string = "Here is the list of countries I know:\n"

	for i := 0; i < len(coutries); i++ {
		coutriesListNames += coutries[i].CountryName + "\n"
	}

	return coutriesListNames
}
