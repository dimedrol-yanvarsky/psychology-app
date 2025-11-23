import React from "react";
import {
    BrowserRouter as Router,
    Routes,
    Route,
    Navigate,
} from "react-router-dom";
import Header from "./components/Header/Header";
import LoginPage from "./pages/LoginPage/LoginPage";
import DashboardPage from "./pages/DashboardPage/DashboardPage";
import RecommendationsPage from "./pages/RecommendationsPage/RecommendationsPage";
import TestsPage from "./pages/TestsPage/TestsPage";
import RegistrationPage from "./pages/RegistrationPage/RegistrationPage";
import "./App.css";
import ReviewsPage from "./pages/ReviewsPage/ReviewsPage";
import Terminal from "./components/Terminal/Terminal";

// Компонент для защищенных маршрутов
// const ProtectedRoute = ({ children }) => {
//     const token = localStorage.getItem("token");
//     return token ? children : <Navigate to="/login" />;
// };

// Компонент для публичных маршрутов (редирект если уже авторизован)
// const PublicRoute = ({ children }) => {
//     const token = localStorage.getItem("token");
//     return !token ? children : <Navigate to="/dashboard" />;
// };

function App() {
    return (
        <Router>
            <div className="App">
                <Header />
                <Routes>
                    <Route path="/" element={<Navigate to="/login" />} />
                    <Route
                        path="/login"
                        element={
                            // <PublicRoute>
                            <LoginPage />
                            // </PublicRoute>
                        }
                    />
                    <Route
                        path="/register"
                        element={
                            // <PublicRoute>
                            <RegistrationPage />
                            // </PublicRoute>
                        }
                    />
                    <Route
                        path="/dashboard"
                        element={
                            // <ProtectedRoute>
                            <DashboardPage />
                            // </ProtectedRoute>
                        }
                    />
                    <Route
                        path="/recommendations"
                        element={
                            // <ProtectedRoute>
                            <RecommendationsPage />
                            // </ProtectedRoute>
                        }
                    />
                    <Route
                        path="/tests"
                        element={
                            // <ProtectedRoute>
                            <TestsPage />
                            // </ProtectedRoute>
                        }
                    />
                    <Route
                        path="/reviews"
                        element={
                            // <ProtectedRoute>
                            <ReviewsPage />
                            // </ProtectedRoute>
                        }
                    />
                    <Route
                        path="/terminal"
                        element={
                            // <ProtectedRoute>
                            <Terminal />
                            // </ProtectedRoute>
                        }
                    />
                </Routes>
            </div>
        </Router>
    );
}

export default App;
