import {
    BrowserRouter as Router,
    Routes,
    Route,
    Navigate,
} from "react-router-dom";
import { Provider } from "react-redux";
import { store } from "./store";
import { Header } from "../widgets/header";
import { AlertMessage } from "../shared/ui/alert-message";
import {
    AuthProvider,
    useAuthContext,
    AlertProvider,
    useAlertContext,
} from "../shared/context";
import { DashboardPage } from "../pages/dashboard";
import { LoginPage } from "../pages/login";
import { RecommendationsPage } from "../pages/recommendations";
import { RegistrationPage } from "../pages/registration";
import { ReviewsPage } from "../pages/reviews";
import { TerminalPage } from "../pages/terminal";
import { TestsPage } from "../pages/tests";
import { TreePage } from "../pages/tree";
import "./styles/App.css";

// Внутренний компонент с маршрутами, использующий контексты.
const AppRoutes = () => {
    const { isAuth } = useAuthContext();
    const { statusAlert, messageAlert } = useAlertContext();

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
        <div className="App">
            <Header />
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
                    element={getRouteElement("/login", <LoginPage />)}
                />
                <Route
                    path="/register"
                    element={getRouteElement(
                        "/register",
                        <RegistrationPage />
                    )}
                />
                <Route
                    path="/account"
                    element={getRouteElement("/account", <DashboardPage />)}
                />
                <Route
                    path="/recommendations"
                    element={getRouteElement(
                        "/recommendations",
                        <RecommendationsPage />
                    )}
                />
                <Route
                    path="/tests"
                    element={getRouteElement("/tests", <TestsPage />)}
                />
                <Route
                    path="/reviews"
                    element={getRouteElement("/reviews", <ReviewsPage />)}
                />
                <Route
                    path="/terminal"
                    element={<TerminalPage />}
                />
                <Route path="/tree" element={<TreePage />} />
            </Routes>
        </div>
    );
};

function App() {
    return (
        <Provider store={store}>
            <Router>
                <AuthProvider>
                    <AlertProvider>
                        <AppRoutes />
                    </AlertProvider>
                </AuthProvider>
            </Router>
        </Provider>
    );
}

export default App;
