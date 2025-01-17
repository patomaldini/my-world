package main

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/TwiN/go-color"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"runtime"
	"strings"
	"time"
)

var start time.Time
var username, password, combolist string
var combo []string

// Function to handle parsing errors
func errorpars(err error) {
	if err != nil {
		log.Panic(err.Error())
	}
}
func main() {
	start = time.Now()
	runtime.GOMAXPROCS(1)
	// Command-line interface setup using urfave/cli
	app := &cli.App{
		Flags: []cli.Flag{
			// Flags for URL, username, password, rate limit, delay, and other options
			// Each flag has a corresponding destination variable to store the value

			&cli.StringFlag{
				Name:        "username",
				Value:       "",
				Aliases:     []string{"u"},
				Destination: &username,
				Usage:       "Enter Username Wordlist",
			},
			&cli.StringFlag{
				Name:        "password",
				Value:       "",
				Aliases:     []string{"p"},
				Destination: &password,
				Usage:       "Enter Password Wordlist",
			},
			&cli.StringFlag{
				Name:        "combolist",
				Value:       "",
				Aliases:     []string{"c"},
				Destination: &combolist,
				Usage:       "Enter Combo Wordlist output",
			},
		},
		Action: func(cCtx *cli.Context) error {
			// Switch case to handle different scenarios based on provided options
			switch {

			case cCtx.String("username") == "":
				fmt.Println(color.Colorize(color.Red, "[-] Please Enter Username Wordlist Address with -u"))
			case cCtx.String("password") == "":
				fmt.Println(color.Colorize(color.Red, "[-] Please Enter Password Wordlist Address with -u"))
			case cCtx.String("combolist") == "":
				fmt.Println(color.Colorize(color.Red, "[-] Please Enter Password Wordlist Address with -u"))

			}
			return nil
		},
	}
	// Run the CLI app
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
	usernamedic, err := os.OpenFile(username, os.O_RDONLY, 0600)
	errorpars(err)
	defer func(usernamedic *os.File) {
		err := usernamedic.Close()
		errorpars(err)
	}(usernamedic)

	passwordsdic, err := os.OpenFile(password, os.O_RDONLY, 0600)
	errorpars(err)
	defer func(passwordsdic *os.File) {
		err := passwordsdic.Close()
		errorpars(err)
	}(passwordsdic)

	linesuser := make([]string, 0)
	linespassw := make([]string, 0)

	done := make(chan struct{})

	// Use bufio.Scanner for reading files
	processFile := func(file *os.File, lines *[]string) {
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			*lines = append(*lines, scanner.Text())
		}
		done <- struct{}{}
	}

	// Read content of wordlist files
	go processFile(usernamedic, &linesuser)
	go processFile(passwordsdic, &linespassw)

	// Wait for both files to finish reading
	<-done
	<-done

	// Add jobs to the queue
	for _, usern := range linesuser {
		for _, passw := range linespassw {
			com := fmt.Sprintf("%s:%s\n", usern, passw)
			combo = append(combo, com)
		}
	}
	writeResults(combo, combolist)
}

// write result to file output
func writeResults(results []string, outname string) {

	path := "output"
	if _, err := os.Stat(path); errors.Is(err, os.ErrNotExist) {
		err := os.Mkdir(path, os.ModePerm)
		if err != nil {
			log.Println(err)
		}
	}

	writeToFile(results, path+"/"+outname)

}

// write text output file
func writeTextFile(file *os.File, results []string) {
	for _, co := range results {

		if _, err := file.WriteString(co); err != nil {
			log.Println(err)

		}
	}
}

// detect file format to save output file
func writeToFile(results []string, filePath string) {
	file, err := os.Create(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Println(err)
		}
	}(file)

	switch {
	case strings.HasSuffix(filePath, ".txt"):
		writeTextFile(file, results)
	default:
		log.Fatal("Unsupported file format just support .txt file")
	}
	elapsed := time.Since(start)
	fmt.Printf("Execution time: %s\n", elapsed)
}
