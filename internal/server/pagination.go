package server

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
)

type (

	PaginationOptions struct {
		Enabled                    bool
		RequestParametersLocation  string // header, query-param
		ResponseParametersLocation string // header, body
		DefaultPageSize            int
	}

	page struct {
		TotalPages       int           `json:"totalPages"`
		TotalElements    int           `json:"totalElements"`
		NumberOfElements int           `json:"numberOfElements"`
		First            bool          `json:"first"`
		Last             bool          `json:"last"`
		Size             int           `json:"size"`
		Number           int           `json:"number"`
		Content          []interface{} `json:"content"`
	}

	pageRequest struct {
		page int
		size int
	}
)

const LocationHeader = "header"
const LocationBody = "body"
const LocationQueryParam = "query-param"

var DefaultPagination = PaginationOptions{
	Enabled:                    true,
	RequestParametersLocation:  LocationQueryParam,
	ResponseParametersLocation: LocationBody,
	DefaultPageSize:            20,
}

func (p *page) writeHeaders(header http.Header) {
	header.Set("page-total-pages", strconv.Itoa(p.TotalPages))
	header.Set("page-total-elements", strconv.Itoa(p.TotalElements))
	header.Set("page-number-of-elements", strconv.Itoa(p.NumberOfElements))
	header.Set("page-first", strconv.FormatBool(p.First))
	header.Set("page-last", strconv.FormatBool(p.First))
	header.Set("page-size", strconv.Itoa(p.Size))
	header.Set("page-number", strconv.Itoa(p.Number))
}

func paginate(data []interface{}, req pageRequest) page {

	var content []interface{}

	start := req.page * req.size
	if start < len(data) {
		end := start + req.size
		if end > len(data) {
			end = len(data)
		}
		content = data[start:end]
	}

	totalPages := int(math.Ceil(float64(len(data)) / float64(req.size)))
	first := req.page == 0
	last := req.page == (totalPages-1) || totalPages == 0

	return page{
		TotalPages:       totalPages,
		TotalElements:    len(data),
		NumberOfElements: len(content),
		First:            first,
		Last:             last,
		Size:             req.size,
		Number:           req.page,
		Content:          content,
	}
}

func getPageRequest(opts PaginationOptions, req *http.Request) (p pageRequest, err error) {
	var pageStr, sizeStr string
	if opts.RequestParametersLocation == LocationQueryParam {
		queries := req.URL.Query()
		pageStr = queries.Get("page")
		sizeStr = queries.Get("size")
	} else if opts.RequestParametersLocation == LocationHeader {
		pageStr = req.Header.Get("page")
		sizeStr = req.Header.Get("size")
	}
	if pageStr != "" {
		p.page, err = strconv.Atoi(pageStr)
		if err != nil || p.page < 0 {
			return p, fmt.Errorf("invalid 'page' parameter [%s]", pageStr)
		}
	} else {
		p.page = 0
	}
	if sizeStr != "" || p.size < 0 {
		p.size, err = strconv.Atoi(sizeStr)
		if err != nil {
			return p, fmt.Errorf("invalid 'size' parameter [%s]", sizeStr)
		}
	} else {
		p.size = opts.DefaultPageSize
	}
	return p, nil
}

