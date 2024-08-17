package rival

import (
	"fmt"
	"strings"
	"time"

	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/guregu/null/v5"
	_ "github.com/ncruces/go-sqlite3/driver"
	_ "github.com/ncruces/go-sqlite3/embed"
)

type RivalTag struct {
	Id        int    `db:"id"`
	RivalId   int    `db:"rival_id"`
	TagName   string `db:"tag_name"`
	Generated bool   `db:"generated"`
	TimeStamp int64  `db:"timestamp"`
	TagSource string `db:"tag_source"` // "LR2" | "Oraja"
}

func (r *RivalTag) String() string {
	t := time.Unix(r.TimeStamp, 0)
	return fmt.Sprintf("%s[%s]", r.TagName, t.Format("2006-01-02 15:04:05"))
}

type RivalTagService interface {
	// ---------- basic methods ----------
	FindRivalTagList(filter RivalTagFilter) ([]*RivalTag, int, error)
	FindRivalTagById(id int) (*RivalTag, error)
	InsertRivalTag(r *RivalTag) error
	DeleteRivalTagById(id int) error
	DeleteRivalTag(filter RivalTagFilter) error

	// Simple wrapper of FindRivalTagList
	// After query, open tui app and wait user select one
	ChooseOneTag(msg string, filter RivalTagFilter) (*RivalTag, error)

	// Build tags for one rival based on passed course data
	// Note: This function can be seen as re-build all generated tags
	// Note: This function must be called after loading rival's data
	BuildTags(r *RivalInfo, courseArr []*difftable.CourseInfo) error

	// Sync generated tags for one rival
	// It's equivilant to replace generated tags with parameter
	SyncGeneratedTags(r *RivalInfo, tags []*RivalTag) error
}

type RivalTagFilter struct {
	// Filtering fields
	Id        null.Int    `db:"id"`
	Name      null.String `db:"name"`
	Generated null.Bool   `db:"generated"`
	RivalId   null.Int    `db:"rival_id"`
}

func (f *RivalTagFilter) GenerateWhereClause() string {
	where := []string{"1=1"}
	if v := f.Id; v.Valid {
		where = append(where, "id=:id")
	}
	if v := f.Name; v.Valid {
		where = append(where, "name=:name")
	}
	if v := f.Generated; v.Valid {
		where = append(where, "generated=:generated")
	}
	if v := f.RivalId; v.Valid {
		where = append(where, "rival_id=:rival_id")
	}
	return strings.Join(where, " AND ")
}
