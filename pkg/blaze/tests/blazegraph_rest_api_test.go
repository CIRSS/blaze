package tests

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"testing"

	"github.com/cirss/blaze/pkg/blaze"
	"github.com/cirss/go-cli/pkg/util"
)

var sparqlEndpoint = blaze.DefaultUrl + "/sparql"

func ExampleBlazegraph_EmptyRequest_HttpGet_StatusOK() {
	response, _ := http.Get(sparqlEndpoint)
	fmt.Println(response.Status)
	// Output:
	// 200 OK
}

func ExampleBlazegraph_EmptyRequest_HttpClientGet_StatusOK() {
	client := &http.Client{}
	response, _ := client.Get(sparqlEndpoint)
	fmt.Println(response.Status)
	// Output:
	// 200 OK
}

func ExampleBlazegraph_EmptyRequest_HttpClientDoGet_StatusOK() {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", sparqlEndpoint, nil)
	response, _ := client.Do(request)
	fmt.Println(response.Status)
	// Output:
	// 200 OK
}

func ExampleBlazegraph_EmptyRequest_HttpClientDoGet_ResponseContentType_RdfXml() {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", sparqlEndpoint, nil)
	response, _ := client.Do(request)
	fmt.Println(response.Header["Content-Type"])
	// Output:
	// [application/rdf+xml]
}

func ExampleBlazegraph_EmptyRequest_HttpClientDoGet_ResponseContentType_JSON() {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", sparqlEndpoint, nil)
	request.Header.Add("Accept", "application/json")
	response, _ := client.Do(request)
	fmt.Println(response.Header["Content-Type"])
	// Output:
	// [application/sparql-results+json]
}

func ExampleBlazegraph_EmptyRequest_Body_First100Bytes() {
	client := &http.Client{}
	request, _ := http.NewRequest("GET", sparqlEndpoint, nil)
	request.Header.Add("Accept", "application/json")
	response, _ := client.Do(request)
	b, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	fmt.Println(string(b[0:100]))
	// Output:
	// {
	//   "head" : {
	//     "vars" : [ "subject", "predicate", "object", "context" ]
	//   },
	//   "results" : {
}

func ExampleBlazegraph_PostSparqlRequest_SelectThreeVariablesWithEmptyWhereClause() {
	deleteDefaultDataset()
	client := &http.Client{}
	request, _ := http.NewRequest("POST", sparqlEndpoint, strings.NewReader(`
		SELECT ?s ?p ?o
		WHERE
		{}
	`))
	request.Header.Add("Content-Type", "application/sparql-query")
	request.Header.Add("Accept", "application/json")
	response, _ := client.Do(request)
	b, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	fmt.Println(string(b))
	// Output:
	// {
	//   "head" : {
	//     "vars" : [ "s", "p", "o" ]
	//   },
	//   "results" : {
	//     "bindings" : [ { } ]
	//   }
	// }
}

func ExampleBlazegraph_PostSparqlRequest_SelectWithWhereClauseMatchingAllTriples_EmptyDataset() {
	deleteDefaultDataset()
	client := &http.Client{}
	request, _ := http.NewRequest("POST", sparqlEndpoint, strings.NewReader(`
		SELECT ?s ?p ?o
		WHERE
		{ ?s ?p ?o }
	`))
	request.Header.Add("Content-Type", "application/sparql-query")
	request.Header.Add("Accept", "application/json")
	response, _ := client.Do(request)
	b, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	fmt.Println(string(b))
	// Output:
	// {
	//   "head" : {
	//     "vars" : [ "s", "p", "o" ]
	//   },
	//   "results" : {
	//     "bindings" : [ ]
	//   }
	// }
}

func TestBlazegraph_PostSparqlRequest_SelectWithWhereClauseMatchingAllTriples_DatasetWithTwoTriples(t *testing.T) {
	deleteDefaultDataset()
	insertTwoTriplesIntoDefaultDataset()
	client := &http.Client{}
	request, _ := http.NewRequest("POST", sparqlEndpoint, strings.NewReader(`SELECT ?s ?p ?o WHERE { ?s ?p ?o }
	`))
	request.Header.Add("Content-Type", "application/sparql-query")
	request.Header.Add("Accept", "application/json")
	response, _ := client.Do(request)
	b, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	util.LineContentsEqual(t, string(b),
		`{
		"head" : {
		  "vars" : [ "s", "p", "o" ]
		},
		"results" : {
		  "bindings" : [ {
			"s" : {
			  "type" : "uri",
			  "value" : "http://learningsparql.com/ns/data#x"
			},
			"p" : {
			  "type" : "uri",
			  "value" : "http://learningsparql.com/ns/addressbook#tag"
			},
			"o" : {
			  "type" : "literal",
			  "value" : "eight"
			}
		  }, {
			"s" : {
			  "type" : "uri",
			  "value" : "http://learningsparql.com/ns/data#y"
			},
			"p" : {
			  "type" : "uri",
			  "value" : "http://learningsparql.com/ns/addressbook#tag"
			},
			"o" : {
			  "type" : "literal",
			  "value" : "seven"
			}
		  } ]
		}
	  }`)
}

func ExampleBlazegraph_PostData_EmptyBody() {
	client := &http.Client{}
	request, _ := http.NewRequest("POST", sparqlEndpoint, strings.NewReader(""))
	request.Header.Add("Content-Type", "application/x-turtle")
	response, _ := client.Do(request)
	b, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	fmt.Println(string(b)[0:54])
	// Output:
	// <?xml version="1.0"?><data modified="0" milliseconds="
}

func ExampleBlazegraph_PostData_NamespaceDeclarationsOnly() {
	client := &http.Client{}
	request, _ := http.NewRequest("POST", sparqlEndpoint, strings.NewReader(`
		@prefix ab:    <http://learningsparql.com/ns/addressbook#> .
		@prefix d:     <http://learningsparql.com/ns/data#> .
	`))
	request.Header.Add("Content-Type", "application/x-turtle")
	response, _ := client.Do(request)
	b, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	fmt.Println(string(b)[0:54])
	// Output:
	// <?xml version="1.0"?><data modified="0" milliseconds="
}

func ExampleBlazegraph_PostData_TwoTriples() {
	deleteDefaultDataset()
	client := &http.Client{}
	request, _ := http.NewRequest("POST", sparqlEndpoint, strings.NewReader(`
		@prefix ab:    <http://learningsparql.com/ns/addressbook#> .
		@prefix d:     <http://learningsparql.com/ns/data#> .

		d:y ab:tag "seven" .
		d:x ab:tag "eight" .
	`))
	request.Header.Add("Content-Type", "application/x-turtle")
	response, _ := client.Do(request)
	b, _ := ioutil.ReadAll(response.Body)
	response.Body.Close()
	fmt.Println(string(b)[0:54])
	// Output:
	// <?xml version="1.0"?><data modified="2" milliseconds="
}

func deleteDefaultDataset() {
	client := &http.Client{}
	request, _ := http.NewRequest("DELETE", sparqlEndpoint, nil)
	client.Do(request)
}

func insertTwoTriplesIntoDefaultDataset() {
	client := &http.Client{}
	request, _ := http.NewRequest("POST", sparqlEndpoint, strings.NewReader(`
		@prefix ab:    <http://learningsparql.com/ns/addressbook#> .
		@prefix d:     <http://learningsparql.com/ns/data#> .
		d:y ab:tag "seven" .
		d:x ab:tag "eight" .
	`))
	request.Header.Add("Content-Type", "application/x-turtle")
	client.Do(request)
}

// func deleteDataset() {
// 	client := &http.Client{}
// 	request, _ := http.NewRequest("DELETE", SparqlEndpoint+"/kb", nil)
// 	response, err := client.Do(request)
// 	if err != nil {
// 		fmt.Println(err)
// 	}
// 	b, _ := ioutil.ReadAll(response.Body)
// 	response.Body.Close()
// 	fmt.Println(string(b))
// 	// Output:
// 	//
// }
