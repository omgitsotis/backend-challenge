package main

import (
    "log"
    "github.com/omgitsotis/backend-challenge/client"
    "github.com/omgitsotis/backend-challenge/dao/sqlite"
)

func main() {
    dao, err := sqlite.NewSQLiteDAO("./fatlama.sqlite3")
    if err != nil {
        panic(err)
    }

    log.Fatal(client.ServeAPI(dao))
}
