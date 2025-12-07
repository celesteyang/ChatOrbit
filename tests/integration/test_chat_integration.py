# test_chat_integration.py
import asyncio
import websockets
import requests
import uuid
import json

AUTH_URL = "http://localhost:8089"
CHAT_URL = "http://localhost:8088"
USER_URL = "http://localhost:8087"
WS_BASE = "ws://localhost:8088/ws/chat"


async def connect(token, room_id):
    url = f"{WS_BASE}?token={token}&room_id={room_id}"
    ws = await websockets.connect(url)
    return ws


async def register(email, password, username):
    """
    Register a new user
    """
    res = requests.post(
        f"{AUTH_URL}/register",
        json={"email": email, "password": password, "username": username},
    )
    if res.status_code == 200:
        return True
    return False


async def login(email, password):
    """
    Login the user and return access token
    """
    res = requests.post(
        f"{AUTH_URL}/login",
        json={"email": email, "password": password},
    )
    token = res.json()["token"]
    return token


async def create_room(room_id):
    """
    Create a new chat room
    """
    res = requests.post(
        f"{CHAT_URL}/chat/rooms",
        json={"room_id": room_id},
    )
    if res.status_code == 200:
        return True
    return False


async def test_chat_flow():
    # Generate 2 random user names:
    name1 = "user_" + str(uuid.uuid4())[:8]
    name2 = "user_" + str(uuid.uuid4())[:8]

    # Generate 2 random passwords:
    password1 = str(uuid.uuid4())[:8]
    password2 = str(uuid.uuid4())[:8]

    # Create the corresponding emails:
    email1 = f"{name1}@gmail.com"
    email2 = f"{name2}@gmail.com"

    # Register users
    success = await register(email1, password1, name1)
    assert success, "Registration failed for user 1"

    success = await register(email2, password2, name2)
    assert success, "Registration failed for user 2"

    # Login as 2 different users
    token1 = await login(email1, password1)
    token2 = await login(email2, password2)

    # Create a new room
    room_id = "room_" + str(uuid.uuid4())[:8]
    success = await create_room(room_id)
    assert success, "Room creation failed"

    # Join same room
    ws1 = await connect(token1, room_id)
    ws2 = await connect(token2, room_id)

    # Define the message
    message_text = "Hello!"

    # send using ws1
    await ws1.send(json.dumps({
        "room_id": room_id,
        "content": message_text
    }))

    # receive using ws2
    response = await ws2.recv()
    message = json.loads(response)

    print(f"Message: {message}")

    assert message["content"] == message_text

    await ws1.close()
    await ws2.close()


if __name__ == "__main__":
    asyncio.run(test_chat_flow())
