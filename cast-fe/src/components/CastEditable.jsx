import React, {Component} from 'react';
import {Badge, Button, Card, Col, Form, Image, Row, Spinner} from "react-bootstrap";
import urls from "../helper/url";
import format from "../helper/format";
import axios from "axios";
import {Prompt} from "react-router-dom";

const resolutions = ["Processing", "180p", "360p", "480p", "720p", "1080p"];

class CastEditable extends Component {
  constructor(props) {
    super(props);
    this.state = {
      title: props.video.title,
      tags: "",
      description: props.video.description,
      error_title: "",
      error_tags: "",
      error_description: "",
      error_edit: "",
      before: {},
      editing: false,
      loading: false
    };
    this.pressEdit = this.pressEdit.bind(this);
    this.pressDelete = this.pressDelete.bind(this);
    this.handleChange = this.handleChange.bind(this);
    this.submitForm = this.submitForm.bind(this);
  }

  pressEdit() {
    if (this.state.editing) {
      this.submitForm()
    } else {
      this.setState({
        before: {
          title: this.state.title,
          tags: this.state.tags,
          description: this.state.description,
        },
        editing: true
      })
    }
  }

  pressDelete() {
    if (this.state.editing) {
      this.setState({
        title: this.state.before.title,
        tags: this.state.before.tags,
        description: this.state.before.description,
        editing: false
      })
    } else {
      console.log("delete")
    }
  }

  handleChange(e) {
    this.setState({error_edit: ""});
    this.setState({[e.target.name]: e.target.value});
    this.validate(e.target.name, e.target.value);
  }

  validate(field, value) {
    switch (field) {
      case "title":
        if (!value.trim()) {
          this.setState({error_title: "Please enter video title"});
          return false;
        }
        this.setState({error_title: ""});
        return true;
      case "tags":
        if (!value.trim()) {
          this.setState({error_tags: "Please enter video tags"});
          return false;
        }
        this.setState({error_tags: ""});
        return true;
      case "description":
        if (!value.trim()) {
          this.setState({error_description: "Please enter video description"});
          return false;
        }
        this.setState({error_description: ""});
        return true;
      default:
        return false;
    }
  }

  submitForm() {
    let ok = true;
    ok &= this.validate("title", this.state.title);
    ok &= this.validate("tags", this.state.tags);
    ok &= this.validate("description", this.state.description);
    if (!ok) return;

    this.setState({loading: true});
    axios.post(urls().upload(), {}
    ).then((response) => {
      this.setState({editing: false, loading: false});
      this.fetchVideos();
    }).catch((error) => {
      console.log(error);
      this.setState({error_edit: "An error has occurred! Please try again.", loading: false});
    });

  }

  render() {
    return (
      <Card body style={style.card}>
        <Row>
          <Col xl={3} lg={4} md={5} sm={12} className={"responsive-fold"}>
            <Image src={urls().thumbnail(this.props.video.hash)} style={style.thumbnail}/>
          </Col>
          <Col md sm={12} className={"responsive-fold"}>
            {this.state.editing ?
              <Form.Group>
                <Form.Control name={"title"} value={this.state.title} onBlur={this.handleChange}
                              onChange={this.handleChange} type={"text"} size={"lg"} style={style.title}
                              isInvalid={this.state.error_title}/>
                <Form.Control.Feedback type={"invalid"}>{this.state.error_title}</Form.Control.Feedback>
              </Form.Group> :
              <h1 style={style.title}>{this.state.title}</h1>
            }
            <p style={{marginTop: 4, marginBottom: 0}}>{format().full_date(this.props.video.created_at)}</p>
            <div style={style.cast_tag_bar}>
              <div>
                <Badge pill style={style.cast_tag_resolution}>{resolutions[this.props.video.resolutions]}
                  {" "} {this.props.video.resolutions === 5 ? "" :
                    <Spinner animation="grow" style={style.spinner}/>}</Badge>
                <Badge pill style={style.cast_tag}>tag</Badge>
                <Badge pill style={style.cast_tag}>another</Badge>
              </div>
            </div>
            {this.state.editing ?
              <Form.Group>
                <Form.Control name={"description"} value={this.state.description} onBlur={this.handleChange}
                              onChange={this.handleChange} as={"textarea"} size={"lg"} style={style.description}
                              isInvalid={this.state.error_description} rows={5}/>
                <Form.Control.Feedback type={"invalid"}>{this.state.error_description}</Form.Control.Feedback>
              </Form.Group> :
              <p style={style.description}>{this.state.description}</p>
            }
          </Col>
          <Col md={"auto"} sm={12}>
            <Button variant={"success"} block size={"sm"} onClick={this.pressEdit} style={style.button}>
              <span className="material-icons">{this.state.editing ? "check" : "edit"}</span></Button>
            <Button variant={this.state.editing ? "info" : "danger"} block size={"sm"}
                    onClick={this.pressDelete} style={style.button}>
              <span className="material-icons">{this.state.editing ? "clear" : "delete"}</span></Button>
          </Col>
        </Row>
        <Prompt
          when={this.state.editing || this.state.loading}
          message={() =>`Are you sure? You changes will be discarded.`}/>
      </Card>
    )
  }
}

let style = {
  card: {
    marginBottom: 16,
    width: "100%"
  },
  title: {
    margin: 0,
    fontSize: "2rem"
  },
  thumbnail: {
    borderRadius: "8px 48px 8px 8px",
    borderWidth: 1,
    borderColor: "lightgray",
    borderStyle: "solid",
    objectFit: "cover",
    width: "100%",
  },
  created_at: {
    fontSize: 16,
    fontWeight: 300,
    margin: 0,
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
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
    fontSize: 14,
    fontWeight: 400,
    marginRight: 8
  },
  cast_tag_resolution: {
    background: "#E84409",
    color: "white",
    marginRight: 8,
    borderRadius: 8,
    fontSize: 14,
    fontWeight: 400
  },
  description: {
    fontSize: 16,
    webkitLineClamp: "3",
    overflow: "hidden",
    webkitBoxOrient: "vertical",
    display: "-webkit-box",
  },
  spinner: {
    width: 10,
    height: 10,
    verticalAlign: "initial"
  }
};

export default CastEditable
