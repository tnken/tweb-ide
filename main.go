package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/takeru56/tcompiler/compiler"
	"github.com/takeru56/tcompiler/parser"
	"github.com/takeru56/tcompiler/token"
)

type SourceJSON struct {
	SourceCode string `json:"source_code"`
}

type BytecodeJSON struct {
	ByteCode string `json:"byte_code"`
}

func main() {
	http.Handle("/", http.FileServer(http.Dir("public")))
	http.HandleFunc("/compile/", compileHandler)
	http.ListenAndServe(":3000", nil)
}

func compileHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		// body := r.Body
		// defer body.Close()

		// buf := new(bytes.Buffer)
		// io.Copy(buf, body)

		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			fmt.Println(err.Error())
			return
		}

		var sj SourceJSON
		json.Unmarshal(body, &sj)

		w.WriteHeader(http.StatusCreated)
		source := sj.SourceCode
		fmt.Println(source)

		tok := token.New(source)
		parser, err := parser.New(tok)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		p, err := parser.Program()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		c := compiler.Exec(p)
		outjson, err := json.Marshal(BytecodeJSON{c.Bytecode()})
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(source)
		fmt.Fprint(w, string(outjson))
	default:
		fmt.Fprint(w, "Method not allowed\n")
	}
}
