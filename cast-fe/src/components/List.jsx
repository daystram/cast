import React, {Component} from 'react';
import {Col, Row, Spinner} from "react-bootstrap";
import InfiniteScroll from 'react-infinite-scroller';
import Cast from "./Cast"
import axios from "axios";
import urls from "../helper/url";
import {VIDEO_LIST_PAGE_SIZE} from "../constants/video";

class List extends Component {
  constructor(props) {
    super(props);
    this.state = {
      page: 0,
      casts: [],
    };
    this.fetchVideos = this.fetchVideos.bind(this);
  }

  fetchVideos() {
    axios.get(urls().list(), {
      params: {
        variant: this.props.variant,
        count: VIDEO_LIST_PAGE_SIZE,
        offset: VIDEO_LIST_PAGE_SIZE * this.state.page,
      }
    }).then((response) => {
      if (response.data.code === 200) {
        if (!response.data.data) {
          this.setState({finished: true});
          return;
        }
        let casts = this.state.casts;
        response.data.data.map(cast => casts.push(cast));
        this.setState({casts, page: this.state.page + 1})
      }
    }).catch((error) => {
      console.log(error);
    });
  }

  render() {
    return (
      <InfiniteScroll
        style={{width: "100%"}}
        pageStart={0}
        hasMore={!this.state.finished}
        loadMore={this.fetchVideos}
        loader={
          <div style={{display: "flex"}}>
            <Spinner style={style.spinner} animation="grow" variant="primary"/>
          </div>
        }
        useWindow={true}
        initialLoad={true}>
        <Row noGutters>
          {this.state.casts.length ? this.state.casts.map(video =>
            <Col xl={3} lg={4} md={6} sm={12} key={video.hash} style={{padding: "0 8px 16px 8px"}}>
              <Cast video={video}/></Col>
          ) : <h5 style={style.h5}>No casts available!</h5>}
          {this.state.casts.length && this.state.finished &&
          <div style={{width: "100%", margin: "32px 64px"}}>
            <h5 style={{...style.h5, textAlign: "center"}}>No more casts!</h5>
          </div>}
        </Row>
      </InfiniteScroll>
    );
  }
}

let style = {
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

export default List
