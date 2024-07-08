package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
)

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ClearTerminal() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		panic("Your platform is unsupported. I can't clear the screen.")
	}
}

func main() {
	tasks := []string{}
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Enter command (add, remove, complete, list, clear, exit)")

		userInput := scanner.Scan()
		if !userInput {
			break
		}

		command := scanner.Text()

		switch command {

		case "add":
			fmt.Println("Enter tasks to add: ")

			input := scanner.Scan()
			if !input {
				break
			}

			task := scanner.Text()
			tasks = append(tasks, task)
			fmt.Println("task added successfully")

		case "remove":
			is := list(tasks)
			if !is {
				continue
			}

			fmt.Println("Enter task number to remove: ")

			input := scanner.Scan()
			if !input {
				break
			}

			tasknumber, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("Invalid input. Please enter a number.")
				continue
			}

			if tasknumber < 1 || tasknumber > len(tasks) {
				fmt.Println("Task number out of range.")
				continue
			}
			tasks = append(tasks[:tasknumber-1], tasks[tasknumber:]...)
			fmt.Println("task removed successfully")

		case "complete":
			is := list(tasks)
			if !is {
				continue
			}

			fmt.Println("Enter task number to mark complete: ")

			input := scanner.Scan()
			if !input {
				break
			}

			tasknumber, err := strconv.Atoi(scanner.Text())
			if err != nil {
				fmt.Println("Invalid input. Please enter a number.")
				continue
			}
			if tasknumber < 1 || tasknumber > len(tasks) {
				fmt.Println("Task number out of range.")
				continue
			}
			fmt.Printf("Task '%s' marked complete.\n", tasks[tasknumber-1])
			tasks[tasknumber-1] = "Completed - " + tasks[tasknumber-1]

		case "list":
			list(tasks)

		case "clear":
			ClearTerminal()

		case "exit":
			fmt.Println("Exiting...")
			return

		default:
			fmt.Println("Invalid command.")
		}
	}
}

func list(tasks []string) bool {
	if len(tasks) == 0 {
		fmt.Println("No tasks.")
		return false
	}

	fmt.Println("Your tasks -")
	for i, task := range tasks {
		fmt.Printf("%d - %s\n", i+1, task)
	}
	return true
}
