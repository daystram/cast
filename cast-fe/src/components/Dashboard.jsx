import React, {Component} from 'react';
import {
  Badge,
  Button,
  Card,
  Col,
  Container,
  Form,
  InputGroup,
  Modal,
  OverlayTrigger,
  Popover,
  Row,
  Spinner
} from "react-bootstrap";
import SidebarProfile from "./SidebarProfile";
import Chat from "./Chat";
import auth from "../helper/auth";
import CastEditable from "./CastEditable";
import axios from "axios";
import urls from "../helper/url";
import format from "../helper/format";

class Dashboard extends Component {
  constructor(props) {
    super(props);
    this.state = {
      loading: true,
    };
    document.title = "Dashboard | cast";
    this.controlStream = this.controlStream.bind(this);
    this.fetchDetail(auth().username())
  }

  fetchDetail(hash) {
    axios.get(urls().cast_details(), {
      params: {
        hash: hash,
        username: auth().username()
      }
    }).then((response) => {
      this.setState({loading: false});
      if (response.data.code === 200) {
        let data = response.data.data;
        this.setState({stream: data, likes: data.likes});
        console.log(data)
      }
    }).catch((error) => {
      console.log(error);
      this.setState({loading: false});
    });
  }

  controlStream() {
    if (this.state.stream && this.state.stream.is_live) {
      // TODO: stop stream
    } else {
      if (this.state.waiting) {
        // TODO: close window
        this.setState({waiting: false})
      } else {
        // TODO: open window
        this.setState({waiting: true})
      }
    }
  }

  handleChange(e) {
    this.setState({error_edit: ""});
    this.setState({[e.target.name]: e.target.value});
    this.validate(e.target.name, e.target.value);
  }

  validate(field, value) {
    switch (field) {
      case "username":
        if (!value.trim()) {
          this.setState({error_username: "Please enter your username"});
          return false;
        }
        this.setState({error_username: ""});
        return true;
      case "password":
        if (!value) {
          this.setState({error_password: "Please enter your password"});
          return false;
        }
        this.setState({error_password: ""});
        return true;
      default:
        return false;
    }
  }

  render() {
    return (
      <>
        <Container fluid style={style.content_container}>
          <Row>
            <Col xl={{span: 2, order: 1}} sm={{span: 6, order: 2}} xs={{span: 12, order: 1}} style={{marginBottom: 32}}>
              <SidebarProfile/>
            </Col>
            <Col xl={{span: 8, order: 2}} sm={{span: 12, order: 1}} xs={{span: 12, order: 2}}>
              <div style={{display: "flex", justifyContent: "space-between", alignItems: "center"}}>
                <h1 style={style.h1}>Livestream {this.state.stream && this.state.stream.is_live &&
                <Badge pill style={style.live_tag}>LIVE</Badge>}</h1>
                {this.state.stream && <Button
                  variant={this.state.stream && (this.state.stream.is_live ? "danger" : this.state.waiting ? "warning" : "success")}
                  size={"md"} style={style.live_button} onClick={this.controlStream}>
                  {this.state.stream && (this.state.stream.is_live ? "End Stream" : this.state.waiting ?
                    <>Waiting for Up-link{"  "}<Spinner animation="grow" style={style.spinner}/></> : "Go Live")}
                </Button>}
              </div>
              <div style={style.card_detail}>
                {this.state.stream && <CastEditable video={this.state.stream}/>}
              </div>
              <h1 style={style.h1}>Metrics</h1>
              {this.state.stream && <Card body style={style.card_detail}>
                <Row>
                  <Col sm={6}>
                    <p
                      style={style.show_count}>{this.state.stream.is_live ? "?" : format().full_date(this.state.stream.created_at)}</p>
                    <p style={style.show_caption}>{this.state.stream.is_live ? "stream duration" : "last stream"}</p>
                  </Col>
                  <Col sm={6}>
                    <p style={style.show_count}>2374</p>
                    <p style={style.show_caption}>{this.state.stream.is_live || "last"} total
                      viewer{this.state.views === 1 ? "" : "s"}</p>
                  </Col>
                </Row>
              </Card>}
              <h1 style={style.h1}>Connection</h1>
              <Card body style={style.card_detail}>
                <Row>
                  <Col md={6} sm={12}>
                    <Form autocomplete="off">
                      <Form.Group>
                        <Form.Label>Server Address</Form.Label>
                        <InputGroup>
                          <Form.Control type="text" value={`rtmp://cast.daystram.com/live`}
                                        ref={ref => this.rtmpField = ref}/>
                          <InputGroup.Append>
                            <OverlayTrigger trigger="click" placement="top" overlay={(
                              <Popover id="popover-basic">
                                <Popover.Content>Copied!</Popover.Content>
                              </Popover>)}>
                              <Button variant="outline-primary" onClick={() => {
                                this.rtmpField.select();
                                document.execCommand("copy");
                                this.rtmpField.blur();
                              }}>Copy</Button>
                            </OverlayTrigger>
                          </InputGroup.Append>
                        </InputGroup>
                      </Form.Group>
                    </Form>
                  </Col>
                  <Col md={6} sm={12}>
                    <Form autocomplete="off">
                      <Form.Group>
                        <Form.Label>Stream Key</Form.Label>
                        <InputGroup>
                          <Form.Control type={this.state.show_key ? "text" : "password"} value={auth().username()}/>
                          <InputGroup.Append>
                            <Button variant={this.state.show_key ? "primary" : "outline-primary"} onClick={() => {
                              this.setState({show_key: !this.state.show_key})
                            }}><span
                              className="material-icons">{this.state.show_key ? "visibility_off" : "visibility"}</span></Button>
                          </InputGroup.Append>
                          <InputGroup.Append>
                            <OverlayTrigger trigger="click" placement="top" overlay={(
                              <Popover id="popover-basic">
                                <Popover.Content>Copied!</Popover.Content>
                              </Popover>)}>
                              <Button variant="outline-primary" onClick={() => {
                                navigator.clipboard.writeText(auth().username());
                              }}>Copy</Button>
                            </OverlayTrigger>
                          </InputGroup.Append>
                        </InputGroup>
                      </Form.Group>
                    </Form>
                  </Col>
                  <Col sm={12}>
                    <Form autocomplete="off">
                      <Form.Group>
                        <Form.Label>Chat Embed</Form.Label>
                        <InputGroup>
                          <Form.Control type="text" value={`https://cast.daystram.com/c/${auth().username()}`}
                                        ref={ref => this.chatField = ref}/>
                          <InputGroup.Append>
                            <OverlayTrigger trigger="click" placement="top" overlay={(
                              <Popover id="popover-basic">
                                <Popover.Content>Copied!</Popover.Content>
                              </Popover>)}>
                              <Button variant="outline-primary" onClick={() => {
                                this.chatField.select();
                                document.execCommand("copy");
                                this.chatField.blur();
                              }}>Copy</Button>
                            </OverlayTrigger>
                          </InputGroup.Append>
                        </InputGroup>
                      </Form.Group>
                    </Form>
                  </Col>
                </Row>
              </Card>
            </Col>
            <Col xl={{span: 2, order: 3}} sm={{span: 6, order: 3}} xs={{span: 12, order: 3}}>
              <Chat height={"80vh"} embedded={true} hash={auth().username()}/>
            </Col>
          </Row>
        </Container>
        <Modal show={this.state.prompt_stop} size={"md"} onHide={() => this.setState({prompt_stop: false})} centered>
          <Modal.Header closeButton>
            <Modal.Title id="contained-modal-title-vcenter">Ending Stream</Modal.Title>
          </Modal.Header>
          <Modal.Body>
            <p>Are you sure you want to stop the stream?</p>
          </Modal.Body>
          <Modal.Footer>
            <Button variant={"info"}>Cancel</Button>
            <Button variant={"danger"}>End Stream</Button>
          </Modal.Footer>
        </Modal>
      </>
    )
  }
}

let style = {
  h1: {
    fontFamily: "Comfortaa",
    display: "inline-block"
  },
  content_container: {
    padding: "36px 0 0 0",
  },
  live_tag: {
    background: "red",
    color: "white",
    marginRight: 8,
    marginBottom: 8,
    borderRadius: 8,
    fontSize: 24,
    fontWeight: 400,
    verticalAlign: "middle"
  },
  live_button: {
    verticalAlign: "top",
    fontWeight: 600,
    height: "fit-content"
  },
  spinner: {
    width: 16,
    height: 16,
    verticalAlign: "initial"
  },
  card_detail: {
    borderRadius: "8px 48px 8px 8px",
    marginBottom: 32
  },
  show_count: {
    textAlign: "center",
    fontSize: 28,
    fontWeight: 800,
    marginBottom: 8
  },
  show_caption: {
    textAlign: "center",
    marginBottom: 0,
    fontWeight: 600
  }
};

export default Dashboard
