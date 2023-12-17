package connections

import (
	"github.com/boltdb/bolt"
	"github.com/minio/madmin-go"
	"github.com/minio/minio-go/v7"
	"github.com/openstadia/openstadia-storage/models"
)

type ConnectionsStore struct {
	connections *models.Connections
}

func CreateConnectionsStore(connections *models.Connections) ConnectionsStore {
	return ConnectionsStore{connections: connections}
}

func (s *ConnectionsStore) GetMinioClient() *minio.Client {
	return s.connections.MinioClient
}

func (s *ConnectionsStore) GetMadminClient() *madmin.AdminClient {
	return s.connections.MadminClient
}

func (s *ConnectionsStore) GetDatabase() *bolt.DB {
	return s.connections.Database
}
