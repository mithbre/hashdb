package main

import (
    "fmt"
    "os"
)


func queryMain(args []string) {
    db, err := DBInit("../test.db")
    if err != nil {
        os.Exit(1)
    }

    switch args[0] {
    case "sha1":
        fmt.Printf("Query sha1: Not Implemented\n")
    case "sha256":
        fmt.Printf("Query sha256: Not Implemented\n")
    case "name":
        fmt.Printf("Query name: Not Implemented\n")
    case "path":
        fmt.Printf("Query path: Not Implemented\n")
    case "size":
        // size {+, -} int{"", "K", "M", "G", "T"}
        fmt.Printf("Query size: Not Implemented\n")
    default:
        os.Exit(1)
    }
}
