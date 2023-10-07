package skimmer

import (
	"os"
	"strings"
	"testing"
)

func TestParseURLList(t *testing.T) {
	src := []byte(`
# This is an example url list. 
https://transitinglosangeles.com/feed/ "~Transiting Los Angeles"
https://laist.com/index.atom "~The LAist"
#https://feeds.mcclatchy.com/sacbee/stories "~Sacramento Bee Stories"
https://apnews.com/world-news.rss "~Associated Press, World News"
https://apnews.com/us-news.rss "~Associated Press, US News"
https://www.theguardian.com/us/rss "~The Guardian US Edition"
`)
	m, err := ParseURLList("test_sample.txt", src)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	expectedKey := "https://laist.com/index.atom"
	expectedVal := `The LAist`
	if val, ok := m[expectedKey]; !ok {
		t.Errorf("expected value for %q, key not found in map", expectedKey)
		t.FailNow()
	} else if strings.Compare(val, expectedVal) != 0 {
		t.Errorf("expected %q, got %q", expectedVal, val)

	}
}

func TestSetup(t *testing.T) {
	app, err := NewSkimmer(os.Stdout, os.Stderr, "test_skimmer")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	appDir := "./test_output"
	if _, err := os.Stat(appDir); err == nil {
		os.RemoveAll(appDir)
	}
	if err := app.Setup(appDir); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestReadUrls(t *testing.T) {
	app, err := NewSkimmer(os.Stdout, os.Stderr, "test_skimmer")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	appDir := "./test_output"
	if err := app.Setup(appDir); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := app.ReadUrls(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestDownload(t *testing.T) {
	app, err := NewSkimmer(os.Stdout, os.Stderr, "test_skimmer")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	appDir := "./test_output"
	if err := app.Setup(appDir); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := app.ReadUrls(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := app.Download(); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestDisplay(t *testing.T) {
	app, err := NewSkimmer(os.Stdout, os.Stderr, "test_skimmer")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	appDir := "./test_output"
	if err := app.Setup(appDir); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := app.ReadUrls(); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := app.Download(); err != nil {
		t.Error(err)
		t.FailNow()
	}

	if err := app.Write(os.Stdout); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
