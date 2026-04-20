package template

import (
	"bytes"
	_ "embed"
	"html/template"
	"time"
)

//go:embed email/verification.html
var verificationHTML string

var verificationTmpl = template.Must(template.New("verification").Parse(verificationHTML))

type VerificationData struct {
	Email           string
	VerificationURL string
	Year            int
}

func RenderVerification(email, verificationURL string) (string, error) {
	var buf bytes.Buffer
	err := verificationTmpl.Execute(&buf, VerificationData{
		Email:           email,
		VerificationURL: verificationURL,
		Year:            time.Now().Year(),
	})
	return buf.String(), err
}
