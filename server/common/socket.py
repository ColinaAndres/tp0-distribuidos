import socket

class Socket:
    def __init__(self, socket):
        self._skt = socket
    
    def receive_all(self, amount_to_receive):
        received_data = b''
        while len(received_data) < amount_to_receive:
            chunk = self._skt.recv(amount_to_receive - len(received_data))
            if not chunk:
                break
            received_data += chunk