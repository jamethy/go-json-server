package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/jamethy/go-json-server/internal/server"
)

func main() {
	for _, arg := range os.Args {
		if arg == "-h" || arg == "help" || arg == "--help" {
			fmt.Printf(`
go-json-server v1.3

Serves json files in a RESTful manner.
Arguments:
  --port PORT_NUMBER: which port to serve on
  --route PATH FILE: add a route under the PATH serving the json in FILE
  --raw-route PATH FILE: add a route under the PATH serving the json in FILE directly
  --base-path PATH: prepend every route with PATH
  --paginated: paginate responses (default false)
  --page-one-indexed: pages start at 1 (default 0)
  --page-request-location LOCATION: where to find page params 'page' and 'size', either 'query-param' or 'header'
  --page-response-location LOCATION: where to send page attributes, either 'body' or 'header'
  --default-page-size SIZE: default pagination size (default to 20)
  --debug/info/error: set log level
  --fake-load: number of seconds to wait before returning each call
`)
		}
	}

	s := parseArgs(os.Args)
	s.Start()
}

func parseArgs(args []string) server.Server {
	opts := server.Server{
		Port:       8080,
		Routes:     nil,
		Pagination: server.DefaultPagination,
		FakeLoad:   0,
	}

	if args == nil {
		return opts
	}

	for i := 1; i < len(args); {

		switch os.Args[i] {
		case "--port":
			assertArgCount(args, i, 1)
			opts.Port = mustParseUnsignedInt(args[i+1])
			i += 2
		case "--route":
			assertArgCount(args, i, 2)
			opts.Routes = append(opts.Routes, server.Route{
				Path:     args[i+1],
				JsonFile: args[i+2],
				IdField:  "id",
			})
			i += 3
		case "--raw-route":
			assertArgCount(args, i, 2)
			opts.Routes = append(opts.Routes, server.Route{
				Path:     args[i+1],
				JsonFile: args[i+2],
				Raw:      true,
			})
			i += 3
		case "--base-path":
			opts.BasePath = args[i+1]
			i += 2
		case "--paginated":
			opts.Pagination.Enabled = true
			i += 1
		case "--page-one-indexed":
			opts.Pagination.ZeroIndexed = false
			i += 1
		case "--page-request-location":
			assertArgCount(args, i, 1)
			opts.Pagination.RequestParametersLocation = args[i+1]
			i += 2
		case "--page-response-location":
			assertArgCount(args, i, 1)
			opts.Pagination.ResponseParametersLocation = args[i+1]
			i += 2
		case "--default-page-size":
			assertArgCount(args, i, 1)
			opts.Pagination.DefaultPageSize = mustParseUnsignedInt(args[i+1])
			i += 2
		case "--debug":
			server.SetLogLevel(server.LogLevelDebug)
			i += 1
		case "--info":
			server.SetLogLevel(server.LogLevelInfo)
			i += 1
		case "--error":
			server.SetLogLevel(server.LogLevelError)
			i += 1
		case "--fake-load":
			assertArgCount(args, i, 1)
			opts.FakeLoad, _ = time.ParseDuration(args[i+1])
			i += 2
		default:
			log.Fatal("unrecognized argument " + args[i])
		}
	}

	return opts
}

func assertArgCount(args []string, index int, count int) {
	message := fmt.Errorf("%s must be followed by %d value(s)", args[index], count)
	if index+count >= len(args) {
		log.Fatal(message)
	}
	for i := index + 1; i < index+count; i++ {
		if strings.HasPrefix(args[i], "-") {
			log.Fatal(message)
		}
	}
}

func mustParseUnsignedInt(str string) int {
	parsedInt, err := strconv.ParseInt(str, 10, 64)
	if err != nil || parsedInt <= 0 {
		log.Fatal("invalid integer provided")
	}
	return int(parsedInt)
}
