#!/bin/bash

COVERAGE_FILE=${PWD}/coverage.txt
COV_SORTED=${PWD}/coverallsorted.out

touch "$COVERAGE_FILE"

function test_package {
  DIR=".$1"
  DEP=$(go list -f '{{ join .Deps "\n" }}' "$DIR" | grep v2ray | tr '\n' ',')
  DEP=${DEP}$DIR
  RND_NAME=$(openssl rand -hex 16)
  COV_PROFILE=${RND_NAME}.out
  go test -coverprofile="$COV_PROFILE" -coverpkg="$DEP" "$DIR" || return
}

TEST_FILES=(./*_test.go)
if [ -f "${TEST_FILES[0]}" ]; then
  test_package ""
fi

# shellcheck disable=SC2044
for DIR in $(find ./* -type d ! -path "*.git*" ! -path "*vendor*" ! -path "*external*"); do
  TEST_FILES=("$DIR"/*_test.go)
  if [ -f "${TEST_FILES[0]}" ]; then
    test_package "/$DIR"
  fi
done

# merge out
while IFS= read -r -d '' OUT_FILE
do
  echo "Merging file ${OUT_FILE}"
  < "${OUT_FILE}" grep -v "mode: set" >> "$COVERAGE_FILE"
done <   <(find ./* -name "*.out" -print0)

< "$COVERAGE_FILE" sort -t: -k1 | grep -vw "testing" | grep -v ".pb.go" | grep -vw "vendor" | grep -vw "external" > "$COV_SORTED"
echo "mode: set" | cat - "${COV_SORTED}" > "${COVERAGE_FILE}"

bash <(curl -s https://codecov.io/bash) || echo 'Codecov failed to upload'