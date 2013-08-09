
NOTES
-------

new conn.

> identify

< [robot | spectator], name, client-type, game ID

> [OK | FULL | NOT AUTH], board size, game params

< robot stats

> [ACK | NOACK]

< READY

!!!LOBBY!!!

3 2 1...

> first turn (4 players)

< first turn

> turn *

< update *




* new game
* wait for conn
    * conn 1
    * conn 2
    5
    4
    3
    2
    * conn 3
    5
    4
    3
    * conn 4
    5
    4
    3
    2
    1
    * start game, 4 players


