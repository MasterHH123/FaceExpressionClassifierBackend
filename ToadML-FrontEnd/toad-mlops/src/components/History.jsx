import "../styles/History.css"

export default function History({ history = [] }) {
  if (!history || history.length === 0) {
    return <section style={{ marginTop: "32px" }}>
      <h2>Historial de Inferencias</h2>
      <p>No hay inferencias registradas.</p>
      </section>;
  }

  const handleDownload = (imgUrl, filename) => {
    const link = document.createElement("a");
    link.href = imgUrl;
    link.download = filename;
    link.click();
  };

  return (
    <section style={{ marginTop: "32px" }}>
      <h2>Historial de Inferencias</h2>
      <table>
        <thead>
          <tr>
            <th>Imagen</th>
            <th>Archivo</th>
            <th>Resultado</th>
            <th>Precisión</th>
            <th>Descargar</th>
          </tr>
        </thead>
        <tbody>
          {history.map((entry, index) => (
            <tr key={index}>
              <td>
                <img
                  src={entry.img}
                  alt={`preview-${index}`}
                />
              </td>
              <td>{entry.imgfile}</td>
              <td>{entry.result}</td>
              <td>{entry.accuracy}%</td>
              <td>
                <button className="download-button"
                  onClick={() => handleDownload(entry.img, entry.imgfile)}
                  
                >
                  Descargar
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}
