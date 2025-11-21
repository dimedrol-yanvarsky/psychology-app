import React, { useEffect, useState } from "react";
import axios from "axios";
import styles from "./ReviewsPage.module.css"; // можно назвать и просто module.css

const api = axios.create({
  baseURL: "http://localhost:8080/api",
});

const ReviewsPage = () => {
  const [reviews, setReviews] = useState([]);
  const [form, setForm] = useState({
    name: "",
    email: "",
    text: "",
  });
  const [submitting, setSubmitting] = useState(false);
  const [toast, setToast] = useState(null); // { type: 'success' | 'error', message: string }
  const [hasSubmitted, setHasSubmitted] = useState(false);
  const [isEditing, setIsEditing] = useState(false);

  useEffect(() => {
    // проверяем, оставлял ли текущий пользователь отзыв (простая реализация через localStorage)
    const storedFlag = localStorage.getItem("hasSubmittedReview") === "true";
    setHasSubmitted(storedFlag);

    fetchReviews();
  }, []);

  const fetchReviews = async () => {
    try {
      const response = await api.get("/reviews");
      if (response.data && Array.isArray(response.data.reviews)) {
        setReviews(response.data.reviews);
      }
    } catch (error) {
      setToast({
        type: "error",
        message: "Не удалось загрузить отзывы. Попробуйте обновить страницу.",
      });
      autoHideToast();
    }
  };

  const autoHideToast = () => {
    setTimeout(() => {
      setToast(null);
    }, 5000);
  };

  const handleChange = (event) => {
    const { name, value } = event.target;
    setForm((prev) => ({
      ...prev,
      [name]: value,
    }));
  };

  const handleSubmit = async (event) => {
    event.preventDefault();
    if (submitting) return;

    setSubmitting(true);
    try {
      const response = await api.post("/reviews", form);

      setToast({
        type: "success",
        message: response.data?.message || "Спасибо! Ваш отзыв отправлен.",
      });
      autoHideToast();

      if (response.data?.review) {
        // добавляем новый отзыв в начало списка
        setReviews((prev) => [response.data.review, ...prev]);
      }

      setHasSubmitted(true);
      localStorage.setItem("hasSubmittedReview", "true");

      setForm({
        name: "",
        email: "",
        text: "",
      });
    } catch (error) {
      const message =
        error.response?.data?.message ||
        "Произошла ошибка при отправке отзыва. Попробуйте ещё раз.";

      setToast({
        type: "error",
        message,
      });
      autoHideToast();
    } finally {
      setSubmitting(false);
    }
  };

  const formatDate = (value) => {
    if (!value) return "";
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return "";
    return date.toLocaleDateString("ru-RU", {
      day: "2-digit",
      month: "2-digit",
      year: "numeric",
    });
  };

  return (
    <div className={styles.page}>
      {/* Верхний уровень: слева уведомление, справа кнопка "редактировать отзывы" */}
      <header className={styles.topBar}>
        <div className={styles.toastContainer}>
          {toast && (
            <div
              className={`${styles.toast} ${
                toast.type === "success"
                  ? styles.toastSuccess
                  : styles.toastError
              }`}
            >
              {toast.message}
            </div>
          )}
        </div>

        <div className={styles.topBarActions}>
          <button
            type="button"
            className={styles.editButton}
            onClick={() => setIsEditing((prev) => !prev)}
          >
            {isEditing ? "Закончить редактирование" : "Редактировать отзывы"}
          </button>
        </div>
      </header>

      {/* Средний и нижний уровни */}
      <main className={styles.main}>
        {/* Средний уровень: форма добавления отзыва (если ещё не добавлен) */}
        <section className={styles.formSection}>
          {!hasSubmitted ? (
            <form className={styles.form} onSubmit={handleSubmit}>
              <h2 className={styles.formTitle}>Оставьте отзыв о сервисе</h2>

              <div className={styles.field}>
                <label htmlFor="name">Имя</label>
                <input
                  id="name"
                  name="name"
                  type="text"
                  value={form.name}
                  onChange={handleChange}
                  placeholder="Например, Иван"
                  required
                />
              </div>

              <div className={styles.field}>
                <label htmlFor="email">E-mail</label>
                <input
                  id="email"
                  name="email"
                  type="email"
                  value={form.email}
                  onChange={handleChange}
                  placeholder="example@mail.ru"
                  required
                />
              </div>

              <div className={styles.field}>
                <label htmlFor="text">Ваш отзыв</label>
                <textarea
                  id="text"
                  name="text"
                  rows="4"
                  value={form.text}
                  onChange={handleChange}
                  placeholder="Поделитесь своим опытом..."
                  required
                />
              </div>

              <button
                type="submit"
                className={styles.submitButton}
                disabled={submitting}
              >
                {submitting ? "Отправка..." : "Отправить отзыв"}
              </button>
            </form>
          ) : (
            <div className={styles.alreadySubmitted}>
              <h2>Спасибо за отзыв!</h2>
              <p>Вы уже оставили отзыв о сервисе.</p>
            </div>
          )}
        </section>

        {/* Нижний уровень: список отзывов */}
        <section className={styles.reviewsSection}>
          <h2 className={styles.sectionTitle}>Отзывы пользователей</h2>
          <div className={styles.reviewsList}>
            {reviews.length === 0 && (
              <p className={styles.emptyState}>
                Пока нет ни одного отзыва. Станьте первым!
              </p>
            )}

            {reviews.map((review) => (
              <article
                key={review.id || review._id}
                className={`${styles.reviewCard} ${
                  isEditing ? styles.reviewCardEditing : ""
                }`}
              >
                <header className={styles.reviewHeader}>
                  <span className={styles.reviewAuthor}>{review.name}</span>
                  <span className={styles.reviewDate}>
                    {formatDate(review.createdAt || review.date)}
                  </span>
                </header>
                <p className={styles.reviewText}>{review.text}</p>
              </article>
            ))}
          </div>
        </section>
      </main>
    </div>
  );
};

export default ReviewsPage;