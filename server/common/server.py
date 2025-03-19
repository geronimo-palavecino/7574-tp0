import socket
import logging
import signal
import sys
from common.utils import *
from common.bet import *
from common.agencia_quiniela_listener import *
from common.agencia_quiniela import *

class Server:
    def __init__(self, listener):
        # Initialize server socket
        self._listener = listener
        self._current_connection = None

        # Setting the signal handler
        signal.signal(signal.SIGTERM, self.__sigterm_handler)

    def __sigterm_handler(self, signal, frame):
        """
        SIGTERM signal handler

        When the application receives a SIGTERM signal, all the file descriptors 
        (welcoming socket, and client current socket) are closed for a graceful shutdown
        """
        self._listener.close()
        logging.info(f'action: graceful_shutdown | result: success | fd: Welcoming socket')
        if self._current_connection != None:
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

    def __handle_client_connection(self, quiniela):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        try:        
            bets = quiniela.get_bets()
            store_bets(bets)
            quiniela.confirm_bets(len(bets))
            logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
        except ReadingError as e:
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(e.decoded_bets)}")
        except Exception as e:
            logging.error(f"action: unexpected_error | result: fail | error: {e}")
        finally:
            quiniela.close()
            self._current_connection = None

    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """
        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        quiniela = self._listener.accept_new_connection()
        self._current_connection = quiniela
        logging.info(f'action: accept_connections | result: success | ip: {quiniela.address()[0]}')
        return quiniela
