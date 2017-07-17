import websocket
import json

ws = websocket.WebSocket()
ws.connect("ws://localhost:8080/entry")
dd = {
	"type":"nick",
	"action":"test",
	"test":"qwerty"
}
ws.send(json.dumps(dd)) 

while 1:
	result = ws.recv()
	print(result)
ws.close()

