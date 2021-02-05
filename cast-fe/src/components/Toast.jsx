import React from "react";
import { useHistory } from "react-router-dom";
import { Button, ButtonGroup, Image } from "react-bootstrap";

import "react-toastify/dist/ReactToastify.min.css";
import "./toast.css";

function Toast(props) {
  let history = useHistory();
  return (
    <>
      <div style={{ display: "flex", alignItems: "center" }}>
        {props.data.username && (
          <Image
            src={""}
            height={42}
            width={42}
            roundedCircle
            style={{ marginRight: 8 }}
          />
        )}
        <p style={{ fontSize: 14, margin: 0, height: "fit-content" }}>
          {props.data.message}
        </p>
      </div>
      <ButtonGroup size={"sm"} style={{ width: "100%", marginTop: 8 }}>
        {props.data.hash && (
          <Button
            variant={"outline-primary"}
            onClick={() => history.push("/w/" + props.data.hash)}
          >
            Watch
          </Button>
        )}
        <Button variant={"outline-primary"}>Close</Button>
      </ButtonGroup>
    </>
  );
}

export default Toast;
