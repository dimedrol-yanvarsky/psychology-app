import { configureStore } from "@reduxjs/toolkit";
import { authReducer, alertReducer } from "../../shared/store";
import dashboardReducer from "../../features/dashboard/model/dashboardSlice";
import testsReducer from "../../features/tests/model/testsSlice";
import reviewsReducer from "../../features/reviews/model/reviewsSlice";
import recommendationsReducer from "../../features/recommendations/model/recommendationsSlice";

// Конфигурация Redux store с объединением всех слайсов.
export const store = configureStore({
    reducer: {
        auth: authReducer,
        alert: alertReducer,
        dashboard: dashboardReducer,
        tests: testsReducer,
        reviews: reviewsReducer,
        recommendations: recommendationsReducer,
    },
});
