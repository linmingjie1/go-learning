GET http://127.0.0.1:8080/
Authorization: Bearer test


POST http://127.0.0.1:8080/json
Content-Type: application/json
Authorization: Bearer test
{
  "id": 1,
  "name": "linmingjie",
  "age": 18
}

GET http://127.0.0.1:8080/json
Authorization: Bearer test
HTTP 405
[Asserts]
body contains "only POST is supported"

POST http://127.0.0.1:8080/json
Content-Type: application/json
`not a json`
HTTP 400

GET http://localhost:8080/timeout
Authorization: Bearer test

GET https://cn.bing.com