# Jokes
## Custom Joke Web Service

The service fetches a random name from a Name API and a random Chuck Norris
joke from a Joke API. The service then replaces Chuck Norris' name with the
random name to create a custom joke and serves it on port 5000.

To run the service, clone or download this repository into you Go workspace. 
```
$ git clone https://github.com/NEPDAVE/jokes.git
```

CD into the `jokes` directory and build the binary
```
$ cd jokes
$ go build
```

Run the program 
```
./jokes
```

To get a joke run `curl localhost:5000` in another Terminal or in your web browser go to `localhost:5000`
