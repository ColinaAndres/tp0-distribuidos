from common.helpers import close_socket

class Socket:
    def __init__(self, socket):
        self._skt = socket
    
    def receive_all(self, amount_to_receive):
        '''
        Receives a specific amount of bytes from the socket, ensuring no short reads
        If the connection is closed before receiving all data, 
        raises a ConnectionError.
        '''
        received_data = b''
        while len(received_data) < amount_to_receive:
            chunk = self._skt.recv(amount_to_receive - len(received_data))
            if not chunk and not received_data:
                return None
            if not chunk:
                raise ConnectionError("Socket connection closed before receiving all data")    
            received_data += chunk
        return received_data
    
    def send_all(self, data):
        """
        sends all the data to the socket, ensuring no short writes
        """
        self._skt.sendall(data)

    def close(self):
        """
        Closes the socket connection
        """
        close_socket(self._skt, "client")
        