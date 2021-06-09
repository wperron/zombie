package config

import (
	"reflect"
	"testing"
)

func TestConfig(t *testing.T) {
	t.Run("simple target", func(t *testing.T) {
		conf, err := Load(simple)
		if err != nil {
			t.Errorf("failed to parse simple config: %s", err)
			t.FailNow()
		}

		if len(conf.Targets) != 1 {
			t.Errorf("expected only one target, got %d", len(conf.Targets))
		}

		target := conf.Targets[0]
		if target.Name != "" {
			t.Errorf("expected target to be unnamed, got name %s", target.Name)
		}

		if target.Delay != 0 {
			t.Errorf("expected target to have default delay, got %d", target.Delay)
		}

		if target.Jitter != 0 {
			t.Errorf("expected target to have default jitter, got %f", target.Jitter)
		}
	})

	t.Run("with headers", func(t *testing.T) {
		conf, err := Load(headers)
		if err != nil {
			t.Errorf("failed to parse simple config: %s", err)
			t.FailNow()
		}

		if len(conf.Targets) != 1 {
			t.Errorf("expected only one target, got %d", len(conf.Targets))
		}

		target := conf.Targets[0]
		if target.Name != "" {
			t.Errorf("expected target to be unnamed, got name %s", target.Name)
		}

		if target.Headers == nil || len(*target.Headers) == 0 {
			t.Errorf("no headers set on target")
		}

		if contentType := target.Headers.Values("Content-Type"); !reflect.DeepEqual(contentType, []string{"application/json"}) {
			t.Errorf("expected content type header `application/json`, got %s", contentType)
		}

		if contentType := target.Headers.Values("Accept"); !reflect.DeepEqual(contentType, []string{
			"*/*",
			"text/plain",
			"text/html",
		}) {
			t.Errorf("expected accept header [*/*, text/plain, text/html], got %+v", contentType)
		}
	})

	t.Run("multiple targets", func(t *testing.T) {
		conf, err := Load(complete)
		if err != nil {
			t.Errorf("failed to parse simple config: %s", err)
			t.FailNow()
		}

		if len(conf.Targets) != 2 {
			t.Errorf("expected two targets, got %d", len(conf.Targets))
		}

		for i := 0; i < 2; i++ {
			target := conf.Targets[i]
			if target.Name == "" {
				t.Errorf("expected target to have a name at index %d", i)
			}

			if i == 1 {
				if target.Headers == nil || len(*target.Headers) == 0 {
					t.Errorf("no headers set on target")
				}

				if contentType := target.Headers.Values("Content-Type"); !reflect.DeepEqual(contentType, []string{"application/json"}) {
					t.Errorf("expected content type header `application/json`, got %s", contentType)
				}

				if target.Delay != 10000 {
					t.Errorf("expected delay of 10000, got %d", target.Delay)
				}
			}

			if i == 2 {
				if target.Delay != 20000 {
					t.Errorf("expected delay of 20000, got %d", target.Delay)
				}
			}
		}
	})
}

var simple = `
targets:
  - url: http://example.org
`

var headers = `
targets:
  - url: http://example.org
    headers:
      "Content-Type":
        - "application/json"
      "Accept":
        - "*/*"
        - "text/plain"
        - "text/html"
`

var complete = `
targets:
  - name: foo
    url: http://foo.org
    duration: 10000
    jitter: 0.1
  - name: var
    url: http://bar.org
    duration: 20000
    jitter: 0.2
    headers:
      "Content-Type":
        - "application/json"
`
