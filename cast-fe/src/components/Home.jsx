import React, {Component} from 'react';
import {Col, Container, Row, Spinner} from "react-bootstrap";
import Cast from "./Cast"
import Sidebar from "./Sidebar";
import axios from "axios";
import urls from "../helper/url";

class Home extends Component {
  constructor(props) {
    super(props);
    this.state = {
      live: [],
      vod: [],
      loading: {
        live: true,
        vod: true,
      }
    }
  }

  componentDidMount() {
    document.title = "cast";
    this.fetchVideos("live");
    this.fetchVideos("vod");
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
        <Container fluid style={style.content_container}>
          <Row>
            <Col md={2} sm={12} style={{marginBottom: 32}}>
              <Sidebar/>
            </Col>
            <Col md={10} sm={12}>
              <h1 style={style.h1}>Live Casts</h1>
              <Row noGutters>
                {!this.state.loading.live && (this.state.live ? this.state.live.map(video =>
                  <Col xl={3} lg={4} md={6} sm={12} style={{padding: "0 8px 16px 8px"}}><Cast video={video}/></Col>
                ) : <h5 style={style.h5}>No live casts today!</h5>)}
                {this.state.loading.live && <Spinner style={style.spinner} animation="grow" variant="primary"/>}
              </Row>
              <hr/>
              <h1 style={style.h1}>Fresh Casts</h1>
              <Row noGutters>
                {!this.state.loading.vod && (this.state.vod ? this.state.vod.map(video =>
                  <Col xl={3} lg={4} md={6} sm={12} style={{padding: "0 8px 16px 8px"}}><Cast video={video}/></Col>
                ) : <h5 style={style.h5}>No fresh casts today!</h5>)}
                {this.state.loading.vod && <Spinner style={style.spinner} animation="grow" variant="primary"/>}
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
  content_container: {
    padding: "36px 0 0 0"
  },
};

export default Home
