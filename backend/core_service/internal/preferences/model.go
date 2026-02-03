package preferences

import "time"

type Preference struct {
	ID             string
	UserID         string
	PreferenceType string
	Value          string
	CreatedAt      time.Time
}
