import React from "react";
import ReactDOM from "react-dom/client";
// Глобальные стили и шрифты, подключаемые при старте приложения.
import "./app/styles/index.css";
import "./shared/assets/fonts/Yandex/stylesheet.css";
import { App } from "./app";

const root = ReactDOM.createRoot(document.getElementById("root"));
root.render(
    <React.StrictMode>
        <App />
    </React.StrictMode>
);
