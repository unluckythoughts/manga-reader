package connector

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// ConnectorTestSuite provides test infrastructure for connector integration tests
type ConnectorTestSuite struct {
	suite.Suite
}

// SetupSuite runs once before all tests in the suite
func (s *ConnectorTestSuite) SetupSuite() {
	// Setup code if needed
}

// SetupTest runs before each test
func (s *ConnectorTestSuite) SetupTest() {
	// Additional setup if needed
}

// TearDownTest runs after each test
func (s *ConnectorTestSuite) TearDownTest() {
	// Cleanup if needed
}

// TearDownSuite runs once after all tests in the suite
func (s *ConnectorTestSuite) TearDownSuite() {
	// Cleanup if needed
}

// Helper function to skip tests when network is unavailable
func (s *ConnectorTestSuite) SkipIfShort(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}
}
