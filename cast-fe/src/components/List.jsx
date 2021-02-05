import React, { Component } from "react";
import { Col, Row, Spinner } from "react-bootstrap";
import InfiniteScroll from "react-infinite-scroller";
import Cast from "./Cast";
import {
  VIDEO_LIST_LIKED,
  VIDEO_LIST_PAGE_SIZE,
  VIDEO_LIST_SUBSCRIBED,
} from "../constants/video";
import api from "../apis/api";

class List extends Component {
  constructor(props) {
    super(props);
    this.state = {
      page: 0,
      casts: [],
      finished: false,
    };
    this.fetchVideos = this.fetchVideos.bind(this);
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.search && this.props.query !== prevProps.query) {
      this.setState({ casts: [], finished: false, page: 0 });
    }
  }

  fetchVideos() {
    let request;
    if (this.props.search) {
      request = api.cast.search({
        query: this.props.query.trim(),
        count: VIDEO_LIST_PAGE_SIZE,
        offset: VIDEO_LIST_PAGE_SIZE * this.state.page,
      });
    } else {
      const params = {
        variant: this.props.variant,
        count: VIDEO_LIST_PAGE_SIZE,
        offset: VIDEO_LIST_PAGE_SIZE * this.state.page,
      };
      switch (this.props.variant) {
        case VIDEO_LIST_LIKED:
          request = api.cast.listCurated(params);
          break;
        case VIDEO_LIST_SUBSCRIBED:
          request = api.cast.listCurated(params);
          break;
        default:
          request = api.cast.list(params);
      }
    }
    request
      .then((response) => {
        if (response.data.code === 200) {
          if (!response.data.data) {
            this.setState({ finished: true });
            return;
          }
          let casts = this.state.casts;
          response.data.data.map((cast) => casts.push(cast));
          this.setState({ casts, page: this.state.page + 1 });
        }
      })
      .catch((error) => {
        console.log(error);
      });
  }

  render() {
    return (
      <InfiniteScroll
        style={{ width: "100%" }}
        pageStart={0}
        hasMore={!this.state.finished}
        loadMore={this.fetchVideos}
        loader={
          <div style={{ display: "flex" }} key={""}>
            <Spinner style={style.spinner} animation="grow" variant="primary" />
          </div>
        }
        useWindow={true}
        initialLoad={true}
      >
        <Row noGutters>
          {!!this.state.casts.length ? (
            this.state.casts.map((video) => (
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
            <h5 style={style.h5}>
              {this.props.emptyMessage || "No casts available!"}
            </h5>
          )}
          {!!this.state.casts.length && this.state.finished && (
            <div style={{ width: "100%", margin: "32px 64px" }}>
              <h5 style={{ ...style.h5, textAlign: "center" }}>
                {this.props.finishedMessage || "No more casts!"}
              </h5>
            </div>
          )}
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
    marginTop: 16,
  },
  spinner: {
    margin: "32px auto 64px auto",
  },
};

export default List;
