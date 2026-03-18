import socket

class Socket:
    def __init__(self, socket):
        self._skt = socket
    
    def receive_all(self, amount_to_receive):
        received_data = b''
        while len(received_data) < amount_to_receive:
            chunk = self._skt.recv(amount_to_receive - len(received_data))
            if not chunk:
                raise ConnectionError("Socket connection closed before receiving all data")    
            received_data += chunk
        return received_data
    
    def send_all(self, data):
        self._skt.sendall(data)
