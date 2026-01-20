package internal

import "testing"

func TestReadConfig(t *testing.T) {
	path := "../test/config.yaml"
	c, err := ReadConfig(path)
	if err != nil {
		t.Fatalf("error reading config: %v", err)
	}

	if len(c.Assets) != 3 {
		t.Fatalf("incorrect length for c.Assets, expected: 3, got: %v", len(c.Assets))
	}

	if c.Assets[0].Glob != "./*.css" {
		t.Fatalf("input.glob failed expectation")
	}

	if c.Assets[1].Id != nil {
		t.Fatalf("input.id expected to be nil")
	}

	if c.Assets[2].Id == nil || *c.Assets[2].Id != "text" {
		t.Fatalf("input.id failed expectation")
	}

	if c.Assets[2].Preload != true {
		t.Fatalf("input.preload failed expectation")
	}

	if c == nil || *c.Out != "./dist" {
		t.Fatalf("c.out failed expectation")
	}
}
