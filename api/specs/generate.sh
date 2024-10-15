#!/bin/bash

DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

set -e

# if npm is not installed, error out
if ! [ -x "$(command -v npm)" ]; then
  echo "### Error: npm is not installed. Please install node & npm and try again."
  exit 1
fi

# if node_modules is not installed run npm install
if [ ! -d "$DIR/node_modules" ]; then
  echo "### Installing dependencies"
  npm install
fi

echo "### Generating OpenAPI spec from TypeScript"
node ./node_modules/.bin/tsp compile main.tsp

echo "### Copying generated spec to swagger-ui"
cp $DIR/tsp-output/@typespec/openapi3/openapi.json $DIR/../../swagger-ui/openapi.json

echo "### Removing tsp-output/"
rm -rf $DIR/tsp-output