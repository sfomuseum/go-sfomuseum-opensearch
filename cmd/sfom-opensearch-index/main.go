// sfom-opensearch-index bulk indexes one or more whosonfirst/go-whosonfirst-iterate/v2 sources in an OpenSearch database.
package main

import (
	_ "github.com/whosonfirst/go-whosonfirst-iterate-git/v2"
)

import (
	"context"
	"log"

	"github.com/sfomuseum/go-flags/flagset"
	"github.com/sfomuseum/go-flags/lookup"
	"github.com/sfomuseum/go-sfomuseum-opensearch/document"
	"github.com/whosonfirst/go-whosonfirst-iterwriter/application/iterwriter"
	os_writer "github.com/whosonfirst/go-whosonfirst-opensearch/writer"
	"github.com/whosonfirst/go-writer/v3"
)

func main() {

	fs := iterwriter.DefaultFlagSet()
	flagset.Parse(fs)

	ctx := context.Background()
	logger := log.Default()

	writer_uri, err := lookup.StringVar(fs, "writer-uri")

	if err != nil {
		logger.Fatalf("Failed to derive writer-uri flag, %v", err)
	}

	wr, err := writer.NewWriter(ctx, writer_uri)

	if err != nil {
		log.Fatalf("Failed to create new writer, %v", err)
	}

	sfom_prepare_func := document.SFOMuseumPrepareDocumentFunc()

	// To do: type/interface checking here...

	err = wr.(os_writer.DocumentWriter).AppendPrepareFunc(ctx, sfom_prepare_func)

	if err != nil {
		log.Fatalf("Failed to append SFOM prepare func, %v", err)
	}

	run_opts := &iterwriter.RunOptions{
		Logger:  logger,
		FlagSet: fs,
		Writer:  wr,
	}

	err = iterwriter.RunWithOptions(ctx, run_opts)

	if err != nil {
		log.Fatalf("Failed to run iterwriter, %v", err)
	}

}
