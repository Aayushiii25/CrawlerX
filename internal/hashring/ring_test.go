package hashring

import (
	"fmt"
	"math"
	"testing"
)

func TestHashRing_BasicAssignment(t *testing.T) {
	ring := New(150)

	ring.AddNode("node-1")
	ring.AddNode("node-2")
	ring.AddNode("node-3")

	domains := []string{"example.com", "google.com", "github.com", "reddit.com"}
	for _, domain := range domains {
		node := ring.GetNode(domain)
		if node == "" {
			t.Errorf("Domain %s mapped to empty node", domain)
		}
		t.Logf("Domain %s → %s", domain, node)
	}
}

func TestHashRing_EmptyRing(t *testing.T) {
	ring := New(150)

	node := ring.GetNode("example.com")
	if node != "" {
		t.Errorf("Expected empty string for empty ring, got %s", node)
	}
}

func TestHashRing_Consistency(t *testing.T) {
	ring := New(150)

	ring.AddNode("node-1")
	ring.AddNode("node-2")
	ring.AddNode("node-3")

	first := ring.GetNode("example.com")
	for i := 0; i < 100; i++ {
		got := ring.GetNode("example.com")
		if got != first {
			t.Errorf("Inconsistent mapping: first=%s, got=%s at iteration %d", first, got, i)
		}
	}
}

func TestHashRing_MinimalRedistribution(t *testing.T) {
	ring := New(150)

	ring.AddNode("node-1")
	ring.AddNode("node-2")
	ring.AddNode("node-3")

	numDomains := 1000
	initialAssignments := make(map[string]string, numDomains)
	for i := 0; i < numDomains; i++ {
		domain := fmt.Sprintf("domain-%d.com", i)
		initialAssignments[domain] = ring.GetNode(domain)
	}

	ring.AddNode("node-4")

	moved := 0
	for domain, oldNode := range initialAssignments {
		newNode := ring.GetNode(domain)
		if newNode != oldNode {
			moved++
		}
	}

	expectedMovement := float64(numDomains) / 4.0
	maxAcceptable := expectedMovement * 2.0

	t.Logf("Keys moved: %d/%d (%.1f%%), expected ~%.0f (%.1f%%)",
		moved, numDomains, float64(moved)/float64(numDomains)*100,
		expectedMovement, expectedMovement/float64(numDomains)*100)

	if float64(moved) > maxAcceptable {
		t.Errorf("Too many keys moved: %d (max acceptable: %.0f)", moved, maxAcceptable)
	}
}

func TestHashRing_Distribution(t *testing.T) {
	ring := New(150)

	numNodes := 5
	for i := 1; i <= numNodes; i++ {
		ring.AddNode(fmt.Sprintf("node-%d", i))
	}

	numDomains := 10000
	counts := make(map[string]int)
	for i := 0; i < numDomains; i++ {
		domain := fmt.Sprintf("domain-%d.com", i)
		node := ring.GetNode(domain)
		counts[node]++
	}

	expected := float64(numDomains) / float64(numNodes)
	for node, count := range counts {
		deviation := math.Abs(float64(count)-expected) / expected * 100
		t.Logf("Node %s: %d keys (%.1f%% deviation from expected %.0f)",
			node, count, deviation, expected)

		if deviation > 30 {
			t.Errorf("Node %s has excessive deviation: %.1f%%", node, deviation)
		}
	}
}

func TestHashRing_AddRemoveNode(t *testing.T) {
	ring := New(150)

	ring.AddNode("node-1")
	ring.AddNode("node-2")

	if ring.Size() != 2 {
		t.Errorf("Expected size 2, got %d", ring.Size())
	}

	ring.RemoveNode("node-1")
	if ring.Size() != 1 {
		t.Errorf("Expected size 1, got %d", ring.Size())
	}

	node := ring.GetNode("example.com")
	if node != "node-2" {
		t.Errorf("Expected node-2 after removing node-1, got %s", node)
	}

	ring.AddNode("node-2")
	if ring.Size() != 1 {
		t.Errorf("Expected size 1 after duplicate add, got %d", ring.Size())
	}
}
