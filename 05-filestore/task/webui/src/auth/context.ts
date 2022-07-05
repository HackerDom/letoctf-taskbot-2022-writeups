import { createContext } from "react";

export interface IAuthContext {
    loggedIn: boolean;
    userId: string;
    signin: (username: string, pass: string, register: boolean) => Promise<Response>;
    logout: () => Promise<Response>;
}

export const authContext = createContext({
    loggedIn: false,
    userId: '',
    signin: () => Promise.resolve(new Response()),
    logout: () => Promise.resolve(new Response())
} as IAuthContext);