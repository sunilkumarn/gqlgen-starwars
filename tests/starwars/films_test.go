package starwars

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/jarcoal/httpmock"
	"github.com/peterhellberg/swapi"

	"gqlgen-starwars/handlers"
	"gqlgen-starwars/tests/testutils"
)

func TestFilmQuery(t *testing.T) {
	t.Run("when film is returned from api", func(t *testing.T) {
		// Arrange
		httpmock.Activate()
		defer httpmock.Deactivate()

		m := testutils.NewMockRequest("GET", "/api/films/1", http.StatusOK)
		m.RespondWith(t, swapi.Film{Title: "Good Movie"})

		c := testutils.SwapiClient(m.URL(t))

		query := `query { film(id: \"1\") { name } }`
		req := testutils.NewGraphQLRequest(t, query)

		response := httptest.NewRecorder()

		// Act
		testutils.PerformGraphQLRequest(response, req, handlers.SwapiClient(c))

		// Assert
		testutils.AssertSuccess(t, response)

		expected := `{"film":{"name":"Good Movie"}}`
		testutils.AssertGraphQLData(t, response, expected)
		testutils.AssertGraphQLErrors(t, response, []string{})
	})

	t.Run("when api returns error", func(t *testing.T) {
		// Arrange
		httpmock.Activate()
		defer httpmock.Deactivate()

		m := testutils.NewMockRequest("GET", "/api/films/1", http.StatusInternalServerError)
		m.RespondWith(t, "")

		c := testutils.SwapiClient(m.URL(t))

		query := `query { film(id: \"1\") { name } }`
		req := testutils.NewGraphQLRequest(t, query)

		response := httptest.NewRecorder()

		// Act
		testutils.PerformGraphQLRequest(response, req, handlers.SwapiClient(c))

		// Assert
		testutils.AssertSuccess(t, response)
		testutils.AssertGraphQLData(t, response, "null")
		testutils.AssertGraphQLErrors(t, response, []string{"Failed to fetch film"})
	})

	t.Run("when query contains both film and character", func(t *testing.T) {
		// Arrange
		httpmock.Activate()
		defer httpmock.Deactivate()

		m1 := testutils.NewMockRequest("GET", "/api/films/1", http.StatusOK)
		m1.RespondWith(t, swapi.Film{Title: "Good Movie"})

		m2 := testutils.NewMockRequest("GET", "/api/people/2", http.StatusOK)
		m2.RespondWith(t, swapi.Person{Name: "John Smith"})

		c := testutils.SwapiClient(m1.URL(t))

		query := `query { film(id: \"1\") { name } character(id: \"2\") { name } }`
		req := testutils.NewGraphQLRequest(t, query)

		response := httptest.NewRecorder()

		// Act
		testutils.PerformGraphQLRequest(response, req, handlers.SwapiClient(c))

		// Assert
		testutils.AssertSuccess(t, response)

		expected := `{"film":{"name":"Good Movie"},"character":{"name":"John Smith"}}`
		testutils.AssertGraphQLData(t, response, expected)
		testutils.AssertGraphQLErrors(t, response, []string{})
	})
}
