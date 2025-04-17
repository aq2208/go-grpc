package main

/*
Service A:
- exposes: StatusService
- calls: GreeterService (on Service B)

Service B:
- exposes: GreeterService
- calls: StatusService (on Service A)
*/

func main() {
	
}