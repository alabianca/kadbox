package http

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"github.com/alabianca/kadbox/core"
	"github.com/alabianca/kadbox/log"
	"io"
	"net/http"
	"os"
	"path"
)

type StorageService struct {
	Node core.NodeClient
	Protocol core.KadProtocolService
}



func (s *StorageService) ServeHTTP(w http.ResponseWriter, r *http.Request) {

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
	file, _, err := r.FormFile("upload")
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

	kbDir, err := core.GetClosestKadboxRepoRelativeToWd()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An Error occured in StorageService.handlePost -> %s (line 57)", err)
	}

	fileHash := hex.EncodeToString(writer.Sum(nil))
	targetFile, err := os.Create(path.Join(kbDir, core.StoreDirectoryName, fileHash))
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

	s.Node.Advertise(core.ProtocolKey(fileHash))

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "File identifier: %s%s\n", core.Scheme, fileHash)
}

func (s *StorageService) handleGet(w http.ResponseWriter, r *http.Request) {
	fileHash := r.URL.Query().Get("key")
	if fileHash == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Key Not Provided"))
		return
	}

	peers, err := s.Node.FindPeers(r.Context(), core.ProtocolKey(fileHash))
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Could not find any peers that advertise %s\n", fileHash)
		return
	}

	// try to connect to each of them and attempt to download the file
	if len(peers) == 0 {
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprintf(w, "Could not find any peers that advertise %s\n", fileHash)
		return
	}
	log.Info("Attemping to open stream to first peer found")

	cm := s.Node.ConnectionManager()
	// try to create a stream first. we may get through. maybe not

	stream, err := cm.NewStream(r.Context(), peers[0], core.Protocol)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "An Error occured in StorageService.handleGet -> %s\n", err)
		return
	}


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
		log.Debug("Copied...")
	}
}