import React, { Component } from "react";
import "./App.css";
import { connect, sendMsg } from "./api";
import Header from "./components/Header/Header";
import ChatHistory from "./components/ChatHistory/ChatHistory";
import ChatInput from "./components/ChatInput/ChatInput";
import FileUpload from "./components/FileUpload/FileUpload";

class App extends Component {
  constructor(props) {
    super(props);
    this.state = {
      chatHistory: []
    }
  }

  componentDidMount() {
    connect((msg) => {
      console.log("New Message", msg);
      this.setState(() => ({
        chatHistory: [...this.state.chatHistory, msg]
      }));
    });
  }

  send(e) {
    if (e.keyCode === 13) {
      sendMsg({text: e.target.value});
      e.target.value = "";
    }
  }

  render() {
    return (
      <div className="App">
        <Header />
        <ChatHistory chatHistory={this.state.chatHistory} />
        <ChatInput send={this.send} />
        <FileUpload sendMsg={sendMsg} />
      </div>
    );
  }
}

export default App;
