package api

import (
	"io"
	"net/http"
	"strings"
	"testing"
)

type categoryTestTransport func(*http.Request) (*http.Response, error)

func (f categoryTestTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	return f(r)
}

func categoryTestResponse(body string) *http.Response {
	return &http.Response{StatusCode: http.StatusOK, Body: io.NopCloser(strings.NewReader(body))}
}

func TestGetCategoriesIncludesSubcategories(t *testing.T) {
	c := &Client{http: &http.Client{Transport: categoryTestTransport(func(r *http.Request) (*http.Response, error) {
		if r.Method != http.MethodGet || r.URL.Path != "/api/v3.0/get_categories" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		return categoryTestResponse(`{"categories":[{"id":1,"name":"Food","subcategories":[{"id":18,"name":"Groceries"}]}]}`), nil
	})}}
	categories, err := c.GetCategories()
	if err != nil {
		t.Fatal(err)
	}
	if len(categories) != 1 || categories[0].Name != "Food" || len(categories[0].Subcategories) != 1 {
		t.Fatalf("unexpected categories: %+v", categories)
	}
	sub := categories[0].Subcategories[0]
	if sub.ID != 18 || sub.Name != "Groceries" {
		t.Fatalf("unexpected subcategory: %+v", sub)
	}
}

func TestCreateExpenseCategoryParameter(t *testing.T) {
	for _, tc := range []struct {
		name string
		id   int
		want string
	}{{"assigned", 18, "18"}, {"unspecified", 0, ""}} {
		t.Run(tc.name, func(t *testing.T) {
			c := &Client{http: &http.Client{Transport: categoryTestTransport(func(r *http.Request) (*http.Response, error) {
				if r.Method != http.MethodPost || r.URL.Path != "/api/v3.0/create_expense" {
					t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				if err := r.ParseForm(); err != nil {
					t.Fatal(err)
				}
				if got := r.PostForm.Get("category_id"); got != tc.want {
					t.Fatalf("category_id = %q, want %q", got, tc.want)
				}
				if tc.id == 0 && r.PostForm.Has("category_id") {
					t.Fatal("unspecified category must be omitted")
				}
				return categoryTestResponse(`{"expenses":[{"id":123,"category":{"id":18,"name":"Groceries"}}],"errors":{}}`), nil
			})}}
			expense, err := c.CreateExpense(CreateExpenseParams{
				Description: "Groceries", Cost: "10.00", GroupID: 1,
				SplitEqually: true, CategoryID: tc.id,
			})
			if err != nil {
				t.Fatal(err)
			}
			if expense.Category == nil || expense.Category.ID != 18 {
				t.Fatalf("missing category in returned expense: %+v", expense)
			}
		})
	}
}
