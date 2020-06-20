#!/bin/sh
#
# This runs from pkg/schema/schema.go go:generate.
#
SRC=../../api/schema.graphql
TGT=schema_graphql.go

echo >  $TGT "// DO NOT EDIT - generated by $(basename $0)"
echo >> $TGT 'package schema'
echo >> $TGT 'const schemaString = `'
cat  >> $TGT $SRC
echo >> $TGT '`'
