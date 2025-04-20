import "../styles/Navbar.css"

export default function Navbar({ user, setUser }) {
    const handleScroll = (id) => {
      document.getElementById(id)?.scrollIntoView({ behavior: "smooth" });
    };
  
    const logout = () => {
      setUser(null);
      localStorage.removeItem("user");
    };
  
    return (
      <nav className="" id="navbar">
        <div className="navbar-left">
          <button className="logo-button" onClick={() => handleScroll("inference")}><img src="/mushroom.webp" alt="Logo" className="navbar-logo" /></button>
          <button className="nav-button" onClick={() => handleScroll("inference")}>Inferencia</button>
          <button className="nav-button" onClick={() => handleScroll("history")}>Historial</button>
          <button className="nav-button" onClick={() => handleScroll("about")}>About</button>
        </div>
        <div>
          <span className="navbar-right">{user?.email}</span>
          <button className="nav-button" onClick={logout}>Cerrar sesión</button>
        </div>
      </nav>
    );
  }
  