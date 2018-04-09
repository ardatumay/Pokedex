First of all, to run code properly, github.com/gorilla/mux package must be downloaded in the machine. Otherwise program will give error.
Other capabilities that the program has are listed below;

-> you can send requests in these two formats
   -http://localhost:8080/list/type=Bug/sortby=BaseAttack
   -http://localhost:8080/list?type=Bug?sortby=BaseAttack
	Both requests will run properly.

List Api
-> you can list all pokemons in all types with sorted or without sorted
	example requests:
	-http://localhost:8080/list/type=all
	-http://localhost:8080/list/type=all/sortby=baseattack

-> you can list pokemons in one type with sorted or without sorted
	example requests:
	-http://localhost:8080/list/type=bug/sortby=baseattack
	-http://localhost:8080/list/type=bug

-> you can list all types
	example request
	-http://localhost:8080/list/types

-> sort optionas are:
	-Base attack
	-Base stamina
	-Base defence
	-Height
	-Weight
	-Name


Get Api
-> you can get the details of a pokemon with its name
	example requests:
	-http://localhost:8080/get/name=raichu
	-http://localhost:8080/get/name=poliwrath

-> you can get the details of a type with its name
	example requsts:
	-http://localhost:8080/get/type=bug
	-http://localhost:8080/get/type=water

-> you can get the details of a move with its name
	example requests:
	-http://localhost:8080/get/move=blizzard
	-http://localhost:8080/get/move=transform

Other notes

-- All other requests will give an invalid request error
-- All requests are non-case sensitive
	example:
	-http://localhost:8080/list/type=Bug/sortby=BaseAttack
	-http://localhost:8080/list/type=bug/sortby=baseattack
	These two are identical to each other
-- '/' and '?' can be used interchangebly in requests
	example:
	-http://localhost:8080/list?type=bug?sortby=baseattack
	http://localhost:8080/list/type=bug/sortby=baseattack