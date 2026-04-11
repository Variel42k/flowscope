package storage

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/flowscope/flowscope/internal/model"
)

func (r *Repository) ListAlertRules(ctx context.Context) ([]model.AlertRule, error) {
	rows, err := r.db.QueryContext(ctx, `
SELECT rule_id, name, rule_type, enabled, threshold_value, window_minutes, severity, created_by, created_at, updated_at
FROM (
	SELECT rule_id,
	       argMax(name, updated_at) AS name,
	       argMax(rule_type, updated_at) AS rule_type,
	       argMax(enabled, updated_at) AS enabled,
	       argMax(threshold_value, updated_at) AS threshold_value,
	       argMax(window_minutes, updated_at) AS window_minutes,
	       argMax(severity, updated_at) AS severity,
	       argMax(created_by, updated_at) AS created_by,
	       min(created_at) AS created_at,
	       max(updated_at) AS updated_at,
	       argMax(deleted, updated_at) AS deleted
	FROM alert_rules
	GROUP BY rule_id
)
WHERE deleted = 0
ORDER BY updated_at DESC`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.AlertRule, 0, 16)
	for rows.Next() {
		var item model.AlertRule
		if err := rows.Scan(
			&item.RuleID,
			&item.Name,
			&item.RuleType,
			&item.Enabled,
			&item.ThresholdValue,
			&item.WindowMinutes,
			&item.Severity,
			&item.CreatedBy,
			&item.CreatedAt,
			&item.UpdatedAt,
		); err != nil {
			return nil, err
		}
		out = append(out, normalizeAlertRule(item))
	}
	return out, nil
}

func (r *Repository) CreateAlertRule(ctx context.Context, rule model.AlertRule) (model.AlertRule, error) {
	now := time.Now().UTC()
	rule = normalizeAlertRule(rule)
	if strings.TrimSpace(rule.RuleID) == "" {
		rule.RuleID = uuid.NewString()
	}
	if rule.CreatedAt.IsZero() {
		rule.CreatedAt = now
	}
	rule.UpdatedAt = now
	_, err := r.db.ExecContext(ctx, `
INSERT INTO alert_rules (
	rule_id, name, rule_type, enabled, threshold_value, window_minutes, severity, created_by, created_at, updated_at, deleted
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rule.RuleID,
		rule.Name,
		rule.RuleType,
		rule.Enabled,
		rule.ThresholdValue,
		rule.WindowMinutes,
		rule.Severity,
		rule.CreatedBy,
		rule.CreatedAt,
		rule.UpdatedAt,
		false,
	)
	return rule, err
}

func (r *Repository) UpdateAlertRule(ctx context.Context, ruleID string, patch model.AlertRule) (model.AlertRule, error) {
	existing, err := r.GetAlertRuleByID(ctx, ruleID)
	if err != nil {
		return model.AlertRule{}, err
	}
	if strings.TrimSpace(patch.Name) != "" {
		existing.Name = strings.TrimSpace(patch.Name)
	}
	if strings.TrimSpace(patch.RuleType) != "" {
		existing.RuleType = strings.TrimSpace(patch.RuleType)
	}
	if patch.WindowMinutes > 0 {
		existing.WindowMinutes = patch.WindowMinutes
	}
	if patch.ThresholdValue > 0 {
		existing.ThresholdValue = patch.ThresholdValue
	}
	if strings.TrimSpace(patch.Severity) != "" {
		existing.Severity = strings.ToLower(strings.TrimSpace(patch.Severity))
	}
	// bool zero value may be intentional, so trust payload when rule id is present
	existing.Enabled = patch.Enabled
	existing.UpdatedAt = time.Now().UTC()
	existing = normalizeAlertRule(existing)
	_, err = r.db.ExecContext(ctx, `
INSERT INTO alert_rules (
	rule_id, name, rule_type, enabled, threshold_value, window_minutes, severity, created_by, created_at, updated_at, deleted
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		existing.RuleID,
		existing.Name,
		existing.RuleType,
		existing.Enabled,
		existing.ThresholdValue,
		existing.WindowMinutes,
		existing.Severity,
		existing.CreatedBy,
		existing.CreatedAt,
		existing.UpdatedAt,
		false,
	)
	return existing, err
}

func (r *Repository) DeleteAlertRule(ctx context.Context, ruleID string, _ string) error {
	existing, err := r.GetAlertRuleByID(ctx, ruleID)
	if err != nil {
		return err
	}
	_, err = r.db.ExecContext(ctx, `
INSERT INTO alert_rules (
	rule_id, name, rule_type, enabled, threshold_value, window_minutes, severity, created_by, created_at, updated_at, deleted
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		existing.RuleID,
		existing.Name,
		existing.RuleType,
		existing.Enabled,
		existing.ThresholdValue,
		existing.WindowMinutes,
		existing.Severity,
		existing.CreatedBy,
		existing.CreatedAt,
		time.Now().UTC(),
		true,
	)
	return err
}

func (r *Repository) GetAlertRuleByID(ctx context.Context, ruleID string) (model.AlertRule, error) {
	row := r.db.QueryRowContext(ctx, `
SELECT rule_id, name, rule_type, enabled, threshold_value, window_minutes, severity, created_by, created_at, updated_at, deleted
FROM (
	SELECT rule_id,
	       argMax(name, updated_at) AS name,
	       argMax(rule_type, updated_at) AS rule_type,
	       argMax(enabled, updated_at) AS enabled,
	       argMax(threshold_value, updated_at) AS threshold_value,
	       argMax(window_minutes, updated_at) AS window_minutes,
	       argMax(severity, updated_at) AS severity,
	       argMax(created_by, updated_at) AS created_by,
	       min(created_at) AS created_at,
	       max(updated_at) AS updated_at,
	       argMax(deleted, updated_at) AS deleted
	FROM alert_rules
	WHERE rule_id = ?
	GROUP BY rule_id
)
LIMIT 1`, ruleID)
	var out model.AlertRule
	var deleted bool
	if err := row.Scan(
		&out.RuleID,
		&out.Name,
		&out.RuleType,
		&out.Enabled,
		&out.ThresholdValue,
		&out.WindowMinutes,
		&out.Severity,
		&out.CreatedBy,
		&out.CreatedAt,
		&out.UpdatedAt,
		&deleted,
	); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.AlertRule{}, errors.New("alert rule not found")
		}
		return model.AlertRule{}, err
	}
	if deleted {
		return model.AlertRule{}, errors.New("alert rule not found")
	}
	return normalizeAlertRule(out), nil
}

func (r *Repository) ListAlertEvents(ctx context.Context, f model.QueryFilter, severity string) (model.PageResult[model.AlertEvent], error) {
	severity = strings.ToLower(strings.TrimSpace(severity))
	where := []string{"detected_at >= ?", "detected_at <= ?"}
	args := []any{f.From, f.To}
	if severity != "" {
		where = append(where, "severity = ?")
		args = append(args, severity)
	}
	whereSQL := strings.Join(where, " AND ")
	var total uint64
	if err := r.db.QueryRowContext(ctx, `SELECT count() FROM alert_events WHERE `+whereSQL, args...).Scan(&total); err != nil {
		return model.PageResult[model.AlertEvent]{}, err
	}
	rows, err := r.db.QueryContext(ctx, `
SELECT event_id, event_key, rule_id, rule_name, rule_type, severity, detected_at, window_from, window_to, node_id, edge_id, description, bytes, flows, metadata_json
FROM alert_events
WHERE `+whereSQL+`
ORDER BY detected_at DESC
LIMIT ? OFFSET ?`, append(args, f.PageSize+1, (f.Page-1)*f.PageSize)...)
	if err != nil {
		return model.PageResult[model.AlertEvent]{}, err
	}
	defer rows.Close()
	items := make([]model.AlertEvent, 0, f.PageSize)
	for rows.Next() {
		var it model.AlertEvent
		var rawMeta string
		if err := rows.Scan(
			&it.EventID,
			&it.EventKey,
			&it.RuleID,
			&it.RuleName,
			&it.RuleType,
			&it.Severity,
			&it.DetectedAt,
			&it.WindowFrom,
			&it.WindowTo,
			&it.NodeID,
			&it.EdgeID,
			&it.Description,
			&it.Bytes,
			&it.Flows,
			&rawMeta,
		); err != nil {
			return model.PageResult[model.AlertEvent]{}, err
		}
		if rawMeta != "" {
			meta := map[string]any{}
			_ = json.Unmarshal([]byte(rawMeta), &meta)
			it.Metadata = meta
		}
		items = append(items, it)
	}
	hasNext := len(items) > f.PageSize
	if hasNext {
		items = items[:f.PageSize]
	}
	return model.PageResult[model.AlertEvent]{
		Data:     items,
		Page:     f.Page,
		PageSize: f.PageSize,
		Total:    total,
		HasNext:  hasNext,
		SortBy:   "detected_at",
		SortDir:  "DESC",
	}, nil
}

func (r *Repository) StoreAlertEvents(ctx context.Context, events []model.AlertEvent) (int, error) {
	if len(events) == 0 {
		return 0, nil
	}
	stmt, err := r.db.PrepareContext(ctx, `
INSERT INTO alert_events (
	event_id, event_key, rule_id, rule_name, rule_type, severity, detected_at,
	window_from, window_to, node_id, edge_id, description, bytes, flows, metadata_json
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return 0, err
	}
	defer stmt.Close()
	inserted := 0
	for _, event := range events {
		if strings.TrimSpace(event.EventID) == "" {
			event.EventID = uuid.NewString()
		}
		if strings.TrimSpace(event.EventKey) == "" {
			event.EventKey = fmt.Sprintf("%s|%s|%s|%s", event.RuleID, event.EdgeID, event.NodeID, event.WindowTo.UTC().Format("200601021504"))
		}
		if event.DetectedAt.IsZero() {
			event.DetectedAt = time.Now().UTC()
		}
		var exists uint64
		if err := r.db.QueryRowContext(ctx, `SELECT count() FROM alert_events WHERE event_key = ? AND detected_at >= ?`, event.EventKey, event.WindowFrom).Scan(&exists); err != nil {
			return inserted, err
		}
		if exists > 0 {
			continue
		}
		metaRaw, _ := json.Marshal(event.Metadata)
		if _, err := stmt.ExecContext(ctx,
			event.EventID,
			event.EventKey,
			event.RuleID,
			event.RuleName,
			event.RuleType,
			event.Severity,
			event.DetectedAt,
			event.WindowFrom,
			event.WindowTo,
			event.NodeID,
			event.EdgeID,
			event.Description,
			event.Bytes,
			event.Flows,
			string(metaRaw),
		); err != nil {
			return inserted, err
		}
		inserted++
	}
	return inserted, nil
}

func (r *Repository) EvaluateAlerts(ctx context.Context, refTs time.Time) ([]model.AlertEvent, error) {
	rules, err := r.ListAlertRules(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]model.AlertEvent, 0, 64)
	for _, rule := range rules {
		if !rule.Enabled {
			continue
		}
		ruleEvents, err := r.evaluateRule(ctx, normalizeAlertRule(rule), refTs)
		if err != nil {
			return nil, err
		}
		out = append(out, ruleEvents...)
	}
	return out, nil
}

func (r *Repository) evaluateRule(ctx context.Context, rule model.AlertRule, refTs time.Time) ([]model.AlertEvent, error) {
	switch strings.ToLower(strings.TrimSpace(rule.RuleType)) {
	case "new_edge":
		return r.evaluateNewEdgeRule(ctx, rule, refTs)
	case "fanout_external":
		return r.evaluateFanoutExternalRule(ctx, rule, refTs)
	case "high_byte_edge":
		return r.evaluateHighByteEdgeRule(ctx, rule, refTs)
	case "port_outlier":
		return r.evaluatePortOutlierRule(ctx, rule, refTs)
	default:
		return nil, nil
	}
}

func (r *Repository) evaluateNewEdgeRule(ctx context.Context, rule model.AlertRule, refTs time.Time) ([]model.AlertEvent, error) {
	from := refTs.Add(-time.Duration(rule.WindowMinutes) * time.Minute).Truncate(time.Minute)
	to := refTs.Truncate(time.Minute)
	rows, err := r.db.QueryContext(ctx, `
SELECT concat(src_node_id, '->', dst_node_id) AS edge_id, sum(bytes) AS b, sum(flows) AS f, min(first_seen) AS first_seen
FROM edges_1m
WHERE minute_bucket >= ? AND minute_bucket <= ?
GROUP BY edge_id
HAVING first_seen >= ? AND f >= ?
ORDER BY b DESC
LIMIT 40`, from, to, from, rule.ThresholdValue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.AlertEvent, 0, 16)
	for rows.Next() {
		var edgeID string
		var b, f uint64
		var firstSeen time.Time
		if err := rows.Scan(&edgeID, &b, &f, &firstSeen); err != nil {
			return nil, err
		}
		out = append(out, buildAlertEvent(rule, from, to, model.AlertEvent{
			EdgeID:      edgeID,
			Description: fmt.Sprintf("New edge observed: %s", edgeID),
			Bytes:       b,
			Flows:       f,
			Metadata: map[string]any{
				"first_seen": firstSeen.Format(time.RFC3339),
			},
		}))
	}
	return out, nil
}

func (r *Repository) evaluateFanoutExternalRule(ctx context.Context, rule model.AlertRule, refTs time.Time) ([]model.AlertEvent, error) {
	from := refTs.Add(-time.Duration(rule.WindowMinutes) * time.Minute).Truncate(time.Minute)
	to := refTs.Truncate(time.Minute)
	rows, err := r.db.QueryContext(ctx, `
SELECT src_node_id, uniqExact(dst_node_id) AS fanout, sum(bytes) AS b, sum(flows) AS f
FROM edges_1m
WHERE minute_bucket >= ? AND minute_bucket <= ? AND src_private = 1 AND dst_private = 0
GROUP BY src_node_id
HAVING fanout >= ?
ORDER BY fanout DESC
LIMIT 40`, from, to, rule.ThresholdValue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.AlertEvent, 0, 16)
	for rows.Next() {
		var nodeID string
		var fanout, b, f uint64
		if err := rows.Scan(&nodeID, &fanout, &b, &f); err != nil {
			return nil, err
		}
		out = append(out, buildAlertEvent(rule, from, to, model.AlertEvent{
			NodeID:      nodeID,
			Description: fmt.Sprintf("Node %s contacted %d external destinations", nodeID, fanout),
			Bytes:       b,
			Flows:       f,
			Metadata: map[string]any{
				"fanout": fanout,
			},
		}))
	}
	return out, nil
}

func (r *Repository) evaluateHighByteEdgeRule(ctx context.Context, rule model.AlertRule, refTs time.Time) ([]model.AlertEvent, error) {
	from := refTs.Add(-time.Duration(rule.WindowMinutes) * time.Minute).Truncate(time.Minute)
	to := refTs.Truncate(time.Minute)
	rows, err := r.db.QueryContext(ctx, `
SELECT concat(src_node_id, '->', dst_node_id) AS edge_id, sum(bytes) AS b, sum(flows) AS f
FROM edges_1m
WHERE minute_bucket >= ? AND minute_bucket <= ?
GROUP BY edge_id
HAVING b >= ?
ORDER BY b DESC
LIMIT 40`, from, to, rule.ThresholdValue)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.AlertEvent, 0, 16)
	for rows.Next() {
		var edgeID string
		var b, f uint64
		if err := rows.Scan(&edgeID, &b, &f); err != nil {
			return nil, err
		}
		out = append(out, buildAlertEvent(rule, from, to, model.AlertEvent{
			EdgeID:      edgeID,
			Description: fmt.Sprintf("High-byte edge detected: %s (%d bytes)", edgeID, b),
			Bytes:       b,
			Flows:       f,
		}))
	}
	return out, nil
}

func (r *Repository) evaluatePortOutlierRule(ctx context.Context, rule model.AlertRule, refTs time.Time) ([]model.AlertEvent, error) {
	from := refTs.Add(-time.Duration(rule.WindowMinutes) * time.Minute).Truncate(time.Minute)
	to := refTs.Truncate(time.Minute)
	rows, err := r.db.QueryContext(ctx, `
SELECT src_node_id, dst_port, sum(bytes) AS b, sum(flows) AS f
FROM edges_1m
WHERE minute_bucket >= ? AND minute_bucket <= ?
GROUP BY src_node_id, dst_port
HAVING f >= ?
ORDER BY f DESC
LIMIT 200`, from, to, maxUint64(rule.ThresholdValue, 1))
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	out := make([]model.AlertEvent, 0, 16)
	baselineFrom := from.Add(-7 * 24 * time.Hour).Truncate(time.Hour)
	baselineTo := from.Truncate(time.Hour)
	for rows.Next() {
		var nodeID string
		var dstPort uint16
		var b, f uint64
		if err := rows.Scan(&nodeID, &dstPort, &b, &f); err != nil {
			return nil, err
		}
		var baseline sql.NullInt64
		if err := r.db.QueryRowContext(ctx, `
SELECT sum(flows)
FROM edges_1h
WHERE hour_bucket >= ? AND hour_bucket < ? AND src_node_id = ? AND dst_port = ?`,
			baselineFrom, baselineTo, nodeID, dstPort).Scan(&baseline); err != nil {
			return nil, err
		}
		if baseline.Valid && baseline.Int64 > 0 {
			continue
		}
		out = append(out, buildAlertEvent(rule, from, to, model.AlertEvent{
			NodeID:      nodeID,
			Description: fmt.Sprintf("Node %s used rare destination port %d", nodeID, dstPort),
			Bytes:       b,
			Flows:       f,
			Metadata: map[string]any{
				"dst_port": dstPort,
			},
		}))
		if len(out) >= 40 {
			break
		}
	}
	return out, nil
}

func normalizeAlertRule(rule model.AlertRule) model.AlertRule {
	rule.RuleType = strings.ToLower(strings.TrimSpace(rule.RuleType))
	rule.Severity = strings.ToLower(strings.TrimSpace(rule.Severity))
	if rule.Severity == "" {
		rule.Severity = "medium"
	}
	if rule.WindowMinutes == 0 {
		rule.WindowMinutes = 15
	}
	if rule.ThresholdValue == 0 {
		switch rule.RuleType {
		case "fanout_external":
			rule.ThresholdValue = 15
		case "high_byte_edge":
			rule.ThresholdValue = 50000000
		case "new_edge":
			rule.ThresholdValue = 1
		case "port_outlier":
			rule.ThresholdValue = 1
		default:
			rule.ThresholdValue = 1
		}
	}
	return rule
}

func buildAlertEvent(rule model.AlertRule, from, to time.Time, event model.AlertEvent) model.AlertEvent {
	event.RuleID = rule.RuleID
	event.RuleName = rule.Name
	event.RuleType = rule.RuleType
	event.Severity = rule.Severity
	event.WindowFrom = from
	event.WindowTo = to
	event.DetectedAt = to
	event.EventKey = fmt.Sprintf("%s|%s|%s|%s|%s", rule.RuleID, event.EdgeID, event.NodeID, to.UTC().Format("200601021504"), event.Description)
	return event
}

func maxUint64(a uint64, b uint64) uint64 {
	if a > b {
		return a
	}
	return b
}

func (r *Repository) EnsureDefaultAlertRules(ctx context.Context, createdBy string) error {
	items, err := r.ListAlertRules(ctx)
	if err != nil {
		return err
	}
	if len(items) > 0 {
		return nil
	}
	defaults := []model.AlertRule{
		{Name: "New Edge Observed", RuleType: "new_edge", Enabled: true, WindowMinutes: 30, ThresholdValue: 1, Severity: "low", CreatedBy: createdBy},
		{Name: "Internal Fanout External", RuleType: "fanout_external", Enabled: true, WindowMinutes: 15, ThresholdValue: 15, Severity: "high", CreatedBy: createdBy},
		{Name: "High Byte Edge", RuleType: "high_byte_edge", Enabled: true, WindowMinutes: 15, ThresholdValue: 50000000, Severity: "medium", CreatedBy: createdBy},
		{Name: "Port Outlier", RuleType: "port_outlier", Enabled: true, WindowMinutes: 30, ThresholdValue: 1, Severity: "low", CreatedBy: createdBy},
	}
	for _, rule := range defaults {
		if _, err := r.CreateAlertRule(ctx, rule); err != nil {
			return err
		}
	}
	return nil
}
