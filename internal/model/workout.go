package model

import "time"

type Workout struct {
	ID                  int       `json:"id"`
	ProfileID           int       `json:"profile_id"`
	Name                string    `json:"name"`
	Description         string    `json:"description"`
	MentalEnergyLevel   int       `json:"mental_energy_level"`
	PhysicalEnergyLevel int       `json:"physical_energy_level"`
	StartDate           time.Time `json:"start_date"`
	EndDate             time.Time `json:"end_date"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}
