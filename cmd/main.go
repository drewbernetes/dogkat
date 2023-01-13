package main

import (
	"github.com/drew-viles/k8s-e2e-tester/pkg/cmd"
	"os"
)

func main() {
	c := cmd.Generate()

	if err := c.Execute(); err != nil {
		os.Exit(1)
	}
}
