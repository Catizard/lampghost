package impl

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/Catizard/lampghost/internal/common"
	"github.com/Catizard/lampghost/internal/config"
	"github.com/Catizard/lampghost/internal/data"
	"github.com/Catizard/lampghost/internal/data/difftable"
	"github.com/Catizard/lampghost/internal/sqlite"
	"github.com/Catizard/lampghost/internal/tui/choose"
	"github.com/charmbracelet/log"
	"github.com/guregu/null/v5"
)

// Ensure service implements interface
var _ difftable.DiffTableHeaderService = (*DiffTableHeaderService)(nil)

// Represents a service component for managing difficult table header
type DiffTableHeaderService struct {
	db *sqlite.DB
}

func NewDiffTableHeaderService(db *sqlite.DB) *DiffTableHeaderService {
	return &DiffTableHeaderService{db: db}
}

func (s *DiffTableHeaderService) FindDiffTableHeaderList(filter data.Filter) ([]*difftable.DiffTableHeader, int, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findDiffTableHeaderList(tx, filter)
}

func (s *DiffTableHeaderService) FindDiffTableHeaderById(id int) (*difftable.DiffTableHeader, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	return findDiffTableHeaderById(tx, id)
}

func (s *DiffTableHeaderService) InsertDiffTableHeader(dth *difftable.DiffTableHeader) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := insertDiffTableHeader(tx, dth); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *DiffTableHeaderService) UpdateDiffTableHeader(id int, upd difftable.DiffTableHeaderUpdate) (*difftable.DiffTableHeader, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	ret, err := updateDiffTableHeader(tx, id, upd)
	if err != nil {
		return ret, err
	} else if err := tx.Commit(); err != nil {
		return ret, err
	}
	return ret, nil
}

func (s *DiffTableHeaderService) DeleteDiffTableHeader(id int) error {
	tx, err := s.db.BeginTx()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteDiffTableHeader(tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

func (s *DiffTableHeaderService) FetchAndSaveDiffTableHeader(url string, alias string) (*difftable.DiffTableHeader, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	// 1) Prepare header
	dth, err := fetchDiffTableFromURL(url)
	if err != nil {
		return nil, err
	}
	dth.Alias = alias
	// 2) Validate
	if ex, err := existsByName(tx, dth.Name); err != nil {
		return nil, err
	} else if ex {
		log.Fatalf(`There is already a table named (or its alias matches) %s
		Use table sync command to sync table.
		Use table del command to delete table`, dth.Name)
	}
	// 3) Download data file
	// Setup data file location
	dth.DataLocation = config.JoinWorkingDirectory(dth.Name + ".json")
	downloadTableData(dth)
	// 4) Insert into database
	if err := insertDiffTableHeader(tx, dth); err != nil {
		return nil, err
	}
	// 5) Save course info
	// Note: step 5 must be executed after step 4 because of the header's id
	if err := saveCourseInfoFromTableHeader(tx, dth); err != nil {
		return nil, err
	}
	if err := tx.Commit(); err != nil {
		return nil, err
	}
	return dth, nil
}

func (s *DiffTableHeaderService) FindDiffTableHeaderListWithChoices(msg string, filter data.Filter) (*difftable.DiffTableHeader, error) {
	tx, err := s.db.BeginTx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	if arr, n, err := findDiffTableHeaderList(tx, filter); err != nil {
		return nil, err
	} else if n == 0 {
		return nil, fmt.Errorf("no table data")
	} else {
		choices := make([]string, 0)
		for _, v := range arr {
			choices = append(choices, v.String())
		}
		i := choose.OpenChooseTuiSkippable(choices, msg)
		return arr[i], nil
	}
}

func insertDiffTableHeader(tx *sqlite.Tx, dth *difftable.DiffTableHeader) error {
	ret, err := tx.NamedExec(`
		INSERT INTO difftable_header(
			data_url,
			data_location,
			last_update,
			name,
			symbol,
			alias
		)
		VALUES (:data_url, :data_location, :last_update, :name, :symbol, :alias)`, dth)
	if id, err := ret.LastInsertId(); err != nil {
		return err
	} else {
		dth.Id = int(id)
	}
	return err
}

func findDiffTableHeaderList(tx *sqlite.Tx, filter data.Filter) (_ []*difftable.DiffTableHeader, _ int, err error) {
	rows, err := tx.NamedQuery(`
		SELECT *
		FROM difftable_header
		WHERE `+filter.GenerateWhereClause(),
		filter)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	ret := make([]*difftable.DiffTableHeader, 0)
	for rows.Next() {
		dth := &difftable.DiffTableHeader{}
		if err := rows.StructScan(dth); err != nil {
			return nil, 0, err
		}

		ret = append(ret, dth)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return ret, len(ret), nil
}

func findDiffTableHeaderById(tx *sqlite.Tx, id int) (*difftable.DiffTableHeader, error) {
	arr, _, err := findDiffTableHeaderList(tx, difftable.DiffTableHeaderFilter{Id: null.IntFrom(int64(id))})
	if err != nil {
		return nil, err
	} else if len(arr) == 0 {
		return nil, fmt.Errorf("panic: no data")
	}
	return arr[0], nil
}

func updateDiffTableHeader(tx *sqlite.Tx, id int, upd difftable.DiffTableHeaderUpdate) (*difftable.DiffTableHeader, error) {
	dth, err := findDiffTableHeaderById(tx, id)
	if err != nil {
		return dth, err
	}
	dth.MergeUpdate(upd)

	if _, err := tx.NamedExec(`
		UPDATE difftable_header
		SET name = :name,
			symbol = :symbol				
		WHERE id = :id
	`, upd); err != nil {
		return dth, err
	}
	return dth, nil
}

func deleteDiffTableHeader(tx *sqlite.Tx, id int) error {
	h, err := findDiffTableHeaderById(tx, id)
	if err != nil {
		return err
	}

	filePath := h.DataLocation
	if err := os.Remove(filePath); err != nil {
		log.Warnf("Removing [%s] failed with [%v]", filePath, err)
		// Supress error
	}
	_, err = tx.Exec("DELETE FROM difftable_header WHERE id=?", id)
	return err
}

func fetchDiffTableFromURL(url string) (*difftable.DiffTableHeader, error) {
	jsonUrl := ""
	if strings.HasSuffix(url, ".html") {
		log.Infof("Fetch difficult table data from %s", url)
		resp, err := http.Get(url)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			if err := scanner.Err(); err != nil {
				return nil, err
			}
			line := strings.Trim(scanner.Text(), " ")
			// TODO: Any other cases?
			// Its pattern should be <meta name="bmstable" content="xxx.json" />
			if strings.HasPrefix(line, "<meta name=\"bmstable\"") {
				startp := strings.Index(line, "content") + len("content=\"") - 1
				if startp == -1 {
					log.Fatalf("Cannot find 'content' field in %s", url)
				}
				endp := -1
				// Finds the end position
				first := false
				for i := startp; i < len(line); i++ {
					if line[i] == '"' {
						if !first {
							first = true
						} else {
							endp = i
							break
						}
					}
				}
				if endp == -1 {
					log.Fatalf("Cannot find 'content' field in %s", url)
				}

				// Construct the json url path
				splitUrl := strings.Split(url, "/")
				splitUrl[len(splitUrl)-1] = line[startp+1 : endp]
				jsonUrl = strings.Join(splitUrl, "/")
				log.Infof("Construct json url [%s] from [%s]", jsonUrl, url)
				break
			}
		}
	} else if strings.HasSuffix(url, ".json") {
		// Okay dokey
		jsonUrl = url
	}
	if jsonUrl == "" {
		log.Fatalf("Cannot fetch %s", url)
	}
	dth := &difftable.DiffTableHeader{}
	common.FetchJson(jsonUrl, &dth)
	return dth, nil
}

func existsByName(tx *sqlite.Tx, name string) (bool, error) {
	filter := difftable.DiffTableHeaderFilter{
		Name: null.StringFrom(name),
	}
	if _, n, err := findDiffTableHeaderList(tx, filter); err != nil {
		return false, err
	} else if n > 0 {
		return true, nil
	}
	return false, nil
}

// Download difficult table's data.json file
// TODO: Use database's transaction to protect download phase?
func downloadTableData(dth *difftable.DiffTableHeader) error {
	if len(dth.DataUrl) == 0 {
		return fmt.Errorf("downloadTableData: no data url")
	}
	// 1) Create data file
	file, err := os.Create(dth.DataLocation)
	if err != nil {
		return err
	}
	defer file.Close()
	// 2) Download and save
	// TODO: if dataUrl is not start with http?
	resp, err := http.Get(dth.DataUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	io.Copy(file, resp.Body)
	return nil
}
