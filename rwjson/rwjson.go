package main

import (
	"encoding/json"
	"io"
	"log"
	"os"

	"github.com/pdk/rozer"
)

func main() {

	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		log.Fatal(err)
	}

	// input := []byte(`{"z":"apple", "x": {"b": "ack", "a": [1,2,3]}, "a":42}`)

	r := rozer.New()
	err = json.Unmarshal(input, r)
	if err != nil {
		log.Fatal(err)
	}

	// log.Printf("r = %#v", r)

	output, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	os.Stdout.Write(output)
	os.Stdout.WriteString("\n")
}
