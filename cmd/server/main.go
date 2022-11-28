package main

import "github.com/ismdeep/chunk-upload-demo/app/server/api"

func main() {
	if err := api.Eng.Run(":9000"); err != nil {
		panic(err)
	}
}
