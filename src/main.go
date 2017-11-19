package main

import (
    "bufio"
    "crypto/sha1"
    "crypto/sha256"
    "fmt"
    "io"
    "os"
    "strings"
    "path/filepath"
)

func usage() {
    fmt.Printf("Add a path to be tracked.\n")
    fmt.Printf("    $ hashdb padd path\n")
}


func checksum(path string, buf []byte) ([]byte, []byte) {
    s1 := sha1.New()
    s2 := sha256.New()

    fi, err := os.Open(path)
    defer fi.Close()
    if err != nil {
        fmt.Println("Failed to open given path.")
        fmt.Println(err)
        os.Exit(1)
    }

    fiRead := bufio.NewReader(fi)

    for {
        bytes, err := fiRead.Read(buf)

        if err == io.EOF {
            break
        } else if err != nil {
            fmt.Println("Read fail during buffering.")
            fmt.Println(err)
            os.Exit(1)
        }

        s1.Write(buf[0:bytes])
        s2.Write(buf[0:bytes])
    }

    return s1.Sum(nil), s2.Sum(nil)
}


func traverseDir(tree map[string] []os.FileInfo, path string) error {
    walk := func(path string, meta os.FileInfo, err error) error {
        if meta == nil {
            fmt.Printf("File or path is unreadable.\n")
            fmt.Println(err)
            fmt.Printf(path)
        }

        if meta.IsDir() {
            if meta.Mode() & os.ModeSymlink != 0 {
                fmt.Println("Special Directory. Skipping ->", path)
                return filepath.SkipDir
            }
            path = filepath.Clean(path)
            tree[path] = make([]os.FileInfo, 0)
        } else if !meta.Mode().IsRegular() {
            fmt.Println("Special file. Skipping ->", path)
            return nil
        } else {
            path, _ = filepath.Split(path)
            path = filepath.Clean(path)
            tree[path] = append(tree[path], meta)
        }
        return nil
    }

    err := filepath.Walk(path, walk)
    if err != nil {
        os.Exit(1)
    }
    return nil
}


func pathAdd(startPath string, alias string) {
    dirMap := make(map[string] []os.FileInfo)
    rootDir := filepath.Clean(startPath)
    rootDir, _ = filepath.Abs(rootDir)

    err := traverseDir(dirMap, rootDir)
    if err != nil {
        os.Exit(1)
    }

    db, err := DBInit("../test.db")
    if err != nil {
        os.Exit(1)
    }

    db, aliasID, _ := InsAlias(db, alias, rootDir)

    buf := make([]byte, 1024*1024*50)
    for path, values := range dirMap {
        //strip rootDir from path
        pathShort := strings.Replace(path[4:], rootDir, "", 1)

        if len(pathShort) == 0 {
            pathShort = `\`
        }

        db, pathID, _ := InsPath(db, pathShort, aliasID)
        tx, _ := db.Begin()
        for _, item := range values {
            sha1, sha2 := checksum(filepath.Join(path, item.Name()), buf)
            tx, err = DBAppend(tx, pathID, item, sha1, sha2)

            if tx == nil {
                os.Exit(1)
            }
        }
        tx.Commit()
    }
}


func main() {
    if len(os.Args) < 2 {
        usage()
        return
    }

    args := os.Args[1:]

    switch args[0] {
    case "add":
        //path, alias
        pathAdd(args[1], args[2])
    case "dupe":
        SelectDupes()
    case "comp":
        CompareAliases(args[1], args[2])
    case "rm":
        RmFile(args[1])
    case "banish":
        Banish(args[1:])
    default:
        usage()
    }
}
