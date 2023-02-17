package testing

import (
	"context"

	"gorm.io/gorm"

	sdkdatabasecontext "github.com/scribd/go-sdk/pkg/context/database"
	"github.com/scribd/go-sdk/pkg/testing/testproto"
)

type (
	TestRecord struct {
		ID   int
		Name string
	}

	TestService struct {
		testproto.UnimplementedTestServiceServer
		db *gorm.DB
	}
)

func NewTestService(db *gorm.DB) *TestService {
	return &TestService{
		db: db,
	}
}

func (s *TestService) Get(ctx context.Context, _ *testproto.GetRequest) (*testproto.GetResponse, error) {
	db, err := sdkdatabasecontext.Extract(ctx)
	if err != nil {
		return nil, err
	}

	var rec TestRecord
	if err := db.First(&rec).Error; err != nil {
		return nil, err
	}

	return &testproto.GetResponse{
		Value: rec.Name,
	}, nil
}

func (s *TestService) GetList(_ *testproto.GetListRequest, stream testproto.TestService_GetListServer) error {
	db, err := sdkdatabasecontext.Extract(stream.Context())
	if err != nil {
		return err
	}

	var rec TestRecord
	rows, err := db.Find(&rec).Rows()
	if err != nil {
		return err
	}
	for rows.Next() {
		err := stream.Send(&testproto.GetResponse{Value: rec.Name})
		if err != nil {
			return err
		}
	}

	return nil
}
