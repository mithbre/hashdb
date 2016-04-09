package main

import (
    "os"
    "log"
    "path/filepath"
    "strconv"
    "strings"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)


func DBInit(path string) (*sql.DB, error) {
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }

    create := `
    CREATE TABLE IF NOT EXISTS tblFile(
    id_file     INTEGER PRIMARY KEY NOT NULL,
    filename    TEXT NOT NULL,
    sha1        BLOB,
    sha256      BLOB,
    leng        INTEGER NOT NULL,
    modtime     INTEGER NOT NULL,
    ext         TEXT,
    path        INTEGER REFERENCES tblPath(id_path) ON UPDATE CASCADE,
    UNIQUE (filename, path) ON CONFLICT IGNORE
    );
    
    CREATE TABLE IF NOT EXISTS tblPath(
    id_path     INTEGER PRIMARY KEY NOT NULL,
    path        TEXT NOT NULL
    );
    `

    _, err = db.Exec(create)
    if err != nil {
        log.Printf("%q\n", err)
    }
    return db, nil
}


func InsPath(db *sql.DB, path string) (*sql.DB, int64, error) {
    /* Inserts path, and returns row ID */
    var row string
    var id int64

    // id is used in lieu of row (int vs str)
    err := db.QueryRow("SELECT id_path FROM tblPath WHERE path = ?",
            path).Scan(&row)

    switch {
    case err == sql.ErrNoRows:
        log.Printf("PATH INSERT: %s\n", path)
        insert, _ := db.Prepare("INSERT INTO tblPath(path) VALUES(?)")
        res, _ := insert.Exec(path)
        id, _ = res.LastInsertId()

    case err != nil:
        log.Printf("%s\n", err)
        os.Exit(3)

    default:
        log.Printf("Path Exists: %s\n", path)
        id, err = strconv.ParseInt(row, 10, 64)
    }

    return db, id, nil
}


func DBAppend(db *sql.DB, path int64, values []os.FileInfo) (*sql.DB, error) {
    /* Inserts Name, Size, UnixTime, PathID */
    insFile, err := db.Prepare(`INSERT INTO tblFile(filename, leng, modtime,
            path, ext) VALUES(?, ?, ?, ?, ?)`)
    if err != nil {
        log.Println(err)
        os.Exit(5)
    }

    for _, item := range values {
        ext := filepath.Ext(item.Name())
        if len(ext) > 1 {
            ext = strings.ToLower(ext)[1:len(ext)]
        }
        _, err := insFile.Exec(item.Name(), item.Size(),
                item.ModTime().Unix(), path, ext)
        if err != nil {
            log.Printf("%s", err)
            os.Exit(1)
        }
    }
    return db, nil
}
