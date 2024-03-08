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

	jcp := logic.MakeJcp(10)
	err := jcp.StartCopy(src, dst)

	if err != nil {
		panic(fmt.Sprintf("Error: %v", err))
	}

	for {
		update, more := <-jcp.ProgressChannel
		if !more {
			fmt.Printf("Transfer complete\n")
			break
		}
		fmt.Printf("Update: %v", update)
		if update.JcpError != nil {
			if update.JcpError == io.EOF {
				break
			}

			panic(fmt.Sprintf("An error occurred: %v", update.JcpError.Error()))
		}
	}

	fmt.Println("Done!")
}
