# Kshoplist-server
Server-side implementation of [gorilla-websocket](https://github.com/gorilla/websocket) for [kshoplist](https://github.com/mkeiji/kshoplist) 

### Start dev database
`make mysql`

### Start local db client (optional)
`make adminer`
> NOTE: use a client like dbeaver / datagrip if you prefer

### Start the server
`go run main.go`
 
### Test local connection to websocket (with browser)
http://localhost:8081/list/1

### Connect to websocket
ws://localhost:8081/ws/1
> NOTE: over ssl you need to use `wss://`

### Add db migration
create two files following the sequece number and format (one for up and one for down)
where `up` do the migration and `down` undo it.
<br>
OR:
use the `cli` with:
```bash
migrate create -ext sql -dir database/migrations -seq name_of_the_file
```
NEXT: fill both files with the desired `sql`
