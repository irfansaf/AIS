# Setup Guide for Bitbucket-Assistant

This guide will walk you through setting up the Bitbucket-Assistant application to automatically comment on pull requests.

## Prerequisites

1. **Bitbucket Server Access**: You need access to a Bitbucket Server instance
2. **User Credentials**: A Bitbucket user account with repository access
3. **Go Environment**: Go 1.24.4 or later installed

## Step-by-Step Setup

### 1. Environment Configuration

Create a `.env` file in the project root with the following configuration:

```env
# Bitbucket Server Configuration
BITBUCKET_URL=https://your-bitbucket-server.com
BITBUCKET_USERNAME=your_username
BITBUCKET_PASSWORD=your_password_or_token

# Webhook Security (generate a random string)
BITBUCKET_WEBHOOK_SECRET=your_secure_random_string

# Server Configuration
PORT=8080
```

#### Important Notes:

- **BITBUCKET_URL**: The base URL of your Bitbucket Server (without trailing slash)
- **BITBUCKET_USERNAME**: Your Bitbucket username
- **BITBUCKET_PASSWORD**: Your password or Personal Access Token (recommended)
- **BITBUCKET_WEBHOOK_SECRET**: A secure random string for webhook validation

### 2. Creating a Personal Access Token (Recommended)

Instead of using your password, create a Personal Access Token:

1. Log into your Bitbucket Server
2. Go to **Personal Settings** → **Personal Access Tokens**
3. Click **Create Token**
4. Set permissions:
   - **Repository**: Read
   - **Pull Request**: Write (to add comments)
5. Copy the generated token and use it as `BITBUCKET_PASSWORD`

### 3. Webhook Configuration

#### Repository-Level Webhook

1. Navigate to your repository in Bitbucket
2. Go to **Repository Settings** → **Webhooks**
3. Click **Create Webhook**
4. Configure:
   - **Name**: `Bitbucket Assistant`
   - **URL**: `http://your-server:8080/api/v1/webhook/bitbucket`
   - **Secret**: Same value as `BITBUCKET_WEBHOOK_SECRET`
   - **Events**: Select "Pull Request" events:
     - Pull request opened
     - Pull request source branch updated
   - **Active**: ✅ Enabled

#### Project-Level Webhook (Optional)

For multiple repositories, you can set up a project-level webhook:

1. Navigate to your project in Bitbucket
2. Go to **Project Settings** → **Webhooks**
3. Follow the same configuration as above

### 4. Testing the Setup

#### 1. Start the Application

```bash
# Development
go run cmd/main.go

# Or build and run
go build -o bitbucket-assistant cmd/main.go
./bitbucket-assistant
```

#### 2. Verify Health Check

```bash
curl http://localhost:8080/health
```

Expected response:
```json
{
  "status": "healthy",
  "service": "bitbucket-assistant",
  "timestamp": "2024-01-01T12:00:00Z"
}
```

#### 3. Test Webhook

1. Create a test pull request in your repository
2. Check the application logs for processing messages
3. Verify that a welcome comment appears on the pull request

### 5. Production Deployment

#### Using Docker

1. Build the Docker image:
   ```bash
   docker compose build
   ```

2. Start the service:
   ```bash
   docker compose up -d
   ```

3. Check logs:
   ```bash
   docker compose logs -f bitbucket-assistant
   ```

#### Using Systemd (Linux)

1. Create a systemd service file `/etc/systemd/system/bitbucket-assistant.service`:
   ```ini
   [Unit]
   Description=Bitbucket Assistant
   After=network.target

   [Service]
   Type=simple
   User=bitbucket
   WorkingDirectory=/opt/bitbucket-assistant
   ExecStart=/opt/bitbucket-assistant/bitbucket-assistant
   EnvironmentFile=/opt/bitbucket-assistant/.env
   Restart=always
   RestartSec=5

   [Install]
   WantedBy=multi-user.target
   ```

2. Enable and start the service:
   ```bash
   sudo systemctl enable bitbucket-assistant
   sudo systemctl start bitbucket-assistant
   ```

### 6. Troubleshooting

#### Common Issues

**Webhook not triggering:**
- Verify the webhook URL is accessible from Bitbucket Server
- Check firewall settings
- Ensure the application is running and listening on the correct port

**Authentication errors:**
- Verify credentials are correct
- Check if Personal Access Token has expired
- Ensure user has repository access permissions

**Comments not appearing:**
- Check application logs for API errors
- Verify user has permission to comment on pull requests
- Ensure Bitbucket API endpoint is correct

#### Log Analysis

The application provides detailed logging:

```bash
# View logs in real-time
tail -f /var/log/bitbucket-assistant.log

# Or with Docker
docker compose logs -f bitbucket-assistant
```

Look for these log patterns:
- `INFO`: Successful operations
- `DEBUG`: Detailed processing information
- `ERROR`: Issues that need attention
- `FATAL`: Critical errors

### 7. Security Considerations

1. **Use HTTPS**: Always use HTTPS in production
2. **Secure Secrets**: Store sensitive data in environment variables
3. **Network Security**: Restrict access to the webhook endpoint
4. **Token Rotation**: Regularly rotate Personal Access Tokens
5. **Monitoring**: Set up monitoring and alerting

### 8. Customization

To customize the welcome comment template, modify the `addWelcomeComment` function in `internal/application/webhook_service.go`.

Example customization:
```go
commentText := fmt.Sprintf(
    "🚀 **Hello %s!**\n\n"+
        "Your custom welcome message here...\n"+
        "Project: %s\n"+
        "Repository: %s",
    event.PullRequest.Author.User.DisplayName,
    projectKey,
    repoSlug,
)
```

## Support

For issues and questions:
1. Check the application logs
2. Review this setup guide
3. Consult the main README.md
4. Create an issue in the project repository