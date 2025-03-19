import datetime

""" agency field byte size """
AGENCY_SIZE = 4
""" first_name length field byte size """
FIRST_NAME_SIZE = 1
""" last_name lenght field byte size """
LAST_NAME_SIZE = 1
""" document field byte size """
DOCUMENT_SIZE = 4
""" birthdate field byte size """
BIRTHDATE_SIZE = 10
""" number field byte size """
NUMBER_SIZE = 4

""" Error representing that a problem occurred while decoding a bet from a bytes encoding """
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
        """
        If possible converts a bytes representation into a Bet structure
        In case of not being posible, a DecodingError exception is raised
        """
        try: 
            pos = 0

            # Decoding the agency id field
            agency = int.from_bytes(bytes(data[pos:pos+AGENCY_SIZE]), "big")
            pos += AGENCY_SIZE

            # Decoding the first name field
            len_first_name = int.from_bytes(bytes(data[pos:pos+FIRST_NAME_SIZE]), "big")
            pos += FIRST_NAME_SIZE
            first_name = bytes(data[pos:pos+len_first_name]).decode()
            pos += len_first_name

            # Decoding the last name field
            len_last_name = int.from_bytes(bytes(data[pos:pos+LAST_NAME_SIZE]), "big")
            pos += LAST_NAME_SIZE
            last_name = bytes(data[pos:pos+len_last_name]).decode()
            pos += len_last_name

            # Decoding the document field
            document = int.from_bytes(bytes(data[pos:pos+DOCUMENT_SIZE]), "big")
            pos += DOCUMENT_SIZE

            # Decoding the birthdate field
            birthdate = bytes(data[pos:pos+BIRTHDATE_SIZE]).decode()
            pos += BIRTHDATE_SIZE

            # Decoding the number field
            number = int.from_bytes(bytes(data[pos:pos+NUMBER_SIZE]), "big")

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