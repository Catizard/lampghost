package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/charmbracelet/log"
)

// TODO: make port configurable
const port = 7391

type InternalServer struct {
	folderService *service.FolderService
}

func NewInternalServer(folderService *service.FolderService) *InternalServer {
	return &InternalServer{
		folderService: folderService,
	}
}

// Setup an internal server for mocking network resource (e.g difficult table)
//
// This functionality contains two interface to mock a favorite folder mechanism
// 1) /table/lampghost.json: mock a difficult table which name is `lampghost` and would be imported by beatoraja
// 2) /table/content.json: convert every user's custom folder to one level folder within the difficult table
func (s *InternalServer) RunServer() error {
	http.HandleFunc("/table/lampghost.json", s.tableHandler)
	http.HandleFunc("/table/content.json", s.tableContentHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return nil
}

func (s *InternalServer) tableHandler(w http.ResponseWriter, r *http.Request) {
	header := dto.NewDiffTableHeaderExportDto("lampghost", fmt.Sprintf("http://127.0.0.1:%d/table/content.json", port))
	data, err := json.Marshal(&header)
	if err != nil {
		log.Errorf("failed to marshal json: %v", err)
		http.Error(w, "Failed to marshal json", 500)
		return
	}
	fmt.Fprintf(w, "%s", data)
}

func (s *InternalServer) tableContentHandler(w http.ResponseWriter, r *http.Request) {
	folders, _, err := s.folderService.FindFolderTree()
	if err != nil {
		http.Error(w, fmt.Sprintf("query folder tree: %v", err), 500)
		return
	}
	contents := make([]*dto.DiffTableDataExportDto, 0)
	for _, folder := range folders {
		for _, folderContent := range folder.Contents {
			def := &dto.DiffTableDataExportDto{
				Level:  folder.FolderName,
				Sha256: folderContent.Sha256,
				Md5:    folderContent.Md5,
				Title:  folderContent.Title,
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
