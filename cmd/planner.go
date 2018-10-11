package main

import (
	"os"
	"fmt"
	"flag"
	"bufio"
	"strings"
	"io/ioutil"
	"encoding/json"
)

const dataPath string = ".planner/tasks"
const jsonPath string = ".planner/tasks.json"

var add string 		= "add"
var remove string 	= "remove"
var show string 	= "show"
var export string 	= "export"

var expectedAddArgs 	= 1
var expectedRemoveArgs 	= 1
var expectedExportArgs 	= 0

var addCommand = flag.NewFlagSet(add, flag.ExitOnError)
var removeCommand = flag.NewFlagSet(remove, flag.ExitOnError)
var exportCommand = flag.NewFlagSet(export, flag.ExitOnError)

type Task struct  {
	Name string
	Duedate string
	Importance string
}

func main() {

	if len(os.Args) == 1 {
		fmt.Println("The Go Planner. For usage, type planner -h")
		os.Exit(1)
	}

	switch prog := os.Args[1]; prog {
	case add:
		cmdAddTask()
	case remove:
		cmdRemoveTask()
	case show:
		cmdShowTasks()
	case export:
		cmdExportTask()
	case "-h":
		fmt.Println("Valid programs: add, remove, show, and export. Use planner <program> -h for more info")
	default:
		fmt.Println("Invalid argument: ", prog)
	}
}

func cmdShowTasks() {
	listTasks()
}

func cmdAddTask() {

	duedate := addCommand.String("d", "", "The due date for the task, which defaults to never.")
	importance := addCommand.String("i", "0", "The importance of the task, defaults to 0")

	parsedArgs := parseArgs(2, expectedAddArgs, addCommand)

	if len(*duedate) <= 0 {
		fmt.Println("No due date given for", parsedArgs[0], ". Continuing...")
	}

	for _, arg := range parsedArgs {
		if strings.Contains(arg, ";") || strings.Contains(arg, "~") {
			fmt.Println("Invalid argument: [you cannot name a task or duedate that]")
			os.Exit(1)
		}
	}
	switch {
	case strings.Contains(*duedate, ";") || strings.Contains(*duedate, "~"):
		fmt.Println("Invalid argument: [you cannot name a task or duedate that]")
	}

	addTask(parsedArgs[0], *duedate, *importance)
}

func cmdRemoveTask() {

	parsedArgs := parseArgs(2, expectedRemoveArgs, removeCommand)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Deleting task [", parsedArgs[0], "]\n Are you sure you want to continue? (y/n) ")
	text, _ := reader.ReadString('\n')

	if text[0] == 'y' {
		removeTask(parsedArgs[0])
	} else  {
		fmt.Println("Aborting.")
	}

}

func cmdExportTask() {
	export_type := exportCommand.String("t", "json", "The type of your exported (default: json)")
	//exportCommand.parseArgs()

	raw_tasks := strings.Split(getRawTaskString(), ";")
	var tasks = make([]Task, 0)

	for _, task := range(raw_tasks) {
		if task == "" {continue;}
		fields := strings.Split(task, "~")

		if len(fields) == 0 {
		continue
		} else if len(fields) !=  3 {
			fmt.Println("Something went wrong trying to print tasks. The data was \n", getRawTaskString())
			os.Exit(1)
		}

		//task := Task{fields[0],fields[1],fields[2]}
		//fmt.Println(task)
		tasks = append(tasks, Task{fields[0],fields[1],fields[2]})
	}

	switch *export_type {
	case "json":
		x, _ := json.Marshal(tasks)
		toWrite := []byte(x)
		err := ioutil.WriteFile(jsonPath, toWrite, 0644)
		check(err)
		fmt.Println(string(x))
	}
}

func addTask(taskName, taskDueDate, taskImportance string) {
	data := getRawTaskString()

	taskData := ";" + taskName + "~" + taskDueDate + "~" + taskImportance

	if strings.Contains(data, taskName) {
		fmt.Print("Already added a task by this name. Would you like to update it? (y/n) ")
		reader := bufio.NewReader(os.Stdin)
		inp, _ := reader.ReadString('\n')

		if inp[0] == 'y' {
			fmt.Println("Overwriting previous task...")

			taskIndex := strings.Index(data, taskName) - 1
			left := data[:taskIndex]

			var i int
			for i = taskIndex + 1; i < len(data) && data[i] != ';'; i++ {}

			right := data[i:]

			data = left + taskData + right

		} else {
			fmt.Println("Aborting")
			os.Exit(1)
		}
	} else {
		data += taskData
	}

	toWrite := []byte(data)
	err := ioutil.WriteFile(dataPath, toWrite, 0644)
	check(err)
}

func removeTask(taskName string) {
	data := getRawTaskString()

	if strings.Contains(data, taskName) {
		fmt.Println("Removing previous task...")

		taskIndex := strings.Index(data, taskName) - 1
		left := data[:taskIndex]

		var i int
		for i = taskIndex + 1; i < len(data) && data[i] != ';'; i++ {}

		right := data[i:]

		data = left + right
	} else {
		fmt.Println("There is no task by that name to remove.")
		os.Exit(1)
	}

	toWrite := []byte(data)
	err := ioutil.WriteFile(dataPath, toWrite, 0644)
	check(err)

}

func listTasks() {

	printBanner()

	tasks := strings.Split(getRawTaskString(), ";")

	var taskName, taskDueDate, taskImportance string

	for index, task := range tasks {

		if task == "" {continue;}
		fields := strings.Split(task, "~")

		if len(fields) == 0 {
		continue
		} else if len(fields) !=  3 {
			fmt.Println("Something went wrong trying to print tasks. The data was \n", getRawTaskString())
			os.Exit(1)
		}

		taskName = fields[0]
		taskDueDate = fields[1]
		taskImportance = fields[2]

		fmt.Println("Task", index, ": |", taskName, "|")
		fmt.Println("Due by:", taskDueDate)
		fmt.Println("Importance is:", taskImportance, "\n")
	}

}

func printBanner() {
fmt.Println( "	 ,----,")
fmt.Println( "      ,/   .`|")
fmt.Println( "    ,`   .'  :                            ,-.")
fmt.Println( "  ;    ;     /                        ,--/ /|")
fmt.Println( ".'___,/    ,'                       ,--. :/ |")
fmt.Println( "|    :     |              .--.--.   :  : ' /  .--.--.")
fmt.Println( ";    |.';  ;  ,--.--.    /  /    '  |  '  /  /  /    '")
fmt.Println( "`----'  |  | /       \\  |  :  /`./  '  |  : |  :  /`./")
fmt.Println( "    '   :  ;.--.  .-. | |  :  ;_    |  |   \\|  :  ;_    ")
fmt.Println( "    |   |  ' \\__\\/: . .  \\  \\    `. '  : |. \\  \\    `. ")
fmt.Println( "    '   :  | , .--.; |   `----.    \\|  | ' \\ \\`----.   \\ ")
fmt.Println( "    |   |.' /  /  ,.  |  /  /`--'  /'  : |--'/  /`--'  /")
fmt.Println( "    '---'  ;  :   .'   \\'--'.     / ;  |,'  '--'.     /")
fmt.Println( "           |  ,     .-./  `--'---'  '--'      `--'---'")
fmt.Println( "            `--`---'                                \n\n\n")
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
