
class BetProtocol:
    """
    Protocol to manage serialization, deserialization sending and reception of bets.
    """

    def __init__(self, client_socket):
        self._client_socket = client_socket
