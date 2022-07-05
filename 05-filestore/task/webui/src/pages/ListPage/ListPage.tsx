import React, { useEffect, useState } from 'react';
import { List, ListItem, ListItemButton, ListItemText } from "@mui/material";
import NavBar, { Page } from '../../components/NavBar/NavBar';
import { NavigateFunction, useNavigate } from 'react-router-dom';

export function ListPage() {
    const [files, setFiles] = useState<string[]>([]);
    const navigate = useNavigate();
    
    useEffect(() => {
        loadFiles().then(f => {
            setFiles(f);
        });
    }, []);

    return (
        <React.Fragment>
            <NavBar page={Page.List} />
            <List>
                { files.map((item, index) => renderItem(navigate, item, index)) }
            </List>
        </React.Fragment>
    );
}

async function loadFiles(): Promise<string[]> {
    const resp = await fetch('/api/list');
    return (await resp.json()).response;
}

function renderItem(navigate: NavigateFunction, item: string, index: number): React.ReactNode {
    return (
        <ListItem key={index} disablePadding>
            <ListItemButton onClick={() => navigate(`/file/${item}`, { replace: true })}>
                <ListItemText primary={item} />
            </ListItemButton>
        </ListItem>
    );
}
