package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
)

func main() {
	var tasks []string

	if len(os.Args) > 1 {
		location := os.Args[1]
		var err error
		tasks, err = load(location)
		if err != nil {
			tasks = []string{}
			fmt.Println("error loading file, initiating empty list.")
		}
	}

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Println("Enter command (add, remove, complete, list, clear, save, exit)")

		userInput := scanner.Scan()
		if !userInput {
			break
		}

		command := strings.TrimSpace(scanner.Text())

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

		case "save":
			fmt.Print("directory: ")
			input := scanner.Scan()
			if !input {
				break
			}
			dir := scanner.Text()

			fmt.Print("filename: ")
			input = scanner.Scan()
			if !input {
				break
			}
			filename := scanner.Text()
			if err := save(tasks, dir, filename); err != nil {
				fmt.Println("error saving.")
			}

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

func load(location string) ([]string, error) {
	readfile, err := os.Open(location)
	if err != nil {
		fmt.Println("couldn't find your file, creating temporary list. Do save the list before exiting.")
		return nil, err
	}
	defer readfile.Close()

	fileScanner := bufio.NewScanner(readfile)
	fileScanner.Split(bufio.ScanLines)
	var tasks []string

	for fileScanner.Scan() {
		tasks = append(tasks, fileScanner.Text())
	}

	return tasks, nil
}

func save(tasks []string, dir string, filename string) error {
	dir = strings.TrimSpace(dir)
	filename = strings.TrimSpace(filename)

	dir, err := expandHomeDir(dir)
	if err != nil {
		fmt.Println("error accessing the directory.")
		return err
	}

	if dir == "" {
		dir = "."
	}
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		fmt.Println("err creating directory.")
		return err
	}

	file, err := os.Create(filepath.Join(dir, filepath.Base(filename)))
	if err != nil {
		fmt.Println("error creating file.")
		return err
	}

	defer file.Close()

	for _, task := range tasks {
		_, err := fmt.Fprintln(file, task)
		if err != nil {
			fmt.Println("error writing to the file.")
			return err
		}
	}

	fmt.Println("file saved successfully.")
	return nil
}

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

func expandHomeDir(dir string) (string, error) {
	if strings.HasPrefix(dir, "~") {
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("could not determine user home directory: %w", err)
		}
		dir = filepath.Join(homeDir, dir[1:])
	}
	return dir, nil
}
