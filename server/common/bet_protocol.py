
from common.socket import Socket
from common.utils import Bet

LENGTH_PREFIX_SIZE = 2
DIVIDER = ','
BYTE_ORDER = 'big'
class BetProtocol:
    """
    Protocol to manage serialization, deserialization sending and reception of bets.
    """

    def __init__(self, client_socket):
        self._client_socket = Socket(client_socket)

    def receive_bet(self) -> Bet:
        """
        Receives a bet from the client socket and deserializes it.
        """
        length_prefix_bytes = self._client_socket.receive_all(LENGTH_PREFIX_SIZE)
        message_length = int.from_bytes(length_prefix_bytes, byteorder=BYTE_ORDER)
        data = self._client_socket.receive_all(message_length).decode('utf-8')
        agency, first_name, last_name, document, birthdate, number = data.split(DIVIDER)
        return Bet(agency, first_name, last_name, document, birthdate, number)
    
    def send_confirmation(self):
        """
        Sends a confirmation message to the client socket.
        """
        self._client_socket.send_all(b'\x01')

    def close(self):
        """
        Closes the client socket.
        """
        self._client_socket.close()
