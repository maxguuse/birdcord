package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

func main() {
	err := os.RemoveAll("generated")
	if err != nil {
		log.Fatal("Error removing directory 'generated':", err)
	}

	files, err := os.ReadDir("proto")
	if err != nil {
		log.Fatal("Error reading directory 'proto':", err)
	}

	for _, file := range files {
		if file.Name() == "google" {
			continue
		}

		name := strings.TrimSuffix(file.Name(), ".proto")

		err = os.MkdirAll("generated/"+name, os.ModePerm)
		if err != nil {
			log.Fatal("Error creating directory 'generated/"+name+"':", err)
		}

		buildString := fmt.Sprintf(
			"protoc "+
				"--go_opt=paths=source_relative "+
				"--go-grpc_opt=paths=source_relative "+
				"--go_out=generated/%s "+
				"--go-grpc_out=generated/%s "+
				"--proto_path=proto "+
				"%s.proto",
			name, name, name,
		)

		var c *exec.Cmd

		switch runtime.GOOS {
		case "windows":
			c = exec.Command("cmd", "/C", buildString)
		default:
			c = exec.Command("sh", "-c", buildString)
		}

		err = c.Run()
		if err != nil {
			log.Fatal("Error running protoc:", err)
		}

		fmt.Printf("✅  Generated %s proto definitions for Go.\n", name)
	}
}
