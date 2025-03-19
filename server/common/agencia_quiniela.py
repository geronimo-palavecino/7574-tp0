from common.bet import *

class ReadingError(Exception):
    def __init__(self, message="An error occurred while reading a bet from the Agencia de Quiniela"):
        self.message = message
        super().__init__(self.message)

class WritingError(Exception):
    def __init__(self, message="An error occurred while writing the confirmation to the Agencia de Quiniela"):
        self.message = message
        super().__init__(self.message)

class AgenciaQuiniela:
    def __init__(self, socket):
        self.socket = socket
    
    def address(self):
        return self.socket.getpeername()
    
    def get_bet(self):
        try:
            bet_len = bytearray(2)
            read_bytes = 0
            while read_bytes < 2:
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
        try:
            self.socket.sendall(code.to_bytes(1, byteorder='big'))
        except Exception as e:
            raise WritingError
    
    def close(self):
        self.socket.close()