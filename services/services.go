package services

import (
	"fmt"
	"myproject/admin-api-gateway/config"
	pbh "myproject/admin-api-gateway/genproto/healthcare-service"
	pbu "myproject/admin-api-gateway/genproto/user-service"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/resolver"
)

type IServiceManager interface {
	UserService() pbu.UserServiceClient
	HealthCareService() pbh.HealthcareServiceClient
}

type serviceManager struct {
	userService       pbu.UserServiceClient
	healthcareService pbh.HealthcareServiceClient
}

func (s *serviceManager) UserService() pbu.UserServiceClient {
	return s.userService
}

func (s *serviceManager) HealthCareService() pbh.HealthcareServiceClient {
	return s.healthcareService
}

func NewServiceManager(cfg config.Config) (IServiceManager, error) {
	resolver.SetDefaultScheme("dns")

	connUser, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.UserServiceHost, cfg.UserServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("user service dial error, %s:%d:%v", cfg.UserServiceHost, cfg.UserServicePort, err)
	}

	connHealthcare, err := grpc.Dial(
		fmt.Sprintf("%s:%d", cfg.HealthcareServiceHost, cfg.HealthcareServicePort),
		grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("healthcare service dial error, %s:%d:%v", cfg.UserServiceHost, cfg.UserServicePort, err)
	}

	return &serviceManager{userService: pbu.NewUserServiceClient(connUser), healthcareService: pbh.NewHealthcareServiceClient(connHealthcare)}, nil
}
