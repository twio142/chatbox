import React, { Component } from "react";
import "./FileUpload.scss";

class FileUpload extends Component {
  handleSubmit = async (event) => {
    event.preventDefault();
    const formData = new FormData(event.target);
    const response = await fetch(`http://${window.API_URL}/upload`, {
        method: "POST",
        body: formData,
    });
    let { fileURL, fileName } = await response.json();
    fileURL = `http://${window.API_URL}${fileURL}`;
    this.props.sendMsg({fileURL, fileName});
  }
  render() {
    return (
      <form action="/upload" method="post" enctype="multipart/form-data" onSubmit = {this.handleSubmit}>
        <input type="file" name="file" id="file" />
        <input type="submit" value="Upload File" />
      </form>
    );
  }
}

export default FileUpload;
