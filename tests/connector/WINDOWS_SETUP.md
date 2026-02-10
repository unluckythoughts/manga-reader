# Windows Setup Guide for Connector Tests

## Problem

When running connector tests on Windows, you may encounter this error:

```
panic: fork/exec leakless.exe: Operation did not complete successfully 
because the file contains a virus or potentially unwanted software.
```

This happens because Windows Defender (or other antivirus software) blocks the `leakless.exe` file that's part of the **go-rod** browser automation library. The connector tests use go-rod to bypass Cloudflare protection on websites.

## Solution 1: Add Windows Defender Exclusions (Recommended)

### Step-by-Step Instructions:

1. **Open Windows Security**
   - Press `Win + I` to open Settings
   - Go to "Privacy & Security" → "Windows Security"
   - Click "Virus & threat protection"

2. **Add Exclusions**
   - Scroll down to "Virus & threat protection settings"
   - Click "Manage settings"
   - Scroll to "Exclusions"
   - Click "Add or remove exclusions"
   - Click "Add an exclusion" → "Folder"

3. **Add These Folders:**
   
   **Folder 1 - Rod Browser Cache:**
   ```
   C:\Users\<YourUsername>\AppData\Roaming\rod
   ```
   
   **Folder 2 - Leakless Temp Directory:**
   ```
   C:\Users\<YourUsername>\AppData\Local\Temp
   ```
   
   Replace `<YourUsername>` with your actual Windows username (in your case: `Vinay`).

4. **Run Tests Again**
   ```powershell
   make test-connector
   ```

## Solution 2: Run Tests in Short Mode

If you don't want to modify Windows Defender settings, you can run tests in "short mode" which skips actual network requests:

```powershell
make test-connector-short
```

This will run all the unit tests that validate the connector configuration without making actual HTTP requests to websites.

## Solution 3: Temporarily Disable Real-Time Protection

⚠️ **Not Recommended** - Only use during testing

1. Open Windows Security
2. Go to "Virus & threat protection"
3. Click "Manage settings"
4. Turn off "Real-time protection" (temporarily)
5. Run your tests
6. **Remember to turn it back on!**

## Why Does This Happen?

The `leakless.exe` file is a legitimate utility used by go-rod to prevent resource leaks when spawning browser processes. However, because it:

- Downloads executable files
- Spawns processes
- Monitors other processes

Antivirus software may flag it as potentially unwanted software (PUP). It's a false positive.

## Verifying the Fix

After adding the exclusions, run:

```powershell
make test-connector
```

You should see output like:

```
=== RUN   TestFreeWebNovelTestSuite
=== RUN   TestFreeWebNovelTestSuite/TestConnectorProperties
=== RUN   TestFreeWebNovelTestSuite/TestGetName
...
```

If tests skip with a message about browser automation being blocked, the exclusions weren't added correctly.

## Alternative Testing Approach

If you're unable to modify Windows Defender settings (e.g., on a corporate machine), consider:

1. **Using WSL (Windows Subsystem for Linux)**
   ```bash
   # In WSL
   cd /mnt/c/Users/Vinay/Documents/workspace/book-reader
   make test-connector
   ```

2. **Running tests in a Docker container**
   ```powershell
   docker run --rm -v ${PWD}:/app -w /app golang:1.24 go test ./tests/connector/... -v
   ```

3. **Using short mode for CI/CD**
   ```powershell
   make test-connector-short
   ```

## Still Having Issues?

If you're still experiencing problems:

1. Check your antivirus software logs
2. Make sure you added exclusions for the correct folders
3. Try restarting your terminal/PowerShell
4. Check the [go-rod documentation](https://go-rod.github.io/) for platform-specific issues

## Security Note

The go-rod and leakless libraries are open-source and widely used in the Go community. You can review their source code:

- go-rod: https://github.com/go-rod/rod
- leakless: https://github.com/ysmood/leakless

Adding exclusions for these specific folders is safe and won't significantly impact your system's security.
