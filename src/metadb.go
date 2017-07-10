package main

import (
    "database/sql"
    "log"
    "os"
    "strconv"
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
    path        INTEGER REFERENCES tblPath(id_path) ON UPDATE CASCADE,
    UNIQUE (filename, path) ON CONFLICT IGNORE
    );
    
    CREATE TABLE IF NOT EXISTS tblPath(
    id_path     INTEGER PRIMARY KEY NOT NULL,
    path        TEXT NOT NULL
    );

    CREATE TABLE IF NOT EXISTS tblBanish(
    sha1    BLOB,
    sha256  BLOB,
    UNIQUE (sha1, sha256) ON CONFLICT IGNORE
    );

    CREATE TABLE IF NOT EXISTS tblDeleted(
    id_file     INTEGER PRIMARY KEY NOT NULL,
    filename    TEXT NOT NULL,
    sha1        BLOB,
    sha256      BLOB,
    leng        INTEGER NOT NULL,
    path        INTEGER REFERENCES tblPath(id_path) ON UPDATE CASCADE,
    UNIQUE (filename, path) ON CONFLICT IGNORE
    )
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

func DBAppend(db *sql.DB, path int64, item os.FileInfo,
    sha1 []byte, sha2 []byte) (*sql.DB, error) {
    /* Inserts Name, Size, UnixTime, PathID */
    insFile, err := db.Prepare(`
        INSERT INTO tblFile(filename, leng, path, modtime, sha1, sha256)
        VALUES(?, ?, ?, ?, ?, ?)`)
    if err != nil {
        log.Println(err)
        os.Exit(5)
    }

    _, err = insFile.Exec(item.Name(), item.Size(),
        path, item.ModTime().Unix(), sha1, sha2)
    if err != nil {
        log.Printf("%s", err)
        os.Exit(1)
    }
    return db, nil
}

func RmFile(id string) (*sql.DB, error) {
    /* Removes specified row from tblFile */
    db, err := DBInit("test.db")
    if err != nil {
        log.Println(err)
        os.Exit(1)
    }

    // Move row information to a temporary table
    mvRow, err := db.Prepare(`
        INSERT INTO tblDeleted(filename,sha1,sha256,leng,path)
        SELECT filename, sha1, sha256, leng, path
        FROM tblFile
        WHERE tblFile.id_file = ?;`)
    if err != nil {
        log.Println(err)
        os.Exit(5)
    }

    // Delete the row from tblFile
    rmFile, err := db.Prepare(`
        DELETE FROM tblFile
        WHERE id_file = ?`)
    if err != nil {
        log.Println(err)
        os.Exit(5)
    }

    _, err = mvRow.Exec(id)
    if err != nil {
        log.Println(err)
        os.Exit(6)
    }
    _, err = rmFile.Exec(id)
    if err != nil {
        log.Println(err)
        os.Exit(6)
    }

    return db, nil
}
