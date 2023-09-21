package datastore

import (
	"context"
	"errors"
	"fmt"
	"go-backend-rinha/model"
	"log"

	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Pessoas interface {
	AddPessoa(ctx context.Context, pessoa *model.Pessoa) error
	BuscaPorId(ctx context.Context, id string) (model.Pessoa, error)
	BuscaPessoasNomeSeguro(ctx context.Context, seguros string) ([]model.Pessoa, error)
	ListPessoas(ctx context.Context) ([]model.Pessoa, error)
}

type PessoaClient struct {
	client *mongo.Client
	cfg    *viper.Viper
	col    *mongo.Collection
}

func getCollection(cfg *viper.Viper, client *mongo.Client, colKey string) *mongo.Collection {
	db := cfg.GetString("mongodb.dbname")
	col := cfg.GetString(colKey)

	return client.Database(db).Collection(col)
}

func NewPessoaClient(client *mongo.Client, cfg *viper.Viper) *PessoaClient {
	return &PessoaClient{
		client: client,
		cfg:    cfg,
		col:    getCollection(cfg, client, "mongodb.dbcollections.pessoas"),
	}
}

func (c *PessoaClient) AddPessoa(ctx context.Context, pessoa *model.Pessoa) error {
	pessoa.ID = primitive.NewObjectID()
	_, err := c.col.InsertOne(ctx, pessoa)
	if err != nil {
		log.Print(fmt.Errorf("Não foi possivel adicionar um nova pessoa: %w", err))
		return err
	}
	return nil
}

func (c *PessoaClient) BuscaPorId(ctx context.Context, id string) (model.Pessoa, error) {
	var pessoa model.Pessoa
	objID, _ := primitive.ObjectIDFromHex(id)
	res := c.col.FindOne(ctx, bson.M{"_id": objID})

	if res.Err() != nil {
		if errors.Is(res.Err(), mongo.ErrNoDocuments) {
			return pessoa, nil
		}
		log.Print(fmt.Errorf("erro ao tentar encontrar a pessoa [%s]: %q", id, res.Err()))
		return pessoa, res.Err()
	}

	if err := res.Decode(&pessoa); err != nil {
		log.Print(fmt.Errorf("error decoding [%s]: %q", id, err))
		return pessoa, err
	}

	return pessoa, nil
}


func (c *PessoaClient) BuscaPessoasNomeSeguro(ctx context.Context, nomePessoaOuSeguro string) ([]model.Pessoa, error) {
	pessoas := make([]model.Pessoa, 0)

	filter := bson.M{
		"$or": []bson.M{
			{"nome": nomePessoaOuSeguro},
			{"seguros": nomePessoaOuSeguro},
		},
	}

	cur, err := c.col.Find(ctx, filter)
	if err != nil {
		log.Print(fmt.Errorf("Erro ao buscar pessoas [%s]: %w", nomePessoaOuSeguro, err))
		return nil, err
	}

	if err := cur.All(ctx, &pessoas); err != nil {
		log.Print(fmt.Errorf("Erro ao parsear: %w", err))
		return nil, err
	}

	return pessoas, nil
}

func (c *PessoaClient) InitPessoas(ctx context.Context) {
	setupIndexes(ctx, c.col, "nome")
}

func setupIndexes(ctx context.Context, collection *mongo.Collection, key string) {
	idxOpt := &options.IndexOptions{}
	idxOpt.SetUnique(true)
	mod := mongo.IndexModel{
		Keys: bson.M{
			key: 1, // index in ascending order
		},
		Options: idxOpt,
	}

	ind, err := collection.Indexes().CreateOne(ctx, mod)
	if err != nil {
		log.Fatal(fmt.Errorf("Indexes().CreateOne() ERROR: %w", err))
	} else {
		log.Printf("CreateOne() index: %s\n", ind)
	}
}

func (c *PessoaClient) ListPessoas(ctx context.Context) ([]model.Pessoa, error) {
	pessoas := make([]model.Pessoa, 0)
	cur, err := c.col.Find(ctx, bson.M{})

	if err != nil {
		log.Print(fmt.Errorf("Não foi possivel pegar todas as pessoas: %w", err))
		return nil, err
	}

	if err = cur.All(ctx, &pessoas); err != nil {
		log.Print(fmt.Errorf("não foi possivel parsear os resultado de pessoas: %w", err))
		return nil, err
	}

	return pessoas, nil
}
