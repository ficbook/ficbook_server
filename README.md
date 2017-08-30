#Ficbook websocket server

## How to compile the project?
First of all, we need to download some libraries. To do this, we use the go get:<br>
`go get github.com/jinzhu/gorm"`<br>
`go get github.com/jinzhu/gorm/dialects/mysql`<br>
`go get github.com/ficbook/ficbook_server/chat`<br>
Now you can start the project or compile it!<br><br>

Available flags:
* -config: `go run main.go -config=config.cfg`. Default: config.cfg <br>
Example: `go run main.go -config=myproject/server/config.cfg`
* -database-init: `go run main.go -database-init=true`. Default: false<br>
Added database migration
* -create-room `go run main.go -create-room=RoomName`. Default: ""<br>
Creates a public room named RoomName
<br><br>
Available Language Packs<br>
| Name  | Language |
|-------|----------|
| en    | English  |
| ru    | Russian  |