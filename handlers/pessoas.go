package handlers

import (
	"context"
	"fmt"
	"go-backend-rinha/datastore"
	"go-backend-rinha/model"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/mongo"
)

type PessoasHandler struct {
	ctx          context.Context
	mongoDBStore datastore.Pessoas
}

func NewPessoasHandler(ctx context.Context, mongoDBStore datastore.Pessoas) *PessoasHandler {
	return &PessoasHandler{
		ctx:          ctx,
		mongoDBStore: mongoDBStore,
	}
}

const ERROR_DUPLICATED_KEY = 11000

func (handler *PessoasHandler) NovaPessoa(ctx *gin.Context) {
	var pessoa *model.Pessoa

	if err := ctx.ShouldBindJSON(&pessoa); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
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
			"error": validationErr.Error(),
		})
		return
	}

	if err := handler.mongoDBStore.AddPessoa(handler.ctx, pessoa); err != nil {

		if writeException, ok := err.(mongo.WriteException); ok {
			for _, writeError := range writeException.WriteErrors {
				if writeError.Code == ERROR_DUPLICATED_KEY {
					ctx.JSON(http.StatusBadRequest, gin.H{
						"error": "Chave duplicada: " + writeError.Message,
					})
					return
				}
			}
		}

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

	if t == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": "parametro t é obrigatório",
		})
		return
	}

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
