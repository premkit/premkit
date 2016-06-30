
#!/bin/bash

# This script tests multiple packages and creates a consolidated cover profile
# See https://gist.github.com/hailiang/0f22736320abe6be71ce for inspiration.
# The list of packages to test is specified in testpackages.txt.

function die() {
  echo $*
  exit 1
}

# Initialize profile.cov
echo "mode: count" > $2

# Initialize error tracking
ERROR=""

# Test each package and append coverage profile info to profile.cov
for pkg in `cat $1`
do
    echo $pkg
    govendor test -v -covermode=count -coverprofile=profile_tmp.cov $pkg || ERROR="Error testing $pkg"
    if [ -f profile_tmp.cov ];
    then
      tail -n +2 profile_tmp.cov >> $2 || die "Unable to append coverage for $pkg"
      rm profile_tmp.cov
    fi
done

if [ ! -z "$ERROR" ]
then
    die "Encountered error, last error was: $ERROR"
fi

echo "Generated coverage profile $2"
