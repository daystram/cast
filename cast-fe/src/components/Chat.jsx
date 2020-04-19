import React, {Component} from 'react';
import {Button, Card, Form, FormControl, InputGroup} from "react-bootstrap";
import urls from "../helper/url";
import auth from "../helper/auth";
import {withRouter} from "react-router-dom";

class Chat extends Component {
  constructor(props) {
    super(props);
    this.state = {
      connection: undefined,
      chat: "",
      loading: true,
      chats: [],
      error_chat: "",
    };
    this.connectChat = this.connectChat.bind(this);
    this.handleSubmit = this.handleSubmit.bind(this);
  }

  componentDidMount() {
    this.connectChat(this.props.embedded ? this.props.hash : this.props.match.params.hash)
  }

  componentDidUpdate(prevProps, prevState, snapshot) {
    if ((this.props.embedded && (this.props.hash !== prevProps.hash)) ||
      (!this.props.embedded && (this.props.match.params.hash !== prevProps.match.params.hash))) {
      this.state.connection.close(1000);
      this.setState({chats: []});
      this.connectChat(this.props.embedded ? this.props.hash : this.props.match.params.hash)
    }
  }

  connectChat(hash) {
    window.WebSocket = window.WebSocket || window.MozWebSocket;
    if (!window.WebSocket) {
      console.log("WebSocket not supported!");
      return;
    }
    const connection = new WebSocket(urls().chat(hash, auth().token()));
    connection.onopen = () => {
      console.log("Live chat connected!");
      this.setState({connection})
    };
    connection.onerror = () => {
      console.log("Cannot connect to server!");
      if (!this.props.embedded) this.props.history.push("/")
    };
    connection.onclose = () => {
      console.log("Server disconnected!");
      if (!this.props.embedded) this.props.history.push("/")
    };
    connection.onmessage = (message) => {
      try {
        let json = JSON.parse(message.data);
        if (json && json.type === "chat") {
          this.setState({chats: [...this.state.chats, json.data]});
        } else {
          console.log("Invalid JSON: ", message.data);
        }
      } catch (e) {
        console.log("Invalid JSON: ", message.data);
      }
    };
  }

  handleSubmit(e) {
    e.preventDefault();
    if (!this.state.chat.trim()) return;
    if (!auth().is_authenticated()) {
      this.props.promptSignup();
      return
    }
    this.state.connection.send(JSON.stringify({
      type: "chat",
      data: this.state.chat.trim()
    }));
    this.setState({chat: ""});
  }

  render() {
    return (
      <Card style={this.props.embedded ? style.live_chat : style.standalone}>
        <Card.Body style={style.live_chat_body}>
          {this.state.chats.length !== 0 && this.state.chats.map(chat => (
            <p style={style.live_chat_item}><b>{chat.author}</b>: {chat.chat}</p>
          ))}
          {this.props.embedded &&
          <Form onSubmit={this.handleSubmit}>
            <InputGroup style={style.live_chat_input}>
              <FormControl type="text" value={this.state.chat} placeholder="Chat"
                           onChange={e => this.setState({chat: e.target.value})}/>
              <InputGroup.Append>
                <Button variant="outline-primary" type={"submit"}><i className="material-icons">send</i></Button>
              </InputGroup.Append>
            </InputGroup>
          </Form>
          }
        </Card.Body>
      </Card>
    )
  }
}

let style = {
  standalone: {
    overflow: "hidden",
    background: "#00000000",
    width: 360,
    border: "none",
    color: "white",
    fontSize: 24
  },
  live_chat: {
    borderRadius: "8px 48px 8px 8px",
    overflow: "hidden",
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
};

export default withRouter(Chat)
