# Rocket-Elevators-Csharp-Controller
🚀Contains the Rocket Elevator files. New solution for the Commercial sector. 📈

This code was developed for the new phase of Rocket Elevators, with cutting edge technology, in this new challenge we will have a super modern elevator.
You will have a central panel on the ground floor, where when placing the desired floor, it will indicate the column and the best elevator of this column for you.
Therefore, there is no need for buttons inside the elevator car.

📌 The program to be developed is a controller set up in a building of 66 floors including 6 basements served by 4 columns of 5 cages each.
All columns will serve the lobby: (B6 to B1) + L, (2 to 20) + L, (21 to 40) + L, (41 to 60) + L.
This controller is capable of supporting two main events:

1. A person presses a call button to request an elevator, the controller selects an
available cage and it is routed to that person based on two parameters provided by
pressing the button:
- a. The floor where the person is
- b. the direction in which he wants to go (Up or Down)

❗ It should be noted that an elevator already in motion (or stopped but still
other requests to be completed) should be prioritized versus an "Idle" elevator.

2. A person on your floor asks for the elevator to go to the RC.
The supplied parameter:
- a. The floor requested
- b. the direction in which he wants to go (Up or Down)

🎯 For this, we use the Go language.
This program contains the following Classes:
- Battery, Column, Elevator, CallButton, FloorRequestButton, Door.
Each class has its own methods.

⚡In the Battery Class, we will have the following methods:
- createBasementColumn: Responsible for creating the column that will serve the basement. 
- createColumns: Responsible for creating the other columns.
- createFloorRequestButtons: Responsible for creating the buttons that will be on the panel, except the basement floors.
- createBasementFloorRequestButtons: Responsible for creating the buttons that will be on the panel, just the basement floors.
- findBestColumn: Responsible for locating which column is capable of serving your floor
- assignElevator: Responsible for handling the demand for an elevator from the central panel.

⚡In the Column Class, we will have the following methods:
- createElevator: Responsible for creating the elevator. 
- createCallButtons: Responsible for creating the buttons floors. 
- requestElevator: Responsible for handling the demand for an elevator from your current floor. 
- findBestElevator: Responsible for analyzing all column elevators and assigning points to the best elevator for that elevator request. 
- checkIfElevatorISBetter: Responsible for checking the points of the previous function and choosing the best among the column elevators.

⚡In the Elevator Class, we will have the following methods:
- movElev: Responsible for moving the elevator.
- sortFloorList: Responsible for organizing the list of floors or elevators requested.
- operateDoors: Responsible for checking the obstruction of a door.
