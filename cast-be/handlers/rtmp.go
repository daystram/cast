package handlers

import (
	"errors"
	"fmt"
	"github.com/nareix/joy4/av/avutil"
	"github.com/nareix/joy4/av/pubsub"
	"github.com/nareix/joy4/format/flv"
	"github.com/nareix/joy4/format/rtmp"
	"gitlab.com/daystram/cast/cast-be/config"
	"io"
	"net/http"
	"path"
	"sync"
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
		// TODO: check if user set LIVE in dashboard
		m.live.mutex.Lock()
		ch := m.live.streams[path.Base(conn.URL.Path)]
		if ch == nil {
			fmt.Printf("[RTMPUpLink] Publish request on %s created\n", path.Base(conn.URL.Path))
			stream, _ := conn.Streams()
			ch = &Stream{}
			ch.queue = pubsub.NewQueue()
			_ = ch.queue.WriteHeader(stream)
			m.live.streams[path.Base(conn.URL.Path)] = ch
			m.live.mutex.Unlock()
		} else {
			fmt.Printf("[RTMPUpLink] Publish request on %s already exists\n", path.Base(conn.URL.Path))
			m.live.mutex.Unlock()
			return
		}
		_ = avutil.CopyPackets(ch.queue, conn)
		_ = ch.queue.Close()
	}
	m.live.uplink.Addr = fmt.Sprintf(":%d", config.AppConfig.RTMPPort)
	go m.live.uplink.ListenAndServe()
	fmt.Printf("[CreateRTMPUpLink] RTMP UpLink created at port %d\n", config.AppConfig.RTMPPort)
}

func (m *module) StreamLive(_ string, w http.ResponseWriter, r *http.Request) (err error) {
	m.live.mutex.RLock()
	ch := m.live.streams[path.Base(r.URL.Path)]
	m.live.mutex.RUnlock()

	if ch != nil {
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
		return errors.New(fmt.Sprintf("[StreamLive] Unable to find stream for %s", path.Base(r.URL.Path)))
	}
	return
}
