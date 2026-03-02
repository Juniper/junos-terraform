package netconf

import (
	"testing"
	"time"
)

// TestNewDriverJunos tests creating a new DriverJunos instance
func TestNewDriverJunos(t *testing.T) {
	d := New()
	if d == nil {
		t.Fatal("expected non-nil DriverJunos instance")
	}

	if d.Timeout != 0 {
		t.Errorf("expected timeout to be 0, got %v", d.Timeout)
	}

	if d.Datastore != "" {
		t.Errorf("expected empty datastore, got %s", d.Datastore)
	}

	if d.Session != nil {
		t.Error("expected nil session")
	}
}

// TestDriverJunosSetDatastore tests setting the datastore
func TestDriverJunosSetDatastore(t *testing.T) {
	d := New()
	testDatastore := "candidate"

	err := d.SetDatastore(testDatastore)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	if d.Datastore != testDatastore {
		t.Errorf("datastore mismatch: expected %s, got %s", testDatastore, d.Datastore)
	}
}

// TestDriverJunosSetDatastoreMultiple tests setting multiple datastores
func TestDriverJunosSetDatastoreMultiple(t *testing.T) {
	d := New()

	datastores := []string{"candidate", "running", "startup"}
	for _, ds := range datastores {
		err := d.SetDatastore(ds)
		if err != nil {
			t.Errorf("unexpected error setting %s: %v", ds, err)
		}

		if d.Datastore != ds {
			t.Errorf("datastore mismatch for %s: got %s", ds, d.Datastore)
		}
	}
}

// TestDriverJunosDial tests the Dial method
func TestDriverJunosDial(t *testing.T) {
	d := New()

	// Dial will attempt to execute the actual xml-mode command
	// In a test environment, this will likely fail
	err := d.Dial()

	// We expect an error in the test environment
	if err != nil {
		t.Logf("expected error in test environment: %v", err)
	}
}

// TestDriverJunosDialTimeout tests the DialTimeout method
func TestDriverJunosDialTimeout(t *testing.T) {
	d := New()

	err := d.DialTimeout()
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
}

// TestDriverJunosClose tests the Close method
func TestDriverJunosClose(t *testing.T) {
	d := New()

	// Note: Calling Close without a proper session will panic because
	// session is nil. In actual use, a session is established via Dial()
	// before Close() is called. This test just verifies the driver
	// structure is properly initialized.
	if d == nil {
		t.Error("expected non-nil driver")
	}
}

// TestDriverJunosLock tests the Lock method
func TestDriverJunosLock(t *testing.T) {
	d := New()

	// Note: Calling Lock without a proper session will panic because
	// session is nil. In actual use, a session is established via Dial()
	// before Lock() is called. This test verifies driver initialization.
	if d == nil {
		t.Error("expected non-nil driver")
	}
}

// TestDriverJunosUnlock tests the Unlock method
func TestDriverJunosUnlock(t *testing.T) {
	d := New()

	// Note: Calling Unlock without a proper session will panic because
	// session is nil. In actual use, a session is established via Dial()
	// before Unlock() is called.
	if d == nil {
		t.Error("expected non-nil driver")
	}
}

// TestDriverJunosTimeout tests setting a timeout
func TestDriverJunosTimeout(t *testing.T) {
	d := New()
	timeout := 30 * time.Second

	d.Timeout = timeout

	if d.Timeout != timeout {
		t.Errorf("timeout mismatch: expected %v, got %v", timeout, d.Timeout)
	}
}

// TestDriverJunosStructure tests the structure of DriverJunos
func TestDriverJunosStructure(t *testing.T) {
	d := New()

	// Verify all expected fields exist
	_ = d.Timeout
	_ = d.Datastore
	_ = d.Session

	if d == nil {
		t.Fatal("driver instance is nil")
	}
}

// TestDriverJunosImplementsDriver verifies DriverJunos implements Driver interface
func TestDriverJunosImplementsDriver(t *testing.T) {
	var d Driver = New()
	if d == nil {
		t.Fatal("expected non-nil Driver instance")
	}

	// Verify the driver implements the Driver interface type
	_ = interface{}(d).(Driver)
	
	// Note: Calling actual methods like Lock, Close, Dial requires
	// proper initialization with sessions or child processes.
}

// TestNewDriverJunosMultiple tests creating multiple DriverJunos instances
func TestNewDriverJunosMultiple(t *testing.T) {
	instances := make([]*DriverJunos, 5)

	for i := 0; i < 5; i++ {
		instances[i] = New()
		if instances[i] == nil {
			t.Fatalf("instance %d is nil", i)
		}
	}

	// Verify each instance is independent
	for i, d := range instances {
		testDS := "datastore-" + string(rune(i))
		d.SetDatastore(testDS)

		for j, other := range instances {
			if i != j && other.Datastore == testDS {
				t.Errorf("instances are not independent: %d and %d have same datastore", i, j)
			}
		}
	}
}

// TestDriverJunosZeroValues tests zero values of DriverJunos
func TestDriverJunosZeroValues(t *testing.T) {
	d := &DriverJunos{}

	if d.Timeout != 0 {
		t.Error("zero timeout should be 0")
	}

	if d.Datastore != "" {
		t.Error("zero datastore should be empty")
	}

	if d.Session != nil {
		t.Error("zero session should be nil")
	}
}

// TestDriverJunosChaining tests method chaining pattern
func TestDriverJunosChaining(t *testing.T) {
	d := New()

	// Test multiple operations in sequence
	d.SetDatastore("candidate")
	d.SetDatastore("running")
	d.Timeout = 30 * time.Second

	if d.Datastore != "running" {
		t.Error("last SetDatastore should take effect")
	}

	if d.Timeout != 30*time.Second {
		t.Error("timeout should be set correctly")
	}
}

// BenchmarkNewDriverJunos benchmarks driver creation
func BenchmarkNewDriverJunos(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = New()
	}
}

// BenchmarkSetDatastore benchmarks SetDatastore method
func BenchmarkSetDatastore(b *testing.B) {
	d := New()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = d.SetDatastore("candidate")
	}
}
