#!/bin/bash

SRC=$(realpath $(cd -P "$(dirname "${BASH_SOURCE[0]}")" && pwd))

TEST=$(basename $SRC)

declare -A DSNS
DSNS+=(
  [mysql]=my://$TEST:$TEST@localhost/$TEST
  [oracle]=or://$TEST:$TEST@localhost/free
  [postgres]=pg://$TEST:$TEST@localhost/$TEST
  [sqlite3]=sq:$TEST.db
  [sqlserver]=ms://$TEST:$TEST@localhost/$TEST
)

APPLY=0
BUILD=0
DATABASES="mysql oracle postgres sqlite3 sqlserver"
ARGS=()

OPTIND=1
while getopts "abd:v" opt; do
case "$opt" in
  a) APPLY=1 ;;
  b) BUILD=1 ;;
  d) DATABASES=$OPTARG ;;
  v) ARGS+=(-v) ;;
esac
done

if [ "$BUILD" = "1" ]; then
  pushd $SRC/../../ &> /dev/null
  (set -x;
    go build
  )
  popd &> /dev/null
fi

DBTPLBIN=$(which dbtpl)
if [ -e $SRC/../../dbtpl ]; then
  DBTPLBIN=$SRC/../../dbtpl
fi
DBTPLBIN=$(realpath $DBTPLBIN)

pushd $SRC &> /dev/null

for TYPE in $DATABASES; do
  DB=${DSNS[$TYPE]}
  if [ -z "$DB" ]; then
    echo "$TYPE has no defined DSN"
    exit 1
  fi
  mkdir -p $TYPE
  rm -f $TYPE/*.dbtpl.*
  echo "------------------------------------------------------"
  echo "$TYPE: $DB"
  if [ "$APPLY" = "1" ]; then
    if [[ "$TYPE" = "sqlite3" && -f $TEST.db ]]; then
      (set -ex;
        rm $TEST.db
      )
    fi
    (set -ex;
      $SRC/../createdb.sh -d $TYPE -n $TEST
      ./init.sh -d $TYPE
    )
    if [ -f sql/${TYPE}_data.sql ]; then
      (set -ex;
        usql -f sql/${TYPE}_data.sql $DB
      )
    fi
    if [ -f sql/${TYPE}_post.sql ]; then
      (set -ex;
        usql -f sql/${TYPE}_post.sql $DB
      )
    fi
  fi
  (set -ex;
    $DBTPLBIN schema $DB -o $TYPE             ${ARGS[@]} --go-initialism ISBN
    $DBTPLBIN schema $DB -o $TYPE -t createdb ${ARGS[@]} --createdb-fmt=""
    $DBTPLBIN schema $DB -o $TYPE -t json     ${ARGS[@]}
    $DBTPLBIN schema $DB -o $TYPE -t yaml     ${ARGS[@]}
    $DBTPLBIN schema $DB -o $TYPE -t dot      ${ARGS[@]}
    go build ./$TYPE
    go build
    ./$TEST -dsn $DB ${ARGS[@]}
  )
done

popd &> /dev/null
