import { useCallback, useEffect, useMemo } from "react";
import { useSelector, useDispatch } from "react-redux";
import {
    createReview,
    decideReview,
    deleteReview,
    getReviews,
    updateReview,
} from "../../../entities/review";
import { useAuthContext } from "../../../shared/context/AuthContext";
import { useAlertContext } from "../../../shared/context/AlertContext";
import {
    fetchReviewsStart,
    fetchReviewsSuccess,
    fetchReviewsError,
    setFormText,
    submitStart,
    submitEnd,
    addReview,
    setHasSubmitted,
    setShowPosts,
    startEdit,
    cancelEdit,
    setEditingText,
    setReviewLoading,
    clearReviewLoading,
    removeReview,
    setReviews,
} from "./reviewsSlice";

const STATUS_MODERATING = "Модерируется";
const STATUS_APPROVED = "Добавлен";
const STATUS_DENIED = "Отклонен";

const getReviewId = (review) => review?._id || review?.id || "";
const getAuthorId = (review) => review?.userID || review?.userId || "";
const getReviewBody = (review) => review?.reviewBody || review?.text || "";

export const useReviewsList = () => {
    const { isAdmin, isAuth, profileData } = useAuthContext();
    const { showAlert } = useAlertContext();
    const reduxDispatch = useDispatch();

    const state = useSelector((s) => s.reviews);

    const currentUserId = (profileData?.id || "").trim();

    const isReviewLoading = (reviewId) => Boolean(state.loadingIds[reviewId]);

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
        reduxDispatch(fetchReviewsStart());
        try {
            const response = await getReviews();
            const loadedReviews = Array.isArray(response?.data?.reviews)
                ? response.data.reviews
                : [];
            reduxDispatch(fetchReviewsSuccess(loadedReviews));
        } catch (error) {
            console.error("Не удалось загрузить отзывы", error);
            reduxDispatch(fetchReviewsError());
            showAlert?.(
                "error",
                error?.response?.data?.message ||
                    "Не удалось загрузить отзывы"
            );
        }
    }, [showAlert, reduxDispatch]);

    useEffect(() => {
        fetchReviews();
    }, [fetchReviews]);

    useEffect(() => {
        const visible = filterVisibleReviews(state.reviews);
        reduxDispatch(setShowPosts(visible.length > 0));

        if (isAuth && currentUserId) {
            const userHasReview = state.reviews.some(
                (review) => getAuthorId(review) === currentUserId
            );
            reduxDispatch(setHasSubmitted(userHasReview));
        } else {
            reduxDispatch(setHasSubmitted(false));
        }
    }, [state.reviews, filterVisibleReviews, isAuth, currentUserId, reduxDispatch]);

    const handleChange = (event) => {
        reduxDispatch(setFormText(event.target.value));
    };

    const handleSubmit = async (event) => {
        event.preventDefault();
        if (state.submitting) return;

        if (!isAuth) {
            showAlert?.("error", "Авторизуйтесь, чтобы оставить отзыв");
            return;
        }

        if (!currentUserId) {
            showAlert?.("error", "Не удалось определить пользователя");
            return;
        }

        const body = state.form.text.trim();
        if (!body) {
            showAlert?.("error", "Заполните текст отзыва");
            return;
        }

        reduxDispatch(submitStart());
        try {
            const response = await createReview({
                userId: currentUserId,
                reviewBody: body,
            });

            if (response?.data?.review) {
                reduxDispatch(addReview(response.data.review));
                showAlert?.("success", "Отзыв отправлен на модерацию");
            }
        } catch (error) {
            showAlert?.(
                "error",
                error?.response?.data?.message ||
                    "Произошла ошибка при отправке отзыва. Попробуйте ещё раз."
            );
        } finally {
            reduxDispatch(submitEnd());
        }
    };

    const handleDelete = async (review) => {
        const reviewId = getReviewId(review);
        if (!reviewId) return;

        reduxDispatch(setReviewLoading(reviewId));
        try {
            await deleteReview({
                reviewId,
                userId: currentUserId,
                isAdmin,
            });

            reduxDispatch(removeReview(reviewId));
            showAlert?.("success", "Отзыв удален");
        } catch (error) {
            showAlert?.(
                "error",
                error?.response?.data?.message || "Не удалось удалить отзыв"
            );
        } finally {
            reduxDispatch(clearReviewLoading(reviewId));
        }
    };

    const handleModerationDecision = async (reviewId, decision) => {
        if (!reviewId) return;
        if (!currentUserId) {
            showAlert?.("error", "Не удалось определить администратора");
            return;
        }

        reduxDispatch(setReviewLoading(reviewId));

        try {
            const response = await decideReview({
                reviewId,
                adminId: currentUserId,
                decision,
            });

            if (response?.data?.review) {
                reduxDispatch(
                    setReviews(
                        state.reviews.map((item) =>
                            getReviewId(item) === reviewId
                                ? { ...item, ...response.data.review }
                                : item
                        )
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
            reduxDispatch(clearReviewLoading(reviewId));
        }
    };

    const handleStartEdit = (review) => {
        const reviewId = getReviewId(review);
        if (!reviewId) return;
        reduxDispatch(
            startEdit({ id: reviewId, text: getReviewBody(review) })
        );
    };

    const handleCancelEdit = () => {
        reduxDispatch(cancelEdit());
    };

    const handleSaveEdit = async (reviewId) => {
        if (!reviewId) return;
        if (!currentUserId) {
            showAlert?.("error", "Не удалось определить пользователя");
            return;
        }

        const body = state.editingText.trim();
        if (!body) {
            showAlert?.("error", "Введите текст отзыва");
            return;
        }

        reduxDispatch(setReviewLoading(reviewId));

        try {
            const response = await updateReview({
                reviewId,
                userId: currentUserId,
                reviewBody: body,
            });

            if (response?.data?.review) {
                reduxDispatch(
                    setReviews(
                        state.reviews.map((item) =>
                            getReviewId(item) === reviewId
                                ? { ...item, ...response.data.review }
                                : item
                        )
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
            reduxDispatch(clearReviewLoading(reviewId));
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
        () => filterVisibleReviews(state.reviews),
        [filterVisibleReviews, state.reviews]
    );

    return {
        STATUS_DENIED,
        STATUS_MODERATING,
        STATUS_APPROVED,
        currentUserId,
        editingReviewId: state.editingReviewId,
        editingText: state.editingText,
        form: state.form,
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
        hasSubmitted: state.hasSubmitted,
        isLoadingReviews: state.isLoadingReviews,
        isReviewLoading,
        setEditingText: (text) => reduxDispatch(setEditingText(text)),
        showPosts: state.showPosts,
        submitting: state.submitting,
        visibleReviews,
    };
};
