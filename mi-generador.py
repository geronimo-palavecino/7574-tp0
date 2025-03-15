COMPOSE = '''
name: tp0
services:
'''

SERVER = '''
  server:                                   
    container_name: server                  
    image: server:latest                    
    entrypoint: python3 /main.py            
    environment:                            
      - PYTHONUNBUFFERED=1                  
      - LOGGING_LEVEL=DEBUG                 
    networks:                               
      - testing_net                         
'''

CLIENT = '''
  client{client_n}:
    container_name: client{client_n}
    image: client:latest
    entrypoint: /client
    environment:
      - CLI_ID={client_n}
      - CLI_LOG_LEVEL=DEBUG
    networks:
      - testing_net
    depends_on:
      - server    
'''

NETWORK = '''
networks:
  testing_net:
    ipam:
      driver: default                       
      config:
        - subnet: 172.25.125.0/24
'''

import sys

def main(file_name, n_clients):
    f = open(file_name, 'w')
    f.write(COMPOSE)
    f.write(SERVER)
    for n in range(n_clients):
        f.write(CLIENT.format(client_n=n+1))
    f.write(NETWORK)
    f.close()

if __name__ == "__main__":
    if len(sys.argv) > 2:
        main(sys.argv[1], int(sys.argv[2]))
    else:
        print("Not enough arguments provided")