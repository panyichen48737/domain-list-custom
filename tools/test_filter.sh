#!/bin/bash
# Regression test for filter.sh.
#
#   bash ./tools/test_filter.sh
#
# filter.sh drops a DOMAIN line when a DOMAIN-SUFFIX line already covers it, and
# passes every other type through untouched. Both halves are pinned here, along
# with the fact that the check does not depend on order: a suffix found later in
# the file still covers a domain seen earlier.
set -u

here=$(cd "$(dirname "$0")" && pwd)
workdir=$(mktemp -d)
trap 'rm -rf "$workdir"' EXIT

failures=0

check() {
    name=$1
    input=$2
    expected=$3

    printf '%s' "$input" > "$workdir/input.list"
    bash "$here/filter.sh" "$workdir/input.list" > "$workdir/actual.list"
    printf '%s' "$expected" > "$workdir/expected.list"

    if diff -u "$workdir/expected.list" "$workdir/actual.list" > "$workdir/diff.txt"; then
        echo "ok   - $name"
    else
        echo "FAIL - $name"
        sed 's/^/       /' "$workdir/diff.txt"
        failures=$((failures + 1))
    fi
}

# `www.example.com`, `deep.sub.example.com` and `example.com` itself are all
# covered by the `example.com` suffix and go away. `other.net` is covered by a
# suffix that appears further down. Keywords and regexes are not domains, so
# they stay even when their text looks like one. The `.cn` line is covered by
# the parent `com.cn`.
check "covered domains are dropped, everything else is kept" \
'DOMAIN-SUFFIX,example.com
DOMAIN,www.example.com
DOMAIN,deep.sub.example.com
DOMAIN,example.com
DOMAIN,other.net
DOMAIN-SUFFIX,other.net
DOMAIN-KEYWORD,example.com
DOMAIN-REGEX,^ads\.example\.com$
DOMAIN,example.com.cn
DOMAIN-SUFFIX,com.cn
DOMAIN,uncovered.org
' \
'DOMAIN-SUFFIX,example.com
DOMAIN-SUFFIX,other.net
DOMAIN-KEYWORD,example.com
DOMAIN-REGEX,^ads\.example\.com$
DOMAIN-SUFFIX,com.cn
DOMAIN,uncovered.org
'

check "an empty file produces no output" '' ''

if [ "$failures" -ne 0 ]; then
    echo
    echo "${failures} test(s) failed"
    exit 1
fi

echo
echo "all tests passed"
