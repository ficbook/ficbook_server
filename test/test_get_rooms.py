import websocket
import json
import settings

ws = websocket.WebSocket()
ws.connect("ws://localhost:8080")
dd = {
	"type":"authorize",
	"login":settings.login,
	"password":settings.password
}
ws.send(json.dumps(dd)) 
result = ws.recv()
print(result)

dd = {
    "type":"rooms",
    "action":"get"
}
ws.send(json.dumps(dd)) 
result = ws.recv()
print(result)

ws.close()
