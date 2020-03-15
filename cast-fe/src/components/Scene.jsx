import React, {Component} from 'react';
import {Container, Row, Col, Card, Image, Badge, FormControl, InputGroup, Button} from "react-bootstrap";
import Cast from "./Cast"
import Navigation from "./Navigation"
import Sidebar from "./Sidebar";
import {HybridPlayer} from "./player";
import abbreviate from "../helper/abbreviate";

class Scene extends Component {
  constructor(props) {
    super(props);
    this.state = {video: this.fetchDetail(this.props.match.params._id)};
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    if (this.props.match.params._id !== prevProps.match.params._id) {
      this.setState({video: this.fetchDetail(this.props.match.params._id)});
    }
  }

  fetchDetail(id) {
    // TODO: fetch from cast-be
    // Mock BE data
    let detail_data = {
      IDLIVE1: {
        _id: "IDLIVE1",
        thumbnail: "https://picsum.photos/seed/cast1/640/360", // TODO: Will not be given
        url: "http://cast.daystram.com/api/test_stream", // TODO: Will not be given
        title: "Livestream 1 Title",
        isLive: true,
        views: 5690,
        likes: 8402,
        isLiked: true,
        description: "This is stream 1 description. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
        CreatedAt: "March 14, 2020",
        author: {
          name: "Danny August",
          image: "https://picsum.photos/seed/profile/64/64", // TODO: Will not be given
          subscribers: 8123905,
          isSubscribed: true
        }
      },
      IDLIVE2: {
        _id: "IDLIVE2",
        thumbnail: "https://picsum.photos/seed/cast4/640/360", // TODO: Will not be given
        url: "http://cast.daystram.com/api/test_stream2", // TODO: Will not be given
        title: "Livestream 2 Title",
        isLive: true,
        views: 9876,
        likes: 2345,
        isLiked: true,
        description: "This is stream 2 description. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur.",
        CreatedAt: "March 14, 2020",
        author: {
          name: "Danny August",
          image: "https://picsum.photos/seed/profile/64/64", // TODO: Will not be given
          subscribers: 8123905,
          isSubscribed: true
        }
      },
      IDVOD1: {
        _id: "IDVOD1",
        thumbnail: "https://picsum.photos/seed/cast2/640/360",  // TODO: Will not be given
        url: "https://storage.googleapis.com/cast-uploaded-videos/HASH02/manifest.mpd", // TODO: Will not be given
        title: "Tokyo City",
        isLive: false,
        views: 5690724,
        likes: 91843,
        isLiked: false,
        description: "This is video 1 description. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat. Duis aute irure dolor in reprehenderit in voluptate velit esse cillum dolore eu fugiat nulla pariatur. Excepteur sint occaecat cupidatat non proident, sunt in culpa qui officia deserunt mollit anim id est laborum.",
        CreatedAt: "March 10, 2020",
        author: {
          name: "Danny August",
          image: "https://picsum.photos/seed/profile/64/64", // TODO: Will not be given
          subscribers: 8123905,
          isSubscribed: true
        }
      },
      IDVOD2: {
        _id: "IDVOD2",
        thumbnail: "https://picsum.photos/seed/cast3/640/360",  // TODO: Will not be given
        url: "https://storage.googleapis.com/cast-uploaded-videos/HASH01/manifest.mpd", // TODO: Will not be given
        title: "Big Buck Bunny",
        isLive: false,
        views: 823730,
        likes: 43223,
        isLiked: false,
        description: "This is video 2 description. Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua. Ut enim ad minim veniam, quis nostrud exercitation ullamco laboris nisi ut aliquip ex ea commodo consequat.",
        CreatedAt: "March 5, 2020",
        author: {
          name: "Danny August",
          image: "https://picsum.photos/seed/profile/64/64", // TODO: Will not be given
          subscribers: 8123905,
          isSubscribed: true
        }
      }
    };
    let detail = detail_data[id];
    document.title = `${detail.title} - ${detail.author.name} | cast`;
    return detail
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
        <Row noGutters style={{padding: "0 0 16px 0"}}>
          <Cast video={Math.random() >= 0.5 ? data.IDLIVE1 : data.IDLIVE2}/>
        </Row>
      )
    }
    let vodSample = [];
    for (let i = 0; i < 8; i++) {
      vodSample.push(
        <Row noGutters style={{padding: "0 0 16px 0"}}>
          <Cast video={Math.random() >= 0.5 ? data.IDVOD1 : data.IDVOD2}/>
        </Row>
      )
    }
    return (
      <>
        <Container fluid style={style.content_container}>
          <Row>
            <Col xl={2} md={12}>
              <Sidebar/>
              <div style={style.cast_list}>
                {liveSample}
              </div>
            </Col>
            <Col xl md={12}>
              <HybridPlayer url={this.state.video.url} thumbnail={this.state.video.thumbnail}
                            live={this.state.video.isLive}/>
              <div style={style.cast_tag_bar}>
                <div>
                  <Badge pill style={style.cast_tag}>tag</Badge>
                  <Badge pill style={style.cast_tag}>another</Badge>
                </div>
                <div>
                  <span style={style.cast_attrib}><i className="material-icons">share</i>{" "}share</span>
                  <span style={style.cast_attrib}>
                    <i className="material-icons">thumb_up</i>
                    {" "}{abbreviate(this.state.video.likes)} likes
                  </span>
                  <span style={style.cast_attrib}>
                    <i className="material-icons">remove_red_eye</i>
                    {" "}{abbreviate(this.state.video.views)} {this.state.video.isLive ? 'viewers' : 'views'}
                  </span>
                </div>
              </div>
              <h1 style={style.title}>{this.state.video.title}</h1>
              <p>{this.state.video.CreatedAt}</p>
              <div style={style.author_bar}>
                <div style={style.author_profile}>
                  <Image src={this.state.video.author.image} height={42} width={42}
                         style={style.profile_image} onClick={this.viewAuthor} roundedCircle/>
                  <div style={style.cast_author_details}>
                    <p style={style.cast_author_name} onClick={this.viewAuthor}>{this.state.video.author.name}</p>
                    <p style={style.cast_author_sub}>{abbreviate(this.state.video.author.subscribers)} subscribers</p>
                  </div>
                </div>
                <div>
                  <Button style={style.tip_button}><i className="material-icons">attach_money</i></Button>
                  <Button style={style.sub_button} disabled={this.state.video.author.isSubscribed}>SUBSCRIBE</Button>
                </div>
              </div>
              <Row noGutters>
                <Col xl={1} sm={0}/>
                <Col>
                  <div style={style.description}>{this.state.video.description}</div>
                </Col>
                <Col xl={1} sm={0}/>
              </Row>
              <hr/>
              <h3>Comments</h3>
              <Row noGutters>
                <Col xl={1} sm={0}/>
                <Col>
                  <InputGroup style={style.comment_input}>
                    <FormControl type="text" placeholder="Comment"/>
                    <InputGroup.Append>
                      <Button variant="outline-info"><i className="material-icons">send</i></Button>
                    </InputGroup.Append>
                  </InputGroup>
                </Col>
                <Col xl={1} sm={0}/>
              </Row>

            </Col>
            <Col xl={2} md={12}>
              <Card style={style.live_chat}>
                <Card.Body style={style.live_chat_body}>
                  <p style={style.live_chat_item}><b>User1</b>: Hello! Hello! Hello! Hello! Hello! Hello! Hello!</p>
                  <p style={style.live_chat_item}><b>User2</b>: Hello!</p>
                  <p style={style.live_chat_item}><b>User1</b>: Hello!</p>
                  <p style={style.live_chat_item}><b>User2</b>: Hello!</p>
                  <p style={style.live_chat_item}><b>User1</b>: Hello!</p>
                  <InputGroup style={style.live_chat_input}>
                    <FormControl type="text" placeholder="Chat"/>
                    <InputGroup.Append>
                      <Button variant="outline-info"><i className="material-icons">send</i></Button>
                    </InputGroup.Append>
                  </InputGroup>
                </Card.Body>
              </Card>
              <div style={style.cast_list}>
                {vodSample}
              </div>
            </Col>
          </Row>
        </Container>
      </>
    );
  }
}

let style = {
  title: {
    margin: 0
  },
  content_container: {
    padding: "36px 0 0 0"
  },
  cast_list: {
    marginTop: 16
  },
  live_chat: {
    borderRadius: "8px 48px 8px 8px",
  },
  live_chat_body: {
    height: 480,
    justifyContent: "flex-end",
    display: "flex",
    flexDirection: "column"
  },
  live_chat_item: {
    marginBottom: 0
  },
  live_chat_input: {
    marginTop: 8
  },
  cast_tag_bar: {
    marginTop: 8,
    marginBottom: 8,
    display: "flex",
    justifyContent: "space-between"
  },
  cast_tag: {
    background: "white",
    borderRadius: 8,
    borderWidth: 1,
    borderColor: "lightgray",
    borderStyle: "solid",
    fontSize: 16,
    fontWeight: 400,
    marginRight: 8
  },
  cast_attrib: {
    marginLeft: 16
  },
  author_bar: {
    display: "flex",
    justifyContent: "space-between"
  },
  author_profile: {
    display: "flex",
  },
  cast_author_details: {
    display: "flex",
    flexDirection: "column"
  },
  cast_author_name: {
    fontSize: 18,
    fontWeight: 600,
    margin: 0,
  },
  cast_author_sub: {
    fontSize: 16,
    fontWeight: 400,
    margin: 0,
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
  },
  profile_image: {
    marginRight: 8,
    alignSelf: "center"
  },
  tip_button: {
    marginRight: 8
  },
  sub_button: {
    fontWeight: 600
  },
  description: {
    // marginLeft: 48,
    marginTop: 16,
    marginBottom: 16,
    // width: "80%"
  },
  comment_input: {
    borderRadius: "8px 48px 8px 8px",
  }
};

export default Scene
