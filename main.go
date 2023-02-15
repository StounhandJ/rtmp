package main

import "fmt"

func main() {

	fmt.Println("Launching server...")
	server := NewServer(1935)

	err := server.Run()
	if err != nil {
		fmt.Println(err)
	}
}
