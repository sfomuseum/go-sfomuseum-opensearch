package document

import (
	"context"
	"fmt"
	"strings"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	es_document "github.com/whosonfirst/go-whosonfirst-elasticsearch/document"
)

func SFOMuseumPrepareEmbeddingsDocumentFunc() es_document.PrepareDocumentFunc {

	sfom_f := func(ctx context.Context, body []byte) ([]byte, error) {

		new_body := []byte(`{}`)

		to_consider := []string{
			"wof:name",
			"sfomuseum:description",
		}

		to_index := make([]string, 0)

		for _, path := range to_consider {

			rsp := gjson.GetBytes(body, path)

			if !rsp.Exists() {
				continue
			}

			to_index = append(to_index, rsp.String())
		}

		str_body := strings.Join(to_index, " ")

		new_body, err := sjson.SetBytes(new_body, "sfomuseum:description", []byte(str_body))

		if err != nil {
			return nil, fmt.Errorf("Failed to assign sfomuseum:description, %w", err)
		}

		return new_body, nil
	}

	return sfom_f
}
