package main

import (
	"bufio"
	"fmt"
	"net/rpc"
	"os"
	"os/exec"

	"./useful"
)

const ADD_CLASS_SCORE = 1
const GET_STUDENT_AVG = 2
const GET_AVG_ALL_CLASSES = 3
const GET_CLASS_AVG = 4
const EXIT = 5

func clearScreen() {
	cmd := exec.Command("cmd", "/c", "cls")
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func client(s *bufio.Reader) {
	c, err := rpc.Dial("tcp", "127.0.0.1:8044")
	var option = 0
	if err != nil {
		fmt.Println(err)
		return
	}

	for option != EXIT {
		option = getMenuOpt(s)

		switch option {
		case ADD_CLASS_SCORE:
			addClassScore(s, c)
		case GET_STUDENT_AVG:
			getStudentAvg(s, c)
		case GET_AVG_ALL_CLASSES:
			getGeneralAvg(s, c)
		case GET_CLASS_AVG:
			getClassAvg(s, c)
		}
		if option != EXIT {
			fmt.Print("Press 'Enter' to continue...")
			fmt.Scanln()
		}
		clearScreen()
	}
}

func addClassScore(s *bufio.Reader, c *rpc.Client) {
	var sGrade useful.StudentGrade
	var response string
	fmt.Println("***You can use spaces on these entries***")
	fmt.Print("Give me the class name: ")
	aux, _, _ := s.ReadLine() //Second param is for delimiter and third for error
	//We need to do this to convert the byte array into a string, alternatively using ReadString('\n') requires an extra line to remove the delimiter from the string
	sGrade.Class = string(aux)
	fmt.Print("Give me the student name: ")
	aux, _, _ = s.ReadLine() //Second param is for delimiter and third for error
	//We need to do this to convert the byte array into a string, alternatively using ReadString('\n') requires an extra line to remove the delimiter from the string
	sGrade.Student = string(aux)
	fmt.Print("Give me the grade for said class: ")
	fmt.Scanln(&sGrade.Grade) //No spaces no need to use the reader
	err := c.Call("Server.AddStudentGrade", sGrade, &response)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(response)
	}
}

func getStudentAvg(s *bufio.Reader, c *rpc.Client) {
	var avgScore float64
	fmt.Print("Give me the student name: ")
	aux, _, _ := s.ReadLine()
	sName := string(aux)
	err := c.Call("Server.GetStudentAverageScore", sName, &avgScore)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Student: %s average score is: %0.4f\n", sName, avgScore)
	}
}

func getGeneralAvg(s *bufio.Reader, c *rpc.Client) {
	var avgScore float64
	var aux string
	err := c.Call("Server.GetGeneralAverageScore", aux, &avgScore)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("The average accross all classes and students is: %0.4f\n", avgScore)
	}
}

func getClassAvg(s *bufio.Reader, c *rpc.Client) {
	var avgScore float64
	fmt.Print("Give me the class name: ")
	aux, _, _ := s.ReadLine()
	cName := string(aux)
	err := c.Call("Server.GetClassAverageScore", cName, &avgScore)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Printf("Class: %s average score is: %0.4f\n", cName, avgScore)
	}
}

func getMenuOpt(s *bufio.Reader) int {
	var option int
	fmt.Println("****School Management****")
	fmt.Println("1. Add a score to a new student")
	fmt.Println("2. Get a student average score")
	fmt.Println("3. Get a general average score")
	fmt.Println("4. Get the class average score")
	fmt.Println("5. Exit")
	fmt.Print("Select an option: ")
	fmt.Scanln(&option)
	return option
}

func main() {
	scanner := bufio.NewReader(os.Stdin)
	client(scanner)

}
