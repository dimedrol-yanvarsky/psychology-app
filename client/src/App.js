import React, { useState, useRef } from "react";
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
import AlertMessage from "./components/AlertMessage/AlertMessage";
import Terminal from "./components/Terminal/Terminal";
import TreePage from "./pages/TreePage/TreePage";

function App() {
    let alertTimerRef = useRef(null);

    const [statusAlert, setStatusAlert] = useState("");
    const [messageAlert, setMessageAlert] = useState("");

    const [isAuth, setIsAuth] = useState(false);
    const [isAdmin, setIsAdmin] = useState(false);

    const privateRoutes = ["/account"];
    const publicRoutes = [
        "/login",
        "/register",
        "/recommendations",
        "/tests",
        "/reviews",
        "/tree",
    ];

    const [profileData, setProfileData] = useState({
        id: null,
        firstName: "",
        email: "",
        status: "",
        psychoType: "",
        date: "",
        isGoogleAdded: false,
        isYandexAdded: false,
    });

    const showAlert = (status, message) => {
        setStatusAlert("");
        setMessageAlert("");
        clearTimeout(alertTimerRef);
        setStatusAlert(status);
        setMessageAlert(message);
        alertTimerRef = setTimeout(() => {
            setStatusAlert("");
            setMessageAlert("");
        }, 3000);

        return true;
    };

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
                            <Terminal
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
