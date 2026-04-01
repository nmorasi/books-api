import { useState, useEffect } from "react";
import { getBooks, createBook, deleteBook } from "./api";
import BookDetail from "./BookDetail";

export default function Books({ user, onLogout }) {
  const [books, setBooks] = useState([]);
  const [loading, setLoading] = useState(true);
  const [showForm, setShowForm] = useState(false);
  const [error, setError] = useState("");
  const [form, setForm] = useState({ title: "", author: "", description: "", year: "" });
  const [submitting, setSubmitting] = useState(false);
  const [selectedBook, setSelectedBook] = useState(null);

  useEffect(() => {
    fetchBooks();
  }, []);

  async function fetchBooks() {
    try {
      const data = await getBooks();
      setBooks(data);
    } catch (err) {
      setError(err.message);
    } finally {
      setLoading(false);
    }
  }

  async function handleCreate(e) {
    e.preventDefault();
    setSubmitting(true);
    try {
      const book = await createBook({
        title: form.title,
        author: form.author,
        description: form.description,
        year: parseInt(form.year) || 0,
      });
      setBooks([book, ...books]);
      setForm({ title: "", author: "", description: "", year: "" });
      setShowForm(false);
    } catch (err) {
      setError(err.message);
    } finally {
      setSubmitting(false);
    }
  }

  async function handleDelete(e, id) {
    e.stopPropagation();
    try {
      await deleteBook(id);
      setBooks(books.filter((b) => b.id !== id));
    } catch (err) {
      setError(err.message);
    }
  }

  if (selectedBook) {
    return <BookDetail book={selectedBook} onBack={() => setSelectedBook(null)} />;
  }

  return (
    <div className="books-container">
      <header className="header">
        <h1 className="logo">📚 Books</h1>
        <div className="header-right">
          <span className="username">Hey, {user.name}</span>
          <button className="btn-outline" onClick={onLogout}>Sign out</button>
        </div>
      </header>

      <main className="main">
        <div className="books-header">
          <h2>My Books <span className="count">{books.length}</span></h2>
          <button className="btn-primary" onClick={() => setShowForm(!showForm)}>
            {showForm ? "Cancel" : "+ Add Book"}
          </button>
        </div>

        {error && <p className="error">{error}</p>}

        {showForm && (
          <form className="book-form" onSubmit={handleCreate}>
            <div className="form-row">
              <div className="field">
                <label>Title *</label>
                <input
                  type="text"
                  placeholder="Book title"
                  value={form.title}
                  onChange={(e) => setForm({ ...form, title: e.target.value })}
                  required
                />
              </div>
              <div className="field">
                <label>Author *</label>
                <input
                  type="text"
                  placeholder="Author name"
                  value={form.author}
                  onChange={(e) => setForm({ ...form, author: e.target.value })}
                  required
                />
              </div>
            </div>
            <div className="form-row">
              <div className="field">
                <label>Description</label>
                <input
                  type="text"
                  placeholder="Short description"
                  value={form.description}
                  onChange={(e) => setForm({ ...form, description: e.target.value })}
                />
              </div>
              <div className="field field-small">
                <label>Year</label>
                <input
                  type="number"
                  placeholder="2024"
                  value={form.year}
                  onChange={(e) => setForm({ ...form, year: e.target.value })}
                />
              </div>
            </div>
            <button type="submit" className="btn-primary" disabled={submitting}>
              {submitting ? "Saving..." : "Save Book"}
            </button>
          </form>
        )}

        {loading ? (
          <p className="empty">Loading...</p>
        ) : books.length === 0 ? (
          <div className="empty-state">
            <p>No books yet.</p>
            <p>Click <strong>+ Add Book</strong> to get started.</p>
          </div>
        ) : (
          <div className="book-list">
            {books.map((book) => (
              <div key={book.id} className="book-card clickable" onClick={() => setSelectedBook(book)}>
                <div className="book-info">
                  <h3>{book.title}</h3>
                  <p className="book-author">{book.author}{book.year ? ` · ${book.year}` : ""}</p>
                  {book.description && <p className="book-desc">{book.description}</p>}
                </div>
                <div className="book-actions">
                  <span className="annotation-hint">Tap to annotate →</span>
                  <button
                    className="btn-delete"
                    onClick={(e) => handleDelete(e, book.id)}
                    title="Delete"
                  >✕</button>
                </div>
              </div>
            ))}
          </div>
        )}
      </main>
    </div>
  );
}
