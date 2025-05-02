import React, { Component } from "react";
import "./FileUpload.scss";

class FileUpload extends Component {
  state = {
    uploading: false,
    progress: 0,
    useChunkedUpload: false
  };

  handleSubmit = async (event) => {
    event.preventDefault();
    this.setState({ uploading: true });

    const fileInput = event.target.querySelector('input[type="file"]');
    const file = fileInput.files[0];

    if (!file) {
      this.setState({ uploading: false });
      return;
    }

    // Use chunked upload for files larger than 100MB
    const useChunked = file.size > 100 * 1024 * 1024;
    this.setState({ useChunkedUpload: useChunked });

    if (useChunked) {
      await this.uploadChunked(file);
    } else {
      await this.uploadStandard(event);
    }

    event.target.reset();
    this.setState({ uploading: false });
  }

  uploadStandard = async (event) => {
    const formData = new FormData(event.target);
    const response = await fetch(`http://${window.API_URL}/upload`, {
      method: "POST",
      body: formData,
    });
    let { fileURL, fileName } = await response.json();
    fileURL = `http://${window.API_URL}${fileURL}`;
    this.props.sendMsg({fileURL, fileName});
  }

  uploadChunked = async (file) => {
    const CHUNK_SIZE = 2 * 1024 * 1024; // 2MB chunks
    const totalChunks = Math.ceil(file.size / CHUNK_SIZE);
    const fileID = Date.now().toString(); // Simple unique ID

    for (let chunkNumber = 0; chunkNumber < totalChunks; chunkNumber++) {
      const start = chunkNumber * CHUNK_SIZE;
      const end = Math.min(file.size, start + CHUNK_SIZE);
      const chunk = file.slice(start, end);

      const formData = new FormData();
      formData.append('chunk', chunk);
      formData.append('chunkNumber', chunkNumber.toString());
      formData.append('totalChunks', totalChunks.toString());
      formData.append('filename', file.name);
      formData.append('fileID', fileID);

      try {
        const response = await fetch(`http://${window.API_URL}/upload/chunk`, {
          method: 'POST',
          body: formData,
        });

        if (!response.ok) {
          throw new Error(`Upload failed: ${response.status}`);
        }

        const result = await response.json();

        // Update progress
        this.setState({ 
          progress: Math.round(((chunkNumber + 1) / totalChunks) * 100) 
        });

        // Check if upload is complete
        if (result.status === 'complete') {
          const fileURL = `http://${window.API_URL}${result.fileURL}`;
          this.props.sendMsg({ fileURL, fileName: result.fileName });
          break;
        }
      } catch (error) {
        console.error('Error uploading chunk:', error);
        // You could implement retry logic here
        break;
      }
    }
  }

  render() {
    const { uploading, progress, useChunkedUpload } = this.state;

    return (
      <div>
        <form action="/upload" method="post" encType="multipart/form-data" onSubmit={this.handleSubmit}>
          <input type="file" name="file" id="file" disabled={uploading} />
          <input type="submit" value="Upload File" disabled={uploading} />
        </form>

        {uploading && (
          <div className="upload-progress">
            {useChunkedUpload ? (
              <div>
                <progress value={progress} max="100" />
                <span>{progress}% (Chunked Upload)</span>
              </div>
            ) : (
              <div>Uploading...</div>
            )}
          </div>
        )}
      </div>
    );
  }
}

export default FileUpload;
