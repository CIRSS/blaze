#!/usr/bin/env bash

run_cell SETUP "INITIALIZE BLAZEGRAPH INSTANCE WITH CITATIONS" << END_CELL

blaze destroy --dataset kb --quiet
blaze create --dataset kb --quiet
blaze import --file ../data/citations.ttl --format ttl

END_CELL


run_cell S1 "EXPORT CITATIONS AS N-TRIPLES" << END_CELL

blaze export --format nt | sort

END_CELL


run_cell S2 "WHICH PAPERS DIRECTLY CITE WHICH PAPERS?" \
    << END_CELL

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>

    SELECT DISTINCT ?citing_paper ?cited_paper
    WHERE {
        ?citing_paper c:cites ?cited_paper .
    }
    ORDER BY ?citing_paper ?cited_paper

END_QUERY

END_CELL


run_cell S3 "WHICH PAPERS DEPEND ON WHICH PRIOR WORK?" \
    << END_CELL

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>

    SELECT DISTINCT ?paper ?prior_work
    WHERE {
        ?paper c:cites+ ?prior_work .
    }
    ORDER BY ?paper ?prior_work

END_QUERY

END_CELL


run_cell S4 "WHICH PAPERS DEPEND ON PAPER A?" \
    << END_CELL

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>
    prefix : <http://learningsparql.com/ns/papers#>

    SELECT DISTINCT ?paper
    WHERE {
        ?paper c:cites+ :paperA .
    }
    ORDER BY ?paper

END_QUERY

END_CELL


run_cell S5 "WHICH PAPERS CITE A PAPER THAT CITES PAPER A?" \
    << END_CELL

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>
    prefix : <http://learningsparql.com/ns/papers#>

    SELECT DISTINCT ?paper
    WHERE {
        ?paper c:cites/c:cites :paperA .
    }
    ORDER BY ?paper

END_QUERY

END_CELL


run_cell S6 "WHICH PAPERS CITE A PAPER CITED BY PAPER D?" \
    << END_CELL

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>
    prefix : <http://learningsparql.com/ns/papers#>

    SELECT DISTINCT ?paper
    WHERE {
        ?paper c:cites/^c:cites :paperD .
        FILTER(?paper != :paperD)
    }
    ORDER BY ?paper

END_QUERY

END_CELL


run_cell S7 "WHAT RESULTS DEPEND DIRECTLY ON RESULTS REPORTED BY PAPER A?" \
    << END_CELL

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>
    prefix : <http://learningsparql.com/ns/papers#>

    SELECT DISTINCT ?result
    WHERE {
        :paperA c:reports/^c:uses/c:reports ?result
    }
    ORDER BY ?result

END_QUERY

END_CELL


run_cell S7 "WHAT RESULTS DEPEND DIRECTLY OR INDIRECTLY ON RESULTS REPORTED BY PAPER A?" \
    << END_CELL

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>
    prefix : <http://learningsparql.com/ns/papers#>

    SELECT DISTINCT ?result
    WHERE {
        :paperA c:reports/(^c:uses/c:reports)+ ?result
    }
    ORDER BY ?result

END_QUERY

END_CELL

