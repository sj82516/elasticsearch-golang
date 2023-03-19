package main

import (
    "bytes"
    "context"
    "encoding/json"
    "fmt"
    "strings"
    
    "github.com/opensearch-project/opensearch-go"
    "github.com/opensearch-project/opensearch-go/opensearchapi"
)

const index = "test"

func main() {
    // connect to elasticsearch
    client, _ := opensearch.NewClient(opensearch.Config{
        Addresses: []string{
            "http://localhost:9200",
        },
    })
    
    // create index
    res, _ := client.Indices.Create(index)
    fmt.Println(res)
    
    // create document
    createThenSearch(res, client)
    // delete index
    deleteIndex(res, client)
    
    // create index with mapping
    res = createMapping(res, client)
    
    createThenSearch(res, client)
    
    // delete index
    deleteIndex(res, client)
}

func createMapping(res *opensearchapi.Response, client *opensearch.Client) *opensearchapi.Response {
    mapping := `{"properties":{"key":{"type":"keyword"}}}`
    res, _ = client.Indices.Create(index)
    req := opensearchapi.IndicesPutMappingRequest{
        Index: []string{index},
        Body:  bytes.NewReader([]byte(mapping)),
    }
    res, _ = req.Do(context.Background(), client.Transport)
    return res
}

func createThenSearch(res *opensearchapi.Response, client *opensearch.Client) {
    // create document
    createDocument(res, client, `{"key": "key1", "title":"Test my first document", "number": 5}`)
    createDocument(res, client, `{"key": "key1.child", "title":"Test my first document", "number": 5}`)
    
    // refresh
    client.Indices.Refresh()
    
    // search document by text field
    r := search(res, client, "title", "document")
    fmt.Println("search by title:", r)
    
    r = search(res, client, "key", "key1")
    fmt.Println("search by key:", r)
}

func deleteIndex(res *opensearchapi.Response, client *opensearch.Client) {
    deleteReq := opensearchapi.IndicesDeleteRequest{
        Index: []string{index},
    }
    res, _ = deleteReq.Do(context.Background(), client.Transport)
}

func createDocument(res *opensearchapi.Response, client *opensearch.Client, document string) {
    req := opensearchapi.IndexRequest{
        Index: index,
        Body:  strings.NewReader(document),
    }
    res, _ = req.Do(context.Background(), client.Transport)
}

type searchResponse struct {
    Hits struct {
        Total struct {
            Value int `json:"value"`
        } `json:"total"`
        Hits []struct {
            Score  float64 `json:"_score"`
            Source struct {
                Key    string `json:"key"`
                Title  string `json:"title"`
                Number int    `json:"number"`
            } `json:"_source"`
        } `json:"hits"`
    } `json:"hits"`
}

func search(res *opensearchapi.Response, client *opensearch.Client, key string, value string) searchResponse {
    s := map[string]interface{}{
        "query": map[string]interface{}{
            "match": map[string]interface{}{
                key: value,
            },
        },
    }
    
    body, _ := json.Marshal(s)
    searchReq := opensearchapi.SearchRequest{
        Index: []string{index},
        Body:  bytes.NewReader(body),
    }
    
    res, _ = searchReq.Do(context.Background(), client.Transport)
    var r searchResponse
    json.NewDecoder(res.Body).Decode(&r)
    return r
}
