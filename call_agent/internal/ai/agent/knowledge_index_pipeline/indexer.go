package knowledge_index_pipeline

import (
	"context"

	"github.com/cloudwego/eino-ext/components/indexer/milvus2"
	"github.com/cloudwego/eino/components/indexer"
	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// newIndexer component initialization function of node 'Indexer' in graph 'KnowledgeIndexing'
func newIndexer(ctx context.Context) (idr indexer.Indexer, err error) {
	embeddingIns, err := newEmbedding(ctx)
	if err != nil {
		return nil, err
	}
	// 写入 Attu 中看到的 default/biz；向量维度需与 embedding 模型一致（vision=2048）。
	config := &milvus2.IndexerConfig{
		ClientConfig: &milvusclient.ClientConfig{
			Address: "localhost:19530",
		},
		Collection: "biz",
		Vector: &milvus2.VectorConfig{
			Dimension:  2048,
			MetricType: milvus2.COSINE,
			VectorField: "vector",
		},
		Embedding: embeddingIns,
	}
	idr, err = milvus2.NewIndexer(ctx, config)
	if err != nil {
		return nil, err
	}
	return idr, nil
}
