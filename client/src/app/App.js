import React, { useState, useRef } from "react";
import {
    BrowserRouter as Router,
    Routes,
    Route,
    Navigate,
} from "react-router-dom";
import { Header } from "../widgets/header";
import { AlertMessage } from "../shared/ui/alert-message";
import { createDefaultProfileData } from "../shared/model/profile";
import { DashboardPage } from "../pages/dashboard";
import { LoginPage } from "../pages/login";
import { RecommendationsPage } from "../pages/recommendations";
import { RegistrationPage } from "../pages/registration";
import { ReviewsPage } from "../pages/reviews";
import { TerminalPage } from "../pages/terminal";
import { TestsPage } from "../pages/tests";
import { TreePage } from "../pages/tree";
import "./styles/App.css";

function App() {
    // Ссылка на таймер автозакрытия уведомлений.
    const alertTimerRef = useRef(null);

    // Данные уведомления для верхней панели.
    const [statusAlert, setStatusAlert] = useState("");
    const [messageAlert, setMessageAlert] = useState("");

    // Состояние авторизации и роли для роутинга и хедера.
    const [isAuth, setIsAuth] = useState(false);
    const [isAdmin, setIsAdmin] = useState(false);

    // Списки маршрутов для проверки доступа.
    const privateRoutes = ["/account"];
    const publicRoutes = [
        "/login",
        "/register",
        "/recommendations",
        "/tests",
        "/reviews",
        "/tree",
    ];

    // Профиль пользователя, общий для нескольких страниц.
    const [profileData, setProfileData] = useState(createDefaultProfileData);

    // Показывает уведомление и очищает его по таймеру.
    const showAlert = (status, message) => {
        setStatusAlert("");
        setMessageAlert("");
        clearTimeout(alertTimerRef.current);
        setStatusAlert(status);
        setMessageAlert(message);
        alertTimerRef.current = setTimeout(() => {
            setStatusAlert("");
            setMessageAlert("");
        }, 3000);

        return true;
    };

    // Применяет простую защиту маршрутов по авторизации.
    const getRouteElement = (path, element) => {
        const isPrivateRoute = privateRoutes.includes(path);
        const isPublicRoute = publicRoutes.includes(path);

        if (isPrivateRoute && !isAuth) {
            return <Navigate to="/login" replace />;
        }

        if (!isPrivateRoute && !isPublicRoute) {
            return <Navigate to="/login" replace />;
        }

        return element;
    };

    return (
        <Router>
            <div className="App">
                <Header isAuth={isAuth} />
                {/* Верхняя панель с уведомлениями */}
                <div className="topBar">
                    <div className="topBarLeft">
                        {messageAlert && (
                            <AlertMessage
                                statusAlert={statusAlert}
                                messageAlert={messageAlert}
                            />
                        )}
                    </div>
                </div>
                {/* Таблица маршрутов приложения */}
                <Routes>
                    <Route path="/" element={<Navigate to="/login" />} />
                    <Route
                        path="/login"
                        element={getRouteElement(
                            "/login",
                            <LoginPage
                                showAlert={showAlert}
                                setIsAdmin={setIsAdmin}
                                setIsAuth={setIsAuth}
                                setProfileData={setProfileData}
                            />
                        )}
                    />
                    <Route
                        path="/register"
                        element={getRouteElement(
                            "/register",
                            <RegistrationPage showAlert={showAlert} />
                        )}
                    />
                    <Route
                        path="/account"
                        element={getRouteElement(
                            "/account",
                            <DashboardPage
                                showAlert={showAlert}
                                isAdmin={isAdmin}
                                setIsAdmin={setIsAdmin}
                                setIsAuth={setIsAuth}
                                profileData={profileData}
                                setProfileData={setProfileData}
                            />
                        )}
                    />
                    <Route
                        path="/recommendations"
                        element={getRouteElement(
                            "/recommendations",
                            <RecommendationsPage
                                showAlert={showAlert}
                                isAdmin={isAdmin}
                            />
                        )}
                    />
                    <Route
                        path="/tests"
                        element={getRouteElement(
                            "/tests",
                            <TestsPage
                                showAlert={showAlert}
                                isAdmin={isAdmin}
                                isAuth={isAuth}
                                profileData={profileData}
                            />
                        )}
                    />
                    <Route
                        path="/reviews"
                        element={getRouteElement(
                            "/reviews",
                            <ReviewsPage
                                showAlert={showAlert}
                                isAdmin={isAdmin}
                                isAuth={isAuth}
                                profileData={profileData}
                            />
                        )}
                    />
                    <Route
                        path="/terminal"
                        element={
                            <TerminalPage
                                showAlert={showAlert}
                                isAdmin={isAdmin}
                                isAuth={isAuth}
                            />
                        }
                    />
                    <Route path="/tree" element={<TreePage />} />
                </Routes>
            </div>
        </Router>
    );
}

export default App;
