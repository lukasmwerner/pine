![Logo](./images/Pine_Logo_Dark.png#gh-dark-mode-only)
![Logo](./images/Pine_Logo_Light.png#gh-light-mode-only)
An all natural tree based router from Oregon
---
After seeing the unfortunate news that [gorilla/mux](https://github.com/gorilla/mux) was archived I needed a new router. After not seeing anything that felt quite like that I decided to build my own.

## Usage
To install run `go get github.com/lukasmwerner/pine`

Example Usage for the router:
```go
p := pine.New()
p.Handle("/hello/{name}", func(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf("Hello, %s", pine.Var(r, "name"))))
})

log.Fatalln(http.ListenAndServe(":8080", p))
```

### Handler Routing
All handlers are parsed in first to last definition order (e.g. if you have a variable in the folder defined first that will run first)


### Middlewares
Gorilla/Mux style middlewares are supported by the type `pine.MiddlewareFunc` which is just a `func(http.Handler) http.Handler` under the hood.

Usage:
```go
p := pine.New()
p.Use(func(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// make sure to call `next` otherwise the handler will not get called
		next.ServeHTTP(w, r)
	})
})
```
Note: Middlewares are run in order of being added

There are a few pre-made middlewares that are defined in the `middlewares` package such as:
* `middlewares.HTTPLogger` Logs the HTTP Request/Responses in the following format: `2023/04/02 12:48:37 host: localhost:8080 method: GET uri: / status: 200 Ok`
* `middlewares.STolinskiTiming` Puts the requests into different time buckets such based on the (slow, middle) durations passed. Inserts a `X-Duration` header to based on those timing buckets


## Scope
* [x] routing based on paths
* [x] variables on paths
* [x] middlewares


### Soundtrack while developing this project
* Firewatch Soundtrack - Chris Remo
* Discovery - Daft Punk
* Tron: Legacy - Daft Punk


## Contributions
While I appreciate any contributions, this is still MY router so if it is not in the spirit of the API I want to use then feel free to make your own fork modify it however you would like.
