package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	wg          sync.WaitGroup
	logger      *log.Logger
	client      = &http.Client{}
	customJokes = []string{} //slice to store customJokes
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

	//web server logic for jokes service
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {

		//creating a joke
		joke := CustomJoke{}.NewCustomJoke()

		//checking the joke and serving a 500 if the joke == "", IE there is no joke
		if joke == "" {
			logger.Printf("Error - Internal Server Error 500: %s\n", err)
			fmt.Fprintln(w, "Internal Server Error 500 - Please Try Again")
		}

		//serving joke
		fmt.Fprintf(w, "%s\n", joke)
		fmt.Println(joke)
		fmt.Println("")

	})

	//fs := http.FileServer(http.Dir("static/"))
	//	http.Handle("/static/", http.StripPrefix("/static/", fs))

	//http server with some custom timeouts set
	s := &http.Server{
		Addr:         ":5000",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  10 * time.Second,
	}

	s.ListenAndServe()

}

/*
***************************
custom joke struct, methods, and cached custom jokes
***************************
*/

//GetCachedJoke retrieves a random joke from the cachedJokes slice, a cached
//joke can be served in the event of a 500 response from Names API or the Jokes API
func GetCachedJoke() string {
	l := len(customJokes)

	//handling edge cases where slice is not fully formed to work with random index
	//selection, IE a slice with 0 or 1 elements
	if l == 0 {
		return "no jokes right now - try again"
	} else if l == 1 {
		return customJokes[0]
	}

	index := rand.Intn(len(customJokes) - 1)
	logger.Println("serving cached joke")
	return customJokes[index]
}

//CustomJoke contains the data and methods to create a custom joke
type CustomJoke struct {
	Name      string
	Surname   string
	Joke      string
	NameError error
	JokeError error
}

//CustomJoke creates a custome joke
func (c CustomJoke) NewCustomJoke() string {
	wg.Add(2)
	//using goroutines to allow waiting on multiple requests concurrently
	//functionality should increase throughput
	go c.NewName()
	go c.NewJoke()
	wg.Wait()

	if c.NameError != nil {
		//logger.Printf("Name Error: %s\n", c.NameError)
		return GetCachedJoke()
	}

	if c.JokeError != nil {
		//logger.Printf("Joke Error: %s\n", c.JokeError)
		return GetCachedJoke()
	}

	//replacing "&quot;" strings from joke with "
	c.Joke = strings.Replace(c.Joke, "&quot;", "\"", -1)

	//replacing "Chuck" strings from with c.Name
	c.Joke = strings.Replace(c.Joke, "Chuck", c.Name, -1)

	//replacing "Norris" strings from with c.Surname
	c.Joke = strings.Replace(c.Joke, "Norris", c.Surname, -1)

	//adding newly created custom joke to slice of cached custom jokes
	customJokes = append(customJokes, c.Joke)
	return c.Joke
}

//NewName gets a name and unmarshals it into a struct
func (c *CustomJoke) NewName() {
	defer wg.Done()

	bytes, err := GetName()

	if err != nil {
		c.NameError = err
		return
	}

	name, err := Name{}.UnmarshalName(bytes)

	if err != nil {
		c.NameError = err
		return
	}

	c.Name = name.Name
	c.Surname = name.Surname
}

//NewJoke gets a joke and unmarshals it into a struct
func (c *CustomJoke) NewJoke() {
	defer wg.Done()

	bytes, err := GetJoke()

	if err != nil {
		c.JokeError = err
		return
	}

	joke, err := Joke{}.UnmarshalJoke(bytes)

	if err != nil {
		c.JokeError = err
		return
	}

	c.Joke = joke.Value.Joke
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
	//capturing panic from any corrupt input
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("Error - UnmarshalName: %v", r)
			logger.Println(err)
		}
	}()

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

	if resp.StatusCode < 200 || resp.StatusCode > 200 {
		err = fmt.Errorf("uimanes.com response code: %d - %s", resp.StatusCode,
			http.StatusText(resp.StatusCode))
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
	//capturing panic from any corrupt input
	defer func() {
		if r := recover(); r != nil {
			err := fmt.Errorf("Error - UnmarshalJoke: %v", r)
			logger.Println(err)
		}
	}()

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

	if resp.StatusCode < 200 || resp.StatusCode > 200 {
		err = fmt.Errorf("api.icndb.com response code: %d - %s", resp.StatusCode,
			http.StatusText(resp.StatusCode))
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
