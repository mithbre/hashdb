package main

import (
    "bufio"
    "crypto/sha1"
    "crypto/sha256"
    "fmt"
    "io"
    "os"
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
        os.Exit(1)
    }

    fiRead := bufio.NewReader(fi)

    for {
        bytes, err := fiRead.Read(buf)

        if err == io.EOF {
            break
        } else if err != nil {
            os.Exit(1)
        }

        s1.Write(buf[0:bytes])
        s2.Write(buf[0:bytes])
    }

    return s1.Sum(nil), s2.Sum(nil)
}


func traverseDir(tree map[string] []os.FileInfo, path string) error {
    walk := func(path string, meta os.FileInfo, err error) error {
        path, _ = filepath.Abs(path)

        if meta.IsDir() {
            // filepath.SkipDir
            path = filepath.Clean(path)
            tree[path] = make([]os.FileInfo, 0)
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


func pathAdd(path string) {
    dirMap := make(map[string] []os.FileInfo)
    walkDir := filepath.Clean(path)

    err := traverseDir(dirMap, walkDir)
    if err != nil {
        os.Exit(1)
    }

    db, err := DBInit("../test.db")
    if err != nil {
        os.Exit(1)
    }

    buf := make([]byte, 1024*1024*50)
    for path, values := range dirMap {
        db, pathId, _ := InsPath(db, path)
        for _, item := range values {
            sha1, sha2 := checksum(filepath.Join(path, item.Name()), buf)
            db, err = DBAppend(db, pathId, item, sha1, sha2)

            if db == nil {
                os.Exit(1)
            }
        }
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
        pathAdd(args[1])
    case "query":
        queryMain(args[1:])
    case "dupe":
        SelectDupes()
    default:
        usage()
    }
}
