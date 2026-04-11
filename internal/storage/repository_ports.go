package storage

import (
	"context"

	"github.com/flowscope/flowscope/internal/model"
)

func (r *Repository) TopPorts(ctx context.Context, f model.QueryFilter, limit int) ([]model.TopItem, error) {
	if limit <= 0 || limit > 100 {
		limit = 10
	}
	where, args := buildFlowWhere(f, "")
	args = append(args, limit)
	rows, err := r.db.QueryContext(ctx, `
SELECT toString(dst_port) AS key, sum(bytes) AS bytes, sum(packets) AS packets, count() AS flows
FROM raw_flow_events
WHERE `+where+`
GROUP BY dst_port
ORDER BY bytes DESC
LIMIT ?`, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanTopItems(rows)
}
