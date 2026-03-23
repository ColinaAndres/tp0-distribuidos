import logging
import socket

def close_socket(skt, socket_name):
        """
        Auxiliar function to close sockets
        """
        if skt is None:
            return
        
        try:
            skt.shutdown(socket.SHUT_RDWR)
            skt.close()
            logging.info(f'action: closing_{socket_name}_socket | result: success')
        except OSError as e:
            logging.error(f'action: closing_{socket_name}_socket | result: fail | error: {e}')
            