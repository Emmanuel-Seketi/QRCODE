package handler

import (
	"context"
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
		if eventName, ok := qr.Content["name"].(string); ok && qr.Content["date"] != nil && qr.Content["end_date"] != nil {
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
			if url, ok := qr.Content["url"].(string); ok {
				dataToEncode = url
			} else {
				dataToEncode = fmt.Sprintf("%s/qr/%d", c.BaseURL(), qr.ID)
			}
		}
	} else if qr.Type == "dynamic" && qr.ShortURL != "" {
		// For dynamic QR codes, encode only the short URL
		dataToEncode = fmt.Sprintf("%s/scan/%s", c.BaseURL(), qr.ShortURL)
	} else if qr.RedirectURL != "" {
		dataToEncode = qr.RedirectURL
	} else if qr.ShortURL != "" {
		dataToEncode = fmt.Sprintf("%s/scan/%s", c.BaseURL(), qr.ShortURL)
	} else {
		if url, ok := qr.Content["url"].(string); ok {
			dataToEncode = url
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

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find QR code by short URL
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

	// Check if QR code is active and not expired
	if !qr.Active {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "QR code is inactive"})
	}
	if qr.ExpiresAt != nil && qr.ExpiresAt.Before(time.Now()) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "QR code has expired"})
	}

	// Track analytics only for dynamic QR codes
	if qr.Analytics && qr.Type == "dynamic" {
		// Get client information
		ipAddress := c.IP()
		userAgent := c.Get("User-Agent")

		// Create a new goroutine with a separate context and copied data
		go func(id int, ip, ua string) {
			// Create a new context for the goroutine
			goCtx, goCancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer goCancel()

			// Get QR code first with error handling
			qr, err := database.DB.QRCode.Get(goCtx, id)
			if err != nil {
				fmt.Printf("Failed to get QR code for analytics: %v\n", err)
				return
			}

			// Create analytics record
			_, err = database.DB.QRCodeAnalytics.Create().
				SetIPAddress(ip).
				SetUserAgent(ua).
				SetScannedAt(time.Now()).
				SetQrCode(qr).
				Save(goCtx)

			if err != nil {
				// Log error but don't fail the request
				fmt.Printf("Failed to track scan: %v\n", err)
			}
		}(qr.ID, ipAddress, userAgent)
	}

	// Remove plain text response for dynamic WiFi QR codes and render a landing page
	if qr.Type == "dynamic" {
		if ssid, ok := qr.Content["ssid"].(string); ok {
			password, _ := qr.Content["password"].(string)
			encryption, _ := qr.Content["encryption"].(string)
			hidden, _ := qr.Content["hidden"].(bool)
			hiddenParam := ""
			if hidden {
				hiddenParam = "H:true;"
			}
			wifiURI := "WIFI:T:" + encryption + ";S:" + ssid + ";P:" + password + ";" + hiddenParam + ";"
			html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Connect to WiFi</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    :root { --primary1: #0c768a; --primary2: #0C8096; --primary3: #26666F; --bg1: #ffffff; --bg2: #fbfbfb; --bg3: #eef2f5; --text1: #424242; --text2: #000000; --border1: #d2d2d2; --border2: #d9d9d9; }
    body { background: var(--bg3); color: var(--text1); font-family: 'Segoe UI', Arial, sans-serif; margin: 0; padding: 0; min-height: 100vh; display: flex; align-items: center; justify-content: center; }
    .container { background: var(--bg1); border-radius: 16px; box-shadow: 0 4px 24px rgba(38, 102, 111, 0.08); padding: 2.5rem 1.5rem 2rem 1.5rem; max-width: 350px; width: 100%; border: 1px solid var(--border2); text-align: center; }
    h2 { color: var(--primary1); margin-bottom: 0.5rem; font-size: 1.6rem; font-weight: 700; }
    .wifi-info { background: var(--bg2); border: 1px solid var(--border1); border-radius: 10px; padding: 1rem; margin: 1.2rem 0 1.5rem 0; text-align: left; font-size: 1.05rem; word-break: break-all; }
    .wifi-info label { color: var(--primary3); font-weight: 600; margin-right: 0.5em; }
    .join-btn, .download-btn { display: block; width: 100%; background: linear-gradient(90deg, var(--primary1), var(--primary2)); color: var(--bg1); font-size: 1.15rem; font-weight: 600; border: none; border-radius: 8px; padding: 0.85rem 0; margin-bottom: 1.2rem; cursor: pointer; transition: background 0.2s; text-decoration: none; }
    .join-btn:hover, .join-btn:focus, .download-btn:hover, .download-btn:focus { background: var(--primary3); color: var(--bg1); }
    .note { font-size: 0.98rem; color: var(--text1); background: var(--bg3); border-radius: 6px; padding: 0.7em 1em; border: 1px solid var(--border1); margin-top: 0.5em; }
    @media (max-width: 480px) { .container { padding: 1.2rem 0.5rem 1.2rem 0.5rem; max-width: 98vw; } h2 { font-size: 1.2rem; } }
  </style>
</head>
<body>
  <div class="container">
    <h2>Connect to WiFi</h2>
    <div class="wifi-info">
      <div><label>Network:</label> <span id="ssid">` + ssid + `</span></div>
      <div><label>Password:</label> <span id="password">` + password + `</span></div>
      <div><label>Security:</label> <span id="security">` + encryption + `</span></div>
    </div>
    <a class="join-btn" id="joinBtn" href="` + wifiURI + `">Join WiFi</a>
    <a class="download-btn" id="downloadWifi" href="data:text/plain;charset=utf-8,` + url.QueryEscape(wifiURI) + `" download="wifi.wifi">Download WiFi File</a>
    <div class="note">
      <b>Tip:</b> On <b>Android</b>, tap "Join WiFi" to connect instantly.<br>
      On <b>iPhone</b>, use the camera to scan a static WiFi QR code for instant connection, or enter the credentials above manually.<br>
      You can also download the WiFi file and open it with a compatible app.
    </div>
  </div>
</body>
</html>`
			return c.Status(fiber.StatusOK).Type("html").SendString(html)
		} else if phone, ok := qr.Content["phone_number"].(string); ok {
			message, _ := qr.Content["message"].(string)
			smsURI := "SMSTO:" + phone + ":" + message
			html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Send SMS</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    :root { --primary1: #0c768a; --primary2: #0C8096; --primary3: #26666F; --bg1: #ffffff; --bg2: #fbfbfb; --bg3: #eef2f5; --text1: #424242; --text2: #000000; --border1: #d2d2d2; --border2: #d9d9d9; }
    body { background: var(--bg3); color: var(--text1); font-family: 'Segoe UI', Arial, sans-serif; margin: 0; padding: 0; min-height: 100vh; display: flex; align-items: center; justify-content: center; }
    .container { background: var(--bg1); border-radius: 16px; box-shadow: 0 4px 24px rgba(38, 102, 111, 0.08); padding: 2.5rem 1.5rem 2rem 1.5rem; max-width: 350px; width: 100%; border: 1px solid var(--border2); text-align: center; }
    h2 { color: var(--primary1); margin-bottom: 0.5rem; font-size: 1.6rem; font-weight: 700; }
    .sms-info { background: var(--bg2); border: 1px solid var(--border1); border-radius: 10px; padding: 1rem; margin: 1.2rem 0 1.5rem 0; text-align: left; font-size: 1.05rem; word-break: break-all; }
    .sms-info label { color: var(--primary3); font-weight: 600; margin-right: 0.5em; }
    .send-btn { display: block; width: 100%; background: linear-gradient(90deg, var(--primary1), var(--primary2)); color: var(--bg1); font-size: 1.15rem; font-weight: 600; border: none; border-radius: 8px; padding: 0.85rem 0; margin-bottom: 1.2rem; cursor: pointer; transition: background 0.2s; text-decoration: none; }
    .send-btn:hover, .send-btn:focus { background: var(--primary3); color: var(--bg1); }
    .note { font-size: 0.98rem; color: var(--text1); background: var(--bg3); border-radius: 6px; padding: 0.7em 1em; border: 1px solid var(--border1); margin-top: 0.5em; }
    @media (max-width: 480px) { .container { padding: 1.2rem 0.5rem 1.2rem 0.5rem; max-width: 98vw; } h2 { font-size: 1.2rem; } }
  </style>
</head>
<body>
  <div class="container">
    <h2>Send SMS</h2>
    <div class="sms-info">
      <div><label>To:</label> <span id="to">` + phone + `</span></div>
      <div><label>Message:</label> <span id="message">` + message + `</span></div>
    </div>
    <a class="send-btn" id="sendBtn" href="` + smsURI + `">Send SMS</a>
    <div class="note">Tap "Send SMS" to open your SMS app. If WhatsApp or another app opens, it may be set as your default for SMS links. This works only in supported browsers and on Android devices.</div>
  </div>
</body>
</html>`
			return c.Status(fiber.StatusOK).Type("html").SendString(html)
		} else if recipient, ok := qr.Content["recipient"].(string); ok {
			// Email landing page
			subject, _ := qr.Content["subject"].(string)
			body, _ := qr.Content["body"].(string)
			mailto := "mailto:" + recipient + "?subject=" + url.QueryEscape(subject) + "&body=" + url.QueryEscape(body)
			html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Send Email</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    :root { --primary1: #0c768a; --primary2: #0C8096; --primary3: #26666F; --bg1: #ffffff; --bg2: #fbfbfb; --bg3: #eef2f5; --text1: #424242; --text2: #000000; --border1: #d2d2d2; --border2: #d9d9d9; }
    body { background: var(--bg3); color: var(--text1); font-family: 'Segoe UI', Arial, sans-serif; margin: 0; padding: 0; min-height: 100vh; display: flex; align-items: center; justify-content: center; }
    .container { background: var(--bg1); border-radius: 16px; box-shadow: 0 4px 24px rgba(38, 102, 111, 0.08); padding: 2.5rem 1.5rem 2rem 1.5rem; max-width: 350px; width: 100%; border: 1px solid var(--border2); text-align: center; }
    h2 { color: var(--primary1); margin-bottom: 0.5rem; font-size: 1.6rem; font-weight: 700; }
    .email-info { background: var(--bg2); border: 1px solid var(--border1); border-radius: 10px; padding: 1rem; margin: 1.2rem 0 1.5rem 0; text-align: left; font-size: 1.05rem; word-break: break-all; }
    .email-info label { color: var(--primary3); font-weight: 600; margin-right: 0.5em; }
    .send-btn { display: block; width: 100%; background: linear-gradient(90deg, var(--primary1), var(--primary2)); color: var(--bg1); font-size: 1.15rem; font-weight: 600; border: none; border-radius: 8px; padding: 0.85rem 0; margin-bottom: 1.2rem; cursor: pointer; transition: background 0.2s; text-decoration: none; }
    .send-btn:hover, .send-btn:focus { background: var(--primary3); color: var(--bg1); }
    .note { font-size: 0.98rem; color: var(--text1); background: var(--bg3); border-radius: 6px; padding: 0.7em 1em; border: 1px solid var(--border1); margin-top: 0.5em; }
    @media (max-width: 480px) { .container { padding: 1.2rem 0.5rem 1.2rem 0.5rem; max-width: 98vw; } h2 { font-size: 1.2rem; } }
  </style>
</head>
<body>
  <div class="container">
    <h2>Send Email</h2>
    <div class="email-info">
      <div><label>To:</label> <span id="to">` + recipient + `</span></div>
      <div><label>Subject:</label> <span id="subject">` + subject + `</span></div>
      <div><label>Body:</label> <span id="body">` + body + `</span></div>
    </div>
    <a class="send-btn" id="sendBtn" href="` + mailto + `">Send Email</a>
    <div class="note">Tap "Send Email" to open your email app with the details filled in.</div>
  </div>
</body>
</html>`
			return c.Status(fiber.StatusOK).Type("html").SendString(html)
		} else if name, ok := qr.Content["name"].(string); ok && (qr.Content["organization"] != nil || qr.Content["title"] != nil || qr.Content["phone"] != nil || qr.Content["email"] != nil || qr.Content["address"] != nil) {
			// vCard landing page
			organization, _ := qr.Content["organization"].(string)
			title, _ := qr.Content["title"].(string)
			phone, _ := qr.Content["phone"].(string)
			emailAddr, _ := qr.Content["email"].(string)
			address, _ := qr.Content["address"].(string)
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
			if emailAddr != "" {
				vcard += "EMAIL:" + emailAddr + "\n"
			}
			if address != "" {
				vcard += "ADR;TYPE=WORK:;;" + address + "\n"
			}
			vcard += "END:VCARD"
			html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Save Contact</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    :root { --primary1: #0c768a; --primary2: #0C8096; --primary3: #26666F; --bg1: #ffffff; --bg2: #fbfbfb; --bg3: #eef2f5; --text1: #424242; --text2: #000000; --border1: #d2d2d2; --border2: #d9d9d9; }
    body { background: var(--bg3); color: var(--text1); font-family: 'Segoe UI', Arial, sans-serif; margin: 0; padding: 0; min-height: 100vh; display: flex; align-items: center; justify-content: center; }
    .container { background: var(--bg1); border-radius: 16px; box-shadow: 0 4px 24px rgba(38, 102, 111, 0.08); padding: 2.5rem 1.5rem 2rem 1.5rem; max-width: 350px; width: 100%; border: 1px solid var(--border2); text-align: center; }
    h2 { color: var(--primary1); margin-bottom: 0.5rem; font-size: 1.6rem; font-weight: 700; }
    .vcard-info { background: var(--bg2); border: 1px solid var(--border1); border-radius: 10px; padding: 1rem; margin: 1.2rem 0 1.5rem 0; text-align: left; font-size: 1.05rem; word-break: break-all; }
    .vcard-info label { color: var(--primary3); font-weight: 600; margin-right: 0.5em; }
    .save-btn { display: block; width: 100%; background: linear-gradient(90deg, var(--primary1), var(--primary2)); color: var(--bg1); font-size: 1.15rem; font-weight: 600; border: none; border-radius: 8px; padding: 0.85rem 0; margin-bottom: 1.2rem; cursor: pointer; transition: background 0.2s; text-decoration: none; }
    .save-btn:hover, .save-btn:focus { background: var(--primary3); color: var(--bg1); }
    .note { font-size: 0.98rem; color: var(--text1); background: var(--bg3); border-radius: 6px; padding: 0.7em 1em; border: 1px solid var(--border1); margin-top: 0.5em; }
    @media (max-width: 480px) { .container { padding: 1.2rem 0.5rem 1.2rem 0.5rem; max-width: 98vw; } h2 { font-size: 1.2rem; } }
  </style>
</head>
<body>
  <div class="container">
    <h2>Save Contact</h2>
    <div class="vcard-info">
      <div><label>Name:</label> <span id="name">` + name + `</span></div>
      <div><label>Organization:</label> <span id="org">` + organization + `</span></div>
      <div><label>Title:</label> <span id="title">` + title + `</span></div>
      <div><label>Phone:</label> <span id="phone">` + phone + `</span></div>
      <div><label>Email:</label> <span id="email">` + emailAddr + `</span></div>
      <div><label>Address:</label> <span id="address">` + address + `</span></div>
    </div>
    <a class="save-btn" id="saveBtn" href="data:text/vcard;charset=utf-8,` + url.QueryEscape(vcard) + `" download="contact.vcf">Save Contact</a>
    <div class="note">Tap "Save Contact" to download and import this contact into your phone.</div>
  </div>
</body>
</html>`
			return c.Status(fiber.StatusOK).Type("html").SendString(html)
		} else if eventName, ok := qr.Content["name"].(string); ok && qr.Content["date"] != nil && qr.Content["end_date"] != nil {
			// Event landing page (iCalendar)
			start, _ := qr.Content["date"].(string)
			end, _ := qr.Content["end_date"].(string)
			location, _ := qr.Content["location"].(string)
			description, _ := qr.Content["description"].(string)
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
			html := `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <title>Add Event</title>
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <style>
    :root { --primary1: #0c768a; --primary2: #0C8096; --primary3: #26666F; --bg1: #ffffff; --bg2: #fbfbfb; --bg3: #eef2f5; --text1: #424242; --text2: #000000; --border1: #d2d2d2; --border2: #d9d9d9; }
    body { background: var(--bg3); color: var(--text1); font-family: 'Segoe UI', Arial, sans-serif; margin: 0; padding: 0; min-height: 100vh; display: flex; align-items: center; justify-content: center; }
    .container { background: var(--bg1); border-radius: 16px; box-shadow: 0 4px 24px rgba(38, 102, 111, 0.08); padding: 2.5rem 1.5rem 2rem 1.5rem; max-width: 350px; width: 100%; border: 1px solid var(--border2); text-align: center; }
    h2 { color: var(--primary1); margin-bottom: 0.5rem; font-size: 1.6rem; font-weight: 700; }
    .event-info { background: var(--bg2); border: 1px solid var(--border1); border-radius: 10px; padding: 1rem; margin: 1.2rem 0 1.5rem 0; text-align: left; font-size: 1.05rem; word-break: break-all; }
    .event-info label { color: var(--primary3); font-weight: 600; margin-right: 0.5em; }
    .add-btn { display: block; width: 100%; background: linear-gradient(90deg, var(--primary1), var(--primary2)); color: var(--bg1); font-size: 1.15rem; font-weight: 600; border: none; border-radius: 8px; padding: 0.85rem 0; margin-bottom: 1.2rem; cursor: pointer; transition: background 0.2s; text-decoration: none; }
    .add-btn:hover, .add-btn:focus { background: var(--primary3); color: var(--bg1); }
    .note { font-size: 0.98rem; color: var(--text1); background: var(--bg3); border-radius: 6px; padding: 0.7em 1em; border: 1px solid var(--border1); margin-top: 0.5em; }
    @media (max-width: 480px) { .container { padding: 1.2rem 0.5rem 1.2rem 0.5rem; max-width: 98vw; } h2 { font-size: 1.2rem; } }
  </style>
</head>
<body>
  <div class="container">
    <h2>Add Event</h2>
    <div class="event-info">
      <div><label>Event:</label> <span id="event">` + eventName + `</span></div>
      <div><label>Start:</label> <span id="start">` + start + `</span></div>
      <div><label>End:</label> <span id="end">` + end + `</span></div>
      <div><label>Location:</label> <span id="location">` + location + `</span></div>
      <div><label>Description:</label> <span id="desc">` + description + `</span></div>
    </div>
    <a class="add-btn" id="addBtn" href="data:text/calendar;charset=utf-8,` + url.QueryEscape(ical) + `" download="event.ics">Add to Calendar</a>
    <div class="note">Tap "Add to Calendar" to download and import this event into your calendar app.</div>
  </div>
</body>
</html>`
			return c.Status(fiber.StatusOK).Type("html").SendString(html)
		}
		// fallback for other dynamic types
		if url, ok := qr.Content["url"].(string); ok {
			return c.Status(fiber.StatusOK).Type("text/plain").SendString(url)
		} else {
			return c.Status(fiber.StatusOK).JSON(qr.Content)
		}
	}

	// If the QR code is a WiFi QR code (dynamic or static), return the WiFi URI string
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
		return c.Status(fiber.StatusOK).Type("text/plain").SendString(wifiURI)
	}

	// Handle static QR codes by returning content directly
	if qr.Type == "static" {
		// Return the content as JSON
		return c.JSON(fiber.Map{"content": qr.Content})
	}

	// For dynamic QR codes, determine redirect URL
	var redirectURL string
	if qr.RedirectURL != "" {
		redirectURL = qr.RedirectURL
	} else if url, ok := qr.Content["url"].(string); ok {
		redirectURL = url
	} else {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "No redirect URL available"})
	}

	// Redirect to the target URL
	return c.Redirect(redirectURL, fiber.StatusFound)
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
