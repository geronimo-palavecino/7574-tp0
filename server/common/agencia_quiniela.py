import socket
from common.bet import *

""" bet size fields size"""
BET_SIZE_SIZE = 2
""" number of bets field size"""
N_BETS_SIZE = 2

""" Error representing that a problem occurred while reading from the socket """
class ReadingError(Exception):
    def __init__(self, decoded_bets=[], message="An error occurred while reading a bet from the Agencia de Quiniela"):
        self.message = message
        self.decoded_bets = decoded_bets
        super().__init__(self.message)

""" Error representing that a problem occurred while writing a confirmation into the socket """
class WritingError(Exception):
    def __init__(self, message="An error occurred while writing the confirmation to the Agencia de Quiniela"):
        self.message = message
        super().__init__(self.message)

""" A connection with an Agencia de Quiniela client """
class AgenciaQuiniela:
    def __init__(self, socket):
        self.socket = socket
    
    def address(self):
        """ Returns the address of the peer connected """
        return self.socket.getpeername()

    def get_bets(self):
        """ 
        Reads a series of bets from the underlying socket
        If the operation is successful the function returns a list with all the bets. If not, a ReadingError exception is raised.
        """
        try:
            # bet_len = int.from_bytes(bytes(bet_len[:]), "big")
            n_bets = int.from_bytes(read_data(self.socket, N_BETS_SIZE))

            bets = []

            for _ in range(n_bets):
                bet_len = int.from_bytes(read_data(self.socket, BET_SIZE_SIZE))
                bet_data = read_data(self.socket, bet_len)
                bet = Bet.from_bytes(bet_data)
                bets.append(bet)

            return bets
        except Exception as e:
            raise ReadingError(decoded_bets=bets)

    def confirm_bets(self, n):
        """ 
        Writes into the underlying connection the amount of bets read
        """
        try:
            self.socket.sendall(n.to_bytes(2, byteorder='big'))
        except Exception as e:
            raise WritingError
    
    def close(self):
        """ 
        Closes the underlying connection
        """
        self.socket.shutdown(socket.SHUT_RDWR)
        self.socket.close()

def read_data(socket, n):
    """ 
    Reads n bytes from the given socket
    """
    # The data is read avoiding short reads
    data = bytearray(n)
    read_bytes = 0
    while read_bytes < n:
        read_bytes += socket.recv_into(data, n - read_bytes)
    
    return data