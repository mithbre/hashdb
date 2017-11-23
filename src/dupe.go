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

func CompareAliases(leftAlias string, rightAlias string) {
    db, err := DBInit("../test.db")
    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }

    // struct at some point
    var leftFn, leftPath string
    var sha1, sha256 []byte

    // query all paths from leftAlias
    resLeft, err := db.Query(
        `SELECT filename, tblPath.path, sha1, sha256
         FROM tblFile
         JOIN tblAlias ON tblPath.alias = tblAlias.id_alias
         JOIN tblPath ON tblFile.path = tblPath.id_path
         WHERE tblAlias.alias = ?
         ORDER BY tblPath.path, filename ASC`, leftAlias)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer resLeft.Close()

    equal, err := db.Prepare(
        `SELECT filename
         FROM tblFile
         JOIN tblAlias ON tblPath.alias = tblAlias.id_alias
         JOIN tblPath ON tblFile.path = tblPath.id_path
         WHERE tblAlias.alias = ?
         AND tblPath.path = ?
         AND sha1 = ?
         AND sha256 = ?`)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer equal.Close()

    var prevLPath string
    for resLeft.Next() {
        err := resLeft.Scan(&leftFn, &leftPath, &sha1, &sha256)
        if err != nil {
            fmt.Println(err)
        }

        if prevLPath == "" {
            prevLPath = leftPath
        } else if leftPath == prevLPath {
            //store files
        } else {
            //print files
            fmt.Println("\nPath:", leftPath)
            prevLPath = leftPath
        }

        res, err := equal.Query(rightAlias, leftPath, sha1, sha256)

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        defer res.Close()

        var rightFn string
        count := 0
        for res.Next() {
            count++
            err := res.Scan(&rightFn)
            if err != nil {
                fmt.Println(err)
            }

            if leftFn == rightFn {
                fmt.Println("[= FILE]", leftFn)
            } else {
                fmt.Println("[! NAME]", leftFn, rightFn)
            }
        }

        if count == 0 {
            fmt.Println("[< FILE]", leftFn)
        }
    }
}
