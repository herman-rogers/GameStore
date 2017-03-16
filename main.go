package main

import (
	"net/http"

	"fmt"

	"github.com/WestCoastOpenSource/GameStore/api"
)

func main() {
	client := api.Start()

	fmt.Printf("GameStore running on port :3000")
	http.ListenAndServe(":3000", client.Handler)
}
