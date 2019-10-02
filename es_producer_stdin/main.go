package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"github.com/elastic/go-elasticsearch"
	"github.com/elastic/go-elasticsearch/esapi"
	"log"
	"os"
	"strings"
)

func main() {
	flag.Parse()
	esIndex := flag.Arg(0)
	if len(esIndex) == 0 {
		log.Panic("Es index cannot be empty")
	}

	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		panic(err)
	}

	sc := bufio.NewScanner(os.Stdin)
	for sc.Scan() {
		input := sc.Text()

		resp, err := esapi.IndexRequest{
			Index: esIndex,
			Body: strings.NewReader(input),
		}.Do(context.Background(), es)

		if resp != nil && err == nil {
			fmt.Printf("indexed document success %d %s \n", resp.StatusCode, input)
		}
	}
}
