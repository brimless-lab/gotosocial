/*
   GoToSocial
   Copyright (C) 2021-2023 GoToSocial Authors admin@gotosocial.org

   This program is free software: you can redistribute it and/or modify
   it under the terms of the GNU Affero General Public License as published by
   the Free Software Foundation, either version 3 of the License, or
   (at your option) any later version.

   This program is distributed in the hope that it will be useful,
   but WITHOUT ANY WARRANTY; without even the implied warranty of
   MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
   GNU Affero General Public License for more details.

   You should have received a copy of the GNU Affero General Public License
   along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package email

import (
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func loadTemplates(templateBaseDir string) (*template.Template, error) {
	if !filepath.IsAbs(templateBaseDir) {
		cwd, err := os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("error getting current working directory: %s", err)
		}
		templateBaseDir = filepath.Join(cwd, templateBaseDir)
	}

	// look for all templates that start with 'email_'
	return template.ParseGlob(filepath.Join(templateBaseDir, "email_*"))
}

// https://datatracker.ietf.org/doc/html/rfc2822
// I did not read the RFC, I just copy and pasted from
// https://pkg.go.dev/net/smtp#SendMail
// and it did seem to work.
func assembleMessage(mailSubject string, mailBody string, mailTo string, mailFrom string) ([]byte, error) {
	if strings.Contains(mailSubject, "\r") || strings.Contains(mailSubject, "\n") {
		return nil, errors.New("email subject must not contain newline characters")
	}

	if strings.Contains(mailFrom, "\r") || strings.Contains(mailFrom, "\n") {
		return nil, errors.New("email from address must not contain newline characters")
	}

	if strings.Contains(mailTo, "\r") || strings.Contains(mailTo, "\n") {
		return nil, errors.New("email to address must not contain newline characters")
	}

	// normalize the message body to use CRLF line endings
	mailBody = strings.ReplaceAll(mailBody, "\r\n", "\n")
	mailBody = strings.ReplaceAll(mailBody, "\n", "\r\n")

	msg := []byte(
		"To: " + mailTo + "\r\n" +
			"Subject: " + mailSubject + "\r\n" +
			"\r\n" +
			mailBody + "\r\n",
	)

	return msg, nil
}

// validateLine checks to see if a line has CR or LF as per RFC 5321
func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}

func SendMail(addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	if err := validateLine(from); err != nil {
		return err
	}
	for _, recp := range to {
		if err := validateLine(recp); err != nil {
			return err
		}
	}

	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	// Here is the key, you need to call tls.Dial instead of smtp.Dial
	// for smtp servers running on 465 that require an ssl connection
	// from the very beginning (no starttls)
	conn, err := tls.Dial("tcp", addr, &tls.Config{ServerName: host})
	if err != nil {
		return err
	}

	c, err := smtp.NewClient(conn, host)
	if err != nil {
		return err
	}
	defer c.Close()

	if a != nil {
		if err = c.Auth(a); err != nil {
			return err
		}
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
