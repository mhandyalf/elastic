package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/elastic/go-elasticsearch/v8"
	"github.com/elastic/go-elasticsearch/v8/esapi"
	"github.com/labstack/echo/v4"
)

func CreateElasticClient() (*elasticsearch.Client, error) {
	cfg := elasticsearch.Config{
		CloudID: "0980c7f861f4459283b73633a21d93dd:dXMtY2VudHJhbDEuZ2NwLmNsb3VkLmVzLmlvJDM2NTgwYmIyNmE0MzQ1MjViOTQwNTMyYjYwYWQ0NTQ5JDE0MDEzMTYyMzM1YTQ4NzNiM2IyODA2YjA5ZDI4MDEy",
		APIKey:  "YlhLdTg0b0J5bl9VMlpvdS1hUUs6aGdpb2t2QzhSTi1CM3FHUGl5bDREdw==",
	}
	es, err := elasticsearch.NewClient(cfg)
	if err != nil {
		return nil, err
	}

	return es, nil
}

func getData(c echo.Context) error {
	// Membuat klien Elasticsearch
	es, err := CreateElasticClient()
	if err != nil {
		return err
	}

	// Contoh query Elasticsearch
	var (
		buf bytes.Buffer
		res *esapi.Response
	)
	// Membuat query Elasticsearch
	query := map[string]interface{}{
		"query": map[string]interface{}{
			"match_all": map[string]interface{}{}, // Ini akan mengambil semua data
		},
	}

	// Mengkonversi query menjadi JSON
	if err := json.NewEncoder(&buf).Encode(query); err != nil {
		return err
	}

	// Mengirimkan permintaan pencarian ke Elasticsearch
	res, err = es.Search(
		es.Search.WithContext(context.Background()),
		es.Search.WithIndex("testing"), // Ganti dengan nama indeks yang sesuai
		es.Search.WithBody(&buf),
		es.Search.WithTrackTotalHits(true),
		es.Search.WithPretty(),
	)

	if err != nil {
		return err
	}
	defer res.Body.Close()

	// Membaca respons dari Elasticsearch
	if res.IsError() {
		return fmt.Errorf("error response: %s", res.Status())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return err
	}

	// Mengembalikan hasil query dari Elasticsearch
	return c.JSON(http.StatusOK, r)
}

func main() {
	e := echo.New()


	e.GET("/data", getData)
	e.Logger.Fatal(e.Start(":8080"))
}
