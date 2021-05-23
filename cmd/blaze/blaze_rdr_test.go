package main

import (
	"strings"
	"testing"

	"github.com/cirss/go-cli/pkg/util"
)

func TestBlazegraphCmd_rdr_query_bobs_age(t *testing.T) {

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
		<<:bob foaf:age 23>> dct:creator <http://example.com/crawlers#c1> ;
							 dct:source <http://example.net/homepage-listing.html> .
	`)

	Program.AssertExitCode(t, "blaze import --format ttlx --dataset rdr", 0)

	assert_query_result := func(query, result string) {
		outputBuffer.Reset()
		Program.InReader = strings.NewReader(query)
		Program.AssertExitCode(t, "blaze query --format table --dataset rdr", 0)
		util.LineContentsEqual(t, outputBuffer.String(), result)
	}

	t.Run("just_bobs_age", func(t *testing.T) {
		assert_query_result(`
			PREFIX bigdata: <http://bigdata.com/>
			PREFIX foaf: <http://xmlns.com/foaf/0.1/>
			PREFIX dct:  <http://purl.org/dc/elements/1.1/>
			
			SELECT ?age 
			WHERE {
				?bob foaf:name "Bob" .
				?bob foaf:age ?age.
			}
			`,
			`age
        	 ==
        	 23  
		`)
	})

	t.Run("bobs_age_and_its_source", func(t *testing.T) {
		assert_query_result(`
			PREFIX bigdata: <http://bigdata.com/>
			PREFIX foaf: <http://xmlns.com/foaf/0.1/>
			PREFIX dct:  <http://purl.org/dc/elements/1.1/>
			
			SELECT ?age ?src
			WHERE {
				?bob foaf:name "Bob" .
				<<?bob foaf:age ?age>> dct:source ?src .
			}
			`,
			`age | src
	         =============================================
    	     23  | http://example.net/homepage-listing.html
		`)
	})

	t.Run("provenance_of_bobs_age", func(t *testing.T) {
		assert_query_result(`
			PREFIX bigdata: <http://bigdata.com/>
			PREFIX foaf: <http://xmlns.com/foaf/0.1/>
			PREFIX dct:  <http://purl.org/dc/elements/1.1/>
			
			SELECT ?predicate ?src
			WHERE {
				?bob foaf:name "Bob" .
				<<?bob foaf:age 23>> ?predicate ?src .
			}
			`,
			`predicate                               | src
			 =================================================================================
			 http://purl.org/dc/elements/1.1/creator | http://example.com/crawlers#c1
			 http://purl.org/dc/elements/1.1/source  | http://example.net/homepage-listing.html
		`)
	})

	t.Run("everyones_ages_and_their_provenance", func(t *testing.T) {
		assert_query_result(`
			PREFIX bigdata: <http://bigdata.com/>
			PREFIX foaf: <http://xmlns.com/foaf/0.1/>
			PREFIX dct:  <http://purl.org/dc/elements/1.1/>

			SELECT ?name ?age ?predicate ?src 
			WHERE {
				?person foaf:name ?name .
				<<?person foaf:age ?age>> ?predicate ?src .
			}
			`,
			`name | age | predicate                               | src
			 ==============================================================================================
			 Bob  | 23  | http://purl.org/dc/elements/1.1/creator | http://example.com/crawlers#c1
			 Bob  | 23  | http://purl.org/dc/elements/1.1/source  | http://example.net/homepage-listing.html
		`)
	})
}

func TestBlazegraphCmd_rdr_query_alice_bob_and_charlies_ages(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	Program.Invoke("blaze destroy --dataset rdr --quiet")
	Program.Invoke("blaze create --quiet --dataset rdr --rdfstar")

	Program.InReader = strings.NewReader(`
		@prefix : <http://bigdata.com/> .
		@prefix foaf: <http://xmlns.com/foaf/0.1/> .
		@prefix dct:  <http://purl.org/dc/elements/1.1/> .

		:alice   foaf:name  "Alice" .
		:bob     foaf:name  "Bob" .
		:charlie foaf:name  "Charlie" .
		:sam	 foaf:name  "Sam" .
		:joe     foaf:name  "Joe" .

		<<:alice foaf:age 21>> dct:source :sam .
		<<:bob foaf:age 23>> dct:source :sam .
		<<:charlie foaf:age 25>> dct:source :sam .
		<<:charlie foaf:age 27>> dct:source :joe .
	`)

	Program.AssertExitCode(t, "blaze import --format ttlx --dataset rdr", 0)

	assert_query_result := func(query, result string) {
		outputBuffer.Reset()
		Program.InReader = strings.NewReader(query)
		Program.AssertExitCode(t, "blaze query --format table --dataset rdr", 0)
		util.LineContentsEqual(t, outputBuffer.String(), result)
	}

	t.Run("everyones_ages_and_their_provenance", func(t *testing.T) {
		assert_query_result(`
			PREFIX bigdata: <http://bigdata.com/>
			PREFIX foaf: <http://xmlns.com/foaf/0.1/>
			PREFIX dct:  <http://purl.org/dc/elements/1.1/>

			SELECT ?name ?age ?source
			WHERE {
				?person foaf:name ?name .
				<<?person foaf:age ?age>> dct:source ?src .
				?src foaf:name ?source .
			}
			`,
			`name    | age | source
			====================
			Charlie | 27  | Joe
			Alice   | 21  | Sam
			Bob     | 23  | Sam
			Charlie | 25  | Sam
			`)
	})

	t.Run("ages_according_to_sam", func(t *testing.T) {
		assert_query_result(`
			PREFIX bigdata: <http://bigdata.com/>
			PREFIX foaf: <http://xmlns.com/foaf/0.1/>
			PREFIX dct:  <http://purl.org/dc/elements/1.1/>

			SELECT ?name ?age 
			WHERE {
				?person foaf:name ?name .
				<<?person foaf:age ?age>> dct:source ?src .
				?src foaf:name "Sam" .
			}
			`,
			`name    | age
			============
			Alice   | 21
			Bob     | 23
			Charlie | 25
			`)
	})
}
