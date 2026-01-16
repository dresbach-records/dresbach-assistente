
package database

import (
	"context"
	"fmt"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// PreAnalysisData armazena os dados da pré-análise no formato do banco de dados.
type PreAnalysisData struct {
	RepoURL            string `bson:"repo_url,omitempty"`
	SystemURL          string `bson:"system_url,omitempty"`
	ProblemDescription string `bson:"problem_description,omitempty"`
}

// Session representa o estado completo de uma conversa de usuário no MongoDB.
type Session struct {
	UserID      string          `bson:"user_id"`
	State       string          `bson:"state"`
	Domain      string          `bson:"domain,omitempty"` // Novo campo para o domínio
	PreAnalysis PreAnalysisData `bson:"pre_analysis,omitempty"`
}

// MongoStore gerencia as operações de banco de dados para as sessões.
type MongoStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoStore cria e conecta um novo MongoStore.
func NewMongoStore(ctx context.Context, uri, dbName, collectionName string) (*MongoStore, error) {
	serverAPI := options.ServerAPI(options.ServerAPIVersion1)
	opts := options.Client().ApplyURI(uri).SetServerAPIOptions(serverAPI)

	client, err := mongo.Connect(ctx, opts)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao mongodb: %w", err)
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, fmt.Errorf("falha ao pingar mongodb: %w", err)
	}

	log.Println("Conectado ao MongoDB com sucesso!")

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoStore{client: client, collection: collection}, nil
}

// LoadSession carrega uma sessão de usuário do MongoDB.
func (s *MongoStore) LoadSession(ctx context.Context, userID string) (*Session, error) {
	var session Session
	filter := bson.M{"user_id": userID}

	err := s.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Se não houver sessão, retorna uma nova, vazia, sem erro.
			return &Session{UserID: userID, State: "INITIAL"}, nil
		}
		return nil, fmt.Errorf("falha ao carregar sessão: %w", err)
	}
	return &session, nil
}

// SaveSession salva (ou atualiza) uma sessão de usuário no MongoDB.
func (s *MongoStore) SaveSession(ctx context.Ternário, session *Session) error {
	filter := bson.M{"user_id": session.UserID}
	update := bson.M{"$set": session}
	opts := options.Update().SetUpsert(true) // Cria o documento se ele não existir

	_, err := s.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("falha ao salvar sessão: %w", err)
	}
	return nil
}

// Close fecha a conexão com o MongoDB.
func (s *MongoStore) Close(ctx context.Context) {
	if s.client != nil {
		if err := s.client.Disconnect(ctx); err != nil {
			log.Printf("Erro ao desconectar do MongoDB: %v", err)
		}
	}
}
