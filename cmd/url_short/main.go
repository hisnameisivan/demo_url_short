package main

import (
	"fmt"

	"github.com/hisnameisivan/demo_url_short/internal/config"
)

func main() {
	cfg := config.MustLoad()

	fmt.Println(cfg)
}
