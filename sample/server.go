package sample

import (
	"context"
	"github.com/jinzhu/gorm"
	"sync"
)

type Server struct {
	UnimplementedSampleServiceServer
	mu         sync.Mutex
	db         *gorm.DB
}

func NewServer(db *gorm.DB) *Server {
	return &Server{db: db}
}

// GetItem returns an item by id.
func (s *Server) GetItem(ctx context.Context, req *Item) (*Item, error) {
	// get item data from its id
	var itemOrm ItemORM
	s.db.Where("id=?", req.Id).First(&itemOrm)
	item, err := itemOrm.ToPB(ctx)
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// ListItems returns all item in database.
func (s *Server) ListItems(req *EmptyParam, srv SampleService_ListItemsServer) error {
	items := make([]*Item, 0)
	s.db.Find(&items)
	if items != nil {
		for _, item := range items {
			if err := srv.Send(item); err != nil {
				return err
			}
		}
	}
	return nil
}

// AddItem adds an item into db.
func (s *Server) AddItem(ctx context.Context, req *Item) (*Item, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// TODO: validate item: id existence, all fields have data or not.
	itemORM, err := req.ToORM(ctx)
	if err != nil {
		return nil, err
	}
	s.db.Create(&itemORM)
	return req, nil
}
