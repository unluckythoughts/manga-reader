package tests

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/unluckythoughts/book-reader/client"
	"github.com/unluckythoughts/go-microservice/tools/logger"
	"github.com/unluckythoughts/go-microservice/tools/psql"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

// TestSuite provides test infrastructure using testify suite
type TestSuite struct {
	suite.Suite
	DB     *gorm.DB
	Logger *zap.Logger
	Client *client.Client
}

func getLogger() *zap.Logger {
	l := logger.New(logger.Options{
		LogLevel: zap.DebugLevel.String(),
	})
	l = l.With(zap.String("env", "TEST"))
	return l
}

func getDB(l *zap.Logger) *gorm.DB {
	return psql.New(psql.Options{
		Logger:   l,
		Host:     "localhost",
		Port:     5432,
		Name:     "book_reader",
		User:     "test",
		Password: "test",
		SSLMode:  "disable",
		Debug:    false,
	})
}

// SetupSuite runs once before all tests in the suite
func (suite *TestSuite) SetupSuite() {
	// Creating logger
	suite.Logger = getLogger()

	// Create database connection for testing
	suite.DB = getDB(suite.Logger)

	// Create HTTP client pointing to the running backend server
	suite.Client = client.NewClient("http://localhost:8080")
}

// SetupTest runs before each test
func (suite *TestSuite) SetupTest() {
	// additional setup if needed
	suite.Client.Log(suite.T().Name())
}

// TearDownTest runs after each test
func (suite *TestSuite) TearDownTest() {
	// Additional cleanup if needed
	suite.Client.Log("")
}

// TearDownSuite runs once after all tests in the suite
func (suite *TestSuite) TearDownSuite() {
	// Close database connection
	if suite.DB != nil {
		sqlDB, err := suite.DB.DB()
		if err == nil {
			sqlDB.Close()
		}
	}
}

// Helper function to create a test suite for individual test files
func NewTestSuite(t *testing.T) *TestSuite {
	testSuite := &TestSuite{}
	suite.Run(t, testSuite)
	return testSuite
}
