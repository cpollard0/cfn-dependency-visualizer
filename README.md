Visualizes relationship between cloudformation imports and exports.

Get some AWS creds in the default location (~/.aws/credentails)

Run get-all-imports-exports.go by eecuting
`go run get-all-imports-exports.go`

That will generate JSON outoutput file

Once processing is done, navigate to http://localhost:8080/visualization.html to view the data.

TODO Items:

1. Add ability to group by tags (color coded circles by tag)
2. Generally make easier to navigate (any feedback??)
3. Use channels in go to fetch data concurrently where possible
