import React, { Component } from "react";
import { Card, Col, Container, Row } from "react-bootstrap";
import SidebarProfile from "./SidebarProfile";
import axios from "axios";
import urls from "../helper/url";
import auth from "../helper/auth";
import MediaQuery from "react-responsive";
import { MOBILE_BP } from "../constants/breakpoint";

class Profile extends Component {
  constructor(props) {
    super(props);
    this.state = {
      user: auth().user(),
      subscribers: 0,
      views: 0,
      video_count: 0,
      before: {},
      error_password2: "",
      loading_info: true,
    };
    this.fetchUser = this.fetchUser.bind(this);
    this.fetchUser();
    document.title = "Profile | cast";
  }

  fetchUser() {
    axios
      .get(urls().user_info(), {})
      .then((response) => {
        this.setState({ loading_info: false });
        if (response.data.code === 200) {
          let data = response.data.data;
          this.setState({
            subscribers: data.subscribers,
            views: data.views,
            video_count: data.uploads,
          });
        }
      })
      .catch((error) => {
        console.log(error);
        this.setState({ loading_info: false });
      });
  }

  render() {
    return (
      <>
        <Container fluid style={style.content_container}>
          <Row>
            <MediaQuery minDeviceWidth={MOBILE_BP}>
              <Col xl={2} xs={12} style={{ marginBottom: 32 }}>
                <Card body style={{ borderRadius: "8px 48px 8px 8px" }}>
                  <SidebarProfile />
                </Card>
              </Col>
            </MediaQuery>
            <Col xl={10} xs={12} className={"mid-container-right"}>
              <Row>
                <Col xl xs={12} className={"responsive-fold"}>
                  <div style={style.profile_bar}>
                    <div
                      className="text-center"
                      style={{
                        width: 128,
                        height: 128,
                        flexShrink: 0,
                        borderRadius: 64,
                        background: "gray",
                        color: "white",
                        fontSize: "64px",
                        lineHeight: "128px",
                        textAlign: "center",
                        textTransform: "capitalize",
                      }}
                    >
                      {this.state.user.given_name[0] +
                        this.state.user.family_name[0]}
                    </div>
                    <div style={style.profile_name}>
                      <h1 style={style.name}>
                        {this.state.user.given_name}{" "}
                        {this.state.user.family_name}
                      </h1>
                      <h3 style={style.sub_count}>
                        {this.state.user.preferred_username}
                      </h3>
                    </div>
                  </div>
                  <h1 style={style.h1}>Showcase</h1>
                  <Card body style={style.profile_detail}>
                    <Row>
                      <Col sm={4}>
                        <p style={style.show_count}>{this.state.subscribers}</p>
                        <p style={style.show_caption}>
                          subscriber{this.state.subscribers === 1 ? "" : "s"}
                        </p>
                      </Col>
                      <Col sm={4}>
                        <p style={style.show_count}>{this.state.views}</p>
                        <p style={style.show_caption}>
                          total view{this.state.views === 1 ? "" : "s"}
                        </p>
                      </Col>
                      <Col sm={4}>
                        <p style={style.show_count}>{this.state.video_count}</p>
                        <p style={style.show_caption}>
                          cast{this.state.video_count === 1 ? "" : "s"} uploaded
                        </p>
                      </Col>
                    </Row>
                  </Card>
                </Col>
              </Row>
            </Col>
          </Row>
        </Container>
      </>
    );
  }
}

let style = {
  h1: {
    fontFamily: "Comfortaa",
  },
  content_container: {},
  name: {
    color: "#EBEBEB",
    margin: 0,
    fontSize: "3.5rem",
    lineHeight: 1.5,
    whiteSpace: "nowrap",
    textOverflow: "ellipsis",
    overflow: "hidden",
  },
  profile_upload: {
    width: 128,
    height: 128,
    background: "#f0f0f088",
    border: "3px dashed #ddddddaa",
    flexShrink: 0,
    borderRadius: 64,
    zIndex: 100,
  },
  profile_upload_modified: {
    background: "#f0f0f022",
    border: "",
  },
  profile_bar: {
    display: "flex",
    marginBottom: 38,
  },
  profile_name: {
    overflow: "hidden",
    marginLeft: 16,
    alignSelf: "center",
    width: "100%",
  },
  sub_count: {
    color: "#DDD",
    fontSize: 22,
  },
  profile_detail: {
    borderRadius: "8px 48px 8px 8px",
    marginBottom: 32,
  },
  email: {
    fontSize: 18,
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

export default Profile;
