package main

import (
	"github.com/svartlfheim/mimisbrunnr/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
