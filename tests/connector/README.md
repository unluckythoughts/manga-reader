# Connector Integration Tests

This folder contains integration tests for various book reader connectors.

## Overview

These tests verify that connectors can properly scrape and parse content from their respective websites. They test the complete workflow:

1. **Book Discovery** - Getting lists of available books
2. **Book Details** - Extracting book metadata (title, cover, synopsis)
3. **Chapter Lists** - Getting chapters for a specific book
4. **Chapter Content** - Extracting readable content from chapters

## Running Tests

### Run All Connector Tests

```bash
go test ./tests/connector/...
```

### Run Tests for a Specific Connector

```bash
# FreeWebNovel connector tests
go test ./tests/connector -run TestFreeWebNovel
```

### Skip Long-Running Tests

Use the `-short` flag to skip actual network requests:

```bash
go test ./tests/connector/... -short
```

### Run with Verbose Output

```bash
go test ./tests/connector/... -v
```

### Run Benchmarks

```bash
go test ./tests/connector -bench=. -benchmem
```

## Test Structure

### FreeWebNovel Connector Tests

The `freewebnovel_test.go` file contains comprehensive tests for the FreeWebNovel connector:

#### Unit Tests
- `TestConnectorProperties` - Verifies connector configuration
- `TestGetName` - Validates connector name
- `TestGetDomain` - Validates connector domain
- `TestGetIconURL` - Validates icon URL
- `TestGetSelectors` - Validates CSS selectors configuration

#### Integration Tests
- `TestGetBooks` - Tests fetching book lists
- `TestGetBookCount` - Tests counting total books
- `TestGetBookChapters` - Tests fetching chapter lists
- `TestGetBookChaptersWithFilter` - Tests chapter filtering by number
- `TestGetChapterContent` - Tests fetching chapter content
- `TestIntegrationFlow` - Tests complete workflow from books to content

#### Benchmark Tests
- `BenchmarkGetBooks` - Measures performance of book list fetching
- `BenchmarkGetBookCount` - Measures performance of book count operation

## Important Notes

### Network Dependency

⚠️ **These are integration tests that make real HTTP requests to external websites.**

- Tests may fail if:
  - The target website is down
  - The website structure has changed
  - Network connectivity issues
  - Rate limiting or bot detection

### Cloudflare Protection

The connectors use **two-layer Cloudflare bypass** for maximum success:

#### Layer 1: cloudflare-bp-go (Preventive)
- Configures proper TLS settings and HTTP headers
- Prevents Cloudflare from triggering challenges
- Fast, lightweight approach
- No browser needed

#### Layer 2: rod Browser (Reactive)  
- Launches real Chromium browser when needed
- Waits for Cloudflare challenges to auto-solve (up to 53 seconds)
- Solves challenges that Layer 1 couldn't prevent
- Slower but very effective

**FreeWebNovel uses both layers** for maximum reliability against aggressive Cloudflare protection.

### Windows-Specific Issues

⚠️ **Windows Defender may block browser automation (go-rod library)**

If you see an error about `leakless.exe` or "virus or potentially unwanted software", see **[WINDOWS_SETUP.md](WINDOWS_SETUP.md)** for detailed setup instructions.

**Quick Fix:**
```powershell
# Run tests without network requests
make test-connector-short
```

**Or** add Windows Defender exclusions for:
- `C:\Users\<YourUsername>\AppData\Roaming\rod`
- `C:\Users\<YourUsername>\AppData\Local\Temp`

The tests will automatically detect and skip with a helpful message if this issue occurs.

### Test Behavior

- Tests will **skip gracefully** if network requests fail
- Use `-short` flag during development to run only quick tests
- Integration tests log detailed information about failures

### Site Changes

If tests start failing consistently:

1. Check if the website is accessible
2. Verify the CSS selectors in the connector configuration
3. Check if the website structure has changed
4. Update selectors in the connector if needed

## Adding New Connector Tests

To add tests for a new connector:

1. Create a new test file: `<connector_name>_test.go`
2. Use the `ConnectorTestSuite` from `suite.go`
3. Implement similar test methods as `freewebnovel_test.go`
4. Follow the pattern:
   - Test basic properties first
   - Test each interface method
   - Add an integration flow test
   - Add benchmarks for performance-critical operations

Example structure:

```go
type NewConnectorTestSuite struct {
    ConnectorTestSuite
    connector models.IConnector
}

func (s *NewConnectorTestSuite) SetupSuite() {
    s.connector = connector.NewYourConnector()
}

func (s *NewConnectorTestSuite) RunTests(t *testing.T) {
    t.Run("TestGetBooks", s.TestGetBooks)
    // Add more tests...
}
```

## Continuous Integration

When running in CI/CD:

```bash
# Run with timeout to prevent hanging tests
go test ./tests/connector/... -timeout 5m -v

# Run only short tests for quick feedback
go test ./tests/connector/... -short
```

## Debugging Tips

1. **Enable verbose mode** to see detailed test output:
   ```bash
   go test ./tests/connector -v -run TestGetChapterContent
   ```

2. **Check specific test logs** to see sample data being fetched

3. **Use browser developer tools** to verify selectors if tests fail

4. **Test individual methods** by running specific test functions:
   ```bash
   go test ./tests/connector -run TestFreeWebNovel/TestGetBooks
   ```

## Coverage

Generate test coverage report:

```bash
go test ./tests/connector/... -coverprofile=coverage.out
go tool cover -html=coverage.out
```
