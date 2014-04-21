package api

import (
	"log"
	"net/http"
	"strings"

	"github.com/mdlayher/wavepipe/data"

	"github.com/go-martini/martini"
	"github.com/martini-contrib/render"
	"github.com/mdlayher/goset"
)

// SearchResponse represents the JSON response for /api/search
type SearchResponse struct {
	Error   *Error        `json:"error"`
	Artists []data.Artist `json:"artists"`
	Albums  []data.Album  `json:"albums"`
	Songs   []data.Song   `json:"songs"`
	Folders []data.Folder `json:"folders"`
	render  render.Render `json:"-"`
}

// RenderError renders a JSON error message with the specified HTTP status code and message
func (s *SearchResponse) RenderError(code int, message string) {
	// Nullify all other fields
	s.Artists = nil
	s.Albums = nil
	s.Songs = nil
	s.Folders = nil

	// Generate error
	s.Error = new(Error)
	s.Error.Code = code
	s.Error.Message = message

	// Render with specified HTTP status code
	s.render.JSON(code, s)
}

// ServerError is a shortcut to render a HTTP 500 with generic "server error" message
func (s *SearchResponse) ServerError() {
	s.RenderError(500, "server error")
	return
}

// GetSearch searches for artists, albums, songs, and folders matching a specified search query,
// and returns a HTTP status and JSON
func GetSearch(req *http.Request, r render.Render, params martini.Params) {
	// Output struct for songs request
	res := SearchResponse{render: r}

	// Check API version
	if version, ok := params["version"]; ok {
		// Check if this API call is supported in the advertised version
		if !apiVersionSet.Has(version) {
			res.RenderError(400, "unsupported API version"+version)
			return
		}
	}

	// Check for a search query
	query, ok := params["query"]
	if !ok {
		res.RenderError(400, "no search query specified")
		return
	}

	// Default list of type to include in results
	defaultTypeSet := set.New("artists", "albums", "songs", "folders")

	// Check for a comma-separated list of data types to search
	types := req.URL.Query().Get("type")
	var typeSet *set.Set
	if types == "" {
		// Search for all types if not specified
		typeSet = defaultTypeSet
	} else {
		// Iterate comma-separated list of types
		typeSet = set.New()
		for _, t := range strings.Split(types, ",") {
			// Add valid types to list
			if defaultTypeSet.Has(t) {
				typeSet.Add(t)
			}
		}
	}

	// If selected, include artists
	if typeSet.Has("artists") {
		// Search for artists which match the search query
		artists, err := data.DB.SearchArtists(query)
		if err != nil {
			log.Println(err)
			res.ServerError()
			return
		}

		// Copy into response
		res.Artists = artists
	}

	// If selected, include albums
	if typeSet.Has("albums") {
		// Search for albums which match the search query
		albums, err := data.DB.SearchAlbums(query)
		if err != nil {
			log.Println(err)
			res.ServerError()
			return
		}

		// Copy into response
		res.Albums = albums
	}

	// If selected, include songs
	if typeSet.Has("songs") {
		// Search for songs which match the search query
		songs, err := data.DB.SearchSongs(query)
		if err != nil {
			log.Println(err)
			res.ServerError()
			return
		}

		// Copy into response
		res.Songs = songs
	}

	// If selected, include folders
	if typeSet.Has("folders") {
		// Search for folders which match the search query
		folders, err := data.DB.SearchFolders(query)
		if err != nil {
			log.Println(err)
			res.ServerError()
			return
		}

		// Copy into response
		res.Folders = folders
	}

	// HTTP 200 OK with JSON
	res.Error = nil
	r.JSON(200, res)
	return
}