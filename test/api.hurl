GET http://127.0.0.1:8080/

POST http://127.0.0.1:8080/json
Content-Type: application/json
{
  "id": 1,
  "name": "linmingjie",
  "age": 18
}

GET http://127.0.0.1:8080/json
HTTP 405
[Asserts]
body contains "only POST is supported"

POST http://127.0.0.1:8080/json
Content-Type: application/json
`not a json`
HTTP 400
