import { useState, useEffect, useRef } from "react";
import { getAnnotations, createAnnotation, deleteAnnotation } from "./api";

export default function BookDetail({ book, onBack }) {
  const [annotations, setAnnotations] = useState([]);
  const [loading, setLoading] = useState(true);
  const [body, setBody] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState("");
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
    try {
      await deleteAnnotation(book.id, id);
      setAnnotations(annotations.filter((a) => a.id !== id));
    } catch (e) {
      setError(e.message);
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

        <form className="annotation-form" onSubmit={handleCreate}>
          <textarea
            ref={textareaRef}
            placeholder="Write an annotation… or use your keyboard mic 🎤"
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
