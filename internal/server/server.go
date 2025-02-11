package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/Catizard/lampghost_wails/internal/dto"
	"github.com/Catizard/lampghost_wails/internal/service"
	"github.com/Catizard/lampghost_wails/internal/vo"
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
func (s *InternalServer) RunServer() error {
	http.HandleFunc("/table/", s.tableHandler)
	http.HandleFunc("/content/", s.tableContentHandler)
	go http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
	return nil
}

func (s *InternalServer) tableHandler(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Path[len("/table/"):]
	folders, n, err := s.folderService.FindFolderList(&vo.FolderVo{
		FolderName: name,
	})
	if err != nil {
		log.Errorf("unable to query difficult table headers: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	if n == 0 {
		http.Error(w, "No such table", 404)
		return
	}
	folder := folders[0]
	header := dto.NewDiffTableHeaderExportDto(folder, fmt.Sprintf("http://127.0.0.1:%d/content/%d.json", port, folder.ID))
	data, err := json.Marshal(&header)
	if err != nil {
		log.Errorf("failed to marshal json: %v", err)
		http.Error(w, "Failed to marshal json", 500)
		return
	}
	fmt.Fprintf(w, "%s", data)
}

func (s *InternalServer) tableContentHandler(w http.ResponseWriter, r *http.Request) {
	id := string(r.URL.Path[len("/content/")])
	int_id, err := strconv.Atoi(id)
	if err != nil {
		log.Errorf("failed to convert id to a number: %v", err)
		http.Error(w, "Failed to convert id to a number", 500)
		return
	}
	rawContents, n, err := s.folderService.FindFolderContentList(&vo.FolderContentVo{FolderID: uint(int_id)})
	if err != nil {
		log.Errorf("unable to query folder contents: %v", err)
		http.Error(w, err.Error(), 500)
		return
	}
	if n == 0 {
		http.Error(w, "No such table", 404)
		return
	}
	data, err := json.Marshal(rawContents)
	if err != nil {
		log.Errorf("failed to marshal json: %v", err)
		http.Error(w, "Failed to marshal json", 500)
		return
	}
	fmt.Fprintf(w, "%s", data)
}
