# Warehouse Simulator

## Compilation and execution

To compile and run the project at the same time, run this command:
```shell
go run WarehouseSimulator <constraints_file>
```

If you want to create an executable, run this command:
````shell
go build WarehouseSimulator
````

## Project organisation

The project was made to be as close as possible to the functioning of a real world warehouse.

All source files are part of the main package. Every file handles a specific component of the warehouse:
* driver.go: manage truck information
* worker.go: manage pallet truck actions
* path.go: implementation of an A* algorithm used by a worker to find a path to a parcel or a truck
* matrix.go: data structure used by the pathfinding algorithm representing a 2-dimensional view of the warehouse

## Strategy

For each pallet truck, we assign a worker to it and execute these tasks until we reach the end of the program or the warehouse is empty:
* find the closest parcel
* follow the route turn by turn
* take the parcel
* find the closest truck
* follow the route turn by turn
* leave the parcel
* if we reach 80 % of the maximum capacity of a truck, we make the truck leave.