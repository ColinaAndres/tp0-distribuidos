
from server.common.utils import Bet


class BetProtocol:
    """
    Protocol to manage serialization, deserialization sending and reception of bets.
    """

    def __init__(self, client_socket):
        self._client_socket = client_socket

    def receive_bet(self) -> Bet:
        """
        Receives a bet from the client socket and deserializes it.
        """
        data = self._client_socket.recv(1024).decode('utf-8')
        agency, first_name, last_name, document, birthdate, number = data.split(',')
        return Bet(agency, first_name, last_name, document, birthdate, number)
    
    def send_confirmation(self)
        """
        Sends a confirmation message to the client socket.
        """
        self._client_socket.sendall(b'\x01')
