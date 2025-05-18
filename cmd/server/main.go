package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/evellyncosta/context-go/internal/database"
	"github.com/evellyncosta/context-go/internal/repository"
)

type CotacaoAPIResponse struct {
	USDBRL struct {
		Code       string `json:"code"`
		CodeIn     string `json:"codein"`
		Name       string `json:"name"`
		High       string `json:"high"`
		Low        string `json:"low"`
		VarBid     string `json:"varBid"`
		PctChange  string `json:"pctChange"`
		Bid        string `json:"bid"`
		Ask        string `json:"ask"`
		Timestamp  string `json:"timestamp"`
		CreateDate string `json:"create_date"`
	} `json:"USDBRL"`
}

type CotacaoResponse struct {
	Bid string `json:"bid"`
}

const (
	PORT             = ":8080"
	API_URL          = "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	API_TIMEOUT      = 2000 * time.Millisecond
	DATABASE_TIMEOUT = 1000 * time.Millisecond
	DATABASE_PATH    = "cotacoes.db"
)

func main() {
	db, err := database.NewDB(DATABASE_PATH)
	if err != nil {
		log.Fatal("Erro ao conectar ao banco de dados:", err)
	}

	repo := repository.NewCotacaoRepository(db)

	http.HandleFunc("/cotacao", func(w http.ResponseWriter, r *http.Request) {
		handleCotacao(w, r, repo)
	})

	log.Println("Servidor iniciado na porta", PORT)
	log.Fatal(http.ListenAndServe(PORT, nil))
}

func handleCotacao(w http.ResponseWriter, r *http.Request, repo *repository.CotacaoRepository) {
	ctx := r.Context()

	ctxAPI, cancelAPI := context.WithTimeout(ctx, API_TIMEOUT)
	defer cancelAPI()

	cotacao, err := fetchCotacao(ctxAPI)
	if err != nil {
		log.Printf("Erro ao buscar cotação: %v", err)
		http.Error(w, "Erro ao buscar cotação", http.StatusInternalServerError)
		return
	}

	ctxDB, cancelDB := context.WithTimeout(ctx, DATABASE_TIMEOUT)
	defer cancelDB()

	if err := repo.Save(ctxDB, cotacao.USDBRL.Bid, "USD-BRL"); err != nil {
		log.Printf("Erro ao salvar cotação no banco: %v", err)
	}

	response := CotacaoResponse{Bid: cotacao.USDBRL.Bid}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func fetchCotacao(ctx context.Context) (*CotacaoAPIResponse, error) {
	req, err := http.NewRequestWithContext(ctx, "GET", API_URL, nil)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		if ctx.Err() == context.DeadlineExceeded {
			log.Println("Erro: timeout ao buscar cotação da API")
		}
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API retornou status: %d", resp.StatusCode)
	}

	var cotacao CotacaoAPIResponse
	if err := json.NewDecoder(resp.Body).Decode(&cotacao); err != nil {
		return nil, err
	}

	return &cotacao, nil
}
