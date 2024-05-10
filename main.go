package main

import (
    "strings"
    "bufio"
    "os"
    "context"
    "fmt"
    "strconv"
    "github.com/raitonoberu/ytmusic"
    "github.com/lrstanley/go-ytdlp"
)

const (
    Reset  = "\033[0m"
    Red    = "\033[31m"
    Green  = "\033[32m"
    Yellow = "\033[33m"
    Blue   = "\033[34m"
    Purple = "\033[35m"
    Cyan   = "\033[36m"
)

type SearchResult struct {
    Tracks    []*TrackItem    `json:"tracks"`
    Artists   []*ArtistItem   `json:"artists"`
    Albums    []*AlbumItem    `json:"albums"`
    Playlists []*PlaylistItem `json:"playlists"`
    Videos    []*VideoItem    `json:"videos"`
}

type VideoItem struct {
   VideoID    string      `json:"videoId"`
    PlaylistID string      `json:"playlistId"`
    Title      string      `json:"title"`
    Artists    []Artist    `json:"artists"`
    Views      string      `json:"views"`
    Duration   int         `json:"duration"`
    Thumbnails []Thumbnail `json:"thumbnails"`
}

type PlaylistItem struct {
    BrowseID   string      `json:"browseId"`
    Title      string      `json:"title"`
    Author     string      `json:"author"`
    ItemCount  string      `json:"itemCount"`
    Thumbnails []Thumbnail `json:"thumbnails"`
}

type ArtistItem struct {
    BrowseID   string      `json:"browseId"`
    Artist     string      `json:"artist"`
    ShuffleID  string      `json:"shuffleId"`
    RadioID    string      `json:"radioId"`
    Thumbnails []Thumbnail `json:"thumbnails"`
}

type TrackItem struct {
    VideoID    string      `json:"videoId"`
    PlaylistID string      `json:"playlistId"`
    Title      string      `json:"title"`
    Artists    []Artist    `json:"artists"`
    Album      Album       `json:"album"`
    Duration   int         `json:"duration"`
    IsExplicit bool        `json:"isExplicit"`
    Thumbnails []Thumbnail `json:"thumbnails"`
}

type Album struct {
    Name string `json:"name"`
    ID   string `json:"id"`
}

type AlbumItem struct {
    BrowseID   string      `json:"browseId"`
    Title      string      `json:"title"`
    Type       string      `json:"type"`
    Artists    []Artist    `json:"artists"`
    Year       string      `json:"year"`
    IsExplicit bool        `json:"isExplicit"`
    Thumbnails []Thumbnail `json:"thumbnails"`
}

type Thumbnail struct {
    URL    string `json:"url"`
    Width  int    `json:"width"`
    Height int    `json:"height"`
}

type Artist struct {
    Name string `json:"name"`
    ID   string `json:"id"`
}


func main() {
    ytdlp.MustInstall(context.TODO(), nil)
    reader := bufio.NewReader(os.Stdin)
    fmt.Print(Yellow + "Enter an artist name or track: " + Reset)
    input, err := reader.ReadString('\n')
    if err != nil {
        panic(err)
        return
    }

    s := ytmusic.TrackSearch(input)

    var limit string
    fmt.Print(Yellow + "How many pages to search: " + Reset)
    _, err = fmt.Scan(&limit)

    if err != nil {
        panic(err)
        return
    }

    int_limit, _ := strconv.Atoi(limit)

    var all_results []TrackItem
    var running_int int

    for {
        // TODO: Instead of iterating this, ask the user if they want to search more after each search
        result, err := s.Next()
        if err != nil || len(all_results) >= int_limit {
            // panic(err)
            break
        }

        for _, track := range result.Tracks {
            running_int += 1
            var artists []Artist
            for _, artist := range track.Artists {
                artists = append(artists, Artist{
                    Name: artist.Name,
                    ID: artist.ID,
                })
            }

            all_results = append(all_results, TrackItem {
                VideoID: track.VideoID,
                Title: track.Title,
                Artists: artists,
                Album: Album{
                    Name: track.Album.Name,
                    ID: track.Album.ID,
                },
                Duration: track.Duration,
            })

            fmt.Print(Yellow + "[", running_int, "]" + Reset)
            fmt.Print(Cyan + track.Artists[0].Name + Reset + ": ")
            fmt.Print(Blue + track.Album.Name + Reset + ": ")
            fmt.Print(Purple + track.Title + Reset + "\n")
        }
    }

    var dl string
    fmt.Print(Yellow + "\nDownload [A]ll, [N]one, [1,2,3,...]: " + Reset)
    _, err = fmt.Scan(&dl)
    if err != nil {
        panic(err)
        return
    }

    if strings.ToLower(dl) == "a" {
        for _, song := range all_results {
            fmt.Println(Green + "Downloading " + Reset + Purple + song.Album.Name + ": " + song.Title + Reset)
            download_track(song.VideoID, song.Album.Name, song.Artists[0].Name)
        }
    } else if strings.ToLower(dl) == "n" {
        return
    } else if len(strings.Split(dl, ",")) > 0 {
        for _, i := range strings.Split(dl, ",") {
            song_int, _ := strconv.Atoi(i)
            for j, song := range all_results {
                if song_int != j {
                   continue 
                }
                fmt.Println(Green + "Downloading " + Reset + Purple + song.Album.Name + ":" + song.Title + Reset)
                download_track(song.VideoID, song.Album.Name, song.Artists[0].Name)
            }
        }
    }
    fmt.Print(Green + "Done!" + Reset)

}

func download_track (video_id string, album string, artist string) {
    var out_dir string = "downloads"
    if _, err := os.Stat(out_dir); os.IsNotExist(err) {
        err := os.Mkdir(out_dir, 0755)
        if err != nil {
            fmt.Println("Error creating directory:", err)
        }
    }

    dl := ytdlp.New().
        PrintJSON().
        NoProgress().
        FormatSort("res,ext:m4a").
        ExtractAudio().
        NoPlaylist().
        NoOverwrites().
        Continue().
        Output(out_dir + "/" + artist + "/" + album + "/" + "%(title)s.%(ext)s")

    // Create artist directory if not exist
    if _, err := os.Stat(out_dir + "/" + artist); os.IsNotExist(err) {
        err := os.Mkdir(out_dir + "/" + artist, 0755)
        if err != nil {
            fmt.Println("Error creating directory:", err)
        }
    }

    // Create album directory if not exist
    if _, err := os.Stat(out_dir + "/" + artist + "/" + album); os.IsNotExist(err) {
        err := os.Mkdir(out_dir + "/" + artist + "/" + album, 0755)
        if err != nil {
            fmt.Println("Error creating directory:", err)
        }
    }

    fmt.Print(video_id)
    var filename string = "https://www.youtube.com/watch?v=" + video_id
    // filename = strings.ReplaceAll(filename, "/", "_")
    _, err := dl.Run(context.TODO(), filename)
    if err != nil {
        panic(err)
    }
}
