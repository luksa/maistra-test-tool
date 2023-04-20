#!/bin/bash

set -e

SUPPORTED_VERSIONS=("v2.2" "v2.3" "v2.4")

log() {
    echo "$*" | tee -a "$LOG_FILE"
}

logHeader() {
    log
    log "====== $*"
    log
}

requireSMCPVERSION() {
    if [ -z "$SMCPVERSION" ]; then
        echo "ERROR: must specify version in SMCPVERSION env var"
        exit 1
    fi
}

runAllTests() {
    requireSMCPVERSION

    local dir="$1"

    echo > "$LOG_FILE"
    if [ -n "$TESTGROUP" ]; then
        logHeader "Executing tests in group '$TESTGROUP' against SMCP $SMCPVERSION"
    else
        logHeader "Executing all tests against SMCP $SMCPVERSION"
    fi

    go test -timeout 2h -v -count 1 -p 1 "$dir/..." 2>&1 \
    | tee -a "$LOG_FILE" >(${GOPATH}/bin/go-junit-report > "$REPORT_FILE")
}

runTest() {
    requireSMCPVERSION

    local dir="$1"
    local testName="$2"

    echo > "$LOG_FILE"
    logHeader "Executing $testName against SMCP $SMCPVERSION"

    go test -timeout 30m -v -count 1 -p 1 -run "^$testName$" "$dir/" 2>&1 \
    | tee -a "$LOG_FILE" >(${GOPATH}/bin/go-junit-report > "$REPORT_FILE")
}

resetCluster() {
    echo
    echo "Resetting cluster by deleting namespaces used in the test suite"
    oc delete namespace istio-system bookinfo foo bar legacy mesh-external --ignore-not-found
}

main() {
    if [ -z "$SMCPVERSION" ]; then
        echo
        echo "Executing tests against all supported versions: ${SUPPORTED_VERSIONS[*]}"
        echo "    Note: To run tests against a specific version, set the SMCPVERSION environment variable."
    else
        echo "Executing tests against SMCP version $SMCPVERSION"
    fi
    echo

    if [ -n "${OCP_CRED_PSW}" ]; then
        oc login -u ${OCP_CRED_USR} -p ${OCP_CRED_PSW} --server=${OCP_API_URL} --insecure-skip-tls-verify=true
    elif [ -n "${OCP_TOKEN}" ]; then
        oc login --token=${OCP_TOKEN} --server=${OCP_API_URL} --insecure-skip-tls-verify=true
    fi

    testName="${TEST_CASE:-$1}"
    if [ -n "$testName" ]; then
        # find the directory containing the specified test
        dir=""
        file=""
        for f in $(find . -type f -name "*_test.go"); do
            if grep -q "func $testName(" "$f"; then
                if [ -z "$dir" ]; then
                    dir=$(dirname "$f")
                    file="$f"
                else
                    echo >&2 "ERROR: Multiple tests with the given name were found. Please ensure test names are unique!"
                    exit 1
                fi
            fi
        done

        if [ -z "$dir" ]; then
            echo >&2 "ERROR: Could not find test $testName"
            exit 1
        fi
        echo "Found $testName in file $file."
    fi

    if [ -z "$SMCPVERSION" ]; then
        for ver in ${SUPPORTED_VERSIONS[@]}; do
            export SMCPVERSION="$ver"
            export LOG_FILE="$PWD/tests/output_${SMCPVERSION}.log"
            export REPORT_FILE="$PWD/tests/report_${SMCPVERSION}.xml"
            if [ -z "$testName" ]; then
                runAllTests "$PWD/pkg/tests"
            else
                runTest "$dir" "$testName"
            fi
            resetCluster
        done

        echo
        echo "=================================================================="
        echo "The JUnit test reports are located in:"
        for ver in ${SUPPORTED_VERSIONS[@]}; do
            echo "    - $PWD/tests/report_${ver}.xml"
        done

    else

        SMCPVERSION="v${SMCPVERSION#v}" # prepend "v" if necessary
        export LOG_FILE="$PWD/tests/output_${SMCPVERSION}.log"
        export REPORT_FILE="$PWD/tests/report_${SMCPVERSION}.xml"
        if [ -z "$testName" ]; then
            runAllTests "$PWD/pkg/tests"
        else
            runTest "$dir" "$testName"
        fi

        echo
        echo "The JUnit test report is located in:"
        echo "    $REPORT_FILE"
    fi
}

time main $@
