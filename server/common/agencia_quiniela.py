import logging
import socket
from common.bet import *

""" bet size fields size """
BET_SIZE_SIZE = 2
""" number of bets field size """
N_BETS_SIZE = 2
""" Bet Batch packet code """
BET_BATCH_CODE = 0
""" Bet Response packet code """
BET_RESPONSE_CODE = 1
""" Winner Request packet code """
WINNER_REQUEST_CODE = 2
""" Winner Response packet code """
WINNER_RESPONSE_CODE = 3

""" Error representing that a problem occurred while reading from the socket """
class ReadingError(Exception):
    def __init__(self, decoded_bets=[], message="An error occurred while reading from the Agencia de Quiniela"):
        self.message = message
        self.decoded_bets = decoded_bets
        super().__init__(self.message)

""" Error representing that a problem occurred while writing a confirmation into the socket """
class WritingError(Exception):
    def __init__(self, message="An error occurred while writing to the Agencia de Quiniela"):
        self.message = message
        super().__init__(self.message)

""" A connection with an Agencia de Quiniela client """
class AgenciaQuiniela:
    def __init__(self, socket):
        self.socket = socket
    
    def address(self):
        """ Returns the address of the peer connected """
        return self.socket.getpeername()

    def recv_message(self):
        """ Reads the packet code of a packet """
        return int.from_bytes(read_data(self.socket, 1), "big")

    def get_bets(self):
        """ 
        Reads a series of bets from a Bets Batch packet from the underlying socket
        If the operation is successful the function returns a list with all the bets. If not, a ReadingError exception is raised.
        """
        bets = []
        
        try:
            n_bets = int.from_bytes(read_data(self.socket, N_BETS_SIZE), "big")

            for _ in range(n_bets):
                bet_len = int.from_bytes(read_data(self.socket, BET_SIZE_SIZE), "big")
                bet_data = read_data(self.socket, bet_len)
                bet = Bet.from_bytes(bet_data)
                bets.append(bet)

            return bets
        except Exception as _:
            raise ReadingError(decoded_bets=bets)
        
    def get_id(self):
        """ Reads the agency id sent in a Winner Request packet from the underlying socket """
        return int.from_bytes(read_data(self.socket, 4), "big")

    def confirm_bets(self, n):
        """ Writes a Bet Response packet with the amount of bets read into the underlying connection """
        try:
            packet = BET_RESPONSE_CODE.to_bytes(1, byteorder='big') + n.to_bytes(2, byteorder='big')
            self.socket.sendall(packet)
        except Exception as _:
            raise WritingError
    
    def send_winners(self, winners):
        """ Writes a Winner Response packet with the documents of the winning bets into the underlying connection """
        try:
            code = WINNER_RESPONSE_CODE.to_bytes(1, byteorder='big')
            length = len(winners).to_bytes(2, byteorder='big')
            documents = b''.join(winner.to_bytes(4, byteorder='big') for winner in winners)
            packet = code + length + documents
            self.socket.sendall(packet)
        except Exception as _:
            raise WritingError
    
    def close(self):
        """ Closes the underlying connection """
        self.socket.shutdown(socket.SHUT_RDWR)
        self.socket.close()

def read_data(socket, n):
    """ Reads n bytes from the given socket """
    # The data is read avoiding short reads
    data = bytearray(n)
    read_bytes = 0
    while read_bytes < n:
        read_bytes += socket.recv_into(data, n - read_bytes)
    
    return data