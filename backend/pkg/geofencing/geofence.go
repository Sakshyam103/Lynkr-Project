package geofencing

import (
	"encoding/json"
	"errors"
	"math"
)

// GeofenceType represents the type of geofence
type GeofenceType string

const (
	// CircleGeofence represents a circular geofence
	CircleGeofence GeofenceType = "circle"
	// PolygonGeofence represents a polygon geofence
	PolygonGeofence GeofenceType = "polygon"
)

// Point represents a geographic point
type Point struct {
	Latitude  float64 `json:"lat"`
	Longitude float64 `json:"lng"`
}

// CircleGeofenceData represents a circular geofence
type CircleGeofenceData struct {
	Center Point   `json:"center"`
	Radius float64 `json:"radius"` // in meters
}

// PolygonGeofenceData represents a polygon geofence
type PolygonGeofenceData struct {
	Points []Point `json:"points"`
}

// GeofenceData represents the data for a geofence
type GeofenceData struct {
	Type    GeofenceType       `json:"type"`
	Circle  *CircleGeofenceData  `json:"circle,omitempty"`
	Polygon *PolygonGeofenceData `json:"polygon,omitempty"`
}

// Geofence represents a geofence for an event
type Geofence struct {
	EventID uint        `json:"event_id"`
	Name    string      `json:"name"`
	Data    GeofenceData `json:"data"`
}

// ParseGeofenceData parses geofence data from JSON
func ParseGeofenceData(data string) (*GeofenceData, error) {
	if data == "" {
		return nil, errors.New("empty geofence data")
	}

	var geofenceData GeofenceData
	err := json.Unmarshal([]byte(data), &geofenceData)
	if err != nil {
		return nil, err
	}

	// Validate geofence data
	switch geofenceData.Type {
	case CircleGeofence:
		if geofenceData.Circle == nil {
			return nil, errors.New("missing circle data for circle geofence")
		}
		if geofenceData.Circle.Radius <= 0 {
			return nil, errors.New("invalid radius for circle geofence")
		}
	case PolygonGeofence:
		if geofenceData.Polygon == nil {
			return nil, errors.New("missing polygon data for polygon geofence")
		}
		if len(geofenceData.Polygon.Points) < 3 {
			return nil, errors.New("polygon geofence must have at least 3 points")
		}
	default:
		return nil, errors.New("invalid geofence type")
	}

	return &geofenceData, nil
}

// IsPointInGeofence checks if a point is inside a geofence
func IsPointInGeofence(point Point, geofenceData *GeofenceData) bool {
	switch geofenceData.Type {
	case CircleGeofence:
		return isPointInCircle(point, geofenceData.Circle)
	case PolygonGeofence:
		return isPointInPolygon(point, geofenceData.Polygon)
	default:
		return false
	}
}

// isPointInCircle checks if a point is inside a circular geofence
func isPointInCircle(point Point, circle *CircleGeofenceData) bool {
	distance := calculateDistance(point, circle.Center)
	return distance <= circle.Radius
}

// isPointInPolygon checks if a point is inside a polygon geofence using the ray casting algorithm
func isPointInPolygon(point Point, polygon *PolygonGeofenceData) bool {
	inside := false
	j := len(polygon.Points) - 1

	for i := 0; i < len(polygon.Points); i++ {
		if ((polygon.Points[i].Latitude > point.Latitude) != (polygon.Points[j].Latitude > point.Latitude)) &&
			(point.Longitude < (polygon.Points[j].Longitude-polygon.Points[i].Longitude)*(point.Latitude-polygon.Points[i].Latitude)/(polygon.Points[j].Latitude-polygon.Points[i].Latitude)+polygon.Points[i].Longitude) {
			inside = !inside
		}
		j = i
	}

	return inside
}

// calculateDistance calculates the distance between two points using the Haversine formula
func calculateDistance(p1, p2 Point) float64 {
	const earthRadius = 6371000 // Earth radius in meters

	lat1 := p1.Latitude * math.Pi / 180
	lat2 := p2.Latitude * math.Pi / 180
	deltaLat := (p2.Latitude - p1.Latitude) * math.Pi / 180
	deltaLng := (p2.Longitude - p1.Longitude) * math.Pi / 180

	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1)*math.Cos(lat2)*
			math.Sin(deltaLng/2)*math.Sin(deltaLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadius * c
}

// CreateCircleGeofence creates a circular geofence
func CreateCircleGeofence(eventID uint, name string, center Point, radius float64) (*Geofence, error) {
	if radius <= 0 {
		return nil, errors.New("radius must be positive")
	}

	return &Geofence{
		EventID: eventID,
		Name:    name,
		Data: GeofenceData{
			Type: CircleGeofence,
			Circle: &CircleGeofenceData{
				Center: center,
				Radius: radius,
			},
		},
	}, nil
}

// CreatePolygonGeofence creates a polygon geofence
func CreatePolygonGeofence(eventID uint, name string, points []Point) (*Geofence, error) {
	if len(points) < 3 {
		return nil, errors.New("polygon must have at least 3 points")
	}

	return &Geofence{
		EventID: eventID,
		Name:    name,
		Data: GeofenceData{
			Type: PolygonGeofence,
			Polygon: &PolygonGeofenceData{
				Points: points,
			},
		},
	}, nil
}

// ToJSON converts a geofence to JSON
func (g *Geofence) ToJSON() (string, error) {
	data, err := json.Marshal(g.Data)
	if err != nil {
		return "", err
	}
	return string(data), nil
}