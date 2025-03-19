import socket
import logging
import signal
import sys
from common.utils import *

class Server:
    def __init__(self, port, listen_backlog):
        # Initialize server socket
        self._server_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._server_socket.bind(('', port))
        self._server_socket.listen(listen_backlog)
        self._current_connection = None

        # Setting the signal handler
        signal.signal(signal.SIGTERM, self.__sigterm_handler)

    def __sigterm_handler(self, signal, frame):
        """
        SIGTERM signal handler

        When the application receives a SIGTERM signal, all the file descriptors 
        (welcoming socket, and client current socket) are closed for a graceful shutdown
        """
        self._server_socket.shutdown(socket.SHUT_RDWR)
        self._server_socket.close()
        logging.info(f'action: graceful_shutdown | result: success | fd: Welcoming socket')
        if self._current_connection != None:
            self._current_connection.shutdown(socket.SHUT_RDWR)
            self._current_connection.close()
            logging.info(f'action: graceful_shutdown | result: success | fd: Client socket')
        sys.exit(0)

    def run(self):
        """
        Dummy Server loop

        Server that accept a new connections and establishes a
        communication with a client. After client with communication
        finishes, servers starts to accept new connections again
        """

        while True:
            client_sock = self.__accept_new_connection()
            self.__handle_client_connection(client_sock)

    def __handle_client_connection(self, client_sock):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:
            bet_len = bytearray(2)
            read_bytes = 0
            while read_bytes < 2:
                read_bytes += client_sock.recv_into(bet_len, 2)
            
            bet_len = int.from_bytes(bytes(bet_len[:]), "big")
            bet_data = bytearray(bet_len)
            read_bytes = 0
            while read_bytes < bet_len:
                read_bytes += client_sock.recv_into(bet_data, bet_len - read_bytes)
        
            bet = Bet.from_bytes(bet_data)

            store_bets([bet])

            client_sock.send(bytes([0x01]))

            logging.info(f'action: apuesta_almacenada | result: success | dni: {bet.document} | numero: {bet.number}')
            
        except Exception as e:
            logging.error("action: apuesta_almacenada | result: fail | error: {e}")
        finally:
            client_sock.close()
            self._current_connection = None

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        c, addr = self._server_socket.accept()
        self._current_connection = c
        logging.info(f'action: accept_connections | result: success | ip: {addr[0]}')
        return c
