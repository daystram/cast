import React, { Component } from "react";
import {
  Alert,
  Badge,
  Button,
  Card,
  Col,
  Container,
  Form,
  InputGroup,
  Modal,
  OverlayTrigger,
  Popover,
  Row,
  Spinner,
} from "react-bootstrap";
import { Cast, HybridPlayer, Sidebar } from "../components";
import abbreviate from "../helper/abbreviate";
import format from "../helper/format";
import { withRouter } from "react-router-dom";
import { authManager } from "../helper/auth";
import TimeAgo from "react-timeago";
import queryString from "query-string";
import Chat from "./Chat";
import MediaQuery from "react-responsive";
import { MOBILE_BP } from "../constants/breakpoint";
import { VIDEO_TYPE_LIVE, VIDEO_TYPE_VOD } from "../constants/video";
import logo from "../components/logo.svg";
import api from "../apis/api";

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
        comment: false,
      },
      liked: false,
      subscribed: false,
      likes: 0,
      comment: "",
      comments: [],
      error_comment: "",
      not_found: false,
    };
    this.handleDownload = this.handleDownload.bind(this);
    this.handleShare = this.handleShare.bind(this);
    this.handleLike = this.handleLike.bind(this);
    this.handleComment = this.handleComment.bind(this);
    this.writeComment = this.writeComment.bind(this);
    this.handleSubscribe = this.handleSubscribe.bind(this);
    this.promptSignup = this.promptSignup.bind(this);
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.match.params.hash !== prevProps.match.params.hash) {
      this.setState({
        comment: "",
        error_comment: "",
        loading: { ...this.state.loading, current: true },
      });
      this.fetchDetail(this.props.match.params.hash);
      window.scroll({ top: 0, left: 0, behavior: "smooth" });
    }
  }

  componentDidMount() {
    this.fetchDetail(this.props.match.params.hash);
    window.scroll({ top: 0, left: 0, behavior: "smooth" });
    this.fetchVideos(VIDEO_TYPE_LIVE);
    this.fetchVideos(VIDEO_TYPE_VOD);
  }

  incrementView(variant, hash) {
    if (this.props.match.params.hash !== hash && variant === VIDEO_TYPE_VOD) {
      let updated = this.state[variant][hash];
      updated.views++;
      this.setState({ [variant]: { ...this.state[variant], [hash]: updated } });
    }
  }

  fetchVideos(variant) {
    api.cast
      .list({
        variant: variant,
        count: 8,
        offset: 0,
      })
      .then((response) => {
        this.setState({ loading: { ...this.state.loading, [variant]: false } });
        if (response.data.code === 200) {
          this.setState({
            [variant]: response.data.data.reduce((map, obj) => {
              map[obj.hash] = obj;
              return map;
            }, {}),
          });
        }
      })
      .catch((error) => {
        console.log(error);
        this.setState({ loading: { ...this.state.loading, [variant]: false } });
      });
  }

  fetchDetail(hash) {
    api.cast
      .detail({
        hash: hash,
        username: authManager.getUser().preferred_username,
      })
      .then((response) => {
        this.setState({ loading: { ...this.state.loading, current: false } });
        if (response.data.code === 200) {
          let data = response.data.data;
          this.setState({
            offline: data.type === VIDEO_TYPE_LIVE && !data.is_live,
            not_found: false,
            video: data,
            likes: data.likes,
            liked: data.liked,
            subscribed: data.subscribed,
            comments: data.comments,
            [data.type]: { ...this.state[data.type], [data.hash]: data },
          });
          document.title = `${data.title} - ${data.author.name} | cast`;
          if ("mediaSession" in navigator) {
            // eslint-disable-next-line no-undef
            navigator.mediaSession.metadata = new MediaMetadata({
              title: data.title,
              artist: data.author.name,
              album: "cast",
              artwork: [
                {
                  src: api.cdn.thumbnail(this.state.video.hash),
                  sizes: "512x512",
                  type: "image/jpg",
                },
              ],
            });
          }
        } else {
          this.setState({ not_found: true });
        }
      })
      .catch((error) => {
        console.log(error);
        this.setState({ loading: { ...this.state.loading, current: false } });
      });
  }

  handleDownload() {
    if (this.state.loading.current) return;
    if (authManager.isAuthenticated()) {
      let link = document.createElement("a");
      link.href = api.cdn.download(this.state.video.hash);
      link.download = `${this.state.video.title} by ${this.state.video.author.name} - cast`;
      document.body.appendChild(link);
      link.click();
      document.body.removeChild(link);
    } else {
      this.promptSignup();
    }
  }

  handleShare() {
    if (this.state.loading.current) return;
    this.setState({ prompt_share: true });
  }

  handleLike() {
    if (this.state.loading.current) return;
    if (authManager.isAuthenticated()) {
      api.cast
        .like({
          hash: this.state.video.hash,
          like: !this.state.liked,
        })
        .then(() => {
          this.setState({
            likes: this.state.likes + (this.state.liked ? -1 : 1),
            liked: !this.state.liked,
          });
        })
        .catch((error) => {
          console.log(error);
        });
    } else {
      this.promptSignup();
    }
  }

  writeComment(e) {
    this.setState({
      comment: e.target.value,
      error_comment: "",
      error_submit: "",
    });
  }

  handleComment(e) {
    e.preventDefault();
    if (this.state.loading.current) return;
    if (authManager.isAuthenticated()) {
      if (!this.state.comment.trim() || this.state.error_comment) {
        this.setState({ error_comment: "Please enter your comment" });
        return;
      }
      if (this.state.loading.comment) return;
      this.setState({ loading: { ...this.state.loading, comment: true } });
      this.setState({ error_submit: "" });
      api.cast
        .comment({
          hash: this.state.video.hash,
          content: this.state.comment.trim(),
        })
        .then((response) => {
          this.setState({
            comment: "",
            comments:
              this.state.comments === null
                ? [response.data.data]
                : [response.data.data, ...this.state.comments],
            loading: { ...this.state.loading, comment: false },
          });
        })
        .catch((error) => {
          console.log(error);
          this.setState({
            error_submit: "Failed submitting comment! Try again later",
          });
        });
    } else {
      this.promptSignup();
    }
  }

  handleSubscribe() {
    if (this.state.loading.current) return;
    if (authManager.isAuthenticated()) {
      api.user
        .subscribe({
          author: this.state.video.author.username,
          subscribe: !this.state.subscribed,
        })
        .then(() => {
          this.setState({
            subscribed: !this.state.subscribed,
            loading: { ...this.state.loading, subscribe: false },
          });
        })
        .catch((error) => {
          console.log(error);
        });
    } else {
      this.promptSignup();
    }
  }

  promptSignup() {
    this.setState({ prompt_auth: true });
  }

  render() {
    if (this.state.loading.live) return <></>;
    return (
      <>
        {!this.state.not_found && !this.state.offline && (
          <Container fluid style={style.content_container}>
            <Row>
              <Col
                xl={{ span: 2, order: 1 }}
                sm={{ span: 6, order: 2 }}
                xs={{ span: 12, order: 2 }}
              >
                <MediaQuery minDeviceWidth={MOBILE_BP}>
                  <Card body style={{ borderRadius: "8px 48px 8px 8px" }}>
                    <Sidebar />
                  </Card>
                </MediaQuery>
                <div style={style.cast_list}>
                  {this.state.live &&
                    Object.values(this.state.live).map((video) => (
                      <Row
                        key={video.hash}
                        noGutters
                        style={{ padding: "0 0 16px 0" }}
                      >
                        <Cast
                          video={video}
                          onClick={(a, b) => this.incrementView(a, b)}
                        />
                      </Row>
                    ))}
                  {this.state.loading.live && (
                    <Spinner
                      style={style.spinner}
                      animation="grow"
                      variant="primary"
                    />
                  )}
                </div>
              </Col>
              <Col
                xl={{ span: 8, order: 2 }}
                sm={{ span: 12, order: 1 }}
                xs={{ span: 12, order: 1 }}
                className={"mid-container"}
              >
                <HybridPlayer
                  url={
                    this.state.video &&
                    (this.state.video.type === VIDEO_TYPE_LIVE
                      ? api.live.stream(this.state.video.hash)
                      : api.cdn.vod(this.state.video.hash))
                  }
                  thumbnail={
                    this.state.video && api.cdn.thumbnail(this.state.video.hash)
                  }
                  live={
                    this.state.video &&
                    this.state.video.type === VIDEO_TYPE_LIVE
                  }
                />
                <Row noGutters style={style.cast_tag_bar}>
                  <Col md={true}>
                    {this.state.video &&
                      this.state.video.tags &&
                      Object.values(this.state.video.tags).map((tag) => (
                        <Badge pill key={tag} style={style.cast_tag}>
                          {tag}
                        </Badge>
                      ))}
                  </Col>
                  <Col
                    md={true}
                    style={{ display: "flex", justifyContent: "flex-end" }}
                  >
                    {this.state.video &&
                      this.state.video.type === VIDEO_TYPE_VOD && (
                        <span
                          style={{ ...style.cast_attrib, ...style.clickable }}
                          onClick={this.handleDownload}
                        >
                          <i className="material-icons">get_app</i> download
                        </span>
                      )}
                    <span
                      style={{ ...style.cast_attrib, ...style.clickable }}
                      onClick={this.handleShare}
                    >
                      <i className="material-icons">share</i> share
                    </span>
                    <span
                      style={{ ...style.cast_attrib, ...style.clickable }}
                      onClick={this.handleLike}
                    >
                      <i
                        style={this.state.liked ? style.liked : {}}
                        className="material-icons"
                      >
                        thumb_up
                      </i>{" "}
                      {(this.state.video &&
                        abbreviate().number(this.state.likes)) ||
                        0}{" "}
                      likes
                    </span>
                    <span style={style.cast_attrib}>
                      <i className="material-icons">remove_red_eye</i>{" "}
                      {(this.state.video &&
                        abbreviate().number(this.state.video.views)) ||
                        0}{" "}
                      {this.state.video &&
                        (this.state.video.type === VIDEO_TYPE_LIVE
                          ? "viewers"
                          : "views")}
                    </span>
                  </Col>
                </Row>
                <h1 style={style.title}>
                  {this.state.video && this.state.video.title}
                </h1>
                <p style={{ marginTop: 4 }}>
                  {this.state.video &&
                    format().date(this.state.video.created_at)}
                </p>
                <div style={style.author_bar}>
                  <div style={style.author_profile}>
                    <div
                      className="text-center"
                      style={{
                        width: 42,
                        height: 42,
                        flexShrink: 0,
                        borderRadius: 21,
                        background: "gray",
                        color: "white",
                        fontSize: "21px",
                        lineHeight: "42px",
                        textAlign: "center",
                        textTransform: "capitalize",
                        ...style.profile_image,
                      }}
                    >
                      {this.state.video.author.username[0]}
                    </div>
                    <div style={style.cast_author_details}>
                      <p style={style.cast_author_name}>
                        {this.state.video && this.state.video.author.name}
                      </p>
                      <p style={style.cast_author_sub}>
                        {(this.state.video &&
                          abbreviate().number(
                            this.state.video.author.subscribers
                          )) ||
                          0}{" "}
                        subscriber
                        {this.state.video &&
                          this.state.video.author.subscribers !== 1 &&
                          "s"}
                      </p>
                    </div>
                  </div>
                  <div>
                    <Button
                      style={style.sub_button}
                      onClick={this.handleSubscribe}
                      variant={
                        this.state.subscribed ? "outline-primary" : "primary"
                      }
                    >
                      {this.state.subscribed ? "SUBSCRIBED" : "SUBSCRIBE"}
                    </Button>
                  </div>
                </div>
                <Row className="justify-content-center" noGutters>
                  <Col xl={10} xs={12}>
                    <div style={style.description}>
                      {this.state.video && this.state.video.description}
                    </div>
                  </Col>
                </Row>
                <hr />
                <h3>Comments</h3>
                <Row
                  className="justify-content-center"
                  noGutters
                  style={{ marginTop: 28 }}
                >
                  <Col xl={10} xs={12}>
                    <Form noValidate onSubmit={this.handleComment}>
                      {this.state.error_submit && (
                        <Alert variant={"danger"}>
                          {this.state.error_submit}
                        </Alert>
                      )}
                      <Form.Group>
                        <InputGroup style={style.comment_input}>
                          <Form.Control
                            type="text"
                            placeholder="Comment"
                            value={this.state.comment}
                            onChange={this.writeComment}
                            isInvalid={!!this.state.error_comment}
                          />
                          <InputGroup.Append>
                            <Button variant="outline-primary" type="submit">
                              <i className="material-icons">send</i>
                            </Button>
                          </InputGroup.Append>
                        </InputGroup>
                      </Form.Group>
                    </Form>
                    <div style={style.comment_list}>
                      {this.state.comments ? (
                        Object.values(this.state.comments).map((comment) => {
                          return (
                            <div
                              key={comment.created_at}
                              style={{
                                ...style.author_profile,
                                ...style.comment_item,
                              }}
                            >
                              <div
                                className="text-center"
                                style={{
                                  width: 42,
                                  height: 42,
                                  flexShrink: 0,
                                  borderRadius: 21,
                                  background: "gray",
                                  color: "white",
                                  fontSize: "21px",
                                  lineHeight: "42px",
                                  textAlign: "center",
                                  textTransform: "capitalize",
                                  ...style.profile_image,
                                  alignSelf: "end",
                                }}
                              >
                                {comment.author.name[0]}
                              </div>
                              <div
                                style={{
                                  ...style.cast_author_details,
                                  minWidth: 0,
                                }}
                              >
                                <p style={style.cast_author_name}>
                                  {comment.author.name}
                                </p>
                                <p style={{ marginBottom: 0, color: "grey" }}>
                                  <TimeAgo date={comment.created_at} />
                                </p>
                                <p
                                  style={{
                                    ...style.cast_author_sub,
                                    whiteSpace: "normal",
                                  }}
                                >
                                  {comment.content}
                                </p>
                              </div>
                            </div>
                          );
                        })
                      ) : (
                        <h5 style={style.h5}>Post the first comment!</h5>
                      )}
                    </div>
                  </Col>
                </Row>
                <MediaQuery maxDeviceWidth={MOBILE_BP}>
                  <hr />
                </MediaQuery>
              </Col>
              <Col
                xl={{ span: 2, order: 3 }}
                sm={{ span: 6, order: 3 }}
                xs={{ span: 12, order: 3 }}
              >
                {this.state.video && (
                  <Chat
                    height={480}
                    embedded={true}
                    hash={this.state.video.hash}
                    promptSignup={this.promptSignup}
                  />
                )}
                <div style={style.cast_list}>
                  {this.state.vod &&
                    Object.values(this.state.vod).map((video) => (
                      <Row
                        key={video.hash}
                        noGutters
                        style={{ padding: "0 0 16px 0" }}
                      >
                        <Cast
                          video={video}
                          onClick={(a, b) => this.incrementView(a, b)}
                        />
                      </Row>
                    ))}
                  {this.state.loading.vod && (
                    <Spinner
                      style={style.spinner}
                      animation="grow"
                      variant="primary"
                    />
                  )}
                </div>
              </Col>
            </Row>
          </Container>
        )}
        {this.state.not_found && (
          <>
            <img
              src={logo}
              height={"120"}
              alt={"Video Unavailable"}
              style={{
                display: "block",
                margin: "auto",
                opacity: 0.25,
                marginTop: "25vh",
              }}
            />
            <h1
              style={{
                fontFamily: "Comfortaa",
                textAlign: "center",
                marginTop: 18,
                opacity: 0.85,
              }}
            >
              Video Unavailable
            </h1>
          </>
        )}
        {this.state.offline && (
          <>
            <img
              src={logo}
              height={"120"}
              alt={"Stream Offline"}
              style={{
                display: "block",
                margin: "auto",
                opacity: 0.25,
                marginTop: "25vh",
              }}
            />
            <h1
              style={{
                fontFamily: "Comfortaa",
                textAlign: "center",
                marginTop: 18,
                opacity: 0.85,
              }}
            >
              Stream Offline
            </h1>
          </>
        )}
        <Modal
          show={this.state.prompt_auth}
          size={"md"}
          onHide={() => this.setState({ prompt_auth: false })}
          centered
        >
          <Modal.Header closeButton>
            <Modal.Title id="contained-modal-title-vcenter">
              Join today!
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <p>
              You need to be logged in to like, comment, chat, download, and
              subscribe. By signing in, you can start sharing your own videos
              and livestream too!
            </p>
            <p>Log In or Sign Up today!</p>
          </Modal.Body>
          <Modal.Footer>
            <Button
              variant={"outline-primary"}
              onClick={() => this.props.history.push("/login")}
            >
              Log In
            </Button>
            <Button
              variant={"primary"}
              onClick={() => this.props.history.push("/signup")}
            >
              Sign Up
            </Button>
          </Modal.Footer>
        </Modal>
        <Modal
          show={this.state.prompt_share}
          size={"md"}
          onHide={() => this.setState({ prompt_share: false })}
          centered
        >
          <Modal.Header closeButton>
            <Modal.Title id="contained-modal-title-vcenter">
              Share Video
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <div
              style={{
                display: "flex",
                justifyContent: "space-evenly",
                marginBottom: 16,
              }}
            >
              <Button
                className={"share-twitter"}
                onClick={() =>
                  window.open(
                    "https://twitter.com/intent/tweet?" +
                      queryString.stringify({
                        text: `Watch ${
                          this.state.video && this.state.video.title
                        } by ${
                          this.state.video && this.state.video.author.username
                        } at cast! ${window.location.href.split("?")[0]}`,
                      }),
                    "Share",
                    `width=600,height=400,left=${
                      (window.outerWidth - 600) / 2
                    },top=${(window.outerHeight - 400) / 2}`
                  )
                }
              >
                <span className={"fa fa-twitter"} /> Tweet
              </Button>
              <Button
                className={"share-facebook"}
                onClick={() =>
                  window.open(
                    "https://www.facebook.com/sharer/sharer.php?display=popup&" +
                      queryString.stringify({
                        u: window.location.href.split("?")[0],
                      }),
                    "Share",
                    `width=600,height=400,left=${
                      (window.outerWidth - 600) / 2
                    },top=${(window.outerHeight - 400) / 2}`
                  )
                }
              >
                <span className={"fa fa-facebook"} /> Share
              </Button>
              <a
                href={`mailto:?body=Watch ${
                  this.state.video && this.state.video.title
                } by ${
                  this.state.video && this.state.video.author.name
                } at cast! ${window.location.href.split("?")[0]}`}
              >
                <Button>
                  <span className={"fa fa-envelope"} /> Email
                </Button>
              </a>
            </div>
            <Form.Group>
              <InputGroup>
                <Form.Control
                  type="text"
                  value={window.location.href.split("?")[0]}
                  ref={(ref) => (this.shareField = ref)}
                />
                <InputGroup.Append>
                  <OverlayTrigger
                    trigger="click"
                    placement="top"
                    overlay={
                      <Popover id="popover-basic">
                        <Popover.Content>Copied!</Popover.Content>
                      </Popover>
                    }
                  >
                    <Button
                      variant="outline-primary"
                      onClick={() => {
                        this.shareField.select();
                        document.execCommand("copy");
                        this.shareField.blur();
                      }}
                    >
                      Copy
                    </Button>
                  </OverlayTrigger>
                </InputGroup.Append>
              </InputGroup>
            </Form.Group>
          </Modal.Body>
        </Modal>
      </>
    );
  }
}

let style = {
  title: {
    color: "#EEE",
    margin: 0,
  },
  content_container: {},
  spinner: {
    display: "block",
    margin: "32px auto 64px auto",
  },
  cast_list: {
    marginTop: 16,
  },
  cast_tag_bar: {
    marginTop: 8,
  },
  cast_tag: {
    background: "#8B2803AA",
    color: "#DDD",
    borderRadius: 8,
    // borderWidth: 1,
    // borderColor: "lightgray",
    // borderStyle: "solid",
    fontSize: 16,
    fontWeight: 400,
    marginRight: 8,
    marginBottom: 8,
  },
  cast_attrib: {
    color: "#DDD",
    marginLeft: 16,
    marginBottom: 8,
  },
  author_bar: {
    display: "flex",
    justifyContent: "space-between",
  },
  author_profile: {
    display: "flex",
  },
  cast_author_details: {
    display: "flex",
    flexDirection: "column",
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
    alignSelf: "center",
  },
  sub_button: {
    fontWeight: 600,
    width: 128,
  },
  description: {
    // marginLeft: 48,
    color: "#DDD",
    marginTop: 16,
    marginBottom: 16,
    // width: "80%"
  },
  comment_input: {
    borderRadius: "8px 48px 8px 8px",
  },
  comment_list: {
    marginTop: 32,
    marginBottom: 16,
  },
  comment_item: {
    marginTop: 16,
    marginBottom: 16,
  },
  clickable: {
    cursor: "pointer",
  },
  liked: {
    color: "#E84409",
  },
  h5: {
    fontFamily: "Open Sans",
    fontSize: 18,
    fontStyle: "italic",
    marginTop: 16,
    textAlign: "center",
  },
};

export default withRouter(Scene);
