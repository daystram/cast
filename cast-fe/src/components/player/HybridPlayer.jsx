import React from 'react';
import videojs from 'video.js'
import 'video.js/dist/video-js.css'
import 'videojs-flvjs-es6'
import 'videojs-contrib-quality-levels'
import 'videojs-http-source-selector'

export default class HybridPlayer extends React.Component {
  componentDidMount() {
    let options = {
      autoplay: this.props.live || false,
      controls: true,
      sources: [{
        src: this.props.url,
        type: this.props.live ? 'video/x-flv' : 'application/dash+xml',
      }],
      flvjs: {
        mediaDataSource: {
          isLive: true,
          cors: true,
          withCredentials: false,
        },
      },
      smoothQualityChange: true
    };
    this.player = videojs(this.videoNode, options);
    this.player.httpSourceSelector();
  }

  componentWillUnmount() {
    if (this.player) {
      this.player.dispose()
    }
  }

  render() {
    return (
      <div>
        <div data-vjs-player>
          <video ref={node => this.videoNode = node} className="video-js"/>
        </div>
      </div>
    )
  }
}