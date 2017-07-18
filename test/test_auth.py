import websocket
import json

ws = websocket.WebSocket()
ws.connect("ws://localhost:8080/entry")
dd = {
	"type":"authorize",
	"login":"ficbook_bot",
	"password":"Ficbook_bot12348"
}
ws.send(json.dumps(dd)) 
result = ws.recv()
print(result)
ws.close()

