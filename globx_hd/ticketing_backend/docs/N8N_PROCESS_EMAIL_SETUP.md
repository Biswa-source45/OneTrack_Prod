# n8n Email-to-Ticket Integration (Backend Gemini Processing)

Complete guide for the new architecture where n8n sends emails to the backend, and the backend handles all Gemini AI processing.

---

## Architecture

```
┌──────────────┐    ┌─────────────────────────────────────────────────────────────────────┐
│    EMAIL     │    │                              n8n                                    │
│    INBOX     │───▶│  ┌──────────┐   ┌──────────────┐   ┌────────────────────────────┐  │
│  (Gmail)     │    │  │  Gmail   │──▶│  Set/Code    │──▶│  HTTP Request              │  │
└──────────────┘    │  │  Trigger │   │  (Format)    │   │  POST /n8n/process-email   │  │
                    │  └──────────┘   └──────────────┘   └────────────────────────────┘  │
                    │                                                                     │
                    │  Output: { "prompt": "...", "id": "email_id" }                      │
                    └─────────────────────────────────────────────────────────────────────┘
                                                    │
                                                    ▼
┌─────────────────────────────────────────────────────────────────────────────────────────┐
│                              GO BACKEND                                                 │
│                                                                                         │
│  POST /n8n/process-email                                                                │
│  ┌───────────────────────────────────────────────────────────────────────────────────┐ │
│  │  1. Parse Email from Prompt                                                       │ │
│  │     • Extract sender email, name, subject, body                                   │ │
│  │                                                                                   │ │
│  │  2. Call Gemini 2.5 Flash API                                                     │ │
│  │     • Send prompt to Gemini                                                       │ │
│  │     • Get JSON response with extracted entities                                   │ │
│  │                                                                                   │ │
│  │  3. Smart Resolver                                                                │ │
│  │     • Waterfall contact resolution (email → phone → name → org → domain)         │ │
│  │     • Account resolution                                                          │ │
│  │     • Product matching                                                            │ │
│  │     • Priority detection                                                          │ │
│  │                                                                                   │ │
│  │  4. Create Ticket                                                                 │ │
│  │     • Generate ticket ID                                                          │ │
│  │     • Save to database                                                            │ │
│  │     • Log activity                                                                │ │
│  │     • Send notifications                                                          │ │
│  └───────────────────────────────────────────────────────────────────────────────────┘ │
│                                                                                         │
│  Response: { "success": true, "id": "email_id", "ticket_id": "001-151225-0001", ... }   │
└─────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## Configuration

### 1. Backend Configuration (`.env.email`)

```bash
# n8n API Key for authentication
N8N_API_KEY=your-secure-api-key-here

# Gemini API Key (get from https://aistudio.google.com/app/apikey)
GEMINI_API_KEY=your-gemini-api-key-here

# SMTP for notifications
EMAIL_SMTP_SERVER=smtp.gmail.com
EMAIL_SMTP_PORT=587
EMAIL_USERNAME=your-email@gmail.com
EMAIL_PASSWORD=your-app-password
EMAIL_NOTIFICATION_ADDRESS=notifications@yourcompany.com
```

### 2. Start Backend

```bash
cd ticketing_backend
go run cmd/main.go
```

### 3. Test Endpoint

```bash
curl -X POST http://localhost:8080/n8n/process-email \
  -H "X-N8N-API-Key: your-api-key" \
  -H "Content-Type: application/json" \
  -d '{
    "prompt": "You are an email parser for a ticketing system. Extract structured data from the email below.  From: \"John Doe\" <john.doe@acmecorp.com> Subject: Dashboard not working Body: Hi Team, Our dashboard is showing error 500. Please help urgently. Contact: 9876543210. Regards, John Doe  Extract and return ONLY a valid JSON object with these fields: - phone_numbers: array of all phone numbers found - person_names: array of all person names found - org_names: array of organization names - product_hints: array of product names mentioned - priority_hints: array of urgency keywords",
    "id": "test-email-123"
  }'
```

---

## API Reference

### POST /n8n/process-email

The primary endpoint for email processing.

**Request:**
```json
{
  "prompt": "You are an email parser for a ticketing system. Extract structured data from the email below.  From: \"Prabir globx\" <prabir.globx@gmail.com> Subject: Urgent: Dell Monitor Not Functioning Body: Dear Team, One of our Dell monitors is not functioning. Contact Person Name- R.K. Sahoo Contact Number-9124668108  Extract and return ONLY a valid JSON object with these fields: - phone_numbers: array - person_names: array - org_names: array - product_hints: array - priority_hints: array",
  "id": "19b0c419f18da347"
}
```

**Success Response (201):**
```json
{
  "success": true,
  "id": "19b0c419f18da347",
  "ticket_id": "001-151225-0001",
  "message": "Ticket 001-151225-0001 created successfully",
  "resolution": {
    "contact": {
      "id": 5,
      "name": "R.K. Sahoo",
      "method": "phone_match",
      "confidence": 95
    },
    "account": {
      "id": 1,
      "name": "District Court Jagatsinghpur",
      "method": "from_contact",
      "confidence": 100
    },
    "product": {
      "id": 2,
      "name": "Dell Monitor",
      "method": "keyword_match",
      "confidence": 80
    }
  },
  "ticket": {
    "id": 42,
    "ticket_id": "001-151225-0001",
    "subject": "Urgent: Dell Monitor Not Functioning",
    "status": "OPEN",
    "priority": "High",
    "created_at": "2025-12-15T12:30:00Z"
  },
  "warnings": []
}
```

**Error Response (400):**
```json
{
  "success": false,
  "id": "19b0c419f18da347",
  "error": "Could not identify contact from email",
  "warnings": ["No contact could be identified from the email"],
  "debug": {
    "parsed_email": {
      "sender_email": "unknown@example.com",
      "sender_name": "Unknown",
      "subject": "...",
      "body_preview": "..."
    },
    "extracted_data": {
      "phone_numbers": [],
      "person_names": [],
      "org_names": [],
      "product_hints": [],
      "priority_hints": []
    }
  }
}
```

---

## n8n Workflow Setup

### Step 1: Gmail Trigger

1. Add **Gmail Trigger** node
2. Configure:
   - **Poll Times**: Every 1 minute
   - **Simple**: OFF (to get full email data)

### Step 2: Set/Code Node (Format Request)

Add a **Set** or **Code** node to format the request:

**Using Code Node:**
```javascript
const email = $input.first().json;

// Build the prompt
const prompt = `You are an email parser for a ticketing system. Extract structured data from the email below.  From: "${email.from?.value?.[0]?.name || ''}" <${email.from?.value?.[0]?.address || ''}> Subject: ${email.subject || ''} Body: ${email.text || email.snippet || ''}  Extract and return ONLY a valid JSON object with these fields: - phone_numbers: array of all phone numbers found in the email - person_names: array of all person names found (sender, recipients, or any names mentioned in the text) - org_names: array of organization or company names mentioned - product_hints: array of product, brand, device, software, or service names mentioned - priority_hints: array of urgency keywords or phrases indicating priority (urgent, asap, immediate, critical, high-priority, important)  If a field has no values, use an empty array []. Return strictly valid JSON only.`;

return {
  prompt: prompt,
  id: email.id || email.messageId || `email-${Date.now()}`
};
```

**Using Set Node:**
- **prompt**: 
  ```
  You are an email parser for a ticketing system. Extract structured data from the email below.  From: "{{ $json.from.value[0].name }}" <{{ $json.from.value[0].address }}> Subject: {{ $json.subject }} Body: {{ $json.text }}  Extract and return ONLY a valid JSON object with these fields: - phone_numbers: array of all phone numbers found - person_names: array of all person names found - org_names: array of organization names - product_hints: array of product names mentioned - priority_hints: array of urgency keywords  If a field has no values, use an empty array []. Return strictly valid JSON only.
  ```
- **id**: `{{ $json.id }}`

### Step 3: HTTP Request Node

Configure HTTP Request to call the backend:

```
Method: POST
URL: http://localhost:8080/n8n/process-email

Headers:
  X-N8N-API-Key: your-api-key-here
  Content-Type: application/json

Body: JSON
  {
    "prompt": "{{ $json.prompt }}",
    "id": "{{ $json.id }}"
  }
```

### Step 4: Workflow Connection

```
Gmail Trigger → Set/Code Node → HTTP Request
```

### Step 5: Activate

Save and activate the workflow.

---

## Contact Resolution Strategies

The backend uses a waterfall strategy to find the best contact match:

| Order | Strategy | Confidence | Description |
|-------|----------|------------|-------------|
| 1 | Email exact match | 100% | Sender email matches `contacts.email` |
| 2 | Phone match | 95% | Phone from email matches `contacts.mobile` |
| 3 | Email username parse | 60-80% | Parse `doejohn@` → search "John Doe" |
| 4 | AI extracted names | 50-90% | Names extracted by Gemini → fuzzy search |
| 5 | Organization match | 50-70% | Find org → get its contacts |
| 6 | Email domain match | 40% | Match sender domain to contacts |

---

## Troubleshooting

### Issue: 401 Unauthorized
- Check `N8N_API_KEY` in `.env.email`
- Check `X-N8N-API-Key` header in n8n HTTP Request

### Issue: Gemini API Error
- Check `GEMINI_API_KEY` in `.env.email`
- Verify API key at https://aistudio.google.com/app/apikey
- Check Gemini API quota/limits

### Issue: Contact Not Found
- Ensure sender email exists as a contact in database
- Or phone number matches an existing contact
- Check `debug` field in response for parsed data

### Issue: Empty Gemini Response
- The system will continue with email parsing only
- Check backend logs for Gemini error details

---

## Files Overview

| File | Purpose |
|------|---------|
| `internal/services/gemini_service.go` | Gemini API client |
| `internal/handlers/n8n_smart_resolver.go` | Smart resolver + ProcessEmailHandler |
| `internal/routes/routes.go` | Route registration |
| `.env.email` | Configuration |

---

## Response to n8n

The response always includes the same `id` that was sent, allowing n8n to track which email was processed:

```json
{
  "success": true,
  "id": "19b0c419f18da347",  // Same ID from request
  "ticket_id": "001-151225-0001",
  "message": "Ticket 001-151225-0001 created successfully"
}
```

You can use this in n8n to:
- Mark email as processed
- Send confirmation notification
- Update external systems
- Handle errors appropriately
