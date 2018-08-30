package todo

import (
	"errors"
	"sync"

	"github.com/rs/xid"
)

//these are our variables
var (
	//array that holds all todo items
	list []Todo
	/* allows you to safely access/manipulate the
	data in this package across different goroutines.*/
	mtx sync.RWMutex
	//assure that a specific operation will run only once
	once sync.Once
)

//func init always executes first
func init() {
	once.Do(initialiseList)
}

//function to create a list of the struct Todo
func initialiseList() {
	list = []Todo{}
}

// Todo data structure for a task with a description of what to do
type Todo struct {
	ID       string `json:"id"`
	Task     string `json:"task"`
	Complete bool   `json:"complete"`
}

//functiom to retrieve list which is a list of the struct Todo
func Get() []Todo {
	return list
}

// Add will add a new todo based on a message
func Add(task string) string {
	t := newTodo(task)
	//t is a new task
	mtx.Lock()
	//list is now list + t which is a new task
	list = append(list, t)
	mtx.Unlock()
	return t.ID
	//this returns the new tasks ID
}

// Delete will remove a Todo from the Todo list
func Delete(id string) error {
	//delete requires a string and returns an error
	location, err := findTodoLocation(id)
	//location is the location of the task id
	if err != nil {
		return err
	}
	removeElementByLocation(location)
	return nil
}

// Complete will set the complete boolean to true, marking a todo as
// completed
func Complete(id string) error {
	//func Complete takes a string and returns an error
	location, err := findTodoLocation(id)
	//sets task location
	if err != nil {
		return err
	}
	//sets task completed based on task id location
	setTodoCompleteByLocation(location)
	return nil
}

func newTodo(tsk string) Todo {
	//
	return Todo{
		ID:       xid.New().String(),
		Task:     tsk,
		Complete: false,
	}
}

func findTodoLocation(id string) (int, error) {
	mtx.RLock()
	defer mtx.RUnlock()
	//iterate over list for matching id and find its location
	for i, t := range list {
		if isMatchingID(t.ID, id) {
			return i, nil
		}
	}
	return 0, errors.New("could not find todo based on id")
}

func removeElementByLocation(i int) {
	mtx.Lock()
	list = append(list[:i], list[i+1:]...)
	mtx.Unlock()
}

func setTodoCompleteByLocation(location int) {
	mtx.Lock()
	list[location].Complete = true
	mtx.Unlock()
}

func isMatchingID(a string, b string) bool {
	return a == b
}
