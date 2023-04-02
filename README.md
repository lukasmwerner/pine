# Pine Router
An all natural tree based router from Oregon

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



## Scope
* [x] routing based on paths
* [x] variables on paths
* [ ] middlewares


### Soundtrack while developing this project
* Firewatch Soundtrack - Chris Remo
* Discovery - Daft Punk
* Tron: Legacy - Daft Punk


## Contributions
While I appreciate any contributions, this is still MY router so if it is not in the spirit of the API I want to use then feel free to make your own fork modify it however you would like.
