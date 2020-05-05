package server

import (
	"reflect"
	"testing"

	"github.com/rs/zerolog"
)

func TestNewLogger(t *testing.T) {
	type args struct {
		develop bool
		level   zerolog.Level
	}
	tests := []struct {
		name string
		args args
		want *zerolog.Logger
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewLogger(tt.args.develop, tt.args.level); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewLogger() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewClipdServer(t *testing.T) {
	type args struct {
		config *Config
	}
	tests := []struct {
		name string
		args args
		want *ClipdServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewClipdServer(tt.args.config); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewClipdServer() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClipdServer_CleanAfterShutdown(t *testing.T) {
	tests := []struct {
		name string
		c    *ClipdServer
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.c.CleanAfterShutdown()
		})
	}
}

func TestClipdServer_ListenAndServe(t *testing.T) {
	tests := []struct {
		name    string
		c       *ClipdServer
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := tt.c.ListenAndServe(); (err != nil) != tt.wantErr {
				t.Errorf("ClipdServer.ListenAndServe() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
