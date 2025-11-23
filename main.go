package main

import (
	"context"
	"fmt"
	"os"

	"github.com/YukiAminaka/cycle-route-backend/db"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/paulmach/orb"
)

func main() {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DATABASE_URL"))
    if err != nil {
            fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
            os.Exit(1)
    }
    defer conn.Close(context.Background())

    q := db.New(conn)

	user_location := orb.Point{139.737763, 35.6646848} 
	// Create a new user	
	user, err := q.CreateUser(context.Background(), db.CreateUserParams{
		Name:                   "JohnDoe",
		Email:                  pgtype.Text{String: "example@com", Valid: true},
		TotalTripDistance:      pgtype.Float8{Float64: 0, Valid: false},
		TotalTripDuration:      pgtype.Float8{Float64: 0, Valid: false},
		TotalTripElevationGain: pgtype.Float8{Float64: 0, Valid: false},
		StMakepoint:            -122.4194, // Longitude
		StMakepoint_2:          37.7749,   // Latitude
		FirstName:              pgtype.Text{String: "John", Valid: true},
		LastName:               pgtype.Text{String: "Doe", Valid: true},
		HasSetLocation:         pgtype.Bool{Bool: false, Valid: false},
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create user: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Created user: %+v\n", user)
}