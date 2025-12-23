# Usage Examples

This document provides practical examples of using the dago-node-planner.

## Table of Contents

- [Simple Examples](#simple-examples)
- [Moderate Examples](#moderate-examples)
- [Complex Examples](#complex-examples)
- [Client SDK Examples](#client-sdk-examples)
- [cURL Examples](#curl-examples)

## Simple Examples

### Example 1: Send Email

**Task**: "Send a welcome email to the new user"

**Request**:
```json
{
  "task": "Send a welcome email to the new user",
  "context": {
    "user_email": "newuser@example.com",
    "user_name": "John Doe"
  },
  "constraints": {
    "max_nodes": 5,
    "preferred_modes": ["tool"]
  }
}
```

**Generated Graph**:
```json
{
  "nodes": [
    {
      "id": "send_welcome_email",
      "type": "executor",
      "config": {
        "mode": "tool",
        "tool_calls": [
          {
            "tool_name": "send_email",
            "parameters": {
              "to": "$.user_email",
              "subject": "Welcome!",
              "body": "Welcome to our platform, $.user_name!"
            }
          }
        ]
      }
    }
  ],
  "edges": [],
  "entry_point": "send_welcome_email"
}
```

### Example 2: API Call

**Task**: "Fetch user profile from API and log the result"

**Generated Graph**:
```json
{
  "nodes": [
    {
      "id": "fetch_profile",
      "type": "executor",
      "config": {
        "mode": "tool",
        "tool_calls": [
          {
            "tool_name": "http_get",
            "parameters": {
              "url": "https://api.example.com/users/$.user_id"
            }
          }
        ],
        "state_output_path": "$.profile"
      }
    },
    {
      "id": "log_result",
      "type": "executor",
      "config": {
        "mode": "tool",
        "tool_calls": [
          {
            "tool_name": "log",
            "parameters": {
              "message": "Profile fetched: $.profile"
            }
          }
        ]
      }
    }
  ],
  "edges": [
    {
      "source": "fetch_profile",
      "target": "log_result"
    }
  ],
  "entry_point": "fetch_profile"
}
```

## Moderate Examples

### Example 3: Text Analysis

**Task**: "Analyze the text, extract key entities, and generate a summary"

**Generated Graph**:
```json
{
  "nodes": [
    {
      "id": "analyze_text",
      "type": "executor",
      "config": {
        "mode": "agent",
        "prompt": "Analyze the following text and extract key entities (people, places, organizations): $.input_text",
        "state_output_path": "$.entities",
        "available_tools": []
      }
    },
    {
      "id": "generate_summary",
      "type": "executor",
      "config": {
        "mode": "llm",
        "prompt": "Generate a concise summary of this text: $.input_text. Key entities found: $.entities",
        "state_output_path": "$.summary"
      }
    }
  ],
  "edges": [
    {
      "source": "analyze_text",
      "target": "generate_summary"
    }
  ],
  "entry_point": "analyze_text"
}
```

### Example 4: Data Transformation

**Task**: "Fetch data from database, transform it to JSON format, and upload to S3"

**Generated Graph**:
```json
{
  "nodes": [
    {
      "id": "fetch_data",
      "type": "executor",
      "config": {
        "mode": "tool",
        "tool_calls": [
          {
            "tool_name": "database_query",
            "parameters": {
              "query": "SELECT * FROM users WHERE active = true"
            }
          }
        ],
        "state_output_path": "$.raw_data"
      }
    },
    {
      "id": "transform_data",
      "type": "executor",
      "config": {
        "mode": "agent",
        "prompt": "Transform this database result to JSON format: $.raw_data",
        "state_output_path": "$.json_data",
        "available_tools": ["json_formatter"]
      }
    },
    {
      "id": "upload_to_s3",
      "type": "executor",
      "config": {
        "mode": "tool",
        "tool_calls": [
          {
            "tool_name": "s3_upload",
            "parameters": {
              "bucket": "data-exports",
              "key": "users-$.timestamp.json",
              "content": "$.json_data"
            }
          }
        ]
      }
    }
  ],
  "edges": [
    {
      "source": "fetch_data",
      "target": "transform_data"
    },
    {
      "source": "transform_data",
      "target": "upload_to_s3"
    }
  ],
  "entry_point": "fetch_data"
}
```

## Complex Examples

### Example 5: Sentiment-Based Routing

**Task**: "Analyze customer feedback sentiment and route to the appropriate team. Escalate negative feedback to manager."

**Generated Graph**: See `examples/complex-task.json`

Key features:
- Sentiment analysis with agent mode
- Categorization with LLM mode
- Deterministic routing based on sentiment score
- Multiple destination nodes
- Escalation logic

### Example 6: Multi-Step Approval Workflow

**Task**: "Review purchase request, check budget, route to approver based on amount, send notification"

**Generated Graph**:
```json
{
  "nodes": [
    {
      "id": "review_request",
      "type": "executor",
      "config": {
        "mode": "agent",
        "prompt": "Review this purchase request and extract: amount, category, justification: $.request",
        "state_output_path": "$.analysis",
        "available_tools": []
      }
    },
    {
      "id": "check_budget",
      "type": "executor",
      "config": {
        "mode": "tool",
        "tool_calls": [
          {
            "tool_name": "database_query",
            "parameters": {
              "query": "SELECT remaining_budget FROM budgets WHERE category = $.analysis.category"
            }
          }
        ],
        "state_output_path": "$.budget_check"
      }
    },
    {
      "id": "amount_router",
      "type": "router",
      "config": {
        "mode": "deterministic",
        "routes": [
          {
            "condition": "$.analysis.amount > 10000",
            "target": "cfo_approval"
          },
          {
            "condition": "$.analysis.amount > 5000",
            "target": "director_approval"
          },
          {
            "condition": "true",
            "target": "manager_approval"
          }
        ]
      }
    },
    {
      "id": "cfo_approval",
      "type": "executor",
      "config": {
        "mode": "tool",
        "tool_calls": [
          {
            "tool_name": "send_email",
            "parameters": {
              "to": "cfo@company.com",
              "subject": "High-value purchase approval needed",
              "body": "Amount: $.analysis.amount, Category: $.analysis.category"
            }
          }
        ]
      }
    },
    {
      "id": "director_approval",
      "type": "executor",
      "config": {
        "mode": "tool",
        "tool_calls": [
          {
            "tool_name": "send_email",
            "parameters": {
              "to": "director@company.com",
              "subject": "Purchase approval needed",
              "body": "Amount: $.analysis.amount, Category: $.analysis.category"
            }
          }
        ]
      }
    },
    {
      "id": "manager_approval",
      "type": "executor",
      "config": {
        "mode": "tool",
        "tool_calls": [
          {
            "tool_name": "send_email",
            "parameters": {
              "to": "manager@company.com",
              "subject": "Purchase approval needed",
              "body": "Amount: $.analysis.amount, Category: $.analysis.category"
            }
          }
        ]
      }
    }
  ],
  "edges": [
    {
      "source": "review_request",
      "target": "check_budget"
    },
    {
      "source": "check_budget",
      "target": "amount_router"
    }
  ],
  "entry_point": "review_request"
}
```

## Client SDK Examples

### Go Client

```go
package main

import (
    "context"
    "fmt"
    "log"

    "github.com/aescanero/dago-node-planner/pkg/client"
    "github.com/aescanero/dago-node-planner/pkg/models"
)

func main() {
    // Create client
    plannerClient := client.NewClient("http://localhost:8080")

    // Create plan request
    req := &models.PlanRequest{
        Task: "Analyze user sentiment and send personalized response",
        Context: map[string]any{
            "user_id": "12345",
        },
        Constraints: &models.Constraints{
            MaxNodes:       10,
            PreferredModes: []string{"agent", "llm"},
            AvailableTools: []string{"send_email", "update_database"},
        },
    }

    // Generate plan
    resp, err := plannerClient.Plan(context.Background(), req)
    if err != nil {
        log.Fatal(err)
    }

    // Print results
    fmt.Printf("Plan ID: %s\n", resp.PlanID)
    fmt.Printf("Iterations: %d\n", resp.Iterations)
    fmt.Printf("Tokens used: %d\n", resp.Metadata.TokensUsed)
    fmt.Printf("Graph: %v\n", resp.Graph)
}
```

### Health Check

```go
// Check if service is healthy
err := plannerClient.Health(context.Background())
if err != nil {
    log.Fatal("Service is unhealthy:", err)
}
```

### Validate Graph

```go
graphJSON := `{"nodes": [...], "edges": [...]}`

result, err := plannerClient.Validate(context.Background(), graphJSON)
if err != nil {
    log.Fatal(err)
}

if !result.Valid {
    fmt.Println("Validation errors:")
    for _, err := range result.Errors {
        fmt.Printf("  - %s\n", err)
    }
}
```

## cURL Examples

### Generate Plan

```bash
curl -X POST http://localhost:8080/api/v1/plan \
  -H "Content-Type: application/json" \
  -d '{
    "task": "Send notification and log the event",
    "constraints": {
      "max_nodes": 5
    }
  }'
```

### Validate Graph

```bash
curl -X POST http://localhost:8080/api/v1/validate \
  -H "Content-Type: application/json" \
  -d '{
    "graph_json": "{\"nodes\": [], \"edges\": [], \"entry_point\": \"start\"}"
  }'
```

### Health Check

```bash
curl http://localhost:8080/health
```

## Tips

### Effective Task Descriptions

✅ Good:
```
"Analyze customer feedback, categorize by topic, and create support tickets for issues"
```

❌ Bad:
```
"Do something with customer data"
```

### Using Context

Provide relevant context data:
```json
{
  "task": "Process order",
  "context": {
    "order_id": "12345",
    "customer_tier": "premium",
    "order_value": 1500
  }
}
```

### Setting Constraints

Guide the planner with constraints:
```json
{
  "task": "...",
  "constraints": {
    "max_nodes": 5,
    "preferred_modes": ["tool"],
    "available_tools": ["send_email", "create_ticket"]
  }
}
```
