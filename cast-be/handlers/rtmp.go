package handlers

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"path"
	"sync"
	"time"

	"github.com/daystram/cast/cast-be/config"
	"github.com/daystram/cast/cast-be/datatransfers"

	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/flv"
	"github.com/nareix/joy4/format/rtmp"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type writeFlusher struct {
	httpflusher http.Flusher
	io.Writer
}

func (w writeFlusher) Flush() error {
	w.httpflusher.Flush()
	return nil
}

// Adapted from https://github.com/nareix/joy4/blob/master/examples/http_flv_and_rtmp_server/main.go
func (m *module) CreateRTMPUpLink() {
	m.live = Live{
		streams: map[string]*Stream{},
		uplink:  &rtmp.Server{},
		mutex:   &sync.RWMutex{},
	}
	m.live.uplink.HandlePublish = func(conn *rtmp.Conn) {
		username := path.Base(conn.URL.Path)
		var video datatransfers.Video
		var err error
		if video, err = m.db.videoOrm.GetOneByHash(username); err != nil {
			fmt.Printf("[RTMPUpLink] Unknown username %s\n", username)
			_ = conn.Close()
			return
		}
		if !video.Pending {
			fmt.Printf("[RTMPUpLink] Window for %s not PENDING\n", username)
			_ = conn.Close()
			return
		}
		m.live.mutex.Lock()
		ch := m.live.streams[path.Base(conn.URL.Path)]
		if ch == nil {
			stream, _ := conn.Streams()
			ch = &Stream{}
			ch.queue = pubsub.NewQueue()
			_ = ch.queue.WriteHeader(stream)
			m.live.streams[path.Base(conn.URL.Path)] = ch
			if err = m.db.videoOrm.SetLive(video.Author.ID, false, true); err != nil {
				fmt.Printf("[RTMPUpLink] Failed setting stream for %s to LIVE. %+v\n", username, err)
				delete(m.live.streams, username)
				return
			}
			fmt.Printf("[RTMPUpLink] UpLink for %s connected\n", username)
			m.BroadcastNotificationSubscriber(video.Author.ID, datatransfers.NotificationOutgoing{
				Message:   fmt.Sprintf("%s just went live! Watch now!", video.Author.Name),
				Username:  video.Author.Username,
				Hash:      video.Hash,
				CreatedAt: time.Now(),
			})
			m.live.mutex.Unlock()
		} else {
			fmt.Printf("[RTMPUpLink] UpLink for %s already exists\n", username)
			m.live.mutex.Unlock()
			_ = ch.queue.Close()
			return
		}
		_ = avutil.CopyPackets(ch.queue, conn)
		_ = ch.queue.Close()
		delete(m.live.streams, username)
		fmt.Printf("[RTMPUpLink] UpLink for %s disconnected. Stopping stream...\n", username)
		if err = m.db.videoOrm.SetLive(video.Author.ID, false, false); err != nil {
			fmt.Printf("[RTMPUpLink] Failed setting stream for %s to STOP. %+v\n", username, err)
		}
	}
	m.live.uplink.Addr = fmt.Sprintf(":%d", config.AppConfig.RTMPPort)
	go m.live.uplink.ListenAndServe()
	fmt.Printf("[CreateRTMPUpLink] RTMP UpLink Window opened at port %d\n", config.AppConfig.RTMPPort)
}

func (m *module) ControlUpLinkWindow(userID primitive.ObjectID, open bool) (err error) {
	var stream datatransfers.Video
	if stream, err = m.db.videoOrm.GetLiveByAuthor(userID); err != nil {
		return errors.New(fmt.Sprintf("[ControlUpLinkWindow] failed retrieving video by %s. %+v", userID.Hex(), err))
	}
	if (stream.IsLive && open) || (!stream.IsLive && !stream.Pending && !open) {
		return errors.New(fmt.Sprintf("[ControlUpLinkWindow] stream window already set for %s", stream.Hash))
	}
	delete(m.live.streams, stream.Hash)
	return m.db.videoOrm.SetLive(userID, open, false)
}

func (m *module) StreamLive(_ string, w http.ResponseWriter, r *http.Request) (err error) {
	m.live.mutex.RLock()
	ch := m.live.streams[path.Base(r.URL.Path)]
	m.live.mutex.RUnlock()

	if ch != nil {
		fmt.Printf("[StreamLive] Streaming request for %s\n", path.Base(r.URL.Path))
		w.Header().Set("Content-Type", "video/x-flv")
		w.Header().Set("Transfer-Encoding", "chunked")
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.WriteHeader(200)
		flusher := w.(http.Flusher)
		flusher.Flush()
		_ = avutil.CopyFile(
			flv.NewMuxerWriteFlusher(writeFlusher{httpflusher: flusher, Writer: w}),
			ch.queue.Latest(),
		)
	} else {
		fmt.Printf("[StreamLive] Unable to find stream for %s\n", path.Base(r.URL.Path))
		return errors.New(fmt.Sprintf("[StreamLive] Unable to find stream for %s", path.Base(r.URL.Path)))
	}
	return
}
