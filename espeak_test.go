// Copyright 2020 djangulo. All rights reserved. Use of this source code is
// governed by an MIT license that can be found in the LICENSE file.
package espeak

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
)

func TestTextToSpeech(t *testing.T) {
	tmp, err := ioutil.TempDir("", "go-espeak-test-*")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmp)
	p := NewParameters(WithDir(tmp))
	t.Run("success", func(t *testing.T) {
		samples, err := TextToSpeech("test speech", nil, "test", p)
		if err != nil {
			t.Errorf("unexpected error: %v", err)
		}
		if samples == 0 {
			t.Errorf("0 samples written")
		}
	})
	t.Run("errors", func(t *testing.T) {
		for _, tt := range []struct {
			name   string
			text   string
			params *Parameters
			voice  *Voice
			want   error
		}{
			{"empty text", "", nil, nil, ErrEmptyText},
		} {
			t.Run(tt.name, func(t *testing.T) {
				s, err := TextToSpeech(tt.text, nil, "test", p)
				if s != 0 {
					t.Errorf("expected return samples 0 got %d", s)
				}
				if !errors.Is(err, tt.want) {

				}
				if err == nil {
					t.Error("expected an error but didn't get one")
				}
			})
		}
	})

}

func TestNewParametersDoesNotMutateDefaults(t *testing.T) {
	before := *DefaultParameters
	p := NewParameters(WithDir("/tmp/espeak-test-dir")).WithRate(123)
	if p.Dir != "/tmp/espeak-test-dir" || p.Rate != 123 {
		t.Errorf("options not applied: Dir=%q Rate=%d", p.Dir, p.Rate)
	}
	if p == DefaultParameters {
		t.Error("NewParameters returned the shared DefaultParameters pointer")
	}
	after := *DefaultParameters
	if before != after {
		t.Errorf("DefaultParameters was mutated by NewParameters: %+v -> %+v", before, after)
	}
}

func TestNewParametersConcurrent(t *testing.T) {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			_ = NewParameters().WithDir(fmt.Sprintf("/tmp/espeak-test-%d", i))
		}(i)
	}
	wg.Wait()
	// defaults must be untouched by the concurrent WithDir calls
	if DefaultParameters.Dir != os.TempDir() {
		t.Errorf("DefaultParameters.Dir changed: %q", DefaultParameters.Dir)
	}
}

func TestInitIdempotent(t *testing.T) {
	t.Run("second init is a no-op", func(t *testing.T) {
		id1, sr1, err := Init(Synchronous, 200, nil, PhonemeEvents)
		if err != nil {
			t.Fatalf("first Init failed: %v", err)
		}
		if sr1 == 0 {
			t.Fatal("first Init returned zero sample rate")
		}
		id2, sr2, err := Init(Playback, -1, nil, 0)
		if err != nil {
			t.Fatalf("second Init failed: %v", err)
		}
		if sr2 != sr1 {
			t.Errorf("sample rate changed after no-op Init: %d != %d", sr2, sr1)
		}
		if id1 == id2 {
			t.Error("expected a fresh data-block id from each Init")
		}
	})
	t.Run("synthesis still works after re-init", func(t *testing.T) {
		samples, err := GenSamples("hello", nil, NewParameters())
		if err != nil {
			t.Fatalf("GenSamples after re-Init failed: %v", err)
		}
		if len(samples) == 0 {
			t.Error("GenSamples returned no samples")
		}
	})
	t.Run("concurrent inits", func(t *testing.T) {
		var wg sync.WaitGroup
		for i := 0; i < 8; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				if _, _, err := Init(Synchronous, 200, nil, PhonemeEvents); err != nil {
					t.Errorf("concurrent Init failed: %v", err)
				}
			}()
		}
		wg.Wait()
	})
	t.Run("re-init after terminate", func(t *testing.T) {
		if err := Terminate(); err != nil {
			t.Fatalf("Terminate failed: %v", err)
		}
		if _, _, err := Init(Synchronous, 200, nil, PhonemeEvents); err != nil {
			t.Fatalf("Init after Terminate failed: %v", err)
		}
	})
}

func TestEnsureWavSuffix(t *testing.T) {
	for _, tt := range []struct {
		in, want string
	}{
		{"outfile.wav", "outfile.wav"},
		{"out", "out.wav"},
		{"out.mp4", "out.mp4.wav"},
		{"out.", "out.wav"},
		{"out....................", "out.wav"},
		{"out...____.", "out...____.wav"},
		{"out.", "out.wav"},
	} {
		t.Run(fmt.Sprintf("ensureWavSuffix(%q)==%q", tt.in, tt.want), func(t *testing.T) {
			got := ensureWavSuffix(tt.in)
			if got != tt.want {
				t.Errorf("expected %q got %q", tt.want, got)
			}
		})
	}
}
