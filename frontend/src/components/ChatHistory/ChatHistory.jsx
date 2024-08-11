import React, { Component } from "react";
import "./ChatHistory.scss";
import Message from "../Message";

class ChatHistory extends Component {
  render() {
    const messages = this.props.chatHistory.map((msg, index) => {
      let {text, fileName, fileURL} = JSON.parse(msg.data);
      return (
        <Message key={index} message={text} attachment={{fileName, fileURL}} />
      );
    });

    return (
      <div className="ChatHistory">
        <h2>Chat History</h2>
        {messages}
      </div>
    );
  }
}

export default ChatHistory;