package main

import (
	"bufio"
	"bytes"
	"errors"
	"html/template"
	"log"
	"net/smtp"
	"os"
)

type Credentials struct {
	Sender   string
	SMTPhost string
	SMTPPort string
	Username string
	Password string
}

func sendMail(subject string, templatePath string, to []string, credentials Credentials) error {
	// Getting Email Body From Template
	var body bytes.Buffer
	template, err := template.ParseFiles(templatePath)
	template.Execute(&body, struct{ Name string }{Name: "Test"})
	auth := smtp.PlainAuth("", credentials.Username, credentials.Password, credentials.SMTPhost)
	if err != nil {
		log.Fatal("Error Connecting SMTP Server")
	}

	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"

	msg := "Subject: " + subject + "\n" + headers + "\n\n" + body.String()

	err = smtp.SendMail(
		credentials.SMTPhost+":"+credentials.SMTPPort,
		auth,
		credentials.Sender,
		to,
		[]byte(msg),
	)
	if err != nil {
		log.Panic(err)
		return errors.New("Failed to send email")
	}
	return err
}

func Shoot(subject string, bodyTemplate string, filePath string, credentials Credentials) error {
	// Open the text file
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// Use a scanner to read line by line
	scanner := bufio.NewScanner(file)

	// // Gouroutines
	// var wg sync.WaitGroup
	// errChan := make(chan error)

	for scanner.Scan() {
		recipient := scanner.Text()
		log.Printf("Sending email to %s", []string{recipient})
		err := sendMail(subject, bodyTemplate, []string{recipient}, credentials)
		if err != nil {
			log.Panicf("Error sending email to %s: %v", recipient, err)
		}
	}
	// Wait for all goroutines to finish
	wg.Wait()
	close(errChan)

	// Handle errors from goroutines
	for err := range errChan {
		if err != nil {
			log.Printf("Error sending email to some recipients: %v", err)
		}
	}
	// Check for any errors during scanning
	if err := scanner.Err(); err != nil {
		return err
	}

	return nil
}

// Gomail

func main() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: M24.exe <subject> <HTML_template_file> <email_list_file>")
	}
	// // Option 1: Environment Variables
	// var credentials Credentials
	// credentials.Sender = os.Getenv("SENDER_EMAIL")
	// credentials.SMTPhost = os.Getenv("SMTP_HOST")
	// credentials.SMTPPort = os.Getenv("SMTP_PORT")
	// credentials.Username = os.Getenv("SMTP_USERNAME")

	credentials := Credentials{
		Sender:   "test@test.com",
		SMTPhost: "smtp.gmail.com",
		SMTPPort: "587",
		Username: "test@test.com",
		Password: "Password",
	}

	var subject string = os.Args[1]
	var bodyTemplate string = os.Args[2]
	var filePath string = os.Args[3]

	if credentials.Sender == "" || credentials.SMTPhost == "" || credentials.SMTPPort == "" || credentials.Username == "" {
		log.Fatal("Missing required environment variables for email credentials")
	}

	err := Shoot(subject, bodyTemplate, filePath, credentials)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Email sent successfully!")

}
