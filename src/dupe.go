package main

import (
    "database/sql"
    "fmt"
    "os"
    _ "github.com/mattn/go-sqlite3"

)

func SelectDupes() {
    db, err := DBInit("../test.db")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    var dfile string
    var dpath string
    var id int
    var leng int
    var prev int

    res, err := db.Query(
       `SELECT id_file, filename, tblPath.path, leng
        FROM tblFile x, tblPath
        JOIN (
            SELECT sha1, sha256
            FROM tblFile
            GROUP BY sha1, sha256
            HAVING COUNT (*) > 1) y
        ON y.sha1 = x.sha1 AND y.sha256 = x.sha256
        WHERE x.path = tblpath.id_path
        ORDER BY leng`)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer res.Close()

    fmt.Printf("Results:\n")
    for res.Next() {
        err := res.Scan(&id, &dfile, &dpath, &leng)
        switch err {
        case sql.ErrNoRows:
           fmt.Printf("No rows returned.\n")
        default:
            if leng != prev {
                fmt.Println("\n")
            }
            prev = leng
            fmt.Println(id, dfile, dpath)
        }
    }
}

func CompareAliases(aAlias string, bAlias string) {
    db, err := DBInit("../test.db")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // goal is to search for all aAlias files within bAlias
    res, err := db.Query(
       `WITH alias_select AS (
            SELECT id_file, filename, tblAlias.alias, tblAlias.root, tblPath.path, sha1, sha256
            FROM tblFile
            JOIN tblPath ON tblFile.path = tblPath.id_path
            JOIN tblAlias ON tblPath.alias = tblAlias.id_alias
            WHERE tblAlias.alias = ? OR tblAlias.alias = ?)
        SELECT id_file, filename, leng, path, alias, root
        FROM alias_select bar
        JOIN (
            SELECT sha1, sha256
            FROM alias_select
            GROUP BY sha1, sha256
            HAVING COUNT (*) > 1) foo
        ON foo.sha1 = bar.sha1 AND foo.sha256 = bar.sha256
        ORDER BY leng`, aAlias, bAlias)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer res.Close()

    var dfile string
    var dpath string
    var alias string
    var root string
    var id int
    var sha1 []byte
    var prev []byte

    fmt.Printf("Results:\n")
    for res.Next() {
        err := res.Scan(&id, &dfile, &dpath, &alias, &root)
        switch err {
        case sql.ErrNoRows:
           fmt.Printf("No rows returned.\n")
        default:
            if sha1[0] != prev[0] && sha1[1] != prev[1] {
                fmt.Println("\n")
            }
            prev = sha1
            fmt.Println(id, dfile, alias, root, dpath)
        }
    }
}

