package commands

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
	"time"

	"github.com/pquerna/otp/totp"
	"github.com/skip2/go-qrcode"
)

func getQrString(name, secret string) string {
	secret = strings.ToUpper(secret)

	// https://github.com/google/google-authenticator/wiki/Key-Uri-Format
	// otpauth://totp/Example:alice@google.com?secret=JBSWY3DPEHPK3PXP&issuer=Example
	return fmt.Sprintf("otpauth://totp/%s?secret=%s&issuer=%s", name, secret, name)
}

func outputQrCode(writer io.Writer, name, secret string) error {
	qrString := getQrString(name, secret)
	q, err := qrcode.New(qrString, qrcode.Medium)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error generating qr code:", err)
		return err
	}
	fmt.Fprint(writer, q.ToSmallString(false))

	return nil
}

func qrCode(writer io.Writer, name, secret string) error {
	if len(name) == 0 {
		fmt.Fprintln(os.Stderr, "Name required for QR code generation")
		return errors.New("name required")
	}

	if len(secret) != 0 {
		_, err := totp.GenerateCode(secret, time.Now())
		if err != nil {
			fmt.Fprintln(os.Stderr, "Invalid secret:", err)
			return err
		}
		return outputQrCode(writer, name, secret)
	}

	c, err := collectionFile.loader()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error loading collection:", err)
		return err
	}

	s, err := c.GetSecret(name)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get collection entry for %s: %s\n", name, err)
		return err
	}

	return outputQrCode(writer, s.Name, s.Value)
}
