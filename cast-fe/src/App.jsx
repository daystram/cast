import React from 'react';
import './App.css';
import {HybridPlayer} from "./components/player";

function App() {
  return (
    <HybridPlayer
      url={"https://storage.googleapis.com/cast-uploaded-videos/HASH03/manifest.mpd"}
      live={false}/>
  );
}

export default App;
