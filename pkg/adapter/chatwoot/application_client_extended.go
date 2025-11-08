// SPDX-License-Identifier: AGPL-3.0-or-later
// Copyright (C) 2025 Controle Digital Ltda

package chatwoot

import (
	"context"
	"fmt"
	"net/http"

	"github.com/click2-run/dictamesh/pkg/adapter"
)

// Inbox Management

// ListInboxes lists all inboxes
func (c *ApplicationClient) ListInboxes(ctx context.Context) ([]Inbox, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inboxes", c.accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload []Inbox `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Payload, nil
}

// GetInbox retrieves an inbox by ID
func (c *ApplicationClient) GetInbox(ctx context.Context, inboxID int64) (*Inbox, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inboxes/%d", c.accountID, inboxID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result Inbox
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateInbox creates a new inbox
func (c *ApplicationClient) CreateInbox(ctx context.Context, inbox *Inbox) (*Inbox, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inboxes", c.accountID)

	resp, err := c.httpClient.Post(ctx, path, inbox, nil)
	if err != nil {
		return nil, err
	}

	var result Inbox
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateInbox updates an inbox
func (c *ApplicationClient) UpdateInbox(ctx context.Context, inboxID int64, inbox *Inbox) (*Inbox, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/inboxes/%d", c.accountID, inboxID)

	resp, err := c.httpClient.Patch(ctx, path, inbox, nil)
	if err != nil {
		return nil, err
	}

	var result Inbox
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteInbox deletes an inbox
func (c *ApplicationClient) DeleteInbox(ctx context.Context, inboxID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/inboxes/%d", c.accountID, inboxID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete inbox", nil)
	}

	return nil
}

// Team Management

// ListTeams lists all teams
func (c *ApplicationClient) ListTeams(ctx context.Context) ([]Team, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams", c.accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result []Team
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateTeam creates a new team
func (c *ApplicationClient) CreateTeam(ctx context.Context, team *Team) (*Team, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams", c.accountID)

	resp, err := c.httpClient.Post(ctx, path, team, nil)
	if err != nil {
		return nil, err
	}

	var result Team
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateTeam updates a team
func (c *ApplicationClient) UpdateTeam(ctx context.Context, teamID int64, team *Team) (*Team, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams/%d", c.accountID, teamID)

	resp, err := c.httpClient.Patch(ctx, path, team, nil)
	if err != nil {
		return nil, err
	}

	var result Team
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteTeam deletes a team
func (c *ApplicationClient) DeleteTeam(ctx context.Context, teamID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/teams/%d", c.accountID, teamID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete team", nil)
	}

	return nil
}

// Label Management

// ListLabels lists all labels
func (c *ApplicationClient) ListLabels(ctx context.Context) ([]Label, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/labels", c.accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload []Label `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Payload, nil
}

// CreateLabel creates a new label
func (c *ApplicationClient) CreateLabel(ctx context.Context, label *Label) (*Label, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/labels", c.accountID)

	resp, err := c.httpClient.Post(ctx, path, label, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload Label `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Payload, nil
}

// DeleteLabel deletes a label
func (c *ApplicationClient) DeleteLabel(ctx context.Context, labelID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/labels/%d", c.accountID, labelID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete label", nil)
	}

	return nil
}

// ListConversationLabels lists labels for a conversation
func (c *ApplicationClient) ListConversationLabels(ctx context.Context, conversationID int64) ([]string, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/labels", c.accountID, conversationID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload []string `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Payload, nil
}

// AddConversationLabels adds labels to a conversation
func (c *ApplicationClient) AddConversationLabels(ctx context.Context, conversationID int64, labels []string) ([]string, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/conversations/%d/labels", c.accountID, conversationID)

	payload := map[string]interface{}{
		"labels": labels,
	}

	resp, err := c.httpClient.Post(ctx, path, payload, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload []string `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Payload, nil
}

// Canned Response Management

// ListCannedResponses lists all canned responses
func (c *ApplicationClient) ListCannedResponses(ctx context.Context) ([]CannedResponse, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/canned_responses", c.accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result []CannedResponse
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateCannedResponse creates a new canned response
func (c *ApplicationClient) CreateCannedResponse(ctx context.Context, cannedResponse *CannedResponse) (*CannedResponse, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/canned_responses", c.accountID)

	resp, err := c.httpClient.Post(ctx, path, cannedResponse, nil)
	if err != nil {
		return nil, err
	}

	var result CannedResponse
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateCannedResponse updates a canned response
func (c *ApplicationClient) UpdateCannedResponse(ctx context.Context, responseID int64, cannedResponse *CannedResponse) (*CannedResponse, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/canned_responses/%d", c.accountID, responseID)

	resp, err := c.httpClient.Patch(ctx, path, cannedResponse, nil)
	if err != nil {
		return nil, err
	}

	var result CannedResponse
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteCannedResponse deletes a canned response
func (c *ApplicationClient) DeleteCannedResponse(ctx context.Context, responseID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/canned_responses/%d", c.accountID, responseID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete canned response", nil)
	}

	return nil
}

// Custom Attribute Management

// ListCustomAttributeDefinitions lists all custom attribute definitions
func (c *ApplicationClient) ListCustomAttributeDefinitions(ctx context.Context) ([]CustomAttributeDefinition, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_attribute_definitions", c.accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result []CustomAttributeDefinition
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateCustomAttributeDefinition creates a new custom attribute definition
func (c *ApplicationClient) CreateCustomAttributeDefinition(ctx context.Context, definition *CustomAttributeDefinition) (*CustomAttributeDefinition, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_attribute_definitions", c.accountID)

	resp, err := c.httpClient.Post(ctx, path, definition, nil)
	if err != nil {
		return nil, err
	}

	var result CustomAttributeDefinition
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetCustomAttributeDefinition retrieves a custom attribute definition
func (c *ApplicationClient) GetCustomAttributeDefinition(ctx context.Context, definitionID int64) (*CustomAttributeDefinition, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_attribute_definitions/%d", c.accountID, definitionID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result CustomAttributeDefinition
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateCustomAttributeDefinition updates a custom attribute definition
func (c *ApplicationClient) UpdateCustomAttributeDefinition(ctx context.Context, definitionID int64, definition *CustomAttributeDefinition) (*CustomAttributeDefinition, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_attribute_definitions/%d", c.accountID, definitionID)

	resp, err := c.httpClient.Patch(ctx, path, definition, nil)
	if err != nil {
		return nil, err
	}

	var result CustomAttributeDefinition
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteCustomAttributeDefinition deletes a custom attribute definition
func (c *ApplicationClient) DeleteCustomAttributeDefinition(ctx context.Context, definitionID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/custom_attribute_definitions/%d", c.accountID, definitionID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete custom attribute definition", nil)
	}

	return nil
}

// Automation Rule Management

// ListAutomationRules lists all automation rules
func (c *ApplicationClient) ListAutomationRules(ctx context.Context) ([]AutomationRule, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/automation_rules", c.accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload []AutomationRule `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Payload, nil
}

// CreateAutomationRule creates a new automation rule
func (c *ApplicationClient) CreateAutomationRule(ctx context.Context, rule *AutomationRule) (*AutomationRule, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/automation_rules", c.accountID)

	resp, err := c.httpClient.Post(ctx, path, rule, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload AutomationRule `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Payload, nil
}

// GetAutomationRule retrieves an automation rule
func (c *ApplicationClient) GetAutomationRule(ctx context.Context, ruleID int64) (*AutomationRule, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/automation_rules/%d", c.accountID, ruleID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload AutomationRule `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Payload, nil
}

// UpdateAutomationRule updates an automation rule
func (c *ApplicationClient) UpdateAutomationRule(ctx context.Context, ruleID int64, rule *AutomationRule) (*AutomationRule, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/automation_rules/%d", c.accountID, ruleID)

	resp, err := c.httpClient.Patch(ctx, path, rule, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload AutomationRule `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Payload, nil
}

// DeleteAutomationRule deletes an automation rule
func (c *ApplicationClient) DeleteAutomationRule(ctx context.Context, ruleID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/automation_rules/%d", c.accountID, ruleID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete automation rule", nil)
	}

	return nil
}

// Webhook Management

// ListWebhooks lists all webhooks
func (c *ApplicationClient) ListWebhooks(ctx context.Context) ([]Webhook, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/webhooks", c.accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload []Webhook `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Payload, nil
}

// CreateWebhook creates a new webhook
func (c *ApplicationClient) CreateWebhook(ctx context.Context, webhook *Webhook) (*Webhook, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/webhooks", c.accountID)

	resp, err := c.httpClient.Post(ctx, path, webhook, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload Webhook `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Payload, nil
}

// UpdateWebhook updates a webhook
func (c *ApplicationClient) UpdateWebhook(ctx context.Context, webhookID int64, webhook *Webhook) (*Webhook, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/webhooks/%d", c.accountID, webhookID)

	resp, err := c.httpClient.Patch(ctx, path, webhook, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload Webhook `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Payload, nil
}

// DeleteWebhook deletes a webhook
func (c *ApplicationClient) DeleteWebhook(ctx context.Context, webhookID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/webhooks/%d", c.accountID, webhookID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete webhook", nil)
	}

	return nil
}

// Integration Management

// ListIntegrationApps lists all available integration apps
func (c *ApplicationClient) ListIntegrationApps(ctx context.Context) ([]IntegrationApp, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/integrations/apps", c.accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload []IntegrationApp `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Payload, nil
}

// ListIntegrations lists all enabled integrations
func (c *ApplicationClient) ListIntegrations(ctx context.Context) ([]IntegrationHook, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/integrations", c.accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload []IntegrationHook `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result.Payload, nil
}

// EnableIntegration enables an integration
func (c *ApplicationClient) EnableIntegration(ctx context.Context, appID string, settings map[string]interface{}) (*IntegrationHook, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/integrations", c.accountID)

	payload := map[string]interface{}{
		"app_id":   appID,
		"settings": settings,
	}

	resp, err := c.httpClient.Post(ctx, path, payload, nil)
	if err != nil {
		return nil, err
	}

	var result struct {
		Payload IntegrationHook `json:"payload"`
	}
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result.Payload, nil
}

// Audit Log Management

// ListAuditLogs lists audit logs for the account
func (c *ApplicationClient) ListAuditLogs(ctx context.Context, page int) (*ListResponse, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/audit_logs", c.accountID)

	if page > 0 {
		path = fmt.Sprintf("%s?page=%d", path, page)
	}

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result ListResponse
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Reports and Analytics

// GetAccountReports retrieves account-level reports
func (c *ApplicationClient) GetAccountReports(ctx context.Context, metricType string, since string, until string) (*Report, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/reports?type=%s&since=%s&until=%s", c.accountID, metricType, since, until)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result Report
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetConversationMetrics retrieves conversation metrics
func (c *ApplicationClient) GetConversationMetrics(ctx context.Context, metricsType string) (*ConversationMetrics, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/reports/conversations?type=%s", c.accountID, metricsType)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result ConversationMetrics
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// Agent Bot Management (Account-level)

// ListAccountAgentBots lists all agent bots in the account
func (c *ApplicationClient) ListAccountAgentBots(ctx context.Context) ([]AgentBot, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agent_bots", c.accountID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result []AgentBot
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// CreateAccountAgentBot creates an agent bot in the account
func (c *ApplicationClient) CreateAccountAgentBot(ctx context.Context, bot *AgentBot) (*AgentBot, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agent_bots", c.accountID)

	resp, err := c.httpClient.Post(ctx, path, bot, nil)
	if err != nil {
		return nil, err
	}

	var result AgentBot
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAccountAgentBot retrieves an agent bot from the account
func (c *ApplicationClient) GetAccountAgentBot(ctx context.Context, botID int64) (*AgentBot, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agent_bots/%d", c.accountID, botID)

	resp, err := c.httpClient.Get(ctx, path, nil)
	if err != nil {
		return nil, err
	}

	var result AgentBot
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateAccountAgentBot updates an agent bot in the account
func (c *ApplicationClient) UpdateAccountAgentBot(ctx context.Context, botID int64, bot *AgentBot) (*AgentBot, error) {
	path := fmt.Sprintf("/api/v1/accounts/%d/agent_bots/%d", c.accountID, botID)

	resp, err := c.httpClient.Patch(ctx, path, bot, nil)
	if err != nil {
		return nil, err
	}

	var result AgentBot
	if err := adapter.ParseJSONResponse(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteAccountAgentBot deletes an agent bot from the account
func (c *ApplicationClient) DeleteAccountAgentBot(ctx context.Context, botID int64) error {
	path := fmt.Sprintf("/api/v1/accounts/%d/agent_bots/%d", c.accountID, botID)

	resp, err := c.httpClient.Delete(ctx, path, nil)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusNoContent {
		return adapter.HTTPErrorToAdapterError(resp.StatusCode, "failed to delete agent bot", nil)
	}

	return nil
}
