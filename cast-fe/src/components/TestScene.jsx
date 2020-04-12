import React, {Component} from 'react';
import {
  Button,
  Col,
  Container,
  Modal,
  Row,
} from "react-bootstrap";
import axios from "axios";
import urls from "../helper/url";
import {Redirect} from "react-router-dom";

import 'dashjs'
import videojs from 'video.js'
import 'videojs-contrib-dash'
import 'video.js/dist/video-js.css'
import 'videojs-flvjs-es6'
import 'videojs-contrib-quality-levels'
import 'videojs-http-source-selector'
import './player/player.css'
import HybridPlayer from "./player/HybridPlayer";

class TestScene extends Component {
  constructor(props) {
    super(props);
    this.state = {
      video: null,
      live: {},
      vod: {},
      loading: {
        current: true,
        live: true,
        vod: true,
        comment: false
      },
      liked: false,
      likes: 0,
      comment: "",
      comments: [],
      error_comment: "",
    };
    this.initPlayer = this.initPlayer.bind(this);
    this.incrementView = this.incrementView.bind(this);
    this.handleShare = this.handleShare.bind(this);
    this.handleLike = this.handleLike.bind(this);
    this.handleComment = this.handleComment.bind(this);
    this.writeComment = this.writeComment.bind(this);
    this.handleSubscribe = this.handleSubscribe.bind(this);
    this.handleTip = this.handleTip.bind(this);
    this.promptSignup = this.promptSignup.bind(this);
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.match.params.hash !== prevProps.match.params.hash) {
      this.setState({comment: "", error_comment: "", loading: {...this.state.loading, current: true}});
      this.fetchDetail(this.props.match.params.hash);
    }
  }

  componentDidMount() {
    this.fetchDetail(this.props.match.params.hash);
    this.fetchVideos("live");
    this.fetchVideos("vod");
  }

  incrementView(variant, hash) {
    if (this.props.match.params.hash !== hash && variant === "vod") {
      let updated = this.state[variant][hash];
      updated.views++;
      this.setState({[variant]: {...this.state[variant], [hash]: updated}})
    }
  }

  fetchVideos(variant) {
    axios.get(urls().list(), {
      params: {
        variant: variant,
        count: 8,
        offset: 0,
      }
    }).then((response) => {
      this.setState({loading: {...this.state.loading, [variant]: false}});
      if (response.data.code === 200) {
        this.setState({
          [variant]: response.data.data.reduce((map, obj) => {
            map[obj.hash] = obj;
            return map
          }, {})
        })
      }
    }).catch((error) => {
      console.log(error);
      this.setState({loading: {...this.state.loading, [variant]: false}});
    });
  }

  fetchDetail(hash) {
    axios.get(urls().cast_details(), {
      params: {
        hash: hash,
        username: localStorage.getItem("username")
      }
    }).then((response) => {
      this.setState({loading: {...this.state.loading, current: false}});
      if (response.data.code === 200) {
        let data = response.data.data;
        console.log("DETAIL RECEIVED")
        this.setState({
          video: data, likes: data.likes, liked: data.liked, comments: data.comments,
          [data.type]: {...this.state[data.type], [data.hash]: data}
        });
        document.title = `${data.title} - ${data.author.name} | cast`;
        console.log("DETAIL READY")
        this.initPlayer();
        // if ('mediaSession' in navigator) {
        //   // eslint-disable-next-line no-undef
        //   navigator.mediaSession.metadata = new MediaMetadata({
        //     title: data.title,
        //     artist: data.author.name,
        //     album: 'cast',
        //     artwork: [
        //       {src: urls().thumbnail(this.state.video.hash), sizes: '512x512', type: 'image/jpg'},
        //     ]
        //   });
        // }
      }
    }).catch((error) => {
      console.log(error);
      this.setState({loading: {...this.state.loading, current: false}});
    });
  }

  handleShare() {
    console.log("share video");
  }

  handleLike() {
    if (this.state.loading.current) return;
    if (localStorage.getItem("username")) {
      axios.get(urls().like(), {
        params: {
          hash: this.state.video.hash,
          like: !this.state.liked,
        }
      }).then(() => {
        this.setState({likes: this.state.likes + (this.state.liked ? -1 : 1), liked: !this.state.liked});
      }).catch((error) => {
        console.log(error);
      });
    } else {
      this.promptSignup();
    }
  }

  writeComment(e) {
    this.setState({comment: e.target.value, error_comment: "", error_submit: ""});
  }

  handleComment(e) {
    e.preventDefault();
    if (this.state.loading.current) return;
    if (localStorage.getItem("username")) {
      if (!this.state.comment.trim() || this.state.error_comment) {
        this.setState({error_comment: "Please enter comment"});
        return;
      }
      if (this.state.loading.comment) return;
      this.setState({loading: {...this.state.loading, comment: true}});
      this.setState({error_submit: ""});
      axios.get(urls().comment(), {
        params: {
          hash: this.state.video.hash,
          content: this.state.comment.trim(),
        }
      }).then((response) => {
        this.setState({
          comment: "",
          comments: this.state.comments === null ? [response.data.data] : [response.data.data, ...this.state.comments],
          loading: {...this.state.loading, comment: false}
        });
      }).catch((error) => {
        console.log(error);
        this.setState({error_submit: "Failed submitting comment! Try again later"})
      });
    } else {
      this.promptSignup();
    }
  }

  handleSubscribe() {
    if (this.state.loading.current) return;
    if (localStorage.getItem("username")) {
      // TODO: subscribe
    } else {
      this.promptSignup();
    }
  }

  handleTip() {
    if (this.state.loading.current) return;
    if (localStorage.getItem("username")) {
      // TODO: subscribe
    } else {
      this.promptSignup();
    }
  }

  promptSignup() {
    this.setState({prompt: true})
  }

  initPlayer() {
    console.log("INIT")
    let options = {
      fluid: true,
      responsive: true,
      aspectRatio: "16:9",
      // liveui: true,
      // preload: "auto",
      controls: true,
      src: {
        src: this.state.video.is_live ? urls().live(this.state.video.hash) : urls().vod(this.state.video.hash),
        type: this.state.video.is_live ? 'video/x-flv' : 'application/dash+xml',
      },
      // flvjs: {
      //   mediaDataSource: {
      //     isLive: true,
      //     cors: true,     // TODO: NOTICE!
      //     withCredentials: false,
      //   },
      // },
      // autoplay: this.props.live,
      // poster: this.props.thumbnail,
    };
    this.player = videojs(this.videoNode, options);
    // this.updatePlayer();
    // this.player.httpSourceSelector();
  }

  // updatePlayer() {
  //   console.log("UPDATE PLAYER")
  //   if (!this.props.url) return;
  //   console.log("UPDATE PLAYER SKIPPED")
  //   this.player.pause();
  //   this.player.src({
  //     src: this.props.url,
  //     type: this.props.live ? 'video/x-flv' : 'application/dash+xml',
  //   });
  //   // this.player.autoplay(this.props.live);
  //   // if (this.props.live) this.player.play();
  //   // this.player.load();
  //   // else this.player.pause();
  //   this.player.reset();
  //   this.player.load();
  //   this.player.poster(this.props.thumbnail);
  // }

  render() {
    return (
      <>
        <Container fluid style={style.content_container}>
          <Row>
            <Col xl={{span: 8, order: 2}} sm={{span: 12, order: 1}} xs={{span: 12, order: 1}}>
              <HybridPlayer
                url={this.state.video && (this.state.video.is_live ? urls().live(this.state.video.hash) : urls().vod(this.state.video.hash))}
                thumbnail={this.state.video && urls().thumbnail(this.state.video.hash)}
                live={this.state.video && this.state.video.is_live}/>
              <br/>
              <br/>
              <br/>
              <div data-vjs-player style={style.player}>
                <video ref={node => this.videoNode = node} className="video-js"/>
              </div>
            </Col>
          </Row>
        </Container>
        <Modal show={this.state.prompt} size={"md"} onHide={() => this.setState({prompt: false})} centered>
          <Modal.Header closeButton>
            <Modal.Title id="contained-modal-title-vcenter">
              Join today!
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <p>You need to be logged in to like, comment, subscribe or tip. By signing in, you can start sharing your
              own videos and livestream too!</p>
            <p>Log In or Sign Up today!</p>
          </Modal.Body>
          <Modal.Footer>
            <Button variant={"outline-primary"}
                    onClick={() => this.setState({redirect: <Redirect to={"/login"}/>})}>Log In</Button>
            <Button variant={"primary"}
                    onClick={() => this.setState({redirect: <Redirect to={"/signup"}/>})}>Sign Up</Button>
          </Modal.Footer>
        </Modal>
        {this.state.redirect}
      </>
    );
  }
}

let style = {
  title: {
    margin: 0
  },
  content_container: {
    padding: "36px 0 0 0"
  },
  spinner: {
    display: "block",
    margin: "32px auto 64px auto",
  },
  cast_list: {
    marginTop: 16
  },
  live_chat: {
    borderRadius: "8px 48px 8px 8px",
  },
  live_chat_body: {
    height: 480,
    justifyContent: "flex-end",
    display: "flex",
    flexDirection: "column"
  },
  live_chat_item: {
    marginBottom: 0
  },
  live_chat_input: {
    marginTop: 8
  },
  cast_tag_bar: {
    marginTop: 8,
    display: "flex",
    justifyContent: "space-between"
  },
  cast_tag: {
    background: "white",
    borderRadius: 8,
    borderWidth: 1,
    borderColor: "lightgray",
    borderStyle: "solid",
    fontSize: 16,
    fontWeight: 400,
    marginRight: 8,
    marginBottom: 8,
  },
  cast_attrib: {
    marginLeft: 16
  },
  author_bar: {
    display: "flex",
    justifyContent: "space-between"
  },
  author_profile: {
    display: "flex",
  },
  cast_author_details: {
    display: "flex",
    flexDirection: "column"
  },
  cast_author_name: {
    fontSize: 18,
    fontWeight: 600,
    margin: 0,
  },
  cast_author_sub: {
    fontSize: 16,
    fontWeight: 400,
    margin: 0,
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
  },
  profile_image: {
    marginRight: 8,
    alignSelf: "center"
  },
  tip_button: {
    marginRight: 8
  },
  sub_button: {
    fontWeight: 600
  },
  description: {
    // marginLeft: 48,
    marginTop: 16,
    marginBottom: 16,
    // width: "80%"
  },
  comment_input: {
    borderRadius: "8px 48px 8px 8px",
  },
  comment_list: {
    marginTop: 32,
    marginBottom: 64
  },
  comment_item: {
    marginTop: 16,
    marginBottom: 16,
  },
  clickable: {
    cursor: "pointer"
  },
  liked: {
    color: "#E84409"
  },
  h5: {
    fontFamily: "Open Sans",
    fontSize: 18,
    fontStyle: "italic",
    marginTop: 16,
    textAlign: "center"
  },
};

export default TestScene
