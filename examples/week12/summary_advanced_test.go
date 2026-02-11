// 详细注释: package week12
package week12

// 详细注释: import "testing"
import "testing"

// TestPortfolioScore 验证完成率计算。
// 详细注释: func TestPortfolioScore(t *testing.T) {
func TestPortfolioScore(t *testing.T) {
	// 详细注释: p := NewPortfolio()
	p := NewPortfolio()
	// 详细注释: p.MarkCompleted("http", "week02 tests")
	p.MarkCompleted("http", "week02 tests")
	// 详细注释: p.MarkCompleted("db", "week05 migration")
	p.MarkCompleted("db", "week05 migration")
	// 详细注释: if got := p.Score(4); got != 50 {
	if got := p.Score(4); got != 50 {
		// 详细注释: t.Fatalf("want 50 got %d", got)
		t.Fatalf("want 50 got %d", got)
		// 详细注释: }
	}
	// 详细注释: }
}

// TestNextStagePlan 验证不同主题的计划输出。
// 详细注释: func TestNextStagePlan(t *testing.T) {
func TestNextStagePlan(t *testing.T) {
	// 详细注释: plan := NextStagePlan("mq")
	plan := NextStagePlan("mq")
	// 详细注释: if len(plan) != 4 {
	if len(plan) != 4 {
		// 详细注释: t.Fatalf("want 4 steps got %d", len(plan))
		t.Fatalf("want 4 steps got %d", len(plan))
		// 详细注释: }
	}
	// 详细注释: if plan[0] != "week1: queue basics" {
	if plan[0] != "week1: queue basics" {
		// 详细注释: t.Fatalf("unexpected plan: %+v", plan)
		t.Fatalf("unexpected plan: %+v", plan)
		// 详细注释: }
	}
	// 详细注释: }
}
