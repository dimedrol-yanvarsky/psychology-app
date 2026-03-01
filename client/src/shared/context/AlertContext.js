import React, {
    createContext,
    useCallback,
    useContext,
    useMemo,
    useRef,
    useState,
} from "react";

// Контекст уведомлений: управляет показом и автоскрытием алертов.
const AlertContext = createContext(null);

export const AlertProvider = ({ children }) => {
    const [statusAlert, setStatusAlert] = useState("");
    const [messageAlert, setMessageAlert] = useState("");
    const alertTimerRef = useRef(null);

    // Показывает уведомление и очищает его по таймеру.
    const showAlert = useCallback((status, message) => {
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
    }, []);

    const value = useMemo(
        () => ({ statusAlert, messageAlert, showAlert }),
        [statusAlert, messageAlert, showAlert]
    );

    return (
        <AlertContext.Provider value={value}>{children}</AlertContext.Provider>
    );
};

export const useAlertContext = () => {
    const context = useContext(AlertContext);
    if (!context) {
        throw new Error(
            "useAlertContext должен использоваться внутри AlertProvider"
        );
    }
    return context;
};
