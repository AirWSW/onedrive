package cache_test

import (
	"testing"
	"time"
)

func Test(t *testing.T) {
	lastModifiedAt, _ := time.Parse("2020-01-05T22:53:19Z", "2020-01-05T22:53:19Z")
	t.Logf("%s", lastModifiedAt)
}
