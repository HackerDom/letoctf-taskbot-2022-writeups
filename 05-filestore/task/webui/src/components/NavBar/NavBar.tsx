import React, { useContext, useState } from 'react';
import { Alert, AppBar, Button, Snackbar, Toolbar, Typography } from '@mui/material';
import { Link } from 'react-router-dom';
import { authContext } from '../../auth/context';

export enum Page {
    List = 'Files',
    Upload = 'Upload',
    Signin = 'Signin',
    Content = 'File'
}

function NavBar({ page }: { page: Page }) {
    enum logoutCodes {
        empty = ''
    }
    const [code, changeCode] = useState(logoutCodes.empty); 
    const ctx = useContext(authContext);

    const onLogout = async() => {
        const r = await ctx.logout();
        if (r.ok) {
            return;
        }

        changeCode((await r.json()).response);
    };

    return (
        <React.Fragment>
        <AppBar position="static">
            <Toolbar>
                <Typography variant="h6" component="div" sx={{ flexGrow: 1 }}>
                    {page.toString()}
                </Typography>
                {page !== Page.List &&
                    <Link to="/">
                        <Button color="primary">files</Button>
                    </Link>
                }
                {page !== Page.Upload && ctx.loggedIn &&
                    <Link to="/upload">
                        <Button color="primary">upload</Button>
                    </Link>
                }
                {page !== Page.Signin && !ctx.loggedIn &&
                    <Link to="/signin">
                        <Button color="primary">signin</Button>
                    </Link>
                }
                {ctx.loggedIn &&
                    <Button color="primary" onClick={onLogout}>logout</Button>
                }
            </Toolbar>
        </AppBar>
        <Snackbar
            open={code !== logoutCodes.empty}
            autoHideDuration={4000}
            onClose={() => changeCode(logoutCodes.empty)}
        >
            <Alert severity="error">{code}</Alert>
        </Snackbar>
        </React.Fragment>
    );
}

export default NavBar;
