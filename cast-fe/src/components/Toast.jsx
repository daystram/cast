import React from "react";
import { useHistory } from "react-router-dom";
import { Button, ButtonGroup, Image } from "react-bootstrap";
import { ProfileImage } from "./index";
import api from "../apis/api";

import "react-toastify/dist/ReactToastify.min.css";
import "../styles/toast.css";

function Toast(props) {
  let history = useHistory();
  return (
    <>
      <div style={{ display: "flex", alignItems: "center" }}>
        {props.data.hash ? (
          <Image
            src={api.cdn.thumbnail(props.data.hash)}
            style={{
              borderRadius: "4px 24px 4px 4px",
              objectFit: "cover",
              height: "42px",
              marginRight: 8,
            }}
          />
        ) : (
          props.data.name && (
            <ProfileImage
              size={42}
              name={props.data.name}
              style={{ marginRight: 8 }}
            />
          )
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
