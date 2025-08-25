package pm

import (
	"fmt"
	"testing"
)

func TestEncodeURL(t *testing.T) {
	data := Pm{
		"abc": "sdfsdf",
		"amt": 100.01,
		"ddd": "",
		"efg": 0,
	}
	fmt.Println(data.EncodeURL())
}

func TestPmBasic(t *testing.T) {
	pm := NewPm(map[string]any{
		"Name":   "Alice",
		"Age":    "25",
		"Height": 1.70,
		"Active": true,
	})

	// --- GetString ---
	if got := pm.GetString("Name"); got != "Alice" {
		t.Errorf("GetString(Name) = %v, want Alice", got)
	}
	if got := pm.GetStringStrict("Height"); got != "" {
		t.Errorf("GetStringStrict(Height) = %v, want ''", got)
	}

	// --- GetInt ---
	if got := pm.GetInt("Age"); got != 25 {
		t.Errorf("GetInt(Age) = %v, want 25", got)
	}
	if v, ok := pm.GetIntOk("Age"); !ok || v != 25 {
		t.Errorf("GetIntOk(Age) = (%v,%v), want (25,true)", v, ok)
	}

	// --- GetFloat64 ---
	if got := pm.GetFloat64("Height"); got != 1.70 {
		t.Errorf("GetFloat64(Height) = %v, want 1.70", got)
	}
	if v, ok := pm.GetFloat64Ok("Height"); !ok || v != 1.70 {
		t.Errorf("GetFloat64Ok(Height) = (%v,%v), want (1.70,true)", v, ok)
	}

	// --- GetInt64 ---
	if got := pm.GetInt64("Age"); got != 25 {
		t.Errorf("GetInt64(Age) = %v, want 25", got)
	}

	// --- GetUint32 ---
	pm.Set("Count", "123")
	if got := pm.GetUint32("Count"); got != 123 {
		t.Errorf("GetUint32(Count) = %v, want 123", got)
	}
	if v, ok := pm.GetUint32Ok("Count"); !ok || v != 123 {
		t.Errorf("GetUint32Ok(Count) = (%v,%v), want (123,true)", v, ok)
	}

	// --- SubPm ---
	child := pm.SubPm("Child")
	child.Set("X", 99)
	if v := pm.SubPm("Child").GetInt("X"); v != 99 {
		t.Errorf("SubPm(Child).GetInt(X) = %v, want 99", v)
	}

	// --- EncodeURL ---
	pm2 := NewPm(map[string]any{
		"b": "world",
		"a": "hello",
	})
	encoded := pm2.EncodeURL()
	if encoded != "a=hello&b=world" {
		t.Errorf("EncodeURL = %v, want 'a=hello&b=world'", encoded)
	}
}
