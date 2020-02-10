# go-json-server
Simple web application to serve json files in a RESTful manner. Can POST, PUT, PATCH, and GET.

Arguments:
 * `--port PORT_NUMBER`: which port to serve on
 * `--route PATH FILE`: add a route under the PATH serving the json in FILE
 * `--raw-route PATH FILE`: add a route under the PATH serving the json in FILE directly
 * `--base-path PATH`: prepend every route with PATH
 * `--paginated`: paginate responses (default false)
 * `--page-one-indexed`: pages start at 1 (default 0)
 * `--page-request-location LOCATION`: where to find page params 'page' and 'size', either 'query-param' or 'header'
 * `--page-response-location LOCATION`: where to send page attributes, either 'body' or 'header'
 * `--default-page-size SIZE`: default pagination size (default to 20)
 * `--debug/info/error`: set log level
 * `--fake-load`: number of seconds to wait before returning each call


Example:
```
    go-json-server --base-path /api \
                   --paginated \
                   --raw-route /auth test_data/auth.json \
                   --route /people test_data/people.json \
                   --route /animals test_data/animals.json
```

Additionally, you can add the `sleep` query parameter with a duration to any call to add a fake load, e.g. `sleep=5s`.
