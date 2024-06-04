package sqlite

import "github.com/Catizard/lampghost/internal/data/difftable"

// Ensure service implements interface
var _ difftable.DiffTableHeaderService = (*DiffTableHeaderService)(nil)

// Represents a service component for managing difficult table header
type DiffTableHeaderService struct {
	db *DB
}

func NewDiffTableHeaderService(db *DB) *DiffTableHeaderService {
	return &DiffTableHeaderService{db: db}
}

func (s *DiffTableHeaderService) InitDiffTableHeaderTable() error {
	// TODO: implement me!
	return nil
}

func (s *DiffTableHeaderService) FindList(filter difftable.DiffTableHeaderFilter) ([]*difftable.DiffTableHeader, error) {
	// TODO: implement me!
	return nil, nil
}