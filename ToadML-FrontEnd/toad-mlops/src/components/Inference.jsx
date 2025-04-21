import "../styles/Sections.css"
import "../styles/Inference.css"
import { useState } from "react";

export default function Inference({ addInferenceToHistory }) {
  const [image, setImage] = useState(null);
  const [previewUrl, setPreviewUrl] = useState(null);
  const [result, setResult] = useState(null);

  if (!previewUrl) {
    setPreviewUrl("/default-image-5-1.jpg");
  }

  const handleImageUpload = (e) => {
    const file = e.target.files[0];
    setImage(file);
    setPreviewUrl(URL.createObjectURL(file));
    setResult(null);
  };

  const handleSubmit = async () => {
    if (!image) return;

    const token = localStorage.getItem("jwt");
    if (!token) {
        setError("Debe iniciar sesión.");
        return;
    }

    const formData = new FormData();
    formData.append("file", image);

    try {
      const response = await fetch("http://3.17.179.120:8080/predict", {
        method: "POST",
        headers: {Authorization: `Bearer ${token}`},
        body: formData,
      });

      const data = await response.json();
      setResult(data);

      // Agregar la inferencia al historial
      const newInference = {
        img: previewUrl,
        imgfile: image.name,
        result: data.class,
        accuracy: data.accuracy,
      };
      addInferenceToHistory(newInference);
    } catch (error) {
      console.error("Error al hacer la inferencia:", error);

      setResult({ class: "Error", accuracy: 0 });
      const newInference = {
        img: previewUrl,
        imgfile: image.name,
        result: "Error",
        accuracy: 0,
      };
      addInferenceToHistory(newInference);
    }
  };

  return (
    <section id="inference" style={{marginTop: "16px"}}>
      <h2>Inference</h2>
      <input id="files" type="file" accept="image/*" onChange={handleImageUpload} />
      {previewUrl && (
        <div style={{ marginTop: "16px" }}>
          <img
            src={previewUrl}
            id="preview-img"
            alt="preview"
            style={{ width: "400px", height: "auto", border: "1px solid #ccc" }}
          />
        </div>
      )}

      <button className="send-img-button" onClick={handleSubmit} style={{ marginTop: "16px", padding: "8px 16px" }}>
        Submit Image 🖼️
      </button>

      {result && (
        <div style={{ marginTop: "16px" }}>
          <p><strong>Class:</strong> {result.class}</p>
          <p><strong>Accuracy:</strong> {result.accuracy}%</p>
        </div>
      )}
    </section>
  );
}
