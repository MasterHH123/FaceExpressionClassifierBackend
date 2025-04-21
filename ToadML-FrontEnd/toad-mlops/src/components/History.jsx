import "../styles/History.css"

export default function History({ history = [] }) {
  if (!history || history.length === 0) {
    return <section className="history-section" style={{ marginTop: "32px" }}>
      <h2>History</h2>
      <p style={{color: "white"}}>No recorded inferences.</p>
      </section>;
  }

  const handleDownload = (imgUrl, filename) => {
    const link = document.createElement("a");
    link.href = imgUrl;
    link.download = filename;
    link.click();
  };

  return (
    <section className="history-section" style={{ marginTop: "16px" }}>
      <h2>History</h2>
      <table>
        <thead>
          <tr>
            <th>Image</th>
            <th>File</th>
            <th>Result</th>
            <th>Accutacy</th>
            <th>Download</th>
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
                  Download 🔽
                </button>
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </section>
  );
}
