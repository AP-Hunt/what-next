package main

import (
	"context"
	"os"

	"github.com/AP-Hunt/what-next/m/cmd"
	cmdContext "github.com/AP-Hunt/what-next/m/context"
)

func main() {

	ctx, err := cmdContext.CreateDefaultCommandContext(context.Background())
	if err != nil {
		panic(err)
	}

	if err := cmd.ExecuteC(ctx); err != nil {
		os.Exit(1)
	}
}
