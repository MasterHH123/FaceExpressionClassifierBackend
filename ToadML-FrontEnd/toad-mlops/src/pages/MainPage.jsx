import Navbar from "../components/Navbar";
import Inference from "../components/Inference";
import History from "../components/History";
import About from "../components/About";
import Footer from "../components/Footer";

import { useState } from "react";

export default function MainPage({ user, setUser }) {

  const [history, setHistory] = useState([]);

  const addInferenceToHistory = (newInference) => {
    setHistory((prevHistory) => [...prevHistory, newInference]);
  };

  return (
    <div>
      <Navbar user={user} setUser={setUser} />
      <div id="inference" className="min-h-screen p-6 scroll-mt-24">
        <Inference addInferenceToHistory={addInferenceToHistory}/>
      </div>
      <div id="history" className="min-h-screen p-6 scroll-mt-24 bg-gray-50">
        <History history={history}/>
      </div>
      <div id="about" className="min-h-screen p-6 scroll-mt-24 bg-white">
        <About />
      </div>
      <Footer />
    </div>
  );
}
