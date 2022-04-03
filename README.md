# meme-cache


## how to run

in `/cmd/meme-cache`, run `go run main.go`


## how to test
```
http://localhost:8080/put?key=1&value=2

http://localhost:8080/get?key=1

// For master server
// add new cache server
http://localhost:8081/addMember?host=127.0.0.1&port=8081&vid=11231&vid=286555555&vid=3885454534235
// get current membership
http://localhost:8081/getMembers
```
