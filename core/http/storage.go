package http

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/alabianca/kadbox/core"
	"io"
	"net/http"
	"os"
	"path"
)

type StorageService struct {
	Node core.NodeClient
	Protocol core.KadProtocolService
}

func enableCors(w http.ResponseWriter) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
}

func (s *StorageService) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	enableCors(w)

	switch r.Method {
	case http.MethodGet:
		s.handleGet(w, r)
	case http.MethodPost:
		s.handlePost(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
	}
}

func (s *StorageService) handlePost(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(100 << 20) // 100mgb
	file, header, err := r.FormFile("upload")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An Error occured in StorageService.handlePost -> %s (line 42)", err)
		return
	}

	defer file.Close()

	buf := new(bytes.Buffer)
	writer := core.NewStorage(buf)
	if _, err := io.Copy(writer, file); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An Error occured in StorageService.handlePost -> %s (line 52)", err)
		return
	}

	fileHash := hex.EncodeToString(writer.Sum(nil))
	targetFile, err := os.Create(path.Join(header.Filename, fileHash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An Error occured in StorageService.handlePost -> %s (line 60)", err)
		return
	}

	defer targetFile.Close()

	if _, err := io.Copy(targetFile, buf); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An Error occured in StorageService.handlePost -> %s (line 68)", err)
		return
	}

	id, err := s.Node.LocalPeerID().Marshal()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An Error occured in StorageService.handlePost -> %s (line 75)", err)
		return
	}

	if err := s.Node.PutValue(r.Context(), core.ProtocolKey(fileHash), id); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An Error occured in StorageService.handlePost -> %s (line 81)", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Published in the network. File hash %s", fileHash)
}

func (s *StorageService) handleGet(w http.ResponseWriter, r *http.Request) {
	fileHash := r.URL.Query().Get("key")
	if fileHash == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Key Not Provided"))
		return
	}

	bts, err := s.Node.GetValue(r.Context(), core.ProtocolKey(fileHash))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An Error occured in StorageService.handleGet -> %s", err)
		return
	}

	id, err := core.PeerIDFromBytes(bts)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An Error occured in StorageService.handleGet -> %s", err)
		return
	}

	stream, err := s.Node.NewStream(r.Context(), id, core.Protocol)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An Error occured in StorageService.handleGet -> %s", err)
		return
	}

	defer stream.Close()


	kadp := s.Protocol.HandleStream(stream)
	reader, errc := kadp.Want(fileHash)

	copied := make(chan struct{})
	go func() {
		defer close(copied)
		buf := new(bytes.Buffer)
		if _, err := io.Copy(buf, reader); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, "An error occured in StorageService.handleGet -> %s", err)
			return
		}
		w.WriteHeader(http.StatusOK)
		io.Copy(w, buf)
	}()

	select {
	case err := <- errc:
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An error occured in StorageService.handleGet -> %s", err)
	case <- copied:
		fmt.Println("Copied...")
	}
}