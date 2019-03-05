package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"sort"
	"strconv"
	"sync"
	"time"
)

type TwitchResponse struct {
	Stream struct {
		ID          *int       `json:"id"`
		Game        *string    `json:"game"`
		Viewers     *int       `json:"viewers"`
		VideoHeight *int       `json:"videoheight"`
		AverageFps  *int       `json:"averggefps"`
		Delay       *int       `json:"delay"`
		CreatedAt   *time.Time `json:"createdAt"`
		IsPlaylist  *bool      `json:"isPlaylist"`
		Preview     struct {
			Small    *string `json:"small"`
			Medium   *string `json:"medium"`
			Large    *string `json:"large"`
			Template *string `json:"template"`
		}
		Channel struct {
			Mature                       bool        `json:"mature"`
			Status                       *string     `json:"status"`
			BroadcasterLanguage          string      `json:"broadcaster_language"`
			DisplayName                  string      `json:"display_name"`
			Game                         string      `json:"game"`
			Language                     string      `json:"language"`
			ID                           int         `json:"_id"`
			Name                         string      `json:"name"`
			CreatedAt                    time.Time   `json:"created_at"`
			UpdatedAt                    time.Time   `json:"updated_at"`
			Partner                      bool        `json:"partner"`
			Logo                         string      `json:"logo"`
			VideoBanner                  string      `json:"video_banner"`
			ProfileBanner                string      `json:"profile_banner"`
			ProfileBannerBackgroundColor interface{} `json:"profile_banner_background_color"`
			URL                          string      `json:"url"`
			Views                        int         `json:"views"`
			Followers                    int         `json:"followers"`
		} `json:"channel"`
	} `json:"stream"`
}

type Newlive struct {
	Name        *string    `json:"name"`
	ImageID     *string    `json:"imageId"`
	ChannelID   *string    `json:"channelId"`
	Title       *string    `json:"title"`
	Description *string    `json:"description"`
	Viewers     *int       `json:"viewers"`
	Likes       *string    `json:"likes"`
	Dislikes    *string    `json:"dislikes"`
	VideoID     *string    `json:"videoId"`
	Thumbnail   Thumbnails `json:"thumbnails"`
	Type        string     `json:"type"`
}

// []string{"hasanabi", "destiny", "invadervie", "richardlewisreports", "hitch", "cjayride", "trainwreckstv"}
var streamers = []Streamer{
	{Name: "Ice Poseidon", ChannelId: "UCv9Edl_WbtbPeURPtFDo-uA", ImageID: "ice", Type: "youtube"},
	{Name: "Hyphonix", ChannelId: "UCaFpm67qMk1W1wJkFhGXucA", ImageID: "hyphonix", Type: "youtube"},
	{Name: "Gary", ChannelId: "UCvxSwu13u1wWyROPlCH-MZg", ImageID: "gary", Type: "youtube"},
	{Name: "Cxnews", ChannelId: "UCStEQ9BjMLjHTHLNA6cY9vg", ImageID: "cxnews", Type: "youtube"},
	{Name: "Voldesad", ChannelId: "UCPkOhci8gkwL7p6hxIJ2WQw", ImageID: "vold", Type: "youtube"},
	{Name: "Cassandra", ChannelId: "UCoQnCN55E25nGavk79Asyng", ImageID: "cass", Type: "youtube"},
	{Name: "Juan Bagnell", ChannelId: "UCvhnYODy6WQ0mw_zi3V1h0g", ImageID: "juan", Type: "youtube"},
	{Name: "Coding Train", ChannelId: "UCvjgXvBlbQiydffZU7m1_aw", ImageID: "coding", Type: "youtube"},
	{Name: "Joe Rogan Podcast", ChannelId: "UCzQUP1qoWDoEbmsQxvdjxgQ", ImageID: "joe", Type: "youtube"},
	{Name: "Mixhound", ChannelId: "UC_jxnWLGJ2eQK4en3UblKEw", ImageID: "mix", Type: "youtube"},
	{Name: "Hasanabi", Type: "twitch", ImageID: "hasanabi"},
	{Name: "Destiny", Type: "twitch", ImageID: "destiny"},
	{Name: "Invadervie", Type: "twitch", ImageID: "invadervie"},
	{Name: "Richardlewisreports", Type: "twitch", ImageID: "richardlewis"},
	{Name: "Cjayride", Type: "twitch", ImageID: "cjayride"},
	{Name: "Hitch", Type: "twitch", ImageID: "hitch"},
	{Name: "Rajjpatel", Type: "twitch", ImageID: "hitch"},
}
var payload = make(chan Newlive)
var done = make(chan bool)

var wait sync.WaitGroup

func GetYoutube() {
	wait.Add(2)
	fmt.Println("getting....")
	ch := make(chan *Islive)
	go func() {
		defer close(ch)
		for i, v := range streamers {
			if v.Type == "twitch" {
				go getTwitch(v, i)
				continue
			}
			url := fmt.Sprintf("https://www.googleapis.com/youtube/v3/search?part=snippet&channelId=%v&eventType=live&type=video&key=%v", v.ChannelId, mykey)
			resp, err := http.Get(url)
			if err != nil || resp.StatusCode != 200 {
				fmt.Println(err)
				continue
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			var streamer Islive
			json.Unmarshal(body, &streamer)
			if streamer.PageInfo.TotalResults == 0 {
				continue
			}
			streamer.Name = v.Name
			streamer.ImageID = v.ImageID
			ch <- &streamer
		}
	}()
	go func() {
		if ch == nil {
			return
		}
		for v := range ch {
			if v == nil {
				fmt.Println("nil value")
				continue
			}
			id := v.Items[0].ID.VideoID
			resp, err := http.Get("https://www.googleapis.com/youtube/v3/videos?part=statistics%2C+snippet%2C+liveStreamingDetails&id=" + id + "&key=" + mykey)
			if err != nil || resp.StatusCode != 200 {
				fmt.Println(err)
				continue
			}
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				fmt.Println(err)
			}
			defer resp.Body.Close()
			var live Livestream
			json.Unmarshal(body, &live)
			name, err := strconv.Atoi(live.Items[0].LiveStreamingDetails.ConcurrentViewers)
			if err != nil {
				fmt.Println(name)
				fmt.Println(err)
				continue
			}
			thumb := Thumbnails{High: live.Items[0].Snippet.Thumbnails.Maxres.URL, Low: live.Items[0].Snippet.Thumbnails.High.URL}
			rz := Newlive{
				Name:        &v.Name,
				ImageID:     &v.ImageID,
				ChannelID:   &live.Items[0].Snippet.ChannelID,
				Title:       &live.Items[0].Snippet.Title,
				Description: &live.Items[0].Snippet.Description,
				Viewers:     &name,
				Likes:       &live.Items[0].Statistics.LikeCount,
				Dislikes:    &live.Items[0].Statistics.DislikeCount,
				VideoID:     &live.Items[0].ID,
				Thumbnail:   thumb,
				Type:        "youtube",
			}
			payload <- rz
		}
		wait.Done()
	}()
	go func() {
		fmt.Println("waiting...")
		wait.Wait()
		fmt.Println("done")
		done <- true
	}()
}

//defer func() {
//	sort.Sort(ByViewers(final))
//	Results = final
// }()

func getTwitch(r Streamer, i int) {
	defer func() {
		if i == len(streamers)-1 {
			wait.Done()
		}
	}()
	url := fmt.Sprintf("https://api.twitch.tv/kraken/streams/%v?client_id=%v", r.Name, os.Getenv("TWITCH"))
	rz, err := http.Get(url)
	if err != nil {
		fmt.Println(err)
	}
	body, err := ioutil.ReadAll(rz.Body)
	defer rz.Body.Close()
	var res TwitchResponse
	json.Unmarshal(body, &res)
	if res.Stream.Channel.Status == nil {
		return
	}
	thumb := Thumbnails{High: *res.Stream.Preview.Large, Low: *res.Stream.Preview.Medium}
	result := Newlive{
		ChannelID: &res.Stream.Channel.Name,
		Name:      &res.Stream.Channel.Name,
		ImageID:   &res.Stream.Channel.Logo,
		VideoID:   &res.Stream.Channel.Name,
		Title:     res.Stream.Channel.Status,
		Viewers:   res.Stream.Viewers,
		Thumbnail: thumb,
		Type:      "twitch",
	}
	payload <- result
}
func Listen() {
	final := []Newlive{}
	for {
		select {
		case request := <-payload:
			final = append(final, request)
		case isDone := <-done:
			if ok := isDone; ok {
				Results = final
				sort.Sort(ByViewers(Results))
				final = nil
				done <- false
			}
		}
	}
}
