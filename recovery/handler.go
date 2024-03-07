package recovery

import (
	"fmt"
)

func Handler() {

	if err := recover(); err != nil {
		fmt.Println("[Recovery] panic:", err)
	}

}
