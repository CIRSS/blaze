package main

import (
	"strings"
	"testing"

	"github.com/cirss/go-cli/pkg/util"
)

func TestBlazegraphCmd_query_rdfstar(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset rdr --quiet")
	Program.Invoke("blaze create --quiet --dataset rdr --rdfstar")

	Program.InReader = strings.NewReader(`
		@prefix : <http://bigdata.com/> .
		@prefix foaf: <http://xmlns.com/foaf/0.1/> .
		@prefix dct:  <http://purl.org/dc/elements/1.1/> .
		:bob foaf:name "Bob" .
		<<:bob foaf:age 23>> dct:creator <http://example.com/crawlers#c2> .
		<<:bob foaf:age 23>> dct:source <http://example.com/crawlers#c2> .
	`)

	Program.AssertExitCode(t, "blaze import --format ttlx --dataset rdr", 0)

	query := `
		PREFIX bigdata: <http://bigdata.com/>
		PREFIX foaf: <http://xmlns.com/foaf/0.1/>
		PREFIX dct:  <http://purl.org/dc/elements/1.1/>
		
		SELECT ?src ?s ?age ?p ?pp WHERE {
		?bob foaf:name "Bob" .
		<<?s ?p ?age>> ?pp ?src .
		}
	`

	outputBuffer.Reset()
	Program.InReader = strings.NewReader(query)
	Program.AssertExitCode(t, "blaze query --format table --dataset rdr", 0)

	util.LineContentsEqual(t, outputBuffer.String(),
		`src                            | s                      | age | p                             | pp
        =========================================================================================================================================
        http://example.com/crawlers#c2 | http://bigdata.com/bob | 23  | http://xmlns.com/foaf/0.1/age | http://purl.org/dc/elements/1.1/creator
        http://example.com/crawlers#c2 | http://bigdata.com/bob | 23  | http://xmlns.com/foaf/0.1/age | http://purl.org/dc/elements/1.1/source
		`)
}
