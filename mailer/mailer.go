package mailer

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

// ErrorResponse response
type ErrorResponse struct {
	Errors []struct {
		Message string `json:"message"`
		Field   string `json:"field"`
		Help    string `json:"help"`
	} `json:"errors"`
}

// ToError transform
func (errs ErrorResponse) ToError() error {
	var msg = ""
	for _, e := range errs.Errors {
		msg += fmt.Sprintf("%v", e)
	}

	return errors.New(msg)
}

// SendMailParams params
type SendMailParams struct {
	TemplateID  string
	Name        string
	Email       string
	Data        map[string]interface{}
	Attachments []*Attachment
}

// Config cfg
type Config struct {
	APIKey     string
	EmailAlias string
	NameAlias  string
	BCCMails   string
}

// Email mail and name
type Email struct {
	Email string
	Name  string
}

// Attachment attachment
type Attachment struct {
	Content     string `json:"content,omitempty"`
	Type        string `json:"type,omitempty"`
	Name        string `json:"name,omitempty"`
	Filename    string `json:"filename,omitempty"`
	Disposition string `json:"disposition,omitempty"`
	ContentID   string `json:"content_id,omitempty"`
}

// Mailer struct
type Mailer struct {
	*Config
}

// New instance
func New(config *Config) *Mailer {
	if config.APIKey == "" {
		panic("Sendgrid API Key is required")
	}
	return &Mailer{
		config,
	}
}

// Send send
func (mailer *Mailer) Send(params SendMailParams) error {
	if params.TemplateID == "" {
		return errors.New("Template ID is required")
	}

	var apiKey = mailer.APIKey
	var addressAlias = mailer.EmailAlias
	var nameAlias = mailer.NameAlias
	var bccMails = mailer.BCCMails
	var userName = params.Name
	var userEmail = params.Email

	var isDRMACMail = strings.Contains(userEmail, "@aol.com") || strings.Contains(userEmail, "@yahoo.com")
	if isDRMACMail {
		var index = strings.Index(userEmail, "@")
		if index != -1 {
			addressAlias = fmt.Sprintf("no-reply%s", userEmail[index:])
		}
	}

	m := mail.NewV3Mail()
	e := mail.NewEmail(nameAlias, addressAlias)
	m.SetFrom(e)

	p := mail.NewPersonalization()
	tos := []*mail.Email{
		mail.NewEmail(userName, userEmail),
	}
	p.AddTos(tos...)

	cc := []*mail.Email{}

	if bccMails != "" {
		var s1 = strings.Split(bccMails, "|")
		if len(s1) > 0 {
			for _, nameMail := range s1 {
				var s2 = strings.Split(nameMail, ",")
				if len(s2) == 2 {
					if s2[1] != userEmail {
						cc = append(cc, mail.NewEmail(s2[0], s2[1]))

					}
				}
			}
		} else {
			var s2 = strings.Split(bccMails, ",")
			if len(s2) == 2 {
				if s2[1] != userEmail {
					cc = append(cc, mail.NewEmail(s2[0], s2[1]))
				}
			}
		}
	}

	if len(cc) > 0 {
		p.AddBCCs(cc...)
	}

	m.AddPersonalizations(p)
	m.SetTemplateID(params.TemplateID)

	for key, value := range params.Data {
		p.SetDynamicTemplateData(key, value)
	}

	// Add attachments
	for _, at := range params.Attachments {
		m.AddAttachment(&mail.Attachment{
			Content:     at.Content,
			ContentID:   at.ContentID,
			Disposition: at.Disposition,
			Filename:    at.Filename,
			Name:        at.Name,
			Type:        at.Type,
		})
	}

	request := sendgrid.GetRequest(apiKey, "/v3/mail/send", "https://api.sendgrid.com")
	request.Method = "POST"
	var Body = mail.GetRequestBody(m)
	request.Body = Body
	response, err := sendgrid.API(request)

	if response.StatusCode >= 300 {
		var errs ErrorResponse
		if err = json.Unmarshal([]byte(response.Body), &errs); err != nil {
			return err
		}
		return errs.ToError()
	}
	return nil
}
