#!/bin/bash

set -e

TOTP=./totp-test
COLLECTION=testcollection.json
TEST_NBR=1

# Build
echo "Building ${TOTP}"
go build -o ${TOTP} -ldflags "-X main.version=$(./scripts/get-version.sh)" ./totp/...

# Basic commands
echo "${TEST_NBR}: Testing basic commands"
${TOTP}
${TOTP} --help

# Test version
RESULT=$(${TOTP} version | head --bytes=12)
echo "Result: ${RESULT}"
if [[ ! ${RESULT} = "totp version" ]]; then
    echo "FAIL: Version command"
    exit 1
fi

# Test secret name or secret required
RESULT=$(${TOTP} 2>&1 | head --lines=1)
echo "Result: ${RESULT}"
if [[ ! ${RESULT} = "Secret name or secret is required." ]]; then
    echo "FAIL: Secret name or secret required"
    exit 1
fi

# Test secret was given so additional arguments are not needed
RESULT=$(${TOTP} --secret seed additional arguments 2>&1 | head --lines=1)
echo "Result: ${RESULT}"
if [[ ! ${RESULT} = "Secret was given so additional arguments are not needed." ]]; then
    echo "FAIL: Secret was given so additional arguments are not needed"
    exit 1
fi

# Test too many arguments. Only one secret name is required
RESULT=$(${TOTP} too many arguments 2>&1 | head --lines=1)
echo "Result: ${RESULT}"
if [[ ! ${RESULT} = "Too many arguments. Only one secret name is required." ]]; then
    echo "FAIL: Too many arguments. Only one secret name is required"
    exit 1
fi

# Test completion output
RESULT=$(${TOTP} config completion | head --lines=1)
echo "Result: ${RESULT}"
if [[ ! ${RESULT} = "# bash completion for totp                                 -*- shell-script -*-" ]]; then
    echo "FAIL: Completion output"
    exit 1
fi

# Test collection reset
((TEST_NBR++))
echo "${TEST_NBR}: Testing config reset"
touch ${COLLECTION}
${TOTP} config reset --file ${COLLECTION} --yes
if test -f "${COLLECTION}"; then
    echo "FAIL: ${TEST_NBR}. Collection file not removed"
    exit 1
fi

# Test generate TOTP with secret
((TEST_NBR++))
echo "${TEST_NBR}: Testing generate TOTP w/ secret on CLI"
TIME="2019-06-23T20:00:00-05:00"
SECRET=seed
RESULT=$(${TOTP} --time ${TIME} --secret ${SECRET})
echo "Result: ${RESULT}"
if [[ ! ${RESULT} =~ ^335072$ ]]; then
    echo "FAIL: ${TEST_NBR}. Incorrect TOTP generated"
    exit 1
fi

# Test add secret
((TEST_NBR++))
ENTRY=entryname
SECRET=seed

echo "${TEST_NBR}: Testing config update for adding"
${TOTP} config update --file ${COLLECTION} ${ENTRY} ${SECRET}
RESULT=$(${TOTP} config list --file ${COLLECTION} | tail -1)
echo "Result: ${RESULT}"
if [[ ! ${RESULT} =~ ^${ENTRY}\ ${SECRET}\  ]]; then
    echo "FAIL: ${TEST_NBR}. Entry not added"
    exit 1
fi

# Test generate TOTP
((TEST_NBR++))
echo "${TEST_NBR}: Testing generate TOTP"
TIME="2019-06-23T20:00:00-05:00"
RESULT=$(${TOTP} --file ${COLLECTION} --time "${TIME}" ${ENTRY})
echo "Result: ${RESULT}"
if [[ ! ${RESULT} =~ ^335072$ ]]; then
    echo "FAIL: ${TEST_NBR}. Incorrect TOTP generated"
    exit 1
fi

# Test generate backward TOTP
((TEST_NBR++))
echo "${TEST_NBR}: Testing backward"
TIME="2019-06-23T20:00:00-05:00"
RESULT=$(${TOTP} --file ${COLLECTION} --time "${TIME}" --backward 300s ${ENTRY})
echo "Result: ${RESULT}"
if [[ ! ${RESULT} =~ ^962630$ ]]; then
    echo "FAIL: Incorrect backward TOTP generated"
    exit 1
fi

# Test generate forward TOTP
((TEST_NBR++))
echo "${TEST_NBR}: Testing forward"
TIME="2019-06-23T20:00:00-05:00"
RESULT=$(${TOTP} --file ${COLLECTION} --time "${TIME}" --forward 300s ${ENTRY})
echo "Result: ${RESULT}"
if [[ ! ${RESULT} =~ ^869438$ ]]; then
    echo "FAIL: ${TEST_NBR}. Incorrect forward TOTP generated"
    exit 1
fi


# Test stdin config
((TEST_NBR++))
echo "${TEST_NBR}: Testing stdin"
ENTRY=entryname
TIME="2019-06-23T20:00:00-05:00"
RESULT=$(cat ${COLLECTION} | ${TOTP} --stdio --time ${TIME} entryname)
echo "Result: ${RESULT}"
if [[ ! ${RESULT} =~ ^335072$ ]]; then
    echo "FAIL: ${TEST_NBR}. Incorrect TOTP generated"
    exit 1
fi

# Test update secret
((TEST_NBR++))
ENTRY=entryname
SECRET=seedseed

echo "${TEST_NBR}: Testing config update for updating"
${TOTP} config update --file ${COLLECTION} ${ENTRY} ${SECRET}
RESULT=$(${TOTP} config list --file ${COLLECTION} | tail -1)
echo "Result: ${RESULT}"
if [[ ! ${RESULT} =~ ^${ENTRY}\ ${SECRET}\  ]]; then
    echo "FAIL: ${TEST_NBR}. Entry not added"
    exit 1
fi

# Test rename secret
((TEST_NBR++))
ENTRY=entryname
NEWENTRY=newentryname
SECRET=seedseed

echo "${TEST_NBR}: Testing config rename"
${TOTP} config rename --file ${COLLECTION} ${ENTRY} ${NEWENTRY}
RESULT=$(${TOTP} config list --file ${COLLECTION} | tail -1)
echo "Result: ${RESULT}"
if [[ ! ${RESULT} =~ ^${NEWENTRY}\ ${SECRET}\  ]]; then
    echo "FAIL: ${TEST_NBR}. Entry not added"
    exit 1
fi

# Test delete secret
((TEST_NBR++))
ENTRY=newentryname

echo "${TEST_NBR}: Testing config delete"
${TOTP} config delete --file ${COLLECTION} --yes ${ENTRY}
RESULT=$(${TOTP} config list --file ${COLLECTION} | wc -l)
echo "Result: ${RESULT}"
if [[ ! ${RESULT} =~ ^2 ]]; then
    echo "FAIL: ${TEST_NBR}. Entry not deleted"
    exit 1
fi

echo "Removing ${COLLECTION}"
${TOTP} config reset --file ${COLLECTION} --yes
echo "Removing ${TOTP}"
rm ${TOTP}

echo Success
