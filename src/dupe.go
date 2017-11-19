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

    var leftPath string
    var leftPathID string

    // query all paths from leftAlias
    resLPath, err := db.Query(
       `Select id_path, path
        FROM tblPath
        JOIN tblAlias
        ON tblPath.alias = tblAlias.id_alias
        WHERE tblAlias.alias = ?
        ORDER BY path ASC`, leftAlias)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer resLPath.Close()

    for resLPath.Next() {
        err := resLPath.Scan(&leftPathID, &leftPath)
        if err != nil {
            fmt.Println(err)
        }

        // Query for this path in rightAlias, return the id_path
        var rightPath string
        err = db.QueryRow(
            `SELECT id_path
             FROM tblPath
             JOIN tblAlias
             ON tblPath.alias = tblAlias.id_alias
             WHERE tblAlias.alias = ?
             AND tblPath.path = ?`, rightAlias, leftPath).Scan(&rightPath)

        if err == sql.ErrNoRows {
            fmt.Println("[< PATH]", leftPath)
            continue
        } else if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }

        fmt.Println("Path:", leftPath)

        // Retrieve all the files in the current leftPath
        res, err := db.Query(
            `SELECT filename, sha1, sha256
             FROM tblFile
             WHERE path = ?`, leftPathID)

        if err != nil {
            fmt.Println(err)
            os.Exit(1)
        }
        defer res.Close()

        var leftFile string
        var sha1, sha256 []byte
        // Change this to a struct at some point.

        for res.Next() {
            err := res.Scan(&leftFile, &sha1, &sha256)
            if err != nil {
                fmt.Println(err)
            }

            // Retrieve filename from rightPath for leftFile matches
            rightFile, err := db.Query(
                `SELECT filename
                 FROM tblFile
                 WHERE path = ?
                 AND sha1 = ?
                 AND sha256 = ?`, rightPath, sha1, sha256)

            //fmt.Println(leftFile)
            if err != nil {
                fmt.Println(err)
                os.Exit(1)
            }

            defer rightFile.Close()

            var rightFn string
            count := 0
            for rightFile.Next() {
                count++
                err := rightFile.Scan(&rightFn)
                if err != nil {
                    fmt.Println(err)
                    os.Exit(1)
                }

                if leftFile == rightFn {
                    fmt.Println("[= FILE]", leftFile)
                } else {
                    fmt.Println("[! NAME]", leftFile, rightFn)
                }
            }

            if count == 0 {
                fmt.Println("[< FILE]", leftFile)
            }

            if rightFile.Err() != nil {
                fmt.Println(rightFile.Err())
            }
        }
    }
}
