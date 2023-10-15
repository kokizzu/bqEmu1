package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"cloud.google.com/go/bigquery"
	"github.com/kokizzu/goproc"
	"github.com/kokizzu/gotro/L"
	"github.com/kokizzu/gotro/S"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func main() {

	// try using binary since embeeded one failed to compile
	ctx := context.Background()
	const (
		projectID = "test"
		datasetID = "dataset1"
	)
	proc := goproc.New()
	wg := &sync.WaitGroup{}
	wg.Add(1)
	cmdId := proc.AddCommand(&goproc.Cmd{
		Program: "./bigquery-emulator-linux-amd64", // download this from: https://github.com/goccy/bigquery-emulator/releases
		Parameters: []string{
			`--project=` + projectID,
			`--dataset=` + datasetID,
			`--database=db.sqlite`,
		},
		PrefixLabel: "[BQE]",
		OnStdout: func(cmd *goproc.Cmd, s string) error {
			if S.StartsWith(s, `[bigquery-emulator] gRPC server listening at 0.0.0.0:9060`) {
				wg.Done()
				wg = nil
			}
			return nil
		},
		OnExit: func(cmd *goproc.Cmd) {
			if wg != nil {
				wg.Done()
				wg = nil
			}
		},
	})
	go proc.StartAllParallel()
	defer func() {
		L.PanicIf(proc.Kill(cmdId), `proc.Kill.defer`)
	}()

	// catch Ctrl+C
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sig
		L.PanicIf(proc.Kill(cmdId), `proc.Kill.signal`)
	}()

	wg.Wait() // wait for ready

	L.Print(`bigquery started`)

	client, err := bigquery.NewClient(
		ctx,
		projectID,
		option.WithEndpoint(`http://127.0.0.1:9050`),
		option.WithoutAuthentication(),
	)
	L.PanicIf(err, `bigquery.NewClient`)
	defer client.Close()
	L.Print(`bigquery connected`)

	it, err := client.Query(fmt.Sprintf(`
SELECT * FROM UNNEST([STRUCT("fruits" AS name, ["apple","orange"] AS items),("cars",["subaru","tesla"])])`)).Read(ctx)
	L.PanicIf(err, `client.Query.Read`)

	L.Print(`bigquery iterator`)

	var row []bigquery.Value
	err = it.Next(&row)
	L.Print(`bigquery next`)
	if err != nil {
		if errors.Is(err, iterator.Done) {
			L.Describe(row)
			return
		}
		L.PanicIf(err, `it.Next`)
	}
	fmt.Println(row[0]) // 30
}
