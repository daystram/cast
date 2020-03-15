import React, {Component} from 'react';
import {Container, Row, Col} from "react-bootstrap";
import Cast from "./Cast"
import Sidebar from "./Sidebar";

class Home extends Component {
  componentDidMount() {
    document.title = "cast";
  }
  render() {
    let data = {
      IDLIVE1: {
        _id: "IDLIVE1",
        thumbnail: "https://picsum.photos/seed/cast1/640/360",
        title: "Livestream 1 Title",
        isLive: true,
        views: 5690,
        author: {
          name: "Danny August",
          image: "https://picsum.photos/seed/profile/64/64"
        }
      },
      IDLIVE2: {
        _id: "IDLIVE2",
        thumbnail: "https://picsum.photos/seed/cast4/640/360",
        title: "Livestream 2 Title",
        isLive: true,
        views: 9876,
        author: {
          name: "Danny August",
          image: "https://picsum.photos/seed/profile/64/64"
        }
      },
      IDVOD1: {
        _id: "IDVOD1",
        thumbnail: "https://picsum.photos/seed/cast2/640/360",
        title: "Tokyo City",
        isLive: false,
        views: 293840,
        author: {
          name: "Danny August",
          image: "https://picsum.photos/seed/profile/64/64"
        }
      },
      IDVOD2: {
        _id: "IDVOD2",
        thumbnail: "https://picsum.photos/seed/cast3/640/360",
        title: "Big Buck Bunny",
        isLive: false,
        views: 823730,
        author: {
          name: "Danny August",
          image: "https://picsum.photos/seed/profile/64/64"
        }
      }
    };

    let liveSample = [];
    for (let i = 0; i < 8; i++) {
      liveSample.push(
        <Col xl={3} lg={4} md={6} sm={12} style={{padding: "0 8px 16px 8px"}}>
          <Cast video={Math.random() >= 0.5 ? data.IDLIVE1 : data.IDLIVE2}/>
        </Col>
      )
    }
    let vodSample = [];
    for (let i = 0; i < 8; i++) {
      vodSample.push(
        <Col xl={3} lg={4} md={6} sm={12} style={{padding: "0 8px 16px 8px"}}>
          <Cast video={Math.random() >= 0.5 ? data.IDVOD1 : data.IDVOD2}/>
        </Col>
      )
    }

    return (
      <>
        <Container fluid style={style.content_container}>
          <Row>
            <Col xs={2} style={style.sidebar_col}>
              <Sidebar/>
            </Col>
            <Col xs={10}>
              <h1 style={style.h1}>Live Casts</h1>
              <Row noGutters>
                {liveSample}
              </Row>
              <hr/>
              <h1 style={style.h1}>Fresh Casts</h1>
              <Row noGutters>
                {vodSample}
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
  content_container: {
    padding: "36px 0 0 0"
  },
};

export default Home
