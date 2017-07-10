package main

import (
    "database/sql"
    "encoding/hex"
    "fmt"
    "os"
    _ "github.com/mattn/go-sqlite3"
)


func qsha1(db *sql.DB, checksum string) {
    var filename string
    var path string

    raw, _ := hex.DecodeString(checksum)

    res, err := db.Query(
        `SELECT tblFile.filename, tblPath.path
         FROM tblFile, tblPath
         WHERE sha1 = ? and tblFile.path = tblPath.id_path`, raw)
    if err != nil {
        fmt.Println(err)
    }
    defer res.Close()

    fmt.Printf("Results:\n")
    for res.Next() {
        err := res.Scan(&filename, &path)
        switch err {
        case sql.ErrNoRows:
            fmt.Printf("Sha1 hash is not in database.\n")
        default:
            fmt.Println(path, filename)
        }
    }
}


func queryMain(args []string) {
    db, err := DBInit("../test.db")
    if err != nil {
        os.Exit(1)
    }

    switch args[0] {
    case "sha1":
        qsha1(db, args[1])
    case "sha256":
        fmt.Printf("Query sha256: Not Implemented\n")
    default:
        os.Exit(1)
    }
}
