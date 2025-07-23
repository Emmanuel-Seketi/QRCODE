package handler

import (
	"context"
	"html/template"

	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

	"qr_backend/ent"
	"qr_backend/ent/qrcode"
	"qr_backend/ent/qrcodeanalytics"
	"qr_backend/internal/database"
	"qr_backend/internal/model"
	"qr_backend/pkg/barcode"
	qrgen "qr_backend/pkg/qrcode"
	"qr_backend/pkg/shorturl"

	"github.com/gofiber/fiber/v2"
	goqrcode "github.com/skip2/go-qrcode"
)

// CreateQRCode generates a new QR code
func CreateQRCode(c *fiber.Ctx) error {
	var req struct {
		Type        model.QRCodeType       `json:"type"`
		Title       string                 `json:"title"`
		Description string                 `json:"description,omitempty"`
		RedirectURL string                 `json:"redirect_url,omitempty"`
		ShortURL    string                 `json:"short_url,omitempty"`
		Content     map[string]interface{} `json:"content"`
		ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
		Analytics   bool                   `json:"analytics"`
		Active      bool                   `json:"active"`
		Tags        []string               `json:"tags,omitempty"`
		Design      map[string]interface{} `json:"design,omitempty"`
		GroupID     *int                   `json:"group_id,omitempty"`
		IsDynamic   bool                   `json:"is_dynamic"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Handle dynamic vs static QR code type
	if req.IsDynamic {
		req.Type = "dynamic"
		shortURL, err := shorturl.Generate()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate short URL"})
		}
		req.ShortURL = shortURL
		req.Analytics = true
	} else {
		req.Type = "static"
		if req.ShortURL == "" && req.Analytics {
			shortURL, err := shorturl.Generate()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate short URL"})
			}
			req.ShortURL = shortURL
		}
	}

	// Create QR code using Ent
	qrBuilder := database.DB.QRCode.Create().
		SetType(string(req.Type)).
		SetTitle(req.Title).
		SetContent(req.Content).
		SetAnalytics(req.Analytics).
		SetActive(req.Active)

	// Set short URL if available
	if req.ShortURL != "" {
		qrBuilder.SetShortURL(req.ShortURL)
	}

	// Set optional fields
	if req.Description != "" {
		qrBuilder.SetDescription(req.Description)
	}
	if req.RedirectURL != "" {
		qrBuilder.SetRedirectURL(req.RedirectURL)
	}
	if req.ExpiresAt != nil {
		qrBuilder.SetExpiresAt(*req.ExpiresAt)
	}
	if len(req.Tags) > 0 {
		qrBuilder.SetTags(req.Tags)
	}
	if len(req.Design) > 0 {
		qrBuilder.SetDesign(req.Design)
	}
	if req.GroupID != nil {
		qrBuilder.SetGroupID(*req.GroupID)
	}

	qr, err := qrBuilder.Save(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create QR code"})
	}

	return c.Status(fiber.StatusCreated).JSON(qr)
}

// GetQRCode retrieves a QR code
func GetQRCode(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid QR code ID"})
	}

	qr, err := database.DB.QRCode.
		Query().
		Where(qrcode.IDEQ(id)).
		WithGroup().
		WithFileRefs().
		WithAnalyticsRecords().
		Only(context.Background())

	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "QR code not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve QR code"})
	}

	return c.JSON(qr)
}

// UpdateQRCode updates a QR code
func UpdateQRCode(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid QR code ID"})
	}

	var req struct {
		Type        string                 `json:"type"`
		Title       string                 `json:"title"`
		Description string                 `json:"description,omitempty"`
		RedirectURL string                 `json:"redirect_url,omitempty"`
		ShortURL    string                 `json:"short_url,omitempty"`
		Content     map[string]interface{} `json:"content"`
		ExpiresAt   *time.Time             `json:"expires_at,omitempty"`
		Analytics   bool                   `json:"analytics"`
		Active      bool                   `json:"active"`
		Tags        []string               `json:"tags,omitempty"`
		Design      map[string]interface{} `json:"design,omitempty"`
		GroupID     *int                   `json:"group_id,omitempty"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Fetch the existing QR code to preserve short_url if not provided
	existingQR, err := database.DB.QRCode.Get(context.Background(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "QR code not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve QR code"})
	}

	updateBuilder := database.DB.QRCode.UpdateOneID(id).
		SetType(req.Type).
		SetTitle(req.Title).
		SetContent(req.Content).
		SetAnalytics(req.Analytics).
		SetActive(req.Active).
		SetUpdatedAt(time.Now())

	// Always set short_url to the existing value if not provided
	shortURLToSet := req.ShortURL
	if shortURLToSet == "" {
		shortURLToSet = existingQR.ShortURL
	}
	updateBuilder.SetShortURL(shortURLToSet)

	// Set optional fields
	if req.Description != "" {
		updateBuilder.SetDescription(req.Description)
	} else {
		updateBuilder.ClearDescription()
	}
	if req.RedirectURL != "" {
		updateBuilder.SetRedirectURL(req.RedirectURL)
	} else {
		updateBuilder.ClearRedirectURL()
	}
	if req.ExpiresAt != nil {
		updateBuilder.SetExpiresAt(*req.ExpiresAt)
	} else {
		updateBuilder.ClearExpiresAt()
	}
	if len(req.Tags) > 0 {
		updateBuilder.SetTags(req.Tags)
	} else {
		updateBuilder.ClearTags()
	}
	if len(req.Design) > 0 {
		updateBuilder.SetDesign(req.Design)
	} else {
		updateBuilder.ClearDesign()
	}
	if req.GroupID != nil {
		updateBuilder.SetGroupID(*req.GroupID)
	} else {
		updateBuilder.ClearGroupID()
	}

	qr, err := updateBuilder.Save(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "QR code not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update QR code"})
	}

	// Always include short_url in the response
	resp := map[string]interface{}{
		"id":         qr.ID,
		"type":       qr.Type,
		"title":      qr.Title,
		"short_url":  qr.ShortURL,
		"content":    qr.Content,
		"created_at": qr.CreatedAt,
		"updated_at": qr.UpdatedAt,
		"analytics":  qr.Analytics,
		"active":     qr.Active,
		"edges":      qr.Edges,
	}
	return c.JSON(resp)
}

// DeleteQRCode deletes a QR code
func DeleteQRCode(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid QR code ID"})
	}

	err = database.DB.QRCode.DeleteOneID(id).Exec(context.Background())
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "QR code not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete QR code"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// BulkDeleteQRCodes deletes all QR codes
func BulkDeleteQRCodes(c *fiber.Ctx) error {
	_, err := database.DB.QRCode.Delete().Exec(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete all QR codes"})
	}
	return c.JSON(fiber.Map{"message": "All QR codes deleted successfully"})
}

// ListQRCodes retrieves all QR codes with pagination
func ListQRCodes(c *fiber.Ctx) error {
	page := c.QueryInt("page", 1)
	limit := c.QueryInt("limit", 20)

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 20
	}

	offset := (page - 1) * limit

	qrs, err := database.DB.QRCode.
		Query().
		WithGroup().
		WithFileRefs().
		Limit(limit).
		Offset(offset).
		Order(ent.Desc(qrcode.FieldCreatedAt)).
		All(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve QR codes"})
	}

	// Get total count for pagination
	total, err := database.DB.QRCode.Query().Count(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to count QR codes"})
	}

	return c.JSON(fiber.Map{
		"data": qrs,
		"pagination": fiber.Map{
			"page":  page,
			"limit": limit,
			"total": total,
		},
	})
}

// DownloadQRCode generates and downloads QR code image
func DownloadQRCode(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid QR code ID"})
	}

	// Get QR code from database
	qr, err := database.DB.QRCode.Get(context.Background(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "QR code not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve QR code"})
	}

	// Check if QR code is active and not expired
	if !qr.Active {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "QR code is inactive"})
	}
	if qr.ExpiresAt != nil && qr.ExpiresAt.Before(time.Now()) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "QR code has expired"})
	}

	// Get query parameters for customization
	format := c.Query("format", "png")
	size := c.QueryInt("size", 256)
	level := c.Query("level", "medium")

	// Validate parameters
	if size < 64 || size > 1024 {
		size = 256
	}

	// Convert level string to goqrcode.RecoveryLevel
	var recoveryLevel goqrcode.RecoveryLevel
	switch strings.ToLower(level) {
	case "low":
		recoveryLevel = goqrcode.Low
	case "medium":
		recoveryLevel = goqrcode.Medium
	case "high":
		recoveryLevel = goqrcode.High
	case "highest":
		recoveryLevel = goqrcode.Highest
	default:
		recoveryLevel = goqrcode.Medium
	}

	// Determine what data to encode
	var dataToEncode string
	if qr.Type == "static" {
		// Encode the info directly for static QR codes
		if text, ok := qr.Content["text"].(string); ok {
			// Plain text content for static text QR codes
			dataToEncode = text
		} else if urlStr, ok := qr.Content["url"].(string); ok {
			// Website URL for static website QR codes
			dataToEncode = urlStr
		} else if eventName, ok := qr.Content["name"].(string); ok && qr.Content["date"] != nil && qr.Content["end_date"] != nil {
			// iCalendar (VEVENT) support for static event QR codes
			start, _ := qr.Content["date"].(string)
			end, _ := qr.Content["end_date"].(string)
			location, _ := qr.Content["location"].(string)
			description, _ := qr.Content["description"].(string)
			// Format dates to YYYYMMDDTHHMMSS (strict 15 chars)
			formatICal := func(dt string) string {
				dt = strings.ReplaceAll(dt, "-", "")
				dt = strings.ReplaceAll(dt, ":", "")
				dt = strings.ReplaceAll(dt, "T", "T")
				if len(dt) >= 15 {
					return dt[:15]
				}
				// Pad with zeros if needed
				return dt + strings.Repeat("0", 15-len(dt))
			}
			crlf := "\r\n"
			ical := "BEGIN:VCALENDAR" + crlf
			ical += "VERSION:2.0" + crlf
			ical += "BEGIN:VEVENT" + crlf
			ical += fmt.Sprintf("SUMMARY:%s%s", eventName, crlf)
			if start != "" {
				ical += fmt.Sprintf("DTSTART:%s%s", formatICal(start), crlf)
			}
			if end != "" {
				ical += fmt.Sprintf("DTEND:%s%s", formatICal(end), crlf)
			}
			if location != "" {
				ical += fmt.Sprintf("LOCATION:%s%s", location, crlf)
			}
			if description != "" {
				ical += fmt.Sprintf("DESCRIPTION:%s%s", description, crlf)
			}
			ical += "END:VEVENT" + crlf + "END:VCALENDAR"
			dataToEncode = ical
		} else if ssid, ok := qr.Content["ssid"].(string); ok {
			password, _ := qr.Content["password"].(string)
			encryption, _ := qr.Content["encryption"].(string)
			hidden, _ := qr.Content["hidden"].(bool)
			hiddenParam := ""
			if hidden {
				hiddenParam = "H:true;"
			}
			dataToEncode = "WIFI:T:" + encryption + ";S:" + ssid + ";P:" + password + ";" + hiddenParam + ";"
		} else if phone, ok := qr.Content["phone_number"].(string); ok {
			message, _ := qr.Content["message"].(string)
			dataToEncode = "SMSTO:" + phone + ":" + message
		} else if recipient, ok := qr.Content["recipient"].(string); ok {
			subject, _ := qr.Content["subject"].(string)
			body, _ := qr.Content["body"].(string)
			dataToEncode = fmt.Sprintf(
				"mailto:%s?subject=%s&body=%s",
				url.QueryEscape(recipient),
				url.QueryEscape(subject),
				url.QueryEscape(body),
			)
		} else if name, ok := qr.Content["name"].(string); ok {
			// vCard support for static QR codes
			organization, _ := qr.Content["organization"].(string)
			title, _ := qr.Content["title"].(string)
			phone, _ := qr.Content["phone"].(string)
			email, _ := qr.Content["email"].(string)
			address, _ := qr.Content["address"].(string)
			// vCard 3.0 format
			vcard := "BEGIN:VCARD\nVERSION:3.0\n"
			vcard += fmt.Sprintf("FN:%s\n", name)
			if organization != "" {
				vcard += fmt.Sprintf("ORG:%s\n", organization)
			}
			if title != "" {
				vcard += fmt.Sprintf("TITLE:%s\n", title)
			}
			if phone != "" {
				vcard += fmt.Sprintf("TEL;TYPE=CELL:%s\n", phone)
			}
			if email != "" {
				vcard += fmt.Sprintf("EMAIL:%s\n", email)
			}
			if address != "" {
				vcard += fmt.Sprintf("ADR;TYPE=WORK:;;%s\n", address)
			}
			vcard += "END:VCARD"
			dataToEncode = vcard
		} else {
			// fallback for other static types
			if urlStr, ok := qr.Content["url"].(string); ok {
				dataToEncode = urlStr
			} else {
				dataToEncode = fmt.Sprintf("%s/qr/%d", c.BaseURL(), qr.ID)
			}
		}
	} else if (qr.Type == "dynamic" || qr.Type == "app" || qr.Type == "business") && qr.ShortURL != "" {
		// For dynamic QR codes (including app and business), encode only the short URL
		dataToEncode = fmt.Sprintf("%s/scan/%s", c.BaseURL(), qr.ShortURL)
	} else if qr.RedirectURL != "" {
		dataToEncode = qr.RedirectURL
	} else if qr.ShortURL != "" {
		dataToEncode = fmt.Sprintf("%s/scan/%s", c.BaseURL(), qr.ShortURL)
	} else {
		if urlStr, ok := qr.Content["url"].(string); ok {
			dataToEncode = urlStr
		} else {
			dataToEncode = fmt.Sprintf("%s/qr/%d", c.BaseURL(), qr.ID)
		}
	}

	// Generate QR code
	var imageData []byte
	if format == "svg" {
		// For SVG, we'll use a simple approach (you might want to use a different library)
		imageData, err = qrgen.GenerateWithLevel(dataToEncode, recoveryLevel, size)
	} else {
		imageData, err = qrgen.GenerateWithLevel(dataToEncode, recoveryLevel, size)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate QR code"})
	}

	// Set appropriate headers
	filename := fmt.Sprintf("%s.%s", qr.Title, format)
	c.Set("Content-Type", fmt.Sprintf("image/%s", format))
	c.Set("Content-Disposition", fmt.Sprintf("attachment; filename=\"%s\"", filename))

	return c.Send(imageData)
}

// ScanQRCode handles QR code scanning and redirection or static content display
func ScanQRCode(c *fiber.Ctx) error {
	shortCode := c.Params("shortcode")
	if shortCode == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Short code is required"})
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	qr, err := database.DB.QRCode.
		Query().
		Where(qrcode.ShortURLEQ(shortCode)).
		Only(ctx)

	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "QR code not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve QR code"})
	}

	if !qr.Active {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "QR code is inactive"})
	}
	if qr.ExpiresAt != nil && qr.ExpiresAt.Before(time.Now()) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "QR code has expired"})
	}

	// Track analytics for all QR codes that have analytics enabled
	if qr.Analytics {
		ipAddress := c.IP()
		userAgent := c.Get("User-Agent")
		go func(id int, ip, ua string) {
			goCtx, goCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer goCancel()
			qr, err := database.DB.QRCode.Get(goCtx, id)
			if err == nil {
				_, _ = database.DB.QRCodeAnalytics.Create().
					SetIPAddress(ip).
					SetUserAgent(ua).
					SetScannedAt(time.Now()).
					SetQrCode(qr).
					Save(goCtx)
			}
		}(qr.ID, ipAddress, userAgent)
	}

	// App QR code landing page
	if qr.Type == "app" {
		appStoreURL, _ := qr.Content["app_store_url"].(string)
		deepLink, _ := qr.Content["deep_link"].(string)
		appName, _ := qr.Content["name"].(string)
		if appName == "" {
			appName = "Mobile App"
		}
		
		data := fiber.Map{
			"AppName":     appName,
			"AppStoreURL": appStoreURL,
			"DeepLink":    deepLink,
			"Title":       "Download App",
		}
		return c.Render("app", data)
	}

	// Business QR code landing page
	if qr.Type == "business" {
		businessName, _ := qr.Content["name"].(string)
		tagline, _ := qr.Content["tagline"].(string)
		website, _ := qr.Content["website"].(string)
		description, _ := qr.Content["description"].(string)
		logoURL, _ := qr.Content["logo_url"].(string)
		contactInfo, _ := qr.Content["contact_info"].(map[string]interface{})
		socialLinks, _ := qr.Content["social_links"].(map[string]interface{})
		
		if businessName == "" {
			businessName = "Business"
		}
		
		data := fiber.Map{
			"BusinessName": businessName,
			"Tagline":      tagline,
			"Website":      website,
			"Description":  description,
			"LogoURL":      logoURL,
			"ContactInfo":  contactInfo,
			"SocialLinks":  socialLinks,
			"Title":        businessName,
		}
		return c.Render("business", data)
	}

	// WiFi QR code landing page
	if qr.Type == "wifi" {
		ssid, _ := qr.Content["ssid"].(string)
		password, _ := qr.Content["password"].(string)
		encryption, ok := qr.Content["encryption"].(string)
		if !ok || encryption == "" {
			encryption = "WPA"
		}
		hidden, _ := qr.Content["hidden"].(bool)
		hiddenParam := ""
		if hidden {
			hiddenParam = "H:true;"
		}
		wifiURI := "WIFI:T:" + encryption + ";S:" + ssid + ";P:" + password + ";" + hiddenParam + ";"
		data := fiber.Map{
			"SSID":        ssid,
			"Password":    password,
			"Encryption":  encryption,
			"WiFiURI":     wifiURI,
			"DownloadURL": "data:text/plain;charset=utf-8," + url.QueryEscape(wifiURI),
			"Title":       "Connect to WiFi",
		}
		return c.Render("wifi", data)
	}

	// SMS QR code landing page
	if phone, ok := qr.Content["phone_number"].(string); ok {
		message, _ := qr.Content["message"].(string)
		smsURI := "SMSTO:" + phone + ":" + message
		data := fiber.Map{
			"Phone":   phone,
			"Message": message,
			"SMSURI":  smsURI,
			"Title":   "Send SMS",
		}
		return c.Render("sms", data)
	}

	// Email QR code landing page
	if recipient, ok := qr.Content["recipient"].(string); ok {
		subject, _ := qr.Content["subject"].(string)
		body, _ := qr.Content["body"].(string)
		mailto := "mailto:" + recipient + "?subject=" + url.QueryEscape(subject) + "&body=" + url.QueryEscape(body)
		data := fiber.Map{
			"Recipient": recipient,
			"Subject":   subject,
			"Body":      body,
			"Mailto":    mailto,
			"Title":     "Send Email",
		}
		return c.Render("email", data)
	}

	// vCard QR code landing page
	if name, ok := qr.Content["name"].(string); ok && (qr.Content["organization"] != nil || qr.Content["title"] != nil || qr.Content["phone"] != nil || qr.Content["email"] != nil || qr.Content["address"] != nil) {
		organization, _ := qr.Content["organization"].(string)
		title, _ := qr.Content["title"].(string)
		phone, _ := qr.Content["phone"].(string)
		emailAddr, _ := qr.Content["email"].(string)
		address, _ := qr.Content["address"].(string)

		// Generate vCard content
		vcard := "BEGIN:VCARD\nVERSION:3.0\n"
		vcard += "FN:" + name + "\n"
		if organization != "" {
			vcard += "ORG:" + organization + "\n"
		}
		if title != "" {
			vcard += "TITLE:" + title + "\n"
		}
		if phone != "" {
			vcard += "TEL:" + phone + "\n"
		}
		if emailAddr != "" {
			vcard += "EMAIL:" + emailAddr + "\n"
		}
		if address != "" {
			vcard += "ADR:" + address + "\n"
		}
		vcard += "END:VCARD"

		// Create data URL for vCard download with template.URL to mark as safe
		vcardURL := template.URL("data:text/x-vcard;charset=utf-8," + url.QueryEscape(vcard))

		data := fiber.Map{
			"Name":         name,
			"Organization": organization,
			"Title":        title,
			"Phone":        phone,
			"Email":        emailAddr,
			"Address":      address,
			"TitlePage":    "Save Contact",
			"VCardURL":     vcardURL,
		}
		return c.Render("vcard", data)
	}

	// Event QR code landing page
	if eventName, ok := qr.Content["name"].(string); ok && qr.Content["date"] != nil && qr.Content["end_date"] != nil {
		start, _ := qr.Content["date"].(string)
		end, _ := qr.Content["end_date"].(string)
		location, _ := qr.Content["location"].(string)
		description, _ := qr.Content["description"].(string)

		// Generate iCal content with proper date formatting
		ical := "BEGIN:VCALENDAR\nVERSION:2.0\nPRODID:-//QR Platform//Event//EN\n"
		ical += "BEGIN:VEVENT\n"
		ical += "SUMMARY:" + eventName + "\n"

		// Format dates for iCal (remove hyphens and add time)
		startFormatted := strings.ReplaceAll(strings.ReplaceAll(start, "-", ""), ":", "")
		endFormatted := strings.ReplaceAll(strings.ReplaceAll(end, "-", ""), ":", "")

		// If dates don't include time, add default times
		if len(startFormatted) == 8 {
			startFormatted += "T090000Z"
		} else if !strings.Contains(startFormatted, "T") {
			startFormatted += "T090000Z"
		}

		if len(endFormatted) == 8 {
			endFormatted += "T170000Z"
		} else if !strings.Contains(endFormatted, "T") {
			endFormatted += "T170000Z"
		}

		ical += "DTSTART:" + startFormatted + "\n"
		ical += "DTEND:" + endFormatted + "\n"

		if location != "" {
			ical += "LOCATION:" + location + "\n"
		}
		if description != "" {
			ical += "DESCRIPTION:" + description + "\n"
		}
		ical += "END:VEVENT\n"
		ical += "END:VCALENDAR"

		// Create data URL for iCal download with template.URL to mark as safe
		icalURL := template.URL("data:text/calendar;charset=utf-8," + url.QueryEscape(ical))

		data := fiber.Map{
			"Event":       eventName,
			"Start":       start,
			"End":         end,
			"Location":    location,
			"Description": description,
			"Title":       "Add Event",
			"ICalURL":     icalURL,
		}
		return c.Render("event", data)
	}

	// PDF QR code landing page
	if qr.Content["type"] == "pdf" || (qr.Content["url"] != nil && strings.HasSuffix(strings.ToLower(qr.Content["url"].(string)), ".pdf")) {
		fileURL, _ := qr.Content["url"].(string)
		filename := "document.pdf"
		if storedFilename, ok := qr.Content["filename"].(string); ok && storedFilename != "" {
			filename = storedFilename
		} else if parts := strings.Split(fileURL, "/"); len(parts) > 0 {
			filename = parts[len(parts)-1]
		}
		data := fiber.Map{
			"FileURL":  fileURL,
			"Filename": filename,
			"Title":    "PDF Document",
		}
		return c.Render("pdf", data)
	}

	// Image QR code landing page
	if qr.Content["type"] == "image" || (qr.Content["url"] != nil && (strings.HasSuffix(strings.ToLower(qr.Content["url"].(string)), ".jpg") || strings.HasSuffix(strings.ToLower(qr.Content["url"].(string)), ".jpeg") || strings.HasSuffix(strings.ToLower(qr.Content["url"].(string)), ".png") || strings.HasSuffix(strings.ToLower(qr.Content["url"].(string)), ".gif"))) {
		fileURL, _ := qr.Content["url"].(string)
		filename := "image"
		if storedFilename, ok := qr.Content["filename"].(string); ok && storedFilename != "" {
			filename = storedFilename
		}
		data := fiber.Map{
			"FileURL":  fileURL,
			"Filename": filename,
			"Title":    "Image File",
		}
		return c.Render("image", data)
	}

	// Data Matrix Barcode QR code landing page
	if qr.Content["type"] == "barcode_2d" || qr.Content["data"] != nil {
		textData, _ := qr.Content["data"].(string)
		if textData == "" {
			textData = "No data available"
		}

		// Get the file URL from content
		fileURL, _ := qr.Content["url"].(string)
		filename := "datamatrix_barcode.png"
		if storedFilename, ok := qr.Content["filename"].(string); ok && storedFilename != "" {
			filename = storedFilename
		}

		data := fiber.Map{
			"TextData":    textData,
			"BarcodeType": "Data Matrix",
			"Size":        "200x200",
			"FileURL":     fileURL,
			"Filename":    filename,
			"Title":       "Data Matrix Barcode",
		}
		return c.Render("barcode", data)
	}

	// Handle dynamic website QR codes with proper redirection
	if qr.Type == "dynamic" {
		// Check if there's a redirect URL set
		if qr.RedirectURL != "" {
			return c.Redirect(qr.RedirectURL, fiber.StatusFound)
		}
		// Fallback to content URL for dynamic QR codes
		if urlStr, ok := qr.Content["url"].(string); ok {
			return c.Redirect(urlStr, fiber.StatusFound)
		}
	}

	// fallback for other dynamic types
	if urlStr, ok := qr.Content["url"].(string); ok {
		return c.Redirect(urlStr, fiber.StatusFound)
	} else {
		return c.Status(fiber.StatusOK).JSON(qr.Content)
	}
}

// GetStaticQRContent handles displaying static QR code content
func GetStaticQRContent(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid QR code ID"})
	}

	// Get QR code from database
	qr, err := database.DB.QRCode.Get(context.Background(), id)
	if err != nil {
		if ent.IsNotFound(err) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "QR code not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve QR code"})
	}

	// Check if QR code is active and not expired
	if !qr.Active {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "QR code is inactive"})
	}
	if qr.ExpiresAt != nil && qr.ExpiresAt.Before(time.Now()) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "QR code has expired"})
	}

	// Handle WiFi QR code specially
	if qr.Type == "wifi" {
		// Extract WiFi details with proper error handling
		ssid, ok := qr.Content["ssid"].(string)
		if !ok || ssid == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid or missing SSID in WiFi QR code"})
		}

		password, _ := qr.Content["password"].(string)
		encryption, ok := qr.Content["encryption"].(string)
		if !ok || encryption == "" {
			encryption = "WPA"
		}
		hidden, _ := qr.Content["hidden"].(bool)

		// Build WiFi URI string
		hiddenParam := ""
		if hidden {
			hiddenParam = "H:true;"
		}
		wifiURI := "WIFI:T:" + encryption + ";S:" + ssid + ";P:" + password + ";" + hiddenParam + ";"
		return c.Status(fiber.StatusOK).Type("text/plain").SendString(wifiURI)
	}

	// For other static QR codes, return their content
	return c.JSON(fiber.Map{"content": qr.Content})
}

// UploadFile handles file uploads for QR codes
func UploadFile(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	// Validate file size (10MB limit)
	if file.Size > 10*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File size exceeds 10MB limit"})
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg":      true,
		"image/png":       true,
		"image/gif":       true,
		"application/pdf": true,
		"text/plain":      true,
	}

	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File type not allowed"})
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create upload directory"})
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := filepath.Join(uploadsDir, filename)

	// Save file
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// Create file reference in database
	fileRef, err := database.DB.FileReference.Create().
		SetFilename(file.Filename).
		SetURL(fmt.Sprintf("/uploads/%s", filename)).
		SetSize(file.Size).
		SetType(contentType).
		Save(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file reference"})
	}

	return c.Status(fiber.StatusCreated).JSON(fileRef)
}

// CreatePDFQRCode handles PDF file upload and creates a QR code for it
func CreatePDFQRCode(c *fiber.Ctx) error {
	// Parse multipart form
	form, err := c.MultipartForm()
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Failed to parse multipart form"})
	}

	// Get the PDF file
	files := form.File["file"]
	if len(files) == 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No PDF file uploaded"})
	}

	file := files[0]

	// Validate file type
	if file.Header.Get("Content-Type") != "application/pdf" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Only PDF files are allowed"})
	}

	// Validate file size (10MB limit)
	if file.Size > 10*1024*1024 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File size exceeds 10MB limit"})
	}

	// Get form values
	title := c.FormValue("title")
	description := c.FormValue("description")
	isDynamic := c.FormValue("is_dynamic") == "true"
	analytics := c.FormValue("analytics") == "true"

	if title == "" {
		title = "PDF QR Code - " + file.Filename
	}

	// Create uploads directory if it doesn't exist
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create upload directory"})
	}

	// Generate unique filename
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := filepath.Join(uploadsDir, filename)

	// Save file
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save PDF file"})
	}

	// Create file reference in database
	fileRef, err := database.DB.FileReference.Create().
		SetFilename(file.Filename).
		SetURL(fmt.Sprintf("/uploads/%s", filename)).
		SetSize(file.Size).
		SetType("application/pdf").
		Save(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file reference"})
	}

	// Determine QR code type and generate short URL if needed
	var qrType string
	var shortURL string

	if isDynamic {
		qrType = "dynamic"
		shortURL, err = shorturl.Generate()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate short URL"})
		}
		analytics = true // Force analytics for dynamic QR codes
	} else {
		qrType = "static"
		if analytics {
			shortURL, err = shorturl.Generate()
			if err != nil {
				return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate short URL"})
			}
		}
	}

	// Create the PDF URL (accessible via your server)
	pdfURL := fmt.Sprintf("/uploads/%s", filename)

	// Create QR code content
	content := map[string]interface{}{
		"type":        "pdf",
		"url":         pdfURL,
		"filename":    file.Filename,
		"file_size":   file.Size,
		"file_ref_id": fileRef.ID,
	}

	// Create QR code using Ent
	qrBuilder := database.DB.QRCode.Create().
		SetType(qrType).
		SetTitle(title).
		SetContent(content).
		SetAnalytics(analytics).
		SetActive(true).
		AddFileRefs(fileRef)

	// Set optional fields
	if description != "" {
		qrBuilder.SetDescription(description)
	}
	if shortURL != "" {
		qrBuilder.SetShortURL(shortURL)
	}
	if isDynamic {
		qrBuilder.SetRedirectURL(pdfURL)
	}

	qr, err := qrBuilder.Save(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create PDF QR code"})
	}

	// Return the created QR code with file reference
	response := fiber.Map{
		"qr_code":        qr,
		"file_reference": fileRef,
		"pdf_url":        pdfURL,
		"message":        "PDF QR code created successfully",
	}

	if shortURL != "" {
		response["short_url"] = shortURL
		response["scan_url"] = fmt.Sprintf("http://localhost:3000/scan/%s", shortURL)
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// CreateImageQRCode handles image file upload and creates a QR code for it
func CreateImageQRCode(c *fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No file uploaded"})
	}

	// Validate file type
	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}
	contentType := file.Header.Get("Content-Type")
	if !allowedTypes[contentType] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "File type not allowed"})
	}

	// Save file
	uploadsDir := "uploads"
	os.MkdirAll(uploadsDir, 0755)
	filename := fmt.Sprintf("%d_%s", time.Now().Unix(), file.Filename)
	filePath := filepath.Join(uploadsDir, filename)
	if err := c.SaveFile(file, filePath); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file"})
	}

	// Create file reference in DB
	fileRef, err := database.DB.FileReference.Create().
		SetFilename(file.Filename).
		SetURL(fmt.Sprintf("/uploads/%s", filename)).
		SetSize(file.Size).
		SetType(contentType).
		Save(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save file reference"})
	}

	// Generate short URL for dynamic QR
	shortURL, err := shorturl.Generate()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to generate short URL"})
	}

	// Create QR code content
	imageURL := fmt.Sprintf("/uploads/%s", filename)
	content := map[string]interface{}{
		"type":        "image",
		"url":         imageURL,
		"filename":    file.Filename,
		"file_ref_id": fileRef.ID,
	}

	// Create QR code in DB
	qr, err := database.DB.QRCode.Create().
		SetType("image").
		SetTitle("Image QR Code - " + file.Filename).
		SetContent(content).
		SetShortURL(shortURL).
		SetAnalytics(true).
		SetActive(true).
		AddFileRefs(fileRef).
		Save(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create image QR code"})
	}

	scanURL := fmt.Sprintf("http://localhost:3000/scan/%s", shortURL)

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"qr_code":        qr,
		"file_reference": fileRef,
		"image_url":      imageURL,
		"short_url":      shortURL,
		"scan_url":       scanURL,
		"message":        "Image QR code created successfully",
	})
}

// CreateBarcodeQRCode handles creating a QR code for Data Matrix barcodes
func CreateBarcodeQRCode(c *fiber.Ctx) error {
	var req struct {
		Title       string `json:"title"`
		Description string `json:"description,omitempty"`
		Data        string `json:"data"`
		IsDynamic   bool   `json:"is_dynamic"`
		Analytics   bool   `json:"analytics"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	if req.Data == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data is required for barcode generation"})
	}

	if req.Title == "" {
		req.Title = "Data Matrix Barcode"
	}

	// Generate the barcode image
	barcodeData, err := barcode.GenerateDataMatrix(barcode.DataMatrixOptions{
		Data: req.Data,
		Size: "auto",
	})
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to generate barcode: " + err.Error(),
		})
	}

	// Ensure uploads directory exists
	uploadsDir := "uploads"
	if err := os.MkdirAll(uploadsDir, 0755); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create uploads directory: " + err.Error(),
		})
	}

	// Create unique filename
	filename := fmt.Sprintf("%d_datamatrix_barcode.png", time.Now().Unix())
	filePath := filepath.Join(uploadsDir, filename)

	// Save barcode file
	if err := os.WriteFile(filePath, barcodeData, 0644); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save barcode file: " + err.Error(),
		})
	}

	// Create file reference in database
	fileRef, err := database.DB.FileReference.Create().
		SetFilename("datamatrix_barcode.png").
		SetURL(fmt.Sprintf("/uploads/%s", filename)).
		SetSize(int64(len(barcodeData))).
		SetType("image/png").
		Save(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to save file reference: " + err.Error(),
		})
	}

	// Generate short URL for analytics/dynamic QR
	var shortURL string
	var qrType string

	if req.IsDynamic {
		qrType = "dynamic"
		req.Analytics = true
	} else {
		qrType = "static"
	}

	if req.Analytics || req.IsDynamic {
		shortURL, err = shorturl.Generate()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "Failed to generate short URL: " + err.Error(),
			})
		}
	}

	// Create QR code content
	barcodeURL := fmt.Sprintf("/uploads/%s", filename)
	content := map[string]interface{}{
		"type":        "barcode_2d",
		"data":        req.Data,
		"url":         barcodeURL,
		"filename":    "datamatrix_barcode.png",
		"file_ref_id": fileRef.ID,
	}

	// Create QR code in database
	qrBuilder := database.DB.QRCode.Create().
		SetType(qrType).
		SetTitle(req.Title).
		SetContent(content).
		SetAnalytics(req.Analytics).
		SetActive(true).
		AddFileRefs(fileRef)

	if req.Description != "" {
		qrBuilder.SetDescription(req.Description)
	}
	if shortURL != "" {
		qrBuilder.SetShortURL(shortURL)
	}

	qr, err := qrBuilder.Save(context.Background())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create QR code: " + err.Error(),
		})
	}

	// Build response
	response := fiber.Map{
		"qr_code":        qr,
		"file_reference": fileRef,
		"barcode_url":    barcodeURL,
		"message":        "Data Matrix barcode QR code created successfully",
	}

	if shortURL != "" {
		response["short_url"] = shortURL
		response["scan_url"] = fmt.Sprintf("http://localhost:3000/scan/%s", shortURL)
	}

	return c.Status(fiber.StatusCreated).JSON(response)
}

// GetQRCodeAnalytics retrieves analytics for a QR code
func GetQRCodeAnalytics(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid QR code ID"})
	}

	// Get analytics records
	analytics, err := database.DB.QRCodeAnalytics.
		Query().
		Where(qrcodeanalytics.HasQrCodeWith(qrcode.IDEQ(id))).
		Order(ent.Desc(qrcodeanalytics.FieldScannedAt)).
		All(context.Background())

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve analytics"})
	}

	// Calculate summary statistics
	totalScans := len(analytics)
	uniqueIPs := make(map[string]bool)
	for _, record := range analytics {
		uniqueIPs[record.IPAddress] = true
	}

	return c.JSON(fiber.Map{
		"total_scans":     totalScans,
		"unique_visitors": len(uniqueIPs),
		"records":         analytics,
	})
}
