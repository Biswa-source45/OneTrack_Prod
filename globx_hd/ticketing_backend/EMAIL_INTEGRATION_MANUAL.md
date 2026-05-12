# Email Integration Manual

## Overview

This document describes how to set up and use the email integration feature for the Globx Ticketing System. This feature allows users to create tickets by sending emails to a designated email address, making ticket creation possible from any email client without needing to log into the system.

## Table of Contents

1. [Setup Instructions](#setup-instructions)
   - [System Requirements](#system-requirements)
   - [Configuration](#configuration)
   - [Running the Email Processor](#running-the-email-processor)
   - [Setting up a Scheduled Task](#setting-up-a-scheduled-task)

2. [Email Format Guidelines](#email-format-guidelines)
   - [Required Email Format](#required-email-format)
   - [Email Subject](#email-subject)
   - [Email Body](#email-body)
   - [Attachments](#attachments)
   - [Example Email](#example-email)

3. [Troubleshooting](#troubleshooting)
   - [Common Issues](#common-issues)
   - [Log Files](#log-files)

## Setup Instructions

### System Requirements

- Gmail account for receiving ticket creation emails
- App Password for Gmail (regular passwords won't work due to Google security measures)
- Go 1.16 or later
- Access to the ticketing system database

### Configuration

1. Copy the `.env.email` template to `.env.email` in the project root directory.

2. Edit the `.env.email` file with your specific settings:

```
# Gmail IMAP Server Settings
EMAIL_IMAP_SERVER=imap.gmail.com
EMAIL_IMAP_PORT=993

# Gmail Credentials
EMAIL_USERNAME=your-mailbox@gmail.com
EMAIL_PASSWORD=your-app-password  # Use App Password

# Allowed sender email addresses (comma-separated)
EMAIL_ALLOWED_SENDERS=sender1@gmail.com,sender2@gmail.com

# Email address for notifications
EMAIL_NOTIFICATION_ADDRESS=notifications@yourdomain.com

# Path for file attachments
UPLOAD_DIR=./uploads

# Database connection (optional)
# DATABASE_URL=host=localhost user=postgres password=postgres dbname=globx_hd port=5432
```

3. How to generate an App Password for Gmail:
   - Go to your Google Account
   - Navigate to Security
   - Under "Signing in to Google," select "App passwords"
   - Generate a new app password for "Mail" and "Windows Computer"
   - Use the generated 16-character password in your configuration

4. Make sure the allowed sender email addresses are correctly specified. Only emails from these addresses will be processed.

### Running the Email Processor

#### Option 1: Manual execution

Run the email processor script:

```
cd path/to/ticketing_backend
./email_processor.ps1
```

This will start the email processor which will check for new emails every 5 minutes.

#### Option 2: Running once

To process emails just once:

```
cd path/to/ticketing_backend
go run ./cmd/email_processor/main.go --once
```

### Setting up a Scheduled Task

To run the email processor automatically:

#### Windows Task Scheduler

1. Open Task Scheduler
2. Create a new Basic Task
3. Name it "Globx Email Processor"
4. Set the trigger to "Daily" and configure it to repeat every 5 minutes
5. Set the action to "Start a program"
6. Program/script: `powershell.exe`
7. Arguments: `-ExecutionPolicy Bypass -File "C:\path\to\ticketing_backend\email_processor.ps1"`
8. Set "Start in" to your project directory

## Email Format Guidelines

### Required Email Format

For an email to be correctly processed, it must follow these format requirements:

1. Must be sent from one of the authorized email addresses configured in `.env.email`
2. Must be sent to the email address configured as `EMAIL_USERNAME`
3. Must contain required fields in the email body

### Email Subject

The email subject will be used as the ticket subject.

### Email Body

The email body must contain the following sections:

```
Account: [name of an existing account in the system]
Contact: [optional: name or email of contact in that account]
Details: [Detailed description of the issue/ticket]
```

**Notes:**
- The `Account` field must match an existing account name in the system
- Account name matching is case-insensitive and will find partial matches
- The `Contact` field is optional - if provided, it will try to find a matching contact in that account
- Contact lookup is smart and supports email, full name, or partial name matching
- If no contact is specified, it will use the first contact found for the account
- The `Details` can be on the same line (e.g., "Details: description") or start on the next line
- The `Details` section can span multiple lines

### Attachments

Any files attached to the email will be automatically attached to the created ticket. The system supports all file types that the main ticketing system supports.

### Example Email

**Subject:** Unable to access account dashboard

**Body:**
```
Account: Acme Corporation
Contact: John Smith
Details: I am unable to access the customer dashboard since yesterday. 
When I try to log in, it accepts my credentials but then shows a blank screen.

I have tried clearing my browser cache and using different browsers with no success.
Please help resolve this issue as it's affecting our operations.

Thank you,
John
```

## Troubleshooting

### Common Issues

#### Emails Not Being Processed

- Check if the email processor is running
- Verify that the sender email is in the allowed list
- Make sure the Gmail credentials are correct
- Check if the App Password is valid and correctly entered
- Ensure that IMAP is enabled in the Gmail account settings

#### Incorrect Ticket Data

- Make sure the email format strictly follows the guidelines
- Check if the contact exists in the system
- Check if the product name is recognized

#### Attachment Issues

- Verify that the upload directory exists and has proper permissions
- Check attachment size (Gmail has a 25MB limit)

### Log Files

The email processor logs its activities to the standard output. If running as a scheduled task, redirect the output to a log file for debugging:

```powershell
# In email_processor.ps1, add:
./bin/email_processor.exe > email_processor.log 2>&1
```

Check this log file for errors or warnings if issues occur.

---

For additional support or questions, please contact the system administrator. 