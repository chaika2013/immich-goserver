package test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetAllJobsStatusNotLoggedIn(t *testing.T) {
	router := NewRouter(t)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs", nil)
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	fmt.Println(w.Body.String())
}

func TestGetAllJobsStatusNotAdmin(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	token := FakeLogin(t, router, 2)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs", nil)
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	fmt.Println(w.Body.String())
}

func TestGetAllJobsStatus(t *testing.T) {
	router := FromScratch(t)
	AddTestUsers(t)
	token := FakeLogin(t, router, 1)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/jobs", nil)
	req.AddCookie(&http.Cookie{Name: "immich_access_token", Value: token})
	router.e.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	assert.JSONEq(t, `
	{
		"metadata-extraction-queue": {
			"jobCounts": {
				"active": 0,
				"completed": 0,
				"failed": 0,
				"delayed": 0,
				"waiting": 0,
				"parsed": 0
			},
			"queueStatus": {
				"isActive": false,
				"isPaused": false
			}
		},
		"storage-template-migration-queue": {
			"jobCounts": {
				"active": 0,
				"completed": 0,
				"failed": 0,
				"delayed": 0,
				"waiting": 0,
				"parsed": 0
			},
			"queueStatus": {
				"isActive": false,
				"isPaused": false
			}
		},
		"thumbnail-generation-queue": {
			"jobCounts": {
				"active": 0,
				"completed": 0,
				"failed": 0,
				"delayed": 0,
				"waiting": 0,
				"parsed": 0
			},
			"queueStatus": {
				"isActive": false,
				"isPaused": false
			}
		},
		"video-conversion-queue": {
			"jobCounts": {
				"active": 0,
				"completed": 0,
				"failed": 0,
				"delayed": 0,
				"waiting": 0,
				"parsed": 0
			},
			"queueStatus": {
				"isActive": false,
				"isPaused": false
			}
		},
		"object-tagging-queue": {
			"jobCounts": {
				"active": 0,
				"completed": 0,
				"failed": 0,
				"delayed": 0,
				"waiting": 0,
				"parsed": 0
			},
			"queueStatus": {
				"isActive": false,
				"isPaused": false
			}
		},
		"clip-encoding-queue": {
			"jobCounts": {
				"active": 0,
				"completed": 0,
				"failed": 0,
				"delayed": 0,
				"waiting": 0,
				"parsed": 0
			},
			"queueStatus": {
				"isActive": false,
				"isPaused": false
			}
		},
		"background-task-queue": {
			"jobCounts": {
				"active": 0,
				"completed": 0,
				"failed": 0,
				"delayed": 0,
				"waiting": 0,
				"parsed": 0
			},
			"queueStatus": {
				"isActive": false,
				"isPaused": false
			}
		},
		"search-queue": {
			"jobCounts": {
				"active": 0,
				"completed": 0,
				"failed": 0,
				"delayed": 0,
				"waiting": 0,
				"parsed": 0
			},
			"queueStatus": {
				"isActive": false,
				"isPaused": false
			}
		},
		"recognize-faces-queue": {
			"jobCounts": {
				"active": 0,
				"completed": 0,
				"failed": 0,
				"delayed": 0,
				"waiting": 0,
				"parsed": 0
			},
			"queueStatus": {
				"isActive": false,
				"isPaused": false
			}
		},
		"sidecar-queue": {
			"jobCounts": {
				"active": 0,
				"completed": 0,
				"failed": 0,
				"delayed": 0,
				"waiting": 0,
				"parsed": 0
			},
			"queueStatus": {
				"isActive": false,
				"isPaused": false
			}
		}
	}`, w.Body.String())
}
