# n8n + Gemini Email-to-Ticket Integration

Complete guide to set up automated email-to-ticket creation using n8n and Google Gemini AI.

---

## Architecture

```
┌──────────────┐    ┌─────────────────────────────────────────────────────────┐    ┌──────────────────┐
│    EMAIL     │    │                         n8n                             │    │   GO BACKEND     │
│    INBOX     │───▶│  ┌──────────┐   ┌──────────────┐   ┌────────────────┐  │───▶│                  │
│  (Gmail)     │    │  │  Gmail   │──▶│   Gemini     │──▶│  HTTP Request  │  │    │ /n8n/smart-ticket│
└──────────────┘    │  │  Trigger │   │  (Extract)   │   │  (Create)      │  │    │                  │
                    │  └──────────┘   └──────────────┘   └────────────────┘  │    │ Smart Resolver:  │
                    │                                                         │    │ • Email match    │
                    │  AI extracts:                                           │    │ • Phone match    │
                    │  • Phone numbers                                        │    │ • Name fuzzy     │
                    │  • Person names                                         │    │ • Domain match   │
                    │  • Org names                                            │    │ • Best candidate │
                    │  • Product hints                                        │    │                  │
                    │  • Priority                                             │    │ Creates ticket   │
                    └─────────────────────────────────────────────────────────┘    └──────────────────┘
```

---

## Prerequisites

1. **n8n** installed and running (http://localhost:5678)
2. **Backend** running with n8n endpoints (http://localhost:8080)
3. **Google Cloud** account with Gemini API enabled
4. **Gmail** account for email trigger

---

## Step 1: Configure Backend

### 1.1 Set API Key in `.env.email`

```bash
# Generate a secure key
openssl rand -hex 32

# Add to .env.email
N8N_API_KEY=your-generated-secure-key-here

# SMTP settings for notifications
EMAIL_SMTP_SERVER=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
EMAIL_NOTIFICATION_ADDRESS=notifications@yourcompany.com
```

### 1.2 Start Backend

```bash
cd ticketing_backend
go run cmd/main.go
```

### 1.3 Test Endpoint

```bash
# Test health
curl -H "X-N8N-API-Key: your-api-key" http://localhost:8080/n8n/health

# Expected response:
# {"service":"n8n-webhook","status":"ok","time":"2025-12-11T14:00:00+05:30"}
```

---

## Step 2: Get Gemini API Key

1. Go to [Google AI Studio](https://aistudio.google.com/app/apikey)
2. Click **"Create API Key"**
3. Copy the key (starts with `AIza...`)
4. Keep it safe - you'll need it in n8n

---

## Step 3: Create n8n Workflow

### 3.1 Open n8n

Open http://localhost:5678 in your browser.

### 3.2 Create New Workflow

1. Click **"New Workflow"**
2. Name it: `Email to Ticket (Gemini AI)`

### 3.3 Add Gmail Trigger Node

1. Click **"+"** to add a node
2. Search for **"Gmail Trigger"**
3. Configure:
   - **Credential**: Create new Gmail OAuth2 credential
   - **Poll Times**: Every 1 minute (or your preference)
   - **Simple**: OFF (to get full email data)
4. Test the trigger by sending a test email

### 3.4 Add Gemini Node

1. Add a new node after Gmail Trigger
2. Search for **"Google Gemini Chat Model"** or **"HTTP Request"**

**Option A: Using HTTP Request (Recommended)**

Configure the HTTP Request node:

```
Method: POST
URL: https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key=YOUR_API_KEY

Headers:
  Content-Type: application/json

Body (JSON):
{
  "contents": [{
    "parts": [{
      "text": "You are an email parser for a ticketing system. Extract structured data from this email.\n\nFrom: {{ $json.from.value[0].name }} <{{ $json.from.value[0].address }}>\nSubject: {{ $json.subject }}\nBody:\n{{ $json.text }}\n\nExtract and return ONLY a valid JSON object (no markdown, no explanation) with these fields:\n- phone_numbers: array of phone numbers found (any format)\n- person_names: array of person names found (from signature, body, etc.)\n- org_names: array of organization/company names found\n- product_hints: array of product or service names mentioned\n- priority_hints: array of urgency keywords (urgent, asap, critical, etc.)\n\nIf a field has no values, use an empty array [].\nReturn ONLY the JSON object, nothing else."
    }]
  }],
  "generationConfig": {
    "temperature": 0.1,
    "maxOutputTokens": 1024
  }
}
```

**Option B: Using Google Gemini Node (if available)**

If you have the Google Gemini node:
1. Add credentials with your API key
2. Set the prompt (same as above)
3. Configure output format as JSON

### 3.5 Add Code Node (Parse Gemini Response)

Add a **Code** node to parse Gemini's response and prepare data:

```javascript
// Get Gemini response
const geminiResponse = $input.first().json;

// Extract the text content from Gemini response
let extractedText = '';
if (geminiResponse.candidates && geminiResponse.candidates[0]) {
  extractedText = geminiResponse.candidates[0].content.parts[0].text;
}

// Parse JSON from Gemini (it might have markdown code blocks)
let parsed = {};
try {
  // Remove markdown code blocks if present
  let cleanJson = extractedText.replace(/```json\n?/g, '').replace(/```\n?/g, '').trim();
  parsed = JSON.parse(cleanJson);
} catch (e) {
  console.log('Failed to parse Gemini response:', extractedText);
  parsed = {
    phone_numbers: [],
    person_names: [],
    org_names: [],
    product_hints: [],
    priority_hints: []
  };
}

// Get original email data from the first item in the workflow
const emailData = $('Gmail Trigger').first().json;

// Build the request for smart-ticket endpoint
return {
  sender_email: emailData.from?.value?.[0]?.address || '',
  sender_name: emailData.from?.value?.[0]?.name || '',
  subject: emailData.subject || 'No Subject',
  body: emailData.text || emailData.snippet || '',
  
  // AI-extracted hints
  phone_numbers: parsed.phone_numbers || [],
  person_names: parsed.person_names || [],
  org_names: parsed.org_names || [],
  product_hints: parsed.product_hints || [],
  priority_hints: parsed.priority_hints || [],
  
  // Attachments (if any)
  attachments: (emailData.attachments || []).map(att => ({
    filename: att.filename || 'attachment',
    content_type: att.contentType || 'application/octet-stream',
    data: att.content ? att.content.toString('base64') : ''
  }))
};
```

### 3.6 Add HTTP Request Node (Create Ticket)

Add an **HTTP Request** node to call your backend:

```
Method: POST
URL: http://localhost:8080/n8n/smart-ticket

Headers:
  X-N8N-API-Key: your-api-key-here
  Content-Type: application/json

Body: {{ $json }}
```

### 3.7 Connect All Nodes

```
Gmail Trigger → HTTP Request (Gemini) → Code → HTTP Request (Backend)
```

### 3.8 Activate Workflow

1. Save the workflow
2. Click **"Activate"** toggle in the top right
3. The workflow is now running!

---

## Step 4: Test the Integration

### 4.1 Send a Test Email

Send an email to your Gmail with content like:

```
Subject: Dashboard not loading - urgent help needed

Hi Team,

We at Acme Corporation are facing issues with the dashboard. 
It's been down since morning and showing error 500.

Please call me at +91 9876543210 if you need more details.

Regards,
John Doe
IT Manager
Acme Corporation
```

### 4.2 Check n8n Execution

1. Go to n8n → Executions
2. You should see a successful execution
3. Check the output of each node

### 4.3 Verify Ticket Created

Check your backend logs or database for the new ticket.

---

## API Reference

### POST /n8n/smart-ticket

**Request:**
```json
{
  "sender_email": "john.doe@acmecorp.com",
  "sender_name": "John Doe",
  "subject": "Dashboard not loading",
  "body": "We at Acme Corporation are facing issues...",
  
  "phone_numbers": ["+919876543210"],
  "person_names": ["John Doe"],
  "org_names": ["Acme Corporation"],
  "product_hints": ["dashboard"],
  "priority_hints": ["urgent"],
  
  "attachments": []
}
```

**Success Response (201):**
```json
{
  "success": true,
  "ticket_id": "001-111225-0001",
  "resolution": {
    "contact": {
      "id": 5,
      "name": "John Doe",
      "method": "email_exact_match",
      "confidence": 100
    },
    "account": {
      "id": 1,
      "name": "Acme Corporation",
      "method": "from_contact",
      "confidence": 100
    },
    "product": {
      "id": 2,
      "name": "Dashboard",
      "method": "keyword_match",
      "confidence": 80
    }
  },
  "ticket": {
    "id": 42,
    "ticket_id": "001-111225-0001",
    "subject": "Dashboard not loading",
    "status": "OPEN",
    "priority": "High",
    "created_at": "2025-12-11T14:30:00Z"
  },
  "warnings": []
}
```

**Error Response (400):**
```json
{
  "success": false,
  "error": "Could not identify contact from email",
  "warnings": [
    "No contact could be identified from the email"
  ]
}
```

---

## Resolution Methods

The Smart Resolver tries these methods in order:

### Contact Resolution (Waterfall)

| Order | Method | Confidence | Description |
|-------|--------|------------|-------------|
| 1 | `email_exact_match` | 100% | Sender email matches contact.email |
| 2 | `phone_match` | 95% | Phone in email matches contact.mobile |
| 3 | `email_username_parse` | 60-80% | Parse "doejohn@" → "John Doe" |
| 4 | `ai_extracted_name` | 50-90% | Names extracted by Gemini |
| 5 | `org_match` | 50-70% | Find org → get contacts |
| 6 | `email_domain_match` | 40% | Match email domain to contacts |

### Account Resolution

| Order | Method | Confidence | Description |
|-------|--------|------------|-------------|
| 1 | `from_contact` | 100% | Use contact's account |
| 2 | `org_name_match` | 80% | Match org names from email |
| 3 | `email_domain_match` | 60% | Match domain to accounts |
| 4 | `no_account` | 100% | Individual contact |

### Product Resolution

| Order | Method | Confidence | Description |
|-------|--------|------------|-------------|
| 1 | `keyword_match` | 80% | Product hints match product name |
| 2 | `content_match` | 70% | Product name found in email |
| 3 | `default` | 50% | Use first/default product |

---

## Troubleshooting

### Issue: 401 Unauthorized

**Cause:** API key mismatch

**Solution:**
1. Check `N8N_API_KEY` in `.env.email`
2. Check `X-N8N-API-Key` header in n8n
3. Make sure they match exactly

### Issue: Contact Not Found

**Cause:** No matching contact in database

**Solution:**
1. The email sender must exist as a contact in your system
2. Or, the phone/name must match an existing contact
3. Check the `warnings` array in response for details

### Issue: Gemini Returns Invalid JSON

**Cause:** AI sometimes adds markdown or explanation

**Solution:** The Code node handles this by removing ```json``` blocks. If still failing:
1. Check Gemini response in n8n execution
2. Adjust the prompt to be more strict
3. Add more error handling in Code node

### Issue: Duplicate Tickets

**Cause:** Email processed multiple times

**Solution:**
1. Add a **Filter** node to check if email was already processed
2. Use Gmail's message ID as a unique identifier
3. Store processed IDs in a database or file

---

## Production Recommendations

### 1. Error Handling

Add an **IF** node after the HTTP Request to check for errors:

```
Condition: {{ $json.success }} equals true
True: Send success notification
False: Send alert to admin
```

### 2. Logging

Add a **Send Email** or **Slack** node to notify on:
- New tickets created
- Errors or low-confidence matches
- Unknown senders

### 3. Rate Limiting

Gemini API has rate limits. Add a **Wait** node if processing many emails.

### 4. Security

1. Use environment variables for API keys in n8n
2. Restrict n8n access with authentication
3. Use HTTPS in production
4. Rotate API keys periodically

---

## Complete n8n Workflow JSON

Import this directly into n8n:

```json
{
  "name": "Email to Ticket (Gemini AI)",
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
      "position": [250, 300],
      "typeVersion": 1
    },
    {
      "parameters": {
        "method": "POST",
        "url": "=https://generativelanguage.googleapis.com/v1beta/models/gemini-1.5-flash:generateContent?key={{ $env.GEMINI_API_KEY }}",
        "sendHeaders": true,
        "headerParameters": {
          "parameters": [
            {"name": "Content-Type", "value": "application/json"}
          ]
        },
        "sendBody": true,
        "specifyBody": "json",
        "jsonBody": "={\n  \"contents\": [{\n    \"parts\": [{\n      \"text\": \"You are an email parser. Extract from this email:\\n\\nFrom: {{ $json.from.value[0].name }} <{{ $json.from.value[0].address }}>\\nSubject: {{ $json.subject }}\\nBody:\\n{{ $json.text }}\\n\\nReturn ONLY JSON with: phone_numbers[], person_names[], org_names[], product_hints[], priority_hints[]\"\n    }]\n  }],\n  \"generationConfig\": {\"temperature\": 0.1}\n}"
      },
      "name": "Gemini Extract",
      "type": "n8n-nodes-base.httpRequest",
      "position": [470, 300],
      "typeVersion": 4.1
    },
    {
      "parameters": {
        "jsCode": "const geminiResponse = $input.first().json;\nlet extractedText = '';\nif (geminiResponse.candidates && geminiResponse.candidates[0]) {\n  extractedText = geminiResponse.candidates[0].content.parts[0].text;\n}\n\nlet parsed = {};\ntry {\n  let cleanJson = extractedText.replace(/```json\\n?/g, '').replace(/```\\n?/g, '').trim();\n  parsed = JSON.parse(cleanJson);\n} catch (e) {\n  parsed = { phone_numbers: [], person_names: [], org_names: [], product_hints: [], priority_hints: [] };\n}\n\nconst emailData = $('Gmail Trigger').first().json;\n\nreturn {\n  sender_email: emailData.from?.value?.[0]?.address || '',\n  sender_name: emailData.from?.value?.[0]?.name || '',\n  subject: emailData.subject || 'No Subject',\n  body: emailData.text || '',\n  phone_numbers: parsed.phone_numbers || [],\n  person_names: parsed.person_names || [],\n  org_names: parsed.org_names || [],\n  product_hints: parsed.product_hints || [],\n  priority_hints: parsed.priority_hints || [],\n  attachments: []\n};"
      },
      "name": "Prepare Request",
      "type": "n8n-nodes-base.code",
      "position": [690, 300],
      "typeVersion": 2
    },
    {
      "parameters": {
        "method": "POST",
        "url": "={{ $env.BACKEND_URL }}/n8n/smart-ticket",
        "sendHeaders": true,
        "headerParameters": {
          "parameters": [
            {"name": "X-N8N-API-Key", "value": "={{ $env.N8N_API_KEY }}"},
            {"name": "Content-Type", "value": "application/json"}
          ]
        },
        "sendBody": true,
        "specifyBody": "json",
        "jsonBody": "={{ $json }}"
      },
      "name": "Create Ticket",
      "type": "n8n-nodes-base.httpRequest",
      "position": [910, 300],
      "typeVersion": 4.1
    }
  ],
  "connections": {
    "Gmail Trigger": {"main": [[{"node": "Gemini Extract", "type": "main", "index": 0}]]},
    "Gemini Extract": {"main": [[{"node": "Prepare Request", "type": "main", "index": 0}]]},
    "Prepare Request": {"main": [[{"node": "Create Ticket", "type": "main", "index": 0}]]}
  },
  "settings": {
    "saveManualExecutions": true,
    "saveDataSuccessExecution": "all"
  }
}
```

### Environment Variables in n8n

Set these in n8n Settings → Variables:
- `GEMINI_API_KEY`: Your Gemini API key
- `BACKEND_URL`: http://localhost:8080
- `N8N_API_KEY`: Your backend API key

---

## Summary

| Component | Purpose |
|-----------|---------|
| **Gmail Trigger** | Receives new emails |
| **Gemini** | Extracts structured data from email |
| **Code Node** | Prepares request for backend |
| **HTTP Request** | Calls `/n8n/smart-ticket` endpoint |
| **Smart Resolver** | Intelligent contact/account matching |

The system handles:
- ✅ Exact email matches
- ✅ Phone number matches
- ✅ Name parsing from emails (doejohn@ → John Doe)
- ✅ Organization name matching
- ✅ Fuzzy name matching with confidence scores
- ✅ Priority detection from keywords
- ✅ Product detection from content
- ✅ Attachments
- ✅ Activity logging
- ✅ Email notifications
