#!/usr/bin/env bash

DOT_RUNNER='../../common/run_dot_examples.sh'
SCRIPT_RUNNER='../../common/run_script_example.sh'


bash ${SCRIPT_RUNNER} SETUP "INITIALIZE BLAZEGRAPH INSTANCE WITH CITATIONS" << END_SCRIPT

blaze destroy --dataset kb --quiet
blaze create --dataset kb --quiet
blaze import --file ../data/citations.ttl --format ttl

END_SCRIPT


bash ${SCRIPT_RUNNER} S1 "EXPORT CITATIONS AS N-TRIPLES" << END_SCRIPT

blaze export --format nt | sort

END_SCRIPT


bash ${SCRIPT_RUNNER} S2 "WHICH PAPERS DIRECTLY CITE WHICH PAPERS?" \
    << END_SCRIPT

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>

    SELECT DISTINCT ?citing_paper ?cited_paper
    WHERE {
        ?citing_paper c:cites ?cited_paper .
    }
    ORDER BY ?citing_paper ?cited_paper

END_QUERY

END_SCRIPT


bash ${SCRIPT_RUNNER} S3 "WHICH PAPERS DEPEND ON WHICH PRIOR WORK?" \
    << END_SCRIPT

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>

    SELECT DISTINCT ?paper ?prior_work
    WHERE {
        ?paper c:cites+ ?prior_work .
    }
    ORDER BY ?paper ?prior_work

END_QUERY

END_SCRIPT


bash ${SCRIPT_RUNNER} S4 "WHICH PAPERS DEPEND ON PAPER A?" \
    << END_SCRIPT

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>
    prefix : <http://learningsparql.com/ns/papers#>

    SELECT DISTINCT ?paper
    WHERE {
        ?paper c:cites+ :paperA .
    }
    ORDER BY ?paper

END_QUERY

END_SCRIPT


bash ${SCRIPT_RUNNER} S5 "WHICH PAPERS CITE A PAPER THAT CITES PAPER A?" \
    << END_SCRIPT

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>
    prefix : <http://learningsparql.com/ns/papers#>

    SELECT DISTINCT ?paper
    WHERE {
        ?paper c:cites/c:cites :paperA .
    }
    ORDER BY ?paper

END_QUERY

END_SCRIPT


bash ${SCRIPT_RUNNER} S6 "WHICH PAPERS CITE A PAPER CITED BY PAPER D?" \
    << END_SCRIPT

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

END_SCRIPT


bash ${SCRIPT_RUNNER} S7 "WHAT RESULTS DEPEND DIRECTLY ON RESULTS REPORTED BY PAPER A?" \
    << END_SCRIPT

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>
    prefix : <http://learningsparql.com/ns/papers#>

    SELECT DISTINCT ?result
    WHERE {
        :paperA c:reports/^c:uses/c:reports ?result
    }
    ORDER BY ?result

END_QUERY

END_SCRIPT


bash ${SCRIPT_RUNNER} S7 "WHAT RESULTS DEPEND DIRECTLY OR INDIRECTLY ON RESULTS REPORTED BY PAPER A?" \
    << END_SCRIPT

blaze query --format table << END_QUERY

    prefix c: <http://learningsparql.com/ns/citations#>
    prefix : <http://learningsparql.com/ns/papers#>

    SELECT DISTINCT ?result
    WHERE {
        :paperA c:reports/(^c:uses/c:reports)+ ?result
    }
    ORDER BY ?result

END_QUERY

END_SCRIPT

