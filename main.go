package main

import (
	"context"
	"fmt"
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
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
