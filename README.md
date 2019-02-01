codimd-hugo-sync
================

Golang program for watching for CodeMD events in its Postgres database and adding labeled posts to a Hugo blog. It does this by using PostgreSQL's NOTIFY/LISTEN pubsub on the "Posts" table for INSERT, UPDATE and DELETE. Posts marked public and locked will be posted and automatically modified when saved in CodiMD.

### Gettings started
1. Edit the settings.json file with relevant database settings
    ```json
    {
        "database": {
            "host": "postgres.ip.address",
            "port": 5432,
            "dbname": "hackmd",
            "user": "hackmd",
            "password": "hackmdpass"
        }
    }
    ```
2. Build project
    ```bash
    $ GOOS=linux GOARCH=amd64 go build -o syncer
    ```
3. Run `./syncer`
