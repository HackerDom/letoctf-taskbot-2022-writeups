import { Alert, Box, Snackbar, Typography } from "@mui/material";
import React, { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import NavBar, { Page } from "../../components/NavBar/NavBar";

export function ContentPage() {
    const { name } = useParams();
    const [file, setFile] = useState('');
    const [owner, setOwner] = useState(''); 
    enum loadCodes {
        empty = ''
    }
    const [code, changeCode] = useState(loadCodes.empty);
    
    useEffect(() => {
        loadFile(name!, setFile, setOwner, changeCode);
    }, []);

    return (
        <React.Fragment>
            <NavBar page={Page.Content} />
            <Box padding={2}>
                <Typography variant="h6">Owner id</Typography>
                {owner}
                <Typography variant="h6">Content</Typography>
                {file}
            </Box>
            <Snackbar
                open={code !== loadCodes.empty}
                autoHideDuration={4000}
                onClose={() => changeCode(loadCodes.empty)}
            >
                <Alert severity="error">{code}</Alert>
            </Snackbar>
        </React.Fragment>
    );
}

async function loadFile(
    filename: string, 
    setFile: React.Dispatch<React.SetStateAction<string>>, 
    setOwner: React.Dispatch<React.SetStateAction<string>>, 
    changeCode: React.Dispatch<React.SetStateAction<any>>
) {
    const fileResp = await fetch(`/api/get?filename=${filename}`, { method: 'GET' });
    if (fileResp.ok) {
        setFile(await fileResp.text());
    } else {
        changeCode((await fileResp.json()).response);
    }

    const ownerResp = await fetch(`/api/owner?filename=${filename}`, { method: 'GET' }); 
    if (ownerResp.ok) {
        setOwner((await ownerResp.json()).response);
    } else {
        changeCode((await ownerResp.json()).response);
    }
}
