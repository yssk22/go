package xbytes

import (
	"fmt"

	"encoding/json"
)

func Example_ByteString() {
	byt := []byte("日本語")
	str := ByteString("日本語")
	bytJSON, _ := json.Marshal(byt)
	strJSON, _ := json.Marshal(&str)
	fmt.Println(string(bytJSON), string(strJSON))
	// Output:
	// "5pel5pys6Kqe" "日本語"
}
