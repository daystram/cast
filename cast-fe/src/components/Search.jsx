import React, {Component} from 'react';
import {Card, Col, Container, Row, Spinner} from "react-bootstrap";
import Cast from "./Cast"
import Sidebar from "./Sidebar";
import axios from "axios";
import urls from "../helper/url";
import {Redirect} from "react-router-dom";
import MediaQuery from "react-responsive";
import {MOBILE_BP} from "../constants/breakpoint";

class Search extends Component {
  constructor(props) {
    super(props);
    this.state = {
      query: new URLSearchParams(this.props.location.search).get("query"),
      results: null,
      loading: false
    }
  }

  componentDidMount() {
    document.title = this.state.query + " - cast";
    if (this.state.query) this.fetchResults(this.state.query);
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    let query = new URLSearchParams(this.props.location.search).get("query").trim();
    if (query !== prevState.query && this.props.location.search !== prevProps.location.search) {
      document.title = this.state.query + " - cast";
      this.setState({query: query});
      if (query) this.fetchResults(query);
    }
  }

  fetchResults(query) {
    axios.get(urls().search(), {
      params: {
        query: query.trim(),
        count: 8,
        offset: 0,
      }
    }).then((response) => {
      this.setState({loading: false});
      if (response.data.code === 200) {
        this.setState({results: response.data.data})
      }
    }).catch((error) => {
      console.log(error);
      this.setState({loading: false});
    });
  }

  render() {
    return (
      <>
        {!this.state.query && <Redirect to={"/"}/>}
        <Container fluid style={style.content_container}>
          <Row>
            <MediaQuery minDeviceWidth={MOBILE_BP}>
              <Col xl={2} xs={12} style={{marginBottom: 32}}>
                <Card body style={{borderRadius: "8px 48px 8px 8px"}}><Sidebar/></Card>
              </Col>
            </MediaQuery>
            <Col md={10} sm={12} className={"mid-container-right"}>
              <h1 style={style.h1}>"{this.state.query}"</h1>
              <Row noGutters>
                {!this.state.loading && (this.state.results ? this.state.results.map(video =>
                  <Col xl={3} lg={4} md={6} sm={12} key={video.hash} style={{padding: "0 8px 16px 8px"}}>
                    <Cast video={video}/></Col>
                ) : <h5 style={style.h5}>No casts found!</h5>)}
                {this.state.loading && <Spinner style={style.spinner} animation="grow" variant="primary"/>}
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
  },
};

export default Search
