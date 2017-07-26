import websocket
import json
import settings

ws = websocket.WebSocket()
ws.connect("ws://192.168.0.105:8080")
dd = {
	"type":"autorize",
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

join_room = {"type":"room","action":"join","room_name":"Nikita"}
ws.send(json.dumps(join_room)) 
result = ws.recv()
print(result)

join_room = {"type":"chat","action":"get","subject":"history", "room_name":"Nikita"}
ws.send(json.dumps(join_room)) 
result = ws.recv()
print(result)

ws.close()
