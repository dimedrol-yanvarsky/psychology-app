import React, { useEffect, useState } from "react";
import axios from "axios";
import styles from "./ReviewsPage.module.css"; // можно назвать и просто module.css

const api = axios.create({
    baseURL: "http://localhost:8080/api",
});

const initialReviews = [
    {
        id: "seed-1",
        name: "Анна Иванова",
        text: "Очень понравился сервис: быстрые ответы поддержки и удобный интерфейс. Буду рекомендовать коллегам!",
        createdAt: "2024-05-10T12:00:00.000Z",
    },
    {
        id: "seed-2",
        name: "Дмитрий Петров",
        text: "Хороший опыт: всё работает стабильно, а новые функции появляются регулярно. Спасибо команде разработчиков!",
        createdAt: "2024-05-05T09:30:00.000Z",
    },
];

const ReviewsPage = (props) => {
    const [reviews, setReviews] = useState(initialReviews);
    const [form, setForm] = useState({
        name: "",
        text: "",
    });
    const [submitting, setSubmitting] = useState(false);
    const [hasSubmitted, setHasSubmitted] = useState(false);
    const [isEditing, setIsEditing] = useState(false);
    const [showPosts, setShowPosts] = useState(initialReviews.length > 0);

    useEffect(() => {
        // проверяем, оставлял ли текущий пользователь отзыв (простая реализация через localStorage)
        const storedFlag =
            localStorage.getItem("hasSubmittedReview") === "true";
        setHasSubmitted(storedFlag);

        fetchReviews();
    }, []);

    const fetchReviews = async () => {
        try {
            const response = await api.get("/reviews");
            if (response.data && Array.isArray(response.data.reviews)) {
                setReviews(response.data.reviews);
                setShowPosts(response.data.reviews.length > 0);
            }
        } catch (error) {
            console.error("Не удалось загрузить отзывы", error);
        }
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

            if (response.data?.review) {
                // добавляем новый отзыв в начало списка
                setReviews((prev) => {
                    const next = [response.data.review, ...prev];
                    setShowPosts(next.length > 0);
                    return next;
                });
            }

            setHasSubmitted(true);
            localStorage.setItem("hasSubmittedReview", "true");

            setForm({
                name: "",
                text: "",
            });
        } catch (error) {
            console.error(
                error.response?.data?.message ||
                    "Произошла ошибка при отправке отзыва. Попробуйте ещё раз."
            );
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
            <header className={styles.topBar}>
                <div className={styles.topBarActions}>
                    <button
                        type="button"
                        className={`${styles.actionButton} ${styles.deleteButton}`}
                        onClick={() => {
                            setReviews([]);
                            setShowPosts(false);
                        }}
                    >
                        Шаблон кнопки №1
                    </button>
                    <button
                        type="button"
                        className={`${styles.actionButton} ${styles.secondaryEditButton}`}
                        onClick={() => setIsEditing((prev) => !prev)}
                    >
                        Шаблон кнопки №2
                    </button>
                </div>
            </header>

            <main className={styles.main}>
                <section className={styles.reviewsSection}>
                    {!showPosts ? (
                        <div className={styles.emptyStateBanner}>
                            Пока нет ни одного отзыва. Станьте первым!
                        </div>
                    ) : (
                        <div className={styles.reviewsList}>
                            {reviews.map((review) => (
                                <article
                                    key={review.id || review._id}
                                    className={`${styles.reviewCard} ${
                                        isEditing
                                            ? styles.reviewCardEditing
                                            : ""
                                    }`}
                                >
                                    <header className={styles.reviewHeader}>
                                        <span className={styles.reviewAuthor}>
                                            {review.name}
                                        </span>
                                        <span className={styles.reviewDate}>
                                            {formatDate(
                                                review.createdAt || review.date
                                            )}
                                        </span>
                                    </header>
                                    <p className={styles.reviewText}>
                                        {review.text}
                                    </p>
                                    <div className={styles.reviewActions}>
                                        {props.isAdmin === true ? (
                                            <>
                                                <button
                                                    type="button"
                                                    className={`${styles.reviewActionButton} ${styles.reviewDeleteButton}`}
                                                    onClick={() => {
                                                        const id =
                                                            review.id ||
                                                            review._id;
                                                        setReviews((prev) => {
                                                            const next =
                                                                prev.filter(
                                                                    (item) =>
                                                                        (item.id ||
                                                                            item._id) !==
                                                                        id
                                                                );
                                                            setShowPosts(
                                                                next.length > 0
                                                            );
                                                            return next;
                                                        });
                                                    }}
                                                >
                                                    Удалить отзыв
                                                </button>
                                                <button
                                                    type="button"
                                                    className={`${styles.reviewActionButton} ${styles.reviewEditButton}`}
                                                    onClick={() =>
                                                        setIsEditing(
                                                            (prev) => !prev
                                                        )
                                                    }
                                                >
                                                    {isEditing
                                                        ? "Закончить редактирование"
                                                        : "Редактировать"}
                                                </button>
                                            </>
                                        ) : (
                                            <></>
                                        )}
                                    </div>
                                </article>
                            ))}
                        </div>
                    )}
                </section>

                {/* Средний уровень: форма добавления отзыва (если ещё не добавлен) */}

                {props.isAuth === true ? (
                    <section className={styles.formSection}>
                        {!hasSubmitted ? (
                            <form
                                className={styles.form}
                                onSubmit={handleSubmit}
                            >
                                <h2 className={styles.formTitle}>
                                    Оставьте отзыв о сервисе
                                </h2>

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
                                    className={styles.reviewActionButton}
                                    disabled={submitting}
                                >
                                    {submitting
                                        ? "Отправка..."
                                        : "Отправить отзыв"}
                                </button>
                            </form>
                        ) : (
                            <div className={styles.alreadySubmitted}>
                                <h2>Спасибо за отзыв!</h2>
                                <p>Вы уже оставили отзыв о сервисе.</p>
                            </div>
                        )}
                    </section>
                ) : (
                    <></>
                )}
            </main>
        </div>
    );
};

export default ReviewsPage;
