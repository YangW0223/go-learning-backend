package memory

import "testing"

func TestTodoStore_CreateAndList(t *testing.T) {
	s := NewTodoStore()

	created, err := s.Create("learn go")
	if err != nil {
		t.Fatalf("create should not fail: %v", err)
	}

	if created.Title != "learn go" {
		t.Fatalf("unexpected title: %s", created.Title)
	}

	items, err := s.List()
	if err != nil {
		t.Fatalf("list should not fail: %v", err)
	}

	if len(items) != 1 {
		t.Fatalf("expected 1 item, got %d", len(items))
	}
}

func TestTodoStore_MarkDone(t *testing.T) {
	s := NewTodoStore()

	created, err := s.Create("learn go")
	if err != nil {
		t.Fatalf("create should not fail: %v", err)
	}

	updated, err := s.MarkDone(created.ID)
	if err != nil {
		t.Fatalf("mark done should not fail: %v", err)
	}

	if !updated.Done {
		t.Fatalf("expected done=true")
	}
}
