import threading
import logging
import signal
import sys
from common.utils import *
from common.bet import *
from common.agencia_quiniela_listener import *
from common.agencia_quiniela import *

class Connections:
    def __init__(self):
        self._id_counter = 0
        self._current_connections = dict()
        self._waiting_agencys = []
        self._lock = threading.Lock()

class Server:
    def __init__(self, listener, n_clients):
        # Initialize server socket
        self._listener = listener
        self._n_clients = n_clients
        self._connections = Connections()

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
        for _, connection in self._connections._current_connections:
            connection.close()
            logging.info(f'action: graceful_shutdown | result: success | fd: Client socket')
        for _, connection in self._connections._waiting_agencys:
            connection.close()
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
            connection_id = self.__accept_new_connection()
            worker = threading.Thread(target=self.__handle_client_connection, args=(connection_id, ))
            worker.start()
    
    def __handle_client_connection(self, connection_id):
        """
        Read message from a specific client socket and closes the socket

        If a problem arises in the communication with the client, the
        client socket will also be closed
        """
        quiniela = self._connections._current_connections[connection_id]
        try:
            while True: 
                message_type = quiniela.recv_message()
                if message_type == BET_BATCH_CODE:
                    _bet_batch_handler(quiniela, connection_id, self._connections)
                elif message_type == WINNER_REQUEST_CODE:
                    _winner_request_handler(quiniela, connection_id, self._connections, self._n_clients)
                    break
                else:
                    logging.error(f"action: unexpected_error | result: fail | error: Unexpected message received")
        except ReadingError as e:
            logging.error(f"action: apuesta_recibida | result: fail | cantidad: {len(e.decoded_bets)}")
        except Exception as e:
            logging.error(f"action: unexpected_error | result: fail | error: {e}")
    
    def __accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """
        # Connection arrived
        logging.info('action: accept_connections | result: in_progress')
        quiniela = self._listener.accept_new_connection()
        with self._connections._lock:
            id = self._connections._id_counter
            self._connections._id_counter += 1
            self._connections._current_connections[id] = quiniela
            logging.info(f'action: accept_connections | result: success | ip: {quiniela.address()[0]}')
            return id


def _bet_batch_handler(quiniela, connection_id, connections):
    logging.info(f'{connection_id} Acquire')
    bets = quiniela.get_bets()
    connections._lock.acquire()
    store_bets(bets)
    quiniela.confirm_bets(len(bets))
    logging.info(f'action: apuesta_recibida | result: success | cantidad: {len(bets)}')
    logging.info(f'{connection_id} Release')
    connections._lock.release()

def _winner_request_handler(quiniela, connection_id, connections, n_clients):
    id = quiniela.get_id()
    connections._lock.acquire()
    logging.info(f'{connection_id} Acquire')
    connections._waiting_agencys.append((id, quiniela))
    connections._current_connections.pop(connection_id)
    if len(connections._waiting_agencys) == n_clients:
        logging.info(f'action: sorteo | result: success')
        winners = [[] for _ in range(n_clients)]
        bets = load_bets()
        for bet in bets:
            if has_won(bet):
                winners[bet.agency - 1].append(int(bet.document))
        for id, agency in connections._waiting_agencys:
            agency.send_winners(winners[id-1])
            agency.close()
        connections._waiting_agencys = []
    logging.info(f'{connection_id} Release')
    connections._lock.release()