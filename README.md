# go-json-server
Simple web application to serve json files in a RESTful manner

Arguments:
 * `--port PORT_NUMBER`: which port to serve on
 * `--route PATH FILE`: add a route under the PATH serving the json in FILE
 * `--base-path PATH`: prepend every route with PATH
 * `--paginated`: paginate responses (default false)
 * `--page-request-location LOCATION`: where to find page params 'page' and 'size', either 'query-param' or 'header'
 * `--page-response-location LOCATION`: where to send page attributes, either 'body' or 'header'
 * `--default-page-size SIZE`: default pagination size (default to 20)

Example:
```
    go-json-server --base-path /api \
                   --paginated \
                   --route /people test_data/people.json \
                   --route /animals test_data/animals.json
```
