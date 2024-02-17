package database

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestContainment(t *testing.T) {
	db := New("abc123")

	entry1 := Entry{
		Title:     "title-123",
		Passwords: []string{"abc", "def"},
	}

	entry2 := Entry{
		Title:     "title-456",
		Passwords: []string{"ghi", "jkl"},
	}

	entry3 := Entry{
		Title:     "title-789",
		Passwords: []string{"mno", "pqr"},
	}

	db.Add(&entry1)
	db.Add(&entry2)

	require.True(t, db.All().contains(&entry1))
	require.True(t, db.All().contains(&entry2))
	require.False(t, db.All().contains(&entry3))
}

func TestGetPassword(t *testing.T) {
	db := New("abc123")

	entry1 := Entry{
		Title:     "title-123",
		Passwords: []string{"abc", "def"},
	}

	entry2 := Entry{
		Title:     "title-456",
		Passwords: []string{"ghi", "jkl"},
	}

	db.Add(&entry1)
	db.Add(&entry2)

	require.Equal(t, db.All()[0].GetPassword(), "def")
	require.Equal(t, db.All()[1].GetPassword(), "jkl")
}

func TestRoundtripJSON(t *testing.T) {
	oldDB := New("abc123")

	entry1 := Entry{
		Title:     "title-123",
		Passwords: []string{"abc", "def"},
	}

	entry2 := Entry{
		Title:     "title-456",
		Passwords: []string{"ghi", "jkl"},
	}

	oldDB.Add(&entry1)
	oldDB.Add(&entry2)

	oldDBAsJSON, err := oldDB.ToJSON()
	require.NoError(t, err, "failed to serialize database")

	newDB, err := FromJSON(oldDBAsJSON)
	require.NoError(t, err, "failed to deserialize database")

	require.Equal(t, oldDB, newDB)
}

func TestActiveEntries(t *testing.T) {
	db := New("abc123")

	entry1 := Entry{
		Title:     "title-123",
		Passwords: []string{"abc", "def"},
		Active:    true,
	}

	entry2 := Entry{
		Title:     "title-456",
		Passwords: []string{"ghi", "jkl"},
		Active:    false,
	}

	db.Add(&entry1)
	db.Add(&entry2)

	require.True(t, db.All().contains(&entry1))
	require.True(t, db.All().contains(&entry2))

	require.True(t, db.Active().contains(&entry1))
	require.False(t, db.Active().contains(&entry2))

	require.False(t, db.Inactive().contains(&entry1))
	require.True(t, db.Inactive().contains(&entry2))
}

func TestFilterByAny(t *testing.T) {
	db := New("abc123")

	entry1 := Entry{
		Title:     "abc def",
		Passwords: []string{"pass-abc", "pass-def"},
	}

	entry2 := Entry{
		Title:     "def ghi",
		Passwords: []string{"pass-ghi", "pass-jkl"},
	}

	entry3 := Entry{
		Title:     "ghi jkl",
		Passwords: []string{"pass-mno", "pass-pqr"},
	}

	db.Add(&entry1)
	db.Add(&entry2)
	db.Add(&entry3)

	require.True(t, db.All().contains(&entry1))
	require.True(t, db.All().contains(&entry2))
	require.True(t, db.All().contains(&entry3))

	require.False(t, db.All().FilterByAny("xyz").contains(&entry1))
	require.True(t, db.All().FilterByAny("abc").contains(&entry1))
	require.True(t, db.All().FilterByAny("abc", "def").contains(&entry1))
	require.True(t, db.All().FilterByAny("abc", "def", "ghi").contains(&entry1))
	require.True(t, db.All().FilterByAny("def", "ghi").contains(&entry1))
	require.False(t, db.All().FilterByAny("ghi").contains(&entry1))

	require.False(t, db.All().FilterByAny("abc").contains(&entry2))
	require.True(t, db.All().FilterByAny("abc", "def").contains(&entry2))
	require.True(t, db.All().FilterByAny("abc", "def", "ghi").contains(&entry2))
	require.True(t, db.All().FilterByAny("def", "ghi").contains(&entry2))
	require.True(t, db.All().FilterByAny("ghi").contains(&entry2))
	require.True(t, db.All().FilterByAny("def").contains(&entry2))
	require.True(t, db.All().FilterByAny("def", "ghi").contains(&entry2))
	require.True(t, db.All().FilterByAny("def", "ghi", "jkl").contains(&entry2))
	require.True(t, db.All().FilterByAny("ghi", "jkl").contains(&entry2))
	require.False(t, db.All().FilterByAny("jkl").contains(&entry2))

	require.False(t, db.All().FilterByAny("def").contains(&entry3))
	require.True(t, db.All().FilterByAny("def", "ghi").contains(&entry3))
	require.True(t, db.All().FilterByAny("def", "ghi", "jkl").contains(&entry3))
	require.True(t, db.All().FilterByAny("ghi", "jkl").contains(&entry3))
	require.True(t, db.All().FilterByAny("jkl").contains(&entry3))
	require.False(t, db.All().FilterByAny("xyz").contains(&entry3))
}

func TestFilterByAll(t *testing.T) {
	db := New("abc123")

	entry1 := Entry{
		Title:     "abc def",
		Passwords: []string{"pass-abc", "pass-def"},
	}

	entry2 := Entry{
		Title:     "def ghi",
		Passwords: []string{"pass-ghi", "pass-jkl"},
	}

	entry3 := Entry{
		Title:     "ghi jkl",
		Passwords: []string{"pass-mno", "pass-pqr"},
	}

	db.Add(&entry1)
	db.Add(&entry2)
	db.Add(&entry3)

	require.True(t, db.All().contains(&entry1))
	require.True(t, db.All().contains(&entry2))
	require.True(t, db.All().contains(&entry3))

	require.False(t, db.All().FilterByAll("xyz").contains(&entry1))
	require.True(t, db.All().FilterByAll("abc").contains(&entry1))
	require.True(t, db.All().FilterByAll("def").contains(&entry1))

	require.False(t, db.All().FilterByAll("xyz", "abc").contains(&entry1))
	require.True(t, db.All().FilterByAll("abc", "def").contains(&entry1))
	require.False(t, db.All().FilterByAll("def", "ghi").contains(&entry1))

	require.False(t, db.All().FilterByAll("abc", "def").contains(&entry2))
	require.True(t, db.All().FilterByAll("def", "ghi").contains(&entry2))
	require.False(t, db.All().FilterByAll("ghi", "jkl").contains(&entry2))
}
