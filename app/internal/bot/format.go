package bot

import (
	"fmt"
	"strings"
	"time"
)

func formatDate(date string) string {
	t, err := time.Parse("2006-01-02", date)
	if err != nil {
		return date
	}
	return t.Format("02.01.2006")
}

func formatPlansList(plans []Plan, title, subtitle string) string {
	if len(plans) == 0 {
		return fmt.Sprintf("📋 *%s*\n\nНет планов.", title)
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("📋 *%s*\n", title))

	if subtitle != "" {
		sb.WriteString(fmt.Sprintf("_%s_\n", subtitle))
	}

	sb.WriteString("\n")

	for i, plan := range plans {
		timeStr := "Весь день"
		if !plan.IsAllDay && plan.Time != "" {
			timeStr = "🕐 " + plan.Time
		}
		sb.WriteString(fmt.Sprintf("*%d.* %s\n", i+1, plan.Title))
		sb.WriteString(fmt.Sprintf("   _%s_\n", plan.Description))
		sb.WriteString(fmt.Sprintf("   %s\n\n", timeStr))
	}

	return sb.String()
}
