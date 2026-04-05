import { useState, useEffect, useRef } from "react";
import { getAnnotations, createAnnotation, deleteAnnotation, getSummary, getCharacters, createCharacter, deleteCharacter } from "./api";

export default function BookDetail({ book, onBack }) {
  const [annotations, setAnnotations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [body, setBody] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
  const [summary, setSummary] = useState(null);
  const [summaryLoading, setSummaryLoading] = useState(false);
  const [summaryError, setSummaryError] = useState("");
  const [expandedChar, setExpandedChar] = useState(null);
  const [characters, setCharacters] = useState([]);
  const [newCharName, setNewCharName] = useState("");
  const [addingChar, setAddingChar] = useState(false);
  const textareaRef = useRef(null);

  useEffect(() => {
    getAnnotations(book.id)
      .then(setAnnotations)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
    getCharacters(book.id)
      .then(setCharacters)
      .catch(() => {});
  }, [book.id]);

  async function handleCreate(e) {
    e.preventDefault();
    if (!body.trim()) return;
    setSubmitting(true);
    try {
      const a = await createAnnotation(book.id, body);
      setAnnotations([a, ...annotations]);
      setBody("");
      textareaRef.current?.focus();
    } catch (e) {
      setError(e.message);
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(id) {
    if (!confirm("Delete this annotation?")) return;
    try {
      await deleteAnnotation(book.id, id);
      setAnnotations(annotations.filter((a) => a.id !== id));
    } catch (e) {
      setError(e.message);
    }
  }

  function insertCharacter(name) {
    const ta = textareaRef.current;
    if (!ta) return;
    const start = ta.selectionStart;
    const end = ta.selectionEnd;
    const before = body.slice(0, start);
    const after = body.slice(end);
    const inserted = (before.length > 0 && !before.endsWith(" ") ? " " : "") + name + " ";
    const newBody = before + inserted + after;
    setBody(newBody);
    setTimeout(() => {
      ta.focus();
      ta.setSelectionRange(before.length + inserted.length, before.length + inserted.length);
    }, 0);
  }

  async function handleAddCharacter(e) {
    e.preventDefault();
    const name = newCharName.trim();
    if (!name) return;
    setAddingChar(true);
    try {
      const c = await createCharacter(book.id, name);
      setCharacters([...characters, c]);
      setNewCharName("");
    } catch (e) {
      setError(e.message);
    } finally {
      setAddingChar(false);
    }
  }

  async function handleDeleteCharacter(id) {
    try {
      await deleteCharacter(book.id, id);
      setCharacters(characters.filter((c) => c.id !== id));
    } catch (e) {
      setError(e.message);
    }
  }

  async function handleSummary() {
    setSummaryLoading(true);
    setSummaryError("");
    setSummary(null);
    setExpandedChar(null);
    try {
      const data = await getSummary(book.id);
      setSummary(data);
    } catch (e) {
      setSummaryError(e.message);
    } finally {
      setSummaryLoading(false);
    }
  }

  function formatDate(iso) {
    const d = new Date(iso);
    return d.toLocaleDateString("en-US", { month: "short", day: "numeric", year: "numeric" }) +
      " · " + d.toLocaleTimeString("en-US", { hour: "2-digit", minute: "2-digit" });
  }

  return (
    <div className="books-container">
      <header className="header">
        <button className="btn-back" onClick={onBack}>← Books</button>
        <div className="header-right" />
      </header>

      <main className="main">
        <div className="book-hero">
          <h2>{book.title}</h2>
          <p className="book-author">{book.author}{book.year ? ` · ${book.year}` : ""}</p>
          {book.description && <p className="book-desc-full">{book.description}</p>}
        </div>

        <div className="annotations-header">
          <h3>Annotations <span className="count">{annotations.length}</span></h3>
        </div>

        {annotations.length > 0 && (
          <div className="summary-section">
            <div className="summary-bar">
              <span className="summary-label">Character summaries</span>
              <button className="btn-outline" onClick={handleSummary} disabled={summaryLoading}>
                {summaryLoading ? "Thinking…" : "✦ Generate"}
              </button>
            </div>
            {summaryError && <p className="error">{summaryError}</p>}
            {summary && summary.characters.length === 0 && (
              <p className="empty">No characters found in annotations.</p>
            )}
            {summary && summary.characters.length > 0 && (
              <div className="character-list">
                {summary.characters.map((c) => (
                  <div key={c.name} className="character-card" onClick={() => setExpandedChar(expandedChar === c.name ? null : c.name)}>
                    <div className="character-header">
                      <span className="character-name">{c.name}</span>
                      <span className="character-chevron">{expandedChar === c.name ? "▲" : "▼"}</span>
                    </div>
                    {expandedChar === c.name && (
                      <p className="character-summary">{c.summary}</p>
                    )}
                  </div>
                ))}
              </div>
            )}
          </div>
        )}

        <form className="annotation-form" onSubmit={handleCreate}>
          {(characters.length > 0 || true) && (
            <div className="char-chips-section">
              <div className="char-chips">
                {characters.map((c) => (
                  <span key={c.id} className="char-chip-wrap">
                    <button type="button" className="char-chip" onClick={() => insertCharacter(c.name)}>
                      {c.name}
                    </button>
                    <button type="button" className="char-chip-del" onClick={() => handleDeleteCharacter(c.id)}>✕</button>
                  </span>
                ))}
                <form className="char-add-form" onSubmit={handleAddCharacter}>
                  <input
                    type="text"
                    className="char-add-input"
                    placeholder="+ Add character"
                    value={newCharName}
                    onChange={(e) => setNewCharName(e.target.value)}
                  />
                  {newCharName.trim() && (
                    <button type="submit" className="char-add-btn" disabled={addingChar}>Add</button>
                  )}
                </form>
              </div>
            </div>
          )}
          <textarea
            ref={textareaRef}
            placeholder="Write an annotation…"
            value={body}
            onChange={(e) => setBody(e.target.value)}
            rows={3}
          />
          <button type="submit" className="btn-primary" disabled={submitting || !body.trim()}>
            {submitting ? "Saving..." : "Save annotation"}
          </button>
        </form>

        {error && <p className="error">{error}</p>}

        {loading ? (
          <p className="empty">Loading...</p>
        ) : annotations.length === 0 ? (
          <div className="empty-state">
            <p>No annotations yet.</p>
            <p>Write your first one above.</p>
          </div>
        ) : (
          <div className="annotation-list">
            {annotations.map((a) => (
              <div key={a.id} className="annotation-card">
                <p className="annotation-body">{a.body}</p>
                <div className="annotation-footer">
                  <span className="annotation-date">{formatDate(a.created_at)}</span>
                  <button className="btn-delete" onClick={() => handleDelete(a.id)} title="Delete">✕</button>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
