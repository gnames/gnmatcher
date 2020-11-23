package rest

import (
	"time"

	"github.com/GeertJohan/go.rice/embedded"
)

func init() {

	// define files
	file4 := &embedded.EmbeddedFile{
		Filename:    "api/v1/openapi.yaml",
		FileModTime: time.Unix(1606168633, 0),

		Content: string("---\ntest: test\n"),
	}

	// define dirs
	dir1 := &embedded.EmbeddedDir{
		Filename:   "",
		DirModTime: time.Unix(1606168848, 0),
		ChildFiles: []*embedded.EmbeddedFile{},
	}
	dir2 := &embedded.EmbeddedDir{
		Filename:   "api",
		DirModTime: time.Unix(1606168598, 0),
		ChildFiles: []*embedded.EmbeddedFile{},
	}
	dir3 := &embedded.EmbeddedDir{
		Filename:   "api/v1",
		DirModTime: time.Unix(1606168633, 0),
		ChildFiles: []*embedded.EmbeddedFile{
			file4, // "api/v1/openapi.yaml"

		},
	}

	// link ChildDirs
	dir1.ChildDirs = []*embedded.EmbeddedDir{
		dir2, // "api"

	}
	dir2.ChildDirs = []*embedded.EmbeddedDir{
		dir3, // "api/v1"

	}
	dir3.ChildDirs = []*embedded.EmbeddedDir{}

	// register embeddedBox
	embedded.RegisterEmbeddedBox(`assets`, &embedded.EmbeddedBox{
		Name: `assets`,
		Time: time.Unix(1606168848, 0),
		Dirs: map[string]*embedded.EmbeddedDir{
			"":       dir1,
			"api":    dir2,
			"api/v1": dir3,
		},
		Files: map[string]*embedded.EmbeddedFile{
			"api/v1/openapi.yaml": file4,
		},
	})
}
