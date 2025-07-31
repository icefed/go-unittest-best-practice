package api

import (
	"go-unittest-best-practice/internal/config"
	"go-unittest-best-practice/internal/store"
	"testing"

	"github.com/stretchr/testify/suite"
	"go.uber.org/mock/gomock"
)

type ServiceTestSuite struct {
	suite.Suite

	ctrl         *gomock.Controller
	conf         *config.Config
	mockUserRepo *store.MockUserRepository
	svc          *Service
}

func (s *ServiceTestSuite) SetupSuite() {
	s.conf = &config.Config{}
}

func (s *ServiceTestSuite) TearDownSuite() {
}

func (s *ServiceTestSuite) SetupTest() {
	s.ctrl = gomock.NewController(s.T())
	s.mockUserRepo = store.NewMockUserRepository(s.ctrl)
	s.svc = NewService(s.mockUserRepo, s.conf)
}

func (s *ServiceTestSuite) TearDownTest() {
	s.ctrl.Finish()
}

func TestService(t *testing.T) {
	suite.Run(t, new(ServiceTestSuite))
}