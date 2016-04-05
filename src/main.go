package main

import (
    //"fmt"
    "os"
    "path/filepath"
)


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


func main() {
    // TODO: parse arguments
    dirMap := make(map[string] []os.FileInfo)
    walkDir := filepath.Clean("../test")

    err := traverseDir(dirMap, walkDir)
    if err != nil {
        os.Exit(1)
    }

    _, err = DBInit("../test.db")
    if err != nil {
        os.Exit(1)
    }

    //db, err = DBAppend(db, dirMap)
}
