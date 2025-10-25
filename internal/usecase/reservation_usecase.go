package usecase

import (
	"fmt"
	"strconv"
	"time"

	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type ReservationUseCase struct {
	reservationRepo domain.ReservationRepository
	driverRepo      domain.DriverRepository
	userRepo        domain.UserRepository
	emailService    domain.EmailService
	logger          *zap.Logger
}

func NewReservationUseCase(
	reservationRepo domain.ReservationRepository,
	driverRepo domain.DriverRepository,
	userRepo domain.UserRepository,
	emailService domain.EmailService,
	logger *zap.Logger,
) *ReservationUseCase {
	return &ReservationUseCase{
		reservationRepo: reservationRepo,
		driverRepo:      driverRepo,
		userRepo:        userRepo,
		emailService:    emailService,
		logger:          logger,
	}
}

func (uc *ReservationUseCase) CreateReservation(req domain.CreateReservationRequest, vehicleType domain.VehicleType, hasSpecialLanguage bool, stops int) (*domain.Reservation, error) {
	uc.logger.Info("Creating reservation")

	// Validate datetime is not in the past
	if req.DateTime.Before(time.Now()) {
		return nil, domain.ErrReservationPastDate
	}

	// Generate reservation ID
	reservationID := uc.reservationRepo.GenerateID()

	// Create reservation
	reservation := &domain.Reservation{
		ID:          reservationID,
		UserID:      req.UserID,
		OrgID:       req.OrgID,
		Pickup:      req.Pickup,
		Destination: req.Destination,
		DateTime:    req.DateTime,
		Passengers:  req.Passengers,
		Status:      domain.ReservationStatusActiva,
		Notes:       req.Notes,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	// Calculate price
	price := reservation.CalculatePrice(vehicleType, hasSpecialLanguage, stops)
	reservation.Amount = &price

	// Calculate distance
	distance := reservation.CalculateDistance()
	reservation.DistanceKM = &distance

	if err := uc.reservationRepo.Create(reservation); err != nil {
		uc.logger.Error("Failed to create reservation", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Add initial timeline event
	timelineEvent := domain.TimelineEvent{
		ReservationID: reservationID,
		Title:         "Reserva creada",
		Description:   "La reserva ha sido creada exitosamente",
		At:            time.Now(),
		Variant:       "success",
		CreatedAt:     time.Now(),
	}

	if err := uc.reservationRepo.AddTimelineEvent(reservationID, timelineEvent); err != nil {
		uc.logger.Warn("Failed to add initial timeline event", zap.Error(err))
		// Don't fail the reservation creation for this
	}

	// Get user information for emails
	user, err := uc.userRepo.GetByID(*req.UserID)
	if err != nil {
		uc.logger.Warn("Failed to get user for email notifications", zap.Error(err))
	} else {
		// Send confirmation email to user
		if err := uc.emailService.SendReservationCreated(user.Email, reservation, user); err != nil {
			uc.logger.Warn("Failed to send reservation confirmation email", zap.Error(err))
			// Don't fail the reservation creation for email errors
		}

		// Send notification email to operations
		operationsEmail := "djaramontenegro@gmail.com"
		if err := uc.emailService.SendReservationNotification(operationsEmail, reservation, user); err != nil {
			uc.logger.Warn("Failed to send reservation notification email", zap.Error(err))
			// Don't fail the reservation creation for email errors
		}
	}

	uc.logger.Info("Reservation created successfully", zap.String("reservation_id", reservation.ID))
	return reservation, nil
}

func (uc *ReservationUseCase) GetReservationByID(id string) (*domain.Reservation, error) {
	reservation, err := uc.reservationRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrReservationNotFound
		}
		uc.logger.Error("Failed to get reservation by ID", zap.Error(err), zap.String("reservation_id", id))
		return nil, domain.ErrInternalError
	}

	return reservation, nil
}

func (uc *ReservationUseCase) ListReservations(req domain.ListReservationsRequest) ([]*domain.Reservation, int, error) {
	// Set default pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}
	if req.Sort == "" {
		req.Sort = "datetime"
	}

	reservations, total, err := uc.reservationRepo.List(req)
	if err != nil {
		uc.logger.Error("Failed to list reservations", zap.Error(err))
		return nil, 0, domain.ErrInternalError
	}

	return reservations, total, nil
}

func (uc *ReservationUseCase) UpdateReservation(id string, req domain.UpdateReservationRequest) (*domain.Reservation, error) {
	uc.logger.Info("Updating reservation", zap.String("reservation_id", id))

	// Check if reservation exists
	existingReservation, err := uc.reservationRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrReservationNotFound
		}
		uc.logger.Error("Failed to get reservation for update", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Validate datetime is not in the past if being updated
	if req.DateTime != nil && req.DateTime.Before(time.Now()) {
		return nil, domain.ErrReservationPastDate
	}

	// Check if reservation can be modified
	if existingReservation.Status == domain.ReservationStatusCompletada || existingReservation.Status == domain.ReservationStatusCancelada {
		return nil, domain.ErrInvalidInput
	}

	reservation, err := uc.reservationRepo.Update(id, req)
	if err != nil {
		uc.logger.Error("Failed to update reservation", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Add timeline event for update
	timelineEvent := domain.TimelineEvent{
		ReservationID: id,
		Title:         "Reserva actualizada",
		Description:   "Los detalles de la reserva han sido actualizados",
		At:            time.Now(),
		Variant:       "info",
		CreatedAt:     time.Now(),
	}

	if err := uc.reservationRepo.AddTimelineEvent(id, timelineEvent); err != nil {
		uc.logger.Warn("Failed to add timeline event for update", zap.Error(err))
	}

	uc.logger.Info("Reservation updated successfully", zap.String("reservation_id", reservation.ID))
	return reservation, nil
}

func (uc *ReservationUseCase) ChangeReservationStatus(id string, req domain.ChangeReservationStatusRequest) (*domain.Reservation, error) {
	uc.logger.Info("Changing reservation status", zap.String("reservation_id", id), zap.String("new_status", string(req.NewStatus)))

	// Check if reservation exists
	existingReservation, err := uc.reservationRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrReservationNotFound
		}
		uc.logger.Error("Failed to get reservation for status change", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Validate status transition
	if !existingReservation.CanTransitionTo(req.NewStatus) {
		uc.logger.Warn("Invalid status transition",
			zap.String("current_status", string(existingReservation.Status)),
			zap.String("new_status", string(req.NewStatus)))
		return nil, domain.ErrInvalidStatusTransition
	}

	// Change status
	if err := uc.reservationRepo.ChangeStatus(id, req.NewStatus); err != nil {
		uc.logger.Error("Failed to change reservation status", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	// Add timeline event
	var title, description, variant string
	switch req.NewStatus {
	case domain.ReservationStatusProgramada:
		title = "Reserva programada"
		description = "La reserva ha sido programada"
		variant = "info"
	case domain.ReservationStatusCompletada:
		title = "Reserva completada"
		description = "El servicio ha sido completado exitosamente"
		variant = "success"
	case domain.ReservationStatusCancelada:
		title = "Reserva cancelada"
		description = "La reserva ha sido cancelada"
		if req.Notes != nil {
			description += ": " + *req.Notes
		}
		variant = "error"
	}

	timelineEvent := domain.TimelineEvent{
		ReservationID: id,
		Title:         title,
		Description:   description,
		At:            time.Now(),
		Variant:       variant,
		CreatedAt:     time.Now(),
	}

	if err := uc.reservationRepo.AddTimelineEvent(id, timelineEvent); err != nil {
		uc.logger.Warn("Failed to add timeline event for status change", zap.Error(err))
	}

	// Get updated reservation
	reservation, err := uc.reservationRepo.GetByID(id)
	if err != nil {
		uc.logger.Error("Failed to get updated reservation", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	uc.logger.Info("Reservation status changed successfully", zap.String("reservation_id", id))
	return reservation, nil
}

func (uc *ReservationUseCase) GetReservationTimeline(id string) ([]domain.TimelineEvent, error) {
	// Check if reservation exists
	_, err := uc.reservationRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.ErrReservationNotFound
		}
		uc.logger.Error("Failed to get reservation for timeline", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	timeline, err := uc.reservationRepo.GetTimeline(id)
	if err != nil {
		uc.logger.Error("Failed to get reservation timeline", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	return timeline, nil
}

// AssignDriver assigns a driver to a reservation
func (uc *ReservationUseCase) AssignDriver(reservationID string, driverID string) (*domain.Reservation, error) {
	uc.logger.Info("Assigning driver to reservation",
		zap.String("reservation_id", reservationID),
		zap.String("driver_id", driverID),
	)

	// Get reservation to check if it exists and is in valid state
	reservation, err := uc.reservationRepo.GetByID(reservationID)
	if err != nil {
		if err == domain.ErrReservationNotFound {
			return nil, domain.ErrReservationNotFound
		}
		uc.logger.Error("Failed to get reservation for driver assignment", zap.Error(err))
		return nil, err
	}

	// Check if reservation can have driver assigned (not completed or cancelled)
	if reservation.Status == domain.ReservationStatusCompletada || reservation.Status == domain.ReservationStatusCancelada {
		uc.logger.Warn("Attempted to assign driver to completed/cancelled reservation",
			zap.String("reservation_id", reservationID),
			zap.String("status", string(reservation.Status)),
		)
		return nil, domain.ErrInvalidInput
	}

	// Get driver information for timeline
	driver, err := uc.driverRepo.GetByID(driverID)
	if err != nil {
		if err == domain.ErrDriverNotFound {
			return nil, domain.ErrDriverNotFound
		}
		uc.logger.Error("Failed to get driver for assignment", zap.Error(err))
		return nil, err
	}

	// Assign driver using repository
	err = uc.reservationRepo.AssignDriver(reservationID, driverID)
	if err != nil {
		if err == domain.ErrDriverNotFound {
			return nil, domain.ErrDriverNotFound
		}
		uc.logger.Error("Failed to assign driver to reservation", zap.Error(err))
		return nil, err
	}

	// Get updated reservation with driver info
	updatedReservation, err := uc.reservationRepo.GetByID(reservationID)
	if err != nil {
		uc.logger.Error("Failed to get updated reservation after driver assignment", zap.Error(err))
		return nil, err
	}

	// Add timeline event with driver name
	timelineEvent := domain.TimelineEvent{
		ReservationID: reservationID,
		Title:         "Conductor asignado",
		Description:   fmt.Sprintf("Conductor %s %s asignado a la reserva", driver.FirstName, driver.LastName),
		At:            time.Now(),
		Variant:       "primary",
		CreatedAt:     time.Now(),
	}

	if err := uc.reservationRepo.AddTimelineEvent(reservationID, timelineEvent); err != nil {
		uc.logger.Warn("Failed to add timeline event for driver assignment", zap.Error(err))
		// Don't fail the whole operation for timeline error
	}

	uc.logger.Info("Driver assigned successfully to reservation",
		zap.String("reservation_id", reservationID),
		zap.String("driver_id", driverID),
		zap.String("driver_name", fmt.Sprintf("%s %s", driver.FirstName, driver.LastName)),
	)

	return updatedReservation, nil
}

func (uc *ReservationUseCase) AddTimelineEvent(id string, req domain.CreateTimelineEventRequest) error {
	// Check if reservation exists
	_, err := uc.reservationRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrReservationNotFound
		}
		uc.logger.Error("Failed to get reservation for timeline event", zap.Error(err))
		return domain.ErrInternalError
	}

	if req.Variant == "" {
		req.Variant = "default"
	}

	timelineEvent := domain.TimelineEvent{
		ReservationID: id,
		Title:         req.Title,
		Description:   req.Description,
		At:            time.Now(),
		Variant:       req.Variant,
		CreatedAt:     time.Now(),
	}

	if err := uc.reservationRepo.AddTimelineEvent(id, timelineEvent); err != nil {
		uc.logger.Error("Failed to add timeline event", zap.Error(err))
		return domain.ErrInternalError
	}

	return nil
}

func (uc *ReservationUseCase) DeleteReservation(id string) error {
	uc.logger.Info("Deleting reservation", zap.String("reservation_id", id))

	// Check if reservation exists
	reservation, err := uc.reservationRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.ErrReservationNotFound
		}
		uc.logger.Error("Failed to get reservation for deletion", zap.Error(err))
		return domain.ErrInternalError
	}

	// Only allow deletion of non-completed reservations
	if reservation.Status == domain.ReservationStatusCompletada {
		return domain.ErrInvalidInput
	}

	if err := uc.reservationRepo.Delete(id); err != nil {
		uc.logger.Error("Failed to delete reservation", zap.Error(err))
		return domain.ErrInternalError
	}

	uc.logger.Info("Reservation deleted successfully", zap.String("reservation_id", id))
	return nil
}

// Helper method to generate reservation ID
func (uc *ReservationUseCase) generateReservationID() string {
	// Simple implementation - in production you might want a more sophisticated approach
	timestamp := time.Now().Unix()
	return fmt.Sprintf("RSV-%s", strconv.FormatInt(timestamp%100000, 10))
}
