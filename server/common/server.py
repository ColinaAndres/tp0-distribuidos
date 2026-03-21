import socket
import logging

from common.bet_protocol import BatchProcessingError, BetProtocol
from common.utils import has_won, load_bets, store_bets
from common.helpers import close_socket
from common.agency_session import AgencySession

class Server:
    def __init__(self, port, listen_backlog, total_agencies):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._agency_sessions = []
        self._total_agencies = total_agencies
        self._doned_agencies = 0
        self._session_active = False
        self._winners = []
    
    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        self._running = True
        while self._running and len(self._agency_sessions) < self._total_agencies:
            client_sock = self.__accept_new_connection()
            if client_sock:
                session = AgencySession(len(self._agency_sessions) + 1, BetProtocol(client_sock))
                self._agency_sessions.append(session)
                
                # Procesa hasta que llega el FinalizationCommand
                self.__handle_client_connection(session)

        if self._running and self._doned_agencies == self._total_agencies:
            logging.info(f"action: sorteo | result: success")
            self.__distribute_winners()

    def graceful_shutdown(self, _signum, _frame):
        """
        Gracefully shutdown the server

        When the server has its _running flag set to false, the server socket 
        is shutdown and closed. If the client socket is still open, it is also shutdown and closed

        """
        logging.info('action: graceful_shutdown | result: in_progress')
        self._running = False
        close_socket(self._server_socket, "server")
        for session in self._agency_sessions:
            session.stop()
        logging.info('action: graceful_shutdown | result: success')

    def __handle_client_connection(self, agency_session):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            self._session_active = True
            while self._running and self._session_active:
                client_request = agency_session.receive_request()
                # TODO: ver de tirar excepcion en el protocolo si se recibe algo no esperado o no
                if client_request:
                    client_request.execute(self)
        except BatchProcessingError as e:
            agency_session.send_error()
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {e.batch_count}")
            agency_session.stop()
        except (OSError, ConnectionError) as e:
            logging.error(f"action: receive_message | result: fail | error: {e}")
            agency_session.stop()

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
        
    def __distribute_winners(self):
        """
        Distribute winners to agencies, every agency should be asking for winners
        """
        self._winners = list(filter(has_won, load_bets()))
        for session in self._agency_sessions:
            self.__handle_client_connection(session)

    def process_bets(self, bets, agency):
        """
        Handle the storing of bets
        """
        store_bets(bets)
        agency.send_confirmation()
        logging.info(f"action: apuesta_recibida | result: success | cantidad: {len(bets)}")

    def finalize_reception_of_bets(self):
        """
        Handle the finalization of the reception of bets
        """
        self._session_active = False
        self._doned_agencies += 1

    def send_winners(self, agency):
        """
        Handle the sending of winners to the client
        """
        agency_winners = list(filter(lambda bet: bet.agency_id == agency._agency_id, self._winners))
        agency.send_winners(agency_winners)
    