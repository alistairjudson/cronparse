package main

import (
	"fmt"
	"log"

	"github.com/alistairjudson/cronparse"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:   "cronparse",
		Short: "a utility for parsing cron strings",
		Long:  "a utility to expand cron strings into the periods that it would run on",
		Run: func(cmd *cobra.Command, args []string) {
			const numArgs = 6
			if len(args) != numArgs {
				log.Fatal("please provide 6 arguments")
			}
			parsed, err := cronparse.CronParser.Parse(args[0:5])
			if err != nil {
				log.Fatal(err)
			}
			for _, part := range parsed {
				fmt.Println(part)
			}
			fmt.Printf("%-14s %s\n", "command", args[5])
		},
	}
	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
