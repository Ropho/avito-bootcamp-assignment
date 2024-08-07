package test_integration

import (
	"context"

	"github.com/ozontech/allure-go/pkg/framework/suite"
)

type DataInsertSuite struct {
	suite.Suite
	ctx context.Context
	// serviceClient service_pb.DevicesAPIServiceClient
	// conn *grpc.ClientConn
}

// func (s *DataInsertSuite) BeforeAll(t provider.T) {
// 	var err error
// 	s.ctx = context.Background()

// 	t.Log("Running down migrations")
// 	testhelpers.RunDBMigrations(t, testhelpers.Migrate_DOWN)
// 	t.Log("Running up migrations")
// 	testhelpers.RunDBMigrations(t, testhelpers.Migrate_UP)

// 	s.conn, err = grpc.DialContext(s.ctx, "localhost:8081", grpc.WithTransportCredentials(insecure.NewCredentials()))
// 	if err != nil {
// 		log.Fatal(err, "fail to connect to monitoring_service with error")
// 	}

// 	s.serviceClient = service_pb.NewDevicesAPIServiceClient(s.conn)
// }

// func (s *DataInsertSuite) AfterAll(provider.T) {
// 	s.conn.Close()
// }
