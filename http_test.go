package mecachis

import (
	"fmt"
	"github.com/sonirico/mecachis/engines"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type state struct {
	group string
	key   string
	value string
}

type action struct {
	endpoint string
	method   string
	payload  string
}

type testCase struct {
	name       string
	action     *action
	want       string
	wantStatus int
	state      *state
}

func prepareHub(t *testing.T, handler http.Handler, actions []action) {
	t.Helper()

	for _, action := range actions {
		req := httptest.NewRequest(action.method, action.endpoint, strings.NewReader(action.payload))
		recorder := httptest.NewRecorder()

		handler.ServeHTTP(recorder, req)

		if recorder.Code < http.StatusOK || recorder.Code > 299 {
			t.Errorf("unexpected status code. want ok, have %d", recorder.Code)
		}
	}
}

func TestHub_ServeHTTP(t *testing.T) {
	tests := []testCase{
		{
			name: "with non-existent group and key",
			action: &action{
				endpoint: "/mecachis/myapp/key",
				method:   http.MethodGet,
			},
			want:       "404 page not found",
			wantStatus: http.StatusNotFound,
		},
		{
			name: "with existent group and key",
			action: &action{
				endpoint: "/mecachis/monitoring/sla",
				method:   http.MethodGet,
			},
			want:       "perfdata: 100%",
			wantStatus: http.StatusOK,
			state:      &state{group: "monitoring", key: "sla", value: "perfdata: 100%"},
		},
		{
			name: "with a wrongly prefixed uri",
			action: &action{
				endpoint: "/unknownnamespace/g/k",
				method:   http.MethodGet,
			},
			want:       "404 page not found",
			wantStatus: http.StatusNotFound,
		},
		{
			name: "with a wrong number of path params ('key' is missing)",
			action: &action{
				endpoint: "/mecachis/myapp/",
				method:   http.MethodGet,
			},
			want:       "404 page not found",
			wantStatus: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			action := test.action
			request := httptest.NewRequest(action.method, action.endpoint, strings.NewReader(action.payload))
			responseRecorder := httptest.NewRecorder()

			hub := NewHub()
			if test.state != nil {
				g, _ := hub.getOrCreateGroup(test.state.group)
				_ = g.Add(test.state.key, MemoryView(test.state.value))
			}
			hub.ServeHTTP(responseRecorder, request)

			if responseRecorder.Code != test.wantStatus {
				t.Errorf("unexpected status code. want %d, have %d",
					test.wantStatus, responseRecorder.Code)
			}

			cleanedBody := strings.TrimSpace(responseRecorder.Body.String())
			if cleanedBody != test.want {
				t.Errorf("unexpected response body. want '%s', have '%s'",
					test.want, cleanedBody)
			}
		})
	}
}

func TestHub_ServeHTTP_LRU_engine_eviction(t *testing.T) {
	var capacity uint64 = 14
	actions := []action{
		{
			method:   http.MethodPost,
			endpoint: fmt.Sprintf("/mecachis/metrics/mem?engi=lru&cap=%d", capacity),
			payload:  "13gb", // +7
		},
		{
			method:   http.MethodPost,
			endpoint: "/mecachis/metrics/ping",
			payload:  "10ms", // +8
		},
	}

	hub := NewHub()
	prepareHub(t, hub, actions)

	group, ok := hub.group("metrics")
	if !ok {
		t.Errorf("want group, have none")
	}
	if group.Ct != engines.LRU {
		t.Errorf("unexpected engine type. want LRU, have '%v'", group.Ct)
	}
	if group.Cap != capacity {
		t.Errorf("unexpected cache capacity. want %d, have %d", capacity, group.Cap)
	}
	val, ok := group.Get("mem")
	if ok {
		t.Errorf("unexpected cache result. expected eviction, have '%s'", val.String())
	}
}
