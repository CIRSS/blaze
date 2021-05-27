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

	t.Run("bobs_age_and_its_source_using_bind", func(t *testing.T) {
		assert_query_result(`
			PREFIX bigdata: <http://bigdata.com/>
			PREFIX foaf: <http://xmlns.com/foaf/0.1/>
			PREFIX dct:  <http://purl.org/dc/elements/1.1/>
			
			SELECT ?age ?src
			WHERE {
				?bob foaf:name "Bob" .
				BIND( <<?bob foaf:age ?age>> AS ?t )
				?t dct:source ?src .
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

	Program.Invoke("blaze destroy --dataset rdr --quiet ")
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

	t.Run("bobs_age_and_its_source", func(t *testing.T) {
		assert_query_result(`
			PREFIX bigdata: <http://bigdata.com/>
			PREFIX foaf: <http://xmlns.com/foaf/0.1/>
			PREFIX dct:  <http://purl.org/dc/elements/1.1/>
			
			SELECT ?age ?source
			WHERE {
				?bob foaf:name "Bob" .
				<<?bob foaf:age ?age>> dct:source ?src .
				?src foaf:name ?source .
			}
			`,
			`age | source
			 ============
			 23  | Sam			 
			`)
	})

	t.Run("charlies_ages_and_their_sources", func(t *testing.T) {
		assert_query_result(`
			PREFIX bigdata: <http://bigdata.com/>
			PREFIX foaf: <http://xmlns.com/foaf/0.1/>
			PREFIX dct:  <http://purl.org/dc/elements/1.1/>
			
			SELECT ?age ?source
			WHERE {
				?charlie foaf:name "Charlie" .
				<<?charlie foaf:age ?age>> dct:source ?src .
				?src foaf:name ?source .
			}
			`,
			`age | source
			 ============
			 27  | Joe			 
			 25  | Sam
			`)
	})

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

	t.Run("everyones_ages_and_their_provenance_using_bind", func(t *testing.T) {
		assert_query_result(`
			PREFIX bigdata: <http://bigdata.com/>
			PREFIX foaf: <http://xmlns.com/foaf/0.1/>
			PREFIX dct:  <http://purl.org/dc/elements/1.1/>

			SELECT ?name ?age ?source
			WHERE {
				?person foaf:name ?name .
				BIND ( <<?person foaf:age ?age>> AS ?t )
				?t dct:source ?src .
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

	t.Run("everything_sam_is_the_source_of", func(t *testing.T) {
		assert_query_result(`
			PREFIX bigdata: <http://bigdata.com/>
			PREFIX foaf: <http://xmlns.com/foaf/0.1/>
			PREFIX dct:  <http://purl.org/dc/elements/1.1/>

			SELECT ?s ?p ?o 
			WHERE {
				?src foaf:name "Sam" .
				<<?s ?p ?o>> dct:source ?src .
			}
			`,
			`s                          | p                             | o
			 ================================================================
			 http://bigdata.com/alice   | http://xmlns.com/foaf/0.1/age | 21
			 http://bigdata.com/bob     | http://xmlns.com/foaf/0.1/age | 23
			 http://bigdata.com/charlie | http://xmlns.com/foaf/0.1/age | 25
			`)
	})

	t.Run("every_triple_that_has_a_triple_for_its_subject", func(t *testing.T) {
		assert_query_result(`
			PREFIX bigdata: <http://bigdata.com/>
			PREFIX foaf: <http://xmlns.com/foaf/0.1/>
			PREFIX dct:  <http://purl.org/dc/elements/1.1/>

			SELECT ?s ?p ?o ?pp ?oo
			WHERE {
				<<?s ?p ?o>> ?pp ?oo .
			}
			`,
			`s                          | p                             | o  | pp                                     | oo
			 ====================================================================================================================================
			 http://bigdata.com/alice   | http://xmlns.com/foaf/0.1/age | 21 | http://purl.org/dc/elements/1.1/source | http://bigdata.com/sam
			 http://bigdata.com/bob     | http://xmlns.com/foaf/0.1/age | 23 | http://purl.org/dc/elements/1.1/source | http://bigdata.com/sam
			 http://bigdata.com/charlie | http://xmlns.com/foaf/0.1/age | 25 | http://purl.org/dc/elements/1.1/source | http://bigdata.com/sam
			 http://bigdata.com/charlie | http://xmlns.com/foaf/0.1/age | 27 | http://purl.org/dc/elements/1.1/source | http://bigdata.com/joe
			`)
	})
}

func TestBlazegraphCmd_rdr_query_two_levels_of_reification(t *testing.T) {

	var outputBuffer strings.Builder
	Program.OutWriter = &outputBuffer
	Program.ErrWriter = &outputBuffer

	insert_data := func(data string) {
		Program.InReader = strings.NewReader(data)
		Program.AssertExitCode(t, "blaze import --format ttlx --dataset rdr", 0)
	}

	assert_query_result := func(query, result string) {
		outputBuffer.Reset()
		Program.InReader = strings.NewReader(query)
		Program.AssertExitCode(t, "blaze query --format table --dataset rdr", 0)
		util.LineContentsEqual(t, outputBuffer.String(), result)
	}

	Program.Invoke("blaze destroy --dataset rdr --quiet ")
	Program.Invoke("blaze create --quiet --dataset rdr --rdfstar")

	insert_data(`
		@prefix x: <http://example/> .
		<x:a> <x:b> <x:c> .
	`)

	assert_query_result(`
		SELECT ?s ?p ?o
		WHERE {
			?s ?p ?o .
		}`,
		`s   | p   | o
         ================
         x:a | x:b | x:c
		`)

	assert_query_result(`
		SELECT ?s ?p ?o
		WHERE {
			?s ?p ?o .
			FILTER NOT EXISTS {
				<<?s ?p ?o>> ?pp ?oo .
			}
		}
		`,
		`s   | p   | o
         ================
         x:a | x:b | x:c
		`)

	insert_data(`
		@prefix x: <http://example/> .
		<<<x:a> <x:b> <x:c>>> <x:d> <x:e> .
	`)

	assert_query_result(`
		PREFIX x: <http://example/>
		SELECT ?s ?p ?o
			WHERE {
				<<?s ?p ?o>> <x:d> <x:e> .
			}
			`,
		`s   | p   | o
		 ================
		 x:a | x:b | x:c
		`)

	assert_query_result(`
		SELECT ?s ?p ?o ?pp ?oo
		WHERE {
			<<?s ?p ?o>> ?pp ?oo .
		}
		`,
		`s   | p   | o   | pp  | oo
		 ==============================
		 x:a | x:b | x:c | x:d | x:e
		`)

	insert_data(`
		@prefix x: <http://example/> .
		<<<<<x:a> <x:b> <x:c>>> <x:d> <x:e>>> <x:f> <x:g> .
	`)

	assert_query_result(`
		SELECT ?s ?p ?o ?pp ?oo ?ppp ?ooo
		WHERE {
			<<<<?s ?p ?o>> ?pp ?oo>> ?ppp ?ooo .
		}
		`,
		`s   | p   | o   | pp  | oo  | ppp | ooo
		 ============================================
		 x:a | x:b | x:c | x:d | x:e | x:f | x:g
		`)

	assert_query_result(`
		PREFIX x: <http://example/>
		SELECT ?s ?p ?o
			WHERE {
				<<?s ?p ?o>> <x:d> <x:e> .
			}
			`,
		`s   | p   | o
		 ================
		 x:a | x:b | x:c
		`)
}
