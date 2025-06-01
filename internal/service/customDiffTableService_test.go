package service

import (
	"testing"

	"github.com/Catizard/lampghost_wails/internal/database"
	"github.com/Catizard/lampghost_wails/internal/entity"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"gorm.io/gorm"
)

// There should always be one custom difficult table which id is 1
// Initialization is done
func TestMigration(t *testing.T) {
	t.Run("SmokeTest", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		var table entity.CustomDiffTable
		if err := db.First(&table, 1).Error; err != nil {
			t.Fatalf("query: %s", err)
		}
	})
}

func TestAddCustomDiffTable(t *testing.T) {
	t.Run("SmokeTest", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		customDiffTableService := NewCustomDiffTableService(db)

		insertParam := &vo.CustomDiffTableVo{
			Name:        "Just a Name",
			Symbol:      "^",
			LevelOrders: "",
		}

		if err := customDiffTableService.AddCustomDiffTable(insertParam); err != nil {
			t.Fatalf("add: %s", err)
		}
	})

	t.Run("FastFailOnMissingFields", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		customDiffTableService := NewCustomDiffTableService(db)
		var tests = []struct {
			name  string
			input *vo.CustomDiffTableVo
		}{
			{"completely nil", nil},
			{"name should be not empty", nil},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				if err := customDiffTableService.AddCustomDiffTable(tt.input); err == nil {
					t.Fatal("expected error, got nothing")
				}
			})
		}
	})

	t.Run("FailOnDuplicateName", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		customDiffTableService := NewCustomDiffTableService(db)
		duplicateName := &vo.CustomDiffTableVo{
			Name: "lampghost",
		}
		if err := customDiffTableService.AddCustomDiffTable(duplicateName); err == nil {
			t.Fatalf("expected error, got nothing")
		}
	})
}

func TestDeleteCustomTable(t *testing.T) {
	t.Run("SmokeTest", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		customDiffTableService := NewCustomDiffTableService(db)

		if err := customDiffTableService.AddCustomDiffTable(&vo.CustomDiffTableVo{
			Name: "Lamb",
		}); err != nil {
			t.Fatalf("add: %s", err)
		}

		if err := customDiffTableService.DeleteCustomDiffTable(2); err != nil {
			t.Fatalf("delete: %s", err)
		}
	})

	t.Run("FailOnDeletingDefaultTable", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		customDiffTableService := NewCustomDiffTableService(db)
		if err := customDiffTableService.DeleteCustomDiffTable(1); err == nil {
			t.Fatalf("assert: default custom table cannot be deleted")
		}
	})
}

func TestFindCustomDiffTableList(t *testing.T) {
	t.Run("SmokeTest", func(t *testing.T) {
		db, err := database.NewMemoryDatabase()
		if err != nil {
			t.Fatalf("db: %s", err)
		}
		customDiffTableService := NewCustomDiffTableService(db)
		var tests = []struct {
			name  string
			input *vo.CustomDiffTableVo
		}{
			{"completely nil", nil},
			{"name", &vo.CustomDiffTableVo{Name: "lampghost"}},
			{"id", &vo.CustomDiffTableVo{Model: gorm.Model{ID: 1}}},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				tables, n, err := customDiffTableService.FindCustomDiffTableList(tt.input)
				if err != nil {
					t.Fatalf("findlist: %s", err)
				}
				if n == 0 {
					t.Fatalf("expected default lampghost table, got nothing")
				}
				if n > 1 {
					t.Fatalf("expected exactly one default table, got multple tables")
				}
				if tables[0].Name != "lampghost" {
					t.Fatalf("expected default lampghost table, got %s", tables[0].Name)
				}
			})
		}
	})
}
