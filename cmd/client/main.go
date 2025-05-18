package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

const (
	SERVER_URL  = "http://localhost:8080/cotacao"
	TIMEOUT     = 300 * time.Millisecond
	OUTPUT_FILE = "cotacao.txt"
)

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), TIMEOUT)
	defer cancel()

	cotacao, err := fetchCotacao(ctx)
	if err != nil {
		log.Fatal("Erro ao buscar cotação:", err)
	}

	if err := saveCotacao(cotacao.Bid); err != nil {
		log.Fatal("Erro ao salvar cotação:", err)
	}

	fmt.Printf("Cotação salva com sucesso: Dólar: %s\n", cotacao.Bid)
}

func fetchCotacao(ctx context.Context) (*CotacaoResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", SERVER_URL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Erro: timeout ao buscar cotação do servidor")
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("servidor retornou status %d: %s", resp.StatusCode, string(body))
	}

	var cotacao CotacaoResponse
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		return nil, err
	}

	return &cotacao, nil
}

func saveCotacao(bid string) error {
	content := fmt.Sprintf("Dólar: %s", bid)
	return os.WriteFile(OUTPUT_FILE, []byte(content), 0644)
}
