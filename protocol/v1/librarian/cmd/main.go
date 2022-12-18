package main

import (
	"github.com/sirupsen/logrus"
	librarian "github.com/tome-gg/librarian/protocol/v1/librarian"
	validator "github.com/tome-gg/librarian/protocol/v1/librarian/validator"
)

func main() {
	directory, err := librarian.Parse("/home/darren/go/src/github.com/tome-gg/librarian/protocol/v1/template")

	if err != nil {
		panic(err)
	}

	validator.Validate(directory)
	
	logrus.Infof("Validating path: %s\n%s", directory.Path, directory.Status())
}