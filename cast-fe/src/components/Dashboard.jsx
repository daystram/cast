import React, { Component } from "react";
import {
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
import SidebarProfile from "./SidebarProfile";
import Chat from "./Chat";
import auth from "../helper/auth";
import CastEditable from "./CastEditable";
import axios from "axios";
import urls from "../helper/url";
import format from "../helper/format";
import Clock from "react-live-clock";
import { MOBILE_BP } from "../constants/breakpoint";
import MediaQuery from "react-responsive";

let interval = null;

class Dashboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
      loading_status: false,
    };
    document.title = "Dashboard | cast";
    this.controlStream = this.controlStream.bind(this);
    this.fetchDetail = this.fetchDetail.bind(this);
    this.setStreamWindow = this.setStreamWindow.bind(this);
    this.loadLive = this.loadLive.bind(this);
    this.fetchDetail(auth().username());
  }

  fetchDetail(hash) {
    axios
      .get(urls().cast_details(), {
        params: {
          hash: hash,
          username: auth().username(),
        },
      })
      .then((response) => {
        this.setState({ loading: false });
        if (response.data.code === 200) {
          let data = response.data.data;
          this.setState({
            stream: data,
            created_at: data.created_at,
            delta: new Date() - new Date(data.created_at),
            likes: data.likes,
            live: data.is_live,
            views: data.views,
            pending: data.pending,
          });
          if (data.pending || data.is_live) this.loadLive();
        }
      })
      .catch((error) => {
        console.log(error);
        this.setState({ loading: false });
      });
  }

  controlStream(e) {
    if (this.state.stream && this.state.live) {
      this.setState({ prompt_stop: true });
    } else {
      if (this.state.pending) {
        this.setStreamWindow(false);
      } else {
        this.setStreamWindow(true);
        e.target.blur();
      }
    }
  }

  loadLive() {
    clearInterval(interval);
    interval = setInterval(() => {
      axios.get(urls().window()).then((response) => {
        if (response.data.code === 200) {
          if (response.data.data && this.state.pending) {
            this.setState({ delta: 1, created_at: new Date().toISOString() });
          }
          if (!response.data.data && !this.state.pending)
            clearInterval(interval);
          this.setState({
            live: response.data.data,
            stream: { ...this.state.stream, is_live: response.data.data },
            pending: this.state.pending && !response.data.data,
          });
        }
      });
    }, 3000);
  }

  setStreamWindow(open) {
    this.setState({ loading_status: true });
    axios
      .put(urls().edit_window(open))
      .then((response) => {
        if (response.data.code === 200) {
          this.setState({ loading_status: false });
          if (!open) {
            clearInterval(interval);
            this.setState({
              live: false,
              pending: false,
              prompt_stop: false,
            });
            if (this.state.live)
              this.setState({ created_at: new Date().toISOString() });
          } else {
            this.loadLive();
            this.setState({ pending: true });
          }
        } else {
          console.log(response.data);
          this.setState({
            error_edit: response.data.error,
            loading_status: false,
          });
        }
      })
      .catch((error) => {
        console.log(error);
        this.setState({
          error_edit: "An error has occurred! Please try again.",
          loading_status: false,
        });
      });
  }

  render() {
    return (
      <>
        <Container fluid style={style.content_container}>
          <Row>
            <MediaQuery minDeviceWidth={MOBILE_BP}>
              <Col xl={{ span: 2, order: 1 }} style={{ marginBottom: 32 }}>
                <Card body style={{ borderRadius: "8px 48px 8px 8px" }}>
                  <SidebarProfile />
                </Card>
              </Col>
            </MediaQuery>
            <Col
              xl={{ span: 8, order: 2 }}
              md={{ span: 6, order: 1 }}
              xs={{ span: 12, order: 2 }}
              className={"mid-container"}
            >
              <Row noGutters>
                <Col>
                  <h1 style={style.h1}>
                    Livestream{" "}
                    {this.state.stream && this.state.live && (
                      <Badge pill style={style.live_tag}>
                        LIVE
                      </Badge>
                    )}
                  </h1>
                </Col>
                <Col md={"auto"} sm={12}>
                  {this.state.stream && (
                    <Button
                      disabled={this.state.loading_status}
                      className={"responsive-fold"}
                      variant={
                        this.state.stream &&
                        (this.state.live
                          ? "danger"
                          : this.state.pending
                          ? "warning"
                          : "primary")
                      }
                      size={"md"}
                      block
                      style={style.live_button}
                      onClick={this.controlStream}
                    >
                      {this.state.stream &&
                        !this.state.loading_status &&
                        (this.state.live ? (
                          "End Stream"
                        ) : this.state.pending ? (
                          <>
                            Waiting for Up-link{"  "}
                            <Spinner animation="grow" style={style.spinner} />
                          </>
                        ) : (
                          "Go Live"
                        ))}
                      {this.state.loading_status && (
                        <Spinner animation="grow" style={style.spinner} />
                      )}
                    </Button>
                  )}
                </Col>
              </Row>
              <div style={style.card_detail}>
                {this.state.stream ? (
                  <CastEditable video={this.state.stream} />
                ) : (
                  <div style={{ display: "flex", justifyContent: "center" }}>
                    <Spinner
                      style={style.spinner}
                      animation="grow"
                      variant="primary"
                    />
                  </div>
                )}
              </div>
              <h1 style={style.h1}>Metrics</h1>
              {this.state.stream ? (
                <Card body style={style.card_detail}>
                  <Row>
                    <Col sm={6}>
                      <p style={style.show_count}>
                        {this.state.live &&
                        this.state.created_at &&
                        this.state.delta ? (
                          <Clock
                            date={this.state.delta}
                            timezone={"UTC"}
                            format={"HH:mm:ss"}
                            ticking={true}
                          />
                        ) : this.state.created_at === "0001-01-01T00:00:00Z" ? (
                          "Never"
                        ) : (
                          format().full_date(this.state.created_at)
                        )}
                      </p>
                      <p style={style.show_caption}>
                        {this.state.live ? "stream duration" : "last stream"}
                      </p>
                    </Col>
                    <Col sm={6}>
                      <p style={style.show_count}>
                        {this.state.views - 1 || 0}
                      </p>
                      <p style={style.show_caption}>
                        {this.state.live || "last"} total viewer
                        {this.state.views - 1 === 1 ? "" : "s"}
                      </p>
                    </Col>
                  </Row>
                </Card>
              ) : (
                <div
                  style={{
                    display: "flex",
                    justifyContent: "center",
                    marginBottom: 32,
                  }}
                >
                  <Spinner
                    style={style.spinner}
                    animation="grow"
                    variant="primary"
                  />
                </div>
              )}
              <h1 style={style.h1}>Connection</h1>
              {this.state.stream ? (
                <Card body style={style.card_detail}>
                  <Row>
                    <Col md={6} sm={12}>
                      <Form autocomplete="off">
                        <Form.Group>
                          <Form.Label>Server Address</Form.Label>
                          <InputGroup>
                            <Form.Control
                              type="text"
                              value={`rtmp://cast.daystram.com/live`}
                              ref={(ref) => (this.rtmpField = ref)}
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
                                    this.rtmpField.select();
                                    document.execCommand("copy");
                                    this.rtmpField.blur();
                                  }}
                                >
                                  Copy
                                </Button>
                              </OverlayTrigger>
                            </InputGroup.Append>
                          </InputGroup>
                        </Form.Group>
                      </Form>
                    </Col>
                    <Col md={6} sm={12}>
                      <Form autocomplete="off">
                        <Form.Group>
                          <Form.Label>Stream Key</Form.Label>
                          <InputGroup>
                            <Form.Control
                              type={this.state.show_key ? "text" : "password"}
                              value={auth().username()}
                            />
                            <InputGroup.Append>
                              <Button
                                variant={
                                  this.state.show_key
                                    ? "primary"
                                    : "outline-primary"
                                }
                                onClick={() => {
                                  this.setState({
                                    show_key: !this.state.show_key,
                                  });
                                }}
                              >
                                <span className="material-icons">
                                  {this.state.show_key
                                    ? "visibility_off"
                                    : "visibility"}
                                </span>
                              </Button>
                            </InputGroup.Append>
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
                                    navigator.clipboard.writeText(
                                      auth().username()
                                    );
                                  }}
                                >
                                  Copy
                                </Button>
                              </OverlayTrigger>
                            </InputGroup.Append>
                          </InputGroup>
                        </Form.Group>
                      </Form>
                    </Col>
                    <Col sm={12}>
                      <Form autocomplete="off">
                        <Form.Group>
                          <Form.Label>Chat Embed</Form.Label>
                          <InputGroup>
                            <Form.Control
                              type="text"
                              value={`https://cast.daystram.com/c/${auth().username()}`}
                              ref={(ref) => (this.chatField = ref)}
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
                                    this.chatField.select();
                                    document.execCommand("copy");
                                    this.chatField.blur();
                                  }}
                                >
                                  Copy
                                </Button>
                              </OverlayTrigger>
                            </InputGroup.Append>
                          </InputGroup>
                        </Form.Group>
                      </Form>
                    </Col>
                  </Row>
                </Card>
              ) : (
                <div
                  style={{
                    display: "flex",
                    justifyContent: "center",
                    marginBottom: 32,
                  }}
                >
                  <Spinner
                    style={style.spinner}
                    animation="grow"
                    variant="primary"
                  />
                </div>
              )}
            </Col>
            <Col
              xl={{ span: 2, order: 3 }}
              md={{ span: 6, order: 2 }}
              xs={{ span: 12, order: 2 }}
              style={{ paddingBottom: 16 }}
            >
              <h1 style={style.h1}>Chat</h1>
              <Chat height={"85vh"} embedded={true} hash={auth().username()} />
            </Col>
          </Row>
        </Container>
        <Modal
          show={this.state.live && this.state.prompt_stop}
          size={"md"}
          centered
        >
          <Modal.Header>
            <Modal.Title id="contained-modal-title-vcenter">
              Ending Stream
            </Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <p>Are you sure you want to stop the stream?</p>
          </Modal.Body>
          <Modal.Footer>
            <Button
              variant={"info"}
              onClick={() => this.setState({ prompt_stop: false })}
            >
              Cancel
            </Button>
            <Button
              variant={"danger"}
              disabled={this.state.loading_status}
              onClick={() => this.setStreamWindow(false)}
            >
              {this.state.loading_status ? (
                <Spinner animation="grow" style={style.spinner} />
              ) : (
                "End Stream"
              )}
            </Button>
          </Modal.Footer>
        </Modal>
      </>
    );
  }
}

let style = {
  h1: {
    fontFamily: "Comfortaa",
    display: "inline-block",
  },
  content_container: {},
  live_tag: {
    background: "red",
    color: "white",
    marginRight: 8,
    marginBottom: 8,
    borderRadius: 8,
    fontSize: 24,
    fontWeight: 400,
    verticalAlign: "middle",
  },
  live_button: {
    verticalAlign: "top",
    fontWeight: 600,
    height: "fit-content",
  },
  spinner: {
    width: 16,
    height: 16,
    verticalAlign: "initial",
  },
  card_detail: {
    borderRadius: "8px 48px 8px 8px",
    marginBottom: 32,
  },
  show_count: {
    textAlign: "center",
    fontSize: 28,
    fontWeight: 800,
    marginBottom: 8,
  },
  show_caption: {
    textAlign: "center",
    marginBottom: 0,
    fontWeight: 600,
  },
};

export default Dashboard;
