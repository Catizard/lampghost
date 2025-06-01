package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
	"github.com/charmbracelet/log"
)

// TODO: make port configurable
// However make it configurable might be broken in the future, since a ir connect jar
// cannot dynamically set port
const port = 7391

type InternalServer struct {
	customDiffTableService *service.CustomDiffTableService
	folderService          *service.FolderService
}

func NewInternalServer(customDiffTableService *service.CustomDiffTableService, folderService *service.FolderService) *InternalServer {
	return &InternalServer{
		customDiffTableService: customDiffTableService,
		folderService:          folderService,
	}
}

// Setup an internal server for mocking network resource (difficult table or IR)
//
// Mock a difficult table distribution server:
//  1. /table/[???].json: return a difficult table metadata which name is '???' and could be imported by beatoraja
//  2. /content/[???].json: return a difficult table's contents data which custom_table_id = '???'
//
// TODO: Mock an IR server for providing version lockable rivals import mechanism
func (s *InternalServer) RunServer() error {
	// TODO: make the name as the parameter
	http.HandleFunc("/table/", s.tableHandler)
	http.HandleFunc("/content/", s.tableContentHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return nil
}

func (s *InternalServer) tableHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	suffix := strings.TrimPrefix(path, "/table/")
	log.Debugf("suffix: %s", suffix)
	if !strings.HasSuffix(suffix, ".json") {
		http.Error(w, "assert: table must be suffixed with .json", 500)
		return
	}
	tableName := suffix[:strings.Index(suffix, ".json")]
	log.Debugf("table name: %s", tableName)
	customTables, _, err := s.customDiffTableService.FindCustomDiffTableList(&vo.CustomDiffTableVo{
		Name: tableName,
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("faile to export table[%s]: %v", tableName, err), 500)
		return
	}
	if len(customTables) != 1 {
		http.Error(w, fmt.Sprintf("cannot determinate table[%s]: %v", tableName, err), 500)
		return
	}
	customTable := customTables[0]
	export := &dto.DiffTableHeaderExportDto{
		Name:      customTable.Name,
		Symbol:    customTable.Symbol,
		HeaderUrl: fmt.Sprintf("http://localhost:%d/table/%s.json", port, tableName),
		DataUrl:   fmt.Sprintf("http://localhost:%d/content/%d.json", port, customTable.ID),
	}
	data, err := json.Marshal(export)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to marshal: %s", err), 500)
		return
	}
	fmt.Fprintf(w, "%s", data)
}

func (s *InternalServer) tableContentHandler(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	suffix := strings.TrimPrefix(path, "/content/")
	if !strings.HasSuffix(suffix, ".json") {
		http.Error(w, "assert: table must be suffixed with .json", 500)
		return
	}
	tableIDStr := suffix[:strings.Index(suffix, ".json")]
	tableID, err := strconv.Atoi(tableIDStr)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to export table[%d]: %s", tableID, err), 500)
		return
	}
	folders, _, err := s.folderService.FindFolderTree(&vo.FolderVo{
		CustomTableID: uint(tableID),
	})
	if err != nil {
		http.Error(w, fmt.Sprintf("query folder tree: %v", err), 500)
		return
	}
	contents := make([]*dto.DiffTableDataExportDto, 0)
	for _, folder := range folders {
		for _, folderContent := range folder.Contents {
			def := &dto.DiffTableDataExportDto{
				Level:   folder.FolderName,
				Sha256:  folderContent.Sha256,
				Md5:     folderContent.Md5,
				Title:   folderContent.Title,
				Comment: folderContent.Comment,
			}
			contents = append(contents, def)
		}
	}
	data, err := json.Marshal(contents)
	if err != nil {
		log.Errorf("failed to marshal json: %v", err)
		http.Error(w, "Failed to marshal json", 500)
		return
	}
	fmt.Fprintf(w, "%s", data)
}
