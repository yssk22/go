package main

import (
	"crypto/sha256"
	"flag"
	"fmt"
	"os"

	"github.com/yssk22/go/x/xcrypto/xhmac"
)

var (
	unsign = flag.Bool("u", false, "unsign the key string")
	key    = flag.String("k", "", "hmac key string")
)

func main() {
	flag.Parse()
	if *key == "" {
		fmt.Fprintf(os.Stderr, "-k [key] must be specified\n")
		os.Exit(1)
	}
	hmac := xhmac.NewBase64([]byte(*key), sha256.New)
	args := flag.Args()
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, "no string is specified\n")
		os.Exit(1)
	}
	if *unsign {
		str, err := hmac.UnsignString(args[0])
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		fmt.Println(str)
	} else {
		fmt.Println(hmac.SignString(args[0]))
	}
}
