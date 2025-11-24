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

    const isAdmin = false;

    const showAlert = (status, message) => {
        setStatusAlert(status);
        setMessageAlert(message);
        setTimeout(() => {
            setStatusAlert("");
            setMessageAlert("");
        }, 3000);

        return true
    };

    return (
        <Router>
            <div className="App">
                <Header />
                <div className="topBar">
                    <div className="topBarLeft">
                        {messageAlert && (
                            <AlertMessage
                                showAlert={showAlert}
                            />
                        )}
                    </div>
                    <div className="topBarRight">
                        {isAdmin && (
                            <>
                                <button className="adminButton">
                                    Редактировать 1
                                </button>
                                <button className="adminButton">
                                    Редактировать 2
                                </button>
                                <button className="adminButton">
                                    Редактировать 3
                                </button>
                            </>
                        )}
                    </div>
                </div>
                <Routes>
                    <Route path="/" element={<Navigate to="/login" />} />
                    <Route
                        path="/login"
                        element={
                            <LoginPage
                                showAlert={showAlert}
                            />
                        }
                    />
                    <Route
                        path="/register"
                        element={
                            <RegistrationPage
                                statusAlert={statusAlert}
                                setStatusAlert={setStatusAlert}
                                messageAlert={messageAlert}
                                setMessageAlert={setMessageAlert}
                            />
                        }
                    />
                    <Route
                        path="/account"
                        element={
                            <DashboardPage
                                statusAlert={statusAlert}
                                setStatusAlert={setStatusAlert}
                                messageAlert={messageAlert}
                                setMessageAlert={setMessageAlert}
                            />
                        }
                    />
                    <Route
                        path="/recommendations"
                        element={
                            <RecommendationsPage
                                statusAlert={statusAlert}
                                setStatusAlert={setStatusAlert}
                                messageAlert={messageAlert}
                                setMessageAlert={setMessageAlert}
                            />
                        }
                    />
                    <Route
                        path="/tests"
                        element={
                            <TestsPage
                                statusAlert={statusAlert}
                                setStatusAlert={setStatusAlert}
                                messageAlert={messageAlert}
                                setMessageAlert={setMessageAlert}
                            />
                        }
                    />
                    <Route
                        path="/reviews"
                        element={
                            <ReviewsPage
                                statusAlert={statusAlert}
                                setStatusAlert={setStatusAlert}
                                messageAlert={messageAlert}
                                setMessageAlert={setMessageAlert}
                            />
                        }
                    />
                </Routes>
            </div>
        </Router>
    );
}

export default App;
