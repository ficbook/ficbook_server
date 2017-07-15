import websocket
import json

ws = websocket.WebSocket()
ws.connect("ws://localhost:8080/entry")
dd = {
	"author":"nick",
	"body":"test"
}
print(json.dumps(dd))
ws.send(json.dumps(dd)) 
result = ws.recv()
print(result)

ws.close()

