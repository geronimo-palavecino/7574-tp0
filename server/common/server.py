import socket
import logging
import signal
import sys
from common.utils import *
from common.bet import *
from common.agencia_quiniela_listener import *
from common.agencia_quiniela import *

class Server:
    def __init__(self, listener, n_clients):
        # Initialize server socket
        self._listener = listener
        self._n_clients = n_clients
        self._current_connection = None
        self._waiting_agencys = []

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
        for _, connection in self._waiting_agencys:
            connection.close()
            logging.info(f'action: graceful_shutdown | result: success | fd: Client socket')

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
            message_type = quiniela.recv_message()
            if message_type == BET_BATCH_CODE:
                bets = quiniela.get_bets()
                store_bets(bets)
                quiniela.confirm_bets(len(bets))
                logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
                quiniela.close()
            elif message_type == WINNER_REQUEST_CODE:
                id = quiniela.get_id()
                self._waiting_agencys.append((id, quiniela))
                if len(self._waiting_agencys) == self._n_clients:
                    logging.info(f'action: sorteo | result: success')
                    winners = [[] for _ in range(self._n_clients)]
                    bets = load_bets()
                    for bet in bets:
                        if has_won(bet):
                            winners[bet.agency - 1].append(int(bet.document))
                    for id, agency in self._waiting_agencys:
                        agency.send_winners(winners[id-1])
                        agency.close()
                    self._waiting_agencys = []
            else:
                logging.error(f"action: unexpected_error | result: fail | error: Unexpected message received")
        except ReadingError as e:
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(e.decoded_bets)}")
        except Exception as e:
            logging.error(f"action: unexpected_error | result: fail | error: {e}")
        finally:
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
