package main

import (
	"fmt"
	"log"
	"os"

	"github.com/learning-go-book/formatter"
	"github.com/shopspring/decimal"
)

// Store one module per repo
// Maintaining two modules means tracking separate versions fo two different projects

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Need two parameters: amount and percent")
	}

	amount, err := decimal.NewFromString(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	percent, err := decimal.NewFromString(os.Args[2])
	if err != nil {
		log.Fatal(err)
	}

	percent = percent.Div(decimal.NewFromInt(100))
	total := amount.Add(amount.Mul(percent)).Round(2)
	fmt.Println(formatter.Space(80, os.Args[1], os.Args[2], total.StringFixed(2)))
}
