package main

import (
	"fmt"
	"log"
	"net/http"

	"encoding/json"
	"io/ioutil"
	"os"
	"strings"
	"github.com/gorilla/mux"
	"sort"
)

type Type struct {
	// Name of the type
	Name string `json:"name"`
	// The effective types, damage multiplize 2x
	EffectiveAgainst []string `json:"effectiveAgainst"`
	// The weak types that against, damage multiplize 0.5x
	WeakAgainst []string `json:"weakAgainst"`
}

type Pokemon struct {
	Number         string   `json:"Number"`
	Name           string   `json:"Name"`
	Classification string   `json:"Classification"`
	TypeI          []string `json:"Type I"`
	TypeII         []string `json:"Type II,omitempty"`
	Weaknesses     []string `json:"Weaknesses"`
	FastAttackS    []string `json:"Fast Attack(s)"`
	Weight         string   `json:"Weight"`
	Height         string   `json:"Height"`
	Candy struct {
		Name     string `json:"Name"`
		FamilyID int    `json:"FamilyID"`
	} `json:"Candy"`
	NextEvolutionRequirements struct {
		Amount int    `json:"Amount"`
		Family int    `json:"Family"`
		Name   string `json:"Name"`
	} `json:"Next Evolution Requirements,omitempty"`
	NextEvolutions []struct {
		Number string `json:"Number"`
		Name   string `json:"Name"`
	} `json:"Next evolution(s),omitempty"`
	PreviousEvolutions []struct {
		Number string `json:"Number"`
		Name   string `json:"Name"`
	} `json:"Previous evolution(s),omitempty"`
	SpecialAttacks      []string `json:"Special Attack(s)"`
	BaseAttack          int      `json:"BaseAttack"`
	BaseDefense         int      `json:"BaseDefense"`
	BaseStamina         int      `json:"BaseStamina"`
	CaptureRate         float64  `json:"CaptureRate"`
	FleeRate            float64  `json:"FleeRate"`
	BuddyDistanceNeeded int      `json:"BuddyDistanceNeeded"`
}

// Move is an attack information. The
type Move struct {
	// The ID of the move
	ID int `json:"id"`
	// Name of the attack
	Name string `json:"name"`
	// Type of attack
	Type string `json:"type"`
	// The damage that enemy will take
	Damage int `json:"damage"`
	// Energy requirement of the attack
	Energy int `json:"energy"`
	// Dps is Damage Per Second
	Dps float64 `json:"dps"`
	// The duration
	Duration int `json:"duration"`
}

// BaseData is a struct for reading data.json
type BaseData struct {
	Types    []Type    `json:"types"`
	Pokemons []Pokemon `json:"pokemons"`
	Moves    []Move    `json:"moves"`
}

func splitByDelimeter(input string, delim string, output *[]string) { //split request by "/"
	runes := []rune(input)
	var sub string
	if (strings.Contains(input, "get")) {
		sub = string(runes[5:])
	} else if (strings.Contains(input, "list")) {
		sub = string(runes[6:])
	}
	*output = strings.Split(sub, delim)
}

func contains(input []string, word string) bool { //checks given word is found given string array
	for _, search := range input {
		if strings.ToLower(search) == strings.ToLower(word) {
			return true
		}
	}
	return false
}

func getQueriesByMap(input []string, delim string, output map[string]string) {
	for i := range input {
		temp := strings.Split(input[i], delim)
		output[temp[0]] = temp[1]
	}
}

func getQueriesByArray(input []string, delim string, output []string) { //splits queries by "="
	count := 0
	outputLen := len(output)
	for i := 0; i < outputLen; i = i + 2 { //increase i by 2 for storing queries in type name="name"
		temp := strings.Split(input[count], delim)
		log.Print(input[count])
		if (len(temp) == 1) { //if query is single, decrease i by 1 and range of loop by 1
			output[i] = temp[0]
			outputLen--
			i = i - 1 //if the argument number is 1 then decrease i by 1 to prevent out of bound exception
		} else if (len(temp) == 2) { //if query is double, put first argument to output[i] and second to output[i+1]
			output[i] = temp[0]
			output[i+1] = temp[1]
		}
		count++
	}
}

func (base *BaseData) sortPokes(control int, queryArr []string, w http.ResponseWriter) { //this is a method for sorting pokemons and printing them on the screen

	if control == 1 {
		var count = 0
		for i := 0; i < len(base.Pokemons); i++ {
			sort.Slice(base.Pokemons, func(k, j int) bool {
				return base.Pokemons[k].BaseAttack < base.Pokemons[j].BaseAttack //sortby base attack
			})
			if ((strings.Compare(strings.ToLower(queryArr[1]), "all") != 0) && ((contains(base.Pokemons[i].TypeI, queryArr[1])) || (contains(base.Pokemons[i].TypeII, queryArr[1])))) { //if request does not include types=all/sortby="specificproperty" but includes specific type for sorting
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "<-\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "\n")
			} else if ((strings.Compare(strings.ToLower(queryArr[1]), "all") == 0)) { //if request includes types=all/sortby="specificproperty"
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "<-\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "\n")
			}
		}
		if count == 0 {//if user want to sort pokemons with invalid type, print warning
			fmt.Fprint(w, "Invalid request.")

		}
	} else if control == 2 {
		var count = 0
		for i := 0; i < len(base.Pokemons); i++ {
			sort.Slice(base.Pokemons, func(k, j int) bool {
				return base.Pokemons[k].BaseDefense < base.Pokemons[j].BaseDefense //sortby base defence
			})
			if ((strings.Compare(strings.ToLower(queryArr[1]), "all") != 0) && ((contains(base.Pokemons[i].TypeI, queryArr[1])) || (contains(base.Pokemons[i].TypeII, queryArr[1])))) { //sort given type of pokemonss
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "<-\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "\n")
			} else if ((strings.Compare(strings.ToLower(queryArr[1]), "all") == 0)) { //sort all pokemons
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "<-\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "\n")
			}
		}
		if count == 0 {//if user want to sort pokemons with invalid type, print warning
			fmt.Fprint(w, "Invalid request.")

		}
	} else if control == 3 {
		var count = 0
		for i := 0; i < len(base.Pokemons); i++ {
			sort.Slice(base.Pokemons, func(k, j int) bool {
				return base.Pokemons[k].BaseStamina < base.Pokemons[j].BaseStamina //sortby base stamina
			})
			if ((strings.Compare(strings.ToLower(queryArr[1]), "all") != 0) && ((contains(base.Pokemons[i].TypeI, queryArr[1])) || (contains(base.Pokemons[i].TypeII, queryArr[1])))) { //sort pokemons with specific type
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "<-\n")
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "\n")
			} else if ((strings.Compare(strings.ToLower(queryArr[1]), "all") == 0)) { //sort all pokemons
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "<-\n")
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "\n")
			}
		}
		if count == 0 {//if user want to sort pokemons with invalid type, print warning
			fmt.Fprint(w, "Invalid request.")
		}
	} else if control == 4 {
		var count = 0
		for i := 0; i < len(base.Pokemons); i++ {
			sort.Slice(base.Pokemons, func(k, j int) bool {
				return base.Pokemons[k].Height < base.Pokemons[j].Height
			})
			if ((strings.Compare(strings.ToLower(queryArr[1]), "all") != 0) && ((contains(base.Pokemons[i].TypeI, queryArr[1])) || (contains(base.Pokemons[i].TypeII, queryArr[1])))) {
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "<-\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "\n")
			} else if ((strings.Compare(strings.ToLower(queryArr[1]), "all") == 0)) {
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "<-\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "\n")
			}
		}
		if count == 0 {//if user want to sort pokemons with invalid type, print warning
			fmt.Fprint(w, "Invalid request.")

		}
	} else if control == 5 {
		var count int = 0
		for i := 0; i < len(base.Pokemons); i++ {
			sort.Slice(base.Pokemons, func(k, j int) bool {
				return base.Pokemons[k].Weight < base.Pokemons[j].Weight
			})
			if ((strings.Compare(strings.ToLower(queryArr[1]), "all") != 0) && ((contains(base.Pokemons[i].TypeI, queryArr[1])) || (contains(base.Pokemons[i].TypeII, queryArr[1])))) {
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "<-\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "\n")
			} else if ((strings.Compare(strings.ToLower(queryArr[1]), "all") == 0)) {
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "<-\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "\n")
			}
		}
		if count == 0 {//if user want to sort pokemons with invalid type, print warning
			fmt.Fprint(w, "Invalid request.")

		}
	} else if control == 6 {
		var count int = 0;
		for i := 0; i < len(base.Pokemons); i++ {
			sort.Slice(base.Pokemons, func(k, j int) bool {
				return base.Pokemons[k].Name < base.Pokemons[j].Name
			})
			if ((strings.Compare(strings.ToLower(queryArr[1]), "all") != 0) && ((contains(base.Pokemons[i].TypeI, queryArr[1])) || (contains(base.Pokemons[i].TypeII, queryArr[1])))) {
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "<-\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "\n")
			} else if ((strings.Compare(strings.ToLower(queryArr[1]), "all") == 0)) {
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "<-\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "\n")
			}
		}
		if count == 0 {//if user want to sort pokemons with invalid type, print warning
			fmt.Fprint(w, "Invalid request.")

		}
	}

}

func (base *BaseData) listHandler(w http.ResponseWriter, r *http.Request) {
	//log.Println("/list url:", r.URL)
	var inputs []string
	//splitByDelimeter(string(r.URL.RawQuery), "?", &inputs)
	var input string;
	input = strings.Replace(r.URL.String(), "?", "/", -1) //if request includes ?, replace it with /
	splitByDelimeter(input, "/", &inputs)                 //split modified input by / to get each query
	/*var queryMap map[string]string
	queryMap = make(map[string]string)
	getQueriesByMap(inputs,"=", queryMap)*/

	var queryArray []string
	queryArray = make([]string, len(inputs)*2) //create slice for queries in the request
	getQueriesByArray(inputs, "=", queryArray) //split queries by = to get each word in query
	if ((strings.Compare(strings.ToLower(queryArray[1]), "") == 0) && (strings.Compare(strings.ToLower(queryArray[0]), "types") == 0)) { //if requests includes only types, print types names on screen
		fmt.Fprint(w, "Types:\n")
		i := 0
		for ; i < len(base.Types); i++ {
			fmt.Fprint(w, "     ", base.Types[i].Name, "\n")
		}
		fmt.Fprint(w, "\nTotal number of types: ", i+1)
	} else if (strings.Compare(strings.ToLower(queryArray[0]), "type") == 0) { //if request includes tpye="typename"
		if ((len(queryArray) > 2) && (strings.Compare(strings.ToLower(queryArray[2]), "sortby") == 0)) { //if request includes both type="typename" and sortby="specificproperty"
			if ((strings.Compare(strings.ToLower(queryArray[3]), "") == 0) && (strings.Compare(strings.ToLower(queryArray[2]), "sortby") == 0)) { //if request is wrong format such as sortby=NULL, print warning
				fmt.Fprint(w, "Invalid request.")
			} else if ((len(queryArray) > 3) && (strings.Compare(strings.ToLower(queryArray[3]), "baseattack") == 0)) { //if request is valid, call another function to sort and print pokemons with given type name and given sortby filter
				base.sortPokes(1, queryArray, w)
			} else if ((len(queryArray) > 3) && (strings.Compare(strings.ToLower(queryArray[3]), "basedefence") == 0)) { //if request is valid, call another function to sort and print pokemons with given type name and given sortby filter
				base.sortPokes(2, queryArray, w)
			} else if ((len(queryArray) > 3) && (strings.Compare(strings.ToLower(queryArray[3]), "basestamina") == 0)) { //if request is valid, call another function to sort and print pokemons with given type name and given sortby filter
				base.sortPokes(3, queryArray, w)
			} else if ((len(queryArray) > 3) && (strings.Compare(strings.ToLower(queryArray[3]), "height") == 0)) { //if request is valid, call another function to sort and print pokemons with given type name and given sortby filter
				base.sortPokes(4, queryArray, w)
			} else if ((len(queryArray) > 3) && (strings.Compare(strings.ToLower(queryArray[3]), "weight") == 0)) { //if request is valid, call another function to sort and print pokemons with given type name and given sortby filter
				base.sortPokes(5, queryArray, w)
			} else if ((len(queryArray) > 3) && (strings.Compare(strings.ToLower(queryArray[3]), "name") == 0)) { //if request is valid, call another function to sort and print pokemons with given type name and given sortby filter
				base.sortPokes(6, queryArray, w)
			} else { //if user want to sort pokemons with invalid property, print warning
				fmt.Fprint(w, "Invalid request.")
			}
		} else {
			if ((len(queryArray) == 2) && strings.Compare(strings.ToLower(queryArray[0]), "type") == 0) { //if user want to print pokemons without sorting them
				if ((strings.Compare(strings.ToLower(queryArray[1]), "all") == 0)) { //if user want to print all pokemons
					for i := 0; i < len(base.Pokemons); i++ {
						fmt.Fprint(w, base.Pokemons[i].Name, "\n")
						//fmt.Fprint(w, "Weight: ", base.Pokemons[i].TypeI[], " kg\n")
						fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
						fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
						fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
						fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
						fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
						fmt.Fprint(w, "     Next evolutions:\n")
						for j := range base.Pokemons[i].NextEvolutions {
							fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
						}
						fmt.Fprint(w, "\n")

					}
				} else {
					for i := 0; i < len(base.Pokemons); i++ {
						if ((contains(base.Pokemons[i].TypeI, queryArray[1])) || (contains(base.Pokemons[i].TypeII, queryArray[1]))) { //if user want to print pokemons with specific types
							fmt.Fprint(w, base.Pokemons[i].Name, "\n")
							//fmt.Fprint(w, "Weight: ", base.Pokemons[i].TypeI[], " kg\n")
							fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
							fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
							fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
							fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
							fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
							fmt.Fprint(w, "     Next evolutions:\n")
							for j := range base.Pokemons[i].NextEvolutions {
								fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
							}
							fmt.Fprint(w, "\n")
						}
					}
				}
			} else if ((len(queryArray) > 2) && (strings.Compare(strings.ToLower(queryArray[1]), "all") == 0) && strings.Compare(strings.ToLower(queryArray[0]), "type") == 0) { //if user want to print all pokemons sorted by specific property
				if ((strings.Compare(strings.ToLower(queryArray[3]), "") == 0) && (strings.Compare(strings.ToLower(queryArray[2]), "sortby") == 0)) { //if request is invalid form
					fmt.Fprint(w, "Invalid request.")
				} else if ((len(queryArray) > 3) && (strings.Compare(strings.ToLower(queryArray[3]), "baseattack") == 0)) {
					base.sortPokes(1, queryArray, w)
				} else if ((len(queryArray) > 3) && (strings.Compare(strings.ToLower(queryArray[3]), "basedefence") == 0)) {
					base.sortPokes(2, queryArray, w)
				} else if ((len(queryArray) > 3) && (strings.Compare(strings.ToLower(queryArray[3]), "basestamina") == 0)) {
					base.sortPokes(3, queryArray, w)
				} else if ((len(queryArray) > 3) && (strings.Compare(strings.ToLower(queryArray[3]), "height") == 0)) {
					base.sortPokes(4, queryArray, w)
				} else if ((len(queryArray) > 3) && (strings.Compare(strings.ToLower(queryArray[3]), "weight") == 0)) {
					base.sortPokes(5, queryArray, w)
				} else if ((len(queryArray) > 3) && (strings.Compare(strings.ToLower(queryArray[3]), "name") == 0)) {
					base.sortPokes(6, queryArray, w)
				} else {
					fmt.Fprint(w, "Invalid request.")
				}
			} else if ((strings.Compare(strings.ToLower(queryArray[0]), "types") == 0) && strings.Compare(strings.ToLower(queryArray[1]), "") != 0) { //if request is invalid form, print error message
				fmt.Fprint(w, "Invalid request.")
			} else {
				fmt.Fprint(w, "Invalid request.")
			}

		}
	} else if ((strings.Compare(strings.ToLower(queryArray[0]), "types") == 0) && strings.Compare(strings.ToLower(queryArray[1]), "") != 0) { //if request is of the form /types/sortby..., print an error message
		fmt.Fprint(w, "Invalid request.")
	} else {
		fmt.Fprint(w, "Invalid request.")
	}

	//fmt.Fprint(w, "The List Handler\n")
}

func (base *BaseData) getHandler(w http.ResponseWriter, r *http.Request) {
	//log.Println("/get url:", r.URL)

	var inputs []string
	var input string;
	input = strings.Replace(r.URL.String(), "?", "/", -1) //if request includes ?, replace it with /
	splitByDelimeter(input, "/", &inputs)                 //split request by / to get queries

	var queryArray []string
	queryArray = make([]string, len(inputs)*2) //create slice to hold queries
	getQueriesByArray(inputs, "=", queryArray) // split queries by = to get each word in a single query

	if (((strings.Compare(strings.ToLower(queryArray[0]), "name") == 0) && strings.Compare(strings.ToLower(queryArray[1]), "") != 0)) { //if query includes name="pokemonname", print details of that pokemon on the screen
		var count int = 0;
		for i := 0; i < len(base.Pokemons); i++ { //traverse base.pokemons to find a pokemon with queired name
			if ((strings.Compare(strings.ToLower(queryArray[1]), strings.ToLower(base.Pokemons[i].Name)) == 0)) {
				count++
				fmt.Fprint(w, base.Pokemons[i].Name, "\n")
				//fmt.Fprint(w, "Weight: ", base.Pokemons[i].TypeI[], " kg\n")
				fmt.Fprint(w, "     Weight: ", base.Pokemons[i].Weight, "\n")
				fmt.Fprint(w, "     Height: ", base.Pokemons[i].Height, "\n")
				fmt.Fprint(w, "     BaseAttack: ", base.Pokemons[i].BaseAttack, "\n")
				fmt.Fprint(w, "     BaseDefense: ", base.Pokemons[i].BaseDefense, "\n")
				fmt.Fprint(w, "     BaseStamina: ", base.Pokemons[i].BaseStamina, "\n")
				fmt.Fprint(w, "     BuddyDistanceNeeded: ", base.Pokemons[i].BuddyDistanceNeeded, "\n")
				fmt.Fprint(w, "     CaptureRate: ", base.Pokemons[i].CaptureRate, "\n")
				fmt.Fprint(w, "     Candy: ", base.Pokemons[i].Candy.Name, "\n")
				fmt.Fprint(w, "     Classification: ", base.Pokemons[i].Classification, "\n")
				fmt.Fprint(w, "     FleeRate: ", base.Pokemons[i].FleeRate, "\n")
				fmt.Fprint(w, "     FleeRate: ", base.Pokemons[i].FleeRate, "\n")
				fmt.Fprint(w, "     Fast Attacks: ", "\n")
				for j := range base.Pokemons[i].FastAttackS {
					fmt.Fprint(w, "          ", base.Pokemons[i].FastAttackS[j], "\n")
				}
				fmt.Fprint(w, "     Type 1: ", "\n")
				for j := range base.Pokemons[i].TypeI {
					fmt.Fprint(w, "          ", base.Pokemons[i].TypeI[j], "\n")
				}
				fmt.Fprint(w, "     Type 2: ", "\n")

				for j := range base.Pokemons[i].TypeII {
					fmt.Fprint(w, "          ", base.Pokemons[i].TypeII[j], "\n")
				}
				fmt.Fprint(w, "     Next evolutions:\n")
				for j := range base.Pokemons[i].NextEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].NextEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "     Previous Evolutions:\n")
				for j := range base.Pokemons[i].PreviousEvolutions {
					fmt.Fprint(w, "          ", base.Pokemons[i].PreviousEvolutions[j].Name, "\n")
				}
				fmt.Fprint(w, "     Special Attacks: ", "\n")
				for j := range base.Pokemons[i].SpecialAttacks {
					fmt.Fprint(w, "          ", base.Pokemons[i].SpecialAttacks[j], "\n")
				}
				fmt.Fprint(w, "     Weaknesses: ", "\n")
				for j := range base.Pokemons[i].Weaknesses {
					fmt.Fprint(w, "          ", base.Pokemons[i].Weaknesses[j], "\n")
				}
				fmt.Fprint(w, "     NextEvolutionRequirements-Name: ", base.Pokemons[i].NextEvolutionRequirements.Name, "\n")
				fmt.Fprint(w, "     NextEvolutionRequirements-Amount: ", base.Pokemons[i].NextEvolutionRequirements.Amount, "\n")
				fmt.Fprint(w, "\n")
			}

		}
		if (count == 0) { //if queried name is not in database or is wrong name print warning
			fmt.Fprint(w, "Invalid request.")
		}
	} else if (((strings.Compare(strings.ToLower(queryArray[0]), "type") == 0) && strings.Compare(strings.ToLower(queryArray[1]), "") != 0)) { //if query includes type="typename", print details of that type on the screen
		var count int = 0
		for i := 0; i < len(base.Types); i++ { //traverse base.types to find a type whose name is matchec with queried tpye name
			if ((strings.Compare(strings.ToLower(queryArray[1]), strings.ToLower(base.Types[i].Name)) == 0)) {
				count++
				fmt.Fprint(w, "Pokemon Type: ", queryArray[1], "\n")
				fmt.Fprint(w, "Effective Against: \n")
				for j := range base.Types[i].EffectiveAgainst {
					fmt.Fprint(w, "          -", base.Types[i].EffectiveAgainst[j], "\n")
				}
				fmt.Fprint(w, "Weak Against: \n")
				for j := range base.Types[i].WeakAgainst {
					fmt.Fprint(w, "          -", base.Types[i].WeakAgainst[j], "\n")
				}
				fmt.Fprint(w, "Example Pokemons: \n")
				var count int = 0
				for j := range base.Pokemons {
					if (((contains(base.Pokemons[j].TypeI, queryArray[1])) || (contains(base.Pokemons[j].TypeII, queryArray[1]))) && count < 2) {
						fmt.Fprint(w, "-", base.Pokemons[j].Name, "\n")
						count++
					}
				}
			}
		}
		if count == 0 { //if queried type name is not in database or is wrong name print warning
			fmt.Fprint(w, "Invalid request.")

		}
	} else if (((strings.Compare(strings.ToLower(queryArray[0]), "move") == 0) && strings.Compare(strings.ToLower(queryArray[1]), "") != 0)) { //if query includes move="movename", print details of that move on the screen
		var count int = 0
		for i := 0; i < len(base.Moves); i++ { //traverse base.move to find a move with queried name
			if ((strings.Compare(strings.ToLower(queryArray[1]), strings.ToLower(base.Moves[i].Name)) == 0)) {
				count++
				fmt.Fprint(w, "Move Name: ", queryArray[1], "\n")
				fmt.Fprint(w, "Move Damage: ", base.Moves[i].Damage, "\n")
				fmt.Fprint(w, "Move DPS: ", base.Moves[i].Dps, "\n")
				fmt.Fprint(w, "Move Duration: ", base.Moves[i].Duration, "\n")
				fmt.Fprint(w, "Move Energy: ", base.Moves[i].Energy, "\n")
				fmt.Fprint(w, "Move Type: ", base.Moves[i].Type, "\n")

				fmt.Fprint(w, "Example Pokemons: \n")
				var count int = 0
				for j := range base.Pokemons {
					if (((contains(base.Pokemons[j].TypeI, base.Moves[i].Type)) || (contains(base.Pokemons[j].TypeII, base.Moves[i].Type))) && count < 2) {
						fmt.Fprint(w, "-", base.Pokemons[j].Name, "\n")
						count++
					}
				}
			}
		}
		if count == 0 { //if queried move name is not in database or is wrong name print warning
			fmt.Fprint(w, "Invalid request.")
		}
	} else { //if request is wrong format print warning
		fmt.Fprint(w, "Invalid request.")
	}

	//fmt.Fprint(w, "The Get Handler\n")
	//fmt.Fprint(w, r.URL)
}

func otherwise(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Invalid request.\n")
}
func getData() BaseData {
	raw, err := ioutil.ReadFile("./data.json") //read file
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	//fmt.Printf("%s\n", string(raw))
	var base BaseData
	json.Unmarshal(raw, &base) //map information read from file to the basedate object where all json structure is hold
	//fmt.Printf("Results: %v\n", base)
	return base
}

func main() {
	//TODO: read data.json to a BaseData
	var Base = getData()
	fmt.Printf("Results: %v\n", Base)

	router := mux.NewRouter()
	router.HandleFunc("/list{sl4:[/]*}{query1:[a-zA-Z0-9=]*}{sl1:[/]*}{query2:[a-zA-Z0-9=]*}{sl2:[/]*}{query3:[a-zA-Z0-9=]*}{sl3:[/]*}", Base.listHandler) //github/gorilla/mux is used to route different requests to the same handler function
	router.HandleFunc("/list{sl4:[?]*}{query1:[a-zA-Z0-9=]*}{sl1:[?]*}{query2:[a-zA-Z0-9=]*}{sl2:[?]*}{query3:[a-zA-Z0-9=]*}{sl3:[?]*}", Base.listHandler) //With opportunity to use regex for catching requests, it is possible to catch many requests with one line regex
	router.HandleFunc("/get{sl4:[/]*}{query1:[a-zA-Z0-9=]*}{sl1:[/]*}{query2:[a-zA-Z0-9=]*}{sl2:[/]*}{query3:[a-zA-Z0-9=]*}{sl3:[/]*}", Base.getHandler)   //there are two types of requests for gethandler. one for '?' and other for '/'
	router.HandleFunc("/get{sl4:[?]*}{query1:[a-zA-Z0-9=]*}{sl1:[?]*}{query2:[a-zA-Z0-9=]*}{sl2:[?]*}{query3:[a-zA-Z0-9=]*}{sl3:[?]*}", Base.getHandler)   //there are two types of requests for listhandler. one for '?' and other for '/'
	router.HandleFunc("/", otherwise)
	http.Handle("/", router)
	//http.HandleFunc("/list", Base.listHandler)
	//http.HandleFunc("/get", getHandler)
	//TODO: add more
	//http.HandleFunc("/", otherwise)
	log.Println("starting server on :8080")
	http.ListenAndServe(":8080", nil)
}

/*
/list/types
/list?types
/list/type=bug
/list?type=bug
/list/type=bug/sortby=BaseAttack
/list?type=bug?sortby=BaseAttack
 */
/*
/get?name=pikachu
/get/name=pikachu
/get?type=bug
/get/type=bug
/get/move=asd
/get?move=asd

 */
