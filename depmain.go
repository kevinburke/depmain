// Package depmain encapsulates a program's external state.
//
// depmain is designed to make it easy to test your program's main function.
// Create a second main() function that accepts a *depmain.Ext and returns an
// integer. Call it from main like this:
//
//	os.Exit(_main(depmain.New()))
//
// Then in tests, you can replace the *depmain.Ext with a bytes.Buffer,
// configure the environment as you see fit, etc.

package depmain

import (
	"io"
	"os"
	"sync"
)

// Ext encapsulates a program's external environment.
type Ext struct {
	Stdin          io.Reader
	Stdout, Stderr io.Writer
	Env            []string
	Args           []string

	env     map[string]int
	envOnce sync.Once
	envMu   sync.RWMutex
}

func (e *Ext) copyenv() {
	e.env = make(map[string]int)
	for i, s := range e.Env {
		for j := 0; j < len(s); j++ {
			if s[j] == '=' {
				key := s[:j]
				if _, ok := e.env[key]; !ok {
					e.env[key] = i // first mention of key
				} else {
					// Clear duplicate keys. This permits Unsetenv to
					// safely delete only the first item without
					// worrying about unshadowing a later one,
					// which might be a security problem.
					e.Env[i] = ""
				}
				break
			}
		}
	}
}

// LookupEnv retrieves the value of the environment variable named
// by the key. If the variable is present in the environment the
// value (which may be empty) is returned and the boolean is true.
// Otherwise the returned value will be empty and the boolean will
// be false.
func (e *Ext) LookupEnv(key string) (string, bool) {
	e.envOnce.Do(e.copyenv)
	if len(key) == 0 {
		return "", false
	}

	e.envMu.RLock()
	defer e.envMu.RUnlock()

	i, ok := e.env[key]
	if !ok {
		return "", false
	}
	s := e.Env[i]
	for i := 0; i < len(s); i++ {
		if s[i] == '=' {
			return s[i+1:], true
		}
	}
	return "", false
}

// Getenv retrieves the value of the environment variable named by the key.
// It returns the value, which will be empty if the variable is not present.
// To distinguish between an empty value and an unset value, use ext.LookupEnv.
func (e *Ext) Getenv(key string) string {
	v, _ := e.LookupEnv(key)
	return v
}

// New creates an Ext with all arguments set to their default values.
func New() *Ext {
	return &Ext{
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
		Env:    os.Environ(),
		Args:   os.Args,
	}
}
