package model

import (
	"time"
)

// QRCode represents the main QR code entity
type QRCode struct {
	ID          string                 `json:"id" db:"id"`
	Type        QRCodeType             `json:"type" db:"type"`
	Title       string                 `json:"title" db:"title"`
	Description string                 `json:"description,omitempty" db:"description"`
	RedirectURL string                 `json:"redirect_url,omitempty" db:"redirect_url"`
	ShortURL    string                 `json:"short_url,omitempty" db:"short_url"`
	Content     map[string]interface{} `json:"content" db:"content"` // Stores type-specific content
	CreatedAt   time.Time              `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at" db:"updated_at"`
	ExpiresAt   *time.Time             `json:"expires_at,omitempty" db:"expires_at"`
	Analytics   bool                   `json:"analytics" db:"analytics"`
	Active      bool                   `json:"active" db:"active"`
	Tags        []string               `json:"tags,omitempty" db:"tags"`
	GroupID     *string                `json:"group_id,omitempty" db:"group_id"`
	Design      QRCodeDesign           `json:"design,omitempty" db:"design"`
	FileRefs    []FileReference        `json:"file_refs,omitempty" db:"file_refs"`
}

// QRCodeType represents the different types of QR codes
type QRCodeType string

const (
	QRTypeWebsite     QRCodeType = "website"
	QRTypeSearch      QRCodeType = "search"
	QRTypeDynamic     QRCodeType = "dynamic"
	QRTypeStatic      QRCodeType = "static"
	QRTypeVirtualCard QRCodeType = "virtual_card"
	QRTypePDF         QRCodeType = "pdf"
	QRTypeSocialMedia QRCodeType = "social_media"
	QRTypeInstagram   QRCodeType = "instagram"
	QRTypeImages      QRCodeType = "images"
	QRTypeApp         QRCodeType = "app"
	QRTypeBusiness    QRCodeType = "business"
	QRTypeEvent       QRCodeType = "event"
	QRTypeBarcode2D   QRCodeType = "barcode_2d"
	QRTypeFeedback    QRCodeType = "feedback"
	QRTypeRating      QRCodeType = "rating"
	QRTypeEmail       QRCodeType = "email"
	QRTypeText        QRCodeType = "text"
	QRTypeWiFi        QRCodeType = "wifi"
	QRTypeSMS         QRCodeType = "sms"
)

// QRCodeDesign represents QR code customization options
type QRCodeDesign struct {
	ForegroundColor string `json:"foreground_color,omitempty"`
	BackgroundColor string `json:"background_color,omitempty"`
	LogoURL         string `json:"logo_url,omitempty"`
	LogoSize        int    `json:"logo_size,omitempty"`
	Shape           string `json:"shape,omitempty"` // square, rounded, circular
}

// FileReference represents uploaded files
type FileReference struct {
	ID       string `json:"id"`
	Filename string `json:"filename"`
	URL      string `json:"url"`
	Size     int64  `json:"size"`
	Type     string `json:"type"`
}

// Type-specific content structures
type WebsiteContent struct {
	URL string `json:"url"`
}

type SearchContent struct {
	Engine string `json:"engine"` // google, bing, youtube
	Query  string `json:"query"`
}

type DynamicContent struct {
	Title       string `json:"title"`
	RedirectURL string `json:"redirect_url"`
}

type VirtualCardContent struct {
	FullName  string `json:"full_name"`
	Phone     string `json:"phone"`
	Email     string `json:"email"`
	Website   string `json:"website"`
	Company   string `json:"company"`
	JobTitle  string `json:"job_title"`
	Address   string `json:"address"`
	PhotoURL  string `json:"photo_url"`
}

type PDFContent struct {
	FileURL string `json:"file_url"`
}

type SocialMediaContent struct {
	Platform string `json:"platform"` // facebook, tiktok, etc.
	URL      string `json:"url"`
	Handle   string `json:"handle"`
}

type InstagramContent struct {
	Handle string `json:"handle"`
	URL    string `json:"url"`
}

type ImagesContent struct {
	GalleryURL string   `json:"gallery_url"`
	ImageURLs  []string `json:"image_urls"`
}

type AppContent struct {
	AppStoreURL string `json:"app_store_url"`
	DeepLink    string `json:"deep_link,omitempty"`
}

type BusinessContent struct {
	Name        string            `json:"name"`
	Tagline     string            `json:"tagline"`
	ContactInfo map[string]string `json:"contact_info"`
	LogoURL     string            `json:"logo_url"`
	Description string            `json:"description"`
	Website     string            `json:"website"`
	SocialLinks map[string]string `json:"social_links"`
}

type EventContent struct {
	Name        string    `json:"name"`
	DateTime    time.Time `json:"date_time"`
	Location    string    `json:"location"`
	Description string    `json:"description"`
	RSVPLink    string    `json:"rsvp_link,omitempty"`
}

type Barcode2DContent struct {
	Data string `json:"data"`
}

type FeedbackContent struct {
	FormURL        string `json:"form_url"`
	ThankYouMsg    string `json:"thank_you_msg,omitempty"`
}

type RatingContent struct {
	FormURL string `json:"form_url"`
	Scale   string `json:"scale"` // stars, emojis, nps
}

type EmailContent struct {
	Recipient string `json:"recipient"`
	Subject   string `json:"subject"`
	Body      string `json:"body"`
}

type TextContent struct {
	Text string `json:"text"`
}

type WiFiContent struct {
	SSID         string `json:"ssid"`
	Password     string `json:"password"`
	Encryption   string `json:"encryption"` // WPA, WEP, None
	Hidden       bool   `json:"hidden"`
}

type SMSContent struct {
	PhoneNumber string `json:"phone_number"`
	Message     string `json:"message"`
}

// QRCodeGroup represents a group of QR codes
type QRCodeGroup struct {
	ID          string    `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	Description string    `json:"description,omitempty" db:"description"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// QRCodeAnalytics represents scan analytics
type QRCodeAnalytics struct {
	ID        string    `json:"id" db:"id"`
	QRCodeID  string    `json:"qr_code_id" db:"qr_code_id"`
	IPAddress string    `json:"ip_address" db:"ip_address"`
	UserAgent string    `json:"user_agent" db:"user_agent"`
	Location  string    `json:"location,omitempty" db:"location"`
	Device    string    `json:"device,omitempty" db:"device"`
	ScannedAt time.Time `json:"scanned_at" db:"scanned_at"`
}
