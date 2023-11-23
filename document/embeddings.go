package document

import (
	"context"
	"fmt"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	es_document "github.com/whosonfirst/go-whosonfirst-elasticsearch/document"
)

func SFOMuseumPrepareEmbeddingsDocumentFunc() es_document.PrepareDocumentFunc {

	sfom_f := func(ctx context.Context, body []byte) ([]byte, error) {

		new_body := []byte(`{}`)

		rsp := gjson.GetBytes(body, "sfomuseum:description")

		if !rsp.Exists() {
			return nil, nil
		}

		new_body, err := sjson.SetBytes(new_body, "sfomuseum:description", []byte(rsp.String()))

		if err != nil {
			return nil, fmt.Errorf("Failed to assign sfomuseum:description, %w", err)
		}

		return new_body, nil
	}

	return sfom_f
}
