
from common.command import Command


class AgencySession:
    """
    Class that represents a session with an agency
    """
    def __init__(self, agency_id, protocol):
        self._agency_id = agency_id
        self._protocol = protocol

    def receive_request(self) -> Command:
        """Receives a request from protocol."""
        return self._protocol.receive_request()
    
    def send_error(self):
        """Sends an error message to the protocol."""
        self._protocol.send_error()

    def stop(self):
        """Closes the protocol connection."""
        self._protocol.close()