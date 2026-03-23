import socket
import logging

from common.bet_protocol import BatchProcessingError, BetProtocol
from common.utils import has_won, load_bets, store_bets
from common.helpers import close_socket
from common.agency_session import AgencySession
from common.lottery_central import Lottery_central

class Server:
    def __init__(self, port, listen_backlog, total_agencies):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._agency_sessions = []
        self._lottery_central = Lottery_central(total_agencies)
    
    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        self._running = True
        while self._running and not self._lottery_central.draw_done():
            client_sock = self.__accept_new_connection()
            if client_sock:
                session = AgencySession(BetProtocol(client_sock), self._lottery_central)
                self._agency_sessions.append(session)
                session.receive_bets()
                self.__remove_stopped_sessions()

        if self._running:
            for session in self._agency_sessions:
                session.receive_winners_request()
        
        self.__cleanup()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made or
        a an exception occurs. Then connection created is printed and returned
        """

        # Connection arrived
        try:
            logging.info('action: accept_connections | result: in_progress')
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c

        except OSError as e:
            logging.error(f'action: accept_connections | result: fail | error: {e}')
            return None
        
    def __remove_stopped_sessions(self):
        """
        Remove stopped sessions
        Function iterates over the list of agency sessions and removes those that are stopped
        """
        self._agency_sessions = list(filter(lambda session: not session.close_if_stoped(), self._agency_sessions))
        
    def __cleanup(self):
        """
        Cleanup the server
        """
        for session in self._agency_sessions:
            session.stop()
        self._agency_sessions = []

    def graceful_shutdown(self, _signum, _frame):
        """
        Gracefully shutdown the server

        When the server has its _running flag set to false, the server socket 
        is shutdown and closed. If the client socket is still open, it is also shutdown and closed

        """
        logging.info('action: graceful_shutdown | result: in_progress')
        self._running = False
        self.__cleanup()
        close_socket(self._server_socket, "server")
        logging.info('action: graceful_shutdown | result: success')


        