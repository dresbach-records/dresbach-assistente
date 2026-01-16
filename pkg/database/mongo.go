
package database

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoStore implementa a persistência de sessão com o MongoDB.
type MongoStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoStore cria e retorna uma nova instância de MongoStore.
func NewMongoStore(ctx context.Context, uri, dbName, collectionName string) (*MongoStore, error) {
	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao MongoDB: %w", err)
	}

	if err := client.Ping(ctx, nil); err != nil {
		return nil, fmt.Errorf("falha ao pingar MongoDB: %w", err)
	}

	collection := client.Database(dbName).Collection(collectionName)
	return &MongoStore{client: client, collection: collection}, nil
}

// Close fecha a conexão com o MongoDB.
func (s *MongoStore) Close(ctx context.Context) {
	if s.client != nil {
		s.client.Disconnect(ctx)
	}
}

// LoadSession carrega a sessão do usuário do MongoDB.
func (s *MongoStore) LoadSession(ctx context.Context, userID string) (*Session, error) {
	var session Session
	filter := bson.M{"user_id": userID}

	err := s.collection.FindOne(ctx, filter).Decode(&session)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Se a sessão não existe, cria uma nova com o estado inicial.
			initialState := "INITIAL"
			return &Session{UserID: userID, State: initialState}, nil
		}
		return nil, fmt.Errorf("falha ao carregar sessão do MongoDB: %w", err)
	}

	// Garante que o estado não seja nulo se o documento existir, mas o campo estiver vazio.
	if session.State == "" {
		session.State = "INITIAL"
	}

	return &session, nil
}

// SaveSession salva a sessão do usuário no MongoDB.
func (s *MongoStore) SaveSession(ctx context.Context, session *Session) error {
	filter := bson.M{"user_id": session.UserID}
	update := bson.M{"$set": session}
	opts := options.Update().SetUpsert(true)

	_, err := s.collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return fmt.Errorf("falha ao salvar sessão no MongoDB: %w", err)
	}
	return nil
}

// Session representa a estrutura de dados da sessão no MongoDB.
type Session struct {
	UserID      string         `bson:"user_id"`
	State       string         `bson:"state"`
	Domain      string         `bson:"domain,omitempty"`
	PreAnalysis PreAnalysisData `bson:"pre_analysis,omitempty"`
}

// PreAnalysisData armazena os dados coletados para análise.
type PreAnalysisData struct {
	RepoURL            string `bson:"repo_url,omitempty"`
	SystemURL          string `bson:"system_url,omitempty"`
	ProblemDescription string `bson:"problem_description,omitempty"`
}
