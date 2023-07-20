package main

import (
	"database/sql"
	"testing"
)

func Test_populationSum_empty(t *testing.T) {
	cities := []City{}
	got := populationSum(&cities)
	want := map[string]int32{}
	if len(got) != 0 {
		t.Errorf("populationSum(%v) = %v, want %v", cities, got, want)
	}
}

func Test_populationSum_one(t *testing.T) {
	cities := []City{
		{
			ID: 1532,
			Name: "Tokyo",
			CountryCode: sql.NullString{
				String: "JPN",
				Valid: true,
			},
			District: "Tokyo-to",
			Population: sql.NullInt32{
				Int32: 7980230,
				Valid: true,
			},
		},
		{
			ID: 1535,
			Name: "Nagoya",
			CountryCode: sql.NullString{
				String: "JPN",
				Valid: true,
			},
			District: "Aichi",
			Population: sql.NullInt32{
				Int32: 2154376,
				Valid: true,
			},
		},
	}

	got := populationSum(&cities)
	want := map[string]int32{
		"JPN": 2154376 + 7980230,
	}
	if len(got) != len(want) {
		t.Errorf("populationSum(%v) = %v, want %v", cities, got, want)
	}
	if got["JPN"] != want["JPN"] {
		t.Errorf("populationSum(%v) = %v, want %v", cities, got, want)
	}
}

func Test_populationSum_two(t *testing.T) {
	cities := []City{
		{
			ID: 456,
			Name: "London",
			CountryCode: sql.NullString{
				String: "GBR",
				Valid: true,
			},
			District: "England",
			Population: sql.NullInt32{
				Int32: 7285000,
				Valid: true,
			},
		},
		{
			ID: 1535,
			Name: "Nagoya",
			CountryCode: sql.NullString{
				String: "JPN",
				Valid: true,
			},
			District: "Aichi",
			Population: sql.NullInt32{
				Int32: 2154376,
				Valid: true,
			},
		},
	}

	got := populationSum(&cities)
	want := map[string]int32{
		"JPN": 2154376,
		"GBR": 7285000,
	}
	if len(got) != len(want) {
		t.Errorf("populationSum(%v) = %v, want %v", cities, got, want)
	}
	if got["GBR"] != want["GBR"] {
		t.Errorf("populationSum(%v) = %v, want %v", cities, got, want)
	}
	if got["JPN"] != want["JPN"] {
		t.Errorf("populationSum(%v) = %v, want %v", cities, got, want)
	}
}

func Test_populationSum_invalid(t *testing.T) {
	cities := []City{
		{
			ID: 456,
			Name: "London",
			CountryCode: sql.NullString{
				String: "GBR",
				Valid: true,
			},
			District: "England",
			Population: sql.NullInt32{
				Int32: 7285000,
				Valid: true,
			},
		},
		{
			ID: 9999,
			Name: "hoge",
			CountryCode: sql.NullString{
				String: "hge",
				Valid: false,
			},
			District: "fuga",
			Population: sql.NullInt32{
				Int32: 998244353,
				Valid: true,
			},
		},
	}

	got := populationSum(&cities)
	want := map[string]int32{
		"GBR": 7285000,
	}
	if len(got) != len(want) {
		t.Errorf("populationSum(%v) = %v, want %v", cities, got, want)
	}
	if got["GBR"] != want["GBR"] {
		t.Errorf("populationSum(%v) = %v, want %v", cities, got, want)
	}
}