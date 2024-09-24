package main

import (
    "fmt"
    "os"

    "github.com/ArthurHydr/SSTI-Revealer/src/ssti"
)

func main() {
    if err := ssti.Run(); err != nil {
        fmt.Println(err)
        os.Exit(1)
    }
}
