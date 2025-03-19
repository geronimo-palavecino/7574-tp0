import socket
from common.agencia_quiniela import *

class AgenciaQuinielaListener:
    def __init__(self, port, listen_backlog):
        #Initialize welcoming socket
        self._welcoming_socket = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
        self._welcoming_socket.bind(('', port))
        self._welcoming_socket.listen(listen_backlog)
    
    def accept_new_connection(self):
        """
        Accept new connections

        Function blocks until a connection to a client is made.
        Then connection created is printed and returned
        """

        # Connection arrived
        c, _ = self._welcoming_socket.accept()
        quiniela = AgenciaQuiniela(c)
        return quiniela
