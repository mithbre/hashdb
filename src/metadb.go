package main

import (
    "os"
    "log"
    "database/sql"
    _ "github.com/mattn/go-sqlite3"
)


func DBInit(path string) (*sql.DB, error) {
    os.Remove(path)
    db, err := sql.Open("sqlite3", path)
    if err != nil {
        os.Exit(1)
    }

    create := `
    CREATE TABLE tblFile(
    id_file     INTEGER PRIMARY KEY NOT NULL,
    filename    TEXT NOT NULL,
    sha1        BLOB,
    sha256      BLOB,
    leng        INTEGER NOT NULL,
    modtime     INTEGER NOT NULL,
    ext         TEXT,
    path        INTEGER REFERENCES tblPath(id_path) ON UPDATE CASCADE
    );
    
    CREATE TABLE tblPath(
    id_path     INTEGER PRIMARY KEY NOT NULL,
    path        TEXT NOT NULL
    );
    `

    _, err = db.Exec(create)
    if err != nil {
        log.Printf("%q: %s\n", err, create)
        os.Exit(1)
    }
    return db, nil
}


//func DBAppend(db *sql.DB, dirMap[string] []os.FileInfo) (*sql.DB, error) {

//}
