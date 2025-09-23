package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
)

type Task struct {
	Title    string
	Status   bool
	Priority int
	Deadline string
}

var intent int
var newtask string
var key int
var List []Task
var priority int
var deadline string
var LastOutPut string

func main() {
	for {
		_ = loadFromFile()
		clearScreen()
		GetIntent()
		switch intent {
		case 1:
			ShowList()
			saveToFile()
		case 2:
			AddTask()
			saveToFile()
		case 3:
			MarkDone()
			saveToFile()
		case 4:
			RemoveTask()
			saveToFile()
		case 5:
			EditTask()
			saveToFile()
		case 6:
			FindTask()
			saveToFile()
		case 7:
			fmt.Println("Good bye!")
			return
		default:
			fmt.Println("Choose from (1-7)!")
			GetIntent()
		}
	}
}

func GetIntent() int {
	fmt.Println("\nWhat you need?")
	fmt.Println("1.Show list")
	fmt.Println("2.Add task")
	fmt.Println("3.Mark done")
	fmt.Println("4.Remove task")
	fmt.Println("5.Edit task")
	fmt.Println("6.Find task")
	fmt.Println("7.Exit")
	fmt.Println("===================================")
	fmt.Println(LastOutPut)
	fmt.Scan(&intent)
	return intent
}

func AddTask() {
	fmt.Println("Write your new task:")
	task := bufio.NewScanner(os.Stdin)
	task.Scan()
	newtask = task.Text()
	fmt.Println("Write your new task priority:")
	fmt.Scan(&priority)
	fmt.Println("Write your new task deadline: 2003.12.12")
	fmt.Scan(&deadline)
	List = append(List, Task{Title: newtask, Priority: priority, Deadline: deadline, Status: false})
	fmt.Println("Task added!")
	LastOutPut = "Task added!"
}

func MarkDone() {
	fmt.Println("Enter number of task you done:")
	fmt.Scan(&key)
	if key >= 0 && key < len(List) {
		List[key].Status = true
		fmt.Println("Marked Done!")
		LastOutPut = "Marked Done!"
	} else {
		fmt.Println("You don't have this task number!")
		LastOutPut = "You don't have this task number!"
	}
}

func RemoveTask() {
	fmt.Println("Enter number of task you want to remove:")
	fmt.Scan(&key)
	if key >= 0 && key < len(List) {
		List = append(List[:key], List[key+1:]...)
		fmt.Println("Task removed!")
		LastOutPut = "Task removed!"
	} else {
		fmt.Println("You don't have this task number!")
		LastOutPut = "You don't have this task number!"
	}
}

func EditTask() {
	fmt.Println("Enter number of task you want to edit:")
	fmt.Scan(&key)
	if key >= 0 && key < len(List) {
		fmt.Println("Enter new task title:")
		task := bufio.NewScanner(os.Stdin)
		task.Scan()
		newtask = task.Text()
		fmt.Println("Enter new priority:")
		fmt.Scan(&priority)
		fmt.Println("Enter new deadline:")
		fmt.Scan(&deadline)
		List[key] = Task{Title: newtask, Deadline: deadline, Priority: priority}
		fmt.Println("Task Edited!")
		LastOutPut = "Task Edited!"
	} else {
		fmt.Println("You don't have this task number!")
		LastOutPut = "You don't have this task number!"
	}
}

func FindTask() {
	var results []Task
	fmt.Println("What you searching for?")
	task := bufio.NewScanner(os.Stdin)
	task.Scan()
	newtask = task.Text()
	for _, t := range List {
		if strings.Contains(strings.ToLower(t.Title), strings.ToLower(newtask)) {
			results = append(results, t)
		}
		// if results set as a global variable we can use results = results[:0] to delete all indexs
	}
	var sb strings.Builder
	for i := range results {
		sb.WriteString(fmt.Sprintf("Task(%d): %s\n", i, results[i].Title))
	}
	if len(results) == 0 {
		sb.WriteString("No matching task found!\n")
	}
	LastOutPut = sb.String()
}

func saveToFile() error {
	data, err := json.MarshalIndent(List, "", "  ")
	if err != nil {
		log.Println("Error while marshaling JSON:", err)
		return err
	}
	if err := os.WriteFile("List.json", data, 0644); err != nil {
		log.Println("Error writing file:", err)
	}
	return err
}

func loadFromFile() error {
	data, err := os.ReadFile("List.json")
	if err != nil {
		return err
	}
	return json.Unmarshal(data, &List)
}

func ShowList() {
	var sb strings.Builder
	sb.WriteString("------------- TO DO LIST -------------\n")
	for i := range List {
		sb.WriteString(fmt.Sprintf(
			"Task(%d): %s\t\t[%v]\nDeadline(%s)\nPriority(%d)\n",
			i, List[i].Title, List[i].Status, List[i].Deadline, List[i].Priority,
		))
	}
	sb.WriteString("--------------------------------------\n")

	LastOutPut = sb.String()
}

func clearScreen() {
	fmt.Print("\033[H\033[2J")
}

//END LINE
