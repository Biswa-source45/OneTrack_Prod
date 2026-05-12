# Email Integration for Globx Ticketing System

## Overview

This feature allows users to create tickets by sending emails to a designated email address. The system automatically processes incoming emails, extracts the ticket information, and creates tickets in the system.

## Features

- Automated ticket creation from emails
- Attachment handling (files attached to emails are added to the ticket)
- Contact and product matching using intelligent fuzzy matching
- Security controls with sender email verification
- Notification emails for successfully created tickets

## Quick Start

1. Install required dependencies:
   ```
   go get github.com/emersion/go-imap
   go get github.com/emersion/go-message
   go get github.com/joho/godotenv
   ```

2. Copy the `.env.email` template to `.env.email` and configure your settings

3. Run the email processor:
   ```
   ./email_processor.ps1
   ```

For detailed instructions, see the [Email Integration Manual](./EMAIL_INTEGRATION_MANUAL.md).

## Email Format

Emails must follow this format:

```
Subject: [Ticket Subject]

Body:
Account: [name of existing account]
Contact: [optional: name or email of contact in that account]
Details: [Detailed description of the issue/ticket]
```

The system will:
- Use the email subject as the ticket subject
- Find the account by name (case-insensitive, partial matches supported)
- Use the account's customer code for proper ticket ID generation
- If Contact is specified, find the matching contact from that account
- If Contact is not specified, use the first contact from that account
- Set default values for other fields

## Architecture

The email integration consists of:

1. **Email Fetcher**: Connects to Gmail via IMAP and fetches unread emails
2. **Email Parser**: Extracts ticket information from emails
3. **Fuzzy Matcher**: Maps email text to master data IDs
4. **Ticket Creator**: Creates tickets using the existing ticket system
5. **Attachment Processor**: Handles file attachments
6. **Scheduler**: Runs the processing at regular intervals

## Security

- Only processes emails from allowed sender addresses
- Only checks the designated inbox
- Emails are marked as read after processing
- Uses Gmail App Passwords for secure authentication

## Logs

The email processor logs its activities to the standard output. Check these logs for troubleshooting.

---

For more information, see the [Email Integration Manual](./EMAIL_INTEGRATION_MANUAL.md).
