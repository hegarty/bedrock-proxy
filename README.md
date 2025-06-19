# bedrock-proxy

A Go-based proxy service that provides access to AWS Bedrock's Claude AI models through both HTTP API and Stdin/Stdout, initially set up to integration with the Zed editor - but now they directly integration with Bedrock.

## Overview

This proxy service allows you to interact with Claude AI models hosted on AWS Bedrock in two modes:
- **HTTP Mode**: RESTful API server for general use
- **Zed Mode**: Direct integration with the Zed code editor using Stdin/Stdout

## Prerequisites

- Go 1.23.10 or later
- AWS with Claude enabled via Bedrock

## Setup

### 1. Clone and Build

```bash
git clone git@github.com:hegarty/bedrock-proxy.git
cd bedrock-proxy
go mod tidy
go build -o claude-proxy
```

### 2. AWS Authentication

This service requires AWS authentication with access to AWS Bedrock. Set up authentication using the following steps:

```bash
# Set the AWS profile
export AWS_PROFILE=[your-profile-here]

# Login to AWS SSO
aws sso login

# Verify access to Bedrock (optional)
aws bedrock list-foundation-models --region us-east-1
```

## Usage

The application automatically detects the execution mode:
- If run with stdin input (piped data), it operates in **Zed mode**
- If run normally, it starts an **HTTP server**

### HTTP Mode

Start the HTTP server:

```bash
./claude-proxy
```

The server will start on port 8080. You should see:
```
Starting HTTP server on :8080
```

#### Testing the HTTP API

Test the service with curl:

```bash
curl -X POST http://localhost:8080/invoke \
  -H "Content-Type: application/json" \
  -d '{
    "input": "what is your name?"
  }'
```

Expected response:
```json
{
  "response": "I'm Claude, an AI assistant created by Anthropic..."
}
```

#### HTTP API Reference

**Endpoint:** `POST /invoke`

**Request Body:**
```json
{
  "input": "Your question or prompt here"
}
```

**Response:**
```json
{
  "response": "Claude's response text"
}
```

**Error Responses:**
- `400 Bad Request`: Invalid JSON in request body
- `500 Internal Server Error`: AWS Bedrock call failed or response parsing error

### Model Configuration

The service is currently configured to use:
- **Model:** `anthropic.claude-3-sonnet-20240229-v1:0`
- **Max Tokens:** 1000
- **API Version:** `bedrock-2023-05-31`

To modify these settings, edit the `bedrock/model.go` file.

### AWS Region

The service uses the default AWS region from your AWS configuration. Ensure your region has access to AWS Bedrock and the Claude models.

## Project Structure

```
bedrock-proxy/
├── main.go              # Application entry point and mode detection
├── bedrock/
│   └── model.go         # AWS Bedrock integration and Claude API calls
├── httpmode/
│   └── http.go          # HTTP server implementation
├── zedmode/
│   └── zedmode.go       # Zed editor integration
├── go.mod               # Go module dependencies
└── README.md           # This file
```

## Dependencies

- `github.com/aws/aws-sdk-go-v2` - AWS SDK for Go v2
- `github.com/aws/aws-sdk-go-v2/config` - AWS configuration loading
- `github.com/aws/aws-sdk-go-v2/service/bedrockruntime` - AWS Bedrock Runtime service

## Troubleshooting

### Common Issues

1. **Authentication Errors**
   - Ensure `AWS_PROFILE=[profile-name]` is set
   - Run `aws sso login` to refresh credentials
   - Verify your AWS profile has Bedrock permissions

2. **Model Access Errors**
   - Ensure your AWS account has access to Claude models in Bedrock
   - Check that the region supports the specified model

3. **JSON Parsing Errors**
   - Verify the request body is valid JSON
   - Ensure the `input` field is present in the request
