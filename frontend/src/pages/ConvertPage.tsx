import { useState } from 'react';

const Converter = () => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [targetFormat, setTargetFormat] = useState<string>('jpeg');
  const [convertedImageUrl, setConvertedImageUrl] = useState<string | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      setConvertedImageUrl(null);
      setError(null);
    }
  };

  const handleConvert = async () => {
    if (!selectedFile) {
      setError('Please select a file first.');
      return;
    }
    setIsLoading(true);
    setError(null);

    const formData = new FormData();
    formData.append('image', selectedFile);

    const apiUrl = `http://localhost:8080/api/convert?format=${targetFormat}`;

    try {
      const response = await fetch(apiUrl, { method: 'POST', body: formData });
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'Conversion failed.');
      }
      const imageBlob = await response.blob();
      const imageUrl = URL.createObjectURL(imageBlob);
      setConvertedImageUrl(imageUrl);
    } catch (err: any) {
      setError(err.message || 'Failed to connect to the server.');
    } finally {
      setIsLoading(false);
    }
  };

  return (
    <div className="card">
      <h2>Image Converter</h2>
      <p>Select an image to convert it to JPEG or PNG format.</p>

      <div className="form-group">
        <label htmlFor="file" className="file-input-label">
          <span>{selectedFile ? selectedFile.name : 'Click or drag file to upload'}</span>
        </label>
        <input type="file" id="file" accept="image/png, image/jpeg" onChange={handleFileChange} />
      </div>

      {selectedFile && (
        <div className="form-group">
          <label>Target Format</label>
          <div style={{ display: 'flex', gap: '2rem' }}>
            <label htmlFor="jpeg">
              <input
                type="radio"
                id="jpeg"
                name="format"
                value="jpeg"
                checked={targetFormat === 'jpeg'}
                onChange={(e) => setTargetFormat(e.target.value)}
              />
              JPEG
            </label>
            <label htmlFor="png">
              <input
                type="radio"
                id="png"
                name="format"
                value="png"
                checked={targetFormat === 'png'}
                onChange={(e) => setTargetFormat(e.target.value)}
              />
              PNG
            </label>
          </div>
        </div>
      )}

      <button onClick={handleConvert} className="btn btn-primary" disabled={!selectedFile || isLoading}>
        {isLoading ? 'Converting...' : 'Convert Image'}
      </button>

      {error && <p className="error-message">Error: {error}</p>}

      {convertedImageUrl && (
        <div className="results-container">
          <img src={convertedImageUrl} alt="Converted result" />
          <div className="results-info">
            <p><strong>Conversion Complete</strong></p>
            <a href={convertedImageUrl} download={`converted-${selectedFile?.name}.${targetFormat}`} className="btn btn-primary">
              Download
            </a>
          </div>
        </div>
      )}
    </div>
  );
};

export default Converter;
