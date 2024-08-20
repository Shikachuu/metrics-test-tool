package main

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"testing"
)

func TestSetLogger(t *testing.T) {
	testCases := []struct {
		desc    string
		logType string
		w       *bytes.Buffer
		isJSON  bool
	}{
		{
			desc:    "text logger",
			w:       bytes.NewBuffer(nil),
			logType: "text",
			isJSON:  false,
		},
		{
			desc:    "json logger",
			w:       bytes.NewBuffer(nil),
			logType: "json",
			isJSON:  true,
		},
	}
	for _, tC := range testCases {
		tC := tC
		t.Run(tC.desc, func(t *testing.T) {
			// cannot use parallel (t.Parallel()) because we are testing the global logger too.

			l := setLogger(tC.logType, tC.w)

			l.Info("info")
			l.Error("error")
			l.Debug("debug")

			if tC.w.Len() == 0 {
				t.Errorf("expected log output, got empty")
			}

			if !bytes.Contains(tC.w.Bytes(), []byte("info")) {
				t.Errorf("expected info log, got %s", tC.w.String())
			}

			if !bytes.Contains(tC.w.Bytes(), []byte("error")) {
				t.Errorf("expected error log, got %s", tC.w.String())
			}

			if bytes.Contains(tC.w.Bytes(), []byte("debug")) {
				t.Errorf("expected no debug log, got %s", tC.w.String())
			}

			tC.w.Reset()
			l.Info("info")

			jv := json.Valid(tC.w.Bytes())
			if tC.isJSON && !jv {
				t.Errorf("expected json log, got text")
			} else if !tC.isJSON && jv {
				t.Errorf("expected text log, got json")
			}

			tC.w.Reset()
			slog.Info("info")
			if tC.w.Len() == 0 {
				t.Errorf("expected log output from default logger, got empty")
			}
		})
	}
}
