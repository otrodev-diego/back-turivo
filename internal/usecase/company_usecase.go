package usecase

import (
	"github.com/google/uuid"
	"go.uber.org/zap"

	"turivo-backend/internal/domain"
)

type CompanyUseCase struct {
	companyRepo domain.CompanyRepository
	logger      *zap.Logger
}

func NewCompanyUseCase(
	companyRepo domain.CompanyRepository,
	logger *zap.Logger,
) *CompanyUseCase {
	return &CompanyUseCase{
		companyRepo: companyRepo,
		logger:      logger,
	}
}

func (uc *CompanyUseCase) ListCompanies(req domain.ListCompaniesRequest) ([]*domain.Company, int, error) {
	// Set default pagination
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 || req.PageSize > 100 {
		req.PageSize = 20
	}
	if req.Sort == "" {
		req.Sort = "created_at"
	}

	companies, total, err := uc.companyRepo.List(req)
	if err != nil {
		uc.logger.Error("Failed to list companies", zap.Error(err))
		return nil, 0, domain.ErrInternalError
	}

	return companies, total, nil
}

func (uc *CompanyUseCase) CreateCompany(req domain.CreateCompanyRequest) (*domain.Company, error) {
	// Check if company with RUT already exists
	existing, err := uc.companyRepo.GetByRUT(req.RUT)
	if err != nil && err != domain.ErrNotFound {
		uc.logger.Error("Failed to check existing company by RUT", zap.Error(err))
		return nil, domain.ErrInternalError
	}
	if existing != nil {
		uc.logger.Warn("Company with RUT already exists", zap.String("rut", req.RUT))
		return nil, domain.ErrAlreadyExists
	}

	company := &domain.Company{
		Name:         req.Name,
		RUT:          req.RUT,
		ContactEmail: req.ContactEmail,
		Status:       req.Status,
		Sector:       req.Sector,
	}

	if err := uc.companyRepo.Create(company); err != nil {
		uc.logger.Error("Failed to create company", zap.Error(err))
		return nil, domain.ErrInternalError
	}

	uc.logger.Info("Company created successfully",
		zap.String("id", company.ID.String()),
		zap.String("name", company.Name))

	return company, nil
}

func (uc *CompanyUseCase) GetCompanyByID(id uuid.UUID) (*domain.Company, error) {
	company, err := uc.companyRepo.GetByID(id)
	if err != nil {
		uc.logger.Error("Failed to get company by ID", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	return company, nil
}

func (uc *CompanyUseCase) UpdateCompany(id uuid.UUID, req domain.UpdateCompanyRequest) (*domain.Company, error) {
	// Check if company exists
	existing, err := uc.companyRepo.GetByID(id)
	if err != nil {
		uc.logger.Error("Failed to get company for update", zap.Error(err), zap.String("id", id.String()))
		return nil, err
	}

	// Check if RUT is being changed and already exists
	if req.RUT != nil && *req.RUT != existing.RUT {
		rutExists, err := uc.companyRepo.GetByRUT(*req.RUT)
		if err != nil && err != domain.ErrNotFound {
			uc.logger.Error("Failed to check existing company by RUT", zap.Error(err))
			return nil, domain.ErrInternalError
		}
		if rutExists != nil {
			uc.logger.Warn("Company with RUT already exists", zap.String("rut", *req.RUT))
			return nil, domain.ErrAlreadyExists
		}
	}

	company, err := uc.companyRepo.Update(id, req)
	if err != nil {
		uc.logger.Error("Failed to update company", zap.Error(err), zap.String("id", id.String()))
		return nil, domain.ErrInternalError
	}

	uc.logger.Info("Company updated successfully", zap.String("id", id.String()))

	return company, nil
}

func (uc *CompanyUseCase) DeleteCompany(id uuid.UUID) error {
	// Check if company exists
	_, err := uc.companyRepo.GetByID(id)
	if err != nil {
		uc.logger.Error("Failed to get company for deletion", zap.Error(err), zap.String("id", id.String()))
		return err
	}

	if err := uc.companyRepo.Delete(id); err != nil {
		uc.logger.Error("Failed to delete company", zap.Error(err), zap.String("id", id.String()))
		return domain.ErrInternalError
	}

	uc.logger.Info("Company deleted successfully", zap.String("id", id.String()))

	return nil
}
