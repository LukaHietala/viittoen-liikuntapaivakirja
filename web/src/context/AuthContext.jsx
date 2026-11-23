import { createContext, useContext } from "solid-js";
import { query, createAsync, revalidate } from "@solidjs/router";
import { authService } from "../services";

const AuthContext = createContext();

const getUserSession = query(async () => {
    try {
        const response = await authService.getSelf();
        if (response && response.authenticated) {
            return response.user;
        }
        return null;
    } catch (err) {
        return null;
    }
}, "userSession");

export function AuthProvider(props) {
    const user = createAsync(() => getUserSession());
    const refetch = () => revalidate(getUserSession.key);

    return <AuthContext.Provider value={{ user, refetch }}>{props.children}</AuthContext.Provider>;
}

export function useAuth() {
    return useContext(AuthContext);
}