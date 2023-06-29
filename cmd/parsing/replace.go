package main

import (
	"fmt"
	"github.com/eskpil/aarhus/internal/slog"
	"github.com/eskpil/aarhus/pkg/contracts"
)

func main() {
	template := contracts.Template{
		StartupCommand: "java -Xmx%MEMORY% -Xmx%CPU% -jar fabricmc-1.20.1.jar nogui",
	}

	replacements := map[string]string{
		"MEMORY": "4",
		"CPU":    "2",
	}

	output, err := template.Apply(replacements)
	if err != nil {
		slog.Fatal("missing variable", err)
	}

	fmt.Println(template.StartupCommand)
	fmt.Println(output)
}
