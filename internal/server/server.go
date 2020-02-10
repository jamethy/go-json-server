package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type (
	Server struct {
		Pagination PaginationOptions
		Routes     []Route
		BasePath   string
		Port       int
		FakeLoad   time.Duration
	}

	Error struct {
		Status  int    `json:"status"`
		Message string `json:"message"`
	}
)

func (s Server) Start() {
	mux := http.NewServeMux()

	for _, route := range s.Routes {

		if route.Raw {
			// get raw json data
			mux.HandleFunc(s.BasePath+route.Path, s.createRawHandler(route))

		} else {
			// get all of route
			mux.HandleFunc(s.BasePath+route.Path, s.createCollectionHandler(route))

			// get specific id
			mux.HandleFunc(s.BasePath+route.Path+"/", s.createItemHandler(route))
		}
	}

	address := fmt.Sprintf(":%d", s.Port)

	err := http.ListenAndServe(address, mux)
	if err != nil {
		log.Fatal(err)
	}
}

func (s Server) createRawHandler(r Route) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		if req.Method != http.MethodOptions {
			s.fakeLoad(req)
			writer.Header().Set("Content-Type", "application/json")

			data, err := r.getRawData()
			if err != nil {
				err = s.writeError(writer, Error{
					Status:  500,
					Message: fmt.Sprintf("Internal server error %v", err),
				})
			} else if data != nil {
				_, err = writer.Write(data)
			}
			if err != nil {
				logError(fmt.Sprintf("error writing data: %v", err))
			}
		}
	}
}

func (s Server) createCollectionHandler(r Route) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		s.fakeLoad(req)
		writer.Header().Set("Content-Type", "application/json")

		var err error
		if req.Method == http.MethodOptions {

		} else if req.Method == http.MethodGet {
			err = s.handleCollectionGet(r, writer, req)

		} else if req.Method == http.MethodPost {
			err = s.handleCollectionUpdate(writer, req, r.addDataItem)

		} else if req.Method == http.MethodPut {
			err = s.handleCollectionUpdate(writer, req, r.putDataItem)

		} else if req.Method == http.MethodPatch {
			err = s.handleCollectionUpdate(writer, req, r.patchDataItem)

		} else {
			err = s.writeError(writer, Error{
				Status:  405,
				Message: fmt.Sprintf("Method Not Allowed"),
			})
		}

		if err != nil {
			err = s.writeError(writer, Error{
				Status:  500,
				Message: fmt.Sprintf("Internal server error %v", err),
			})
			if err != nil {
				logError(fmt.Sprintf("error writing data: %v", err))
			}
		}
	}
}

func (s Server) createItemHandler(r Route) http.HandlerFunc {
	return func(writer http.ResponseWriter, req *http.Request) {
		var writingError error
		if req.Method == http.MethodOptions {
			return
		}

		if req.Method != http.MethodGet {
			writingError = s.writeError(writer, Error{
				Status:  405,
				Message: fmt.Sprintf("Method Not Allowed"),
			})
		} else {
			s.fakeLoad(req)
			writer.Header().Set("Content-Type", "application/json")

			data, err := s.getItemData(r, req)

			if err != nil {
				writingError = s.writeError(writer, Error{
					Status:  500,
					Message: fmt.Sprintf("Internal server error %v", err),
				})
			} else {
				_, writingError = writer.Write(data)
			}
		}
		if writingError != nil {
			logError(fmt.Sprintf("error writing data: %v", writingError))
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
				Status:  400,
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

func (s Server) handleCollectionGet(r Route, writer http.ResponseWriter, req *http.Request) error {
	data, err := s.getCollectionData(r, writer, req)
	if err != nil {
		return err
	} else if data != nil {
		_, err = writer.Write(data)
	}
	return err
}

type updaterFunc func(item map[string]interface{}) (map[string]interface{}, error)

func (s Server) handleCollectionUpdate(writer http.ResponseWriter, req *http.Request, updater updaterFunc) error {
	item, err := getWholeBodyAsObject(req)
	if err != nil {
		return err
	}
	item, err = updater(item)
	if err != nil {
		return err
	}

	b, _ := json.Marshal(item)
	_, err = writer.Write(b)
	return err
}

func logClose(c io.Closer) {
	if err := c.Close(); err != nil {
		log.Printf("close error: %v", err)
	}
}

func (s Server) fakeLoad(req *http.Request) {
	fakeLoad := s.FakeLoad

	sleep := req.URL.Query().Get("sleep")
	if sleep != "" {
		fakeLoad, _ = time.ParseDuration(sleep)
	}

	if fakeLoad > 0 {
		logDebug(fmt.Sprintf("Sleeping for %v seconds for fake load", fakeLoad))
		time.Sleep(fakeLoad)
	}
}

func getWholeBodyAsObject(req *http.Request) (map[string]interface{}, error) {
	b, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}
	logClose(req.Body)

	item := make(map[string]interface{})
	err = json.Unmarshal(b, &item)
	return item, err
}
