package server

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

type (
	Server struct {
		Pagination PaginationOptions
		Routes     []Route
		BasePath   string
		Port       int
	}

	Error struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
)

func (s Server) Start() {
	mux := http.NewServeMux()

	for _, route := range s.Routes {

		// get all of route
		mux.HandleFunc(s.BasePath+route.Path, s.createCollectionHandler(route))

		// get specific id
		mux.HandleFunc(s.BasePath+route.Path+"/", s.createItemHandler(route))
	}

	address := fmt.Sprintf(":%d", s.Port)

	err := http.ListenAndServe(address, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func (s Server) createCollectionHandler(r Route) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		data, err := s.getCollectionData(r, writer, req)

		if err != nil {
			err = s.writeError(writer, Error{
				Status:  500,
				Message: fmt.Sprintf("Internal server error %v", err),
			})
		} else if data != nil {
			_, err = writer.Write(data)
		}
		if err != nil {
			log.Printf("error writing data: %v", err)
		}
	}
}

func (s Server) createItemHandler(r Route) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		writer.Header().Set("Content-Type", "application/json")

		data, err := s.getItemData(r, req)

		if err != nil {
			_ = s.writeError(writer, Error{
				Status:  500,
				Message: fmt.Sprintf("Internal server error %v", err),
			})
		} else {
			_, _ = writer.Write(data)
		}
	}
}

func (s Server) writeError(writer http.ResponseWriter, e Error) error {
	writer.WriteHeader(e.Status)

	data, err := json.Marshal(e)
	if err != nil {
		return err
	}

	_, err = writer.Write(data)
	return err
}

func (s Server) getCollectionData(r Route, writer http.ResponseWriter, req *http.Request) ([]byte, error) {

	if s.Pagination.Enabled {
		p, err := getPageRequest(s.Pagination, req)
		if err != nil {
			return json.Marshal(Error{
				Status: 400,
				Message: err.Error(),
			})
		}
		page, err := r.getPage(p)
		if err != nil {
			return nil, err
		}
		if s.Pagination.ResponseParametersLocation == LocationHeader {
			page.writeHeaders(writer.Header())
			return json.Marshal(page.Content)
		} else {
			return json.Marshal(page)
		}
	} else {
		return r.getRawData()
	}
}

func (s Server) getItemData(r Route, req *http.Request) ([]byte, error) {
	id := req.RequestURI[len(r.Path)+1:]
	item, err := r.getDataItem(id)
	if err == nil && item != nil {
		return json.Marshal(item)
	}
	return json.Marshal(Error{
		Status:  404,
		Message: fmt.Sprintf("Object with id %s not found", id),
	})
}

func logClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("close error: %v", err)
	}
}
