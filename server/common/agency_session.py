import logging
from common.bet_protocol import BatchProcessingError

class AgencySession:
    """
    Class that represents a session with an agency
    """
    def __init__(self, protocol, lottery_central, agency_id=None):
        self.agency_id = agency_id
        self._protocol = protocol
        self._lottery_central = lottery_central
        self._running = True

    def receive_bets(self):
        """
        Receives bets from the client until a finalization command is received.
        """
        try:
            while self._running:
                request = self._protocol.receive_request()
                if request is None:
                    break
                done = request.execute(self._lottery_central, self)
                if done:  # FinalizationCommand retorna True
                    break
                self._protocol.send_confirmation()
        except BatchProcessingError as e:
            self.__batch_processing_error_handler(self, e)
        except ConnectionError:
            self.__connection_error_handler()

    def receive_winners_request(self):
        """
        Waits for client to request winners.
        """
        try:
            request = self._protocol.receive_request()
            if request is None:
                self.__connection_error_handler()
                return
            request.execute(self._lottery_central, self)
        except ConnectionError:
            self.__connection_error_handler()

    def send_winners(self, winners):
        """Sends the winners to the protocol."""
        try:
            self._protocol.send_winners(winners)
        except ConnectionError:
            self.__connection_error_handler()

    def stop(self):
        """Closes the protocol connection."""
        self._running = False
        self._protocol.close()

    def close_if_stoped(self) -> bool:
        """
        Closes the protocol connection if the session is stopped.
        Returns True if the session is stopped, False otherwise.
        """
        if not self._running:
            self.stop()
        return not self._running

    def __batch_processing_error_handler(self, error):
        """
        Handles BatchProcessingError by logging the error and closing the protocol connection.
        """
        try:
            self._protocol.send_error()
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {error.batch_count}")
        except ConnectionError:
            self.__connection_error_handler()
    
    def __connection_error_handler(self):
        """
        Handles ConnectionError by logging the error.
        """
        self._running = False
        logging.error(f"action: apuesta_recibida | result: fail | reason: client disconnected")