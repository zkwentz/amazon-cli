# Troubleshooting Guide

This guide covers common issues you may encounter when using the Amazon CLI tool and their solutions.

## Common Issues

### Authentication Failed

**Problem:** You receive an authentication failed error when trying to execute commands.

**Solution:** Re-run the authentication login command to refresh your credentials:

```bash
amazon-cli auth login
```

This will prompt you to log in again and update your stored credentials.

---

### Rate Limited

**Problem:** You're receiving rate limit errors from Amazon's servers.

**Solution:** Wait for the rate limit period to expire and then retry your request. Amazon imposes rate limits to protect their services. The rate limit typically resets after a short period (usually a few minutes to an hour depending on the endpoint).

You can:
- Wait a few minutes before retrying
- Reduce the frequency of your requests
- Consider batching operations if possible

---

### CAPTCHA Required

**Problem:** A CAPTCHA challenge is required to complete your request.

**Solution:** Complete the CAPTCHA verification in your web browser:

1. Open your web browser and navigate to the Amazon website
2. Log in with your Amazon account credentials
3. Complete the CAPTCHA challenge when prompted
4. After successful completion, return to the CLI and retry your command

Note: CAPTCHA challenges are typically triggered by unusual activity patterns or security measures.

---

### Command Not Found

**Problem:** When running `amazon-cli` commands, you receive a "command not found" error.

**Solution:** Check that the installation directory is in your system's PATH:

**For Unix/Linux/macOS:**
```bash
# Check if the binary is in your PATH
which amazon-cli

# If not found, add the installation directory to your PATH
# Add this line to your ~/.bashrc, ~/.zshrc, or equivalent shell config file
export PATH="$PATH:/path/to/amazon-cli"

# Then reload your shell configuration
source ~/.bashrc  # or source ~/.zshrc
```

**For Windows:**
1. Open System Properties → Advanced → Environment Variables
2. Under "User variables" or "System variables", find the PATH variable
3. Add the directory containing `amazon-cli.exe` to the PATH
4. Restart your command prompt or terminal

---

### Permission Denied on Config

**Problem:** You receive a "permission denied" error when the tool tries to access configuration files.

**Solution:** Check and fix the file permissions for your configuration directory and files:

**For Unix/Linux/macOS:**
```bash
# Check current permissions
ls -la ~/.amazon-cli/

# Fix permissions for the config directory (owner read/write/execute only)
chmod 700 ~/.amazon-cli/

# Fix permissions for config files (owner read/write only)
chmod 600 ~/.amazon-cli/config
chmod 600 ~/.amazon-cli/credentials

# Ensure you own the files
# Replace 'username' with your actual username
chown -R username:username ~/.amazon-cli/
```

**For Windows:**
1. Right-click on the config directory (typically in `%USERPROFILE%\.amazon-cli\`)
2. Select "Properties" → "Security" tab
3. Ensure your user account has "Full control" permissions
4. Apply the changes to all files in the directory

---

## Still Having Issues?

If you continue to experience problems after trying these solutions:

1. Check the [GitHub Issues](https://github.com/yourusername/amazon-cli/issues) page for similar problems
2. Review the application logs for more detailed error messages
3. Ensure you're using the latest version of the tool
4. Create a new issue with detailed information about your problem, including:
   - The exact command you're running
   - The complete error message
   - Your operating system and version
   - The tool version (`amazon-cli version`)
