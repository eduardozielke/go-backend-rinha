package handlers

import (
	"context"
	"go-backend-rinha/datastore"

	"github.com/gin-gonic/gin"
)


func Setup(ctx context.Context, router *gin.Engine, mongoPessoaClient datastore.Pessoas) *gin.Engine {
	pessoaHandler := NewPessoasHandler(ctx, mongoPessoaClient)


	router.POST("/pessoa", pessoaHandler.NovaPessoa)
	router.GET("/pessoa/:id", pessoaHandler.BuscaPorId)
	router.GET("/pessoa", pessoaHandler.BuscaPessoasNomeSeguro)
	
	router.GET("/pessoas", pessoaHandler.ListPessoas)

	
	return router
}