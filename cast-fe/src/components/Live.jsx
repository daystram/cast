import React, {Component} from 'react';
import {Card, Col, Container, Row, Spinner} from "react-bootstrap";
import Cast from "./Cast"
import Sidebar from "./Sidebar";
import axios from "axios";
import urls from "../helper/url";
import {MOBILE_BP} from "../constants/breakpoint";
import MediaQuery from "react-responsive";
import List from "./List";

class Trending extends Component {
  constructor(props) {
    super(props);
    this.state = {
      trending: null,
      loading: {
        trending: true,
      }
    }
  }

  componentDidMount() {
    document.title = "Trending Casts | cast";
    this.fetchVideos("trending");
  }

  fetchVideos(variant) {
    axios.get(urls().list(), {
      params: {
        variant: variant,
        count: 16,
        offset: 0,
      }
    }).then((response) => {
      this.setState({loading: {...this.state.loading, [variant]: false}});
      if (response.data.code === 200) {
        this.setState({[variant]: response.data.data})
      }
    }).catch((error) => {
      console.log(error);
      this.setState({loading: {...this.state.loading, [variant]: false}});
    });
  }

  render() {
    return (
      <>
        <Container fluid>
          <Row>
            <MediaQuery minDeviceWidth={MOBILE_BP}>
              <Col xl={2} xs={12} style={{marginBottom: 32}}>
                <Card body style={{borderRadius: "8px 48px 8px 8px"}}><Sidebar/></Card>
              </Col>
            </MediaQuery>
            <Col xl={10} xs={12} className={"mid-container-right"}>
              <h1 style={style.h1}>Trending Casts</h1>
              <Row noGutters>
                <List variant={"trending"}/>
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
    fontFamily: "Comfortaa"
  },
  h5: {
    fontFamily: "Open Sans",
    fontSize: 18,
    fontStyle: "italic",
    marginTop: 16
  },
  spinner: {
    margin: "32px auto 64px auto",
  },
};

export default Trending
