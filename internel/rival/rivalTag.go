package rival

import (
	"github.com/Catizard/lampghost/internel/common"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type RivalTag struct {
	TagName   string
	Generated bool
	TimeStamp int64
}

func InitRivalTagTable() error {
	db, err := sqlx.Open("sqlite3", common.DBFileName)
	if err != nil {
		return err
	}
	defer db.Close()
	_, err = db.Exec("DROP TABLE IF EXISTS 'rival_tag';CREATE TABLE rival_tag (id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT, rival_id INTEGER NOT NULL, tag_name TEXT(255) NOT NULL, 'generated' INTEGER DEFAULT (0) NOT NULL, 'timestamp' TEXT NOT NULL)")
	return err
}
