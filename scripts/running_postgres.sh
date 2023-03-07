NAME=acmhomepage-postgres
DOCKER_RUN=true
SHOW_DOC=true

if [[ "$1" == "export" ]]; then
    DOCKER_RUN=false
    SHOW_DOC=false
    shift
fi
if [[ "$1" == "--quiet" ]]; then
    SHOW_DOC=false
    shift
fi

# To have a running database.
if [[ "$DOCKER_RUN" == true ]]; then
    docker run --name "$NAME" -d -p 5432:5432 -e POSTGRES_PASSWORD=pw postgres:15.1
fi

# Use the default value to connect.
POSTGRES_URL=postgres://postgres:pw@localhost:5432/postgres?sslmode=disable
export POSTGRES_URL=$POSTGRES_URL

# Show document
if [[ "$SHOW_DOC" == true ]]; then
    echo
    echo "    Now we have a running docker container named $NAME, which runing a"
    echo "    PostgreSQL service. Use docker to handle it."
    echo
    echo "    We export the enviroment variable POSTGRES_URL. Use command below to"
    echo "    re-export if you want:"
    echo
    echo "        . ./running_postgres.sh export"
    echo
    echo "    That means you should run the command above in your shell if you don't run"
    echo "    me by a dot command."
    echo
fi
