package main

import (
    "database/sql"
    "fmt"
    "log"
    "os"
    _ "github.com/mattn/go-sqlite3"
)

func SearchBanish(db *sql.DB) {
    var dfile string
    var dpath string
    var id int

    // simply compare
    res, err := db.Query(
       `SELECT id_file, filename, tblPath.path
        FROM tblFile x, tblPath
        JOIN (
            SELECT sha1, sha256
            FROM tblBanish) y
        ON y.sha1 = x.sha1 AND y.sha256 = x.sha256
        WHERE x.path = tblpath.id_path`)

    if err != nil {
        log.Println(err)
        os.Exit(6)
    }

    defer res.Close()

    fmt.Printf("Results:\n")
    for res.Next() {
        err := res.Scan(&id, &dfile, &dpath)
        switch err {
        case sql.ErrNoRows:
           fmt.Printf("No rows returned.\n")
        default:
            fmt.Println(id, dfile, dpath)

        }
    }
}

func InsertBanish(db *sql.DB, id string) {
    var sha1 []byte
    var sha256 []byte

    err := db.QueryRow("SELECT sha1, sha256 FROM tblFile WHERE id_file = ?",
        id).Scan(&sha1, &sha256)


    insBanish, _ := db.Prepare(`
        INSERT INTO tblBanish(sha1, sha256) VALUES(?,?);`)

    rmFile, _ := db.Prepare(`
        DELETE FROM tblFile WHERE id_file = ?;`)

    _, err = insBanish.Exec(sha1, sha256)
    if err != nil {
        log.Printf("%s", err)
        os.Exit(1)
    }

    _, err = rmFile.Exec(id)
    if err != nil {
        log.Printf("%s", err)
        os.Exit(1)
    }
}

func Banish(args []string) {
    if len(args) > 1 {
        fmt.Println("Too many values")
        os.Exit(1)
    }

    db, err := DBInit("../test.db")
    if err != nil {
        os.Exit(1)
    }

    if len(args) == 0 {
        SearchBanish(db)
    } else {
        // verify that it's a number
        InsertBanish(db, args[0])
    }
}
