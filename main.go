package main

import (
	"bufio"
	"errors"
	"fmt"
	"net"
	"net/mail"
	"os"
	"strings"
)

var ErrNoSPF = errors.New("no SPF records")
var ErrNoDMARC = errors.New("no DMARC records")

func main() {
	fmt.Printf("Please write an email address: ")
	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		email := scanner.Text()
		if !isValid(email) {
			fmt.Printf("Invalid syntax for the email.\n")
		} else {
			at := strings.LastIndex(email, "@") // email can have "@"s in the username but not in the domain
			username, domain := email[:at], email[at+1:]
			fmt.Printf("Valid syntax. Username: %v, Domain: %v\n", username, domain)

			hasMX, err := hasMX(domain)
			if err != nil {
				fmt.Printf("Error when checking MX: %v\n", err)
			} else {
				fmt.Printf("Has MX: %v\n", hasMX)
			}

			spfRecord, err := getSpfRecord(domain)
			if err != nil {
				if err == ErrNoSPF {
					fmt.Printf("No SPF record found\n")
				} else {
					fmt.Printf("Error when checking SPF record: %v\n", err)
				}
			} else {
				fmt.Printf("SPF record: %v\n", spfRecord)
			}

			dmarcRecord, err := getDmarcRecords(domain)
			if err != nil {
				if err == ErrNoDMARC {
					fmt.Printf("No DMARC record found\n")
				} else {
					fmt.Printf("Error when checking DMARC record: %v\n", err)
				}
			} else {
				fmt.Printf("DMARC record: %v\n", dmarcRecord)
			}
		}

		fmt.Printf("\nTry another one or exit with CTRL+C: ")
	}

	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: could not read from input: %v\n", err)
	}
}

func hasMX(domain string) (bool, error) {

	mxRecords, err := net.LookupMX(domain)

	if err != nil {
		return false, err
	}

	return len(mxRecords) > 0, nil
}

func getSpfRecord(domain string) (string, error) {

	txtRecords, err := net.LookupTXT(domain)
	if err != nil {
		return "", err
	}

	for _, record := range txtRecords {
		if strings.HasPrefix(record, "v=spf1") {
			return record, nil
		}
	}

	return "", ErrNoSPF
}

func getDmarcRecords(domain string) (string, error) {

	dmarcRecords, err := net.LookupTXT("_dmarc." + domain)
	if err != nil {
		return "", err
	}

	for _, record := range dmarcRecords {
		if strings.HasPrefix(record, "v=DMARC1") {
			return record, nil
		}
	}

	return "", ErrNoDMARC
}

func isValid(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}
