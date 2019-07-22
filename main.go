package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
)

var (
	wg     sync.WaitGroup
	logger *log.Logger
	client = &http.Client{}
)

func main() {
	//creating the logger object
	logFile, err := os.Create("log.txt")
	defer logFile.Close()

	if err != nil {
		logger.Println(err)
	}

	logger = log.New(logFile, "jokes ", log.LstdFlags|log.Lshortfile)
	logger.Println("Program Start!")

	nameByte, err := GetName()
	name, _ := Name{}.UnmarshalName(nameByte)
	fmt.Println(name)

	jokeByte, err := GetJoke()
	joke, _ := Joke{}.UnmarshalJoke(jokeByte)
	fmt.Println(joke)
}

/*
***************************
name struct methods and funtions
***************************
*/

//Name reprents the data structure returned by the uinames.com Web Service
type Name struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Gender  string `json:"gender"`
	Region  string `json:"region"`
}

//UnmarshalName unmarshals the JSON returned by the uinames.com Web Service
func (n Name) UnmarshalName(nameByte []byte) (*Name, error) {

	err := json.Unmarshal(nameByte, &n)

	if err != nil {
		logger.Println(err)
		return &n, err
	}

	return &n, nil
}

//GetName makes a Get HTTP request to the uinames.com Web Service and returns
//a random name and an error
func GetName() ([]byte, error) {
	nameURL := "http://uinames.com/api/"
	req, err := http.NewRequest("GET", nameURL, nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "Keep-Alive")

	if err != nil {
		logger.Println(err)
		return []byte{}, err
	}

	resp, err := client.Do(req)

	if err != nil {
		logger.Println(err)
		return []byte{}, err
	}

	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		logger.Println(err)
		return []byte{}, err
	}

	return bytes, err
}

/*
***************************
joke struct methods and functions
***************************
*/

//Joke reprents the data structure returned by the api.icndb.com Web Service
type Joke struct {
	Type  string `json:"type"`
	Value Value  `json:"value"`
}

//Value reprents the embedded Value field of the Joke struct
type Value struct {
	ID         int      `json:"id"`
	Joke       string   `json:"joke"`
	Categories []string `json:"categories"`
}

//UnmarshalJoke unmarshals the JSON returned by the api.icndb.com Web Service
func (j Joke) UnmarshalJoke(jokeByte []byte) (*Joke, error) {

	err := json.Unmarshal(jokeByte, &j)

	if err != nil {
		logger.Println(err)
		return &j, err
	}

	return &j, nil
}

//GetJoke makes a Get HTTP request to api.icndb.com Web Service
//and returns a random Chuck Norris joke and an error
func GetJoke() ([]byte, error) {
	jokeURL := "http://api.icndb.com/jokes/random"

	//ABOUT NAME API:
	//possible to add extra query parameter for example "limitTo=nerdy"
	//by not including the limitTo parameter a Chuck Norris joke is returned
	//also possible to include "firstName=Larry" and "lastName=Smith" parameters

	req, err := http.NewRequest("GET", jokeURL+"/random", nil)

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Connection", "Keep-Alive")

	if err != nil {
		logger.Println(err)
		return []byte{}, err
	}

	resp, err := client.Do(req)

	if err != nil {
		logger.Println(err)
		return []byte{}, err
	}

	defer resp.Body.Close()

	bytes, _ := ioutil.ReadAll(resp.Body)

	if err != nil {
		logger.Println(err)
		return []byte{}, err
	}

	return bytes, err
}
