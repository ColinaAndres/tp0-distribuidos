
import logging

from common.command import Command
from common.bet_protocol import BatchProcessingError


class AgencySession:
    """
    Class that represents a session with an agency
    """
    def __init__(self, protocol, lottery_central, agency_id=None):
        self.agency_id = agency_id
        self._protocol = protocol
        self._lottery_central = lottery_central

    def receive_bets(self):
        """
        Receives bets from the client until a finalization command is received.
        """
        try:
            while True:
                request = self._protocol.receive_request()
                if request is None:
                    break
                done = request.execute(self._lottery_central, self)
                if done:  # FinalizationCommand retorna True
                    break
                self._protocol.send_confirmation()
        except BatchProcessingError as e:
            self._protocol.send_error()
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {e.batch_count}")

    def receive_winners_request(self):
        """
        Waits for client to request winners.
        """
        request = self._protocol.receive_request()
        if request is None:
            logging.error(f"action: send_winners_phase | result: fail | reason: client disconnected")
            return
        request.execute(self._lottery_central, self)

    def send_winners(self, winners):
        """Sends the winners to the protocol."""
        self._protocol.send_winners(winners)

    def stop(self):
        """Closes the protocol connection."""
        self._protocol.close()