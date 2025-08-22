/*
    Skimmer is a package for working with feeds and rendering Link Blogs
	Copyright (C) 2025  R. S. Doiel

    This program is free software: you can redistribute it and/or modify
    it under the terms of the GNU Affero General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.

    This program is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU Affero General Public License for more details.

    You should have received a copy of the GNU Affero General Public License
    along with this program.  If not, see <https://www.gnu.org/licenses/>.
*/
package skimmer

import (
	"os"
	"strings"
	"testing"
	"database/sql"
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
	} else if strings.Compare(val.Label, expectedVal) != 0 {
		t.Errorf("expected %q, got %q", expectedVal, val.Label)

	}
}

func TestSetup(t *testing.T) {
	app, err := NewSkimmer("test_skimmer")
	app.out = os.Stdout
	app.eout = os.Stderr
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
	app, err := NewSkimmer("test_skimmer")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	app.out = os.Stdout
	app.eout = os.Stderr
	appDir := "./test_output"
	fName := "test.urls"
	if err := app.Setup(appDir); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := app.ReadUrls(fName); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestDownload(t *testing.T) {
	app, err := NewSkimmer("test_skimmer")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	app.out = os.Stdout
	app.eout = os.Stderr
	appDir := "./test_output"
	fName := "test.urls"
	dsn := app.DBName
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer db.Close()
	if err := app.Setup(appDir); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := app.ReadUrls(fName); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := app.Download(db); err != nil {
		t.Error(err)
		t.FailNow()
	}
}

func TestDisplay(t *testing.T) {
	app, err := NewSkimmer("test_skimmer")
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	app.out = os.Stdout
	app.eout = os.Stderr
	appDir := "./test_output"
	fName := "test.urls"
	if err := app.Setup(appDir); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := app.ReadUrls(fName); err != nil {
		t.Error(err)
		t.FailNow()
	}
	dsn := app.DBName
	db, err := sql.Open("sqlite",dsn)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	defer db.Close()
	if err := app.Download(db); err != nil {
		t.Error(err)
		t.FailNow()
	}
	if err := app.Write(db); err != nil {
		t.Error(err)
		t.FailNow()
	}
}
