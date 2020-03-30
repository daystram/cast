import React, {Component} from 'react';
import {
  Alert,
  Badge,
  Button,
  Card,
  Col,
  Container,
  Form,
  FormControl,
  Image,
  InputGroup,
  Modal,
  Row,
  Spinner
} from "react-bootstrap";
import Cast from "./Cast"
import Sidebar from "./Sidebar";
import {HybridPlayer} from "./player";
import abbreviate from "../helper/abbreviate";
import axios from "axios";
import urls from "../helper/url";
import format from "../helper/format";
import {Redirect} from "react-router-dom";

class Scene extends Component {
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
        this.setState({
          video: data, likes: data.likes, liked: data.liked, comments: data.comments,
          [data.type]: {...this.state[data.type], [data.hash]: data}
        });
        document.title = `${data.title} - ${data.author.name} | cast`;
        console.log("DETAIL READY")
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

  render() {
    return (
      <>
        <Container fluid style={style.content_container}>
          <Row>
            <Col xl={{span: 2, order: 1}} sm={{span: 6, order: 2}} xs={{span: 12, order: 2}}>
              <Sidebar/>
              <div style={style.cast_list}>
                {this.state.live && Object.values(this.state.live).map(video =>
                  <Row noGutters style={{padding: "0 0 16px 0"}}>
                    <Cast video={video} onClick={(a, b) => this.incrementView(a, b)}/>
                  </Row>
                )}
                {this.state.loading.live && <Spinner style={style.spinner} animation="grow" variant="primary"/>}
              </div>
            </Col>
            <Col xl={{span: 8, order: 2}} sm={{span: 12, order: 1}} xs={{span: 12, order: 1}}>
              <HybridPlayer
                url={this.state.video && (this.state.video.is_live ? urls().live(this.state.video.hash) : urls().vod(this.state.video.hash))}
                thumbnail={this.state.video && urls().thumbnail(this.state.video.hash)}
                live={this.state.video && this.state.video.is_live}/>
              <div style={style.cast_tag_bar}>
                <div>
                  {this.state.video && this.state.video.tags && Object.values(this.state.video.tags).map(tag =>
                    <Badge pill style={style.cast_tag}>{tag}</Badge>
                  )}
                </div>
                <div>
                  <span style={{...style.cast_attrib, ...style.clickable}} onClick={this.handleShare}>
                    <i className="material-icons">share</i>{" "}share</span>
                  <span style={{...style.cast_attrib, ...style.clickable}} onClick={this.handleLike}>
                    <i style={this.state.liked ? style.liked : {}} className="material-icons">thumb_up</i>
                    {" "}{(this.state.video && abbreviate().number(this.state.likes)) || 0} likes
                  </span>
                  <span style={style.cast_attrib}>
                    <i className="material-icons">remove_red_eye</i>
                    {" "}{(this.state.video && abbreviate().number(this.state.video.views)) || 0} {this.state.video && (this.state.video.is_live ? 'viewers' : 'views')}
                  </span>
                </div>
              </div>
              <h1 style={style.title}>{this.state.video && this.state.video.title}</h1>
              <p style={{marginTop: 4}}>{this.state.video && format().date(this.state.video.created_at)}</p>
              <div style={style.author_bar}>
                <div style={style.author_profile}>
                  <Image src={this.state.video && urls().profile(this.state.video.author.username)}
                         width={42} height={42} style={style.profile_image} roundedCircle/>
                  <div style={style.cast_author_details}>
                    <p style={style.cast_author_name}>{this.state.video && this.state.video.author.name}</p>
                    <p style={style.cast_author_sub}>
                      {(this.state.video && abbreviate().number(this.state.video.author.subscribers)) || 0} subscribers
                    </p>
                  </div>
                </div>
                <div>
                  <Button style={style.tip_button} onClick={this.handleTip}><i
                    className="material-icons">attach_money</i></Button>
                  <Button style={style.sub_button} onClick={this.handleSubscribe}
                          disabled={this.state.video && this.state.video.author.isSubscribed}>SUBSCRIBE</Button>
                </div>
              </div>
              <Row noGutters>
                <Col xl={1} sm={0}/>
                <Col>
                  <div style={style.description}>{this.state.video && this.state.video.description}</div>
                </Col>
                <Col xl={1} sm={0}/>
              </Row>
              <hr/>
              <h3>Comments</h3>
              <Row noGutters style={{marginTop: 28}}>
                <Col xl={1} sm={0}/>
                <Col>
                  <Form noValidate onSubmit={this.handleComment}>
                    {this.state.error_submit && <Alert variant={"danger"}>{this.state.error_submit}</Alert>}
                    <Form.Group>
                      <InputGroup style={style.comment_input}>
                        <Form.Control type="text" placeholder="Comment" value={this.state.comment}
                                      onChange={this.writeComment} isInvalid={!!this.state.error_comment}/>
                        <InputGroup.Append>
                          <Button variant="outline-primary" type="submit">
                            <i className="material-icons">send</i></Button>
                        </InputGroup.Append>
                      </InputGroup>
                    </Form.Group>
                  </Form>
                  <div style={style.comment_list}>
                    {this.state.comments ? Object.values(this.state.comments).map(comment => {
                      let duration = abbreviate().time(Date.now() - new Date(comment.created_at));
                      return (
                        <div style={{...style.author_profile, ...style.comment_item}}>
                          <Image src={urls().profile(comment.author.username)} height={42} width={42}
                                 style={{...style.profile_image, alignSelf: "end"}} roundedCircle/>
                          <div style={{...style.cast_author_details, minWidth: 0}}>
                            <p style={style.cast_author_name}>{comment.author.name}</p>
                            <p style={{marginBottom: 0, color: "grey"}}>
                              {`${duration} ${duration === "Yesterday" ? "" : "ago"}`}
                            </p>
                            <p style={{...style.cast_author_sub, whiteSpace: "normal"}}>{comment.content}</p>
                          </div>
                        </div>
                      )
                    }) : <h5 style={style.h5}>Post the first comment!</h5>}
                  </div>
                </Col>
                <Col xl={1} sm={0}/>
              </Row>
            </Col>
            <Col xl={{span: 2, order: 3}} sm={{span: 6, order: 3}} xs={{span: 12, order: 3}}>
              <Card style={style.live_chat}>
                <Card.Body style={style.live_chat_body}>
                  <p style={style.live_chat_item}><b>User1</b>: Hello! Hello! Hello! Hello! Hello! Hello! Hello!</p>
                  <p style={style.live_chat_item}><b>User2</b>: Hello!</p>
                  <p style={style.live_chat_item}><b>User1</b>: Hello!</p>
                  <p style={style.live_chat_item}><b>User2</b>: Hello!</p>
                  <p style={style.live_chat_item}><b>User1</b>: Hello!</p>
                  <InputGroup style={style.live_chat_input}>
                    <FormControl type="text" placeholder="Chat"/>
                    <InputGroup.Append>
                      <Button variant="outline-primary"><i className="material-icons">send</i></Button>
                    </InputGroup.Append>
                  </InputGroup>
                </Card.Body>
              </Card>
              <div style={style.cast_list}>
                {this.state.vod && Object.values(this.state.vod).map(video =>
                  <Row noGutters style={{padding: "0 0 16px 0"}}>
                    <Cast video={video} onClick={(a, b) => this.incrementView(a, b)}/>
                  </Row>
                )}
                {this.state.loading.vod && <Spinner style={style.spinner} animation="grow" variant="primary"/>}
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

export default Scene
