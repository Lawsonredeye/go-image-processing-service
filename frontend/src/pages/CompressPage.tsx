import { useState } from 'react';

const CompressPage = () => {
  const [selectedFile, setSelectedFile] = useState<File | null>(null);
  const [quality, setQuality] = useState<number>(75);
  const [compressedImageUrl, setCompressedImageUrl] = useState<string | null>(null);
  const [originalSize, setOriginalSize] = useState<number | null>(null);
  const [compressedSize, setCompressedSize] = useState<number | null>(null);
  const [isLoading, setIsLoading] = useState<boolean>(false);
  const [error, setError] = useState<string | null>(null);

  const handleFileChange = (event: React.ChangeEvent<HTMLInputElement>) => {
    const file = event.target.files?.[0];
    if (file) {
      setSelectedFile(file);
      setOriginalSize(file.size);
      setCompressedImageUrl(null);
      setCompressedSize(null);
      setError(null);
    }
  };

  const handleCompress = async () => {
    if (!selectedFile) {
      setError('Please select a file first.');
      return;
    }
    setIsLoading(true);
    setError(null);
    const formData = new FormData();
    formData.append('image', selectedFile);
    const apiUrl = `http://localhost:8080/api/compress?quality=${quality}`;
    try {
      const response = await fetch(apiUrl, { method: 'POST', body: formData });
      if (!response.ok) {
        const errorText = await response.text();
        throw new Error(errorText || 'An unknown error occurred.');
      }
      const imageBlob = await response.blob();
      const imageUrl = URL.createObjectURL(imageBlob);
      setCompressedImageUrl(imageUrl);
      setCompressedSize(imageBlob.size);
    } catch (err: any) {
      setError(err.message || 'Failed to connect to the server.');
    } finally {
      setIsLoading(false);
    }
  };

  const formatBytes = (bytes: number, decimals = 2) => {
    if (bytes === 0) return '0 Bytes';
    const k = 1024;
    const dm = decimals < 0 ? 0 : decimals;
    const sizes = ['Bytes', 'KB', 'MB', 'GB'];
    const i = Math.floor(Math.log(bytes) / Math.log(k));
    return parseFloat((bytes / Math.pow(k, i)).toFixed(dm)) + ' ' + sizes[i];
  };

  return (
    <div className="card">
      <h2>Image Compressor</h2>
      <p>Select an image (PNG or JPEG) and a quality level to compress it.</p>
      
      <div className="form-group">
        <label htmlFor="file" className="file-input-label">
          <span>{selectedFile ? selectedFile.name : 'Click or drag file to upload'}</span>
        </label>
        <input type="file" id="file" accept="image/png, image/jpeg" onChange={handleFileChange} />
      </div>

      {selectedFile && (
        <div className="form-group">
          <label htmlFor="quality">Compression Quality: {quality}</label>
          <input
            type="range"
            id="quality"
            min="1"
            max="100"
            value={quality}
            onChange={(e) => setQuality(parseInt(e.target.value, 10))}
          />
        </div>
      )}

      <button onClick={handleCompress} className="btn btn-primary" disabled={!selectedFile || isLoading}>
        {isLoading ? 'Compressing...' : 'Compress Image'}
      </button>

      {error && <p className="error-message">Error: {error}</p>}

      {compressedImageUrl && originalSize && compressedSize && (
        <div className="results-container">
          <img src={compressedImageUrl} alt="Compressed result" />
          <div className="results-info">
            <p><strong>Original:</strong> {formatBytes(originalSize)}</p>
            <p><strong>Compressed:</strong> {formatBytes(compressedSize)}</p>
            <p><strong>Reduction:</strong> {(((originalSize - compressedSize) / originalSize) * 100).toFixed(2)}%</p>
            <a href={compressedImageUrl} download={`compressed-${selectedFile?.name}.jpg`} className="btn btn-primary">
              Download
            </a>
          </div>
        </div>
      )}
    </div>
  );
};

export default CompressPage;
