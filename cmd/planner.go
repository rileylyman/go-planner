package main

import (
	"os"
	"fmt"
	"flag"
	"bufio"
	"strings"
	"io/ioutil"
)

const dataPath string = ".planner/tasks"

var add string = "add"
var remove string = "remove"

var expectedAddArgs = 1
var expectedRemoveArgs = 1

var addCommand = flag.NewFlagSet(add, flag.ExitOnError)
var removeCommand = flag.NewFlagSet(remove, flag.ExitOnError)

func main() {

	switch prog := os.Args[1]; prog {
	case add:
		cmdAddTask()
	case remove:
		cmdRemoveTask()
	default:
		fmt.Println("Invalid argument: ", prog)
	}
}

func cmdAddTask() {

	duedate := addCommand.String("d", "", "The due date for the task, which defaults to never.")
	importance := addCommand.Int("i", 0, "The importance of the task, defaults to 0")

	parsedArgs := parseArgs(2, expectedAddArgs, addCommand)

	if len(*duedate) <= 0 {
		fmt.Println("No due date given for", parsedArgs[0], ". Continuing...")
	}

	for _, arg := range parsedArgs {
		if strings.Contains(arg, ";") || strings.Contains(arg, "`") {
			fmt.Println("Invalid argument: [you cannot name a task or duedate that]")
			os.Exit(1)
		}
	}
	switch {
	case strings.Contains(*duedate, ";") || strings.Contains(*duedate, "`"):
		fmt.Println("Invalid argument: [you cannot name a task or duedate that]")
	}

	addTask(parsedArgs[0], *duedate, *importance)
	fmt.Println(getRawTaskString())
}

func cmdRemoveTask() {

	parsedArgs := parseArgs(2, expectedRemoveArgs, removeCommand)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Deleting task", parsedArgs[0], "are you sure you want to continue? (y/n) ")
	text, _ := reader.ReadString('\n')

	if text[0] == 'y' {
		//removeTask(parsedArgs[0])
	} else  {
		fmt.Println("Aborting.")
	}

}

func addTask(taskName, taskDueDate string, taskImportance int) {
	data := getRawTaskString()
	data += taskName + ";" + taskDueDate + ";" + string(taskImportance)

	toWrite := []byte(data)
	err := ioutil.WriteFile(dataPath, toWrite, 0644)
	check(err)
}

func listTasks() {
	fmt.Println(getRawTaskString())
}

func parseArgs(argsStart, expectedArgs int, flagSet *flag.FlagSet) (parsedArgs []string) {

	if (os.Args[argsStart][0] != '-') {
		//The flags must be coming second
		parsedArgs = os.Args[argsStart : argsStart + expectedArgs]
		flagSet.Parse(os.Args[argsStart + expectedArgs:])
	} else {
		flagSet.Parse(os.Args[argsStart:])

		args := addCommand.Args()

		if len(args) < expectedArgs {
			fmt.Println("Too few arguments.")
			os.Exit(1)
		} else if len(args) > expectedAddArgs {
			fmt.Println("Too many arguments given. Maybe there are misplaced flags?")
			os.Exit(1)
		}
		parsedArgs = args[0 : expectedArgs]
	}
	return
}

func getRawTaskString() string {
	b, err := ioutil.ReadFile(dataPath)
	check(err)

	return string(b)
}

func check(e error) { if e != nil { panic(e) } }
