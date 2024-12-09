package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
	"bufio"
	"os"
	"strconv"
)

type Task struct {
	ID        int       `json:"id"`
	Title     string    `json:"title"`
	UpdatedAt time.Time `json:"updatedAt"`
	CreatedAt time.Time `json:"createdAt"`
	Status    string    `json:"status"`
}



func main () {
	if len(os.Args) < 2 || os.Args[1] != "start" {
		fmt.Println("You start cli with 'start' command")
		os.Exit(1)
	}


	fmt.Println("Program started. Type 'help' for commands.")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("> ")

		scanner.Scan()

		command := strings.TrimSpace(scanner.Text())

		switch  {
		case command == "start":
			fmt.Println("Run the program")
		case  command == "list":
			if _, err := os.Stat("list.json"); os.IsNotExist(err) {
				createEmptyFile()
				fmt.Println("Empty file created")
			} else {
				handleListFile()
			}
		case strings.HasPrefix(command, "add "):
			handleAddTask(command[4:])
		case strings.HasPrefix(command, "delete "):
			handleDeleteFile(command[7:])
		case strings.HasPrefix(command, "update "):
			handleUpdateFile(command[7:])
		case command == "help":
			fmt.Println("Available commands:")
			fmt.Println("help  - Show all commands")
			fmt.Println("exit  - Exit the program")
			fmt.Println("list  - Show tasks(id, title, status)")
			fmt.Println("add <title>  - Create new task")
			fmt.Println("update <id>  - Update task status")
			fmt.Println("delete <id>  - Delete task")
		case command == "exit":
			fmt.Println("Exiting the program...")
			os.Exit(0)
		default:
			fmt.Printf("Unknown command: %s\n", command)
		}
	}
}

func handleUpdateFile(id string) {
	intId, err := strconv.Atoi(id)

	if err != nil {
		fmt.Println("Invalid ID")
		return
	}
	
	var tasks []Task = readTasksFromFile()

	if len(tasks) == 0 {
		fmt.Println("List is empty")
		return
	}

	var updated = false;

	for i := range tasks {
		if tasks[i].ID == intId {
			if tasks[i].Status == "todo" {
				tasks[i].Status = "complete"
			} else {
				tasks[i].Status = "todo"
			}

			tasks[i].UpdatedAt = time.Now()
			updated = true
		}
	}

	if updated {
		writeTasksToFile(tasks)
		fmt.Printf("Task witn id:%d updated \n", intId)
	} else {
		println("Id not found")
	}
}

func handleDeleteFile(id string) {
	intId, err := strconv.Atoi(id)
	if err != nil {
		fmt.Println("Invalid ID")
		return
	}

	var tasks []Task = readTasksFromFile()

	if len(tasks) == 0 {
		fmt.Println("List is empty")
		return
	}

	var formatingTasks []Task = []Task{}

	for _, t := range tasks {
		if t.ID != intId {
			formatingTasks = append(formatingTasks, t)
		}
	}

	if len(formatingTasks) == len(tasks) {
		println("Id not found")
		return
	}

	writeTasksToFile(formatingTasks)
	fmt.Printf("Task witn id:%d deleted \n", intId)
}

func handleListFile() {
	file,err := os.Open("list.json") 
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}

	defer file.Close()
 
	var tasks[]Task = readTasksFromFile() 

	if len(tasks) == 0 {
		fmt.Println("Tasks list is empty")
	} else {
		for _, task := range tasks {
			fmt.Printf("%d. %s (Status: %s, Created: %s, Updated: %s)\n", 
				task.ID, task.Title, task.Status, task.CreatedAt.Format(time.RFC3339), task.UpdatedAt.Format(time.RFC3339))
		}
	}

}

func createEmptyFile () {
	file, err := os.Create("list.json")
	if err != nil {
		fmt.Println("Error creating file:", err)
		return;
	}

	defer file.Close()

	var tasks []Task

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", " ")
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Println("Error writing:", err)
	}
}

func readTasksFromFile() []Task {
	file, err := os.Open("list.json")
	if err != nil {
		fmt.Println("Reading file error:", err)
	}
	defer file.Close()

	var tasks []Task
	decoder := json.NewDecoder(file)
	err = decoder.Decode(&tasks)
	if err != nil {
		fmt.Println("Decoding error:", err)
	}

	return tasks
}

func generateId() int {
	tasks := readTasksFromFile()
	if len(tasks) == 0 {
		return 1 
	}

	var maxId int
	for _, t := range tasks {
		if t.ID > maxId {
			maxId = t.ID
		}
	}
	return maxId + 1 
}

func writeTasksToFile(tasks []Task) {
	file, err := os.OpenFile("list.json", os.O_RDWR|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error opening file for writing:", err)
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	err = encoder.Encode(tasks)
	if err != nil {
		fmt.Println("Error writing to file:", err)
	}
}

func handleAddTask(title string) {
	if title == "" {
		fmt.Println("Error task title cannot be empty.")
	}

	var tasks []Task

	if _, err := os.Stat("list.json"); os.IsNotExist(err) {
		createEmptyFile()
		tasks = []Task{}
	} else {
		tasks = readTasksFromFile()
	}

	newTask := Task{
		ID: generateId(),
		Title: title,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Status: "todo",
	}
	
	tasks = append(tasks, newTask)
	writeTasksToFile(tasks)


	fmt.Printf("Task '%s' added with ID %d.\n", title, newTask.ID)
}

