package main

import (
	"net/http"

	"fmt"

	"github.com/WestCoastOpenSource/GameStore/api"
)

func main() {
	client := api.Start()

	fmt.Println("GameStore running on port :3000")

	if err := http.ListenAndServe(":3000", client.Handler); err != nil {
		fmt.Printf(err.Error())
	}
}
