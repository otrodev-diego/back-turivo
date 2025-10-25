package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"

	"turivo-backend/internal/infrastructure/config"
)

func main() {
	demo := flag.Bool("demo", false, "Generate demo data aligned with frontend")
	flag.Parse()

	// Load configuration
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// Connect to database
	conn, err := pgx.Connect(context.Background(), cfg.DB.DSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer conn.Close(context.Background())

	if *demo {
		log.Println("Generating demo data...")
		if err := generateDemoData(conn); err != nil {
			log.Fatalf("Failed to generate demo data: %v", err)
		}
		log.Println("Demo data generated successfully!")
	} else {
		log.Println("Use -demo flag to generate demo data")
	}
}

func generateDemoData(conn *pgx.Conn) error {
	ctx := context.Background()

	// Start transaction
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	// Create companies
	companies := []struct {
		id           uuid.UUID
		name         string
		rut          string
		contactEmail string
		sector       string
	}{
		{uuid.New(), "Turivo", "76.123.456-7", "contacto@turivo.com", "TURISMO"},
		{uuid.New(), "Turismo Andes", "76.234.567-8", "info@turismoandes.cl", "TURISMO"},
		{uuid.New(), "Minera del Norte", "96.345.678-9", "contacto@mineranorte.cl", "MINERIA"},
	}

	for _, c := range companies {
		_, err := tx.Exec(ctx, `
			INSERT INTO companies (id, name, rut, contact_email, status, sector, created_at, updated_at)
			VALUES ($1, $2, $3, $4, 'ACTIVE', $5, NOW(), NOW())
		`, c.id, c.name, c.rut, c.contactEmail, c.sector)
		if err != nil {
			return fmt.Errorf("failed to create company %s: %w", c.name, err)
		}
	}

	// Create hotels
	hotels := []struct {
		id           uuid.UUID
		name         string
		city         string
		contactEmail string
	}{
		{uuid.New(), "Hotel Miramar", "Valparaíso", "reservas@hotelmiramar.cl"},
		{uuid.New(), "Hotel Andes", "Santiago", "contacto@hotelandes.cl"},
		{uuid.New(), "Hotel Patagonia", "Puerto Montt", "info@hotelpatagonia.cl"},
	}

	for _, h := range hotels {
		_, err := tx.Exec(ctx, `
			INSERT INTO hotels (id, name, city, contact_email, created_at, updated_at)
			VALUES ($1, $2, $3, $4, NOW(), NOW())
		`, h.id, h.name, h.city, h.contactEmail)
		if err != nil {
			return fmt.Errorf("failed to create hotel %s: %w", h.name, err)
		}
	}

	// Create users
	// Usuarios por defecto comentados por seguridad
	// Descomenta y modifica las credenciales según sea necesario
	/*
		users := []struct {
			id           uuid.UUID
			name         string
			email        string
			passwordHash string
			role         string
			orgID        *uuid.UUID
		}{
			{uuid.New(), "Admin Sistema", "admin@turivo.com", "$2a$10$9pjpOYuBh0O/loXHGwBOg.l6hrmZHtsae/2UqhD23ff4O5nTwgymC", "ADMIN", nil},                      // password
			{uuid.New(), "Juan Pérez", "juan@turivo.com", "$2a$10$9pjpOYuBh0O/loXHGwBOg.l6hrmZHtsae/2UqhD23ff4O5nTwgymC", "COMPANY", &companies[0].id},           // password
			{uuid.New(), "María González", "maria@turismoandes.cl", "$2a$10$9pjpOYuBh0O/loXHGwBOg.l6hrmZHtsae/2UqhD23ff4O5nTwgymC", "COMPANY", &companies[1].id}, // password
			{uuid.New(), "Cliente Demo", "cliente@demo.com", "$2a$10$9pjpOYuBh0O/loXHGwBOg.l6hrmZHtsae/2UqhD23ff4O5nTwgymC", "USER", nil},                        // password
		}
	*/

	// Array vacío para no crear usuarios por defecto
	users := []struct {
		id           uuid.UUID
		name         string
		email        string
		passwordHash string
		role         string
		orgID        *uuid.UUID
	}{}

	for _, u := range users {
		_, err := tx.Exec(ctx, `
			INSERT INTO users (id, name, email, password_hash, role, status, org_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, 'ACTIVE', $6, NOW(), NOW())
		`, u.id, u.name, u.email, u.passwordHash, u.role, u.orgID)
		if err != nil {
			return fmt.Errorf("failed to create user %s: %w", u.name, err)
		}
	}

	// Create drivers with aligned IDs
	drivers := []struct {
		id        string
		firstName string
		lastName  string
		rutOrDNI  string
		phone     string
		email     string
		status    string
	}{
		{"CON-001", "Carlos", "Mendoza", "12.345.678-9", "+56912345678", "carlos@turivo.com", "ACTIVE"},
		{"CON-002", "Ana", "Silva", "98.765.432-1", "+56987654321", "ana@turivo.com", "ACTIVE"},
		{"CON-003", "Pedro", "Ramírez", "11.222.333-4", "+56911222333", "pedro@turivo.com", "ACTIVE"},
	}

	for _, d := range drivers {
		// Insert driver
		_, err := tx.Exec(ctx, `
			INSERT INTO drivers (id, first_name, last_name, rut_or_dni, phone, email, status, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, NOW(), NOW())
		`, d.id, d.firstName, d.lastName, d.rutOrDNI, d.phone, d.email, d.status)
		if err != nil {
			return fmt.Errorf("failed to create driver %s: %w", d.id, err)
		}

		// Insert driver license
		_, err = tx.Exec(ctx, `
			INSERT INTO driver_licenses (driver_id, number, class, issued_at, expires_at)
			VALUES ($1, $2, 'A3', $3, $4)
		`, d.id, fmt.Sprintf("LIC%s", d.id[4:]), time.Now().AddDate(-2, 0, 0), time.Now().AddDate(3, 0, 0))
		if err != nil {
			return fmt.Errorf("failed to create license for driver %s: %w", d.id, err)
		}

		// Insert background check
		_, err = tx.Exec(ctx, `
			INSERT INTO driver_background_checks (driver_id, status, checked_at)
			VALUES ($1, 'APPROVED', NOW())
		`, d.id)
		if err != nil {
			return fmt.Errorf("failed to create background check for driver %s: %w", d.id, err)
		}

		// Insert availability
		_, err = tx.Exec(ctx, `
			INSERT INTO driver_availability (driver_id, regions, days, time_ranges, updated_at)
			VALUES ($1, '["RM", "V"]', '["monday", "tuesday", "wednesday", "thursday", "friday"]', '[{"from": "08:00", "to": "18:00"}]', NOW())
		`, d.id)
		if err != nil {
			return fmt.Errorf("failed to create availability for driver %s: %w", d.id, err)
		}
	}

	// Create sample reservations
	reservations := []struct {
		id          string
		userID      *uuid.UUID
		pickup      string
		destination string
		datetime    time.Time
		passengers  int
		status      string
		amount      float64
	}{
		{"RSV-1001", &users[3].id, "Hotel Miramar, Valparaíso", "Aeropuerto SCL, Santiago", time.Now().AddDate(0, 0, 7), 2, "ACTIVA", 120000},
		{"RSV-1002", &users[3].id, "Hotel Andes, Santiago", "Mall Costanera Center", time.Now().AddDate(0, 0, 1), 4, "PROGRAMADA", 80000},
		{"RSV-1003", nil, "Oficina Central", "Mina El Teniente", time.Now().AddDate(0, 0, -1), 8, "COMPLETADA", 250000},
	}

	for _, r := range reservations {
		_, err := tx.Exec(ctx, `
			INSERT INTO reservations (id, user_id, pickup, destination, datetime, passengers, status, amount, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, NOW(), NOW())
		`, r.id, r.userID, r.pickup, r.destination, r.datetime, r.passengers, r.status, r.amount)
		if err != nil {
			return fmt.Errorf("failed to create reservation %s: %w", r.id, err)
		}

		// Add timeline events
		events := []struct {
			title       string
			description string
			variant     string
			at          time.Time
		}{
			{"Reserva creada", "La reserva ha sido creada exitosamente", "success", time.Now().Add(-24 * time.Hour)},
		}

		if r.status == "PROGRAMADA" {
			events = append(events, struct {
				title       string
				description string
				variant     string
				at          time.Time
			}{"Reserva programada", "La reserva ha sido programada", "info", time.Now().Add(-12 * time.Hour)})
		}

		if r.status == "COMPLETADA" {
			events = append(events,
				struct {
					title       string
					description string
					variant     string
					at          time.Time
				}{"Reserva programada", "La reserva ha sido programada", "info", time.Now().Add(-12 * time.Hour)},
				struct {
					title       string
					description string
					variant     string
					at          time.Time
				}{"Servicio completado", "El servicio ha sido completado exitosamente", "success", time.Now().Add(-1 * time.Hour)})
		}

		for _, event := range events {
			_, err := tx.Exec(ctx, `
				INSERT INTO reservation_timeline (id, reservation_id, title, description, at, variant, created_at)
				VALUES ($1, $2, $3, $4, $5, $6, NOW())
			`, uuid.New(), r.id, event.title, event.description, event.at, event.variant)
			if err != nil {
				return fmt.Errorf("failed to create timeline event for reservation %s: %w", r.id, err)
			}
		}
	}

	// Create sample payments
	for _, r := range reservations {
		if r.status == "COMPLETADA" {
			paymentID := uuid.New()
			_, err := tx.Exec(ctx, `
				INSERT INTO payments (id, reservation_id, gateway, amount, currency, status, transaction_ref, payload, created_at)
				VALUES ($1, $2, 'WEBPAY_PLUS', $3, 'CLP', 'APPROVED', $4, '{"vci": "TSY", "status": "AUTHORIZED"}', NOW())
			`, paymentID, r.id, r.amount, fmt.Sprintf("WP_%d", time.Now().Unix()))
			if err != nil {
				return fmt.Errorf("failed to create payment for reservation %s: %w", r.id, err)
			}
		}
	}

	// Create sample feedback
	feedbackData := []struct {
		tripID        string
		passengerName string
		rating        int
		comment       string
	}{
		{"RSV-1003", "Roberto Martínez", 5, "Excelente servicio, muy puntual y profesional"},
		{"RSV-1003", "Carmen López", 4, "Buen conductor, viaje cómodo"},
	}

	for _, f := range feedbackData {
		_, err := tx.Exec(ctx, `
			INSERT INTO driver_feedback (id, trip_id, passenger_name, rating, comment, created_at)
			VALUES ($1, $2, $3, $4, $5, NOW())
		`, uuid.New(), f.tripID, f.passengerName, f.rating, f.comment)
		if err != nil {
			return fmt.Errorf("failed to create feedback: %w", err)
		}
	}

	// Commit transaction
	return tx.Commit(ctx)
}
