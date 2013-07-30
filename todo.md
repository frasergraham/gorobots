# TODO List

# general
* time matches

## Client
* Set Robot properties, have way to validate
* Messages from server:
    * X was killed etc
    * Your config is invalid etc.
    * Player dropped
    * player joined

* We should allow a "reconfigure" instruction to have the robot change itself. You can have N per game and this would be how we initially let them pick strengths and weaknesses and then they can adapt later based on their enemies.

* active scan grows each turn you maintain it?

* event on hit and death calling a user function
* on death show whole map

### Tweakable Robots Traits

* speed
* weapon speed
* weapon radius, damage will always fall off with distance from center
* weapon power
* shots at once (how many in flight at a time)
* scanner radius (scalaing multiplier for passive and active scans)
* armor (subtracts X from all damage)
* health / size (bigger = more health but bigger target)

Maybe...

* Number of open radio channels, starts at 0 perhaps (i.e. having a radio is not guranteed)
* Do we want to track ammo?


## Server
* support observer only connections
* add ids to turns
* keep track of who gets a response in, and only send state to those clients
* create the player on ws connect but dont create the robot until we get a config message with it's stats
* balancing and validation of config settings
