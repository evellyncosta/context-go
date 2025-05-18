package repository

import (
	"context"
	"log"
	"strconv"

	"github.com/evellyncosta/context-go/internal/models"
	"gorm.io/gorm"
)

type CotacaoRepository struct {
	db *gorm.DB
}

func NewCotacaoRepository(db *gorm.DB) *CotacaoRepository {
	return &CotacaoRepository{db: db}
}

func (r *CotacaoRepository) Save(ctx context.Context, valorStr string, currency string) error {
	valor, err := strconv.ParseFloat(valorStr, 64)
	if err != nil {
		return err
	}
	
	cotacao := &models.Cotacao{
		Valor:    valor,
		Currency: currency,
	}
	
	err = r.db.WithContext(ctx).Create(cotacao).Error
	if err != nil {
		if ctx.Err() != nil {
			log.Println("Erro ao salvar cotação: timeout excedido")
			return ctx.Err()
		}
		return err
	}
	
	return nil
}