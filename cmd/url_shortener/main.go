package main

import (
	"fmt"
	"log"

	"github.com/RomanenkoViktor080/go_url_shortener_service/internal/config"
)

func main() {
	app := config.Init()

	if err := app.Run(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(app)
}
