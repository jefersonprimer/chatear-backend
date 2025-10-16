package infrastructure

import (
	"bytes"
	"context"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"

	"github.com/jefersonprimer/chatear-backend/internal/notification/domain"
)

type SMTPSender struct {
	host          string
	port          string
	user          string
	password      string
	from          string
	templates     map[string]*template.Template
}

func NewSMTPSender() (*SMTPSender, error) {
	templates := make(map[string]*template.Template)
	
	// Load all templates from the templates directory
	templateDir := "internal/notification/infrastructure/templates"
	templateFiles := []string{"email.txt", "welcome.html", "magic_link.html"}
	
	for _, templateFile := range templateFiles {
		templatePath := filepath.Join(templateDir, templateFile)
		if _, err := os.Stat(templatePath); err == nil {
			t, err := template.ParseFiles(templatePath)
			if err != nil {
				return nil, fmt.Errorf("failed to parse template %s: %w", templateFile, err)
			}
			templateName := templateFile[:len(templateFile)-len(filepath.Ext(templateFile))]
			templates[templateName] = t
		}
	}

	return &SMTPSender{
		host:          os.Getenv("SMTP_HOST"),
		port:          os.Getenv("SMTP_PORT"),
		user:          os.Getenv("SMTP_USER"),
		password:      os.Getenv("SMTP_PASS"), // Fixed: use SMTP_PASS instead of SMTP_PASSWORD
		from:          os.Getenv("SMTP_FROM"),
		templates:     templates,
	}, nil
}

func (s *SMTPSender) Send(ctx context.Context, emailSend *domain.EmailSend) error {
	auth := smtp.PlainAuth("", s.user, s.password, s.host)

	var body bytes.Buffer
	
	// Use template if specified, otherwise use the body directly
	if emailSend.TemplateName != "" {
		if template, exists := s.templates[emailSend.TemplateName]; exists {
			if err := template.Execute(&body, emailSend); err != nil {
				return fmt.Errorf("failed to execute template %s: %w", emailSend.TemplateName, err)
			}
		} else {
			// Fallback to default template or plain body
			body.WriteString(emailSend.Body)
		}
	} else {
		body.WriteString(emailSend.Body)
	}

	// Determine content type based on template
	contentType := "text/plain"
	if emailSend.TemplateName != "" && filepath.Ext(emailSend.TemplateName) == ".html" {
		contentType = "text/html"
	}

	msg := []byte("To: " + emailSend.Recipient + "\r\n" +
		"Subject: " + emailSend.Subject + "\r\n" +
		"Content-Type: " + contentType + "; charset=UTF-8\r\n" +
		"\r\n" +
		body.String() + "\r\n")

	addr := fmt.Sprintf("%s:%s", s.host, s.port)

	return smtp.SendMail(addr, auth, s.from, []string{emailSend.Recipient}, msg)
}
