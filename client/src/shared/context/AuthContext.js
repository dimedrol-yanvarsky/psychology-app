import React, { createContext, useContext, useMemo, useState } from "react";
import { createDefaultProfileData } from "../../entities/user";

// Контекст авторизации: хранит состояние сессии и профиль пользователя.
const AuthContext = createContext(null);

export const AuthProvider = ({ children }) => {
    const [isAuth, setIsAuth] = useState(false);
    const [isAdmin, setIsAdmin] = useState(false);
    const [profileData, setProfileData] = useState(createDefaultProfileData);

    const value = useMemo(
        () => ({
            isAuth,
            setIsAuth,
            isAdmin,
            setIsAdmin,
            profileData,
            setProfileData,
        }),
        [isAuth, isAdmin, profileData]
    );

    return (
        <AuthContext.Provider value={value}>{children}</AuthContext.Provider>
    );
};

export const useAuthContext = () => {
    const context = useContext(AuthContext);
    if (!context) {
        throw new Error("useAuthContext должен использоваться внутри AuthProvider");
    }
    return context;
};
