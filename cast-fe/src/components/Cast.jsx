import React from 'react';
import {Badge, Image} from "react-bootstrap";
import {useHistory} from "react-router-dom";
import abbreviate from '../helper/abbreviate'

function Cast(props) {
  const history = useHistory();

  function playVideo(e) {
    history.push("/w/" + props.video._id);
    console.log(`Play: ${props.video._id}`)
  }

  function viewAuthor(e) {
    e.stopPropagation();
    console.log(`View: ${props.video.author.name}`)
  }

  return (
    <div style={style.cast_card} onClick={playVideo}>
      <Image src={props.video.thumbnail}
             style={style.cast_thumbnail}/>
      <div style={style.cast_tag_bar}>
        {props.video.isLive && <Badge pill style={style.cast_live_tag}>LIVE</Badge>}
        <Badge pill
               style={style.cast_viewer_tag}>{abbreviate(props.video.views)} {props.video.isLive ? 'viewers' : 'views'}</Badge>
      </div>
      <div style={style.cast_detail}>
        <Image src={props.video.author.image} height={42} width={42}
               style={style.profile_image} onClick={viewAuthor} roundedCircle/>
        <div style={style.cast_author_details}>
          <p style={style.cast_title}>{props.video.title}</p>
          <p style={style.cast_author} onClick={viewAuthor}>{props.video.author.name}</p>
          <p style={style.cast_duration}>Streaming for 3 hours</p>
        </div>
      </div>
    </div>
  )
}

let style = {
  cast_card: {},
  cast_thumbnail: {
    borderRadius: "8px 48px 8px 8px",
    borderWidth: 1,
    borderColor: "lightgray",
    borderStyle: "solid",
    objectFit: "cover",
    width: "100%",
  },
  cast_tag_bar: {
    marginTop: 8,
    marginBottom: 8,
  },
  cast_live_tag: {
    background: "red",
    color: "white",
    marginRight: 8,
    borderRadius: 8,
    fontSize: 16,
    fontWeight: 400
  },
  cast_viewer_tag: {
    background: "white",
    borderRadius: 8,
    borderWidth: 1,
    borderColor: "lightgray",
    borderStyle: "solid",
    fontSize: 16,
    fontWeight: 400
  },
  cast_detail: {
    display: "flex",
  },
  cast_author_details: {
    display: "grid"
  },
  cast_title: {
    fontSize: 18,
    fontWeight: 600,
    margin: 0
  },
  cast_author: {
    fontSize: 16,
    fontWeight: 400,
    margin: 0
  },
  cast_duration: {
    fontSize: 16,
    fontWeight: 300,
    margin: 0,
    whiteSpace: "nowrap",
    overflow: "hidden",
    textOverflow: "ellipsis",
  },
  profile_image: {
    marginRight: 8,
    alignSelf: "start"
  }
};

export default Cast
