import React, { Component } from "react";
import "./Message.scss";

const Image = ({ fileURL, fileName }) => (
  <div>
    <img
      className="Message__image"
      src={fileURL}
      alt={fileName}
    />
    <a href={fileURL} download={fileName}>Download</a>
  </div>
);

class Message extends Component {
  constructor(props) {
    super(props);
    this.state = {
      message: this.props.message,
      attachment: this.props.attachment,
    };
  }

  render() {
    return (
      <div className="Message">
        {this.state.message}
        {this.state.attachment?.fileURL && <Image fileURL={this.state.attachment.fileURL} fileName={this.state.attachment.fileName} />}
      </div>
    );
  }
}

export default Message;
