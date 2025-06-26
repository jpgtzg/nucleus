# Nucleus

Nucleus is the central control service of the platform. It synchronizes billing data from Stripe with user identity data from Clerk, manages user metadata, and determines which AI agents each user has access to. Acting as the source of truth, Nucleus ensures all user entitlements, usage limits, and service-level permissions are consistently maintained and distributed across the system.

## Features

- **Stripe Webhook Handler**: Processes Stripe webhook events for payment processing
- **Payment Intent Management**: Handles successful payment intents and payment method attachments
- **Secure Webhook Verification**: Validates webhook signatures to ensure requests come from Stripe
- **Environment-based Configuration**: Uses environment variables for secure configuration management

## Prerequisites

- Go 1.19 or higher
- Stripe account with webhook endpoint configured
- Clerk account for user identity management

## Installation

1. Clone the repository:
```bash
git clone <repository-url>
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
```

## Configuration

### Environment Variables

| Variable | Description | Required |
|----------|-------------|----------|
| `STRIPE_KEY` | Your Stripe secret key | Yes |
| `STRIPE_WEBHOOK_SECRET` | Webhook endpoint secret from Stripe | Yes |
| `PORT` | Port number for the server to listen on | No (defaults to 8080) |

### Stripe Webhook Setup

1. Go to your Stripe Dashboard
2. Navigate to Developers > Webhooks
3. Create a new webhook endpoint with your server URL (e.g., `https://yourdomain.com/webhook`)
4. Select the following events to listen for:
   - `payment_intent.succeeded`
   - `payment_method.attached`
5. Copy the webhook signing secret and add it to your `.env` file

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
- `payment_intent.succeeded`: Triggered when a payment intent is successfully completed
- `payment_method.attached`: Triggered when a payment method is attached to a customer

**Security Features:**
- Request body size limit: 64KB
- Webhook signature verification
- JSON payload validation

**Response Codes:**
- `200 OK`: Event processed successfully
- `400 Bad Request`: Invalid webhook signature or malformed JSON
- `503 Service Unavailable`: Error reading request body

## Development

### Project Structure

```
nucleus/
├── main.go          # Main application entry point
├── go.mod           # Go module dependencies
├── go.sum           # Go module checksums
├── README.md        # This file
├── clerk/           # Clerk integration (future)
└── stripe/          # Stripe utilities (future)
```

### Adding New Webhook Events

To handle additional Stripe webhook events:

1. Add a new case in the switch statement in `handleWebhook()`:
```go
case "your.event.type":
    var eventData YourEventStruct
    err := json.Unmarshal(event.Data.Raw, &eventData)
    if err != nil {
        // Handle error
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
- Use HTTPS in production
- Keep your Stripe keys secure and never commit them to version control
- Implement proper error handling and logging
- Consider rate limiting for webhook endpoints

## Dependencies

- `github.com/joho/godotenv`: Environment variable management
- `github.com/stripe/stripe-go/v82`: Stripe Go SDK

## License

[Add your license information here]

## Contributing

[Add contribution guidelines here]