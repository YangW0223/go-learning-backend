// 详细注释: package week12
package week12

// 详细注释: import (
import (
	// 详细注释: "fmt"
	"fmt"
	// 详细注释: "sort"
	"sort"
	// 详细注释: "strings"
	"strings"
	// 详细注释: )
)

// Portfolio 用于汇总 12 周能力产出。
// 详细注释: type Portfolio struct {
type Portfolio struct {
	// 详细注释: completed map[string]bool
	completed map[string]bool
	// 详细注释: evidence  map[string]string
	evidence map[string]string
	// 详细注释: }
}

// NewPortfolio 创建作品集。
// 详细注释: func NewPortfolio() *Portfolio {
func NewPortfolio() *Portfolio {
	// 详细注释: return &Portfolio{
	return &Portfolio{
		// 详细注释: completed: make(map[string]bool),
		completed: make(map[string]bool),
		// 详细注释: evidence:  make(map[string]string),
		evidence: make(map[string]string),
		// 详细注释: }
	}
	// 详细注释: }
}

// MarkCompleted 标记能力项已完成。
// 详细注释: func (p *Portfolio) MarkCompleted(capability, proof string) {
func (p *Portfolio) MarkCompleted(capability, proof string) {
	// 详细注释: capability = strings.TrimSpace(capability)
	capability = strings.TrimSpace(capability)
	// 详细注释: proof = strings.TrimSpace(proof)
	proof = strings.TrimSpace(proof)
	// 详细注释: if capability == "" {
	if capability == "" {
		// 详细注释: return
		return
		// 详细注释: }
	}
	// 详细注释: p.completed[capability] = true
	p.completed[capability] = true
	// 详细注释: p.evidence[capability] = proof
	p.evidence[capability] = proof
	// 详细注释: }
}

// Score 返回完成率百分比。
// 详细注释: func (p *Portfolio) Score(total int) int {
func (p *Portfolio) Score(total int) int {
	// 详细注释: if total <= 0 {
	if total <= 0 {
		// 详细注释: return 0
		return 0
		// 详细注释: }
	}
	// 详细注释: return len(p.completed) * 100 / total
	return len(p.completed) * 100 / total
	// 详细注释: }
}

// Capabilities 返回按字典序排序的已完成能力列表。
// 详细注释: func (p *Portfolio) Capabilities() []string {
func (p *Portfolio) Capabilities() []string {
	// 详细注释: keys := make([]string, 0, len(p.completed))
	keys := make([]string, 0, len(p.completed))
	// 详细注释: for k := range p.completed {
	for k := range p.completed {
		// 详细注释: keys = append(keys, k)
		keys = append(keys, k)
		// 详细注释: }
	}
	// 详细注释: sort.Strings(keys)
	sort.Strings(keys)
	// 详细注释: return keys
	return keys
	// 详细注释: }
}

// NextStagePlan 给出下一阶段 4 周最小计划。
// 详细注释: func NextStagePlan(topic string) []string {
func NextStagePlan(topic string) []string {
	// 详细注释: topic = strings.ToLower(strings.TrimSpace(topic))
	topic = strings.ToLower(strings.TrimSpace(topic))
	// 详细注释: switch topic {
	switch topic {
	// 详细注释: case "microservice", "microservices":
	case "microservice", "microservices":
		// 详细注释: return []string{"week1: service split", "week2: service discovery", "week3: gateway", "week4: distributed tracing"}
		return []string{"week1: service split", "week2: service discovery", "week3: gateway", "week4: distributed tracing"}
		// 详细注释: case "mq", "message-queue":
	case "mq", "message-queue":
		// 详细注释: return []string{"week1: queue basics", "week2: retry and dead-letter", "week3: outbox", "week4: idempotent consumer"}
		return []string{"week1: queue basics", "week2: retry and dead-letter", "week3: outbox", "week4: idempotent consumer"}
		// 详细注释: case "ddd":
	case "ddd":
		// 详细注释: return []string{"week1: bounded context", "week2: aggregate", "week3: domain service", "week4: application service"}
		return []string{"week1: bounded context", "week2: aggregate", "week3: domain service", "week4: application service"}
		// 详细注释: default:
	default:
		// 详细注释: return []string{"week1: choose one topic", "week2: build poc", "week3: add tests", "week4: retrospective"}
		return []string{"week1: choose one topic", "week2: build poc", "week3: add tests", "week4: retrospective"}
		// 详细注释: }
	}
	// 详细注释: }
}

// BuildNarrative 生成可演示的项目叙述。
// 详细注释: func BuildNarrative(projectName string, capabilities []string) string {
func BuildNarrative(projectName string, capabilities []string) string {
	// 详细注释: return fmt.Sprintf("project=%s capabilities=%s", projectName, strings.Join(capabilities, ","))
	return fmt.Sprintf("project=%s capabilities=%s", projectName, strings.Join(capabilities, ","))
	// 详细注释: }
}
