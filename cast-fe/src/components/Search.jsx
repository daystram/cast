import React, { Component } from "react";
import { Card, Col, Container, Row } from "react-bootstrap";
import Sidebar from "./Sidebar";
import { Redirect } from "react-router-dom";
import MediaQuery from "react-responsive";
import { MOBILE_BP } from "../constants/breakpoint";
import List from "./List";

class Search extends Component {
  constructor(props) {
    super(props);
    this.state = {
      query: new URLSearchParams(this.props.location.search).get("query"),
      loading: false,
    };
  }

  componentDidMount() {
    document.title = this.state.query + " - cast";
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    let query = new URLSearchParams(this.props.location.search)
      .get("query")
      .trim();
    if (
      query !== prevState.query &&
      this.props.location.search !== prevProps.location.search
    ) {
      document.title = this.state.query + " - cast";
      this.setState({ query: query });
    }
  }

  render() {
    return (
      <>
        {!this.state.query && <Redirect to={"/"} />}
        <Container fluid style={style.content_container}>
          <Row>
            <MediaQuery minDeviceWidth={MOBILE_BP}>
              <Col xl={2} xs={12} style={{ marginBottom: 32 }}>
                <Card body style={{ borderRadius: "8px 48px 8px 8px" }}>
                  <Sidebar />
                </Card>
              </Col>
            </MediaQuery>
            <Col md={10} sm={12} className={"mid-container-right"}>
              <h1 style={style.h1}>Results for "{this.state.query}"</h1>
              <List
                search={true}
                query={this.state.query}
                emptyMessage={"No casts found!"}
                finishedMessage={"No more results!"}
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
  content_container: {},
};

export default Search;
