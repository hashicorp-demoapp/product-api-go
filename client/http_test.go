package client

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func setupHTTPTests(t *testing.T) (*HTTP, func()) {
	ts := httptest.NewServer(
		http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			if r.URL.Path == "/coffees" {
				rw.Write([]byte(`
					[{"id": 1, "name": "test1"}, {"id": 2, "name": "test2"}]
				`))
				return
			}

			if strings.HasSuffix(r.URL.Path, "/ingredients") {
				rw.Write([]byte(`
					[{"id": 1, "name": "test2", "quantity": 3, "unit": "g"}, {"id": 2, "name": "ingredient2", "quantity": 6, "unit": "g"}]
				`))
				return
			}

			if strings.HasPrefix(r.URL.Path, "/coffees/") {
				rw.Write([]byte(`
					{"id": 2, "name": "test2"}
				`))
				return
			}

		},
		))

	return NewHTTP(ts.URL), func() {
		ts.Close()
	}
}

func TestGetsCoffees(t *testing.T) {
	c, cleanup := setupHTTPTests(t)
	defer cleanup()

	cof, err := c.GetCoffees()

	assert.NoError(t, err)
	assert.Len(t, cof, 2)
}

func TestGetsCoffee(t *testing.T) {
	c, cleanup := setupHTTPTests(t)
	defer cleanup()

	cof, err := c.GetCoffee(2)

	assert.NoError(t, err)
	assert.Equal(t, 2, cof.ID)
}

func TestGetsIngredients(t *testing.T) {
	c, cleanup := setupHTTPTests(t)
	defer cleanup()

	in, err := c.GetIngredientsForCoffee(2)

	assert.NoError(t, err)
	assert.Len(t, in, 2)
}
