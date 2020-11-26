package main

import (
	"errors"
	"fmt"
	"net"
	"net/rpc"
	"strings"

	"./useful"
)

type Server struct{}

var classes = make(map[string]map[string]float64)
var students = make(map[string]map[string]float64)

func (s *Server) AddStudentGrade(g useful.StudentGrade, response *string) error {
	//Doing this to avoid different entries that have the same name like "Algoritmia" and "algoritmia" Technically both are different keys but are the same class
	class := strings.Title(g.Class)
	fmt.Println("Received request for class: " + class)
	student := strings.Title(g.Student)
	grade := g.Grade //Invoking an attribute multiple times is more costly than assigning the value to a variable

	_, classExists := classes[class] //Checking if the class already exists
	if classExists {                 //Class found checking if student exists next
		_, studentExists := classes[class][student]
		if studentExists { //A student was found returning error message
			return errors.New("Student " + student + " already has a grade for: " + class)
		}
		//We don't need to create a new hash for this entry since we did found something in the map already meaning there's at least one entry in here.
		classes[class][student] = grade //else not required here since if it goes in the previous statement the control would end there.
	} else {
		//Creating hash for class since this is the first entry
		classes[class] = make(map[string]float64)
		classes[class][student] = grade
	}
	_, studentExists := students[student] //Checking if the student exists
	if !studentExists {                   //If it doesn't we need to create the hash for the classes since this is the first entry
		students[student] = make(map[string]float64)
	}
	students[student][class] = grade
	*response = "Grade added correctly"
	return nil
}

func (s *Server) GetStudentAverageScore(student string, averageScore *float64) error {
	fmt.Println("Received request to get student average score")
	studentVal := strings.Title(student)
	sClasses, studentExists := students[studentVal]
	if studentExists { //Checking if the student exists first
		avgScore := float64(0)
		for _, grade := range sClasses { //Iterating through all of the clases for that student
			avgScore += grade
		}
		*averageScore = avgScore / float64(len(sClasses)) //Dividing the amount of classes for said student with the grade sum of score across all classes.
	} else {
		return errors.New("Student: " + studentVal + " doesn't exist") //Student was not found
	}
	return nil
}

func (s *Server) GetGeneralAverageScore(unused string, averageScore *float64) error {
	fmt.Println("Received request to get general average score")
	avgScore := float64(0)
	classesForStudent := float64(0)
	for _, sClasses := range students { //Iterating through all of the students
		for _, grade := range sClasses { //Iterating through all of the classes for a specific student
			avgScore += grade
		}
		classesForStudent += float64(len(sClasses)) //Adding the # of classes for one student to get our global average
	}
	*averageScore = avgScore / classesForStudent
	return nil
}

func (s *Server) GetClassAverageScore(class string, averageScore *float64) error {
	fmt.Println("Received request to get class average score")
	classVal := strings.Title(class)
	cStudents, classExists := classes[classVal]
	if classExists { //Checking if the class exists first
		avgScore := float64(0)
		for _, grade := range cStudents {
			avgScore += grade
		}
		*averageScore = avgScore / float64(len(cStudents)) //Dividing the amount of students for said class.
	} else {
		return errors.New("Class: " + classVal + " doesn't exist") //Class was not found
	}
	return nil
}

func server() {
	rpc.Register(new(Server))
	ln, err := net.Listen("tcp", ":8044")
	if err != nil {
		fmt.Println(err)
	}
	for {
		c, err := ln.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		go rpc.ServeConn(c)
	}
}

func main() {
	go server()
	fmt.Println("Press Enter to finalize...")
	fmt.Scanln()
}
