package week12

import (
	"fmt"
	"sort"
	"strings"
)

// Portfolio 用于汇总 12 周能力产出。
type Portfolio struct {
	completed map[string]bool
	evidence  map[string]string
}

// NewPortfolio 创建作品集。
func NewPortfolio() *Portfolio {
	return &Portfolio{
		completed: make(map[string]bool),
		evidence:  make(map[string]string),
	}
}

// MarkCompleted 标记能力项已完成。
func (p *Portfolio) MarkCompleted(capability, proof string) {
	capability = strings.TrimSpace(capability)
	proof = strings.TrimSpace(proof)
	if capability == "" {
		return
	}
	p.completed[capability] = true
	p.evidence[capability] = proof
}

// Score 返回完成率百分比。
func (p *Portfolio) Score(total int) int {
	if total <= 0 {
		return 0
	}
	return len(p.completed) * 100 / total
}

// Capabilities 返回按字典序排序的已完成能力列表。
func (p *Portfolio) Capabilities() []string {
	keys := make([]string, 0, len(p.completed))
	for k := range p.completed {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// NextStagePlan 给出下一阶段 4 周最小计划。
func NextStagePlan(topic string) []string {
	topic = strings.ToLower(strings.TrimSpace(topic))
	switch topic {
	case "microservice", "microservices":
		return []string{"week1: service split", "week2: service discovery", "week3: gateway", "week4: distributed tracing"}
	case "mq", "message-queue":
		return []string{"week1: queue basics", "week2: retry and dead-letter", "week3: outbox", "week4: idempotent consumer"}
	case "ddd":
		return []string{"week1: bounded context", "week2: aggregate", "week3: domain service", "week4: application service"}
	default:
		return []string{"week1: choose one topic", "week2: build poc", "week3: add tests", "week4: retrospective"}
	}
}

// BuildNarrative 生成可演示的项目叙述。
func BuildNarrative(projectName string, capabilities []string) string {
	return fmt.Sprintf("project=%s capabilities=%s", projectName, strings.Join(capabilities, ","))
}
