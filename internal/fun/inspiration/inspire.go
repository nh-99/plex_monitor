package inspiration

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"time"
)

// Quotes is a list of Quote types
type Quotes struct {
	Quotes []Quote
}

// Quote stores info about an inspirational quote
type Quote struct {
	QuoteText   string
	QuoteAuthor string
}

func openJSONFile(fileName string) []byte {

	// opening our json file
	jsonFile, err := os.Open(fileName)
	// if there was an error while opening our json file we print it
	if err != nil {
		fmt.Println(err)
	}

	//fmt.Println("Successfully Opened users.json")

	// reading our json file as an byte array
	byteValue, _ := io.ReadAll(jsonFile)

	return byteValue
}

func parseJSON(byteValue []byte) Quotes {
	// initializing our quote array
	var quotes Quotes

	// we unmarshal our byteArray which contains our
	// jsonFile's content and assign it into 'quotes' which we defined above
	// unmarshal takes our json data as bytearray and we give it a pointer to our struct
	json.Unmarshal(byteValue, &quotes)

	return quotes
}

func getRandomQuoteAndPrint(quotes Quotes) (ret string) {
	// randomly picking a quote

	// randomly generating seed very fast, which the rand funtion uses
	// default of seed is one 1 and when we generate Intn we will get the same output every time
	// so that's why we changed the seed so we get random output every time we run the program
	rand.Seed(time.Now().UnixNano())

	// Intn works by generating a number between 0 and N
	randomQuote := quotes.Quotes[rand.Intn(len(quotes.Quotes)-1)]
	ret = randomQuote.QuoteText

	// for cases when QuoteAuthor is empty
	if randomQuote.QuoteAuthor != "" {
		ret = fmt.Sprintf("%s - <em>%s</em>", ret, randomQuote.QuoteAuthor)
	} else {
		ret = fmt.Sprintf("%s - <em>Unkown</em>", ret)
	}

	return ret
}

// GetInspirationalQuote loads the data and gets a random quote
func GetInspirationalQuote() string {
	var byteValue = openJSONFile("./data/inspirational_quotes.json")
	var quotes = parseJSON(byteValue)
	return getRandomQuoteAndPrint(quotes)
}
