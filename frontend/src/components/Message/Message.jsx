import { Component } from 'react';
import './Message.scss';

const Image = ({ fileURL, fileName }) => {
  const isImage = fileName && /\.(jpg|jpeg|png|gif|bmp|svg|webp)$/i.test(fileName);
  return (
    <div>
      {isImage ? (
        <img
          className="Message__image"
          src={fileURL}
          alt={fileName}
        />
      ) : (
        <div className="Message__file-icon">
          <img src="/file-icon.svg" alt="File" width="16" height="16" />
          <span>{fileName}</span>
        </div>
      )}
      <a href={fileURL} download={fileName}>Download</a>
    </div>
  );
};

class Message extends Component {
  constructor(props) {
    super(props);
    const SENDERS = {
      1: 'me',
      2: 'other',
    };
    this.state = {
      message: this.props.message,
      attachment: this.props.attachment,
      sender: SENDERS[this.props.type] || '',
    };
  }

  render() {
    return (
      <div className={ `Message ${this.state.sender}` }>
        {this.state.message}
        {this.state.attachment?.fileURL && <Image fileURL={this.state.attachment.fileURL} fileName={this.state.attachment.fileName} />}
      </div>
    );
  }
}

export default Message;
