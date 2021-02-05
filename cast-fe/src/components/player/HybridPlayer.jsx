import React from "react";
import "dashjs";
import videojs from "video.js";
import "videojs-contrib-quality-levels";
import "videojs-http-source-selector";
import "videojs-contrib-dash";
import "video.js/dist/video-js.css";
import "videojs-flvjs-es6";
import "./player.css";

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
      this.player.dispose();
    }
  }

  initPlayer() {
    console.log("[HybridPlayer] Initialize");
    let options = {
      fluid: true,
      responsive: true,
      aspectRatio: "16:9",
      // liveui: true,
      preload: "true",
      controls: true,
      userActions: { hotkeys: true },
      plugins: {
        httpSourceSelector: {
          default: "auto",
        },
      },
      flvjs: {
        mediaDataSource: {
          isLive: true,
          cors: true, // TODO: NOTICE!
          withCredentials: false,
        },
      },
      // autoplay: this.props.live,
      // poster: this.props.thumbnail,
    };
    this.player = videojs(this.videoNode, options);
    this.player.httpSourceSelector();
    this.player.load();
  }

  updatePlayer() {
    console.log(this.props.url);
    if (!this.props.url) return;
    this.player.pause();
    this.player.src({
      src: this.props.url,
      type: this.props.live ? "video/x-flv" : "application/dash+xml",
    });
    this.player.autoplay(this.props.live);
    if (this.props.live) this.player.play();
    // else this.player.pause();
    this.player.load();
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
