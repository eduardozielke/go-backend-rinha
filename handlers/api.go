package handlers

import (
	"context"
	"go-backend-rinha/datastore"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)


func Setup(ctx context.Context, cfg *viper.Viper, router *gin.Engine, mongoPessoaClient datastore.Pessoas) *gin.Engine {
	pessoaHandler := NewPessoasHandler(ctx, cfg, mongoPessoaClient)


	router.POST("/pessoa", pessoaHandler.NovaPessoa)
	router.GET("/pessoa/:id", pessoaHandler.BuscaPorId)
	router.GET("/pessoa", pessoaHandler.BuscaPessoasNomeSeguro)
	router.GET("/pessoas", pessoaHandler.ListPessoas)

	
	return router
}