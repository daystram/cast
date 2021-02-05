import React, { Component } from "react";
import {
  Alert,
  Button,
  Card,
  Col,
  Container,
  Row,
  Spinner,
} from "react-bootstrap";
import { MOBILE_BP } from "../constants/breakpoint";
import MediaQuery from "react-responsive";
import { Link } from "react-router-dom";
import { VIDEO_LIST_PAGE_SIZE } from "../constants/video";
import api from "../apis/api";
import { Cast, Sidebar } from "../components";
import logo from "../components/logo.svg";

class Home extends Component {
  constructor(props) {
    super(props);
    this.state = {
      trending: null,
      live: null,
      vod: null,
      loading: {
        trending: true,
        live: true,
        vod: true,
      },
    };
  }

  componentDidMount() {
    document.title = "cast";
    this.fetchVideos("trending");
    this.fetchVideos("live");
    this.fetchVideos("vod");
  }

  fetchVideos(variant) {
    api.cast
      .list({
        variant: variant,
        count: VIDEO_LIST_PAGE_SIZE,
        offset: 0,
      })
      .then((response) => {
        this.setState({ loading: { ...this.state.loading, [variant]: false } });
        if (response.data.code === 200) {
          if (response.data.data && response.data.data.length === 8)
            response.data.data.pop();
          this.setState({ [variant]: response.data.data });
        }
      })
      .catch((error) => {
        console.log(error);
        this.setState({ loading: { ...this.state.loading, [variant]: false } });
      });
  }

  render() {
    let viewMore = (target) => (
      <Col
        xl={3}
        lg={4}
        md={6}
        sm={12}
        style={{
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
          padding: "32px 0",
        }}
      >
        <Link className={"inline"} to={target}>
          <p style={{ margin: 0, textAlign: "center", fontSize: 20 }}>
            see more{" "}
            <span
              className="material-icons"
              style={{ fontSize: 20, verticalAlign: "middle" }}
            >
              arrow_forward_ios
            </span>
          </p>
        </Link>
      </Col>
    );

    return (
      <Container fluid style={style.content_container}>
        <Row>
          <MediaQuery minDeviceWidth={MOBILE_BP}>
            <Col xl={2} xs={12} style={{ marginBottom: 32 }}>
              <Card body style={{ borderRadius: "8px 48px 8px 8px" }}>
                <Sidebar />
              </Card>
            </Col>
          </MediaQuery>
          <Col xl={10} xs={12} className={"mid-container-right"}>
            <Alert variant="banner">
              <Alert.Heading>
                <img
                  src={logo}
                  height="30"
                  className="d-inline-block align-top cast-logo"
                  alt="cast"
                />
              </Alert.Heading>
              <p>
                cast is a DASH video-streaming and RTMP live-streaming platform.
              </p>
              <p>
                It uses <b>DASH (Dynamic Adaptive Streaming over HTTP)</b> to
                smoothly stream videos with ranging qualities to the clients.{" "}
                <b>RTMP (Real-Time Messaging Protocol)</b> is the method used to
                stream live content from the client (e.g. OBS Studio) to cast's
                servers.
              </p>
              <p>
                Uploaded videos are ingested by cast's <b>transcoding nodes</b>{" "}
                powered with <b>NVIDIA CUDA GPU hardware acceleration</b>. These
                transcoders are <b>highly scalable</b> and can be deployed with
                a high number of replica either on- or off-premise.
              </p>
              <p>
                User and authentication is managed by{" "}
                <b>
                  <a
                    className={"inline"}
                    href={"https://ratify.daystram.com"}
                    target={"_blank"}
                  >
                    Ratify
                  </a>
                </b>
                .
              </p>
              <hr />
              <div className="d-flex justify-content-end">
                <Button
                  href={"https://github.com/daystram/cast"}
                  target={"_blank"}
                  variant="outline-primary"
                >
                  View on GitHub <i className="fab fa-github" />
                </Button>
              </div>
            </Alert>
            <h1 style={style.h1}>Trending Casts</h1>
            <Row noGutters>
              {!this.state.loading.trending &&
                (this.state.trending ? (
                  this.state.trending.map((video) => (
                    <Col
                      xl={3}
                      lg={4}
                      md={6}
                      sm={12}
                      key={video.hash}
                      style={{ padding: "0 8px 16px 8px" }}
                    >
                      <Cast video={video} />
                    </Col>
                  ))
                ) : (
                  <h5 style={style.h5}>No casts uploaded yet!</h5>
                ))}
              {!this.state.loading.trending &&
                this.state.trending &&
                viewMore("/trending")}
              {this.state.loading.trending && (
                <Spinner
                  style={style.spinner}
                  animation="grow"
                  variant="primary"
                />
              )}
            </Row>
            <hr />
            <h1 style={style.h1}>Live Casts</h1>
            <Row noGutters>
              {!this.state.loading.live &&
                (this.state.live ? (
                  this.state.live.map((video) => (
                    <Col
                      xl={3}
                      lg={4}
                      md={6}
                      sm={12}
                      key={video.hash}
                      style={{ padding: "0 8px 16px 8px" }}
                    >
                      <Cast video={video} />
                    </Col>
                  ))
                ) : (
                  <h5 style={style.h5}>No live casts today!</h5>
                ))}
              {!this.state.loading.live && this.state.live && viewMore("/live")}
              {this.state.loading.live && (
                <Spinner
                  style={style.spinner}
                  animation="grow"
                  variant="primary"
                />
              )}
            </Row>
            <hr />
            <h1 style={style.h1}>Fresh Casts</h1>
            <Row noGutters>
              {!this.state.loading.vod &&
                (this.state.vod ? (
                  this.state.vod.map((video) => (
                    <Col
                      xl={3}
                      lg={4}
                      md={6}
                      sm={12}
                      key={video.hash}
                      style={{ padding: "0 8px 16px 8px" }}
                    >
                      <Cast video={video} />
                    </Col>
                  ))
                ) : (
                  <h5 style={style.h5}>No casts uploaded yet!</h5>
                ))}
              {!this.state.loading.vod && this.state.vod && viewMore("/fresh")}
              {this.state.loading.vod && (
                <Spinner
                  style={style.spinner}
                  animation="grow"
                  variant="primary"
                />
              )}
            </Row>
          </Col>
        </Row>
      </Container>
    );
  }
}

let style = {
  h1: {
    fontFamily: "Comfortaa",
  },
  h5: {
    fontFamily: "Open Sans",
    fontSize: 18,
    fontStyle: "italic",
    marginTop: 16,
  },
  spinner: {
    margin: "32px auto 64px auto",
  },
  content_container: {},
};

export default Home;
