package assetservice

import "fin4-core/server/dbservice"

//Service defines asset service type
type Service struct {
	*dbservice.DB
}

//NewService factory for Service
func NewService(db *dbservice.DB) *Service {
	return &Service{db}
}
