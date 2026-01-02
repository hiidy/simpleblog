package main

import (
	"fmt"

	"github.com/hiidy/simpleblog/pkg/errorsx"
)

func main() {
	errx := errorsx.New(500, "InternalError.DBConn", "Something went wrong: %s", "DB Connection failed")

	fmt.Println(errx)
}
