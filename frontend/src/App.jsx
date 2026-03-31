import { useState } from "react";
import Auth from "./Auth";
import Books from "./Books";
import "./index.css";

export default function App() {
  const [user, setUser] = useState(() => {
    const stored = localStorage.getItem("user");
    return stored ? JSON.parse(stored) : null;
  });

  function handleAuth(user) {
    setUser(user);
  }

  function handleLogout() {
    localStorage.removeItem("token");
    localStorage.removeItem("user");
    setUser(null);
  }

  if (!user) return <Auth onAuth={handleAuth} />;
  return <Books user={user} onLogout={handleLogout} />;
}
