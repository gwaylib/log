package behavior

import "fmt"

// TODO: save to logserver
func Put(e *Event) {
	fmt.Println(string(e.ToJson()))
}
