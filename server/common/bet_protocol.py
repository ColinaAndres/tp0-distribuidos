
from common.socket import Socket
from common.utils import Bet

class BatchProcessingError(Exception):
    """
    Custom exception for errors during batch processing of bets.
    """
    def __init__(self, message, batch_count):
        super().__init__(message)
        self.batch_count = batch_count

class BetProtocol:
    """
    Protocol to manage serialization, deserialization sending and reception of bets.
    """
    LENGTH_OF_TYPE = 1
    LENGTH_PREFIX_SIZE = 2
    BET_DIVIDER = ','
    BATCH_DIVIDER = '|'
    BYTE_ORDER = 'big'
    AMOUNT_OF_BET_ATRIBUTES = 6
    OK_CODE = b'\x01'
    ERROR_CODE = b'\x00'
    FINALIZATION_BYTE = b'\x00'
    BATCH_SENDING_BYTE = b'\x01'

    def __init__(self, client_socket):
        self._client_socket = Socket(client_socket)

    def receive_bet(self) -> Bet:
        """
        Receives a bet from the client socket and deserializes it.
        """
        length_prefix_bytes = self._client_socket.receive_all(self.LENGTH_PREFIX_SIZE)
        message_length = int.from_bytes(length_prefix_bytes, byteorder=self.BYTE_ORDER)
        data = self.__decode_to_utf8(self._client_socket.receive_all(message_length))
        return self.__deserialize_bet(data)
    
    def receive_bets(self) -> list[Bet]:
        """
        Receives multiple bets from the client socket and deserializes them.
        """
        length_prefix_bytes = self._client_socket.receive_all(self.LENGTH_PREFIX_SIZE)
        if not length_prefix_bytes:
            return []
        
        batch_length = int.from_bytes(length_prefix_bytes, byteorder=self.BYTE_ORDER)
        batch_data = self.__decode_to_utf8(self._client_socket.receive_all(batch_length))
        return self.__deserialize_bets(batch_data)
    
    def receive_request(self):
        type_of_request = self._client_socket.receive_all(self.LENGTH_OF_TYPE)
        match type_of_request:
            case self.FINALIZATION_BYTE:
                pass
            case self.BATCH_SENDING_BYTE:
                pass
            case _ :
                return None
    
    def send_confirmation(self):
        """
        Sends a confirmation message to the client socket.
        """
        self._client_socket.send_all(self.OK_CODE)

    def send_error(self):
        """
        Sends an error message to the client socket.
        """
        self._client_socket.send_all(self.ERROR_CODE)

    def close(self):
        """
        Closes the client socket.
        """
        self._client_socket.close()

    def __decode_to_utf8(self, data) -> str:
        """
        Decodes bytes data to a UTF-8 string.
        """
        return data.decode('utf-8')

    def __deserialize_bets(self, batch_data) -> list[Bet]:
        """
        Deserializes a batch of bet data into a list of Bet objects.
        """
        try:
            stringed_bets = batch_data.split(self.BATCH_DIVIDER)
            return list(map(self.__deserialize_bet, stringed_bets))
        except ValueError as e:
            raise BatchProcessingError(f'Error deserializing bets: {e}', batch_count=len(stringed_bets))

    def __deserialize_bet(self, bet_data) -> Bet :
        """
        Deserializes a single bet data string into a Bet object.
        Raise Value error if a bet is not completed
        """
        bet_atributes = bet_data.split(self.BET_DIVIDER)
        if len(bet_atributes) != self.AMOUNT_OF_BET_ATRIBUTES:
            raise ValueError(f'invalid bet information: {bet_atributes}')
        agency, first_name, last_name, document, birthdate, number = bet_atributes
        return Bet(agency, first_name, last_name, document, birthdate, number)
