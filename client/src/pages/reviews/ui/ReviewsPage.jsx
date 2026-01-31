import React from "react";
import styles from "./ReviewsPage.module.css";
import { useReviewsList } from "../../../features/reviews";

const ReviewsPage = ({
    showAlert,
    isAdmin = false,
    isAuth = false,
    profileData = {},
}) => {
    const {
        STATUS_MODERATING,
        currentUserId,
        editingReviewId,
        editingText,
        form,
        formatDate,
        getAuthorId,
        getReviewBody,
        getReviewId,
        handleCancelEdit,
        handleChange,
        handleDelete,
        handleModerationDecision,
        handleSaveEdit,
        handleStartEdit,
        handleSubmit,
        hasSubmitted,
        isLoadingReviews,
        isReviewLoading,
        setEditingText,
        showPosts,
        submitting,
        visibleReviews,
    } = useReviewsList({ showAlert, isAdmin, isAuth, profileData });

    return (
        <div className={styles.page}>
            <main className={styles.main}>
                {/* Лента отзывов и модерация */}
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
                                            <span className={styles.reviewDate}>
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

                {/* Форма отправки отзыва для авторизованных пользователей */}
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

                {/* Сообщение после отправки */}
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
