package constants

type RelationType uint8 // uint8 is more than enough (255 types)

// Fast validation map (O(1))
var validRelationTypes = map[RelationType]struct{}{
	RelationRelated:   {},
	RelationUpsell:    {},
	RelationCrossSell: {},
}

const (
	RelationRelated RelationType = iota + 1
	RelationUpsell
	RelationCrossSell
	// Add new ones here → no DB change ever
)

// String returns the DB-friendly value (you can change the strings anytime)
func (r RelationType) String() string {
	switch r {
	case RelationRelated:
		return "related"
	case RelationUpsell:
		return "upsell"
	case RelationCrossSell:
		return "cross_sell"
	default:
		return "unknown"
	}
}

// For API responses / frontend you can have a display name
func (r RelationType) DisplayName() string {
	switch r {
	case RelationRelated:
		return "Related Products"
	case RelationUpsell:
		return "Upsell"
	case RelationCrossSell:
		return "Cross-sell"
	default:
		return "Unknown"
	}
}

// IsValid is what you call in your handlers / service layer
func (r RelationType) IsValid() bool {
	_, ok := validRelationTypes[r]
	return ok
}

// ParseFromString is useful when you receive JSON/API input as string
func ParseRelationType(s string) (RelationType, bool) {
	switch s {
	case "related":
		return RelationRelated, true
	case "upsell":
		return RelationUpsell, true
	case "cross_sell":
		return RelationCrossSell, true
	default:
		return 0, false
	}
}
