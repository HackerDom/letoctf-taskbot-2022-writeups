import { Alert, Box, Button, Input, Snackbar } from "@mui/material";
import React, { useContext, useState } from "react";
import NavBar, { Page } from "../../components/NavBar/NavBar";
import { Navigate } from "react-router-dom"
import { authContext } from "../../auth/context";

export function SigninPage() {
    enum signinCodes {
        empty = ''
    }
    const ctx = useContext(authContext);
    const [code, changeCode] = useState(signinCodes.empty); 
    const [username, changeUsername] = useState('');
    const [pass, changePass] = useState('');

    if (ctx.loggedIn) {
        return <Navigate replace to="/" />;
    }

    const onSignin = async(register: boolean) => {
        const r = await ctx.signin(username, pass, register);
        if (r.ok) {
            return;
        }

        changeCode((await r.json()).response);
    };

    return (
        <React.Fragment>
            <NavBar page={Page.Signin} />
            <Box
                display="flex"
                justifyContent="center"
                alignItems="center"
                height="100%"
                flexDirection="column"
                gap={2}
            >
                <Input placeholder="Username" value={username} onChange={(e: any) => changeUsername(e.target.value)} />
                <Input type="password" placeholder="Password" value={pass} onChange={(e: any) => changePass(e.target.value)} />
                <Button onClick={() => onSignin(false)}>Login</Button>
                <Button onClick={() => onSignin(true)}>Register</Button>
            </Box>
            <Snackbar
                open={code !== signinCodes.empty}
                autoHideDuration={4000}
                onClose={() => changeCode(signinCodes.empty)}
            >
                <Alert severity="error">{code}</Alert>
            </Snackbar>
        </React.Fragment>
    );
}
