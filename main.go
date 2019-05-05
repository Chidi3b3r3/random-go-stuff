package main

import (
	// "github.com/chidi3b3r3/gocode/barycenter"
	"fmt"
	"time"

	// "github.com/chidi3b3r3/gocode/reddithackerclient"
	"github.com/chidi3b3r3/gocode/reddithackerserver"
)

func main() {
	// barycenter.Compute()
	startTime := time.Now()
	// reddithackerclient.Run()
	reddithackerserver.Run()
	fmt.Println(time.Since(startTime))
}
