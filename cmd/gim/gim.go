package main

import (
	"os"

	"github.com/padi2312/govim/pkg"
)

func main() {
	var filePath string
	if len(os.Args) >= 2 {
		filePath = os.Args[1:2][0]
	}

	editor := pkg.NewEditor()
	editor.ReadFile(filePath)
	editor.Run()
}
