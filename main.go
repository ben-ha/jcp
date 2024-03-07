package main

import (
	"fmt"
	"io"
	"os"

	"github.com/ben-ha/jcp/logic"
)

func main() {
	src := os.Args[1]
	dst := os.Args[2]

	progressChannel, err := logic.StartCopy(src, dst, 10)

	if err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}

	for {
		update, more := <-progressChannel
		if !more {
			fmt.Printf("Transfer complete\n")
			break
		}
		fmt.Printf("Update: %v", update)
		if update.Error != nil {
			if update.Error == io.EOF {
				break
			}

			panic(fmt.Sprintf("An error occurred: %v", update.Error.Error()))
		}
	}

	fmt.Println("Done!")
}
