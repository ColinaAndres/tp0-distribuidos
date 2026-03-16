import socket
import logging
import signal


class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._timeout_time = 0.5
        self._server_socket.settimeout(self._timeout_time)

        #register signal handler for SIGTERM
        signal.signal(signal.SIGTERM, self.__handle_sigterm)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communucation
        finishes, servers starts to accept new connections again
        """

        self._running = True
        while self._running:
            self._client_sock = self.__accept_new_connection()
            if self._client_sock is not None:
                self.__handle_client_connection()

        self.__gracefull_shutdown()

    def __gracefull_shutdown(self):
        """
        Gracefully shutdown the server

        When the server has its _running flag set to false, the server socket 
        is shutdown and closed. If the client socket is still open, it is also shutdown and closed

        """
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()

        # Shutsdown client socket if it is still open, possibly necesary on future changes
        if self._client_sock is not None:
            self._client_sock.shutdown(socket.SHUT_RDWR)
            self._client_sock.close()

        return 

    def __handle_client_connection(self):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            # TODO: Modify the receive to avoid short-reads
            msg = self._client_sock.recv(1024).rstrip().decode('utf-8')
            addr = self._client_sock.getpeername()
            logging.info(f'action: receive_message | result: success | ip: {addr[0]} | msg: {msg}')
            # TODO: Modify the send to avoid short-writes
            self._client_sock.send("{}\n".format(msg).encode('utf-8'))
        except OSError as e:
            logging.error("action: receive_message | result: fail | error: {e}")
        finally:
            self._client_sock.close()

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made or
        a timeout occurs. Then connection created is printed and returned
        """

        # Connection arrived
        try:
            logging.info('action: accept_connections | result: in_progress')
            c, addr = self._server_socket.accept()
            logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
            return c

        except socket.timeout:
            return None

    def __handle_sigterm(self, _signum, _frame):
        """ 
        Handle SIGTERM signal

        When SIGTERM signal is received, the server running flag is set to false
        """

        logging.info('action: receive_sigterm | result: success')
        self._running = False