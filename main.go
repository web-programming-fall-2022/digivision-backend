package main

import (
	"fmt"
	"os"

	"github.com/web-programming-fall-2022/digivision-backend/cmd"
)

func main() {
	if err := cmd.New().Execute(); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
}
