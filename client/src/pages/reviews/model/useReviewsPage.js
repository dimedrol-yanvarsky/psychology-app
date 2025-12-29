import { useCallback, useEffect, useMemo, useState } from "react";
import {
    createReview,
    decideReview,
    deleteReview,
    getReviews,
    updateReview,
} from "../api/reviewsApi";

const STATUS_MODERATING = "Модерируется";
const STATUS_APPROVED = "Добавлен";
const STATUS_DENIED = "Отклонен";

const getReviewId = (review) => review?._id || review?.id || "";
const getAuthorId = (review) => review?.userID || review?.userId || "";
const getReviewBody = (review) => review?.reviewBody || review?.text || "";

export const useReviewsPage = ({
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

    const filterVisibleReviews = useCallback(
        (list) =>
            Array.isArray(list)
                ? isAdmin
                    ? list
                    : list.filter(
                          (review) =>
                              (review?.status || "").trim() === STATUS_APPROVED
                      )
                : [],
        [isAdmin]
    );

    const fetchReviews = useCallback(async () => {
        setIsLoadingReviews(true);
        try {
            const response = await getReviews();
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
    }, [reviews, filterVisibleReviews, isAuth, currentUserId]);

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
            const response = await createReview({
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
            await deleteReview({
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
            const response = await decideReview({
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
            const response = await updateReview({
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

    const visibleReviews = useMemo(
        () => filterVisibleReviews(reviews),
        [filterVisibleReviews, reviews]
    );

    return {
        STATUS_DENIED,
        STATUS_MODERATING,
        STATUS_APPROVED,
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
    };
};
