package main

import (
    "bufio"
    "flag"
    "fmt"
    "github.com/bclement/textgen"
    "os"
)


var maxlen = flag.Uint("max", 32, "max length of generated text in words")
var chainlen = flag.Uint("prefix", 2, "length of markov chain prefix")

func main() {
    flag.Parse()
    args := flag.Args()
    if len(args) < 1 {
        fmt.Printf("Usage: %v [-max <int>] [-len <int>] files...", os.Args[0])
        os.Exit(1)
    }

    gen := textgen.NewGenerator(*chainlen)
    var err error
    for i := 0; err == nil && i < len(args); i += 1 {
        var f *os.File
        f, err := os.Open(args[i])
        if err == nil {
            r := bufio.NewReader(f)
            err = gen.Load(r)
        }
    }


    if err != nil {
        fmt.Printf("Problem reading files: %v\n", err)
        os.Exit(1)
    }

    w := bufio.NewWriter(os.Stdout)
    err = gen.Generate(w, *maxlen)
    fmt.Println()

    if err != nil {
        fmt.Printf("Problem generating text: %v\n", err)
        os.Exit(1)
    }
}

