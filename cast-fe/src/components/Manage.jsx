import React, {Component} from 'react';
import {Alert, Button, Col, Container, Form, ProgressBar, Row, Spinner} from "react-bootstrap";
import axios from "axios";
import bsCustomFileInput from 'bs-custom-file-input'
import SidebarProfile from "./SidebarProfile";

class Manage extends Component {
  constructor(props) {
    super(props);
    this.state = {
      title: "",
      description: "",
      tags: "",
      thumbnail: "",
      video: "",
      error_title: "",
      error_description: "",
      error_tags: "",
      error_thumbnail: "",
      error_video: "",
      error_upload: "",
      progress: 0,
      loading: false,
      success: false,
    };
    this.handleChange = this.handleChange.bind(this);
    this.handleChangeFile = this.handleChangeFile.bind(this);
    this.submitForm = this.submitForm.bind(this);
  }

  componentDidMount() {
    bsCustomFileInput.init();
  }

  handleChange(e) {
    this.setState({error_login: ""});
    this.setState({[e.target.name]: e.target.value});
    this.validate(e.target.name, e.target.value);
  }


  handleChangeFile(e) {
    this.setState({error_login: ""});
    this.setState({[e.target.name]: e.target.files[0]});
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
      case "description":
        if (!value.trim()) {
          this.setState({error_description: "Please enter video description"});
          return false;
        }
        this.setState({error_description: ""});
        return true;
      case "tags":
        if (!value.trim()) {
          this.setState({error_tags: "Please enter video tags"});
          return false;
        }
        this.setState({error_tags: ""});
        return true;
      case "thumbnail":
        if (!value) {
          this.setState({error_thumbnail: "Please select thumbnail image"});
          return false;
        }
        this.setState({error_thumbnail: ""});
        return true;
      case "video":
        if (!value) {
          this.setState({error_video: "Please select video file"});
          return false;
        }
        this.setState({error_video: ""});
        return true;
      default:
        return false;
    }
  }

  submitForm(e) {
    e.preventDefault();
    let ok = true;
    ok &= this.validate("title", this.state.title);
    ok &= this.validate("description", this.state.description);
    ok &= this.validate("tags", this.state.tags);
    ok &= this.validate("thumbnail", this.state.thumbnail);
    ok &= this.validate("video", this.state.video);
    if (!ok) return;

    this.setState({loading: true});
    const form = new FormData();
    form.append("title", this.state.title);
    form.append("description", this.state.description);
    form.append("tags", this.state.tags);
    form.append("thumbnail", this.state.thumbnail);
    form.append("video", this.state.video);
    axios.post('/p/video/upload', form, {
        headers: {
          "Access-Control-Allow-Origin": "*",
          "Content-Type": "multipart/form-data"
        },
        onUploadProgress: (progressEvent) => {
          this.setState({progress: (progressEvent.loaded * 100) / progressEvent.total})
        }
      }
    ).then((response) => {
      this.setState({
        title: "",
        description: "",
        tags: "",
        thumbnail: "",
        video: "",
        progress: 0,
        loading: false,
        success: true
      });
      // this.setState({error_upload: "An error has occurred!"});
    }).catch((error) => {
      console.log(error);
      this.setState({loading: false});
    });
  }

  render() {
    return (
      <>
        <Container fluid style={style.content_container}>
          <Row>
            <Col md={2} sm={12} style={{marginBottom: 32}}>
              <SidebarProfile/>
            </Col>
            <Col md={10} sm={12}>
              <h1 style={style.h1}>Upload</h1>
              {this.state.error_upload && <Alert variant={"danger"}>{this.state.error_upload}</Alert>}
              {this.state.success &&
              <Alert variant="success" onClose={() => this.setState({success: false})} dismissible>
                <Alert.Heading>Video Successfully Uploaded!</Alert.Heading>
                <p>
                  Your video is now currently being processed and will be available in a few minutes. Video may first
                  appear in lower quality, as higher resolutions are still being ingested.
                </p>
              </Alert>}
              <Form noValidate autocomplete={"off"} onSubmit={this.submitForm}>
                <Form.Row>
                  <Col>
                    <Form.Group>
                      <Form.Label>Title</Form.Label>
                      <Form.Control name={"title"} value={this.state.title} onBlur={this.handleChange}
                                    onChange={this.handleChange} type={"text"} size={"lg"} style={{fontSize: "2rem"}}
                                    isInvalid={!!this.state.error_title} disabled={this.state.loading}/>
                      <Form.Control.Feedback type={"invalid"}>{this.state.error_title}</Form.Control.Feedback>
                    </Form.Group>
                    <Form.Group>
                      <Form.Label>Description</Form.Label>
                      <Form.Control name={"description"} value={this.state.description} onBlur={this.handleChange}
                                    onChange={this.handleChange} as={"textarea"} rows={7}
                                    isInvalid={!!this.state.error_description} disabled={this.state.loading}/>
                      <Form.Control.Feedback type={"invalid"}>{this.state.error_description}</Form.Control.Feedback>
                    </Form.Group>
                  </Col>
                  <Col>
                    <Form.Group>
                      <Form.Label>Tags</Form.Label>
                      <Form.Control name={"tags"} value={this.state.tags} onBlur={this.handleChange}
                                    onChange={this.handleChange} type={"text"}
                                    isInvalid={!!this.state.error_tags} disabled={this.state.loading}/>
                      <Form.Control.Feedback type={"invalid"}>{this.state.error_tags}</Form.Control.Feedback>
                    </Form.Group>
                    <Form.Group>
                      <Form.Label>Thumbnail</Form.Label>
                      <div className={"custom-file"}>
                        <Form.Control name={"thumbnail"} type={"file"} className={"custom-file-input"}
                                      onChange={this.handleChangeFile} accept="image/*"
                                      isInvalid={!!this.state.error_thumbnail} disabled={this.state.loading}/>
                        <Form.Control.Feedback type={"invalid"}>{this.state.error_thumbnail}</Form.Control.Feedback>
                        <Form.Label className="custom-file-label">{this.state.thumbnail.name}</Form.Label>
                      </div>
                    </Form.Group>
                    <Form.Group>
                      <Form.Label>Video</Form.Label>
                      <div className={"custom-file"}>
                        <Form.Control name={"video"} type={"file"} className={"custom-file-input"}
                                      onChange={this.handleChangeFile} accept="video/*"
                                      isInvalid={!!this.state.error_video} disabled={this.state.loading}/>
                        <Form.Control.Feedback type={"invalid"}>{this.state.error_video}</Form.Control.Feedback>
                        <Form.Label className="custom-file-label">{this.state.video.name}</Form.Label>
                      </div>
                    </Form.Group>
                    <ProgressBar animated now={this.state.progress} style={style.progress}/>
                    <Button variant="primary" type="submit" block disabled={this.state.loading}>
                      Upload{" "}
                      {this.state.loading &&
                      <Spinner style={{verticalAlign: "initial"}} as="span" animation="grow"
                               size="sm" role="status" aria-hidden="true"/>}
                    </Button>
                  </Col>
                </Form.Row>
              </Form>
              <hr/>
              <h1 style={style.h1}>Manage</h1>
              Manage card
            </Col>
          </Row>
        </Container>
      </>
    )
  }
}

let style = {
  h1: {
    fontFamily: "Comfortaa",
  },
  content_container: {
    padding: "36px 0 0 0",
  },
  progress: {
    marginBottom: 16
  }
};

export default Manage