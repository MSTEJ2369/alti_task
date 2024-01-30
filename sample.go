package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

const (
	lastfmkey = "YOUR_LASTFM_API_KEY"
	musixKey  = "YOUR_MUSIXMATCH_API_KEY"
	lastfmurl = "https://ws.audioscrobbler.com/2.0/"
	musixurl  = "https://api.musixmatch.com/ws/1.1/"
)

type mmicroserv struct{}

type lfmResponse struct {
	Tracks struct {
		Track []struct {
			Name   string `json:"name"`
			Artist struct {
				Name string `json:"name"`
			} `json:"artist"`
		} `json:"track"`
	} `json:"tracks"`
}

type musixmatchlyricsRes struct {
	Message struct {
		Body struct {
			Lyrics struct {
				LyricsBody string `json:"lyrics_body"`
			} `json:"lyrics"`
		} `json:"body"`
	} `json:"message"`
}

type lastFMArtistRes struct {
	Artist struct {
		Name  string `json:"name"`
		Image []struct {
			Text string `json:"#text"`
			Size string `json:"size"`
		} `json:"image"`
	} `json:"artist"`
}

func (ms *mmicroserv) getTopTrack(region string) (string, string, error) {
	url := fmt.Sprintf("%s?method=geo.getTopTracks&country=%s&api_key=%s&format=json", lastfmurl, region, lastfmkey)
	resp, err := http.Get(url)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	var trackResponse lfmResponse
	err = json.NewDecoder(resp.Body).Decode(&trackResponse)
	if err != nil {
		return "", "", err
	}

	topTrack := trackResponse.Tracks.Track[0]
	return topTrack.Name, topTrack.Artist.Name, nil
}

func (ms *mmicroserv) getLyrics(trackName, artistName string) (string, error) {
	url := fmt.Sprintf("%s?apikey=%s&q_track=%s&q_artist=%s", musixurl+"matcher.lyrics.get", musixKey, strings.ReplaceAll(trackName, " ", "%20"), strings.ReplaceAll(artistName, " ", "%20"))
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var lyricsResponse musixmatchlyricsRes
	err = json.NewDecoder(resp.Body).Decode(&lyricsResponse)
	if err != nil {
		return "", err
	}

	lyrics := lyricsResponse.Message.Body.Lyrics.LyricsBody
	return lyrics, nil
}

func (ms *mmicroserv) getArtist(artistName string) (string, error) {
	url := fmt.Sprintf("%s?method=artist.getinfo&artist=%s&api_key=%s&format=json", lastfmurl, strings.ReplaceAll(artistName, " ", "%20"), lastfmkey)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var artistResponse lastFMArtistRes
	err = json.NewDecoder(resp.Body).Decode(&artistResponse)
	if err != nil {
		return "", err
	}

	return artistResponse.Artist.Name, nil
}

func main() {
	microservice := mmicroserv{}
	region := "USA" // Example region
	topTrack, artistName, err := microservice.getTopTrack(region)
	if err != nil {
		fmt.Println("Error getting top track:", err)
		return
	}

	lyrics, err := microservice.getLyrics(topTrack, artistName)
	if err != nil {
		fmt.Println("Error getting lyrics:", err)
		return
	}

	artistInfo, err := microservice.getArtist(artistName)
	if err != nil {
		fmt.Println("Error getting artist info:", err)
		return
	}

	fmt.Println("Top Track:", topTrack)
	fmt.Println("Artist:", artistName)
	fmt.Println("Lyrics:", lyrics)
	fmt.Println("Artist Info:", artistInfo)
	// need google API for google image search
}
