package subsonic

import (
	"encoding/xml"
	"log"
	"net/http"
	"path"

	"github.com/mdlayher/wavepipe/config"

	"github.com/martini-contrib/render"
)

// MusicFoldersContainer contains a list of emulated Subsonic music folders
type MusicFoldersContainer struct {
	// Container name
	XMLName xml.Name `xml:"musicFolders,omitempty"`

	// Music folders
	MusicFolders []MusicFolder `xml:"musicFolder"`
}

// GetMusicFolders is used in Subsonic to return a list of random songs
func GetMusicFolders(req *http.Request, res http.ResponseWriter, r render.Render) {
	// Load name of media folder from config
	conf, err := config.C.Load()
	if err != nil {
		log.Println(err)
		r.XML(200, ErrGeneric)
		return
	}

	// Create a new response container
	c := newContainer()
	c.MusicFolders = &MusicFoldersContainer{
		MusicFolders: []MusicFolder{MusicFolder{0, path.Base(conf.Media())}},
	}

	// Write response
	r.XML(200, c)
}
