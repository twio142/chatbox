import React, { Component } from "react";
import "./ChatHistory.scss";
import Message from "../Message";

class ChatHistory extends Component {
  render() {
    const messages = this.props.chatHistory.map((msg, index) => {
      const {text, fileName, fileURL, type} = JSON.parse(msg.data);
      return (
        <Message key={index} message={text} attachment={{fileName, fileURL}} type={type} />
      );
    });

    return (
      <div className="ChatHistory">
        {messages}
      </div>
    );
  }
}

export default ChatHistory;
