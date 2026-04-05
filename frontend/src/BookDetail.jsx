import { useState, useEffect, useRef } from "react";
import { getAnnotations, createAnnotation, deleteAnnotation, getSummary } from "./api";

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
  const textareaRef = useRef(null);

  useEffect(() => {
    getAnnotations(book.id)
      .then(setAnnotations)
      .catch((e) => setError(e.message))
      .finally(() => setLoading(false));
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
