from common.bet import *

""" packet size fields size"""
PACKET_SIZE_SIZE = 2

""" Error representing that a problem occurred while reading from the socket """
class ReadingError(Exception):
    def __init__(self, message="An error occurred while reading a bet from the Agencia de Quiniela"):
        self.message = message
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
    
    def get_bet(self):
        """ 
        Reads a bet bytes representation from the underlying socket.
        Then the bytes are decoded into a Bet which is later returned
        """
        try:
            bet_len = bytearray(PACKET_SIZE_SIZE)
            read_bytes = 0
            while read_bytes < PACKET_SIZE_SIZE:
                read_bytes += self.socket.recv_into(bet_len, 2)
            
            bet_len = int.from_bytes(bytes(bet_len[:]), "big")
            bet_data = bytearray(bet_len)
            read_bytes = 0
            while read_bytes < bet_len:
                read_bytes += self.socket.recv_into(bet_data, bet_len - read_bytes)
        
            bet = Bet.from_bytes(bet_data)

            return bet
        except Exception as e:
            raise ReadingError

    def confirm_bet(self, code):
        """ 
        Writes into the underlying connection the amount of bets read
        """
        try:
            self.socket.sendall(code.to_bytes(1, byteorder='big'))
        except Exception as e:
            raise WritingError
    
    def close(self):
        """ 
        Closes the underlying connection
        """
        self.socket.close()