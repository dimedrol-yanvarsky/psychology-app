// Начальное состояние списка отзывов, объединяющее данные, форму и состояния редактирования.
export const getInitialReviewsState = () => ({
    // Список отзывов
    reviews: [],
    isLoadingReviews: false,

    // Форма создания отзыва
    form: {
        text: "",
    },
    submitting: false,
    hasSubmitted: false,

    // Отображение постов
    showPosts: false,

    // Редактирование отзыва
    editingReviewId: "",
    editingText: "",

    // Индивидуальные загрузки отзывов
    loadingIds: {},
});

// Редьюсер списка отзывов, обрабатывающий все переходы состояний.
export const reviewsListReducer = (state, action) => {
    switch (action.type) {
        // --- Загрузка отзывов ---
        case "FETCH_REVIEWS_START":
            return { ...state, isLoadingReviews: true };
        case "FETCH_REVIEWS_SUCCESS":
            return {
                ...state,
                isLoadingReviews: false,
                reviews: action.payload,
            };
        case "FETCH_REVIEWS_ERROR":
            return { ...state, isLoadingReviews: false };

        // --- Форма создания отзыва ---
        case "SET_FORM_TEXT":
            return {
                ...state,
                form: { ...state.form, text: action.payload },
            };

        // --- Отправка отзыва ---
        case "SUBMIT_START":
            return { ...state, submitting: true };
        case "SUBMIT_END":
            return { ...state, submitting: false };
        case "ADD_REVIEW":
            return {
                ...state,
                reviews: [action.payload, ...state.reviews],
                form: { text: "" },
            };

        // --- Состояние отправки ---
        case "SET_HAS_SUBMITTED":
            return { ...state, hasSubmitted: action.payload };
        case "SET_SHOW_POSTS":
            return { ...state, showPosts: action.payload };

        // --- Редактирование отзыва ---
        case "START_EDIT":
            return {
                ...state,
                editingReviewId: action.payload.id,
                editingText: action.payload.text,
            };
        case "CANCEL_EDIT":
            return {
                ...state,
                editingReviewId: "",
                editingText: "",
            };
        case "SET_EDITING_TEXT":
            return { ...state, editingText: action.payload };

        // --- Индивидуальные загрузки ---
        case "SET_REVIEW_LOADING":
            return {
                ...state,
                loadingIds: { ...state.loadingIds, [action.payload]: true },
            };
        case "CLEAR_REVIEW_LOADING": {
            const updated = { ...state.loadingIds };
            delete updated[action.payload];
            return { ...state, loadingIds: updated };
        }

        // --- Обновление и удаление отзывов ---
        case "UPDATE_REVIEW":
            return {
                ...state,
                reviews: state.reviews.map((review) =>
                    review.id === action.payload.id
                        ? { ...review, text: action.payload.text }
                        : review
                ),
            };
        case "REMOVE_REVIEW":
            return {
                ...state,
                reviews: state.reviews.filter(
                    (review) => review.id !== action.payload
                ),
            };

        // --- Установка всех отзывов ---
        case "SET_REVIEWS":
            return { ...state, reviews: action.payload };

        default:
            return state;
    }
};
