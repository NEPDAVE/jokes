# Jokes
## Custom Joke Web Service

The service fetches a random name from a Name API and a random Chuck Norris
joke from a Joke API. The service then replaces Chuck Norris' name with the
random name to create a custom joke and serves it on port 5000.

## To run the service, use one of two options.

### Option One
Clone or download this repository into you Go workspace. 
```
$ git clone https://github.com/NEPDAVE/jokes.git
```

CD into the `jokes` directory build, the binary, and run the program
```
$ cd jokes
$ go build
$./jokes
```
To get a joke run `curl localhost:5000` in another Terminal or in your web browser go to `localhost:5000`

### Option Two
Use go get to fetch the package
```
$ go get github.com/nepdave/jokes
$ cd $HOME/go/bin
$ ./jokes
```
To get a joke run `curl localhost:5000` in another Terminal or in your web browser go to `localhost:5000`


