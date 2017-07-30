import websocket
import json
import settings
import random

ws = websocket.WebSocket()
ws.connect("ws://192.168.0.105:7070")
dd = {
	"type":"autorize",
	"login":settings.login,
	"password":settings.password
}
ws.send(json.dumps(dd)) 
result = ws.recv()
print(result)

while 1:
	join_room = {"type":"chat","action":"send","subject":"message", "room_name":"Test", "message":str(random.random())}
	ws.send(json.dumps(join_room)) 
	result = ws.recv()
	print(result)



ws.close()
