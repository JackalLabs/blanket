package main

import (
	"fmt"
	"os"

	"github.com/JackalLabs/blanket/blanket"
)

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		fmt.Printf("Usage: %s {provider_url}\n", os.Args[0])
		return
	}

	url := os.Args[1]

	blanket.CmdRunBlanket(url, "https://api.jackalprotocol.com")

}
