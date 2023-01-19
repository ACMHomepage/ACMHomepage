# Backend for ACMHomepage

## Run for development

Firstly, we can get the running PostgreSQL by docker.

```bash
. ../scripts/running_postgres.sh
```

It will export some enviroment variables. If you changed your shell, you can
run the command below to re-export:

```bash
. ../scripts/running_postgres.sh export
```

If you have a running PostgreSQL, run this command to have a running backend:

```bash
go run .
```

