
from common.command import Command


class AgencySession:
    """
    Class that represents a session with an agency
    """
    def __init__(self, protocol, agency_id=None):
        self.agency_id = agency_id
        self._protocol = protocol

    def receive_request(self) -> Command:
        """Receives a request from protocol."""
        request = self._protocol.receive_request()
        return self._protocol.receive_request()
    
    def send_error(self):
        """Sends an error message to the protocol."""
        self._protocol.send_error()

    def send_confirmation(self):
        """Sends a confirmation message to the protocol."""
        self._protocol.send_confirmation()

    def send_winners(self, winners):
        """Sends the winners to the protocol."""
        self._protocol.send_winners(winners)

    def stop(self):
        """Closes the protocol connection."""
        self._protocol.close()