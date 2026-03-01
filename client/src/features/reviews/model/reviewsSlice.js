import { createSlice } from "@reduxjs/toolkit";

const initialState = {
    reviews: [],
    isLoadingReviews: false,
    form: { text: "" },
    submitting: false,
    hasSubmitted: false,
    showPosts: false,
    editingReviewId: "",
    editingText: "",
    loadingIds: {},
};

// Слайс отзывов: список, создание, модерация, редактирование, удаление.
const reviewsSlice = createSlice({
    name: "reviews",
    initialState,
    reducers: {
        // Загрузка
        fetchReviewsStart(state) {
            state.isLoadingReviews = true;
        },
        fetchReviewsSuccess(state, action) {
            state.isLoadingReviews = false;
            state.reviews = action.payload;
        },
        fetchReviewsError(state) {
            state.isLoadingReviews = false;
        },

        // Форма
        setFormText(state, action) {
            state.form.text = action.payload;
        },

        // Отправка
        submitStart(state) {
            state.submitting = true;
        },
        submitEnd(state) {
            state.submitting = false;
        },
        addReview(state, action) {
            state.reviews.unshift(action.payload);
            state.form.text = "";
        },

        // Состояние отправки
        setHasSubmitted(state, action) {
            state.hasSubmitted = action.payload;
        },
        setShowPosts(state, action) {
            state.showPosts = action.payload;
        },

        // Редактирование
        startEdit(state, action) {
            state.editingReviewId = action.payload.id;
            state.editingText = action.payload.text;
        },
        cancelEdit(state) {
            state.editingReviewId = "";
            state.editingText = "";
        },
        setEditingText(state, action) {
            state.editingText = action.payload;
        },

        // Индивидуальные загрузки
        setReviewLoading(state, action) {
            state.loadingIds[action.payload] = true;
        },
        clearReviewLoading(state, action) {
            delete state.loadingIds[action.payload];
        },

        // Обновление и удаление
        updateReview(state, action) {
            const index = state.reviews.findIndex(
                (r) =>
                    (r._id || r.id) === action.payload.id
            );
            if (index !== -1) {
                state.reviews[index] = {
                    ...state.reviews[index],
                    ...action.payload,
                };
            }
        },
        removeReview(state, action) {
            state.reviews = state.reviews.filter(
                (r) => (r._id || r.id) !== action.payload
            );
        },

        // Установка всех отзывов
        setReviews(state, action) {
            state.reviews = action.payload;
        },
    },
});

export const {
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
    updateReview,
    removeReview,
    setReviews,
} = reviewsSlice.actions;

export default reviewsSlice.reducer;
