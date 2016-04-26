package main

import (
    "path/filepath"
    "fmt"
    "os"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"

)

func trackedFileDupe(db *sql.DB, file string, path string) {
    var fid string
    var sha1 []byte
    var sha2 []byte
    var dfile string
    var dpath string

    absFile, _ := filepath.Abs(file)
    absPath, _ := filepath.Abs(path)

    fi := filepath.Clean(absFile)
    pa := filepath.Clean(absPath)

    fiPath, filename := filepath.Split(fi)
    fiPath = filepath.Clean(fiPath)

    err := db.QueryRow(
        `SELECT id_file, sha1, sha256
         FROM tblFile
         INNER JOIN tblPath
         ON tblFile.path = tblPath.id_path
         WHERE tblPath.path = ?
         AND tblFile.filename = ?`, fiPath, filename).Scan(&fid, &sha1, &sha2)

    switch {
    case err == sql.ErrNoRows:
        fmt.Printf("File NOT found.\n")
        os.Exit(1)
    case err != nil:
        fmt.Printf("%s\n", err)
        os.Exit(1)
    default:
        fmt.Printf("Found File.\n")
    }

    res, err := db.Query(
        `SELECT filename, tblPath.path
         FROM tblFile
         INNER JOIN tblPath
         ON tblFile.path = tblPath.id_path
         WHERE tblPath.path LIKE ?
         AND sha1 = ?
         AND sha256 = ?
         AND id_file != ?`, pa+"%", sha1, sha2, fid)

    if err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
    defer res.Close()

    fmt.Printf("Results:\n")
    for res.Next() {
        err := res.Scan(&dfile, &dpath)
        switch err {
        case sql.ErrNoRows:
           fmt.Printf("No rows returned.\n")
        default:
            fmt.Println(dfile, dpath)
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
