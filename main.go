package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocarina/gocsv"
)

type Tags []string

func (s *Tags) String() string {
	return strings.Join(*s, ",")
}

func (t *Tags) Set(value string) error {
	parts := strings.Split(value, ",")
	*t = append(*t, parts...)
	return nil
}

func (t Tags) MarshalText() ([]byte, error) {
	return []byte(strings.Join(t, ",")), nil
}

func (t *Tags) UnmarshalText(text []byte) error {
	if len(text) == 0 {
		*t = []string{}
		return nil
	}
	parts := strings.Split(string(text), ",")
	*t = parts
	return nil
}

type Todo struct {
	Name          string    `csv:"todo_name"`
	Tags          Tags      `csv:"tags"`
	CreatedAt     time.Time `csv:"created_at"`
	ComplitedTime time.Time `csv:"completed_at"`
	IsCompleted   bool      `csv:"is_completed"`
}

func GetTodoFile() (*os.File, error) {
	file, err := os.OpenFile("todos.csv", os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, err
	}

	info, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if info.Size() == 0 {
		header := "todo_name,tags,created_at,completed_at,is_completed\n"
		if _, err := file.WriteString(header); err != nil {
			return nil, err
		}
		file.Seek(0, 0)
	}

	return file, nil
}

func PrintTodos(todos []Todo) {
	for idx, todo := range todos {
		fmt.Printf("S.No: %d | Name: %s | Tags: %v | Completed: %v\n",
			idx, todo.Name, todo.Tags, todo.IsCompleted)
	}
}

func ListAllTodos() error {
	todoFile, err := GetTodoFile()

	if err != nil {
		panic(err)
	}

	defer todoFile.Close()

	var todoList []Todo

	if err := gocsv.UnmarshalFile(todoFile, &todoList); err != nil {
		panic(err)
	}

	if len(todoList) == 0 {
		fmt.Println("No todo left")
		return nil
	}

	PrintTodos(todoList)

	return nil
}

func AddTodo(name string, completionTime time.Time, tags Tags) error {
	file, err := GetTodoFile()

	if err != nil {
		panic(err)
	}

	defer file.Close()

	todo := &Todo{
		Name:          name,
		Tags:          tags,
		CreatedAt:     time.Now(),
		ComplitedTime: time.Now().AddDate(0, 0, 1),
	}

	var todos []*Todo
	// If file not empty, read existing todos
	if stat, _ := file.Stat(); stat.Size() > 0 {
		if err := gocsv.UnmarshalFile(file, &todos); err != nil {
			return err
		}
	}

	todos = append(todos, todo)

	// Truncate file and write back everything
	file.Truncate(0)
	file.Seek(0, 0)

	if err := gocsv.MarshalFile(&todos, file); err != nil {
		return err
	}

	fmt.Printf("Succefully added your New Todo: %s \n", todo.Name)

	return nil
}

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 0 {
		fmt.Println("Error: no command provided")
		os.Exit(1)
	}

	switch argsWithoutProg[0] {
	case "ls":
		ListAllTodos()

	case "create":
		createCmd := flag.NewFlagSet("create", flag.ExitOnError)

		nameOfTheTask := createCmd.String("name", "", "name of todo")
		var tags Tags
		createCmd.Var(&tags, "tags", "tags for create todo")

		// Parse *after* the subcommand
		createCmd.Parse(argsWithoutProg[1:])

		if *nameOfTheTask == "" {
			fmt.Println("Error: --name is required")
			createCmd.Usage()
			os.Exit(1)
		}

		if err := AddTodo(*nameOfTheTask, time.Now().AddDate(0, 0, 1), tags); err != nil {
			panic(err)
		}

		ListAllTodos()
	}
}
