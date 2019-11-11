#!/bin/bash

MYTMPDIR=$(mktemp -d /tmp/addLicenseHeaderToSourceFiles.XXXX)
LICENSE_HEADER_FILE="utils/LICENSE_header.txt"

rm_tmpdir() {
  rm -rf "${MYTMPDIR}"
}

sha1file() {
  sha1sum "${1}" | cut -d ' ' -f1
}

addLicenseHeaderToFile() {
  tmpFile="${1}.license_tmp"
  cat "${LICENSE_HEADER_FILE}" "${1}" > "${tmpFile}"
  mv "${tmpFile}" "${1}"
}

trap "rm_tmpdir" EXIT

find . -type f \
  -not -path './vendor/*' \
  -name '*.go' > ${MYTMPDIR}/gofiles.txt


refsha=$(sha1file "${LICENSE_HEADER_FILE}")

for file in $(< ${MYTMPDIR}/gofiles.txt)
do

  sanitizedFilename=$(echo "${file}" | sed -e 's#/#_#g')

  hfile="${MYTMPDIR}/${sanitizedFilename}.headers"
  head -30 ${file} > ${hfile}
  headsha=$(sha1file "${hfile}")

  if [ "${headsha}" != "${refsha}" ]
  then
    echo "Adding License headers to ${file}"
    addLicenseHeaderToFile "${file}"
  fi
done
