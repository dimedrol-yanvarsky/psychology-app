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

function App() {
    const [statusAlert, setStatusAlert] = useState("");
    const [messageAlert, setMessageAlert] = useState("");

    const [isAuth, setIsAuth] = useState(false);
    const user = {
        name: "Димедрол",
        surname: "Январский",
        email: "jerrystreet@example.com",
        status: "admin",
    };

    const [isAdmin, setIsAdmin] = useState(false);

    const showAlert = (status, message) => {
        setStatusAlert(status);
        setMessageAlert(message);
        setTimeout(() => {
            setStatusAlert("");
            setMessageAlert("");
        }, 3000);

        return true;
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
                    <div className="topBarRight"></div>
                </div>
                <Routes>
                    <Route path="/" element={<Navigate to="/login" />} />
                    <Route
                        path="/login"
                        element={<LoginPage showAlert={showAlert} />}
                    />
                    <Route
                        path="/register"
                        element={<RegistrationPage showAlert={showAlert} />}
                    />
                    <Route
                        path="/account"
                        element={<DashboardPage showAlert={showAlert} />}
                    />
                    <Route
                        path="/recommendations"
                        element={
                            <RecommendationsPage
                                showAlert={showAlert}
                                isAdmin={isAdmin}
                            />
                        }
                    />
                    <Route
                        path="/tests"
                        element={
                            <TestsPage
                                showAlert={showAlert}
                                isAdmin={isAdmin}
                                isAuth={isAuth}
                            />
                        }
                    />
                    <Route
                        path="/reviews"
                        element={
                            <ReviewsPage
                                showAlert={showAlert}
                                isAdmin={isAdmin}
                                isAuth={isAuth}
                            />
                        }
                    />
                </Routes>
            </div>
        </Router>
    );
}

export default App;
