package main

import (
    "fmt"
    "os"
    "database/sql"
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


func dupeMain(args []string) {
    db, err := DBInit("../test.db")
    if err != nil {
        os.Exit(1)
    }

    switch args[0] {
    case "u", "untracked":
        fmt.Printf("Untracked mode: Not Implemented\n")
    default:
        trackedFileDupe(db, args[0], args[1])
    }
}
