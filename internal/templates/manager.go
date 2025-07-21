package templates

import (
	"bytes"
	"html/template"
	"net/url"
	"strings"
)

type TemplateManager struct {
	templates *template.Template
}

func NewTemplateManager() (*TemplateManager, error) {
	tmpl, err := template.ParseGlob("templates/*.html")
	if err != nil {
		return nil, err
	}
	
	tmpl, err = tmpl.ParseGlob("templates/landing/*.html")
	if err != nil {
		return nil, err
	}
	
	return &TemplateManager{templates: tmpl}, nil
}

type PDFData struct {
	Title    string
	Filename string
	FileURL  string
}

type WiFiData struct {
	Title       string
	SSID        string
	Password    string
	Encryption  string
	WiFiURI     string
	DownloadURL string
}

type SMSData struct {
	Title       string
	PhoneNumber string
	Message     string
	SMSURI      string
}

type EmailData struct {
	Title     string
	Recipient string
	Subject   string
	Body      string
	MailtoURL string
}

type VCardData struct {
	Title        string
	Name         string
	Organization string
	Title_       string
	Phone        string
	Email        string
	Address      string
	VCardURL     string
}

type EventData struct {
	Title       string
	EventName   string
	StartDate   string
	EndDate     string
	Location    string
	Description string
	ICalURL     string
}

func (tm *TemplateManager) RenderPDF(data PDFData) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, "base.html", data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (tm *TemplateManager) RenderWiFi(data WiFiData) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, "base.html", data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (tm *TemplateManager) RenderSMS(data SMSData) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, "base.html", data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (tm *TemplateManager) RenderEmail(data EmailData) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, "base.html", data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (tm *TemplateManager) RenderVCard(data VCardData) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, "base.html", data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

func (tm *TemplateManager) RenderEvent(data EventData) (string, error) {
	var buf bytes.Buffer
	err := tm.templates.ExecuteTemplate(&buf, "base.html", data)
	if err != nil {
		return "", err
	}
	return buf.String(), nil
}

// Helper functions for generating URIs and data URLs
func GenerateWiFiURI(ssid, password, encryption string, hidden bool) string {
	hiddenParam := ""
	if hidden {
		hiddenParam = "H:true;"
	}
	return "WIFI:T:" + encryption + ";S:" + ssid + ";P:" + password + ";" + hiddenParam + ";"
}

func GenerateSMSURI(phone, message string) string {
	return "SMSTO:" + phone + ":" + message
}

func GenerateMailtoURL(recipient, subject, body string) string {
	return "mailto:" + recipient + "?subject=" + url.QueryEscape(subject) + "&body=" + url.QueryEscape(body)
}

func GenerateVCardURL(name, organization, title, phone, email, address string) string {
	vcard := "BEGIN:VCARD\nVERSION:3.0\nFN:" + name + "\n"
	if organization != "" {
		vcard += "ORG:" + organization + "\n"
	}
	if title != "" {
		vcard += "TITLE:" + title + "\n"
	}
	if phone != "" {
		vcard += "TEL;TYPE=CELL:" + phone + "\n"
	}
	if email != "" {
		vcard += "EMAIL:" + email + "\n"
	}
	if address != "" {
		vcard += "ADR;TYPE=WORK:;;" + address + "\n"
	}
	vcard += "END:VCARD"
	return "data:text/vcard;charset=utf-8," + url.QueryEscape(vcard)
}

func GenerateICalURL(eventName, start, end, location, description string) string {
	formatICal := func(dt string) string {
		dt = strings.ReplaceAll(dt, "-", "")
		dt = strings.ReplaceAll(dt, ":", "")
		dt = strings.ReplaceAll(dt, "T", "T")
		if len(dt) >= 15 {
			return dt[:15]
		}
		return dt + strings.Repeat("0", 15-len(dt))
	}
	
	crlf := "\r\n"
	ical := "BEGIN:VCALENDAR" + crlf
	ical += "VERSION:2.0" + crlf
	ical += "BEGIN:VEVENT" + crlf
	ical += "SUMMARY:" + eventName + crlf
	if start != "" {
		ical += "DTSTART:" + formatICal(start) + crlf
	}
	if end != "" {
		ical += "DTEND:" + formatICal(end) + crlf
	}
	if location != "" {
		ical += "LOCATION:" + location + crlf
	}
	if description != "" {
		ical += "DESCRIPTION:" + description + crlf
	}
	ical += "END:VEVENT" + crlf + "END:VCALENDAR"
	
	return "data:text/calendar;charset=utf-8," + url.QueryEscape(ical)
}