import datetime

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
            first_name = bytes(data[pos:pos+len_first_name]).decode()
            pos += len_first_name
            len_last_name = int.from_bytes(bytes(data[pos:pos+1]), "big")
            pos += 1
            last_name = bytes(data[pos:pos+len_last_name]).decode()
            pos += len_last_name
            document = int.from_bytes(bytes(data[pos:pos+4]), "big")
            pos += 4
            birthdate = bytes(data[pos:pos+10]).decode()
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