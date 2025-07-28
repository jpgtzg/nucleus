# Nucleus

Nucleus is the central control service of the platform that synchronizes billing data from Stripe with user identity data from Clerk, manages organization metadata, and provides subscription-based access control. Acting as the source of truth, Nucleus ensures all user entitlements, usage limits, and service-level permissions are consistently maintained and distributed across the system.

## Features

- **Stripe Webhook Handler**: Processes Stripe webhook events for subscription management
- **Clerk Webhook Handler**: Processes Clerk webhook events for organization lifecycle management
- **Organization Management**: Automatic creation and deletion of Stripe customers when organizations are created/deleted in Clerk
- **Subscription-Based Access Control**: Comprehensive subscription tracking and management
- **JWT Authentication**: Secure middleware for verifying Clerk JWT tokens
- **Dynamic IP Validation**: Fetches current Stripe webhook IPs dynamically for enhanced security
- **MongoDB Integration**: Persistent storage for organization-customer mappings
- **Asynchronous Processing**: Webhook events are processed asynchronously for better performance
- **Docker Support**: Containerized deployment with optimized build process

## Architecture

Nucleus acts as a bridge between three key services:

1. **Clerk**: User identity and organization management
2. **Stripe**: Payment processing and subscription management  
3. **MongoDB**: Persistent storage for organization-customer mappings

### Data Flow

```
Clerk Organization Created → Nucleus → Create Stripe Customer → Store in MongoDB
Stripe Subscription Event → Nucleus → Update Organization Metadata in Clerk
User Request → Nucleus (JWT Auth) → Return User's Active Subscriptions
```

## Prerequisites

- Go 1.24.4 or higher (for local development)
- Docker (for containerized deployment)
- Stripe account with webhook endpoint configured
- Clerk account for user identity management
- MongoDB database for data persistence

## Installation

### Local Development

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
CLERK_SECRET_KEY=your_clerk_secret_key
CLERK_WEBHOOK_SECRET=whsec_your_clerk_webhook_secret
MONGO_URI=mongodb://localhost:27017
MONGO_DATABASE=nucleus
MONGO_COLLECTION_SYNC=organizations
PORT=8080
```

### Docker Deployment

1. Clone the repository:
```bash
git clone https://github.com/SAMLA-io/nucleus.git
cd nucleus
```

2. Create a `.env` file with your configuration (same as above)

3. Build the Docker image:
```bash
docker build -t nucleus .
```

4. Run the container:
```bash
docker run -p 8080:8080 --env-file .env nucleus
```

The server will start listening on port 8080 and be accessible at `http://localhost:8080`.

## Configuration

### Stripe Webhook Setup

1. Go to your Stripe Dashboard
2. Navigate to Developers > Webhooks
3. Create a new webhook endpoint with your server URL (e.g., `https://yourdomain.com/stripe/webhook`)
4. Select the following events to listen for:
   - `customer.subscription.created`
   - `customer.subscription.updated`
   - `customer.subscription.deleted`
5. Copy the webhook signing secret and add it to your `.env` file

### Clerk Webhook Setup

1. Go to your Clerk Dashboard
2. Navigate to Webhooks
3. Create a new webhook endpoint with your server URL (e.g., `https://yourdomain.com/clerk/webhook`)
4. Select the following events to listen for:
   - `organization.created`
   - `organization.updated`
   - `organization.deleted`
5. Copy the webhook signing secret and add it to your `.env` file

### MongoDB Setup

1. Set up a MongoDB database (local or cloud-based like MongoDB Atlas)
2. Create a database with your preferred name
3. Create a collection with your preferred name
4. The collection will automatically store documents with the following structure:
```json
{
  "_id": "ObjectId",
  "clerk_organization_id": "string",
  "stripe_customer_id": "string"
}
```
5. Add your MongoDB connection URI to your `.env` file

## Usage

### Running the Server

#### Local Development
```bash
go run main.go
```

#### Docker
```bash
# Build the image
docker build -t nucleus .

# Run the container
docker run -p 8080:8080 --env-file .env nucleus

# Run in detached mode
docker run -d -p 8080:8080 --env-file .env --name nucleus-app nucleus

# View logs
docker logs nucleus-app

# Stop the container
docker stop nucleus-app
```

The server will start listening on the configured port (default: 8080).

### Environment Variables

The application supports both local `.env` files and Docker environment variable injection:

- **Local Development**: Uses `godotenv` to load `.env` file
- **Docker**: Uses `--env-file .env` to pass environment variables to the container
- **Fallback**: If `.env` file is not found, the application uses system environment variables

## API Endpoints

### Stripe Webhook

#### POST `/stripe/webhook`

Handles incoming Stripe webhook events.

**Supported Events:**
- `customer.subscription.created`: Creates new subscription in organization metadata
- `customer.subscription.updated`: Updates existing subscription information
- `customer.subscription.deleted`: Removes subscription from organization metadata

**Security Features:**
- Request body size limit: 64KB
- Webhook signature verification using Stripe SDK
- Dynamic IP address validation against Stripe's current webhook IP list
- JSON payload validation
- Asynchronous event processing

**Response Codes:**
- `200 OK`: Event processed successfully
- `400 Bad Request`: Invalid webhook signature or malformed JSON
- `403 Forbidden`: Request from non-Stripe IP address
- `503 Service Unavailable`: Error reading request body

### Clerk Webhook

#### POST `/clerk/webhook`

Handles incoming Clerk webhook events.

**Supported Events:**
- `organization.created`: Creates new Stripe customer and stores mapping in MongoDB
- `organization.updated`: Handles organization updates (currently no-op)
- `organization.deleted`: Deletes Stripe customer and removes mapping from MongoDB

**Security Features:**
- Webhook signature verification using Svix
- JSON payload validation
- Asynchronous event processing

**Response Codes:**
- `200 OK`: Event processed successfully
- `400 Bad Request`: Malformed JSON or invalid webhook
- `401 Unauthorized`: Invalid webhook signature
- `405 Method Not Allowed`: Non-POST requests

### User API

#### GET `/user/subscriptions`

Returns the active subscriptions for the authenticated user's organization.

**Authentication:**
- Requires valid Clerk JWT token in Authorization header
- Format: `Bearer <token>`

**Response:**
```json
[
  {
    "id": "sub_123",
    "status": "active",
    "current_period_end": 1753903901,
    "product_id": "prod_premium",
    "price_id": "price_123"
  }
]
```

**Response Codes:**
- `200 OK`: Subscriptions returned successfully
- `401 Unauthorized`: Invalid or missing JWT token
- `500 Internal Server Error`: Error retrieving user data

#### GET `/user/stripe-customer-id`

Returns the Stripe customer ID for the authenticated user's organization.

**Authentication:**
- Requires valid Clerk JWT token in Authorization header
- Format: `Bearer <token>`

**Response:**
```json
{
  "stripe_customer_id": "cus_123"
}
```

**Response Codes:**
- `200 OK`: Stripe customer ID returned successfully
- `401 Unauthorized`: Invalid or missing JWT token
- `500 Internal Server Error`: Error retrieving user data

## Organization Management

### Automatic Customer Creation

When a new organization is created in Clerk:

1. Nucleus receives the `organization.created` webhook
2. Creates a new Stripe customer with the organization name
3. Stores the mapping between Clerk organization ID and Stripe customer ID in MongoDB
4. This mapping enables subscription events to be properly routed to the correct organization

### Subscription Metadata Structure

Organization metadata in Clerk includes comprehensive subscription information:

```json
{
  "stripe": {
    "subscriptions": [
      {
        "id": "sub_123",
        "status": "active",
        "current_period_end": 1753903901,
        "product_id": "prod_premium",
        "price_id": "price_123"
      }
    ]
  }
}
```

### Access Control Functions

The service provides helper functions for checking organization access:

```go
// Get active subscriptions for an organization
subscriptions := clerk.GetActiveSubscriptionsByOrganizationID(organizationID)

// Get active subscriptions by customer ID
subscriptions := clerk.GetActiveSubscriptionsByCustomerID(customerID)

// Add subscription to organization metadata
clerk.AddSubscriptionToOrganizationMetadata(customerID, subscription)

// Update subscription in organization metadata
clerk.UpdateSubscriptionInOrganizationMetadata(customerID, subscription)

// Remove subscription from organization metadata
clerk.RemoveSubscriptionFromOrganizationMetadata(customerID, subscriptionID)
```

## Authentication

### JWT Verification

The service includes middleware for verifying Clerk JWT tokens:

```go
// Middleware that verifies JWT and extracts user ID
auth.VerifyingMiddleware(next http.Handler)

// Extract user ID from request context
userID, ok := auth.GetUserID(r)
```

### Organization Resolution

For authenticated requests, the service automatically resolves the user's organization:

```go
// Get user's organization ID
organizationID, err := clerk.GetUserOrganizationId(userID)

// Get organization metadata
metadata, err := clerk.GetOrganizationPublicMetadata(organizationID)

// Update organization metadata
err = clerk.UpdateOrganizationPublicMetadata(organizationID, metadata)
```

## Development

### Project Structure

```
nucleus/
├── main.go                    # Main application entry point
├── Dockerfile                 # Docker container configuration
├── go.mod                     # Go module dependencies
├── go.sum                     # Go module checksums
├── README.md                  # This file
├── api/
│   └── handlers.go            # User API handlers
├── auth/
│   └── auth.go                # JWT authentication middleware
├── clerk/
│   ├── handlers.go            # Clerk webhook handlers
│   ├── organizations.go       # Organization management
│   ├── subscription.go        # Subscription metadata management
│   └── webhook.go            # Clerk webhook processing
├── stripe/
│   ├── address.go             # Dynamic webhook IP validation
│   ├── handlers.go            # Stripe event handlers
│   └── webhook.go            # Stripe webhook processing
├── mongodb/
│   └── sync.go               # Database operations
└── types/
    ├── clerk/
    │   └── types.go          # Clerk webhook event types
    └── mongodb/
        └── organizations.go   # Database model types
```

### Docker Configuration

The project includes optimized Docker configuration:

- **Multi-stage build**: Uses `golang:1.24-alpine` for efficient builds
- **Security**: Excludes sensitive files from build context
- **Environment handling**: Graceful fallback from `.env` files to system environment variables
- **Port exposure**: Exposes port 8080 for webhook endpoints

### Webhook IP Validation

The service dynamically fetches the current list of Stripe webhook IPs from Stripe's official endpoint (`https://stripe.com/files/ips/ips_webhooks.json`) to ensure the most up-to-date security validation.

### Adding New Webhook Events

To handle additional webhook events:

#### For Stripe Events:
1. Add a new case in the switch statement in `stripe/webhook.go`:
```go
case "your.event.type":
    // Handle the event
```

2. Implement the corresponding handler function in `stripe/handlers.go`

#### For Clerk Events:
1. Add a new case in the switch statement in `clerk/webhook.go`:
```go
case "your.event.type":
    return HandleYourEvent(event)
```

2. Implement the corresponding handler function in `clerk/handlers.go`

3. Update your Clerk webhook configuration to listen for the new event

### Testing

For local development, you can use tools like:
- [Stripe CLI](https://stripe.com/docs/stripe-cli) for webhook forwarding
- [ngrok](https://ngrok.com/) for exposing your local server

Example with Stripe CLI:
```bash
stripe listen --forward-to localhost:8080/stripe/webhook
```

For Docker testing:
```bash
# Run with Stripe CLI forwarding
stripe listen --forward-to localhost:8080/stripe/webhook

# In another terminal, run the container
docker run -p 8080:8080 --env-file .env nucleus
```

## Security Considerations

- Always verify webhook signatures using the provided secrets
- Dynamic webhook IP validation ensures requests only come from current Stripe IPs
- Use HTTPS in production
- Keep your Stripe, Clerk, and MongoDB connection strings secure and never commit them to version control
- Implement proper error handling and logging
- Consider rate limiting for webhook endpoints
- JWT tokens are verified using Clerk's official SDK
- Environment variables are passed securely to containers without embedding in images

## Dependencies

- `github.com/clerk/clerk-sdk-go/v2`: Clerk user management SDK
- `github.com/joho/godotenv`: Environment variable management
- `github.com/stripe/stripe-go/v82`: Stripe Go SDK
- `go.mongodb.org/mongo-driver/v2`: MongoDB Go driver
- `github.com/svix/svix-webhooks`: Clerk webhook verification

## License

This project is licensed under the Apache License 2.0. See the [LICENSE](LICENSE) file for details.

## Contributors

- [@jpgtzg](https://github.com/jpgtzg)