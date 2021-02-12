package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

var columnID = 1
var floorRequestButtonID = 1
var elevatorID = 1
var callButtonID = 1
var callButtonFloor = 1

// Battery Struct
type Battery struct {
	ID                      int
	amountOfColumns         int
	amountOfFloors          int
	amountOfBasements       int
	status                  string
	columnsList             []Column
	floorRequestButtonsList []FloorRequestButton
}

//Column Struct
type Column struct {
	ID                int
	name              string
	status            string
	amountOfElevators int
	servedFloors      []int
	elevatorsList     []Elevator
	callButtonList    []CallButton
}

// Elevator Struct
type Elevator struct {
	ID               int
	name             string
	status           string
	currentFloor     int
	screenDisplay    int
	direction        string
	door             Door
	floorRequestList []int
}

// CallButton Struct
type CallButton struct {
	ID        int
	status    string
	floor     int
	direction string
}

// FloorRequestButton Struct
type FloorRequestButton struct {
	ID        int
	status    string
	floor     int
	direction string
}

// Door Struct
type Door struct {
	ID     int
	status string
}

// BestElevatorInfo Struct
type BestElevatorInfo struct {
	bestElevator Elevator
	bestScore    int
	referanceGap int
}

//-----------------------------------------------------"  Battery  "-----------------------------------------------------

func createBattery(_id int, _amountOfColumns int, _status string, _amountOfFloors int, _amountOfBasements int, _amountOfElevatorPerColumn int) *Battery {

	battery := Battery{_id, _amountOfColumns, _amountOfFloors, _amountOfBasements, _status, []Column{}, []FloorRequestButton{}}

	//Checks whether basement exists
	if _amountOfBasements > 0 {
		//Create and add the Basement's columns
		battery.createBasementColumn(_amountOfBasements, _amountOfElevatorPerColumn)
		//Create and add the Basement's floorRequestButton to the List.
		battery.createBasementFloorRequestButtons(_amountOfBasements)
		_amountOfColumns--
	}
	battery.createColumns(_amountOfColumns, _amountOfFloors, _amountOfElevatorPerColumn)
	battery.createFloorRequestButtons(_amountOfFloors)

	return &battery
}

func (battery *Battery) createBasementColumn(_amountOfBasements int, _amountOfElevatorPerColumn int) {

	servedFloorsList := []int{}
	servedFloorsList = append(servedFloorsList, 1)
	floor := -1

	for i := 1; i <= _amountOfBasements; i++ {
		servedFloorsList = append(servedFloorsList, floor)
		floor--
	}

	column := Column{columnID, "A", "online", _amountOfElevatorPerColumn, servedFloorsList, []Elevator{}, []CallButton{}}
	column.createElevators(_amountOfBasements, _amountOfElevatorPerColumn)
	column.createCallButtons(_amountOfBasements, true)
	battery.columnsList = append(battery.columnsList, column)
	columnID++
}

func (battery *Battery) createColumns(_amountOfColumns int, _amountOfFloors int, _amountOfElevatorPerColumn int) {

	columnNameList := []string{"B", "C", "D"}

	floorsPerColumn := float64(_amountOfFloors / _amountOfColumns)
	amountOfFloorsPerColumn := int(math.Ceil(floorsPerColumn))
	floor := 1

	for i := 1; i <= _amountOfColumns; i++ {

		servedFloorsList := []int{}

		for i := 1; i <= amountOfFloorsPerColumn; i++ {

			if floor <= _amountOfFloors {
				servedFloorsList = append(servedFloorsList, floor)
				floor++
			}
		}
		searchedFloor := Find(servedFloorsList, 1)
		if searchedFloor == false {
			servedFloorsList = append(servedFloorsList, 1)
		}
		sort.Ints(servedFloorsList)
		column := Column{columnID, columnNameList[i-1], "online", _amountOfElevatorPerColumn, servedFloorsList, []Elevator{}, []CallButton{}}
		column.createElevators(amountOfFloorsPerColumn, _amountOfElevatorPerColumn)
		column.createCallButtons(amountOfFloorsPerColumn, false)
		battery.columnsList = append(battery.columnsList, column)
		columnID++
	}
}

func (battery *Battery) createBasementFloorRequestButtons(_amountOfBasements int) {

	buttonFloor := -1
	for i := 1; i <= _amountOfBasements; i++ {

		floorRequestButton := FloorRequestButton{floorRequestButtonID, "off", buttonFloor, "down"}
		battery.floorRequestButtonsList = append(battery.floorRequestButtonsList, floorRequestButton)
		buttonFloor--
		floorRequestButtonID++
	}
}

func (battery *Battery) createFloorRequestButtons(_amountOfFloors int) {

	for buttonFloor := 1; buttonFloor <= _amountOfFloors; buttonFloor++ {
		floorRequestButton := FloorRequestButton{floorRequestButtonID, "off", buttonFloor, "up"}
		battery.floorRequestButtonsList = append(battery.floorRequestButtonsList, floorRequestButton)
		floorRequestButtonID++
	}
}

func (battery *Battery) findBestColumn(_requestedFloor int) *Column {

	foundColumn := Column{}
	for _, column := range battery.columnsList {

		searchedColumn := Find(column.servedFloors, _requestedFloor)
		if searchedColumn {
			foundColumn = column
		}
	}
	return &foundColumn
}

func (battery *Battery) assignElevator(_requestedFloor int, _direction string) {
	column := battery.findBestColumn(_requestedFloor)
	fmt.Println("- Selected Column: ", column.name)
	elevator := column.findElevator(1, _direction)
	fmt.Println("- Selected Elevator: ", elevator.name)
	elevator.floorRequestList = append(elevator.floorRequestList, _requestedFloor)
	elevator.move()
}

//-----------------------------------------------------"  Column  "-----------------------------------------------------

func (column *Column) createElevators(_amountOfFloors int, _amountOfElevators int) {

	for i := 1; i <= _amountOfElevators; i++ {
		elevator := Elevator{elevatorID, column.name + strconv.Itoa(i), "idle", 1, 1, "null", Door{}, []int{}}
		column.elevatorsList = append(column.elevatorsList, elevator)
		elevatorID++
	}
}

func (column *Column) createCallButtons(_amountOfFloors int, _isBasement bool) {

	if _isBasement {

		buttonBasement := -1
		for i := 1; i <= _amountOfFloors; i++ {

			callButton := CallButton{callButtonID, "off", buttonBasement, "up"}
			column.callButtonList = append(column.callButtonList, callButton)
			buttonBasement--
			callButtonID++
		}
	} else {

		for i := 1; i <= _amountOfFloors; i++ {

			callButton := CallButton{callButtonID, "off", callButtonFloor, "down"}
			column.callButtonList = append(column.callButtonList, callButton)
			callButtonFloor++
			callButtonID++
		}
	}
}

func (column *Column) requestElevator(_requestedFloor int, _direction string) *Elevator {

	fmt.Println("- Current column: ", column.name)
	elevator := column.findElevator(_requestedFloor, _direction)
	fmt.Println("- Selected Elevator: ", elevator.name)
	elevator.floorRequestList = append(elevator.floorRequestList, _requestedFloor)
	elevator.sortFloorList()
	elevator.move()

	return &elevator
}

func (column *Column) findElevator(_requestedFloor int, _direction string) Elevator {

	bestElevatorInfo := BestElevatorInfo{Elevator{}, 6, 10000000}

	if _requestedFloor == 1 {
		for _, elevator := range column.elevatorsList {

			if _requestedFloor == elevator.currentFloor && elevator.status == "stopped" {
				//The elevator is at the lobby and already has some requests. It is about to leave but has not yet departe
				bestElevatorInfo = column.checkIfElevatorISBetter(1, elevator, bestElevatorInfo, _requestedFloor)

			} else if _requestedFloor == elevator.currentFloor && elevator.status == "idle" {
				//The elevator is at the lobby and has no requests

				bestElevatorInfo = column.checkIfElevatorISBetter(2, elevator, bestElevatorInfo, _requestedFloor)

			} else if _requestedFloor > elevator.currentFloor && elevator.direction == "up" {
				//The elevator is lower than me and is coming up. It means that I am requesting an elevator to go to a basement, and the elevator is on it's way to me.

				bestElevatorInfo = column.checkIfElevatorISBetter(3, elevator, bestElevatorInfo, _requestedFloor)

			} else if _requestedFloor < elevator.currentFloor && elevator.direction == "down" {
				//The elevator is above me and is coming down. It means that I'm requesting an elevator to go to a floor, and the elevator is on it's way to me

				bestElevatorInfo = column.checkIfElevatorISBetter(3, elevator, bestElevatorInfo, _requestedFloor)

			} else if elevator.status == "idle" {
				//The elevator is not at the first floor, but doesn't have any request

				bestElevatorInfo = column.checkIfElevatorISBetter(4, elevator, bestElevatorInfo, _requestedFloor)
			} else {
				//The elevator is not available, but still could take the call if nothing better is foun
				bestElevatorInfo = column.checkIfElevatorISBetter(5, elevator, bestElevatorInfo, _requestedFloor)
			}
		}
	} else {
		for _, elevator := range column.elevatorsList {

			if _requestedFloor == elevator.currentFloor && elevator.status == "stopped" && _direction == elevator.direction {
				//The elevator is at the same level as me, and is about to depart to the first floor

				bestElevatorInfo = column.checkIfElevatorISBetter(1, elevator, bestElevatorInfo, _requestedFloor)
			} else if _requestedFloor > elevator.currentFloor && elevator.direction == "up" && _direction == "up" {
				//The elevator is lower than me and is going up. I'm on a basement, and the elevator can pick me up on it's way

				bestElevatorInfo = column.checkIfElevatorISBetter(2, elevator, bestElevatorInfo, _requestedFloor)
			} else if _requestedFloor < elevator.currentFloor && elevator.direction == "down" && _direction == "down" {
				//The elevator is higher than me and is going down. I'm on a floor, and the elevator can pick me up on it's way

				bestElevatorInfo = column.checkIfElevatorISBetter(2, elevator, bestElevatorInfo, _requestedFloor)
			} else if elevator.status == "idle" {
				//The elevator is idle and has no requests

				bestElevatorInfo = column.checkIfElevatorISBetter(4, elevator, bestElevatorInfo, _requestedFloor)
			} else {
				//The elevator is not available, but still could take the call if nothing better is found

				bestElevatorInfo = column.checkIfElevatorISBetter(5, elevator, bestElevatorInfo, _requestedFloor)
			}
		}
	}
	return bestElevatorInfo.bestElevator
}

func (column *Column) checkIfElevatorISBetter(scoreToCheck int, newElevator Elevator, bestElevatorInfo BestElevatorInfo, _requestedFloor int) BestElevatorInfo {
	if scoreToCheck < bestElevatorInfo.bestScore {

		bestElevatorInfo.bestScore = scoreToCheck
		bestElevatorInfo.bestElevator = newElevator
		bestElevatorInfo.referanceGap = int(math.Abs(float64(newElevator.currentFloor - _requestedFloor)))

	} else if bestElevatorInfo.bestScore == scoreToCheck {

		gap := int(math.Abs(float64(newElevator.currentFloor - _requestedFloor)))

		if bestElevatorInfo.referanceGap > gap {
			bestElevatorInfo.bestScore = scoreToCheck
			bestElevatorInfo.bestElevator = newElevator
			bestElevatorInfo.referanceGap = gap
		}
	}
	return bestElevatorInfo
}

//-----------------------------------------------------"  Elevator  "-----------------------------------------------------

func (elevator *Elevator) move() {
	for len(elevator.floorRequestList) != 0 {

		destination := elevator.floorRequestList[0]
		elevator.operateDoors("closed")

		if elevator.door.status == "closed" {

			fmt.Println("Status door:", elevator.door.status)
			elevator.status = "moving"
			elevator.screenDisplay = elevator.currentFloor
			fmt.Println("Elevator Status: ", elevator.status, " ||  Elevator Display: ", elevator.screenDisplay)

			if elevator.currentFloor < destination {

				elevator.direction = "up"

				for elevator.currentFloor < destination {

					elevator.currentFloor++

					if elevator.currentFloor != 0 {

						elevator.screenDisplay = elevator.currentFloor
						fmt.Println("Elevator Status: ", elevator.status, " ||  Elevator Display: ", elevator.screenDisplay)
					}
				}
			} else if elevator.currentFloor > destination {

				elevator.direction = "down"

				for elevator.currentFloor > destination {

					elevator.currentFloor--
					elevator.screenDisplay = elevator.currentFloor
					fmt.Println("Elevator Status: ", elevator.status, " ||  Elevator Display: ", elevator.screenDisplay)
				}
			}
			elevator.status = "stopped"
			fmt.Println("Elevator Status: ", elevator.status)
			elevator.operateDoors("openned")
			fmt.Println("Status door:", elevator.door.status)
		}
		elevator.floorRequestList = RemoveIndex(elevator.floorRequestList, 0)
	}
	elevator.status = "idle"
}

func (elevator *Elevator) sortFloorList() {
	if elevator.direction == "up" {

		sort.Slice(elevator.floorRequestList, func(i, j int) bool {
			return elevator.floorRequestList[i] < elevator.floorRequestList[j]
		})
	} else {

		sort.Slice(elevator.floorRequestList, func(i, j int) bool {
			return elevator.floorRequestList[i] > elevator.floorRequestList[j]
		})
	}
}

func (elevator *Elevator) operateDoors(_command string) {

	sensorDoor := false
	if sensorDoor == false {
		elevator.door.status = _command
	} else {
		fmt.Println("Blocked door")
	}
}

//-----------------------------------------------------"  Auxiliary functions  "-----------------------------------------------------

// Find func
func Find(slice []int, val int) bool {
	for _, a := range slice {
		if a == val {
			return true
		}
	}
	return false
}

//RemoveIndex func
func RemoveIndex(s []int, index int) []int {
	return append(s[:index], s[index+1:]...)
}

//-----------------------------------------------------"  Tests  "-----------------------------------------------------

//-------------------------------------"    Create Battery   "-------------------------------------
func main() {

	battery1 := createBattery(1, 4, "onLine", 60, 6, 5)

	colorReset := "\033[0m"
	colorYellow := "\033[33m"
	fmt.Println("_")
	fmt.Println(string(colorYellow), "=======================| Creating the Battery |=======================", string(colorReset))
	fmt.Println("_")
	fmt.Println("New  Battery ID = ", battery1.ID, " || Status =  ", battery1.status, " || Number of Columns =  ", battery1.amountOfColumns, " || Number of Floors =  ", battery1.amountOfFloors, " || Number of Basements =  ", battery1.amountOfBasements)
	fmt.Println("_")
	fmt.Println(string(colorYellow), "=======================| Creating the Columns |=======================", string(colorReset))
	fmt.Println("_")
	for _, column := range battery1.columnsList {
		fmt.Println("Column: ", column.name, "  ||  "+"Status: ", column.status, " || Floors served = ", strings.Trim(strings.Join(strings.Fields(fmt.Sprint(column.servedFloors)), ", "), "[]"))
	}
	battery1.scenario1()
	battery1.scenario2()
	battery1.scenario3()
	battery1.scenario4()
}

//-------------------------------------"    Scenario 1   "-------------------------------------

func (battery *Battery) scenario1() {

	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	fmt.Println("_")
	fmt.Println(string(colorRed), "=============================================| Scenario 1 |=============================================", string(colorReset))
	fmt.Println("_")
	fmt.Println(string(colorGreen), "Someone at RC wants to go to the 20th floor", string(colorReset))

	battery.columnsList[1].elevatorsList[0].currentFloor = 20
	battery.columnsList[1].elevatorsList[0].direction = "down"
	battery.columnsList[1].elevatorsList[0].status = "moving"
	battery.columnsList[1].elevatorsList[0].floorRequestList = append(battery.columnsList[1].elevatorsList[0].floorRequestList, 5)

	battery.columnsList[1].elevatorsList[1].currentFloor = 3
	battery.columnsList[1].elevatorsList[1].direction = "up"
	battery.columnsList[1].elevatorsList[1].status = "moving"
	battery.columnsList[1].elevatorsList[1].floorRequestList = append(battery.columnsList[1].elevatorsList[1].floorRequestList, 15)

	battery.columnsList[1].elevatorsList[2].currentFloor = 13
	battery.columnsList[1].elevatorsList[2].direction = "down"
	battery.columnsList[1].elevatorsList[2].status = "moving"
	battery.columnsList[1].elevatorsList[2].floorRequestList = append(battery.columnsList[1].elevatorsList[2].floorRequestList, 1)

	battery.columnsList[1].elevatorsList[3].currentFloor = 15
	battery.columnsList[1].elevatorsList[3].direction = "down"
	battery.columnsList[1].elevatorsList[3].status = "moving"
	battery.columnsList[1].elevatorsList[3].floorRequestList = append(battery.columnsList[1].elevatorsList[3].floorRequestList, 2)

	battery.columnsList[1].elevatorsList[4].currentFloor = 6
	battery.columnsList[1].elevatorsList[4].direction = "down"
	battery.columnsList[1].elevatorsList[4].status = "moving"
	battery.columnsList[1].elevatorsList[4].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 1)

	battery.assignElevator(20, "up")
}

//-------------------------------------"    Scenario 2   "-------------------------------------

func (battery *Battery) scenario2() {

	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	fmt.Println("_")
	fmt.Println(string(colorRed), "=============================================| Scenario 2 |=============================================", string(colorReset))
	fmt.Println("_")
	fmt.Println(string(colorGreen), "Someone at RC wants to go to the 36th floor", string(colorReset))

	battery.columnsList[2].elevatorsList[0].currentFloor = 1
	battery.columnsList[2].elevatorsList[0].direction = "up"
	battery.columnsList[2].elevatorsList[0].status = "stopped"
	battery.columnsList[2].elevatorsList[0].floorRequestList = append(battery.columnsList[2].elevatorsList[0].floorRequestList, 21)

	battery.columnsList[2].elevatorsList[1].currentFloor = 23
	battery.columnsList[2].elevatorsList[1].direction = "up"
	battery.columnsList[2].elevatorsList[1].status = "moving"
	battery.columnsList[2].elevatorsList[1].floorRequestList = append(battery.columnsList[2].elevatorsList[1].floorRequestList, 28)

	battery.columnsList[2].elevatorsList[2].currentFloor = 33
	battery.columnsList[2].elevatorsList[2].direction = "down"
	battery.columnsList[2].elevatorsList[2].status = "moving"
	battery.columnsList[2].elevatorsList[2].floorRequestList = append(battery.columnsList[2].elevatorsList[2].floorRequestList, 1)

	battery.columnsList[2].elevatorsList[3].currentFloor = 40
	battery.columnsList[2].elevatorsList[3].direction = "down"
	battery.columnsList[2].elevatorsList[3].status = "moving"
	battery.columnsList[2].elevatorsList[3].floorRequestList = append(battery.columnsList[2].elevatorsList[3].floorRequestList, 24)

	battery.columnsList[2].elevatorsList[4].currentFloor = 39
	battery.columnsList[2].elevatorsList[4].direction = "down"
	battery.columnsList[2].elevatorsList[4].status = "moving"
	battery.columnsList[2].elevatorsList[4].floorRequestList = append(battery.columnsList[2].elevatorsList[4].floorRequestList, 1)

	battery.assignElevator(36, "up")
}

//-------------------------------------"    Scenario 3   "-------------------------------------

func (battery *Battery) scenario3() {

	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	fmt.Println("_")
	fmt.Println(string(colorRed), "=============================================| Scenario 3 |=============================================", string(colorReset))
	fmt.Println("_")
	fmt.Println(string(colorGreen), "Someone at 54e floor wants to go to RC", string(colorReset))

	battery.columnsList[3].elevatorsList[0].currentFloor = 58
	battery.columnsList[3].elevatorsList[0].direction = "down"
	battery.columnsList[3].elevatorsList[0].status = "moving"
	battery.columnsList[3].elevatorsList[0].floorRequestList = append(battery.columnsList[3].elevatorsList[0].floorRequestList, 1)

	battery.columnsList[3].elevatorsList[1].currentFloor = 50
	battery.columnsList[3].elevatorsList[1].direction = "up"
	battery.columnsList[3].elevatorsList[1].status = "moving"
	battery.columnsList[3].elevatorsList[1].floorRequestList = append(battery.columnsList[3].elevatorsList[1].floorRequestList, 60)

	battery.columnsList[3].elevatorsList[2].currentFloor = 46
	battery.columnsList[3].elevatorsList[2].direction = "up"
	battery.columnsList[3].elevatorsList[2].status = "moving"
	battery.columnsList[3].elevatorsList[2].floorRequestList = append(battery.columnsList[3].elevatorsList[2].floorRequestList, 58)

	battery.columnsList[3].elevatorsList[3].currentFloor = 1
	battery.columnsList[3].elevatorsList[3].direction = "up"
	battery.columnsList[3].elevatorsList[3].status = "moving"
	battery.columnsList[3].elevatorsList[3].floorRequestList = append(battery.columnsList[3].elevatorsList[3].floorRequestList, 54)

	battery.columnsList[3].elevatorsList[4].currentFloor = 60
	battery.columnsList[3].elevatorsList[4].direction = "down"
	battery.columnsList[3].elevatorsList[4].status = "moving"
	battery.columnsList[3].elevatorsList[4].floorRequestList = append(battery.columnsList[3].elevatorsList[4].floorRequestList, 1)

	battery.columnsList[3].requestElevator(54, "down")
}

//-------------------------------------"    Scenario 4   "-------------------------------------

func (battery *Battery) scenario4() {

	colorReset := "\033[0m"
	colorRed := "\033[31m"
	colorGreen := "\033[32m"
	fmt.Println("_")
	fmt.Println(string(colorRed), "=============================================| Scenario 4 |=============================================", string(colorReset))
	fmt.Println("_")
	fmt.Println(string(colorGreen), "Someone at SS3 wants to go to RC", string(colorReset))

	battery.columnsList[0].elevatorsList[0].currentFloor = -4
	battery.columnsList[0].elevatorsList[0].status = "idle"

	battery.columnsList[0].elevatorsList[1].currentFloor = 1
	battery.columnsList[0].elevatorsList[1].status = "idle"

	battery.columnsList[0].elevatorsList[2].currentFloor = -3
	battery.columnsList[0].elevatorsList[2].direction = "down"
	battery.columnsList[0].elevatorsList[2].status = "moving"
	battery.columnsList[0].elevatorsList[2].floorRequestList = append(battery.columnsList[0].elevatorsList[2].floorRequestList, -5)

	battery.columnsList[0].elevatorsList[3].currentFloor = -6
	battery.columnsList[0].elevatorsList[3].direction = "up"
	battery.columnsList[0].elevatorsList[3].status = "moving"
	battery.columnsList[0].elevatorsList[3].floorRequestList = append(battery.columnsList[0].elevatorsList[3].floorRequestList, 1)

	battery.columnsList[0].elevatorsList[4].currentFloor = -1
	battery.columnsList[0].elevatorsList[4].direction = "down"
	battery.columnsList[0].elevatorsList[4].status = "moving"
	battery.columnsList[0].elevatorsList[4].floorRequestList = append(battery.columnsList[0].elevatorsList[4].floorRequestList, -6)

	battery.columnsList[0].requestElevator(-3, "up")
}
