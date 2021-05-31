package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	"github.com/peterbourgon/ff/v3/ffcli"
)

var (
	rootFlagSet = flag.NewFlagSet("sqs-pub", flag.ExitOnError)
	replayer    = NewSQSMessageReplayer()
)

func init() {
	rootFlagSet.StringVar(&replayer.cfg.from, "from", "queue-name-source", "sqs queue from where messages will be sourced from")
	rootFlagSet.StringVar(&replayer.cfg.to, "to", "queue-name-destination", "sqs queue where messages will be pushed to")
	rootFlagSet.StringVar(&replayer.cfg.filters, "filters", "10104211111292", "comma separted text that can be used a message body filter")
	rootFlagSet.BoolVar(&replayer.cfg.deleteFromSource, "delete", true, "delete messages from source after successfuly pushed to destination queue")
	rootFlagSet.BoolVar(&replayer.cfg.dryrun, "dryrun", false, "a flag to run the replay in dry run mode.")
}

func main() {
	root := &ffcli.Command{
		Name:       "replay",
		ShortUsage: "sqs_pub [-from queue1 - to queue2 -filter text1,text2,...]",
		ShortHelp:  "Source message from the given queue then push to the destination queue",
		FlagSet:    rootFlagSet,
		Exec:       replayer.replay,
	}

	if err := root.ParseAndRun(context.Background(), os.Args[1:]); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

}
