# Nucleus

Nucleus is the central control service of the platform. It synchronizes billing data from Stripe with user identity data from Clerk, manages user metadata, and determines which AI agents each user has access to. Acting as the source of truth, Nucleus ensures all user entitlements, usage limits, and service-level permissions are consistently maintained and distributed across the system.

## Features

- **Stripe Webhook Handler**: Processes Stripe webhook events for payment processing
- **Invoice Payment Processing**: Handles successful invoice payments and updates user product access
- **Product Access Management**: Automatically adds purchased product IDs to user metadata
- **Secure Webhook Verification**: Validates webhook signatures and IP addresses to ensure requests come from Stripe
- **Dynamic IP Validation**: Fetches current Stripe webhook IPs dynamically for enhanced security
- **Event Deduplication**: In-memory cache system prevents duplicate webhook processing with automatic cleanup
- **Clerk User Management**: Full CRUD operations for user identity management
- **Environment-based Configuration**: Uses environment variables for secure configuration management
- **Asynchronous Processing**: Webhook events are processed asynchronously for better performance

## Prerequisites

- Go 1.24.4 or higher
- Stripe account with webhook endpoint configured
- Clerk account for user identity management

## Installation

1. Clone the repository:
```bash
git clone https://github.com/SAMLA-io/nucleus.git
cd nucleus
```

2. Install dependencies:
```bash
go mod download
```

3. Create a `.env` file in the root directory with the following variables:
```env
STRIPE_KEY=sk_test_your_stripe_secret_key
STRIPE_WEBHOOK_SECRET=whsec_your_webhook_secret
PORT=8080
CLERK_SECRET_KEY=your_clerk_secret_key
```

## Configuration

### Stripe Webhook Setup

1. Go to your Stripe Dashboard
2. Navigate to Developers > Webhooks
3. Create a new webhook endpoint with your server URL (e.g., `https://yourdomain.com/webhook`)
4. Select the following events to listen for:
   - `invoice.paid`
5. Copy the webhook signing secret and add it to your `.env` file

### Clerk Setup

1. Go to your Clerk Dashboard
2. Navigate to API Keys
3. Copy your Secret Key and add it to your `.env` file

## Usage

### Running the Server

```bash
go run main.go
```

The server will start listening on the configured port (default: 8080).

### Webhook Endpoints

#### POST `/webhook`

Handles incoming Stripe webhook events.

**Supported Events:**
- `invoice.paid`: Triggered when an invoice is successfully paid, automatically adds the product ID to user metadata

**Security Features:**
- Request body size limit: 64KB
- Webhook signature verification
- Dynamic IP address validation against Stripe's current webhook IP list
- JSON payload validation
- Event deduplication with 30-second TTL and automatic cleanup

**Response Codes:**
- `200 OK`: Event processed successfully
- `400 Bad Request`: Invalid webhook signature or malformed JSON
- `403 Forbidden`: Request from non-Stripe IP address
- `503 Service Unavailable`: Error reading request body

### Product Access Management

When an invoice is paid, the service automatically:

1. Extracts the product ID from the invoice line items
2. Retrieves the current user metadata from Clerk
3. Adds the product ID to the user's `stripe.products_id` array in metadata
4. Updates the user's metadata in Clerk

This ensures users have immediate access to purchased products across the platform.

### Clerk User Management

The service includes full user management capabilities through Clerk:

- **Get User Metadata**: `clerk.GetUserMetadata(userId)`
- **Update User Metadata**: `clerk.UpdateUserMetadata(userId, metadata)`

## Development

### Project Structure

```
nucleus/
├── main.go              # Main application entry point
├── go.mod               # Go module dependencies
├── go.sum               # Go module checksums
├── README.md            # This file
├── .gitignore           # Git ignore rules
├── cache/
│   └── store.go         # Event cache management with automatic cleanup
├── clerk/
│   └── metadata.go      # Clerk user metadata management
├── stripe/
│   ├── address.go       # Dynamic webhook IP validation
│   └── webhook.go       # Stripe webhook handler and invoice processing
└── types/
    └── cache/
        └── cache_types.go # Cache data structures and cleanup logic
```

### Cache System

The application includes a sophisticated in-memory cache system for webhook event deduplication:

- **Automatic Cleanup**: Expired entries are removed every 30 seconds via background goroutine
- **Thread-safe**: Uses read-write mutex for concurrent access
- **TTL**: 30-second expiration for cache entries
- **Statistics**: Get cache size and entry information
- **Memory Efficient**: Automatically removes expired entries to prevent memory leaks

### Webhook IP Validation

The service dynamically fetches the current list of Stripe webhook IPs from Stripe's official endpoint (`https://stripe.com/files/ips/ips_webhooks.json`) to ensure the most up-to-date security validation.

### Adding New Webhook Events

To handle additional Stripe webhook events:

1. Add a new case in the switch statement in `processWebhookEvent()`:
```go
case "your.event.type":
    var eventData YourEventStruct
    jsonData, err := json.Marshal(event.Data.Object)
    if err != nil {
        log.Printf("Error marshaling webhook data: %v", err)
        return
    }
    err = json.Unmarshal(jsonData, &eventData)
    if err != nil {
        log.Printf("Error parsing webhook JSON: %v", err)
        return
    }
    // Handle the event
```

2. Implement the corresponding handler function
3. Update your Stripe webhook configuration to listen for the new event

### Testing

For local development, you can use tools like:
- [Stripe CLI](https://stripe.com/docs/stripe-cli) for webhook forwarding
- [ngrok](https://ngrok.com/) for exposing your local server

Example with Stripe CLI:
```bash
stripe listen --forward-to localhost:8080/webhook
```

## Security Considerations

- Always verify webhook signatures using the provided secret
- Dynamic webhook IP validation ensures requests only come from current Stripe IPs
- Use HTTPS in production
- Keep your Stripe and Clerk keys secure and never commit them to version control
- Implement proper error handling and logging
- Consider rate limiting for webhook endpoints
- Event deduplication prevents replay attacks
- Automatic cache cleanup prevents memory exhaustion

## Dependencies

- `github.com/clerk/clerk-sdk-go/v2`: Clerk user management SDK
- `github.com/joho/godotenv`: Environment variable management
- `github.com/stripe/stripe-go/v82`: Stripe Go SDK

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Contributors

- [@jpgtzg](https://github.com/jpgtzg)