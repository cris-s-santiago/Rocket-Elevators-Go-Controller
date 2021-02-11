package main

import (
	"fmt"
	"math"
	"sort"
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
	status            string
	amountOfElevators int
	servedFloors      []int
	elevatorsList     []Elevator
	callButtonList    []CallButton
}

// Elevator Struct
type Elevator struct {
	ID     int
	status string
	//amountOfFloors    int		//*******Does not exist in the requeriment
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

// Teste Struct
type Teste struct {
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

	column := Column{columnID, "online", _amountOfElevatorPerColumn, servedFloorsList, []Elevator{}, []CallButton{}}
	battery.columnsList = append(battery.columnsList, column)
	columnID++

	column.createElevators(_amountOfBasements, _amountOfElevatorPerColumn)
	column.createCallButtons(_amountOfBasements, true)
}

func (battery *Battery) createColumns(_amountOfColumns int, _amountOfFloors int, _amountOfElevatorPerColumn int) {

	FloorsPerColumn := float64(_amountOfFloors / _amountOfColumns)
	amountOfFloorsPerColumn := int(math.Ceil(FloorsPerColumn))
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
		column := Column{columnID, "online", _amountOfElevatorPerColumn, servedFloorsList, []Elevator{}, []CallButton{}}
		battery.columnsList = append(battery.columnsList, column)
		columnID++

		column.createElevators(amountOfFloorsPerColumn, _amountOfElevatorPerColumn)
		column.createCallButtons(amountOfFloorsPerColumn, false)
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
	fmt.Println("- Selected Column: ", column.ID)
	elevator := column.findElevator(1, _direction)
	fmt.Println("- Selected Elevator: ", elevator.ID)
	elevator.floorRequestList = append(elevator.floorRequestList, _requestedFloor)
	elevator.move()
}

//-----------------------------------------------------"  Column  "-----------------------------------------------------

func (column *Column) createElevators(_amountOfFloors int, _amountOfElevators int) {

	for i := 1; i <= _amountOfElevators; i++ {
		elevator := Elevator{elevatorID, "idle", 1, 1, "null", Door{}, []int{}}
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

	fmt.Println("- Current column: ", columnID)
	elevator := column.findElevator(_requestedFloor, _direction)
	fmt.Println("- Selected Elevator: ", elevator.ID)
	elevator.floorRequestList = append(elevator.floorRequestList, _requestedFloor)
	elevator.sortFloorList()
	elevator.move()

	return &elevator
}

func (column *Column) findElevator(_requestedFloor int, _direction string) Elevator {

	bestElevatorInfo := BestElevatorInfo{Elevator{}, 7, 10000000}

	for _, elevator := range column.elevatorsList {

		if _requestedFloor == elevator.currentFloor && elevator.status == "stopped" {
			//The elevator is at the lobby and already has some requests. It is about to leave but has not yet departed

			bestElevatorInfo = column.checkIfElevatorISBetter(1, elevator, bestElevatorInfo, _requestedFloor)
		} else if _requestedFloor == elevator.currentFloor && elevator.status == "idle" {
			//The elevator is at the lobby and has no requests

			bestElevatorInfo = column.checkIfElevatorISBetter(2, elevator, bestElevatorInfo, _requestedFloor)

		} else if _requestedFloor > elevator.currentFloor && elevator.direction == "up" && elevator.direction == _direction {
			//The elevator is lower than me and is coming up. It means that I'm requesting an elevator to go to a basement, and the elevator ion it's way to me.

			bestElevatorInfo = column.checkIfElevatorISBetter(3, elevator, bestElevatorInfo, _requestedFloor)

		} else if _requestedFloor < elevator.currentFloor && elevator.direction == "down" && elevator.direction == _direction {
			//The elevator is above me and is coming down. It means that I'm requesting an elevator to go to a floor, and the elevator is oit's way to me

			bestElevatorInfo = column.checkIfElevatorISBetter(3, elevator, bestElevatorInfo, _requestedFloor)

		} else if _requestedFloor > elevator.currentFloor && elevator.direction == "up" {

			bestElevatorInfo = column.checkIfElevatorISBetter(4, elevator, bestElevatorInfo, _requestedFloor)

		} else if _requestedFloor < elevator.currentFloor && elevator.direction == "down" {

			bestElevatorInfo = column.checkIfElevatorISBetter(4, elevator, bestElevatorInfo, _requestedFloor)

		} else if elevator.status == "idle" {
			//The elevator is not at the first floor, but doesn't have any request

			bestElevatorInfo = column.checkIfElevatorISBetter(5, elevator, bestElevatorInfo, _requestedFloor)

		} else {
			//The elevator is not available, but still could take the call if nothing better is found

			bestElevatorInfo = column.checkIfElevatorISBetter(6, elevator, bestElevatorInfo, _requestedFloor)
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

//-------------------------------------"    Scenario 1   "-------------------------------------

// func (battery *Battery) scenario1() {

// 	battery.columnsList[1].elevatorsList[0].currentFloor = 20
// 	battery.columnsList[1].elevatorsList[0].direction = "down"
// 	battery.columnsList[1].elevatorsList[0].status = "moving"
// 	battery.columnsList[1].elevatorsList[0].floorRequestList = append(battery.columnsList[1].elevatorsList[0].floorRequestList, 5)

// 	battery.columnsList[1].elevatorsList[1].currentFloor = 3
// 	battery.columnsList[1].elevatorsList[1].direction = "up"
// 	battery.columnsList[1].elevatorsList[1].status = "moving"
// 	battery.columnsList[1].elevatorsList[1].floorRequestList = append(battery.columnsList[1].elevatorsList[1].floorRequestList, 15)

// 	battery.columnsList[1].elevatorsList[2].currentFloor = 13
// 	battery.columnsList[1].elevatorsList[2].direction = "down"
// 	battery.columnsList[1].elevatorsList[2].status = "moving"
// 	battery.columnsList[1].elevatorsList[2].floorRequestList = append(battery.columnsList[1].elevatorsList[2].floorRequestList, 1)

// 	battery.columnsList[1].elevatorsList[3].currentFloor = 15
// 	battery.columnsList[1].elevatorsList[3].direction = "down"
// 	battery.columnsList[1].elevatorsList[3].status = "moving"
// 	battery.columnsList[1].elevatorsList[3].floorRequestList = append(battery.columnsList[1].elevatorsList[3].floorRequestList, 2)

// 	battery.columnsList[1].elevatorsList[4].currentFloor = 6
// 	battery.columnsList[1].elevatorsList[4].direction = "down"
// 	battery.columnsList[1].elevatorsList[4].status = "moving"
// 	battery.columnsList[1].elevatorsList[4].floorRequestList = append(battery.columnsList[1].elevatorsList[4].floorRequestList, 1)

// 	battery.assignElevator(20, "up")
// }

func main() {
	battery1 := createBattery(1, 4, "onLine", 60, 6, 5)
	//battery.scenario1()
	//fmt.Println(baterry1)

	battery1.assignElevator(20, "up")

	//controller := NewController(1)
	//controller.TestScenario1()
	//controller.TestScenario2()
	//controller.TestScenario3()
	//controller.TestScenario4()
}
