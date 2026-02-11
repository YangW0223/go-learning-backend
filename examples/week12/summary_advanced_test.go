package week12

import "testing"

// TestPortfolioScore 验证完成率计算。
func TestPortfolioScore(t *testing.T) {
	p := NewPortfolio()
	p.MarkCompleted("http", "week02 tests")
	p.MarkCompleted("db", "week05 migration")
	if got := p.Score(4); got != 50 {
		t.Fatalf("want 50 got %d", got)
	}
}

// TestNextStagePlan 验证不同主题的计划输出。
func TestNextStagePlan(t *testing.T) {
	plan := NextStagePlan("mq")
	if len(plan) != 4 {
		t.Fatalf("want 4 steps got %d", len(plan))
	}
	if plan[0] != "week1: queue basics" {
		t.Fatalf("unexpected plan: %+v", plan)
	}
}
