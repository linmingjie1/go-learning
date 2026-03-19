package demo

import (
	"encoding/json"
	"fmt"
)

func RunJson() {
	jsonStr := `{"name": "Alice", "age": 18}`

	// 定义一个结构体来匹配 JSON 数据的结构
	type Person struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	person := Person{}
	json.Unmarshal([]byte(jsonStr), &person)
	fmt.Println(person)

	bytes, _ := json.Marshal(person)
	fmt.Println(string(bytes))

	bytes2, _ := json.MarshalIndent(person, "", "  ")
	fmt.Println(string(bytes2))
}
