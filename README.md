# Bedrock-Bitbucket-Assistant

A Go-based webhook listener for Bitbucket Server that automatically adds welcome comments to newly opened pull requests. Built using hexagonal architecture principles with clean separation of concerns.

## Features

- 🎯 **Webhook Processing**: Listens to Bitbucket Server webhook events
- 💬 **Automatic Comments**: Adds personalized welcome comments to new pull requests
- 🔒 **Secure**: HMAC signature validation for webhook security
- 🏗️ **Clean Architecture**: Implements hexagonal architecture with SOLID principles
- 🚀 **High Performance**: Built with Fiber web framework and Sonic JSON parser

## Architecture

This application follows hexagonal architecture (ports and adapters) pattern:

```
├── cmd/                    # Application entry point
├── internal/
│   ├── core/
│   │   ├── domain/        # Business entities and domain logic
│   │   └── ports/         # Interfaces (ports)
│   ├── application/       # Use cases and business logic
│   └── adapters/          # External adapters
│       ├── http/          # HTTP handlers (driving adapter)
│       └── bitbucket/     # Bitbucket API client (driven adapter)
└── pkg/                   # Shared utilities
```

## Prerequisites

- Go 1.24.4 or later
- Bitbucket Server instance
- Valid Bitbucket user credentials with repository access

## Setup

### 1. Clone and Install Dependencies

```bash
git clone <repository-url>
cd AIS
go mod download
```

### 2. Environment Configuration

Copy the sample environment file and configure your settings:

```bash
cp .env.sample .env
```

Edit `.env` with your Bitbucket Server details:

```env
# Bitbucket Server Configuration
BITBUCKET_URL=https://your-bitbucket-server.com
BITBUCKET_USERNAME=your_username
BITBUCKET_PASSWORD=your_password_or_personal_access_token

# Webhook Security
BITBUCKET_WEBHOOK_SECRET=your_webhook_secret

# Server Configuration
PORT=8080
```

### 3. Bitbucket Webhook Configuration

1. Navigate to your Bitbucket repository settings
2. Go to **Webhooks** section
3. Create a new webhook with:
   - **URL**: `http://your-server:8080/api/v1/webhook/bitbucket`
   - **Secret**: Same as `BITBUCKET_WEBHOOK_SECRET` in your `.env`
   - **Events**: Select "Pull Request" events
   - **SSL Verification**: Enable if using HTTPS

## Running the Application

### Development

```bash
go run cmd/main.go
```

### Production

```bash
go build -o bitbucket-assistant cmd/main.go
./bitbucket-assistant
```

### Using Docker

```bash
docker compose up --build
```

The application will be available at `http://localhost:8080`

## API Endpoints

### Webhook Endpoint

```
POST /api/v1/webhook/bitbucket
```

**Headers:**
- `X-Event-Key`: Bitbucket event type (e.g., `pr:opened`)
- `X-Hub-Signature`: HMAC signature for security validation
- `Content-Type`: `application/json`

**Supported Events:**
- `pr:opened` - Automatically adds welcome comment
- `pr:from_ref_updated` - Logs update information

## Welcome Comment Template

When a new pull request is opened, the application automatically adds a personalized comment:

```markdown
🎉 **Welcome [Author Name]!**

Thank you for opening this pull request. Here are some quick reminders:

✅ **Code Review Checklist:**
- [ ] Code follows project coding standards
- [ ] Tests are included and passing
- [ ] Documentation is updated if needed
- [ ] No sensitive information is exposed

📋 **PR Details:**
- **Title:** [PR Title]
- **Source Branch:** `[source-branch]`
- **Target Branch:** `[target-branch]`

A reviewer will be assigned shortly. Thank you for your contribution! 🚀
```

## Security

- **HMAC Validation**: All webhook requests are validated using HMAC-SHA256
- **Environment Variables**: Sensitive data stored in environment variables
- **Basic Authentication**: Bitbucket API calls use basic authentication
- **Timeout Protection**: HTTP client configured with 30-second timeout

## Extending the Application

### Adding New Event Handlers

1. Update the `ProcessPullRequestEvent` method in `internal/application/webhook_service.go`
2. Add new case statements for additional event types
3. Implement corresponding business logic

### Adding New Bitbucket API Operations

1. Extend the `BitbucketClient` interface in `internal/core/ports/bitbucket_client.go`
2. Implement the new methods in `internal/adapters/bitbucket/client.go`
3. Use the new operations in your application services

## Troubleshooting

### Common Issues

1. **Webhook not triggering**:
   - Verify webhook URL is accessible from Bitbucket Server
   - Check webhook secret matches environment variable
   - Ensure correct events are selected in webhook configuration

2. **Authentication errors**:
   - Verify Bitbucket credentials are correct
   - Ensure user has appropriate repository permissions
   - Check if using personal access token instead of password

3. **Comment not appearing**:
   - Check application logs for API errors
   - Verify Bitbucket API endpoint URL format
   - Ensure user has permission to comment on pull requests

### Logging

The application uses structured logging. Key log levels:
- `INFO`: Successful operations
- `DEBUG`: Detailed processing information
- `ERROR`: Error conditions that don't stop processing
- `FATAL`: Critical errors that stop the application

## Contributing

1. Follow Go coding standards and conventions
2. Maintain hexagonal architecture principles
3. Add tests for new functionality
4. Update documentation for new features
5. Ensure SOLID principles are followed

## License

This project is licensed under the MIT License.