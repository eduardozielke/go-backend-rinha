package handlers

import (
	"context"
	"fmt"
	"go-backend-rinha/datastore"
	"go-backend-rinha/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type PessoasHandler struct {
	ctx          context.Context
	cfg          *viper.Viper
	mongoDBStore datastore.Pessoas
}

func NewPessoasHandler(ctx context.Context, mongoDBStore datastore.Pessoas) *PessoasHandler {
	return &PessoasHandler{
		ctx:          ctx,
		mongoDBStore: mongoDBStore,
	}
}

func (handler *PessoasHandler) NovaPessoa(ctx *gin.Context) {
	var pessoa *model.Pessoa

	if err := ctx.ShouldBindJSON(&pessoa); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"inputInvalido": err.Error(),
		})
		return
	}

	v := validator.New()
	validationErr := v.Struct(pessoa)

	if validationErr != nil {
		for _, e := range validationErr.(validator.ValidationErrors) {
			fmt.Println(e)
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"inputInvalido": validationErr.Error(),
		})
		return
	}

	if err := handler.mongoDBStore.AddPessoa(handler.ctx, pessoa); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusCreated, pessoa)
}

func (handler *PessoasHandler) BuscaPorId(ctx *gin.Context) {
	id := ctx.Param("id")

	pessoa, err := handler.mongoDBStore.BuscaPorId(handler.ctx, id)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, pessoa)
}

func (handler *PessoasHandler) BuscaPessoasNomeSeguro(ctx *gin.Context) {
	t := ctx.Query("t")

	pessoas, err := handler.mongoDBStore.BuscaPessoasNomeSeguro(handler.ctx, t)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, pessoas)
}

func (handler *PessoasHandler) ListPessoas(ctx *gin.Context) {

	pessoas, err := handler.mongoDBStore.ListPessoas(handler.ctx)

	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, pessoas)
}
