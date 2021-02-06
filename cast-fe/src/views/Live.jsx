import React, { Component } from "react";
import { Card, Col, Container, Row } from "react-bootstrap";
import { MOBILE_BP } from "../constants/breakpoint";
import MediaQuery from "react-responsive";
import { VIDEO_TYPE_LIVE } from "../constants/video";
import { List, Sidebar } from "../components";

class Live extends Component {
  componentDidMount() {
    document.title = "Live Casts | cast";
  }

  render() {
    return (
      <>
        <Container fluid>
          <Row>
            <MediaQuery minDeviceWidth={MOBILE_BP}>
              <Col xl={2} xs={12} style={{ marginBottom: 32 }}>
                <Card body style={{ borderRadius: "8px 48px 8px 8px" }}>
                  <Sidebar />
                </Card>
              </Col>
            </MediaQuery>
            <Col xl={10} xs={12} className={"mid-container-right"}>
              <h1 style={style.h1}>Live Casts</h1>
              <List
                variant={VIDEO_TYPE_LIVE}
                emptyMessage={"No live casts today!"}
              />
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
  h5: {
    fontFamily: "Open Sans",
    fontSize: 18,
    fontStyle: "italic",
    marginTop: 16,
  },
  spinner: {
    margin: "32px auto 64px auto",
  },
};

export default Live;
