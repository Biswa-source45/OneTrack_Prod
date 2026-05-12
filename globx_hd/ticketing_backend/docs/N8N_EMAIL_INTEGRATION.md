# n8n Email-to-Ticket Integration Guide

This document explains how to set up n8n to automatically create tickets from incoming emails.

---

## 🚀 RECOMMENDED: AI Agent with Tools

The most powerful approach uses n8n's **AI Agent** node with **Tools** (function calling). The AI autonomously queries your database and sends **final resolved IDs** - no hints, no backend lookups needed!

### How It Works

```
┌──────────────────────────────────────────────────────────────────────────────┐
│                        n8n AI AGENT WORKFLOW                                 │
│                                                                              │
│  ┌─────────┐    ┌─────────────────────────────────────────────────────────┐ │
│  │  Email  │───▶│                  AI AGENT (GPT-4/Claude)                │ │
│  │ Trigger │    │                                                         │ │
│  └─────────┘    │  1. "Let me find the contact by sender email..."        │ │
│                 │     🔧 TOOL: search_contact_by_email                    │ │
│                 │     ✅ Found: contact_id=5, account_id=1                │ │
│                 │                                                         │ │
│                 │  2. "Now finding the product mentioned..."              │ │
│                 │     🔧 TOOL: search_product                             │ │
│                 │     ✅ Found: product_id=2                              │ │
│                 │                                                         │ │
│                 │  3. "Creating final ticket data..."                     │ │
│                 └─────────────────────────────────────────────────────────┘ │
│                                        │                                     │
│                                        ▼                                     │
│                 ┌─────────────────────────────────────────────────────────┐ │
│                 │  OUTPUT (Final resolved IDs - ready to use!)            │ │
│                 │  {                                                      │ │
│                 │    "contact_id": 5,                                     │ │
│                 │    "account_id": 1,                                     │ │
│                 │    "product_id": 2,                                     │ │
│                 │    "subject": "Dashboard error",                        │ │
│                 │    "details": "User reports error 500...",              │ │
│                 │    "priority": "High"                                   │ │
│                 │  }                                                      │ │
│                 └─────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────────────────────┘
                                        │
                                        ▼
┌──────────────────────────────────────────────────────────────────────────────┐
│  BACKEND: Simple POST /n8n/ticket - NO LOOKUPS NEEDED!                       │
└──────────────────────────────────────────────────────────────────────────────┘
```

### Available AI Tools

| Tool | Endpoint | Description |
|------|----------|-------------|
| `search_contact_by_email` | `/n8n/tools/contact/by-email?email=` | Find contact by exact email (highest confidence) |
| `search_contact_by_phone` | `/n8n/tools/contact/by-phone?phone=` | Find contact by phone number |
| `search_contact_by_name` | `/n8n/tools/contact/by-name?name=&account_id=` | Find contact by name (optionally in specific account) |
| `search_account_by_name` | `/n8n/tools/account/by-name?name=` | Find account by organization name |
| `search_account_by_domain` | `/n8n/tools/account/by-domain?domain=` | Find account by email domain |
| `get_account_contacts` | `/n8n/tools/account/contacts?account_id=` | Get all contacts in an account |
| `search_product` | `/n8n/tools/product/search?query=` | Find product by keyword |
| `list_products` | `/n8n/tools/products` | List all available products |

### Tool Response Examples

**search_contact_by_email - Found:**
```json
{
  "found": true,
  "contact_id": 5,
  "first_name": "John",
  "last_name": "Doe",
  "email": "john.doe@acmecorp.com",
  "mobile": "+919876543210",
  "account_id": 1,
  "account_name": "Acme Corporation",
  "customer_code": "001"
}
```

**search_contact_by_email - Not Found:**
```json
{
  "found": false,
  "message": "No contact found with this email"
}
```

### n8n AI Agent Setup (Step-by-Step)

#### Step 1: Create New Workflow
1. Open n8n (http://localhost:5678)
2. Click "New Workflow"
3. Name it "Email to Ticket (AI Agent)"

#### Step 2: Add Gmail Trigger Node
1. Add node → Search "Gmail Trigger"
2. Configure OAuth2 credentials
3. Set poll interval or use push

#### Step 3: Add AI Agent Node
1. Add node → Search "AI Agent"
2. Connect Gmail Trigger → AI Agent
3. Configure:
   - **Chat Model**: OpenAI GPT-4 or Anthropic Claude
   - **System Prompt**: (see below)

#### Step 4: Configure Tools in AI Agent
In the AI Agent node, add **HTTP Request Tool** for each lookup:

**Tool 1: search_contact_by_email**
```
Name: search_contact_by_email
Description: Search for a contact by their email address. Returns contact_id, account_id if found.
Method: GET
URL: http://localhost:8080/n8n/tools/contact/by-email
Query Parameters: email (the email to search)
Headers: X-N8N-API-Key: your-api-key
```

**Tool 2: search_contact_by_phone**
```
Name: search_contact_by_phone
Description: Search for a contact by phone number. Handles various formats.
Method: GET
URL: http://localhost:8080/n8n/tools/contact/by-phone
Query Parameters: phone (the phone number)
Headers: X-N8N-API-Key: your-api-key
```

**Tool 3: search_account_by_name**
```
Name: search_account_by_name
Description: Search for an account/organization by name.
Method: GET
URL: http://localhost:8080/n8n/tools/account/by-name
Query Parameters: name (organization name)
Headers: X-N8N-API-Key: your-api-key
```

**Tool 4: search_account_by_domain**
```
Name: search_account_by_domain
Description: Find account by email domain (e.g., acmecorp.com).
Method: GET  
URL: http://localhost:8080/n8n/tools/account/by-domain
Query Parameters: domain (email domain without @)
Headers: X-N8N-API-Key: your-api-key
```

**Tool 5: get_account_contacts**
```
Name: get_account_contacts
Description: Get all contacts belonging to an account.
Method: GET
URL: http://localhost:8080/n8n/tools/account/contacts
Query Parameters: account_id (the account ID)
Headers: X-N8N-API-Key: your-api-key
```

**Tool 6: search_product**
```
Name: search_product
Description: Search for a product by name or keyword.
Method: GET
URL: http://localhost:8080/n8n/tools/product/search
Query Parameters: query (product name or keyword)
Headers: X-N8N-API-Key: your-api-key
```

#### Step 5: AI Agent System Prompt

```
You are an intelligent ticket creation assistant. Your job is to analyze incoming support emails and create tickets with the correct customer information.

WORKFLOW:
1. First, extract the sender's email address and search for them using search_contact_by_email
2. If not found, extract any phone numbers from the email and try search_contact_by_phone
3. If still not found, extract organization names and try search_account_by_name, then use get_account_contacts to find a contact
4. If still not found, try search_account_by_domain using the sender's email domain
5. For products, extract any product/service mentions and use search_product

IMPORTANT RULES:
- ALWAYS try search_contact_by_email first with the sender's email
- If a tool returns found=false, try the next lookup method
- If multiple results are returned, pick the most relevant one
- For products, if nothing matches, use the first/default product
- Extract priority from keywords like "urgent", "critical", "asap" = High; "when possible", "low priority" = Low; otherwise Medium

OUTPUT FORMAT (JSON only, no markdown):
{
  "contact_id": <number>,
  "account_id": <number or null>,
  "product_id": <number>,
  "subject": "<brief subject line>",
  "details": "<full issue description from email>",
  "priority": "High" | "Medium" | "Low",
  "resolution_notes": "<brief note on how contact was identified>"
}

If you absolutely cannot find a contact after trying all methods, return:
{
  "error": true,
  "message": "Could not identify contact",
  "extracted_data": { ... all extracted info ... }
}
```

#### Step 6: Add HTTP Request Node for Ticket Creation
1. Add node → "HTTP Request"
2. Connect AI Agent → HTTP Request
3. Configure:
   - Method: POST
   - URL: http://localhost:8080/n8n/ticket
   - Headers: X-N8N-API-Key: your-api-key
   - Body: `{{ $json }}` (pass AI Agent output directly)

#### Step 7: Activate Workflow
1. Save workflow
2. Click "Activate" toggle
3. Send a test email to verify

---

## Architecture Overview

```
┌─────────────────┐     ┌────────────────────────────────────────────────┐     ┌──────────────────┐
│   Email Inbox   │────▶│                     n8n                        │────▶│  Backend API     │
│  (Gmail/IMAP)   │     │  ┌─────────┐  ┌─────────┐  ┌───────────────┐   │     │  /n8n/ticket     │
└─────────────────┘     │  │ Email   │─▶│ AI/Code │─▶│ HTTP Request  │   │     └──────────────────┘
                        │  │ Trigger │  │  Node   │  │    Node       │   │
                        │  └─────────┘  └─────────┘  └───────────────┘   │
                        └────────────────────────────────────────────────┘
```

## Backend API Endpoints

### Authentication
All `/n8n/*` endpoints use API key authentication instead of JWT.

**Header:** `X-N8N-API-Key: your-api-key`  
**Or Query Parameter:** `?api_key=your-api-key`

### Endpoints

| Endpoint | Method | Description |
|----------|--------|-------------|
| `/n8n/health` | GET | Health check |
| `/n8n/ticket` | POST | Create ticket |
| `/n8n/lookup/accounts` | GET | Search accounts |
| `/n8n/lookup/contacts` | GET | Search contacts |
| `/n8n/lookup/products` | GET | Search/list products |

---

## API Reference

### 1. Health Check
```http
GET /n8n/health
X-N8N-API-Key: your-api-key
```

**Response:**
```json
{
  "status": "ok",
  "service": "n8n-webhook",
  "time": "2025-12-11T11:00:00Z"
}
```

### 2. Create Ticket
```http
POST /n8n/ticket
X-N8N-API-Key: your-api-key
Content-Type: application/json
```

**Request Body:**
```json
{
  // OPTION 1: Using IDs (if you did lookups first)
  "account_id": 1,
  "contact_id": 5,
  "product_id": 2,
  
  // OPTION 2: Using names (system will lookup)
  "account_name": "Acme Corporation",
  "contact_name": "John Doe",
  "contact_email": "john@acme.com",
  "product_name": "Premium Support",
  
  // Required fields
  "subject": "Cannot access dashboard",
  "details": "User reports error 500 when accessing the main dashboard...",
  
  // Optional fields
  "priority": "High",           // High, Medium, Low (default: Medium)
  "sender_email": "john@acme.com",
  "sender_user_id": 3,          // For activity logging
  
  // Attachments (base64 encoded)
  "attachments": [
    {
      "filename": "screenshot.png",
      "content_type": "image/png",
      "data": "iVBORw0KGgoAAAANSUhEUgAA..."
    }
  ],
  
  // Raw email data (for audit)
  "raw_email_subject": "Original email subject",
  "raw_email_from": "John Doe <john@acme.com>",
  "raw_email_date": "2025-12-11T10:30:00Z"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "ticket_id": "001-111225-0001",
  "ticket": {
    "id": 42,
    "ticket_id": "001-111225-0001",
    "subject": "Cannot access dashboard",
    "status": "OPEN",
    "priority": "High",
    "account": "Acme Corporation",
    "contact": "John Doe",
    "product": "Premium Support",
    "created_at": "2025-12-11T11:00:00Z"
  }
}
```

**Error Response (400):**
```json
{
  "error": "Account not found",
  "details": "Could not find account: Acme Corp",
  "hint": "Use /n8n/lookup/accounts endpoint to find the correct account"
}
```

### 3. Lookup Accounts
```http
GET /n8n/lookup/accounts?search=acme
X-N8N-API-Key: your-api-key
```

**Response:**
```json
{
  "count": 2,
  "accounts": [
    {
      "id": 1,
      "account_name": "Acme Corporation",
      "customer_code": "001",
      "account_owner": "Sales Team",
      "score": 100
    },
    {
      "id": 5,
      "account_name": "Acme Industries",
      "customer_code": "005",
      "account_owner": "Enterprise Team",
      "score": 80
    }
  ]
}
```

### 4. Lookup Contacts
```http
GET /n8n/lookup/contacts?search=john&account_id=1
X-N8N-API-Key: your-api-key
```

**Response:**
```json
{
  "count": 1,
  "contacts": [
    {
      "id": 5,
      "first_name": "John",
      "last_name": "Doe",
      "email": "john@acme.com",
      "mobile": "+1234567890",
      "account_id": 1,
      "account_name": "Acme Corporation",
      "score": 100
    }
  ]
}
```

### 5. Lookup Products
```http
GET /n8n/lookup/products?search=premium
X-N8N-API-Key: your-api-key
```

**Response:**
```json
{
  "count": 2,
  "products": [
    {
      "id": 2,
      "product_name": "Premium Support",
      "product_description": "24/7 priority support"
    }
  ]
}
```

---

## n8n Workflow Setup

### Prerequisites
1. n8n installed and running (https://n8n.io)
2. Backend server running with n8n endpoints enabled
3. API key configured in `.env.email`

### Step 1: Create Gmail Trigger Node

1. Add **Gmail Trigger** node (or IMAP Email Trigger)
2. Configure credentials (OAuth2 for Gmail)
3. Set **Poll Times** to check every minute or use webhook
4. Filter by labels or sender if needed

**Output fields you'll use:**
- `from.value[0].address` - Sender email
- `from.value[0].name` - Sender name
- `subject` - Email subject
- `text` - Plain text body
- `html` - HTML body (optional)
- `attachments` - Array of attachments

### Step 2: Add AI/Code Node for Data Extraction

**Option A: Using OpenAI/Claude Node**

Add an **OpenAI** or **AI Agent** node with this prompt:

```
Extract structured data from this support email. Return JSON with these fields:
- account_name: Company/organization name mentioned
- contact_name: Person's name (sender or mentioned)
- contact_email: Email address
- subject: Brief summary (max 100 chars)
- details: Full issue description
- priority: "High" if urgent/critical, "Low" if minor, otherwise "Medium"

Email Subject: {{ $json.subject }}
Email From: {{ $json.from.value[0].name }} <{{ $json.from.value[0].address }}>
Email Body:
{{ $json.text }}

Return ONLY valid JSON, no markdown or explanation.
```

**Option B: Using Code Node (JavaScript)**

```javascript
const emailFrom = $input.item.json.from?.value?.[0] || {};
const subject = $input.item.json.subject || '';
const body = $input.item.json.text || '';

// Extract company name from email domain
const emailDomain = emailFrom.address?.split('@')[1]?.split('.')[0] || '';
const companyName = emailDomain.charAt(0).toUpperCase() + emailDomain.slice(1);

// Determine priority from keywords
const urgentKeywords = ['urgent', 'critical', 'emergency', 'asap', 'immediately'];
const lowKeywords = ['question', 'inquiry', 'when you have time'];

let priority = 'Medium';
const lowerBody = body.toLowerCase();
const lowerSubject = subject.toLowerCase();

if (urgentKeywords.some(kw => lowerBody.includes(kw) || lowerSubject.includes(kw))) {
  priority = 'High';
} else if (lowKeywords.some(kw => lowerBody.includes(kw) || lowerSubject.includes(kw))) {
  priority = 'Low';
}

return {
  account_name: companyName || 'Unknown',
  contact_name: emailFrom.name || 'Unknown',
  contact_email: emailFrom.address || '',
  subject: subject.substring(0, 255),
  details: body,
  priority: priority,
  raw_email_from: `${emailFrom.name} <${emailFrom.address}>`,
  raw_email_subject: subject
};
```

### Step 3: (Optional) Add Lookup Nodes

If you want to verify accounts/contacts exist before creating:

**HTTP Request Node - Lookup Account:**
```
Method: GET
URL: http://localhost:8080/n8n/lookup/accounts
Query Parameters:
  - search: {{ $json.account_name }}
Headers:
  - X-N8N-API-Key: your-api-key
```

Then use an **IF** node to check if `{{ $json.count }} > 0`

### Step 4: Create Ticket via HTTP Request

**HTTP Request Node:**
```
Method: POST
URL: http://localhost:8080/n8n/ticket
Headers:
  - X-N8N-API-Key: your-api-key
  - Content-Type: application/json
Body (JSON):
{
  "account_name": "{{ $json.account_name }}",
  "contact_email": "{{ $json.contact_email }}",
  "subject": "{{ $json.subject }}",
  "details": "{{ $json.details }}",
  "priority": "{{ $json.priority }}",
  "raw_email_from": "{{ $json.raw_email_from }}",
  "raw_email_subject": "{{ $json.raw_email_subject }}"
}
```

### Step 5: Handle Attachments (Optional)

If emails have attachments, add a **Code Node** before the HTTP request:

```javascript
const attachments = $input.item.json.attachments || [];
const processedAttachments = [];

for (const att of attachments) {
  if (att.content) {
    processedAttachments.push({
      filename: att.filename || 'attachment',
      content_type: att.contentType || 'application/octet-stream',
      data: att.content.toString('base64')
    });
  }
}

return {
  ...$input.item.json,
  attachments: processedAttachments
};
```

---

## Complete n8n Workflow JSON

Import this workflow directly into n8n:

```json
{
  "name": "Email to Ticket",
  "nodes": [
    {
      "parameters": {
        "pollTimes": {
          "item": [{"mode": "everyMinute"}]
        },
        "simple": false
      },
      "name": "Gmail Trigger",
      "type": "n8n-nodes-base.gmailTrigger",
      "position": [250, 300]
    },
    {
      "parameters": {
        "jsCode": "const emailFrom = $input.item.json.from?.value?.[0] || {};\nconst subject = $input.item.json.subject || '';\nconst body = $input.item.json.text || '';\n\nconst emailDomain = emailFrom.address?.split('@')[1]?.split('.')[0] || '';\nconst companyName = emailDomain.charAt(0).toUpperCase() + emailDomain.slice(1);\n\nlet priority = 'Medium';\nconst lowerBody = body.toLowerCase();\nif (lowerBody.includes('urgent') || lowerBody.includes('critical')) priority = 'High';\nif (lowerBody.includes('question') || lowerBody.includes('inquiry')) priority = 'Low';\n\nreturn {\n  account_name: companyName,\n  contact_email: emailFrom.address,\n  subject: subject.substring(0, 255),\n  details: body,\n  priority: priority,\n  raw_email_from: `${emailFrom.name} <${emailFrom.address}>`,\n  raw_email_subject: subject\n};"
      },
      "name": "Extract Data",
      "type": "n8n-nodes-base.code",
      "position": [470, 300]
    },
    {
      "parameters": {
        "method": "POST",
        "url": "http://localhost:8080/n8n/ticket",
        "sendHeaders": true,
        "headerParameters": {
          "parameters": [
            {"name": "X-N8N-API-Key", "value": "={{ $env.N8N_API_KEY }}"}
          ]
        },
        "sendBody": true,
        "bodyParameters": {
          "parameters": [
            {"name": "account_name", "value": "={{ $json.account_name }}"},
            {"name": "contact_email", "value": "={{ $json.contact_email }}"},
            {"name": "subject", "value": "={{ $json.subject }}"},
            {"name": "details", "value": "={{ $json.details }}"},
            {"name": "priority", "value": "={{ $json.priority }}"},
            {"name": "raw_email_from", "value": "={{ $json.raw_email_from }}"}
          ]
        }
      },
      "name": "Create Ticket",
      "type": "n8n-nodes-base.httpRequest",
      "position": [690, 300]
    }
  ],
  "connections": {
    "Gmail Trigger": {"main": [[{"node": "Extract Data", "type": "main", "index": 0}]]},
    "Extract Data": {"main": [[{"node": "Create Ticket", "type": "main", "index": 0}]]}
  }
}
```

---

## Configuration Checklist

### Backend (.env.email)
- [ ] Set `N8N_API_KEY` to a secure random value
- [ ] Configure `EMAIL_SMTP_SERVER` and `EMAIL_SMTP_PORT` for notifications
- [ ] Set `EMAIL_NOTIFICATION_ADDRESS` for ticket creation alerts

### n8n
- [ ] Configure Gmail/Email credentials
- [ ] Set `N8N_API_KEY` environment variable (or hardcode in headers)
- [ ] Test workflow with a sample email
- [ ] Enable workflow for production

### Testing
```bash
# Test health endpoint
curl -H "X-N8N-API-Key: your-api-key" http://localhost:8080/n8n/health

# Test account lookup
curl -H "X-N8N-API-Key: your-api-key" "http://localhost:8080/n8n/lookup/accounts?search=acme"

# Test ticket creation
curl -X POST http://localhost:8080/n8n/ticket \
  -H "X-N8N-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "account_name": "Test Account",
    "contact_email": "test@example.com",
    "subject": "Test ticket from n8n",
    "details": "This is a test ticket created via the n8n webhook endpoint."
  }'
```

---

## Troubleshooting

### Common Issues

| Issue | Solution |
|-------|----------|
| 401 Unauthorized | Check `X-N8N-API-Key` header matches `N8N_API_KEY` in .env.email |
| Account not found | Use lookup endpoint first, or check account name spelling |
| Contact not found | System will use first contact from account as fallback |
| No products | At least one product must exist in database |

### Logs
Check backend logs for detailed error messages:
```bash
# All n8n related logs are prefixed with [n8n]
grep "\[n8n\]" app.log
```

---

## Disabling Legacy IMAP Processor

Once n8n is working, you can disable the legacy IMAP email processor:

1. Remove or comment out `go startEmailProcessor()` in `cmd/main.go`
2. Or simply don't configure `EMAIL_USERNAME`/`EMAIL_PASSWORD` in environment

The n8n integration is completely independent of the legacy IMAP processor.
