package tags

import (
	"encoding/base64"
	"fmt"

	"github.com/TerraDharitri/drt-go-chain-es-indexer/data"
	"github.com/TerraDharitri/drt-go-chain-es-indexer/process/elasticproc/converters"
)

// Serialize will serialize tagsCount in a way that Elasticsearch expects a bulk request
func (tc *tagsCount) Serialize(buffSlice *data.BufferSlice, index string) error {
	for tag, count := range tc.tags {
		if tag == "" {
			continue
		}

		base64Tag := base64.StdEncoding.EncodeToString([]byte(tag))
		if len(base64Tag) > converters.MaxIDSize {
			base64Tag = base64Tag[:converters.MaxIDSize]
		}
		meta := []byte(fmt.Sprintf(`{ "update" : {"_index":"%s", "_id" : "%s" } }%s`, index, converters.JsonEscape(base64Tag), "\n"))

		codeToExecute := `
			ctx._source.count += params.count; 
			ctx._source.tag = params.tag
`

		serializedDataStr := fmt.Sprintf(`{"script": {"source": "%s","lang": "painless","params": {"count": %d, "tag": "%s"}},"upsert": {"count": %d, "tag":"%s"}}`,
			converters.FormatPainlessSource(codeToExecute), count, converters.JsonEscape(tag), count, converters.JsonEscape(tag),
		)

		err := buffSlice.PutData(meta, []byte(serializedDataStr))
		if err != nil {
			return err
		}
	}

	return nil
}
