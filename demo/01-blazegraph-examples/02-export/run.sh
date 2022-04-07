#!/usr/bin/env bash

# *****************************************************************************

run_cell SETUP "INITIALIZE BLAZEGRAPH INSTANCE WITH ADDRESS BOOK" << END_CELL

blaze destroy --dataset kb --quiet
blaze create --dataset kb --quiet
blaze import --file ../data/address-book.jsonld --format jsonld

END_CELL

# *****************************************************************************

run_cell S1 "EXPORT ADDRESS BOOK AS JSON-LD" << END_CELL

blaze export --format jsonld

END_CELL

# *****************************************************************************

run_cell S1 "EXPORT ADDRESS BOOK AS TURTLE" << END_CELL

blaze export --format ttl

END_CELL

# *****************************************************************************

run_cell S1 "EXPORT ADDRESS BOOK AS N-TRIPLES" << END_CELL

blaze export --format nt | sort

END_CELL

# *****************************************************************************

run_cell S1 "EXPORT ADDRESS BOOK AS RDF-XML" << END_CELL

blaze export --format xml

END_CELL
