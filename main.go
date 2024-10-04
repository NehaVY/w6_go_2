package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

type Employee struct {
	EMPID       int    `json:"empid"`
	Name        string `json:"name"`
	Designation string `json:"designation"`
	Income      int    `json:"income"`
}

var employees = []Employee{}
var nextEMPID = 1

func main() {

	http.HandleFunc("/employees", getMyEmployee)
	http.HandleFunc("/employees/", handleMyEmployeeByID)

	fmt.Println("Server running on port 5656...")
	http.ListenAndServe(":5656", nil)
}

func getMyEmployee(w http.ResponseWriter, r *http.Request) {
	//id, _ := strconv.Atoi(r.URL.Path[len("/employees"):])

	if r.Method == "GET" {
		getAllMyEmployees(w)
	} else {
		http.Error(w, "Invalid GET method", http.StatusMethodNotAllowed)
	}
}

func handleMyEmployeeByID(w http.ResponseWriter, r *http.Request) {
	id, _ := strconv.Atoi(r.URL.Path[len("/employees/"):])

	if r.Method == "GET" {
		getMyEmployeeByID(w, id)
	} else if r.Method == "PUT" {
		updateMyEmployee(w, r)

	} else if r.Method == "POST" {
		createMyEmployee(w, r)
	} else if r.Method == "DELETE" {
		deleteMyEmployee(w, id)
	} else {
		http.Error(w, "Invalid handleMyEmployeeById Method", http.StatusMethodNotAllowed)
	}
}

func createMyEmployee(w http.ResponseWriter, r *http.Request) {
	var newEmployee Employee
	json.NewDecoder(r.Body).Decode(&newEmployee)
	newEmployee.EMPID = nextEMPID
	nextEMPID++
	employees = append(employees, newEmployee)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(newEmployee)
	fmt.Fprintf(w, "Finally new Employee is Created")
}

func getAllMyEmployees(w http.ResponseWriter) {
	json.NewEncoder(w).Encode(employees)
}

func getMyEmployeeByID(w http.ResponseWriter, id int) {
	for _, employee := range employees {
		if employee.EMPID == id {
			json.NewEncoder(w).Encode(employee)
			return
		}
	}
	http.Error(w, "Employee not found", http.StatusNotFound)
}

func updateMyEmployee(w http.ResponseWriter, r *http.Request) {
	var updatedEmployee Employee
	if err := json.NewDecoder(r.Body).Decode(&updatedEmployee); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	for i, employee := range employees {
		if employee.EMPID == updatedEmployee.EMPID {
			employees[i] = updatedEmployee
			json.NewEncoder(w).Encode(employees[i])
			return
		}
	}

	http.Error(w, "Employee not found", http.StatusNotFound)
}

func deleteMyEmployee(w http.ResponseWriter, id int) {
	for i, employee := range employees {
		if employee.EMPID == id {
			employees = append(employees[:i], employees[i+1:]...)
			//w.WriteHeader(http.StatusNoContent)
			fmt.Fprintf(w, "Employee with ID %d is deleted", employee.EMPID)
			return
		}
	}
	http.Error(w, "Employee not found", http.StatusNotFound)
}
