package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/chibx/vuecom/backend/services/gateway/internal/utils"
	"go.uber.org/zap"
)

type continent struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

type country struct {
	ID          int    `json:"id"`
	ContinentID int    `json:"continentId"`
	Name        string `json:"name"`
	Code        string `json:"code"`
	Phone       string `json:"phone"`
	Currency    string `json:"currency"`
}

type state struct {
	ID        int    `json:"id"`
	CountryID int    `json:"countryId"`
	Name      string `json:"name"`
}

type city struct {
	ID      int    `json:"id"`
	StateID int    `json:"stateId"`
	Name    string `json:"name"`
}

type transformedData struct {
	Continents []continent `json:"continents"`
	Countries  []country   `json:"countries"`
	States     []state     `json:"states"`
	Cities     []city      `json:"cities"`
}

func escapeSQL(s string) string {
	return strings.ReplaceAll(s, "'", "''")
}

func generateContinents(continents []continent) string {
	if len(continents) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("INSERT INTO continents (id, name) VALUES")
	for i, item := range continents {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("(%d,'%s')", item.ID, escapeSQL(item.Name)))
	}
	sb.WriteString(";\n")
	return sb.String()
}

func generateCountries(countries []country) string {
	if len(countries) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("INSERT INTO countries (id, code, continent_id, phone, currency, name) VALUES")
	for i, item := range countries {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("(%d,'%s',%d,'%s','%s','%s')",
			item.ID,
			escapeSQL(item.Code),
			item.ContinentID,
			escapeSQL(item.Phone),
			escapeSQL(item.Currency),
			escapeSQL(item.Name),
		))
	}
	sb.WriteString(";\n")
	return sb.String()
}

func generateStates(states []state) string {
	if len(states) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("INSERT INTO states (id, country_id, name) VALUES")
	for i, item := range states {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("(%d,%d,'%s')", item.ID, item.CountryID, escapeSQL(item.Name)))
	}
	sb.WriteString(";\n")
	return sb.String()
}

func generateCities(cities []city) string {
	if len(cities) == 0 {
		return ""
	}
	var sb strings.Builder
	sb.WriteString("INSERT INTO cities (id, state_id, name) VALUES")
	for i, item := range cities {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("(%d,%d,'%s')", item.ID, item.StateID, escapeSQL(item.Name)))
	}
	sb.WriteString(";\n")
	return sb.String()
}

func SeedWorldData(ctx context.Context, db *sql.DB) error {
	const sql = `SELECT
    (EXISTS (SELECT 1 FROM continents LIMIT 1))
    AND (EXISTS (SELECT 1 FROM countries LIMIT 1))
    AND (EXISTS (SELECT 1 FROM states LIMIT 1))
    AND (EXISTS (SELECT 1 FROM cities LIMIT 1))
    AS all_tables_populated;
    `
	var allExists bool
	logger := utils.Logger()
	row := db.QueryRowContext(ctx, sql)
	var err = row.Err()
	if err != nil {
		logger.Error("Error while checking to see if rows exists in the specified tables!!!", zap.Error(err))
		return err
	}

	row.Scan(&allExists)
	if allExists {
		return nil
	}

	data, err := os.ReadFile("./db-ready-world.json")
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		return err
	}

	var td transformedData
	if err := json.Unmarshal(data, &td); err != nil {
		fmt.Printf("Error parsing JSON: %v\n", err)
		return err
	}

	var out strings.Builder
	out.WriteString(generateContinents(td.Continents))
	out.WriteString(generateCountries(td.Countries))
	out.WriteString(generateStates(td.States))
	out.WriteString(generateCities(td.Cities))

	_, err = db.ExecContext(ctx, out.String())

	if err != nil {
		logger.Error("[Seeding]: Error inserting fields into columns", zap.Error(err))
		return err
	}

	return nil
}
