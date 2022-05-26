// Product Files API
//
// Documentation for Product Files API
//
//  Schemes: http
//  BasePath: /
//  Version: 1.0.0
//
//  Consumes:
//  - application/json
//
//  Produces:
//  - application/json
// swagger:meta
package handlers

import (
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/badasukerubin/go-microservices/files"
)

type Files struct {
	l     *log.Logger
	store files.Storage
}

func NewFiles(s files.Storage, l *log.Logger) *Files {
	return &Files{store: s, l: l}
}

func (f *Files) UploadMultipart(rw http.ResponseWriter, r *http.Request) {
	err := r.ParseMultipartForm(128 * 1024)

	if err != nil {
		f.l.Fatal("Bad request", "error", err)
		http.Error(rw, "Expected multipart form data", http.StatusBadRequest)
		return
	}

	id, idErr := strconv.Atoi(r.FormValue("id"))
	if idErr != nil {
		f.l.Fatal("Bad Request", "error", err)
		http.Error(rw, "Expected integer id", http.StatusBadRequest)
	}
	f.l.Print("Process form for id", id)

	ff, mh, err := r.FormFile("file")
	if err != nil {
		f.l.Fatal("Bad Request", "error", err)
		http.Error(rw, "Expected file", http.StatusBadRequest)
		return
	}
	f.saveFile(r.FormValue("id"), mh.Filename, rw, ff)
}

func (f *Files) saveFile(id, path string, rw http.ResponseWriter, r io.ReadCloser) {
	f.l.Print("Save file for product", "id", id, "path", path)

	fp := filepath.Join(id, path)
	err := f.store.Save(fp, r)
	if err != nil {
		f.l.Panic("Unable to save file", "error", err)
		http.Error(rw, "Unable to save file", http.StatusInternalServerError)
	}
}
