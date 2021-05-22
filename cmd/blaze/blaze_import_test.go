package main

import (
	"strings"
	"testing"

	"github.com/cirss/go-cli/pkg/util"
)

func TestBlazegraphCmd_import_two_triples(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	t.Run("import_nt", func(t *testing.T) {

		Program.Invoke("blaze destroy --dataset kb --quiet")
		Program.Invoke("blaze create --quiet --dataset kb")

		Program.InReader = strings.NewReader(`
			<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
			<http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
		`)

		Program.AssertExitCode(t, "blaze import --format nt", 0)

		outputBuffer.Reset()
		Program.Invoke("blaze export --format nt --sort=true")
		util.LineContentsEqual(t, outputBuffer.String(),
			`<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
			 <http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
			`)
	})

	t.Run("import_ttl", func(t *testing.T) {

		Program.Invoke("blaze destroy --dataset kb --quiet")
		Program.Invoke("blaze create --quiet --dataset kb")

		Program.InReader = strings.NewReader(
			`@prefix data: <http://tmcphill.net/data#> .
			 @prefix tags: <http://tmcphill.net/tags#> .

			 data:y tags:tag "eight" .
			 data:x tags:tag "seven" .
			`)

		Program.AssertExitCode(t, "blaze import --format ttl", 0)

		outputBuffer.Reset()
		Program.Invoke("blaze export --format nt --sort=true")
		util.LineContentsEqual(t, outputBuffer.String(),
			`<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
			 <http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
			`)
	})

	t.Run("import_jsonld", func(t *testing.T) {

		Program.Invoke("blaze destroy --dataset kb --quiet")
		Program.Invoke("blaze create --quiet --dataset kb")

		Program.InReader = strings.NewReader(
			`
			[
				{
					"@id": "http://tmcphill.net/data#x",
					"http://tmcphill.net/tags#tag": "seven"
				},
				{
					"@id": "http://tmcphill.net/data#y",
					"http://tmcphill.net/tags#tag": "eight"
				}
			]
			`)

		Program.AssertExitCode(t, "blaze import --format jsonld", 0)

		outputBuffer.Reset()
		Program.Invoke("blaze export --format nt --sort=true")
		util.LineContentsEqual(t, outputBuffer.String(),
			`<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven"^^<http://www.w3.org/2001/XMLSchema#string> .
			 <http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight"^^<http://www.w3.org/2001/XMLSchema#string> .
			`)
	})

	t.Run("import_ttl", func(t *testing.T) {

		Program.Invoke("blaze destroy --dataset kb --quiet")
		Program.Invoke("blaze create --quiet --dataset kb")

		Program.InReader = strings.NewReader(
			`@prefix data: <http://tmcphill.net/data#> .
			 @prefix tags: <http://tmcphill.net/tags#> .

			 data:y tags:tag "eight" .
			 data:x tags:tag "seven" .
			`)

		Program.AssertExitCode(t, "blaze import --format ttl", 0)

		outputBuffer.Reset()
		Program.Invoke("blaze export --format nt --sort=true")
		util.LineContentsEqual(t, outputBuffer.String(),
			`<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
			 <http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
			`)
	})

	t.Run("import_xml", func(t *testing.T) {

		Program.Invoke("blaze destroy --dataset kb --quiet")
		Program.Invoke("blaze create --quiet --dataset kb")

		Program.InReader = strings.NewReader(
			`<rdf:RDF xmlns:rdf="http://www.w3.org/1999/02/22-rdf-syntax-ns#">

			 <rdf:Description rdf:about="http://tmcphill.net/data#y">
				<tag xmlns="http://tmcphill.net/tags#">eight</tag>
		 	 </rdf:Description>

 			 <rdf:Description rdf:about="http://tmcphill.net/data#x">
				<tag xmlns="http://tmcphill.net/tags#">seven</tag>
  			 </rdf:Description>

			 </rdf:RDF>
			`)

		Program.AssertExitCode(t, "blaze import --format xml", 0)

		outputBuffer.Reset()
		Program.Invoke("blaze export --format nt --sort=true")
		util.LineContentsEqual(t, outputBuffer.String(),
			`<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
			 <http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
			`)
	})
}

func TestBlazegraphCmd_import_rdfstar_two_rdf_triples(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset kb --quiet")
	Program.Invoke("blaze create --quiet --dataset kb --rdfstar")

	Program.InReader = strings.NewReader(`
		<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
		<http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
	`)

	Program.AssertExitCode(t, "blaze import --format nt", 0)

	outputBuffer.Reset()
	Program.Invoke("blaze export --format nt --sort=true")
	util.LineContentsEqual(t, outputBuffer.String(),
		`<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
			<http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
		`)
}

//http://127.0.0.1:9999/blazegraph/namespace/kb/sparql

// func TestBlazegraphCmd_import_rdfstar_import_reified_triple(t *testing.T) {

// 	var outputBuffer strings.Builder
// 	Program.OutWriter = &outputBuffer
// 	Program.ErrWriter = &outputBuffer

// 	Program.Invoke("blaze destroy --dataset kb --quiet")
// 	Program.Invoke("blaze create --quiet --dataset kb --rdfstar")

// 	Program.InReader = strings.NewReader(`
// 		@prefix : <http://bigdata.com/> .
// 		@prefix foaf: <http://xmlns.com/foaf/0.1/> .
// 		@prefix dct:  <http://purl.org/dc/elements/1.1/> .
// 		:bob foaf:name "Bob" .
// 		<<:bob foaf:age 23>> dct:creator <http://example.com/crawlers#c1> .

// 	`)

// 	Program.AssertExitCode(t, "blaze import --format ttlx", 0)

// 	outputBuffer.Reset()
// 	Program.AssertExitCode(t, "blaze export --format json", 0)
// 	util.LineContentsEqual(t, outputBuffer.String(),
// 		`<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
// 			<http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
// 		`)
// }

func TestBlazegraphCmd_import_specific_dataset(t *testing.T) {

	triples_ttl :=
		`<http://tmcphill.net/data#x> <http://tmcphill.net/tags#tag> "seven" .
		 <http://tmcphill.net/data#y> <http://tmcphill.net/tags#tag> "eight" .
		`
	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	// t.Run("default", func(t *testing.T) {
	// 	outputBuffer.Reset()
	// 	Program.Invoke("blaze destroy --silent")
	// 	Program.Invoke("blaze create --quiet")
	// 	Program.InReader = strings.NewReader(triples_ttl)
	// 	Program.AssertExitCode(t, "blaze import", 0)
	// 	Program.Invoke("blaze export --sort=true")
	// 	util.LineContentsEqual(t, outputBuffer.String(), triples_ttl)
	// })

	t.Run("single_custom", func(t *testing.T) {
		outputBuffer.Reset()
		Program.Invoke("blaze destroy --dataset foo --silent")
		Program.Invoke("blaze create --dataset foo --quiet")
		Program.InReader = strings.NewReader(triples_ttl)
		Program.AssertExitCode(t, "blaze import --dataset foo", 0)
		Program.Invoke("blaze export --dataset foo --sort=true")
		util.LineContentsEqual(t, outputBuffer.String(), triples_ttl)
	})

}

var expectedImportHelpOutput = string(
	`blaze import: Imports triples in the specified format into an RDF dataset.

	usage: blaze import [<flags>]

	flags:
		-dataset name
				name of RDF dataset to import triples into (default "kb")
		-file string
				File containing triples to import (default "-")
		-format string
				Format of triples to import [jsonld, nt, ttl, ttlx, or xml] (default "ttl")
		-instance URL
				URL of Blazegraph instance (default "http://127.0.0.1:9999/blazegraph")
		-quiet
				Discard normal command output
		-silent
				Discard normal and error command output
	`)

func TestBlazegraphCmd_import_help(t *testing.T) {
	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer
	Program.AssertExitCode(t, "blaze import help", 0)
	util.LineContentsEqual(t, outputBuffer.String(), expectedImportHelpOutput)
}

func TestBlazegraphCmd_help_import(t *testing.T) {
	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer
	Program.AssertExitCode(t, "blaze help import", 0)
	util.LineContentsEqual(t, outputBuffer.String(), expectedImportHelpOutput)
}

func TestBlazegraphCmd_import_bad_flag(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.AssertExitCode(t, "blaze import --not-a-flag", 1)

	util.LineContentsEqual(t, outputBuffer.String(),
		`blaze import: flag provided but not defined: -not-a-flag

		usage: blaze import [<flags>]

		flags:
			-dataset name
					name of RDF dataset to import triples into (default "kb")
			-file string
					File containing triples to import (default "-")
			-format string
					Format of triples to import [jsonld, nt, ttl, ttlx, or xml] (default "ttl")
			-instance URL
					URL of Blazegraph instance (default "http://127.0.0.1:9999/blazegraph")
			-quiet
					Discard normal command output
			-silent
					Discard normal and error command output
	`)
}
