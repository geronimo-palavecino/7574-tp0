# TP0: Docker + Comunicaciones + Concurrencia

En el presente repositorio se provee un esqueleto básico de cliente/servidor, en donde todas las dependencias del mismo se encuentran encapsuladas en containers. Los alumnos deberán resolver una guía de ejercicios incrementales, teniendo en cuenta las condiciones de entrega descritas al final de este enunciado.

 El cliente (Golang) y el servidor (Python) fueron desarrollados en diferentes lenguajes simplemente para mostrar cómo dos lenguajes de programación pueden convivir en el mismo proyecto con la ayuda de containers, en este caso utilizando [Docker Compose](https://docs.docker.com/compose/).

## Instrucciones de uso
El repositorio cuenta con un **Makefile** que incluye distintos comandos en forma de targets. Los targets se ejecutan mediante la invocación de:  **make \<target\>**. Los target imprescindibles para iniciar y detener el sistema son **docker-compose-up** y **docker-compose-down**, siendo los restantes targets de utilidad para el proceso de depuración.

Los targets disponibles son:

| target  | accion  |
|---|---|
|  `docker-compose-up`  | Inicializa el ambiente de desarrollo. Construye las imágenes del cliente y el servidor, inicializa los recursos a utilizar (volúmenes, redes, etc) e inicia los propios containers. |
| `docker-compose-down`  | Ejecuta `docker-compose stop` para detener los containers asociados al compose y luego  `docker-compose down` para destruir todos los recursos asociados al proyecto que fueron inicializados. Se recomienda ejecutar este comando al finalizar cada ejecución para evitar que el disco de la máquina host se llene de versiones de desarrollo y recursos sin liberar. |
|  `docker-compose-logs` | Permite ver los logs actuales del proyecto. Acompañar con `grep` para lograr ver mensajes de una aplicación específica dentro del compose. |
| `docker-image`  | Construye las imágenes a ser utilizadas tanto en el servidor como en el cliente. Este target es utilizado por **docker-compose-up**, por lo cual se lo puede utilizar para probar nuevos cambios en las imágenes antes de arrancar el proyecto. |
| `build` | Compila la aplicación cliente para ejecución en el _host_ en lugar de en Docker. De este modo la compilación es mucho más veloz, pero requiere contar con todo el entorno de Golang y Python instalados en la máquina _host_. |

### Servidor

Se trata de un "echo server", en donde los mensajes recibidos por el cliente se responden inmediatamente y sin alterar. 

Se ejecutan en bucle las siguientes etapas:

1. Servidor acepta una nueva conexión.
2. Servidor recibe mensaje del cliente y procede a responder el mismo.
3. Servidor desconecta al cliente.
4. Servidor retorna al paso 1.


### Cliente
 se conecta reiteradas veces al servidor y envía mensajes de la siguiente forma:
 
1. Cliente se conecta al servidor.
2. Cliente genera mensaje incremental.
3. Cliente envía mensaje al servidor y espera mensaje de respuesta.
4. Servidor responde al mensaje.
5. Servidor desconecta al cliente.
6. Cliente verifica si aún debe enviar un mensaje y si es así, vuelve al paso 2.

### Ejemplo

Al ejecutar el comando `make docker-compose-up`  y luego  `make docker-compose-logs`, se observan los siguientes logs:

```
client1  | 2024-08-21 22:11:15 INFO     action: config | result: success | client_id: 1 | server_address: server:12345 | loop_amount: 5 | loop_period: 5s | log_level: DEBUG
client1  | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:14 DEBUG    action: config | result: success | port: 12345 | listen_backlog: 5 | logging_level: DEBUG
server   | 2024-08-21 22:11:14 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:15 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°1
server   | 2024-08-21 22:11:15 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:20 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:20 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°2
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°3
client1  | 2024-08-21 22:11:25 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°3
server   | 2024-08-21 22:11:25 INFO     action: accept_connections | result: in_progress
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:30 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:30 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°4
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: success | ip: 172.25.125.3
server   | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | ip: 172.25.125.3 | msg: [CLIENT 1] Message N°5
client1  | 2024-08-21 22:11:35 INFO     action: receive_message | result: success | client_id: 1 | msg: [CLIENT 1] Message N°5
server   | 2024-08-21 22:11:35 INFO     action: accept_connections | result: in_progress
client1  | 2024-08-21 22:11:40 INFO     action: loop_finished | result: success | client_id: 1
client1 exited with code 0
```


## Parte 1: Introducción a Docker
En esta primera parte del trabajo práctico se plantean una serie de ejercicios que sirven para introducir las herramientas básicas de Docker que se utilizarán a lo largo de la materia. El entendimiento de las mismas será crucial para el desarrollo de los próximos TPs.

### Ejercicio N°1:
Definir un script de bash `generar-compose.sh` que permita crear una definición de Docker Compose con una cantidad configurable de clientes.  El nombre de los containers deberá seguir el formato propuesto: client1, client2, client3, etc. 

El script deberá ubicarse en la raíz del proyecto y recibirá por parámetro el nombre del archivo de salida y la cantidad de clientes esperados:

`./generar-compose.sh docker-compose-dev.yaml 5`

Considerar que en el contenido del script pueden invocar un subscript de Go o Python:

```
#!/bin/bash
echo "Nombre del archivo de salida: $1"
echo "Cantidad de clientes: $2"
python3 mi-generador.py $1 $2
```

En el archivo de Docker Compose de salida se pueden definir volúmenes, variables de entorno y redes con libertad, pero recordar actualizar este script cuando se modifiquen tales definiciones en los sucesivos ejercicios.

#### Resolución

**<ins>Aclaración</ins>**: Aunque para levantar el container del cliente se necesita que el container del servidor este levantado, por las diferentes velocidades en las que se ejecutan los procesos hay veces en las que el cliente intenta conectarse al servidor antes de que el servidor haya podido bindear el socket de bienvenida. Luego de consultar en clase (Clase presencial del dia 25/3) se me indicó junto con un par de compañeros mas, que para evitar este problema en los test y en la demo de la entrega, se agregara un mecanismo de re-conexión al cliente. Debido a esto, para la conexión del cliente se agrego un mecanismo de 3 re-intentos de re-conexión antes de mostrar el log de fallo en la conexión

Para la resolución de este ejercicio se optó por crear un script de bash `generar-compose.sh` que recibe como argumentos los parámetros indicados en la consigna, y luego ejecuta un script de python llamado `mi-generador.py` utilizando dichos parámetros. Esto es debido a que considero que es mucho mas fácil (Y poseo mayor familiaridad) el manejo de archivos en python que con el uso de bash.

Para el script de python se definieron 4 constantes las cuales representan los diferentes artefactos dentro del docker-compose a ser creado:
- COMPOSE: Define el nombre que va a tener el compose, como así el header de los servicios que se declararán a continuación
- SERVER: Contiene toda la información para poder crear un container de un servidor dentro del proyecto
- CLIENT: Contiene toda la información para poder crear un container de un cliente dentro del proyecto
    - Dentro de esta constante se definió un *format specifier* para poder indicar el numero de cliente que será (Esto afecta a la id del mismo, tanto como al nombre del container)
- NETWORK: Contiene toda la información para crear las redes utilizadas dentro del proyecto

Para resolver la tarea plantead, el script de python, primero abre el archivo en el que se desea especificar el docker-compose. A continuación escribe los artefactos de COMPOSE y SERVER, para luego escribir CLIENTES tantas veces como haya sido especificado. Finalmente, se escribe NETWORK y se cierra el archivo

Para la ejecución de este ejercicio se debe utilizar el comando: `./generar-compose.sh <NOMBRE_DEL_ARCHIVO_DE_SALIDA> <CANTIDAD_DE_CLIENTES>`

#### Resolución

Para la resolución de este ejercicio se agrego el atributo de `volumes` a los artefactos SERVER y CLIENT. En dicho atributo se configuró que cada uno de los artefactos debe tomar su respectivo archivo de configuración como volume, provocando que la información de los respectivos archivos se persista por fuera del contenedor.

La ejecución de este ejercicio se realiza de la misma forma que el ejercicio anterior: `./generar-compose.sh <NOMBRE_DEL_ARCHIVO_DE_SALIDA> <CANTIDAD_DE_CLIENTES>`

### Ejercicio N°2:
Modificar el cliente y el servidor para lograr que realizar cambios en el archivo de configuración no requiera reconstruír las imágenes de Docker para que los mismos sean efectivos. La configuración a través del archivo correspondiente (`config.ini` y `config.yaml`, dependiendo de la aplicación) debe ser inyectada en el container y persistida por fuera de la imagen (hint: `docker volumes`).

#### Resolución

Para la resolución de este ejercicio se agrego el atributo de `volumes` a los artefactos SERVER y CLIENT. En dicho atributo se configuró que cada uno de los artefactos debe tomar su respectivo archivo de configuración como volume, provocando que la información de los respectivos archivos se persista por fuera del contenedor.

La ejecución de este ejercicio se realiza de la misma forma que el ejercicio anterior: `./generar-compose.sh <NOMBRE_DEL_ARCHIVO_DE_SALIDA> <CANTIDAD_DE_CLIENTES>`

### Ejercicio N°3:
Crear un script de bash `validar-echo-server.sh` que permita verificar el correcto funcionamiento del servidor utilizando el comando `netcat` para interactuar con el mismo. Dado que el servidor es un echo server, se debe enviar un mensaje al servidor y esperar recibir el mismo mensaje enviado.

En caso de que la validación sea exitosa imprimir: `action: test_echo_server | result: success`, de lo contrario imprimir:`action: test_echo_server | result: fail`.

El script deberá ubicarse en la raíz del proyecto. Netcat no debe ser instalado en la máquina _host_ y no se pueden exponer puertos del servidor para realizar la comunicación (hint: `docker network`). `

#### Resolución

Para la resolución de este ejercicio se creo un script de bash `validar-echo-server.sh` el cual crea un contenedor utilizando una imagen de alpine (*) y lo conecta a la red del docker-compose `tp0_testing_net`. A continuación se le indica al container que envíe un mensaje al servidor utilizando netcat. Luego se verifica si la respuesta obtenida es igual al mensaje enviado y se imprime por pantalla si la operación tuvo éxito o no. Finalmente se para y remueve el contenedor.

Para la ejecución de este ejercicio se debe utilizar el siguiente comando: `./validar-echo-server.sh`

(*) Alpine es una distribución liviana de ubuntu, que posee pre-instalado netcat. Debido a estas dos características, fue seleccionado como imagen para la creación del container

### Ejercicio N°4:
Modificar servidor y cliente para que ambos sistemas terminen de forma _graceful_ al recibir la signal SIGTERM. Terminar la aplicación de forma _graceful_ implica que todos los _file descriptors_ (entre los que se encuentran archivos, sockets, threads y procesos) deben cerrarse correctamente antes que el thread de la aplicación principal muera. Loguear mensajes en el cierre de cada recurso (hint: Verificar que hace el flag `-t` utilizado en el comando `docker compose down`).

#### Resolución

La primera parte de la resolución de este ejercicio consistió en incorporar un handler de la señal SIGTERM al servidor para cerrar el socket de bienvenida y en el caso de ser necesario, una conexión activa con un cliente en dicho momento. Para esto, se utilizo la biblioteca `signal` la cual nos permite con la función `signal.signal()` declarar una función (`__sigterm_handler()`) a ejecutar cuando se detecte este tipo de señal. Finalmente, para cerrar los _file descriptors_ anteriormente mencionados, se pasa a `__sigterm_handler()` una referencia al servidor para que cierre los mismo

Por el lado del cliente se utilizo la biblioteca os.Signal, la cual nos permite crear un canal el cual recibirá las señales dirigidas al programa. De esta forma, en el bucle de ejecución del cliente, utilizando un `select` verificamos si se recibe una señal, ante lo cual se corta la ejecución del cliente y se termina el programa. Cabe mencionar que en este caso, no se cierra ningún filedescriptor, ya que el único que se abre durante la ejecución (el socket de comunicación) se cierra al final de cada loop de envío de mensajes.

Para la ejecución del ejercicio primero se debe generar el archivo de docker-compose utilizando el commando `./generar-compose.sh <NOMBRE_DEL_ARCHIVO_DE_SALIDA> <CANTIDAD_DE_CLIENTES>`. A continuación se deben generar las imágenes de cada uno de los servicio utilizando el comando `make docker-image`. Finalmente se puede ejecutar el compose completo utilizando el commando `make docker-compose-up`

## Parte 2: Repaso de Comunicaciones

Las secciones de repaso del trabajo práctico plantean un caso de uso denominado **Lotería Nacional**. Para la resolución de las mismas deberá utilizarse como base el código fuente provisto en la primera parte, con las modificaciones agregadas en el ejercicio 4.

### Ejercicio N°5:
Modificar la lógica de negocio tanto de los clientes como del servidor para nuestro nuevo caso de uso.

#### Cliente
Emulará a una _agencia de quiniela_ que participa del proyecto. Existen 5 agencias. Deberán recibir como variables de entorno los campos que representan la apuesta de una persona: nombre, apellido, DNI, nacimiento, numero apostado (en adelante 'número'). Ej.: `NOMBRE=Santiago Lionel`, `APELLIDO=Lorca`, `DOCUMENTO=30904465`, `NACIMIENTO=1999-03-17` y `NUMERO=7574` respectivamente.

Los campos deben enviarse al servidor para dejar registro de la apuesta. Al recibir la confirmación del servidor se debe imprimir por log: `action: apuesta_enviada | result: success | dni: ${DNI} | numero: ${NUMERO}`.



#### Servidor
Emulará a la _central de Lotería Nacional_. Deberá recibir los campos de la cada apuesta desde los clientes y almacenar la información mediante la función `store_bet(...)` para control futuro de ganadores. La función `store_bet(...)` es provista por la cátedra y no podrá ser modificada por el alumno.
Al persistir se debe imprimir por log: `action: apuesta_almacenada | result: success | dni: ${DNI} | numero: ${NUMERO}`.

#### Comunicación:
Se deberá implementar un módulo de comunicación entre el cliente y el servidor donde se maneje el envío y la recepción de los paquetes, el cual se espera que contemple:
* Definición de un protocolo para el envío de los mensajes.
* Serialización de los datos.
* Correcta separación de responsabilidades entre modelo de dominio y capa de comunicación.
* Correcto empleo de sockets, incluyendo manejo de errores y evitando los fenómenos conocidos como [_short read y short write_](https://cs61.seas.harvard.edu/site/2018/FileDescriptors/).

#### Resolución

La resolución de este ejercicio consistió en 4 partes:
- Agregado de variables de entorno al generador de docker-compose
- Estructura *Bet* del lado del cliente
- Protocolo de comunicación
- Lógica del dominio

Para la primera parte, como su nombre lo indica, se agregaron las variable de entorno especificadas en la consigna al *script* de creación del docker-compose para disponer de esta información cuando se iniciara el cliente.

A continuación se procedió a crear un objeto *Bet* dentro del cliente el cual contiene toda la información de una apuesta (Agencia desde donde se realiza la apuesta, Nombre del apostante, Apellido del apostante, Documento del apostante, Fecha de nacimiento del apostante y Numero apostado), junto con una función que permite a dicha estructura representarse en formato de bytes

El formato en bytes de las *Bet*s fue definido de la siguiente forma:
1. 4 bytes para el número de la agencia
2. 1 byte que indican la longitud del nombre del apostante
3. Una tira de entre 0 y 255 bytes que contiene el nombre del apostante (Esta longitud esta determinada por el "campo" anterior)
4. 1 byte que indican la longitud del apellido del apostante
5. Una tira de entre 0 y 255 bytes que contiene el apellido del apostante (Esta longitud esta determinada por el "campo" anterior)
6. 4 bytes para el numero de documento del apostante
7. 10 bytes para la fecha de nacimiento del apostante
8. 4 bytes para el numero apostado

Para la tercera parte del ejercicio se diseño el protocolo de comunicaciones entre los clientes y el servidor. Para mantener la sencillez de la  funcionalidad de nuestro sistema, donde el cliente envía una apuesta al servidor para que sea almacenada, se diseño un protocolo igual de sencillo:
1. El cliente envía un paquete de *Bet* al servidor
2. El servidor responde con un paquete de *Bet Response* al cliente

De esta forma, se definieron los paquetes *Bet* y *Bet Response* los cuales son descritos a continuación:
- ***Bet***: Contiene toda la información necesaria para realizar una apuesta en el servidor. Para esto se definen dos campos:
    1. *Bet_size* (2 bytes): Indica la longitud en bytes de la representación en bytes de la apuesta enviada
    2. *Bet* (0-534 bytes): Contiene la representación en bytes de la apuesta a realizar
- ***Bet Response***: Contiene un código de confirmación indicando si se pudo realizar la apuesta o no. Para esto se define un campo:
    - *Response_code* (2 bytes): 0 si no se logró concretar la apuesta. De lo contrario cualquier otro valor

No obstante, para administrar el correcto manejo de los paquetes y la verificación de la correctitud de la información en los mismos, se diseñaron tres objetos, dos en el lado del servidor (`AgenciaQuiniela` y `AgenciaQuinielaListener`) y otro en el cliente (`CentralLoteriaNacional`), los cuales cumplen el rol de mantener la comunicación entre las aplicaciones.

Por el lado del cliente, al crear la `CentralLoteriaNacional` se le pasa la dirección del servidor. Luego este objeto posee la interfaz `SendBet()` la cual se encarga de realizar la conexión con el servidor, crear el paquete *Bet* con la información de la apuesta a ser enviada, enviar el paquete, esperar la confirmación del servidor y finalmente cerrar la conexión.

Por el lado del servidor, al crear un `AgenciaQuinielaListener` se le pasa el puerto en el que se va a estar escuchando por nuevas conexiones. Para realizar esta tarea, el objeto posee la función `accept_new_connection()` la cual se queda esperando a que entre una conexión y crea un objeto `AgenciaQuiniela` el cual contiene el socket a utilizar para realizar la comunicación con el nuevo cliente.

Por su parte, `AgenciaQuiniela` provee la función `get_bet()` la cual permite leer del socket la información de una apuesta y en caso de ser posible crear un objeto *Bet* (Clase ya existente en el esqueleto del trabajo) con la información enviada. Adicionalmente se provee la función `confirm_bet()` la cual se utiliza para confirmar el estado de recepción de una apuesta por parte del servidor.

Finalmentem, teniendo todas las herramientas necesarias para hacerlo, se implemento la lógica de dominio pedida en la consigna. Para esto, al crear un cliente, se le pasa un `CentralLoteriaNacional` el cual utilizara para enviar la apuesta generada a partir de los valores de las variables de entorno. Por otro lado, el servidor al ser creado, ahora se le pasa un `AgenciaQuinielaListener` el cual utilizara para recibir las conexiones de los clientes como `AgenciaQuiniela`s, y así leer las apuestas enviadas, guardarlas utilizando la función provista por la cátedra `store_bet()`, y enviar la confirmación de dicha apuesta

Para la ejecución del ejercicio primero se debe generar el archivo de docker-compose utilizando el commando `./generar-compose.sh <NOMBRE_DEL_ARCHIVO_DE_SALIDA> <CANTIDAD_DE_CLIENTES>`. A continuación se deben generar las imágenes de cada uno de los servicio utilizando el comando `make docker-image`. Finalmente se puede ejecutar el compose completo utilizando el commando `make docker-compose-up`.

### Ejercicio N°6:
Modificar los clientes para que envíen varias apuestas a la vez (modalidad conocida como procesamiento por _chunks_ o _batchs_). 
Los _batchs_ permiten que el cliente registre varias apuestas en una misma consulta, acortando tiempos de transmisión y procesamiento.

La información de cada agencia será simulada por la ingesta de su archivo numerado correspondiente, provisto por la cátedra dentro de `.data/datasets.zip`.
Los archivos deberán ser inyectados en los containers correspondientes y persistido por fuera de la imagen (hint: `docker volumes`), manteniendo la convencion de que el cliente N utilizara el archivo de apuestas `.data/agency-{N}.csv` .

En el servidor, si todas las apuestas del *batch* fueron procesadas correctamente, imprimir por log: `action: apuesta_recibida | result: success | cantidad: ${CANTIDAD_DE_APUESTAS}`. En caso de detectar un error con alguna de las apuestas, debe responder con un código de error a elección e imprimir: `action: apuesta_recibida | result: fail | cantidad: ${CANTIDAD_DE_APUESTAS}`.

La cantidad máxima de apuestas dentro de cada _batch_ debe ser configurable desde config.yaml. Respetar la clave `batch: maxAmount`, pero modificar el valor por defecto de modo tal que los paquetes no excedan los 8kB. 

Por su parte, el servidor deberá responder con éxito solamente si todas las apuestas del _batch_ fueron procesadas correctamente.

### Resolución

Para este ejercicio primero se decidió modificar el protocolo de comunicaciones entre el servidor y el cliente. Para esto se realizaron las siguientes modificaciones a los paquetes utilizados:
- ***Bet***: Debe contener la información de multiples apuestas. Para esto, al comienzo del paquete se agrego el campo *amount_bets* (2 bytes) el cual indica la cantidad de apuestas que contiene el paquete. A continuación se van a encontrar *amount_bets* concatenaciones de los campos *bet_size* y *bet* explicados en la versión anterior del protocolo
- ***Bet Response***: Ahora, en vez de transportar el código de confirmación de la versión anterior, contiene la cantidad de bets que pudieron ser procesadas

Junto con estas modificaciones se realizaron las modificaciones pertinentes a los objetos que controlan la comunicaciones entre los servicios. Para esto se modificaron `CentralLoteriaNacional` y `AgenciaQuiniela` para que puedan serializar y deserializar una lista de apuestas, respectivamente.

Finalmente se implemento la estructura `BetRepository` la cual se encarga de manejar las apuestas de los archivos montados a los contenedores de los clientes. Para esto, al crear el objeto se le pasa el path al archivo a leer, y se proveen la función `GetBets()` para la obtención de la cantidad de apuestas especificada por parámetro.

Para la ejecución del ejercicio primero se debe generar el archivo de docker-compose utilizando el commando `./generar-compose.sh <NOMBRE_DEL_ARCHIVO_DE_SALIDA> <CANTIDAD_DE_CLIENTES>`. A continuación se deben generar las imágenes de cada uno de los servicio utilizando el comando `make docker-image`. Los archivos a ser utilizados para obtener las apuestas se deben encontrar en el directorio `.data`. Finalmente se puede ejecutar el compose completo utilizando el commando `make docker-compose-up`.

### Ejercicio N°7:

Modificar los clientes para que notifiquen al servidor al finalizar con el envío de todas las apuestas y así proceder con el sorteo.
Inmediatamente después de la notificacion, los clientes consultarán la lista de ganadores del sorteo correspondientes a su agencia.
Una vez el cliente obtenga los resultados, deberá imprimir por log: `action: consulta_ganadores | result: success | cant_ganadores: ${CANT}`.

El servidor deberá esperar la notificación de las 5 agencias para considerar que se realizó el sorteo e imprimir por log: `action: sorteo | result: success`.
Luego de este evento, podrá verificar cada apuesta con las funciones `load_bets(...)` y `has_won(...)` y retornar los DNI de los ganadores de la agencia en cuestión. Antes del sorteo no se podrán responder consultas por la lista de ganadores con información parcial.

Las funciones `load_bets(...)` y `has_won(...)` son provistas por la cátedra y no podrán ser modificadas por el alumno.

No es correcto realizar un broadcast de todos los ganadores hacia todas las agencias, se espera que se informen los DNIs ganadores que correspondan a cada una de ellas.

#### Resolución

Para comenzar, se modificó el *script* que genera el docker-compose agregando al servidor una variable de entorno que le indica la cantidad de clientes que se conectarán. Esto se utilizado mas adelante para la espera de que todos los clientes hayan pedido los ganadores para enviar los mismos


Debido a que ahora los clientes deben poder solicitar a los ganadores del sorteo, el protocolo de comunicaciones se vio afectado de la siguiente forma:
1. El cliente envía un paquete *Bet* con un listado de apuestas
2. El servidor response con la cantidad de apuestas que se pudieron procesar
3. Esto se repite hasta que el cliente no tiene mas apuestas para enviar
4. El cliente envía un paquete *Winner Request* al servidor pidiéndole los ganadores correspondientes a esa agencia
5. Cuando el servidor recibe los paquetes *Winner Request* de todos los clientes, responde a cada uno con un paquete *Winner Response* el cual contiene los ganadores de la agencia a los que están siendo enviados

Debido a que ahora tanto el cliente como el servidor pueden enviar dos tipos de mensajes diferentes, se debido modificar los paquetes anteriores para que estos contengan un campo *code* (1 byte) al comienzo el cual señaliza que tipo de paquete se trata. De esta forma se definieron los siguientes códigos:
- *Bet Batch*: 0
- *Bet Response*: 1
- *Winner Request*: 2
- *Winner Response*: 3

Por ultimo respecto al protocolo de comunicaciones, se definió la estructura de los dos nuevos paquetes creados siendo las siguientes:
- ***Winner Request***: Contiene los siguientes campos:
    - *code* (1 byte): Indica que el paquete se trata de una *Winner Request*
    - *agency* (4 byets): Indica que agencia esta realizando el pedido
- ***Winner Response***: Debe trasportar los documentos de los ganadores de la agencia que los pidió. Para esto posee los siguientes campos:
    - *code* (1 byte): Indica que el paquete se trata de una *Winner Response*
    - *amount* (2 bytes): Indica la cantidad de documentos que contiene el paquete
    - *winner* (4 bytes): Documento de un ganador. Este campo se repite *amount* veces

Para realizar el manejo de estos nuevos requerimientos del protocolo en el domino se modificaron las estructuras pre-existentes:
- `CentralLoteriaNacional`:
    - `GetWinners()`: Esta función agregada se encarga de iniciar la comunicación con el servidor, crear y enviar el paquete de *Winners Request* y esperar la respuesta del servidor. En caso de recibir un paquete *Winners Response* válido, retorna los documentos de los ganadores
- `AgenciaQuiniela`:
    - `recv_message()`: Esta función permite leer el campo *code* de un paquete para identificar de que tipo se trata. Esta función se utiliza en conjunto con el resto de funciones para poder identificar que tipo de paquete es y saber a que función de procesamiento llamar
    - `get_id()`: Esta función permite leer el campo *agency* de un paquete *Winner Request*
    - `send_winners()`: Esta función serializa la lista de documentos de ganadores, crea un paquete *Winner Response* y se lo envía a la agencia

Finalmente, para poder afrontar los nuevos requerimientos de la lógica de dominio, se agrego del lado del cliente que una vez que no se pueden obtener mas apuestas de `BetRepository` se llame a la función `GetWinners()` de `CentralLoteriaNacional`. Por el lado del servidor, se incorporó una atributo interno al servidor el cual posee una lista de las clientes que están esperando a los ganadores. De esta forma, cada vez que una agencia pide a los ganadores se lo agrega a dicha lista y se verifica si la cantidad de clientes esperando es igual al total de los clientes. En dicho caso, se cargan todas las apuestas y se verifican los ganadores, agregándolos a listas correspondientes con cada una de las agencias. Una vez que este proceso se termina, se envía cada lista a su agencia correspondiente.

Para la ejecución del ejercicio primero se debe generar el archivo de docker-compose utilizando el commando `./generar-compose.sh <NOMBRE_DEL_ARCHIVO_DE_SALIDA> <CANTIDAD_DE_CLIENTES>`. A continuación se deben generar las imágenes de cada uno de los servicio utilizando el comando `make docker-image`. Los archivos a ser utilizados para obtener las apuestas se deben encontrar en el directorio `.data`. Finalmente se puede ejecutar el compose completo utilizando el commando `make docker-compose-up`.

## Parte 3: Repaso de Concurrencia
En este ejercicio es importante considerar los mecanismos de sincronización a utilizar para el correcto funcionamiento de la persistencia.

### Ejercicio N°8:

Modificar el servidor para que permita aceptar conexiones y procesar mensajes en paralelo. En caso de que el alumno implemente el servidor en Python utilizando _multithreading_,  deberán tenerse en cuenta las [limitaciones propias del lenguaje](https://wiki.python.org/moin/GlobalInterpreterLock).

#### Resolución

<ins>**Aclaración**</ins>: Luego de leer las limitaciones del lenguaje al utilizar multithreading, y tener en consideración otras opciones como multiprocessing, se tomó la decisión de utilizar multithreading de todas formas por las siguientes razones:
- A causa del GIL (Global Interpreter Lock) de python, solo una única tarea se puede ejecutar en el procesador en un momento dado, por lo que no se aprovecha correctamente las capacidades computacionales del procesador. Sin embargo la librearía provee un modelo que no tiene problemas en ejecutar operaciones I/O en simultaneo. Debido a que las tareas que van a estar realizando cada uno de los threads en nuestro sistema es fuerte en I/O (Leer/Escribir un socket, leer/escribir un archivo), no sería totalmente des-provechoso utilizar este modelo
- Si se utilizase otro modelo de concurrencia, como utilizar *multiprocessing* se estaría invocando un proceso por cada una de las conexiones con los clientes. Debido a las tareas que se van a estar realizando, sería muy alto el uso de los recursos en comparación a lo que se desea realizar.

Para la parelelización de las tareas del servidor, se utilizó la biblioteca *threading* de python. Con esto, se crea un thread por cada una de las conexiones de los clientes. Para la sincronización de los threads, se agregó un lock el cual permite la serialización del acceso del archivo donde se escriben las apuestas, como así a la lista de conexiones activas y clientes que esperan a los ganadores

Debido a la sicronización generada por el uso del lock, solamente el último thread en encolar un cliente para esperar a los ganadores tendrá que revisar todas las apuestas y enviar los ganadores. Esta implementación, en comparación con otros métodos de sincronización como puede haber sido el uso de barriers entre los diferentes threads, trae las siguientes dos ventajas:
- Debido a que los threads no se deben quedar esperando a que todos los clientes deban pedir a los ganadores, los recursos de los mismos se pueden ir liberando lo mas rápido posible.
- Las apuestas se deben leer una única vez para construir todas las listas de ganadores a ser enviadas. De esta forma se ahorra leer n_clientes-1 veces el archivo de apuestas y los threads no se deben "pelear" por el tiempo en el cpu para realizar las operaciones y determinar a los ganadores.

Para la ejecución del ejercicio primero se debe generar el archivo de docker-compose utilizando el commando `./generar-compose.sh <NOMBRE_DEL_ARCHIVO_DE_SALIDA> <CANTIDAD_DE_CLIENTES>`. A continuación se deben generar las imágenes de cada uno de los servicio utilizando el comando `make docker-image`. Los archivos a ser utilizados para obtener las apuestas se deben encontrar en el directorio `.data`. Finalmente se puede ejecutar el compose completo utilizando el commando `make docker-compose-up`.

## Condiciones de Entrega
Se espera que los alumnos realicen un _fork_ del presente repositorio para el desarrollo de los ejercicios y que aprovechen el esqueleto provisto tanto (o tan poco) como consideren necesario.

Cada ejercicio deberá resolverse en una rama independiente con nombres siguiendo el formato `ej${Nro de ejercicio}`. Se permite agregar commits en cualquier órden, así como crear una rama a partir de otra, pero al momento de la entrega deberán existir 8 ramas llamadas: ej1, ej2, ..., ej7, ej8.
 (hint: verificar listado de ramas y últimos commits con `git ls-remote`)

Se espera que se redacte una sección del README en donde se indique cómo ejecutar cada ejercicio y se detallen los aspectos más importantes de la solución provista, como ser el protocolo de comunicación implementado (Parte 2) y los mecanismos de sincronización utilizados (Parte 3).

Se proveen [pruebas automáticas](https://github.com/7574-sistemas-distribuidos/tp0-tests) de caja negra. Se exige que la resolución de los ejercicios pase tales pruebas, o en su defecto que las discrepancias sean justificadas y discutidas con los docentes antes del día de la entrega. El incumplimiento de las pruebas es condición de desaprobación, pero su cumplimiento no es suficiente para la aprobación. Respetar las entradas de log planteadas en los ejercicios, pues son las que se chequean en cada uno de los tests.

La corrección personal tendrá en cuenta la calidad del código entregado y casos de error posibles, se manifiesten o no durante la ejecución del trabajo práctico. Se pide a los alumnos leer atentamente y **tener en cuenta** los criterios de corrección informados  [en el campus](https://campusgrado.fi.uba.ar/mod/page/view.php?id=73393).
