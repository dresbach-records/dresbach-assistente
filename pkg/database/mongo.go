package database

import (
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// MongoStore implementa a persistência de sessão com o MongoDB.
type MongoStore struct {
	client     *mongo.Client
	collection *mongo.Collection
}

// NewMongoStore cria e retorna uma nova instância de MongoStore.
func NewMongoStore(ctx context.Context, uri, dbName, collectionName string) (*MongoStore, error) {
	// PASSO 3 da sua recomendação: Usar um contexto com timeout e pingar o primary.
	connectCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(connectCtx, options.Client().ApplyURI(uri))
	if err != nil {
		return nil, fmt.Errorf("falha ao conectar ao MongoDB: %w", err)
	}

	// Pinga o nó primário para verificar se a conexão foi estabelecida.
	if err := client.Ping(connectCtx, readpref.Primary()); err != nil {
		// É uma boa prática tentar desconectar se o ping falhar.
		client.Disconnect(context.Background())
		return nil, fmt.Errorf("falha ao pingar MongoDB (primary): %w", err)
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
