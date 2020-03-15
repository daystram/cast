import React, {Component} from 'react';
import {Container} from "react-bootstrap";

class SignUp extends Component {
  render() {
    return (
      <>
        <Container fluid style={style.content_container}>
          SignUp
        </Container>
      </>
    )
  }
}

let style = {
  content_container: {
    padding: "36px 0 0 0"
  },
};

export default SignUp
