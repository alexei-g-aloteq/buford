package payload

import (
	"encoding/json"

	"github.com/alexei-g-aloteq/buford/payload/badge"
)

// https://developer.apple.com/documentation/usernotifications/setting_up_a_remote_notification_server/generating_a_remote_notification?language=objc

// APS is Apple's reserved namespace.
// Use it for payloads destined to mobile devices (iOS).
type APS struct {
	// Alert dictionary.
	Alert Alert

	// Badge to display on the app icon.
	// Set to badge.Preserve (default), badge.Clear
	// or a specific value with badge.New(n).
	Badge badge.Badge

	// Details how to play alert.
	Sound Sound

	// Thread identifier to create notification groups in iOS 12 or newer.
	ThreadID string

	// Category identifier for custom actions in iOS 8 or newer.
	Category string

	// Content available is for silent notifications
	// with no alert, sound, or badge.
	ContentAvailable bool

	// Mutable is used for Service Extensions introduced in iOS 10.
	MutableContent bool

	// Content identifier.
	TargetContentID string

	// URL arguments for Safari pushes: https://developer.apple.com/library/archive/documentation/NetworkingInternet/Conceptual/NotificationProgrammingGuideForWebsites/PushNotifications/PushNotifications.html#//apple_ref/doc/uid/TP40013225-CH3-SW17
	SafariURLArgs []string

	InterruptionLevel InterruptionLevel

	RelevanceScore float32
}

type InterruptionLevel string

const (
	InterruptionLevelPassive 		InterruptionLevel = "passive"
	InterruptionLevelActive 		InterruptionLevel = "active"
	InterruptionLevelTimeSensitive 	InterruptionLevel = "time-sensitive"
	InterruptionLevelCritical 		InterruptionLevel = "critical"
)

// Alert dictionary.
type Alert struct {
	// Title is a short string shown briefly on Apple Watch in iOS 8.2 or newer.
	Title        string   `json:"title,omitempty"`
	TitleLocKey  string   `json:"title-loc-key,omitempty"`
	TitleLocArgs []string `json:"title-loc-args,omitempty"`

	// Subtitle added in iOS 10
	Subtitle        string   `json:"subtitle,omitempty"`
	SubtitleLocKey  string   `json:"subtitle-loc-key,omitempty"`
	SubtitleLocArgs []string `json:"subtitle-loc-args,omitempty"`

	// Body text of the alert message.
	Body    string   `json:"body,omitempty"`
	LocKey  string   `json:"loc-key,omitempty"`
	LocArgs []string `json:"loc-args,omitempty"`

	// Key for localized string for "View" button.
	ActionLocKey string `json:"action-loc-key,omitempty"`

	// Image file to be used when user taps or slides the action button.
	LaunchImage string `json:"launch-image,omitempty"`

	// String for "View" button on Safari.
	SafariAction string `json:"action,omitempty"`
}

// Sound dictionary.
type Sound struct {
	SoundName      string  `json:"name,omitempty"`
	IsCritical     int     `json:"critical,omitempty"`
	CriticalVolume float32 `json:"volume,omitempty"`
}

// isSimple alert with only Body set.
func (a *Alert) isSimple() bool {
	return len(a.Title) == 0 && len(a.Subtitle) == 0 &&
		len(a.LaunchImage) == 0 &&
		len(a.TitleLocKey) == 0 && len(a.TitleLocArgs) == 0 &&
		len(a.LocKey) == 0 && len(a.LocArgs) == 0 && len(a.ActionLocKey) == 0
}

// isZero if no Alert fields are set.
func (a *Alert) isZero() bool {
	return len(a.Body) == 0 && a.isSimple()
}

// Map returns the payload as a map that you can customize
// before serializing it to JSON.
func (a *APS) Map() map[string]interface{} {
	aps := make(map[string]interface{}, 5)

	if !a.Alert.isZero() {
		if a.Alert.isSimple() {
			aps["alert"] = a.Alert.Body
		} else {
			aps["alert"] = a.Alert
		}
	}
	if n, ok := a.Badge.Number(); ok {
		aps["badge"] = n
	}
	if a.Sound.SoundName != "" {
		aps["sound"] = a.Sound
	}
	if a.ContentAvailable {
		aps["content-available"] = 1
	}
	if a.Category != "" {
		aps["category"] = a.Category
	}
	if a.MutableContent {
		aps["mutable-content"] = 1
	}
	if a.ThreadID != "" {
		aps["thread-id"] = a.ThreadID
	}
	if a.TargetContentID != "" {
		aps["target-content-id"] = a.TargetContentID
	}
	if len(a.SafariURLArgs) > 0 {
		aps["url-args"] = a.SafariURLArgs
	}
	if a.InterruptionLevel != "" {
		aps["interruption-level"] = a.InterruptionLevel
	}
	if a.RelevanceScore > 0 {
		aps["relevance-score"] = a.RelevanceScore
	}

	// wrap in "aps" to form the final payload
	return map[string]interface{}{"aps": aps}
}

// MarshalJSON allows you to json.Marshal(aps) directly.
func (a APS) MarshalJSON() ([]byte, error) {
	return json.Marshal(a.Map())
}

// Validate that a payload has the correct fields.
func (a *APS) Validate() error {
	if a == nil {
		return ErrIncomplete
	}

	// must have a body or a badge (or custom data)
	if len(a.Alert.Body) == 0 && a.Badge == badge.Preserve {
		return ErrIncomplete
	}
	return nil
}
