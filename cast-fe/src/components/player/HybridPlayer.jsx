import React from 'react';
import videojs from 'video.js'
import 'video.js/dist/video-js.css'
import 'videojs-flvjs-es6'
import 'videojs-contrib-quality-levels'
import 'videojs-http-source-selector'
import './player.css'

class HybridPlayer extends React.Component {
  componentDidMount() {
    this.initPlayer()
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.url !== prevProps.url) {
      this.updatePlayer();
    }
  }

  componentWillUnmount() {
    if (this.player) {
      this.player.dispose()
    }
  }

  initPlayer() {
    let options = {
      fluid: true,
      responsive: true,
      aspectRatio: "16:9",
      // liveui: true,
      preload: true,
      controls: true,
      flvjs: {
        mediaDataSource: {
          isLive: true,
          cors: true,     // TODO: NOTICE!
          withCredentials: false,
        },
      },
      smoothQualityChange: true,
      autoplay: this.props.live,
      poster: this.props.thumbnail,
    };
    this.player = videojs(this.videoNode, options);
    this.player.httpSourceSelector();
  }

  updatePlayer() {
    this.player.pause();
    if (!this.props.url) return;
    this.player.poster(this.props.thumbnail);
    this.player.src({
      src: this.props.url,
      type: this.props.live ? 'video/x-flv' : 'application/dash+xml',
    });
    this.player.autoplay(this.props.live);
    if (this.props.live) this.player.play();
    else this.player.pause();
    this.player.load();
  }

  render() {
    return (
      <div>
        <div data-vjs-player style={style.player}>
          <video ref={node => this.videoNode = node} className="video-js"/>
        </div>
      </div>
    )
  }
}

let style = {
  player: {
    borderRadius: 8
  },
};

export default HybridPlayer
