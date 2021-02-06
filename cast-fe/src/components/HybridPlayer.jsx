import React from "react";
import "dashjs";
import videojs from "video.js";
// import "videojs-contrib-dash" // must disable for quality selector to appear;
import "videojs-flvjs-es6";
import "videojs-contrib-quality-levels";
import "videojs-http-source-selector";

import "video.js/dist/video-js.css";
import "../styles/player.css";

class HybridPlayer extends React.Component {
  componentDidMount() {
    this.initPlayer();
    this.updatePlayer();
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.url !== prevProps.url) {
      console.log("[HybridPlayer] Update");
      this.updatePlayer();
    }
  }

  componentWillUnmount() {
    if (this.player) {
      console.log("[HybridPlayer] Dismount");
      // this.player.dispose(); // causes SourceBufferSink errors
    }
  }

  initPlayer() {
    console.log("[HybridPlayer] Initialize");
    let options = {
      fluid: true,
      responsive: true,
      aspectRatio: "16:9",
      preload: "false",
      controls: true,
      userActions: { hotkeys: true },
      // liveui: true,
      plugins: {
        httpSourceSelector: {
          default: "auto",
        },
      },
      flvjs: {
        mediaDataSource: {
          isLive: true,
          cors: true,
          withCredentials: false,
        },
      },
    };
    this.player = videojs(this.videoNode, options);
    this.player.qualityLevels();
    this.player.httpSourceSelector();
  }

  updatePlayer() {
    if (!this.props.url) return;
    this.player.pause();
    this.player.reset();
    this.player.src({
      src: this.props.url,
      type: this.props.live ? "video/x-flv" : "application/dash+xml",
    });
    this.player.load();
    this.player.autoplay(this.props.live);
    if (this.props.live) this.player.play();
    this.player.poster(this.props.thumbnail);
  }

  render() {
    return (
      <div>
        <div data-vjs-player style={style.player}>
          <video ref={(ref) => (this.videoNode = ref)} className="video-js" />
        </div>
      </div>
    );
  }
}

let style = {
  player: {
    borderRadius: 8,
  },
};

export default HybridPlayer;
