// sfom-opensearch-index bulk indexes one or more whosonfirst/go-whosonfirst-iterate/v2 sources in an OpenSearch database.
package main

import (
	_ "github.com/whosonfirst/go-whosonfirst-iterate-git/v2"
	_ "github.com/whosonfirst/go-whosonfirst-opensearch/writer"
)

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/sfomuseum/go-flags/lookup"
	"github.com/sfomuseum/go-flags/flagset"	
	"github.com/sfomuseum/go-sfomuseum-instagram/media"
	"github.com/whosonfirst/go-whosonfirst-iterwriter/application/iterwriter"	
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
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

	// START OF move in to library
	
	sfom_f := func(ctx context.Context, body []byte) ([]byte, error) {

		im_rsp := gjson.GetBytes(body, "millsfield:images")

		if im_rsp.Exists() {

			count := len(im_rsp.Array())

			body, err = sjson.SetBytes(body, "millsfield:count_images", count)

			if err != nil {
				return nil, fmt.Errorf("Failed to assign millsfield:count_images, %w", err)
			}
		}

		// Instagram stuff
		// tl;dr is "convert IG's goofy datetime strings in RFC3339 so that Elasticsearch isn't sad"
		// See also: sfomuseum/go-sfomuseum-instagram and sfomuseum/go-sfomuseum-instagram-publish

		ig_rsp := gjson.GetBytes(body, "instagram:post")

		if ig_rsp.Exists() {

			taken_rsp := gjson.GetBytes(body, "instagram:post.taken_at")

			t, err := time.Parse(media.TIME_FORMAT, taken_rsp.String())

			if err != nil {
				return nil, fmt.Errorf("Failed to parse '%s', %w", taken_rsp.String(), err)
			}

			body, err = sjson.SetBytes(body, "instagram:post.taken_at", t.Format(time.RFC3339))

			if err != nil {
				return nil, err
			}

			tags_rsp := gjson.GetBytes(body, "instagram:post.caption.hashtags")

			if tags_rsp.Exists() {

				hashtags := make([]string, 0)

				for _, t := range tags_rsp.Array() {
					hashtags = append(hashtags, strings.ToLower(t.String()))
				}

				body, err = sjson.SetBytes(body, "instagram:post.caption.hashtags", hashtags)

				if err != nil {
					return nil, fmt.Errorf("Failed to update IG hash tags, %w", err)
				}

			}

		}

		return body, nil
	}

	// END OF move in to library

	err = wr.AppendPrepareFunc(sfom_f)

	if err != nil {
		log.Fatalf("Failed to append SFOM prepare func, %v", err)
	}

	run_opts := &iterwriter.RunOptions{
		Logger: logger,
		FlagSet: fs,
		Writer: wr,
	}

	err = iterwriter.RunWithOptions(ctx, run_opts)

	if err != nil {
		log.Fatalf("Failed to run iterwriter, %v", err)
	}

}
