import websocket
import json

ws = websocket.WebSocket()
ws.connect("ws://localhost:8080")
dd = {
	"type":"authorize",
	"login":"login",
	"password":"password"
}
ws.send(json.dumps(dd)) 

while 1:
	result = ws.recv()
	print(result)
ws.close()

