package bot

import (
	"testing"
)

func TestConversationStoreGetNew(t *testing.T) {
	store := newConversationStore()
	c := store.get(123)
	if c.step != stepIdle {
		t.Errorf("new conversation should be idle, got %d", c.step)
	}
}

func TestConversationStoreSetAndGet(t *testing.T) {
	store := newConversationStore()
	c := &conversation{step: stepTitle, plan: PlanRequest{Title: "test"}}
	store.set(123, c)

	got := store.get(123)
	if got.step != stepTitle {
		t.Errorf("step = %d, want %d", got.step, stepTitle)
	}
	if got.plan.Title != "test" {
		t.Errorf("title = %q, want 'test'", got.plan.Title)
	}
}

func TestConversationStoreReset(t *testing.T) {
	store := newConversationStore()
	store.set(123, &conversation{step: stepTitle})
	store.reset(123)

	c := store.get(123)
	if c.step != stepIdle {
		t.Errorf("after reset step = %d, want stepIdle", c.step)
	}
}

func TestConversationStoreIsolation(t *testing.T) {
	store := newConversationStore()
	store.set(1, &conversation{step: stepTitle, plan: PlanRequest{Title: "one"}})
	store.set(2, &conversation{step: stepDesc, plan: PlanRequest{Title: "two"}})

	if store.get(1).plan.Title != "one" {
		t.Error("chat 1 corrupted")
	}
	if store.get(2).plan.Title != "two" {
		t.Error("chat 2 corrupted")
	}

	store.reset(1)
	if store.get(2).step != stepDesc {
		t.Error("resetting chat 1 affected chat 2")
	}
}
