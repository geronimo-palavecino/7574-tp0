import csv
import datetime
import time


""" Bets storage location. """
STORAGE_FILEPATH = "./bets.csv"
""" Simulated winner number in the lottery contest. """
LOTTERY_WINNER_NUMBER = 7574

class DecodingError(Exception):
    def __init__(self, message="An error occurred while decoding the bet"):
        self.message = message
        super().__init__(self.message)


""" A lottery bet registry. """
class Bet:
    def __init__(self, agency: str, first_name: str, last_name: str, document: str, birthdate: str, number: str):
        """
        agency must be passed with integer format.
        birthdate must be passed with format: 'YYYY-MM-DD'.
        number must be passed with integer format.
        """
        self.agency = int(agency)
        self.first_name = first_name
        self.last_name = last_name
        self.document = document
        self.birthdate = datetime.date.fromisoformat(birthdate)
        self.number = int(number)
    
    @classmethod
    def from_bytes(cls, data):
        try: 
            pos = 0
            agency = int.from_bytes(bytes(data[pos:pos+4]), "big")
            pos += 4
            len_first_name = int.from_bytes(bytes(data[pos:pos+1]), "big")
            pos += 1
            first_name = bytes(data[pos:pos+len_first_name])
            pos += len_first_name
            len_last_name = int.from_bytes(bytes(data[pos:pos+1]), "big")
            pos += 1
            last_name = bytes(data[pos:pos+len_last_name]).decode()
            pos += len_last_name
            document = int.from_bytes(bytes(data[pos:pos+4]), "big")
            pos += 1
            birthdate = datetime.date.fromisoformat(bytes(data[pos:pos+10]).decode())
            pos += 10
            number = int.from_bytes(bytes(data[pos:pos+4]), "big")

            return cls(
                agency,
                first_name,
                last_name,
                document,
                birthdate,
                number,
            )
        except Exception as e: 
            raise DecodingError

""" Checks whether a bet won the prize or not. """
def has_won(bet: Bet) -> bool:
    return bet.number == LOTTERY_WINNER_NUMBER

"""
Persist the information of each bet in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def store_bets(bets: list[Bet]) -> None:
    with open(STORAGE_FILEPATH, 'a+') as file:
        writer = csv.writer(file, quoting=csv.QUOTE_MINIMAL)
        for bet in bets:
            writer.writerow([bet.agency, bet.first_name, bet.last_name,
                             bet.document, bet.birthdate, bet.number])

"""
Loads the information all the bets in the STORAGE_FILEPATH file.
Not thread-safe/process-safe.
"""
def load_bets() -> list[Bet]:
    with open(STORAGE_FILEPATH, 'r') as file:
        reader = csv.reader(file, quoting=csv.QUOTE_MINIMAL)
        for row in reader:
            yield Bet(row[0], row[1], row[2], row[3], row[4], row[5])

