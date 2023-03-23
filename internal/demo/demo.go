package demo

import (
	"embed"
	"io"
	"net/http"
	"time"
)

//go:embed files
var files embed.FS

func ListenAndServe(addr string) (*http.Server, error) {
	fs, err := files.ReadDir("files")
	if err != nil {
		return nil, err
	}

	mux := http.NewServeMux()
	for _, f := range fs {
		h, err := files.Open("files/" + f.Name())
		if err != nil {
			return nil, err
		}
		bs, err := io.ReadAll(h)
		if err != nil {
			return nil, err
		}
		mux.HandleFunc("/"+f.Name(), func(w http.ResponseWriter, _ *http.Request) {
			_, err := w.Write(bs)
			if err != nil {
				panic(err)
			}
		})
	}

	s := http.Server{
		Addr:           addr,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() {
		if err := s.ListenAndServe(); err != nil {
			if err != http.ErrServerClosed {
				panic(err)
			}
		}
	}()
	return &s, nil
}
