import React, { useCallback, useEffect, useState } from "react";
import axios from "axios";
import styles from "./ReviewsPage.module.css";

const api = axios.create({
    baseURL: "http://localhost:8080/api",
});

const STATUS_MODERATING = "Модерируется";
const STATUS_APPROVED = "Добавлен";
const STATUS_DENIED = "Отклонен";

const getReviewId = (review) => review?._id || review?.id || "";
const getAuthorId = (review) => review?.userID || review?.userId || "";
const getReviewBody = (review) => review?.reviewBody || review?.text || "";

const ReviewsPage = ({
    showAlert,
    isAdmin = false,
    isAuth = false,
    profileData = {},
}) => {
    const [reviews, setReviews] = useState([]);
    const [form, setForm] = useState({ text: "" });
    const [submitting, setSubmitting] = useState(false);
    const [hasSubmitted, setHasSubmitted] = useState(false);
    const [showPosts, setShowPosts] = useState(false);
    const [editingReviewId, setEditingReviewId] = useState("");
    const [editingText, setEditingText] = useState("");
    const [loadingIds, setLoadingIds] = useState({});
    const [isLoadingReviews, setIsLoadingReviews] = useState(false);

    const currentUserId = (profileData?.id || "").trim();

    const setReviewLoading = (reviewId, value) => {
        setLoadingIds((prev) => {
            const next = { ...prev };
            if (value) {
                next[reviewId] = true;
            } else {
                delete next[reviewId];
            }
            return next;
        });
    };

    const isReviewLoading = (reviewId) => Boolean(loadingIds[reviewId]);

    const filterVisibleReviews = (list) =>
        Array.isArray(list)
            ? isAdmin
                ? list
                : list.filter(
                        (review) =>
                            (review?.status || "").trim() === STATUS_APPROVED
                  )
            : [];

    const fetchReviews = useCallback(async () => {
        setIsLoadingReviews(true);
        try {
            const response = await api.get("/reviews/getReviews");
            const loadedReviews = Array.isArray(response?.data?.reviews)
                ? response.data.reviews
                : [];
            setReviews(loadedReviews);
        } catch (error) {
            console.error("Не удалось загрузить отзывы", error);
            showAlert?.(
                "error",
                error?.response?.data?.message ||
                    "Не удалось загрузить отзывы"
            );
        } finally {
            setIsLoadingReviews(false);
        }
    }, [showAlert]);

    useEffect(() => {
        fetchReviews();
    }, [fetchReviews]);

    useEffect(() => {
        const visible = filterVisibleReviews(reviews);
        setShowPosts(visible.length > 0);

        if (isAuth && currentUserId) {
            const userHasReview = reviews.some(
                (review) => getAuthorId(review) === currentUserId
            );
            setHasSubmitted(userHasReview);
        } else {
            setHasSubmitted(false);
        }
    }, [reviews, isAdmin, isAuth, currentUserId]);

    const handleChange = (event) => {
        setForm({ text: event.target.value });
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        if (submitting) return;

        if (!isAuth) {
            showAlert?.("error", "Авторизуйтесь, чтобы оставить отзыв");
            return;
        }

        if (!currentUserId) {
            showAlert?.("error", "Не удалось определить пользователя");
            return;
        }

        const body = form.text.trim();
        if (!body) {
            showAlert?.("error", "Заполните текст отзыва");
            return;
        }

        setSubmitting(true);
        try {
            const response = await api.post("/reviews/createReview", {
                userId: currentUserId,
                reviewBody: body,
            });

            if (response?.data?.review) {
                setReviews((prev) => [response.data.review, ...prev]);
                setForm({ text: "" });
                showAlert?.("success", "Отзыв отправлен на модерацию");
            }
        } catch (error) {
            showAlert?.(
                "error",
                error?.response?.data?.message ||
                    "Произошла ошибка при отправке отзыва. Попробуйте ещё раз."
            );
        } finally {
            setSubmitting(false);
        }
    };

    const handleDelete = async (review) => {
        const reviewId = getReviewId(review);
        if (!reviewId) return;

        setReviewLoading(reviewId, true);
        try {
            await api.post("/reviews/deleteReview", {
                reviewId,
                userId: currentUserId,
                isAdmin,
            });

            setReviews((prev) =>
                prev.filter((item) => getReviewId(item) !== reviewId)
            );
            showAlert?.("success", "Отзыв удален");
        } catch (error) {
            showAlert?.(
                "error",
                error?.response?.data?.message || "Не удалось удалить отзыв"
            );
        } finally {
            setReviewLoading(reviewId, false);
        }
    };

    const handleModerationDecision = async (reviewId, decision) => {
        if (!reviewId) return;
        if (!currentUserId) {
            showAlert?.("error", "Не удалось определить администратора");
            return;
        }

        setReviewLoading(reviewId, true);

        try {
            const response = await api.post("/reviews/approveOrDeny", {
                reviewId,
                adminId: currentUserId,
                decision,
            });

            if (response?.data?.review) {
                setReviews((prev) =>
                    prev.map((item) =>
                        getReviewId(item) === reviewId
                            ? { ...item, ...response.data.review }
                            : item
                    )
                );
                showAlert?.(
                    "success",
                    decision === "approve"
                        ? "Отзыв принят"
                        : "Отзыв отклонен"
                );
            }
        } catch (error) {
            showAlert?.(
                "error",
                error?.response?.data?.message ||
                    "Не удалось обновить статус отзыва"
            );
        } finally {
            setReviewLoading(reviewId, false);
        }
    };

    const handleStartEdit = (review) => {
        const reviewId = getReviewId(review);
        if (!reviewId) return;
        setEditingReviewId(reviewId);
        setEditingText(getReviewBody(review));
    };

    const handleCancelEdit = () => {
        setEditingReviewId("");
        setEditingText("");
    };

    const handleSaveEdit = async (reviewId) => {
        if (!reviewId) return;
        if (!currentUserId) {
            showAlert?.("error", "Не удалось определить пользователя");
            return;
        }

        const body = editingText.trim();
        if (!body) {
            showAlert?.("error", "Введите текст отзыва");
            return;
        }

        setReviewLoading(reviewId, true);

        try {
            const response = await api.post("/reviews/updateReview", {
                reviewId,
                userId: currentUserId,
                reviewBody: body,
            });

            if (response?.data?.review) {
                setReviews((prev) =>
                    prev.map((item) =>
                        getReviewId(item) === reviewId
                            ? { ...item, ...response.data.review }
                            : item
                    )
                );
                showAlert?.("success", "Отзыв обновлен");
                handleCancelEdit();
            }
        } catch (error) {
            showAlert?.(
                "error",
                error?.response?.data?.message || "Не удалось обновить отзыв"
            );
        } finally {
            setReviewLoading(reviewId, false);
        }
    };

    const formatDate = (value) => {
        if (!value) return "";
        const trimmed = String(value).trim();

        if (/^\d{2}\.\d{2}\.\d{4}$/.test(trimmed)) {
            return trimmed;
        }

        const parsed = new Date(trimmed);
        if (Number.isNaN(parsed.getTime())) {
            return trimmed;
        }

        return parsed.toLocaleDateString("ru-RU", {
            day: "2-digit",
            month: "2-digit",
            year: "numeric",
        });
    };

    const visibleReviews = filterVisibleReviews(reviews);

    return (
        <div className={styles.page}>
            <main className={styles.main}>
                <section className={styles.reviewsSection}>
                    {isLoadingReviews ? (
                        <div className={styles.emptyStateBanner}>
                            Загружаем отзывы...
                        </div>
                    ) : !showPosts ? (
                        <div className={styles.emptyStateBanner}>
                            Пока нет ни одного отзыва. Станьте первым!
                        </div>
                    ) : (
                        <div className={styles.reviewsList}>
                            {visibleReviews.map((review) => {
                                const reviewId = getReviewId(review);
                                const authorId = getAuthorId(review);
                                const isOwnReview =
                                    currentUserId && authorId === currentUserId;
                                const inModeration =
                                    (review?.status || "").trim() ===
                                    STATUS_MODERATING;
                                const isEditingCurrent =
                                    editingReviewId === reviewId;
                                const isActionLoading =
                                    isReviewLoading(reviewId);

                                return (
                                    <article
                                        key={reviewId}
                                        className={`${styles.reviewCard} ${
                                            isAdmin && inModeration
                                                ? styles.reviewCardModerating
                                                : ""
                                        }`}
                                    >
                                        <header className={styles.reviewHeader}>
                                            <span
                                                className={styles.reviewAuthor}
                                            >
                                                {review?.firstName || "Аноним"}
                                            </span>
                                            <span
                                                className={styles.reviewDate}
                                            >
                                                {formatDate(
                                                    review?.date ||
                                                        review?.createdAt
                                                )}
                                            </span>
                                        </header>

                                        {isEditingCurrent ? (
                                            <div className={styles.editRow}>
                                                <input
                                                    type="text"
                                                    value={editingText}
                                                    onChange={(event) =>
                                                        setEditingText(
                                                            event.target.value
                                                        )
                                                    }
                                                    className={styles.editInput}
                                                />
                                                <button
                                                    type="button"
                                                    className={`${styles.reviewActionButton} ${styles.reviewEditButton}`}
                                                    onClick={() =>
                                                        handleSaveEdit(
                                                            reviewId
                                                        )
                                                    }
                                                    disabled={isActionLoading}
                                                >
                                                    Сохранить
                                                </button>
                                                <button
                                                    type="button"
                                                    className={`${styles.reviewActionButton} ${styles.reviewDeleteButton}`}
                                                    onClick={handleCancelEdit}
                                                    disabled={isActionLoading}
                                                >
                                                    Отмена
                                                </button>
                                            </div>
                                        ) : (
                                            <p className={styles.reviewText}>
                                                {getReviewBody(review)}
                                            </p>
                                        )}

                                        <div className={styles.reviewActions}>
                                            {isAdmin && inModeration ? (
                                                <>
                                                    <button
                                                        type="button"
                                                        className={`${styles.reviewActionButton} ${styles.reviewEditButton}`}
                                                        onClick={() =>
                                                            handleModerationDecision(
                                                                reviewId,
                                                                "approve"
                                                            )
                                                        }
                                                        disabled={
                                                            isActionLoading
                                                        }
                                                    >
                                                        Принять отзыв
                                                    </button>
                                                    <button
                                                        type="button"
                                                        className={`${styles.reviewActionButton} ${styles.reviewDeleteButton}`}
                                                        onClick={() =>
                                                            handleModerationDecision(
                                                                reviewId,
                                                                "deny"
                                                            )
                                                        }
                                                        disabled={
                                                            isActionLoading
                                                        }
                                                    >
                                                        Отклонить
                                                    </button>
                                                </>
                                            ) : (
                                                <>
                                                    {(isAdmin || isOwnReview) && (
                                                        <button
                                                            type="button"
                                                            className={`${styles.reviewActionButton} ${styles.reviewDeleteButton}`}
                                                            onClick={() =>
                                                                handleDelete(
                                                                    review
                                                                )
                                                            }
                                                            disabled={
                                                                isActionLoading
                                                            }
                                                        >
                                                            Удалить отзыв
                                                        </button>
                                                    )}
                                                    {isOwnReview && (
                                                        <button
                                                            type="button"
                                                            className={`${styles.reviewActionButton} ${styles.reviewEditButton}`}
                                                            onClick={() =>
                                                                handleStartEdit(
                                                                    review
                                                                )
                                                            }
                                                            disabled={
                                                                isActionLoading
                                                            }
                                                        >
                                                            Редактировать
                                                        </button>
                                                    )}
                                                </>
                                            )}
                                        </div>
                                    </article>
                                );
                            })}
                        </div>
                    )}
                </section>

                {isAuth && !hasSubmitted && (
                    <section className={styles.formSection}>
                        <form className={styles.form} onSubmit={handleSubmit}>
                            <h2 className={styles.formTitle}>
                                Оставьте отзыв о сервисе
                            </h2>

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
                                    className={styles.wideTextarea}
                                />
                            </div>

                            <button
                                type="submit"
                                className={styles.submitButton}
                                disabled={submitting}
                            >
                                {submitting
                                    ? "Отправка..."
                                    : "Отправить отзыв"}
                            </button>
                        </form>
                    </section>
                )}

                {isAuth && hasSubmitted && (
                    <div className={styles.alreadySubmitted}>
                        <h2>Спасибо за отзыв!</h2>
                        <p>Вы уже оставили отзыв о сервисе.</p>
                    </div>
                )}
            </main>
        </div>
    );
};

export default ReviewsPage;
