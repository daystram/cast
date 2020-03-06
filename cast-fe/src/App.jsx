import React from 'react';
import './App.css';
import {HybridPlayer} from "./components/player";

function App() {
  return (
    <HybridPlayer
      url={"http://cast.daystram.com/api/test_stream"}
      live={true}/>
  );
}

export default App;
