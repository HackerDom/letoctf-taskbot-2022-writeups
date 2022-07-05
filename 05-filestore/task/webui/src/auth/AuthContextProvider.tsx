import React, { useEffect, useState } from "react";
import { authContext, IAuthContext } from "./context";

export function AuthContextProvider(props: React.PropsWithChildren) {
    const [context, changeContext] = useState({
        loggedIn: false,
        userId: '',
        signin: (u: string, p: string, r: boolean) => Promise.resolve(new Response()),
        logout: () => Promise.resolve(new Response())
    });

    const logout = async() => {
        const r = await fetch('/api/logout', { method: 'POST' }); 

        changeContext({...context, loggedIn: false});
        return r;
    };

    const signin = async(username: string, pass: string, register: boolean) => {
        let url = '/api/login';
        if (register) {
            url = '/api/register';
        }

        const r = await fetch(url, {
            method: 'POST',
            body: JSON.stringify({ username, pass }),
            headers: {
                'Content-Type': 'application/json'
            }
        });
        if (!r.ok) {
            return r;
        }

        const userIdResp = await fetch('/api/userid');
        if (!userIdResp.ok) {
            return userIdResp;
        }
        const userId = (await userIdResp.json()).response;

        changeContext({ ...context, userId, loggedIn: true });
        return r; 
    };

    useEffect(() => {
        checkSigned(context, changeContext);
    }, []);

    return (
        <authContext.Provider value={{...context, signin, logout}}>
            {props.children}
        </authContext.Provider>
    );
}

async function checkSigned(context: IAuthContext, changeContext: React.Dispatch<React.SetStateAction<IAuthContext>>) {
    const r = await fetch('/api/userid');
    if (!r.ok) {
        return;
    }
    const userId = (await r.json()).response;

    changeContext({ ...context, userId, loggedIn: true })
} 