#!/usr/bin/env bash

function join_by {
  local d=${1-} f=${2-}
  if shift 2; then
    printf %s "$f" "${@/#/$d}"
  fi
}

autoloads=();

for x in $(find src/rdpcloud/lib/proto/* -mindepth 1 -maxdepth 1 -not \( -path "src/rdpcloud/lib/proto/GPBMetadata/*" -prune \) -type d -print); do
	IFS='/'; read -ra arr <<< "$x"; unset IFS;
	autoloads+=("\"${arr[4]}\\\\${arr[5]}\\\\\": \"${arr[2]}/${arr[3]}/${arr[4]}/${arr[5]}\"")
done

for x in $(find src/rdpcloud/lib/proto/GPBMetadata/* -mindepth 1 -maxdepth 1 -type d -print); do
	IFS='/'; read -ra arr <<< "$x"; unset IFS;
	autoloads+=("\"${arr[4]}\\\\${arr[5]}\\\\${arr[6]}\\\\\": \"${arr[2]}/${arr[3]}/${arr[4]}/${arr[5]}/${arr[6]}\"")
done

cat >src/rdpcloud/composer.json <<EOL
{
	"require": {
		"grpc/grpc": "^v1.42.0",
		"google/protobuf": "^v3.21.6"
	},
	"autoload": {
		"psr-4": {
			$(join_by $', \n\t\t\t' "${autoloads[@]}")
		}
	}
}
EOL
